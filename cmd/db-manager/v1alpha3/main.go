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

	// Check if connection to katib db driver is okay since otherwise manager could not serve most of its methods.
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
