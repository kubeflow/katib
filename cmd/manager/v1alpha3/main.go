package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"

	health_pb "github.com/kubeflow/katib/pkg/apis/manager/health"
	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
	db "github.com/kubeflow/katib/pkg/db/v1alpha3"
	"github.com/kubeflow/katib/pkg/db/v1alpha3/common"
	"k8s.io/klog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = "0.0.0.0:6789"
)

var dbIf common.KatibDBInterface

type server struct {
}

// Register a Experiment to DB.
func (s *server) RegisterExperiment(ctx context.Context, in *api_pb.RegisterExperimentRequest) (*api_pb.RegisterExperimentReply, error) {
	err := dbIf.RegisterExperiment(in.Experiment)
	return &api_pb.RegisterExperimentReply{}, err
}

// Register a Experiment to DB.
func (s *server) PreCheckRegisterExperiment(ctx context.Context, in *api_pb.RegisterExperimentRequest) (*api_pb.PreCheckRegisterExperimentReply, error) {
	can_register, err := dbIf.PreCheckRegisterExperiment(in.Experiment)
	return &api_pb.PreCheckRegisterExperimentReply{
		CanRegister: can_register,
	}, err
}

// Delete a Experiment from DB by name.
func (s *server) DeleteExperiment(ctx context.Context, in *api_pb.DeleteExperimentRequest) (*api_pb.DeleteExperimentReply, error) {
	err := dbIf.DeleteExperiment(in.ExperimentName)
	return &api_pb.DeleteExperimentReply{}, err
}

// Get a Experiment from DB by name.
func (s *server) GetExperiment(ctx context.Context, in *api_pb.GetExperimentRequest) (*api_pb.GetExperimentReply, error) {
	exp, err := dbIf.GetExperiment(in.ExperimentName)
	return &api_pb.GetExperimentReply{
		Experiment: exp,
	}, err
}

// Get a summary list of Experiment from DB.
// The summary includes name and condition.
func (s *server) GetExperimentList(ctx context.Context, in *api_pb.GetExperimentListRequest) (*api_pb.GetExperimentListReply, error) {
	expList, err := dbIf.GetExperimentList()
	return &api_pb.GetExperimentListReply{
		ExperimentSummaries: expList,
	}, err
}

// Update Status of a experiment.
func (s *server) UpdateExperimentStatus(ctx context.Context, in *api_pb.UpdateExperimentStatusRequest) (*api_pb.UpdateExperimentStatusReply, error) {
	err := dbIf.UpdateExperimentStatus(in.ExperimentName, in.NewStatus)
	return &api_pb.UpdateExperimentStatusReply{}, err
}

// Update AlgorithmExtraSettings.
// The ExtraSetting is created if it does not exist, otherwise it is overwrited.
func (s *server) UpdateAlgorithmExtraSettings(ctx context.Context, in *api_pb.UpdateAlgorithmExtraSettingsRequest) (*api_pb.UpdateAlgorithmExtraSettingsReply, error) {
	err := dbIf.UpdateAlgorithmExtraSettings(in.ExperimentName, in.ExtraAlgorithmSettings)
	return &api_pb.UpdateAlgorithmExtraSettingsReply{}, err
}

// Get all AlgorithmExtraSettings.
func (s *server) GetAlgorithmExtraSettings(ctx context.Context, in *api_pb.GetAlgorithmExtraSettingsRequest) (*api_pb.GetAlgorithmExtraSettingsReply, error) {
	eas, err := dbIf.GetAlgorithmExtraSettings(in.ExperimentName)
	return &api_pb.GetAlgorithmExtraSettingsReply{
		ExtraAlgorithmSettings: eas,
	}, err
}

// Register a Trial to DB.
// ID will be filled by manager automatically.
func (s *server) RegisterTrial(ctx context.Context, in *api_pb.RegisterTrialRequest) (*api_pb.RegisterTrialReply, error) {
	err := dbIf.RegisterTrial(in.Trial)
	return &api_pb.RegisterTrialReply{}, err
}

// Delete a Trial from DB by ID.
func (s *server) DeleteTrial(ctx context.Context, in *api_pb.DeleteTrialRequest) (*api_pb.DeleteTrialReply, error) {
	err := dbIf.DeleteTrial(in.TrialName)
	return &api_pb.DeleteTrialReply{}, err
}

// Get a list of Trial from DB by name of a Experiment.
func (s *server) GetTrialList(ctx context.Context, in *api_pb.GetTrialListRequest) (*api_pb.GetTrialListReply, error) {
	trList, err := dbIf.GetTrialList(in.ExperimentName, in.Filter)
	return &api_pb.GetTrialListReply{
		Trials: trList,
	}, err
}

// Get a Trial from DB by ID of Trial.
func (s *server) GetTrial(ctx context.Context, in *api_pb.GetTrialRequest) (*api_pb.GetTrialReply, error) {
	tr, err := dbIf.GetTrial(in.TrialName)
	return &api_pb.GetTrialReply{
		Trial: tr,
	}, err
}

// Update Status of a trial.
func (s *server) UpdateTrialStatus(ctx context.Context, in *api_pb.UpdateTrialStatusRequest) (*api_pb.UpdateTrialStatusReply, error) {
	err := dbIf.UpdateTrialStatus(in.TrialName, in.NewStatus)
	return &api_pb.UpdateTrialStatusReply{}, err
}

// Report a log of Observations for a Trial.
// The log consists of timestamp and value of metric.
// Katib store every log of metrics.
// You can see accuracy curve or other metric logs on UI.
func (s *server) ReportObservationLog(ctx context.Context, in *api_pb.ReportObservationLogRequest) (*api_pb.ReportObservationLogReply, error) {
	err := dbIf.RegisterObservationLog(in.TrialName, in.ObservationLog)
	return &api_pb.ReportObservationLogReply{}, err
}

// Get all log of Observations for a Trial.
func (s *server) GetObservationLog(ctx context.Context, in *api_pb.GetObservationLogRequest) (*api_pb.GetObservationLogReply, error) {
	ol, err := dbIf.GetObservationLog(in.TrialName, in.MetricName, in.StartTime, in.EndTime)
	return &api_pb.GetObservationLogReply{
		ObservationLog: ol,
	}, err
}

// Delete all log of Observations for a Trial.
func (s *server) DeleteObservationLog(ctx context.Context, in *api_pb.DeleteObservationLogRequest) (*api_pb.DeleteObservationLogReply, error) {
	err := dbIf.DeleteObservationLog(in.TrialName)
	return &api_pb.DeleteObservationLogReply{}, err
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
	dbNameEnvName := common.DBNameEnvName
	dbName := os.Getenv(dbNameEnvName)
	if dbName == "" {
		klog.Fatal("DB_NAME env is not set. Exiting")
	}
	dbIf, err = db.NewKatibDBInterface(dbName)
	if err != nil {
		klog.Fatalf("Failed to open db connection: %v", err)
	}
	dbIf.DBInit()
	listener, err := net.Listen("tcp", port)
	if err != nil {
		klog.Fatalf("Failed to listen: %v", err)
	}

	size := 1<<31 - 1
	klog.Infof("Start Katib manager: %s", port)
	s := grpc.NewServer(grpc.MaxRecvMsgSize(size), grpc.MaxSendMsgSize(size))
	api_pb.RegisterManagerServer(s, &server{})
	health_pb.RegisterHealthServer(s, &server{})
	reflection.Register(s)
	if err = s.Serve(listener); err != nil {
		klog.Fatalf("Failed to serve: %v", err)
	}
}
