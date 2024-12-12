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
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	health_pb "github.com/kubeflow/katib/pkg/apis/manager/health"
	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	db "github.com/kubeflow/katib/pkg/db/v1beta1"
	"github.com/kubeflow/katib/pkg/db/v1beta1/common"
	"k8s.io/klog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	defaultListenAddress  = "0.0.0.0:6789"
	defaultConnectTimeout = time.Second * 60
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
	var connectTimeout time.Duration
	var listenAddress string
	flag.DurationVar(&connectTimeout, "connect-timeout", defaultConnectTimeout, "Timeout before calling error during database connection. (e.g. 120s)")
	flag.StringVar(&listenAddress, "listen-address", defaultListenAddress, "The network interface or IP address to receive incoming connections. (e.g. 0.0.0.0:6789)")
	flag.Parse()

	var err error
	dbNameEnvName := common.DBNameEnvName
	dbName := os.Getenv(dbNameEnvName)
	if dbName == "" {
		klog.Fatal("DB_NAME env is not set. Exiting")
	}
	dbIf, err = db.NewKatibDBInterface(dbName, connectTimeout)
	if err != nil {
		klog.Fatalf("Failed to open db connection: %v", err)
	}
	dbIf.DBInit()
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		klog.Fatalf("Failed to listen: %v", err)
	}

	size := 1<<31 - 1
	klog.Infof("Start Katib manager: %s", listenAddress)
	s := grpc.NewServer(grpc.MaxRecvMsgSize(size), grpc.MaxSendMsgSize(size))
	api_pb.RegisterDBManagerServer(s, &server{})
	health_pb.RegisterHealthServer(s, &server{})
	reflection.Register(s)
	if err = s.Serve(listener); err != nil {
		klog.Fatalf("Failed to serve: %v", err)
	}
}
