package main

import (
	"context"
	"net"

	health_pb "github.com/kubeflow/katib/pkg/apis/manager/health"
	"github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
	suggestion "github.com/kubeflow/katib/pkg/suggestion/v1alpha3/goptuna"
	"google.golang.org/grpc"
	"k8s.io/klog"
)

const (
	address = "0.0.0.0:6789"
)

type healthService struct {
}

func (s *healthService) Check(ctx context.Context, in *health_pb.HealthCheckRequest) (*health_pb.HealthCheckResponse, error) {
	return &health_pb.HealthCheckResponse{
		Status: health_pb.HealthCheckResponse_SERVING,
	}, nil
}

func main() {
	l, err := net.Listen("tcp", address)
	if err != nil {
		klog.Fatalf("Failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	api_v1_alpha3.RegisterSuggestionServer(srv, suggestion.NewSuggestionService())
	health_pb.RegisterHealthServer(srv, &healthService{})

	klog.Infof("Start Goptuna suggestion service: %s", address)
	err = srv.Serve(l)
	if err != nil {
		klog.Fatalf("Failed to serve: %v", err)
	}
}
