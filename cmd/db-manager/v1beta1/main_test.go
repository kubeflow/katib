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
	"testing"

	"go.uber.org/mock/gomock"

	health_pb "github.com/kubeflow/katib/pkg/apis/manager/health"
	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	mockdb "github.com/kubeflow/katib/pkg/mock/v1beta1/db"
)

func TestReportObservationLog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := &server{}
	mockDB := mockdb.NewMockKatibDBInterface(ctrl)
	dbIf = mockDB

	req := &api_pb.ReportObservationLogRequest{
		TrialName: "test1-trial1",
		ObservationLog: &api_pb.ObservationLog{
			MetricLogs: []*api_pb.MetricLog{
				{
					TimeStamp: "2019-02-03T04:05:06+09:00",
					Metric: &api_pb.Metric{
						Name:  "f1_score",
						Value: "88.95",
					},
				},
				{
					TimeStamp: "2019-02-03T04:05:06+09:00",
					Metric: &api_pb.Metric{
						Name:  "loss",
						Value: "0.5",
					},
				},
				{
					TimeStamp: "2019-02-03T04:05:06+09:00",
					Metric: &api_pb.Metric{
						Name:  "precision",
						Value: "88.7",
					},
				},
				{
					TimeStamp: "2019-02-03T04:05:06+09:00",
					Metric: &api_pb.Metric{
						Name:  "recall",
						Value: "89.2",
					},
				},
			},
		},
	}
	mockDB.EXPECT().RegisterObservationLog(req.TrialName, req.ObservationLog).Return(nil)
	_, err := s.ReportObservationLog(context.Background(), req)
	if err != nil {
		t.Fatalf("ReportObservationLog Error %v", err)
	}
}

func TestGetObservationLog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := &server{}
	mockDB := mockdb.NewMockKatibDBInterface(ctrl)
	dbIf = mockDB

	req := &api_pb.GetObservationLogRequest{
		TrialName: "test1-trial1",
		StartTime: "2019-02-03T03:05:06+09:00",
		EndTime:   "2019-02-03T05:05:06+09:00",
	}

	obs := &api_pb.ObservationLog{
		MetricLogs: []*api_pb.MetricLog{
			{
				TimeStamp: "2019-02-03T04:05:06+09:00",
				Metric: &api_pb.Metric{
					Name:  "f1_score",
					Value: "88.95",
				},
			},
			{
				TimeStamp: "2019-02-03T04:05:06+09:00",
				Metric: &api_pb.Metric{
					Name:  "loss",
					Value: "0.5",
				},
			},
			{
				TimeStamp: "2019-02-03T04:05:06+09:00",
				Metric: &api_pb.Metric{
					Name:  "precision",
					Value: "88.7",
				},
			},
			{
				TimeStamp: "2019-02-03T04:05:06+09:00",
				Metric: &api_pb.Metric{
					Name:  "recall",
					Value: "89.2",
				},
			},
		},
	}

	mockDB.EXPECT().GetObservationLog(req.TrialName, req.MetricName, req.StartTime, req.EndTime).Return(obs, nil)
	ret, err := s.GetObservationLog(context.Background(), req)
	if err != nil {
		t.Fatalf("GetObservationLog Error %v", err)
	}
	if len(obs.MetricLogs) != len(ret.ObservationLog.MetricLogs) {
		t.Fatalf("GetObservationLog Test fail expect metrics number %d got %d", len(obs.MetricLogs), len(ret.ObservationLog.MetricLogs))
	}
}

func TestDeleteObservationLog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := &server{}
	mockDB := mockdb.NewMockKatibDBInterface(ctrl)
	dbIf = mockDB

	req := &api_pb.DeleteObservationLogRequest{
		TrialName: "test1-trial1",
	}
	mockDB.EXPECT().DeleteObservationLog(req.TrialName).Return(nil)
	_, err := s.DeleteObservationLog(context.Background(), req)
	if err != nil {
		t.Fatalf("DeleteExperiment Error %v", err)
	}
}

func TestCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := &server{}
	mockDB := mockdb.NewMockKatibDBInterface(ctrl)
	dbIf = mockDB
	testCases := []struct {
		Request        *health_pb.HealthCheckRequest
		ExpectedStatus health_pb.HealthCheckResponse_ServingStatus
		Name           string
	}{
		{
			Request: &health_pb.HealthCheckRequest{
				Service: "grpc.health.v1.Health",
			},
			ExpectedStatus: health_pb.HealthCheckResponse_SERVING,
			Name:           "Valid Request",
		},
		{
			Request: &health_pb.HealthCheckRequest{
				Service: "grpc.health.v1.1.Health",
			},
			ExpectedStatus: health_pb.HealthCheckResponse_UNKNOWN,
			Name:           "Invalid service name",
		},
	}

	mockDB.EXPECT().SelectOne().Return(nil)

	for _, tc := range testCases {
		response, _ := s.Check(context.Background(), tc.Request)
		if response.Status != tc.ExpectedStatus {
			t.Fatalf("Case %v failed. ExpectedStatus %v, got %v", tc.Name, tc.ExpectedStatus, response.Status)
		}
	}
}
