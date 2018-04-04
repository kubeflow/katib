package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"github.com/mlkube/katib/manager/worker_interface"
	dlkwif "github.com/mlkube/katib/manager/worker_interface/dlk"
	k8swif "github.com/mlkube/katib/manager/worker_interface/kubernetes"
	nvdwif "github.com/mlkube/katib/manager/worker_interface/nvdocker"

	tbif "github.com/mlkube/katib/manager/visualise/tensorboard"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	vdb "github.com/mlkube/katib/db"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	pb "github.com/mlkube/katib/api"
)

const (
	k8s_namespace = "katib"
	port          = "0.0.0.0:6789"
)

var init_db = flag.Bool("init", false, "Initialize DB")
var worker = flag.String("w", "kubernetes", "Worker Typw")
var dbIf vdb.VizierDBInterface

type studyCh struct {
	stopCh       chan bool
	addMetricsCh chan string
}
type server struct {
	wIF         worker_interface.WorkerInterface
	StudyChList map[string]studyCh
}

func (s *server) saveResult(study_id string) error {
	var result string
	c := s.wIF.GetCompletedTrials(study_id)
	if len(c) == 0 {
		log.Printf("Study %v has no completed Trials", study_id)
		return nil
	}
	result += "TrialID"
	for _, p := range c[0].ParameterSet {
		result += fmt.Sprintf("\t%v", p.Name)
	}
	result += "\tObjectiveValue"
	if len(c) > 0 {
		if len(c[0].EvalLogs) > 0 {
			for _, m := range c[0].EvalLogs[len(c[0].EvalLogs)-1].Metrics {
				result += fmt.Sprintf("\t%v", m.Name)
			}
		}
	}
	result += "\tTiem_cost"
	result += "\n"

	for _, ct := range c {
		result += fmt.Sprintf("%v", ct.TrialId)
		for _, p := range ct.ParameterSet {
			result += fmt.Sprintf("\t%v", p.Value)
		}
		result += fmt.Sprintf("\t%v", ct.ObjectiveValue)
		for _, m := range ct.EvalLogs[len(ct.EvalLogs)-1].Metrics {
			result += fmt.Sprintf("\t%v", m.Value)
		}
		st, _ := time.Parse(time.RFC3339, ct.EvalLogs[0].Time)
		et, _ := time.Parse(time.RFC3339, ct.EvalLogs[len(ct.EvalLogs)-1].Time)
		result += fmt.Sprintf("\t%v", et.Sub(st))
		result += "\n"
	}
	ioutil.WriteFile("/conf/result", []byte(result), os.ModePerm)
	return nil
}

func (s *server) trialIteration(conf *pb.StudyConfig, study_id string, sCh studyCh) error {
	defer delete(s.StudyChList, study_id)
	defer s.wIF.CleanWorkers(study_id)
	tm := time.NewTimer(1 * time.Second)
	log.Printf("Study %v start.", study_id)
	log.Printf("Study conf %v", conf)
	for {
		select {
		case <-tm.C:
			err := s.wIF.CheckRunningTrials(study_id, conf.ObjectiveValueName, conf.Metrics)
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
				//s.saveResult(study_id)
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
	if in.StudyConfig.ObjectiveValueName == "" {
		return &pb.CreateStudyReply{}, errors.New("Objective_Value_Name is required.")
	}

	study_id, err := dbIf.CreateStudy(in.StudyConfig)

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
	sCh := studyCh{stopCh: make(chan bool), addMetricsCh: make(chan string)}
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

func (s *server) GetStudys(ctx context.Context, in *pb.GetStudysRequest) (*pb.GetStudysReply, error) {
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
	return &pb.GetStudysReply{StudyInfos: ss}, nil
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

func (s *server) CompleteTrial(context.Context, *pb.CompleteTrialRequest) (*pb.CompleteTrialReply, error) {
	return nil, errors.New("not implemented")
}
func (s *server) ShouldTrialStop(context.Context, *pb.ShouldTrialStopRequest) (*pb.ShouldTrialStopReply, error) {
	return nil, errors.New("not implemented")
}
func (s *server) GetObjectValue(context.Context, *pb.GetObjectValueRequest) (*pb.GetObjectValueReply, error) {
	return nil, errors.New("not implemented")
}

func (s *server) AddMeasurementToTrials(context.Context, *pb.AddMeasurementToTrialsRequest) (*pb.AddMeasurementToTrialsReply, error) {

	return &pb.AddMeasurementToTrialsReply{}, nil
}

func main() {
	flag.Parse()

	//	if *init_db {
	var err error
	dbIf = vdb.New()

	dbIf.DB_Init()
	//	} else {
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
		pb.RegisterManagerServer(s, &server{wIF: k8swif.NewKubernetesWorkerInterface(clientset, dbIf), StudyChList: make(map[string]studyCh)})
		// XXX Is this useful?
	case "dlk":
		log.Printf("Worker: dlk\n")
		pb.RegisterManagerServer(s, &server{wIF: dlkwif.NewDlkWorkerInterface("http://dlk-manager:1323", k8s_namespace), StudyChList: make(map[string]studyCh)})
	case "nv-docker":
		log.Printf("Worker: nv-docker\n")
		pb.RegisterManagerServer(s, &server{wIF: nvdwif.NewNvDockerWorkerInterface(), StudyChList: make(map[string]studyCh)})
	default:
		log.Fatalf("Unknown worker")
	}
	reflection.Register(s)
	if err = s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	//	}
}
