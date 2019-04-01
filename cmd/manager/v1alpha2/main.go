package main

import (
	"context"
	"flag"
	"log"
	"net"

	api_pb "github.com/kubeflow/katib/pkg/api/v1alpha2"
	health_pb "github.com/kubeflow/katib/pkg/api/v1alpha2/health"
	kdb "github.com/kubeflow/katib/pkg/db/v1alpha2"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = "0.0.0.0:6789"
)

var dbIf kdb.KatibDBInterface

type server struct {
}

// Register a Experiment to DB.
func (s *server) RegisterExperiment(context.Context, *api_pb.RegisterExperimentRequest) (*api_pb.RegisterExperimentReply, error) {
	return nil, nil
}

// Delete a Experiment from DB by name.
func (s *server) DeleteExperiment(context.Context, *api_pb.DeleteExperimentRequest) (*api_pb.DeleteExperimentReply, error) {
	return nil, nil
}

// Get a Experiment from DB by name.
func (s *server) GetExperiment(context.Context, *api_pb.GetExperimentRequest) (*api_pb.GetExperimentReply, error) {
	return nil, nil
}

// Get a summary list of Experiment from DB.
// The summary includes name and condition.
func (s *server) GetExperimentList(context.Context, *api_pb.GetExperimentListRequest) (*api_pb.GetExperimentListReply, error) {
	return nil, nil
}

// Update Status of a experiment.
func (s *server) UpdateExperimentStatus(context.Context, *api_pb.UpdateExperimentStatusRequest) (*api_pb.UpdateExperimentStatusReply, error) {
	return nil, nil
}

// Update AlgorithmExtraSettings.
// The ExtraSetting is created if it does not exist, otherwise it is overwrited.
func (s *server) UpdateAlgorithmExtraSettings(context.Context, *api_pb.UpdateAlgorithmExtraSettingsRequest) (*api_pb.UpdateAlgorithmExtraSettingsReply, error) {
	return nil, nil
}

// Get all AlgorithmExtraSettings.
func (s *server) GetAlgorithmExtraSettings(context.Context, *api_pb.GetAlgorithmExtraSettingsRequest) (*api_pb.GetAlgorithmExtraSettingsReply, error) {
	return nil, nil
}

// Register a Trial to DB.
// ID will be filled by manager automatically.
func (s *server) RegisterTrial(context.Context, *api_pb.RegisterTrialRequest) (*api_pb.RegisterTrialReply, error) {
	return nil, nil
}

// Delete a Trial from DB by ID.
func (s *server) DeleteTrial(context.Context, *api_pb.DeleteTrialRequest) (*api_pb.DeleteTrialReply, error) {
	return nil, nil
}

// Get a list of Trial from DB by name of a Experiment.
func (s *server) GetTrialList(context.Context, *api_pb.GetTrialListRequest) (*api_pb.GetTrialListReply, error) {
	return nil, nil
}

// Get a Trial from DB by ID of Trial.
func (s *server) GetTrial(context.Context, *api_pb.GetTrialRequest) (*api_pb.GetTrialReply, error) {
	return nil, nil
}

// Update Status of a trial.
func (s *server) UpdateTrialStatus(context.Context, *api_pb.UpdateTrialStatusRequest) (*api_pb.UpdateTrialStatusReply, error) {
	return nil, nil
}

// Report a log of Observations for a Trial.
// The log consists of timestamp and value of metric.
// Katib store every log of metrics.
// You can see accuracy curve or other metric logs on UI.
func (s *server) ReportObservationLog(context.Context, *api_pb.ReportObservationLogRequest) (*api_pb.ReportObservationLogReply, error) {
	return nil, nil
}

// Get all log of Observations for a Trial.
func (s *server) GetObservationLog(context.Context, *api_pb.GetObservationLogRequest) (*api_pb.GetObservationLogReply, error) {
	return nil, nil
}

// Get Suggestions from a Suggestion service.
func (s *server) GetSuggestions(context.Context, *api_pb.GetSuggestionsRequest) (*api_pb.GetSuggestionsReply, error) {
	return nil, nil
}

// Validate AlgorithmSettings in an Experiment.
// Suggestion service should return INVALID_ARGUMENT Error when the parameter is invalid
func (s *server) ValidateAlgorithmSettings(context.Context, *api_pb.ValidateAlgorithmSettingsRequest) (*api_pb.ValidateAlgorithmSettingsReply, error) {
	return nil, nil
}

func (s *server) Check(context.Context, *health_pb.HealthCheckRequest) (*health_pb.HealthCheckResponse, error) {
	return nil, nil
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
