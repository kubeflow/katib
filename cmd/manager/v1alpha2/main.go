package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	health_pb "github.com/kubeflow/katib/pkg/api/health"
	api_pb "github.com/kubeflow/katib/pkg/api/v1alpha2"
	dbif "github.com/kubeflow/katib/pkg/api/v1alpha2/dbif"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port        = "0.0.0.0:6789"
	dbIfaddress = "mysql-db-backend:6789"
)

var dbIf dbif.DBIFClient

type server struct {
}

// Register a Experiment to DB.
func (s *server) RegisterExperiment(ctx context.Context, in *dbif.RegisterExperimentRequest) (*dbif.RegisterExperimentReply, error) {
	return dbIf.RegisterExperiment(ctx, in)
}

// Delete a Experiment from DB by name.
func (s *server) DeleteExperiment(ctx context.Context, in *dbif.DeleteExperimentRequest) (*dbif.DeleteExperimentReply, error) {
	return dbIf.DeleteExperiment(ctx, in)
}

// Get a Experiment from DB by name.
func (s *server) GetExperiment(ctx context.Context, in *dbif.GetExperimentRequest) (*dbif.GetExperimentReply, error) {
	return dbIf.GetExperiment(ctx, in)
}

// Get a summary list of Experiment from DB.
// The summary includes name and condition.
func (s *server) GetExperimentList(ctx context.Context, in *dbif.GetExperimentListRequest) (*dbif.GetExperimentListReply, error) {
	return dbIf.GetExperimentList(ctx, in)
}

// Update Status of a experiment.
func (s *server) UpdateExperimentStatus(ctx context.Context, in *dbif.UpdateExperimentStatusRequest) (*dbif.UpdateExperimentStatusReply, error) {
	return dbIf.UpdateExperimentStatus(ctx, in)
}

// Update AlgorithmExtraSettings.
// The ExtraSetting is created if it does not exist, otherwise it is overwrited.
func (s *server) UpdateAlgorithmExtraSettings(ctx context.Context, in *dbif.UpdateAlgorithmExtraSettingsRequest) (*dbif.UpdateAlgorithmExtraSettingsReply, error) {
	return dbIf.UpdateAlgorithmExtraSettings(ctx, in)
}

// Get all AlgorithmExtraSettings.
func (s *server) GetAlgorithmExtraSettings(ctx context.Context, in *dbif.GetAlgorithmExtraSettingsRequest) (*dbif.GetAlgorithmExtraSettingsReply, error) {
	return dbIf.GetAlgorithmExtraSettings(ctx, in)
}

// Register a Trial to DB.
// ID will be filled by manager automatically.
func (s *server) RegisterTrial(ctx context.Context, in *dbif.RegisterTrialRequest) (*dbif.RegisterTrialReply, error) {
	return dbIf.RegisterTrial(ctx, in)
}

// Delete a Trial from DB by ID.
func (s *server) DeleteTrial(ctx context.Context, in *dbif.DeleteTrialRequest) (*dbif.DeleteTrialReply, error) {
	return dbIf.DeleteTrial(ctx, in)
}

// Get a list of Trial from DB by name of a Experiment.
func (s *server) GetTrialList(ctx context.Context, in *dbif.GetTrialListRequest) (*dbif.GetTrialListReply, error) {
	return dbIf.GetTrialList(ctx, in)
}

// Get a Trial from DB by ID of Trial.
func (s *server) GetTrial(ctx context.Context, in *dbif.GetTrialRequest) (*dbif.GetTrialReply, error) {
	return dbIf.GetTrial(ctx, in)
}

// Update Status of a trial.
func (s *server) UpdateTrialStatus(ctx context.Context, in *dbif.UpdateTrialStatusRequest) (*dbif.UpdateTrialStatusReply, error) {
	return dbIf.UpdateTrialStatus(ctx, in)
}

// Report a log of Observations for a Trial.
// The log consists of timestamp and value of metric.
// Katib store every log of metrics.
// You can see accuracy curve or other metric logs on UI.
func (s *server) ReportObservationLog(ctx context.Context, in *dbif.ReportObservationLogRequest) (*dbif.ReportObservationLogReply, error) {
	return dbIf.ReportObservationLog(ctx, in)
}

// Get all log of Observations for a Trial.
func (s *server) GetObservationLog(ctx context.Context, in *dbif.GetObservationLogRequest) (*dbif.GetObservationLogReply, error) {
	return dbIf.GetObservationLog(ctx, in)
}

func (s *server) getSuggestionServiceConnection(algoName string) (*grpc.ClientConn, error) {
	if algoName == "" {
		return nil, errors.New("No algorithm name is specified")
	}
	return grpc.Dial("katib-suggestion-"+algoName+":6789", grpc.WithInsecure())
}

// Get Suggestions from a Suggestion service.
func (s *server) GetSuggestions(ctx context.Context, in *api_pb.GetSuggestionsRequest) (*api_pb.GetSuggestionsReply, error) {
	conn, err := s.getSuggestionServiceConnection(in.AlgorithmName)
	if err != nil {
		return &api_pb.GetSuggestionsReply{}, err
	}
	defer conn.Close()
	c := api_pb.NewSuggestionClient(conn)
	r, err := c.GetSuggestions(ctx, in)
	if err != nil {
		return &api_pb.GetSuggestionsReply{}, err
	}
	return r, nil
}

// Validate AlgorithmSettings in an Experiment.
// Suggestion service should return INVALID_ARGUMENT Error when the parameter is invalid
func (s *server) ValidateAlgorithmSettings(ctx context.Context, in *api_pb.ValidateAlgorithmSettingsRequest) (*api_pb.ValidateAlgorithmSettingsReply, error) {
	conn, err := s.getSuggestionServiceConnection(in.AlgorithmName)
	if err != nil {
		return &api_pb.ValidateAlgorithmSettingsReply{}, err
	}
	defer conn.Close()
	c := api_pb.NewSuggestionClient(conn)
	return c.ValidateAlgorithmSettings(ctx, in)
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

	// Check if connection to katib-db is okay since otherwise manager could not serve most of its methods.
	_, err := dbIf.SelectOne(ctx, &dbif.SelectOneRequest{})
	if err != nil {
		resp.Status = health_pb.HealthCheckResponse_NOT_SERVING
		return &resp, fmt.Errorf("Failed to execute `SELECT 1` probe: %v", err)
	}

	return &resp, nil
}

func main() {
	flag.Parse()
	var err error

	conn, err := grpc.Dial(dbIfaddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to DBIF service: %v", err)
	}
	defer conn.Close()
	dbIf = dbif.NewDBIFClient(conn)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	size := 1<<31 - 1
	log.Printf("Start Katib manager: %s", port)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, err = dbIf.RegisterExperiment(ctx, &dbif.RegisterExperimentRequest{Experiment: &dbif.Experiment{
			Name: "testExp",
			ExperimentSpec: &dbif.ExperimentSpec{
				ParameterSpecs: &dbif.ExperimentSpec_ParameterSpecs{
					Parameters: []*dbif.ParameterSpec{},
				},
				Objective: &dbif.ObjectiveSpec{
					Type:                   dbif.ObjectiveType_UNKNOWN,
					Goal:                   0.99,
					ObjectiveMetricName:    "f1_score",
					AdditionalMetricNames:  []string{"loss", "precision", "recall"},
				},
				Algorithm:          &dbif.AlgorithmSpec{},
				TrialTemplate:      "",
				ParallelTrialCount: 10,
				MaxTrialCount:      100,
			},
			ExperimentStatus: &dbif.ExperimentStatus{
				Condition:      dbif.ExperimentStatus_CREATED,
				StartTime:      "2019-02-03T04:05:06+09:00",
				CompletionTime: "2019-02-03T04:05:06+09:00",
			},
		},})
	if err != nil {
		log.Fatalf("could not create experiment: %v", err)
	}
	log.Printf("Experiment created with id: %s")
	s := grpc.NewServer(grpc.MaxRecvMsgSize(size), grpc.MaxSendMsgSize(size))
	api_pb.RegisterManagerServer(s, &server{})
	health_pb.RegisterHealthServer(s, &server{})
	reflection.Register(s)
	if err = s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
