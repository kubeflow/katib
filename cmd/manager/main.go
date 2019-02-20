package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	api_pb "github.com/kubeflow/katib/pkg/api"
	health_pb "github.com/kubeflow/katib/pkg/api/health"
	kdb "github.com/kubeflow/katib/pkg/db"
	"github.com/kubeflow/katib/pkg/manager/modelstore"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = "0.0.0.0:6789"
)

var dbIf kdb.VizierDBInterface

type server struct {
	msIf modelstore.ModelStore
}

func (s *server) CreateStudy(ctx context.Context, in *api_pb.CreateStudyRequest) (*api_pb.CreateStudyReply, error) {
	if in == nil || in.StudyConfig == nil {
		return &api_pb.CreateStudyReply{}, errors.New("StudyConfig is missing.")
	}

	studyID, err := dbIf.CreateStudy(in.StudyConfig)
	if err != nil {
		return &api_pb.CreateStudyReply{}, err
	}

	return &api_pb.CreateStudyReply{StudyId: studyID}, nil
}

func (s *server) DeleteStudy(ctx context.Context, in *api_pb.DeleteStudyRequest) (*api_pb.DeleteStudyReply, error) {
	if in == nil || in.StudyId == "" {
		return &api_pb.DeleteStudyReply{}, errors.New("StudyId is missing.")
	}
	err := dbIf.DeleteStudy(in.StudyId)
	if err != nil {
		return &api_pb.DeleteStudyReply{}, err
	}
	return &api_pb.DeleteStudyReply{StudyId: in.StudyId}, nil
}

func (s *server) GetStudy(ctx context.Context, in *api_pb.GetStudyRequest) (*api_pb.GetStudyReply, error) {

	sc, err := dbIf.GetStudy(in.StudyId)

	if err != nil {
		return &api_pb.GetStudyReply{}, err
	}
	return &api_pb.GetStudyReply{StudyConfig: sc}, err
}

func (s *server) GetStudyList(ctx context.Context, in *api_pb.GetStudyListRequest) (*api_pb.GetStudyListReply, error) {
	sl, err := dbIf.GetStudyList()
	if err != nil {
		return &api_pb.GetStudyListReply{}, err
	}
	result := make([]*api_pb.StudyOverview, len(sl))
	for i, id := range sl {
		sc, err := dbIf.GetStudy(id)
		if err != nil {
			return &api_pb.GetStudyListReply{}, err
		}
		result[i] = &api_pb.StudyOverview{
			Name:  sc.Name,
			Owner: sc.Owner,
			Id:    id,
		}
	}
	return &api_pb.GetStudyListReply{StudyOverviews: result}, err
}

func (s *server) CreateTrial(ctx context.Context, in *api_pb.CreateTrialRequest) (*api_pb.CreateTrialReply, error) {
	err := dbIf.CreateTrial(in.Trial)
	return &api_pb.CreateTrialReply{TrialId: in.Trial.TrialId}, err
}

func (s *server) GetTrials(ctx context.Context, in *api_pb.GetTrialsRequest) (*api_pb.GetTrialsReply, error) {
	tl, err := dbIf.GetTrialList(in.StudyId)
	return &api_pb.GetTrialsReply{Trials: tl}, err
}

func (s *server) GetTrial(ctx context.Context, in *api_pb.GetTrialRequest) (*api_pb.GetTrialReply, error) {
	t, err := dbIf.GetTrial(in.TrialId)
	return &api_pb.GetTrialReply{Trial: t}, err
}

func (s *server) GetSuggestions(ctx context.Context, in *api_pb.GetSuggestionsRequest) (*api_pb.GetSuggestionsReply, error) {
	if in.SuggestionAlgorithm == "" {
		return &api_pb.GetSuggestionsReply{Trials: []*api_pb.Trial{}}, errors.New("No suggest algorithm specified")
	}
	conn, err := grpc.Dial("vizier-suggestion-"+in.SuggestionAlgorithm+":6789", grpc.WithInsecure())
	if err != nil {
		return &api_pb.GetSuggestionsReply{Trials: []*api_pb.Trial{}}, err
	}

	defer conn.Close()
	c := api_pb.NewSuggestionClient(conn)
	r, err := c.GetSuggestions(ctx, in)
	if err != nil {
		return &api_pb.GetSuggestionsReply{Trials: []*api_pb.Trial{}}, err
	}
	return r, nil
}

func (s *server) RegisterWorker(ctx context.Context, in *api_pb.RegisterWorkerRequest) (*api_pb.RegisterWorkerReply, error) {
	wid, err := dbIf.CreateWorker(in.Worker)
	return &api_pb.RegisterWorkerReply{WorkerId: wid}, err
}

func (s *server) GetWorkers(ctx context.Context, in *api_pb.GetWorkersRequest) (*api_pb.GetWorkersReply, error) {
	var ws []*api_pb.Worker
	var err error
	if in.WorkerId == "" {
		ws, err = dbIf.GetWorkerList(in.StudyId, in.TrialId)
	} else {
		var w *api_pb.Worker
		w, err = dbIf.GetWorker(in.WorkerId)
		ws = append(ws, w)
	}
	return &api_pb.GetWorkersReply{Workers: ws}, err
}

func (s *server) GetShouldStopWorkers(ctx context.Context, in *api_pb.GetShouldStopWorkersRequest) (*api_pb.GetShouldStopWorkersReply, error) {
	if in.EarlyStoppingAlgorithm == "" {
		return &api_pb.GetShouldStopWorkersReply{}, errors.New("No EarlyStopping Algorithm specified")
	}
	conn, err := grpc.Dial("vizier-earlystopping-"+in.EarlyStoppingAlgorithm+":6789", grpc.WithInsecure())
	if err != nil {
		return &api_pb.GetShouldStopWorkersReply{}, err
	}
	defer conn.Close()
	c := api_pb.NewEarlyStoppingClient(conn)
	return c.GetShouldStopWorkers(context.Background(), in)
}

func (s *server) GetMetrics(ctx context.Context, in *api_pb.GetMetricsRequest) (*api_pb.GetMetricsReply, error) {
	var mNames []string
	if in.StudyId == "" {
		return &api_pb.GetMetricsReply{}, errors.New("StudyId should be set")
	}
	sc, err := dbIf.GetStudy(in.StudyId)
	if err != nil {
		return &api_pb.GetMetricsReply{}, err
	}
	if len(in.MetricsNames) > 0 {
		mNames = in.MetricsNames
	} else {
		mNames = sc.Metrics
	}
	if err != nil {
		return &api_pb.GetMetricsReply{}, err
	}
	mls := make([]*api_pb.MetricsLogSet, len(in.WorkerIds))
	for i, w := range in.WorkerIds {
		wr, err := s.GetWorkers(ctx, &api_pb.GetWorkersRequest{
			StudyId:  in.StudyId,
			WorkerId: w,
		})
		if err != nil {
			return &api_pb.GetMetricsReply{}, err
		}
		mls[i] = &api_pb.MetricsLogSet{
			WorkerId:     w,
			MetricsLogs:  make([]*api_pb.MetricsLog, len(mNames)),
			WorkerStatus: wr.Workers[0].Status,
		}
		for j, m := range mNames {
			ls, err := dbIf.GetWorkerLogs(w, &kdb.GetWorkerLogOpts{Name: m})
			if err != nil {
				return &api_pb.GetMetricsReply{}, err
			}
			mls[i].MetricsLogs[j] = &api_pb.MetricsLog{
				Name:   m,
				Values: make([]*api_pb.MetricsValueTime, len(ls)),
			}
			for k, l := range ls {
				mls[i].MetricsLogs[j].Values[k] = &api_pb.MetricsValueTime{
					Value: l.Value,
					Time:  l.Time.UTC().Format(time.RFC3339Nano),
				}
			}
		}
	}
	return &api_pb.GetMetricsReply{MetricsLogSets: mls}, nil
}

func (s *server) ReportMetricsLogs(ctx context.Context, in *api_pb.ReportMetricsLogsRequest) (*api_pb.ReportMetricsLogsReply, error) {
	for _, mls := range in.MetricsLogSets {
		err := dbIf.StoreWorkerLogs(mls.WorkerId, mls.MetricsLogs)
		if err != nil {
			return &api_pb.ReportMetricsLogsReply{}, err
		}

	}
	return &api_pb.ReportMetricsLogsReply{}, nil
}

func (s *server) UpdateWorkerState(ctx context.Context, in *api_pb.UpdateWorkerStateRequest) (*api_pb.UpdateWorkerStateReply, error) {
	err := dbIf.UpdateWorker(in.WorkerId, in.Status)
	return &api_pb.UpdateWorkerStateReply{}, err
}

func (s *server) GetWorkerFullInfo(ctx context.Context, in *api_pb.GetWorkerFullInfoRequest) (*api_pb.GetWorkerFullInfoReply, error) {
	return dbIf.GetWorkerFullInfo(in.StudyId, in.TrialId, in.WorkerId, in.OnlyLatestLog)
}

func (s *server) SetSuggestionParameters(ctx context.Context, in *api_pb.SetSuggestionParametersRequest) (*api_pb.SetSuggestionParametersReply, error) {
	var err error
	var id string
	if in.ParamId == "" {
		id, err = dbIf.SetSuggestionParam(in.SuggestionAlgorithm, in.StudyId, in.SuggestionParameters)
	} else {
		id = in.ParamId
		err = dbIf.UpdateSuggestionParam(in.ParamId, in.SuggestionParameters)
	}
	return &api_pb.SetSuggestionParametersReply{ParamId: id}, err
}

func (s *server) SetEarlyStoppingParameters(ctx context.Context, in *api_pb.SetEarlyStoppingParametersRequest) (*api_pb.SetEarlyStoppingParametersReply, error) {
	var err error
	var id string
	if in.ParamId == "" {
		id, err = dbIf.SetEarlyStopParam(in.EarlyStoppingAlgorithm, in.StudyId, in.EarlyStoppingParameters)
	} else {
		id = in.ParamId
		err = dbIf.UpdateEarlyStopParam(in.ParamId, in.EarlyStoppingParameters)
	}
	return &api_pb.SetEarlyStoppingParametersReply{ParamId: id}, err
}

func (s *server) GetSuggestionParameters(ctx context.Context, in *api_pb.GetSuggestionParametersRequest) (*api_pb.GetSuggestionParametersReply, error) {
	ps, err := dbIf.GetSuggestionParam(in.ParamId)
	return &api_pb.GetSuggestionParametersReply{SuggestionParameters: ps}, err
}

func (s *server) GetSuggestionParameterList(ctx context.Context, in *api_pb.GetSuggestionParameterListRequest) (*api_pb.GetSuggestionParameterListReply, error) {
	pss, err := dbIf.GetSuggestionParamList(in.StudyId)
	return &api_pb.GetSuggestionParameterListReply{SuggestionParameterSets: pss}, err
}

func (s *server) GetEarlyStoppingParameters(ctx context.Context, in *api_pb.GetEarlyStoppingParametersRequest) (*api_pb.GetEarlyStoppingParametersReply, error) {
	ps, err := dbIf.GetEarlyStopParam(in.ParamId)
	return &api_pb.GetEarlyStoppingParametersReply{EarlyStoppingParameters: ps}, err
}

func (s *server) GetEarlyStoppingParameterList(ctx context.Context, in *api_pb.GetEarlyStoppingParameterListRequest) (*api_pb.GetEarlyStoppingParameterListReply, error) {
	pss, err := dbIf.GetEarlyStopParamList(in.StudyId)
	return &api_pb.GetEarlyStoppingParameterListReply{EarlyStoppingParameterSets: pss}, err
}

func (s *server) SaveStudy(ctx context.Context, in *api_pb.SaveStudyRequest) (*api_pb.SaveStudyReply, error) {
	var err error
	if s.msIf != nil {
		err = s.msIf.SaveStudy(in)
	}
	return &api_pb.SaveStudyReply{}, err
}

func (s *server) SaveModel(ctx context.Context, in *api_pb.SaveModelRequest) (*api_pb.SaveModelReply, error) {
	if s.msIf != nil {
		err := s.msIf.SaveModel(in)
		if err != nil {
			log.Printf("Save Model failed %v", err)
			return &api_pb.SaveModelReply{}, err
		}
	}
	return &api_pb.SaveModelReply{}, nil
}

func (s *server) GetSavedStudies(ctx context.Context, in *api_pb.GetSavedStudiesRequest) (*api_pb.GetSavedStudiesReply, error) {
	ret := []*api_pb.StudyOverview{}
	var err error
	if s.msIf != nil {
		ret, err = s.msIf.GetSavedStudies()
	}
	return &api_pb.GetSavedStudiesReply{Studies: ret}, err
}

func (s *server) GetSavedModels(ctx context.Context, in *api_pb.GetSavedModelsRequest) (*api_pb.GetSavedModelsReply, error) {
	ret := []*api_pb.ModelInfo{}
	var err error
	if s.msIf != nil {
		ret, err = s.msIf.GetSavedModels(in)
	}
	return &api_pb.GetSavedModelsReply{Models: ret}, err
}

func (s *server) GetSavedModel(ctx context.Context, in *api_pb.GetSavedModelRequest) (*api_pb.GetSavedModelReply, error) {
	var ret *api_pb.ModelInfo = nil
	var err error
	if s.msIf != nil {
		ret, err = s.msIf.GetSavedModel(in)
	}
	return &api_pb.GetSavedModelReply{Model: ret}, err
}

func (s *server) Check(ctx context.Context, in *health_pb.HealthCheckRequest) (*health_pb.HealthCheckResponse, error) {
	resp := health_pb.HealthCheckResponse{
		Status: health_pb.HealthCheckResponse_SERVING,
	}

	// We only accept optional service name only if it's set to suggested format.
	if in != nil && in.Service != "" && in.Service != "grpc.health.v1.Health" {
		resp.Status = health_pb.HealthCheckResponse_UNKNOWN
		return &resp, fmt.Errorf("grpc.health.v1.Health can only be accepted if you specify service name.")
	}

	// Check if connection to vizier-db is okay since otherwise manager could not serve most of its methods.
	err := dbIf.SelectOne()
	if err != nil {
		resp.Status = health_pb.HealthCheckResponse_NOT_SERVING
		return &resp, fmt.Errorf("Failed to execute `SELECT 1` probe: %v", err)
	}

	return &resp, nil
}

func main() {
	flag.Parse()
	var err error

	dbIf, err = kdb.New()
	if err != nil {
		log.Fatalf("Failed to open db connection: %v", err)
	}
	dbIf.DBInit()
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	size := 1<<31 - 1
	log.Printf("Start Katib manager: %s", port)
	s := grpc.NewServer(grpc.MaxRecvMsgSize(size), grpc.MaxSendMsgSize(size))
	api_pb.RegisterManagerServer(s, &server{})
	health_pb.RegisterHealthServer(s, &server{})
	reflection.Register(s)
	if err = s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
