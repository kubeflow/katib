package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"

	health_pb "github.com/kubeflow/katib/pkg/apis/manager/health"
	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	db "github.com/kubeflow/katib/pkg/db/v1beta1"
	"github.com/kubeflow/katib/pkg/db/v1beta1/common"
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
	dbPasswordEnvName := common.DBPasswordEnvName
	dbPassword := os.Getenv(dbPasswordEnvName)
	if dbPassword == "" {
		klog.Fatal("DB_PASSWORD env is not set or empty. Exiting")
	}
	dbIf, err = db.NewKatibDBInterface()
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
	api_pb.RegisterDBManagerServer(s, &server{})
	health_pb.RegisterHealthServer(s, &server{})
	reflection.Register(s)
	if err = s.Serve(listener); err != nil {
		klog.Fatalf("Failed to serve: %v", err)
	}
}
