package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	health_pb "github.com/kubeflow/katib/pkg/apis/manager/health"
	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
	"github.com/kubeflow/katib/pkg/db/v1alpha3/common"
	mockdb "github.com/kubeflow/katib/pkg/mock/v1alpha3/db"
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
	req := &health_pb.HealthCheckRequest{
		Service: "grpc.health.v1.Health",
	}
	servingResp := &health_pb.HealthCheckResponse{
		Status: health_pb.HealthCheckResponse_SERVING,
	}

	mockDB.EXPECT().SelectOne().Return(nil)

	resp, err := s.Check(context.Background(), req)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	if resp.Status != servingResp.Status {
		t.Fatalf("Check must return serving status, but returned %v", resp.Status)
	}
}

var mock sqlmock.Sqlmock

func TestMain(m *testing.M) {

	dbNameEnvName := common.DBNameEnvName
	os.Setenv(dbNameEnvName, "mysql")
	_, sm, err := sqlmock.New()
	mock = sm
	if err != nil {
		fmt.Printf("error opening db: %v\n", err)
		os.Exit(1)
	}
	main()
	os.Exit(m.Run())
}
