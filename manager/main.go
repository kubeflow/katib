package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"strconv"
	"time"

	pb "github.com/kubeflow/hp-tuning/api"
	kdb "github.com/kubeflow/hp-tuning/db"
	"github.com/kubeflow/hp-tuning/manager/modelstore"
	tbif "github.com/kubeflow/hp-tuning/manager/visualise/tensorboard"
	"github.com/kubeflow/hp-tuning/manager/worker_interface"
	dlkwif "github.com/kubeflow/hp-tuning/manager/worker_interface/dlk"
	k8swif "github.com/kubeflow/hp-tuning/manager/worker_interface/kubernetes"
	nvdwif "github.com/kubeflow/hp-tuning/manager/worker_interface/nvdocker"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	k8s_namespace            = "katib"
	port                     = "0.0.0.0:6789"
	defaultEarlyStopInterval = 60
	defaultStoreInterval     = 30
)

var init_db = flag.Bool("init", false, "Initialize DB")
var worker = flag.String("w", "kubernetes", "Worker Typw")
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

func (s *server) trialIteration(conf *pb.StudyConfig, study_id string, sCh studyCh) error {
	defer delete(s.StudyChList, study_id)
	defer s.wIF.CleanWorkers(study_id)
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
	strtm := time.NewTimer(defaultStoreInterval * time.Second)
	log.Printf("Study %v start.", study_id)
	log.Printf("Study conf %v", conf)
	for {
		select {
		case <-tm.C:
			if conf.SuggestAlgorithm != "" {
				err := s.wIF.CheckRunningTrials(study_id, conf.ObjectiveValueName)
				if err != nil {
					return err
				}
				r, err := s.SuggestTrials(context.Background(), &pb.SuggestTrialsRequest{StudyId: study_id, SuggestAlgorithm: conf.SuggestAlgorithm, Configs: conf})
				if err != nil {
					log.Printf("SuggestTrials failed %v", err)
					return err
				}
				if r.Completed {
					log.Printf("Study %v completed.", study_id)
					return nil
				} else if len(r.Trials) > 0 {
					for _, trial := range r.Trials {
						trial.Status = pb.TrialState_PENDING
						trial.StudyId = study_id
						err = dbIf.CreateTrial(trial)
						if err != nil {
							log.Printf("CreateTrial failed %v", err)
							return err
						}
					}
					err = s.wIF.SpawnWorkers(r.Trials, study_id)
					if err != nil {
						log.Printf("SpawnWorkers failed %v", err)
						return err
					}
					for _, t := range r.Trials {
						err = tbif.SpawnTensorBoard(study_id, t.TrialId, k8s_namespace, conf.Mount)
						if err != nil {
							log.Printf("SpawnTB failed %v", err)
							return err
						}
					}
				}
				tm.Reset(1 * time.Second)
			}
		case <-strtm.C:
			ret, err := s.GetStoredModels(context.Background(), &pb.GetStoredModelsRequest{StudyName: conf.Name})
			if err != nil {
				log.Printf("GetStoredModels Err %v", err)
			}
			ts, err := dbIf.GetTrialList(study_id)
			if err != nil {
				log.Printf("GetTrials Err %v", err)
			}
			for _, t := range ts {
				tid := t.TrialId
				tst, err := dbIf.GetTrialStatus(tid)
				if err != nil {
					log.Printf("GetTrialStatus Err %v", err)
					continue
				}
				if tst == pb.TrialState_COMPLETED {
					isin := false
					for _, m := range ret.Models {
						if m.TrialId == tid {
							isin = true
						}
					}
					met := make([]*pb.Metrics, len(conf.Metrics))
					for i, mn := range conf.Metrics {
						l, _ := dbIf.GetTrialLogs(tid, &kdb.GetTrialLogOpts{Name: mn})
						met[i] = &pb.Metrics{Name: mn, Value: l[len(l)-1].Value}
					}
					if !isin {
						t, _ := dbIf.GetTrial(tid)
						s.StoreModel(context.Background(), &pb.StoreModelRequest{
							Model: &pb.ModelInfo{
								StudyName:  conf.Name,
								TrialId:    tid,
								Parameters: t.ParameterSet,
								Metrics:    met,
							},
						})

					}
				}
			}
			strtm.Reset(defaultStoreInterval * time.Second)

		case <-estm.C:
			ret, err := s.EarlyStopping(context.Background(), &pb.EarlyStoppingRequest{StudyId: study_id, EarlyStoppingAlgorithm: conf.EarlyStoppingAlgorithm})
			if err != nil {
				log.Printf("Early Stopping Error: %v", err)
			} else {
				if len(ret.Trials) > 0 {
					for _, t := range ret.Trials {
						s.CompleteTrial(context.Background(), &pb.CompleteTrialRequest{StudyId: study_id, TrialId: t.TrialId, IsComplete: false})
					}
				}
			}
			estm.Reset(time.Duration(ei) * time.Second)
		case <-sCh.stopCh:
			log.Printf("Study %v is stopped.", study_id)
			for _, t := range s.wIF.GetRunningTrials(study_id) {
				t.Status = pb.TrialState_KILLED
			}
			return nil
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

	study_id, err := dbIf.CreateStudy(in.StudyConfig)
	if in.StudyConfig.SuggestAlgorithm != "" {
		_, err = s.InitializeSuggestService(
			ctx,
			&pb.InitializeSuggestServiceRequest{
				StudyId:              study_id,
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
				StudyId:                 study_id,
				EarlyStoppingParameters: in.StudyConfig.EarlyStoppingParameters,
			},
		)
		if err != nil {
			return &pb.CreateStudyReply{}, err
		}
	}
	sCh := studyCh{stopCh: make(chan bool), addMetricsCh: make(chan string)}
	_, err = s.StoreStudy(ctx, &pb.StoreStudyRequest{StudyName: in.StudyConfig.Name})
	if err != nil {
		return &pb.CreateStudyReply{}, err
	}
	go s.trialIteration(in.StudyConfig, study_id, sCh)
	s.StudyChList[study_id] = sCh
	return &pb.CreateStudyReply{StudyId: study_id}, nil
}

func (s *server) StopStudy(ctx context.Context, in *pb.StopStudyRequest) (*pb.StopStudyReply, error) {
	sc, ok := s.StudyChList[in.StudyId]
	if !ok {
		return &pb.StopStudyReply{}, errors.New("Study Id not found")
	}
	sc.stopCh <- false
	return &pb.StopStudyReply{}, nil
}

func spawn_worker(study_task string, params string) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	template, err := clientset.CoreV1().PodTemplates(k8s_namespace).Get(study_task, metav1.GetOptions{})
	if err != nil {
		return err
	}

	_, err = clientset.BatchV1().Jobs(k8s_namespace).Create(&batchv1.Job{
		Spec: batchv1.JobSpec{
			Template: template.Template,
		},
	})

	// TODO: Update worker status
	return err
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
	var suggest_algo string

	// TODO: only a few columns are needed but GetStudyConfig does a full retrieval
	study, err := dbIf.GetStudyConfig(in.StudyId)
	if err != nil {
		return nil, err
	}

	if in.SuggestAlgorithm != "" {
		suggest_algo = in.SuggestAlgorithm
	} else if study.SuggestAlgorithm != "" {
		suggest_algo = study.SuggestAlgorithm
	} else {
		return &pb.SuggestTrialsReply{Completed: false}, errors.New("No suggest algorithm specified")
	}

	conn, err := grpc.Dial("vizier-suggestion-"+suggest_algo+":6789", grpc.WithInsecure())
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

func (s *server) StoreStudy(ctx context.Context, in *pb.StoreStudyRequest) (*pb.StoreStudyReply, error) {
	err := s.msIf.StoreStudy(in)
	return &pb.StoreStudyReply{}, err
}

func (s *server) StoreModel(ctx context.Context, in *pb.StoreModelRequest) (*pb.StoreModelReply, error) {
	err := s.msIf.StoreModel(in)
	return &pb.StoreModelReply{}, err
}

func (s *server) GetStoredStudies(ctx context.Context, in *pb.GetStoredStudiesRequest) (*pb.GetStoredStudiesReply, error) {
	ret, err := s.msIf.GetStoredStudies()
	return &pb.GetStoredStudiesReply{Studies: ret}, err
}

func (s *server) GetStoredModels(ctx context.Context, in *pb.GetStoredModelsRequest) (*pb.GetStoredModelsReply, error) {
	ret, err := s.msIf.GetStoredModels(in)
	return &pb.GetStoredModelsReply{Models: ret}, err
}

func (s *server) GetStoredModel(ctx context.Context, in *pb.GetStoredModelRequest) (*pb.GetStoredModelReply, error) {
	ret, err := s.msIf.GetStoredModel(in)
	return &pb.GetStoredModelReply{Model: ret}, err
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
	switch *worker {
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
		pb.RegisterManagerServer(s, &server{wIF: dlkwif.NewDlkWorkerInterface("http://dlk-manager:1323", k8s_namespace), msIf: modelstore.NewModelDB("modeldb-backend", "6543"), StudyChList: make(map[string]studyCh)})
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
