package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"strings"
	"time"

	pb "github.com/kubeflow/katib/pkg/api"
	kdb "github.com/kubeflow/katib/pkg/db"
	"github.com/kubeflow/katib/pkg/manager/modelstore"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = "0.0.0.0:6789"
)

var ingressHost = flag.String("i", "kube-cluster.example.net", "Ingress host")
var dbIf kdb.VizierDBInterface

type server struct {
	msIf modelstore.ModelStore
}

func (s *server) CreateStudy(ctx context.Context, in *pb.CreateStudyRequest) (*pb.CreateStudyReply, error) {
	if in == nil || in.StudyConfig == nil {
		return &pb.CreateStudyReply{}, errors.New("StudyConfig is missing.")
	}
	studyID, err := dbIf.CreateStudy(in.StudyConfig)
	if err != nil {
		return &pb.CreateStudyReply{}, err
	}
	s.SaveStudy(ctx, &pb.SaveStudyRequest{
		StudyName:   in.StudyConfig.Name,
		Owner:       in.StudyConfig.Owner,
		Description: "StudyID: " + studyID,
	})
	return &pb.CreateStudyReply{StudyId: studyID}, nil
}

func (s *server) GetStudy(ctx context.Context, in *pb.GetStudyRequest) (*pb.GetStudyReply, error) {
	sc, err := dbIf.GetStudyConfig(in.StudyId)
	return &pb.GetStudyReply{StudyConfig: sc}, err
}

func (s *server) GetStudyList(ctx context.Context, in *pb.GetStudyListRequest) (*pb.GetStudyListReply, error) {
	sl, err := dbIf.GetStudyList()
	if err != nil {
		return &pb.GetStudyListReply{}, err
	}
	result := make([]*pb.StudyOverview, len(sl))
	for i, id := range sl {
		sc, err := dbIf.GetStudyConfig(id)
		if err != nil {
			return &pb.GetStudyListReply{}, err
		}
		result[i] = &pb.StudyOverview{
			Name:  sc.Name,
			Owner: sc.Owner,
			Id:    id,
		}
	}
	return &pb.GetStudyListReply{StudyOverviews: result}, err
}

func (s *server) CreateTrial(ctx context.Context, in *pb.CreateTrialRequest) (*pb.CreateTrialReply, error) {
	err := dbIf.CreateTrial(in.Trial)
	return &pb.CreateTrialReply{TrialId: in.Trial.TrialId}, err
}

func (s *server) GetTrials(ctx context.Context, in *pb.GetTrialsRequest) (*pb.GetTrialsReply, error) {
	tl, err := dbIf.GetTrialList(in.StudyId)
	return &pb.GetTrialsReply{Trials: tl}, err
}

func (s *server) GetSuggestions(ctx context.Context, in *pb.GetSuggestionsRequest) (*pb.GetSuggestionsReply, error) {
	if in.SuggestionAlgorithm == "" {
		return &pb.GetSuggestionsReply{Trials: []*pb.Trial{}}, errors.New("No suggest algorithm specified")
	}
	conn, err := grpc.Dial("vizier-suggestion-"+in.SuggestionAlgorithm+":6789", grpc.WithInsecure())
	if err != nil {
		return &pb.GetSuggestionsReply{Trials: []*pb.Trial{}}, err
	}

	defer conn.Close()
	c := pb.NewSuggestionClient(conn)
	r, err := c.GetSuggestions(ctx, in)
	if err != nil {
		return &pb.GetSuggestionsReply{Trials: []*pb.Trial{}}, err
	}
	return r, nil
}

func (s *server) CreateWorker(ctx context.Context, in *pb.CreateWorkerReauest) (*pb.CreateWorkerReply, error) {
	wid, err := dbIf.CreateWorker(in.Worker)
	return &pb.CreateWorkerReply{WorkerId: wid}, err
}

func (s *server) GetWorkers(ctx context.Context, in *pb.GetWorkersRequest) (*pb.GetWorkersReply, error) {
	var ws []*pb.Worker
	var err error
	if in.WorkerId == "" {
		ws, err = dbIf.GetWorkerList(in.StudyId, in.TrialId)
	} else {
		var w *pb.Worker
		w, err = dbIf.GetWorker(in.WorkerId)
		ws = append(ws, w)
	}
	return &pb.GetWorkersReply{Workers: ws}, err
}

func (s *server) GetShouldStopWorkers(ctx context.Context, in *pb.GetShouldStopWorkersRequest) (*pb.GetShouldStopWorkersReply, error) {
	if in.EarlyStoppingAlgorithm == "" {
		return &pb.GetShouldStopWorkersReply{}, errors.New("No EarlyStopping Algorithm specified")
	}
	conn, err := grpc.Dial("vizier-earlystopping-"+in.EarlyStoppingAlgorithm+":6789", grpc.WithInsecure())
	if err != nil {
		return &pb.GetShouldStopWorkersReply{}, err
	}
	defer conn.Close()
	c := pb.NewEarlyStoppingClient(conn)
	return c.GetShouldStopWorkers(context.Background(), in)
}

func (s *server) GetMetrics(ctx context.Context, in *pb.GetMetricsRequest) (*pb.GetMetricsReply, error) {
	var mNames []string
	if in.StudyId == "" {
		return &pb.GetMetricsReply{}, errors.New("StudyId should be set")
	}
	sc, err := dbIf.GetStudyConfig(in.StudyId)
	if err != nil {
		return &pb.GetMetricsReply{}, err
	}
	if len(in.MetricsNames) > 0 {
		mNames = in.MetricsNames
	} else {
		mNames = sc.Metrics
	}
	if err != nil {
		return &pb.GetMetricsReply{}, err
	}
	mls := make([]*pb.MetricsLogSet, len(in.WorkerIds))
	for i, w := range in.WorkerIds {
		wr, err := s.GetWorkers(ctx, &pb.GetWorkersRequest{
			StudyId:  in.StudyId,
			WorkerId: w,
		})
		if err != nil {
			return &pb.GetMetricsReply{}, err
		}
		mls[i] = &pb.MetricsLogSet{
			WorkerId:     w,
			MetricsLogs:  make([]*pb.MetricsLog, len(mNames)),
			WorkerStatus: wr.Workers[0].Status,
		}
		for j, m := range mNames {
			ls, err := dbIf.GetWorkerLogs(w, &kdb.GetWorkerLogOpts{Name: m})
			if err != nil {
				return &pb.GetMetricsReply{}, err
			}
			mls[i].MetricsLogs[j] = &pb.MetricsLog{
				Name:   m,
				Values: make([]*pb.MetricsValueTime, len(ls)),
			}
			for k, l := range ls {
				mls[i].MetricsLogs[j].Values[k] = &pb.MetricsValueTime{
					Value: l.Value,
					Time:  l.Time.UTC().Format(time.RFC3339Nano),
				}
			}
		}
	}
	return &pb.GetMetricsReply{MetricsLogSets: mls}, nil
}

func (s *server) ReportMetrics(ctx context.Context, in *pb.ReportMetricsRequest) (*pb.ReportMetricsReply, error) {
	sc, err := dbIf.GetStudyConfig(in.StudyId)
	if err != nil {
		return &pb.ReportMetricsReply{}, err
	}
	for _, mls := range in.MetricsLogSets {
		w, err := dbIf.GetWorker(mls.WorkerId)
		if err != nil {
			return &pb.ReportMetricsReply{}, err
		}
		trial, err := dbIf.GetTrial(w.TrialId)
		if err != nil {
			return &pb.ReportMetricsReply{}, err
		}
		err = dbIf.StoreWorkerLogs(mls.WorkerId, mls.MetricsLogs)
		if err != nil {
			return &pb.ReportMetricsReply{}, err
		}
		mets := []*pb.Metrics{}
		for _, ml := range mls.MetricsLogs {
			if ml != nil {
				if len(ml.Values) > 0 {
					mets = append(mets, &pb.Metrics{
						Name:  ml.Name,
						Value: ml.Values[len(ml.Values)-1].Value,
					})
				}
			}
		}
		if len(mets) > 0 {
			mi := &pb.ModelInfo{
				StudyName:  sc.Name,
				WorkerId:   mls.WorkerId,
				Parameters: trial.ParameterSet,
				Metrics:    mets,
			}
			smr := &pb.SaveModelRequest{
				Model:   mi,
				DataSet: &pb.DataSetInfo{},
			}
			_, err = s.SaveModel(ctx, smr)
			if err != nil {
				return &pb.ReportMetricsReply{}, err
			}
			err = dbIf.UpdateWorker(mls.WorkerId, mls.WorkerStatus)
		}
	}
	return &pb.ReportMetricsReply{}, nil

}

func (s *server) SetSuggestionParameters(ctx context.Context, in *pb.SetSuggestionParametersRequest) (*pb.SetSuggestionParametersReply, error) {
	var err error
	var id string
	if in.ParamId == "" {
		id, err = dbIf.SetSuggestionParam(in.SuggestionAlgorithm, in.StudyId, in.SuggestionParameters)
	} else {
		id = in.ParamId
		err = dbIf.UpdateSuggestionParam(in.ParamId, in.SuggestionParameters)
	}
	return &pb.SetSuggestionParametersReply{ParamId: id}, err
}

func (s *server) SetEarlyStoppingParameters(ctx context.Context, in *pb.SetEarlyStoppingParametersRequest) (*pb.SetEarlyStoppingParametersReply, error) {
	var err error
	var id string
	if in.ParamId == "" {
		id, err = dbIf.SetEarlyStopParam(in.EarlyStoppingAlgorithm, in.StudyId, in.EarlyStoppingParameters)
	} else {
		id = in.ParamId
		err = dbIf.UpdateEarlyStopParam(in.ParamId, in.EarlyStoppingParameters)
	}
	return &pb.SetEarlyStoppingParametersReply{ParamId: id}, err
}

func (s *server) GetSuggestionParameters(ctx context.Context, in *pb.GetSuggestionParametersRequest) (*pb.GetSuggestionParametersReply, error) {
	ps, err := dbIf.GetSuggestionParam(in.ParamId)
	return &pb.GetSuggestionParametersReply{SuggestionParameters: ps}, err
}

func (s *server) GetSuggestionParameterList(ctx context.Context, in *pb.GetSuggestionParameterListRequest) (*pb.GetSuggestionParameterListReply, error) {
	pss, err := dbIf.GetSuggestionParamList(in.StudyId)
	return &pb.GetSuggestionParameterListReply{SuggestionParameterSets: pss}, err
}

func (s *server) GetEarlyStoppingParameters(ctx context.Context, in *pb.GetEarlyStoppingParametersRequest) (*pb.GetEarlyStoppingParametersReply, error) {
	ps, err := dbIf.GetEarlyStopParam(in.ParamId)
	return &pb.GetEarlyStoppingParametersReply{EarlyStoppingParameters: ps}, err
}

func (s *server) GetEarlyStoppingParameterList(ctx context.Context, in *pb.GetEarlyStoppingParameterListRequest) (*pb.GetEarlyStoppingParameterListReply, error) {
	pss, err := dbIf.GetEarlyStopParamList(in.StudyId)
	return &pb.GetEarlyStoppingParameterListReply{EarlyStoppingParameterSets: pss}, err
}

func (s *server) SaveStudy(ctx context.Context, in *pb.SaveStudyRequest) (*pb.SaveStudyReply, error) {
	err := s.msIf.SaveStudy(in)
	return &pb.SaveStudyReply{}, err
}

func (s *server) SaveModel(ctx context.Context, in *pb.SaveModelRequest) (*pb.SaveModelReply, error) {
	err := s.msIf.SaveModel(in)
	if err != nil {
		log.Printf("Save Model failed %v", err)
		return &pb.SaveModelReply{}, err
	}
	return &pb.SaveModelReply{}, nil
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
	log.Printf("Start Katib manager: %s", port)
	s := grpc.NewServer(grpc.MaxRecvMsgSize(size), grpc.MaxSendMsgSize(size))
	pb.RegisterManagerServer(s, &server{msIf: modelstore.NewModelDB("modeldb-backend", "6543")})
	reflection.Register(s)
	if err = s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
