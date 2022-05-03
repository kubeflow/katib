/*
Copyright 2022 The Kubeflow Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"net"

	health_pb "github.com/kubeflow/katib/pkg/apis/manager/health"
	api_v1_beta1 "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	suggestion "github.com/kubeflow/katib/pkg/suggestion/v1beta1/goptuna"
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
	api_v1_beta1.RegisterSuggestionServer(srv, suggestion.NewSuggestionService())
	health_pb.RegisterHealthServer(srv, &healthService{})

	klog.Infof("Start Goptuna suggestion service: %s", address)
	err = srv.Serve(l)
	if err != nil {
		klog.Fatalf("Failed to serve: %v", err)
	}
}
