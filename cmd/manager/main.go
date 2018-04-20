package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"strconv"
	"time"

	pb "github.com/kubeflow/katib/api"
	kdb "github.com/kubeflow/katib/pkg/db"
	"github.com/kubeflow/katib/pkg/manager/modelstore"
	tbif "github.com/kubeflow/katib/pkg/manager/visualise/tensorboard"
	"github.com/kubeflow/katib/pkg/manager/worker_interface"
	dlkwif "github.com/kubeflow/katib/pkg/manager/worker_interface/dlk"
	k8swif "github.com/kubeflow/katib/pkg/manager/worker_interface/kubernetes"
	nvdwif "github.com/kubeflow/katib/pkg/manager/worker_interface/nvdocker"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	namespace                = "katib"
	port                     = "0.0.0.0:6789"
	defaultEarlyStopInterval = 60
	defaultSaveInterval      = 30
)

var workerType = flag.String("w", "kubernetes", "Worker Type")
var ingressHost = flag.String("i", "kube-cluster.example.net", "Ingress host for TensorBoard visualize")
var dbIf kdb.VizierDBInterface

type studyCh struct {
	stopCh       chan bool
	addMetricsCh chan string
}
type server struct {
	wIF         worker_interface.WorkerInterface
	msIf        modelstore.ModelStore
	StudyChList map[string]studyCh
}

func (s *server) saveCompletedModels(studyID string, conf *pb.StudyConfig) error {
	ret, err := s.GetSavedModels(context.Background(), &pb.GetSavedModelsRequest{StudyName: conf.Name})
	if err != nil {
		log.Printf("GetSavedModels Err %v", err)
		return err
	}
	ts, err := dbIf.GetTrialList(studyID)
	if err != nil {
		log.Printf("GetTrials Err %v", err)
		return err
	}
	for _, t := range ts {
		tid := t.TrialId
		tst, err := dbIf.GetTrialStatus(tid)
		if err != nil {
			log.Printf("GetTrialStatus Err %v", err)
			continue
		}
		isin := false
		if tst == pb.TrialState_COMPLETED {
			for _, m := range ret.Models {
				if m.TrialId == tid {
					isin = true
					break
				}
			}
			if !isin {
				met := make([]*pb.Metrics, len(conf.Metrics))
				for i, mn := range conf.Metrics {
					l, _ := dbIf.GetTrialLogs(tid, &kdb.GetTrialLogOpts{Name: mn})
					if len(l) > 0 {
						met[i] = &pb.Metrics{Name: mn, Value: l[len(l)-1].Value}
					}
				}
				t, _ := dbIf.GetTrial(tid)
				s.SaveModel(context.Background(), &pb.SaveModelRequest{
					Model: &pb.ModelInfo{
						StudyName:  conf.Name,
						TrialId:    tid,
						Parameters: t.ParameterSet,
						Metrics:    met,
					},
				})
				log.Printf("Trial %v in Study %v is saved", tid, conf.Name)
			}
		}
	}
	return nil
}

func (s *server) trialIteration(conf *pb.StudyConfig, studyID string, sCh studyCh) error {
	defer delete(s.StudyChList, studyID)
	defer s.wIF.CleanWorkers(studyID)
	tm := time.NewTimer(1 * time.Second)
	ei := 0
	var err error
	for _, ec := range conf.EarlyStoppingParameters {
		if ec.Name == "CheckInterval" {
			ei, err = strconv.Atoi(ec.Value)
			if err != nil {
				ei = 0
			}
		}
	}
	if ei == 0 {
		ei = defaultEarlyStopInterval
	}
	estm := time.NewTimer(time.Duration(ei) * time.Second)
	strtm := time.NewTimer(defaultSaveInterval * time.Second)
	log.Printf("Study %v start.", studyID)
	log.Printf("Study conf %v", conf)
	for {
		select {
		case <-tm.C:
			if conf.SuggestAlgorithm != "" {
				err := s.wIF.CheckRunningTrials(studyID, conf.ObjectiveValueName)
				if err != nil {
					return err
				}
				r, err := s.SuggestTrials(context.Background(), &pb.SuggestTrialsRequest{StudyId: studyID, SuggestAlgorithm: conf.SuggestAlgorithm, Configs: conf})
				if err != nil {
					log.Printf("SuggestTrials failed %v", err)
					return err
				}
				if r.Completed {
					log.Printf("Study %v completed.", studyID)
					return s.saveCompletedModels(studyID, conf)
				} else if len(r.Trials) > 0 {
					for _, trial := range r.Trials {
						trial.Status = pb.TrialState_PENDING
						trial.StudyId = studyID
						err = dbIf.CreateTrial(trial)
						if err != nil {
							log.Printf("CreateTrial failed %v", err)
							return err
						}
					}
					err = s.wIF.SpawnWorkers(r.Trials, studyID)
					if err != nil {
						log.Printf("SpawnWorkers failed %v", err)
						return err
					}
					for _, t := range r.Trials {
						err = tbif.SpawnTensorBoard(studyID, t.TrialId, conf.Name, namespace, conf.Mount, ingressHost)
						if err != nil {
							log.Printf("SpawnTB failed %v", err)
							return err
						}
					}
				}
				tm.Reset(1 * time.Second)
			}
		case <-strtm.C:
			s.saveCompletedModels(studyID, conf)
			strtm.Reset(defaultSaveInterval * time.Second)

		case <-estm.C:
			ret, err := s.EarlyStopping(context.Background(), &pb.EarlyStoppingRequest{StudyId: studyID, EarlyStoppingAlgorithm: conf.EarlyStoppingAlgorithm})
			if err != nil {
				log.Printf("Early Stopping Error: %v", err)
			} else {
				if len(ret.Trials) > 0 {
					for _, t := range ret.Trials {
						s.CompleteTrial(context.Background(), &pb.CompleteTrialRequest{StudyId: studyID, TrialId: t.TrialId, IsComplete: false})
					}
				}
			}
			estm.Reset(time.Duration(ei) * time.Second)
		case <-sCh.stopCh:
			log.Printf("Study %v is stopped.", studyID)
			for _, t := range s.wIF.GetRunningTrials(studyID) {
				t.Status = pb.TrialState_KILLED
			}
			return s.saveCompletedModels(studyID, conf)
		case m := <-sCh.addMetricsCh:
			conf.Metrics = append(conf.Metrics, m)
		}
	}
	return nil
}

func (s *server) CreateStudy(ctx context.Context, in *pb.CreateStudyRequest) (*pb.CreateStudyReply, error) {
	if in == nil || in.StudyConfig == nil {
		return &pb.CreateStudyReply{}, errors.New("StudyConfig is missing.")
	}

	if in.StudyConfig.ObjectiveValueName == "" {
		return &pb.CreateStudyReply{}, errors.New("Objective_Value_Name is required.")
	}

	studyID, err := dbIf.CreateStudy(in.StudyConfig)
	if in.StudyConfig.SuggestAlgorithm != "" {
		_, err = s.InitializeSuggestService(
			ctx,
			&pb.InitializeSuggestServiceRequest{
				StudyId:              studyID,
				SuggestAlgorithm:     in.StudyConfig.SuggestAlgorithm,
				SuggestionParameters: in.StudyConfig.SuggestionParameters,
				Configs:              in.StudyConfig,
			},
		)
		if err != nil {
			return &pb.CreateStudyReply{}, err
		}
	} else {
		log.Printf("Suggestion Algorithm is not set.")
	}

	if in.StudyConfig.EarlyStoppingAlgorithm != "" && in.StudyConfig.EarlyStoppingAlgorithm != "none" {
		conn, err := grpc.Dial("vizier-earlystopping-"+in.StudyConfig.EarlyStoppingAlgorithm+":6789", grpc.WithInsecure())
		if err != nil {
			return &pb.CreateStudyReply{}, err
		}
		defer conn.Close()
		c := pb.NewEarlyStoppingClient(conn)
		_, err = c.SetEarlyStoppingParameter(
			context.Background(),
			&pb.SetEarlyStoppingParameterRequest{
				StudyId:                 studyID,
				EarlyStoppingParameters: in.StudyConfig.EarlyStoppingParameters,
			},
		)
		if err != nil {
			return &pb.CreateStudyReply{}, err
		}
	}
	sCh := studyCh{stopCh: make(chan bool), addMetricsCh: make(chan string)}
	_, err = s.SaveStudy(ctx, &pb.SaveStudyRequest{StudyName: in.StudyConfig.Name})
	if err != nil {
		return &pb.CreateStudyReply{}, err
	}
	go s.trialIteration(in.StudyConfig, studyID, sCh)
	s.StudyChList[studyID] = sCh
	return &pb.CreateStudyReply{StudyId: studyID}, nil
}

func (s *server) StopStudy(ctx context.Context, in *pb.StopStudyRequest) (*pb.StopStudyReply, error) {
	sc, ok := s.StudyChList[in.StudyId]
	if !ok {
		return &pb.StopStudyReply{}, errors.New("Study Id not found")
	}
	sc.stopCh <- false
	return &pb.StopStudyReply{}, nil
}

func (s *server) GetStudies(ctx context.Context, in *pb.GetStudiesRequest) (*pb.GetStudiesReply, error) {
	ss := make([]*pb.StudyInfo, len(s.StudyChList))
	i := 0
	for sid := range s.StudyChList {
		sc, _ := dbIf.GetStudyConfig(sid)
		ss[i] = &pb.StudyInfo{
			StudyId:           sid,
			Name:              sc.Name,
			Owner:             sc.Owner,
			RunningTrialNum:   int32(len(s.wIF.GetRunningTrials(sid))),
			CompletedTrialNum: int32(len(s.wIF.GetCompletedTrials(sid))),
		}
		i++
	}
	return &pb.GetStudiesReply{StudyInfos: ss}, nil
}

func (s *server) InitializeSuggestService(ctx context.Context, in *pb.InitializeSuggestServiceRequest) (*pb.InitializeSuggestServiceReply, error) {
	conn, err := grpc.Dial("vizier-suggestion-"+in.SuggestAlgorithm+":6789", grpc.WithInsecure())
	if err != nil {
		log.Printf("could not connect: %v", err)
		return &pb.InitializeSuggestServiceReply{}, err
	}
	defer conn.Close()
	c := pb.NewSuggestionClient(conn)
	req := &pb.SetSuggestionParametersRequest{StudyId: in.StudyId, SuggestionParameters: in.SuggestionParameters, Configs: in.Configs}
	_, err = c.SetSuggestionParameters(context.Background(), req)
	if err != nil {
		log.Printf("Set Suggestion Parameter failed: %v", err)
	}
	return &pb.InitializeSuggestServiceReply{}, err
}

func (s *server) SuggestTrials(ctx context.Context, in *pb.SuggestTrialsRequest) (*pb.SuggestTrialsReply, error) {
	var suggestAlgorithm string

	// TODO: only a few columns are needed but GetStudyConfig does a full retrieval
	study, err := dbIf.GetStudyConfig(in.StudyId)
	if err != nil {
		return nil, err
	}

	if in.SuggestAlgorithm != "" {
		suggestAlgorithm = in.SuggestAlgorithm
	} else if study.SuggestAlgorithm != "" {
		suggestAlgorithm = study.SuggestAlgorithm
	} else {
		return &pb.SuggestTrialsReply{Completed: false}, errors.New("No suggest algorithm specified")
	}

	conn, err := grpc.Dial("vizier-suggestion-"+suggestAlgorithm+":6789", grpc.WithInsecure())
	if err != nil {
		return &pb.SuggestTrialsReply{Completed: false}, err
	}

	defer conn.Close()
	c := pb.NewSuggestionClient(conn)
	cts := s.wIF.GetCompletedTrials(in.StudyId)
	rts := s.wIF.GetRunningTrials(in.StudyId)
	req := &pb.GenerateTrialsRequest{StudyId: in.StudyId, Configs: in.Configs, CompletedTrials: cts, RunningTrials: rts}
	r, err := c.GenerateTrials(context.Background(), req)
	if err != nil {
		return &pb.SuggestTrialsReply{Completed: false}, err
	}

	// TODO: do async
	return &pb.SuggestTrialsReply{Trials: r.Trials, Completed: r.Completed}, nil
}

func (s *server) CompleteTrial(ctx context.Context, in *pb.CompleteTrialRequest) (*pb.CompleteTrialReply, error) {
	err := s.wIF.CompleteTrial(in.StudyId, in.TrialId, in.IsComplete)
	return &pb.CompleteTrialReply{}, err
}

func (s *server) EarlyStopping(ctx context.Context, in *pb.EarlyStoppingRequest) (*pb.EarlyStoppingReply, error) {
	if in.EarlyStoppingAlgorithm == "" || in.EarlyStoppingAlgorithm == "none" {
		return &pb.EarlyStoppingReply{}, nil
	}
	conn, err := grpc.Dial("vizier-earlystopping-"+in.EarlyStoppingAlgorithm+":6789", grpc.WithInsecure())
	if err != nil {
		return &pb.EarlyStoppingReply{}, err
	}

	defer conn.Close()
	c := pb.NewEarlyStoppingClient(conn)
	req := &pb.ShouldTrialStopRequest{StudyId: in.StudyId}
	r, err := c.ShouldTrialStop(context.Background(), req)
	if err != nil {
		return &pb.EarlyStoppingReply{}, err
	}

	// TODO: do async
	return &pb.EarlyStoppingReply{Trials: r.Trials}, nil
}

func (s *server) GetObjectValue(context.Context, *pb.GetObjectValueRequest) (*pb.GetObjectValueReply, error) {
	return nil, errors.New("not implemented")
}

func (s *server) AddMeasurementToTrials(context.Context, *pb.AddMeasurementToTrialsRequest) (*pb.AddMeasurementToTrialsReply, error) {

	return &pb.AddMeasurementToTrialsReply{}, nil
}

func (s *server) SaveStudy(ctx context.Context, in *pb.SaveStudyRequest) (*pb.SaveStudyReply, error) {
	err := s.msIf.SaveStudy(in)
	return &pb.SaveStudyReply{}, err
}

func (s *server) SaveModel(ctx context.Context, in *pb.SaveModelRequest) (*pb.SaveModelReply, error) {
	err := s.msIf.SaveModel(in)
	return &pb.SaveModelReply{}, err
}

func (s *server) GetSavedStudies(ctx context.Context, in *pb.GetSavedStudiesRequest) (*pb.GetSavedStudiesReply, error) {
	ret, err := s.msIf.GetSavedStudies()
	return &pb.GetSavedStudiesReply{Studies: ret}, err
}

func (s *server) GetSavedModels(ctx context.Context, in *pb.GetSavedModelsRequest) (*pb.GetSavedModelsReply, error) {
	ret, err := s.msIf.GetSavedModels(in)
	return &pb.GetSavedModelsReply{Models: ret}, err
}

func (s *server) GetSavedModel(ctx context.Context, in *pb.GetSavedModelRequest) (*pb.GetSavedModelReply, error) {
	ret, err := s.msIf.GetSavedModel(in)
	return &pb.GetSavedModelReply{Model: ret}, err
}

func main() {
	flag.Parse()
	var err error
	dbIf = kdb.New()
	dbIf.DB_Init()
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	size := 1<<31 - 1
	s := grpc.NewServer(grpc.MaxRecvMsgSize(size), grpc.MaxSendMsgSize(size))
	switch *workerType {
	case "kubernetes":
		log.Printf("Worker: kubernetes\n")
		kc, err := clientcmd.BuildConfigFromFlags("", "/conf/kubeconfig")
		if err != nil {
			log.Fatal(err)
		}
		clientset, err := kubernetes.NewForConfig(kc)
		if err != nil {
			log.Fatal(err)
		}
		pb.RegisterManagerServer(s, &server{wIF: k8swif.NewKubernetesWorkerInterface(clientset, dbIf), msIf: modelstore.NewModelDB("modeldb-backend", "6543"), StudyChList: make(map[string]studyCh)})
	case "dlk":
		log.Printf("Worker: dlk\n")
		pb.RegisterManagerServer(s, &server{wIF: dlkwif.NewDlkWorkerInterface("http://dlk-manager:1323", namespace), msIf: modelstore.NewModelDB("modeldb-backend", "6543"), StudyChList: make(map[string]studyCh)})
	case "nv-docker":
		log.Printf("Worker: nv-docker\n")
		pb.RegisterManagerServer(s, &server{wIF: nvdwif.NewNvDockerWorkerInterface(), msIf: modelstore.NewModelDB("modeldb-backend", "6543"), StudyChList: make(map[string]studyCh)})
	default:
		log.Fatalf("Unknown worker")
	}
	reflection.Register(s)
	if err = s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
