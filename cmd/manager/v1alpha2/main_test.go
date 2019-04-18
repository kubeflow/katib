package main

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	api_pb "github.com/kubeflow/katib/pkg/api/v1alpha2"
	mockdb "github.com/kubeflow/katib/pkg/mock/v1alpha2/db"
)

func TestRegisterExperiment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := &server{}
	mockDB := mockdb.NewMockKatibDBInterface(ctrl)
	dbIf = mockDB

	req := &api_pb.RegisterExperimentRequest{
		Experiment: &api_pb.Experiment{
			Name: "testExp",
			ExperimentSpec: &api_pb.ExperimentSpec{
				ParameterSpecs: &api_pb.ExperimentSpec_ParameterSpecs{
					Parameters: []*api_pb.ParameterSpec{},
				},
				Objective: &api_pb.ObjectiveSpec{
					Type:                   api_pb.ObjectiveType_UNKNOWN,
					Goal:                   0.99,
					ObjectiveMetricName:    "f1_score",
					AdditionalMetricsNames: []string{"loss", "precision", "recall"},
				},
				Algorithm:          &api_pb.AlgorithmSpec{},
				TrialTemplate:      "",
				ParallelTrialCount: 10,
				MaxTrialCount:      100,
			},
			ExperimentStatus: &api_pb.ExperimentStatus{
				Condition:      api_pb.ExperimentStatus_CREATED,
				StartTime:      "2019-02-03T04:05:06+09:00",
				CompletionTime: "",
			},
		},
	}
	mockDB.EXPECT().RegisterExperiment(req.Experiment).Return(nil)
	_, err := s.RegisterExperiment(context.Background(), req)
	if err != nil {
		t.Fatalf("RegisterExperiment Error %v", err)
	}
}

func TestDeleteExperiment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := &server{}
	mockDB := mockdb.NewMockKatibDBInterface(ctrl)
	dbIf = mockDB

	req := &api_pb.DeleteExperimentRequest{
		ExperimentName: "testExp",
	}
	mockDB.EXPECT().DeleteExperiment(req.ExperimentName).Return(nil)
	_, err := s.DeleteExperiment(context.Background(), req)
	if err != nil {
		t.Fatalf("DeleteExperiment Error %v", err)
	}
}

func TestGetExperiment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := &server{}
	mockDB := mockdb.NewMockKatibDBInterface(ctrl)
	dbIf = mockDB
	req := &api_pb.GetExperimentRequest{
		ExperimentName: "testExp",
	}

	testExp := api_pb.Experiment{
		Name: "testExp",
		ExperimentSpec: &api_pb.ExperimentSpec{
			ParameterSpecs: &api_pb.ExperimentSpec_ParameterSpecs{
				Parameters: []*api_pb.ParameterSpec{},
			},
			Objective: &api_pb.ObjectiveSpec{
				Type:                   api_pb.ObjectiveType_UNKNOWN,
				Goal:                   0.99,
				ObjectiveMetricName:    "f1_score",
				AdditionalMetricsNames: []string{"loss", "precision", "recall"},
			},
			Algorithm:          &api_pb.AlgorithmSpec{},
			TrialTemplate:      "",
			ParallelTrialCount: 10,
			MaxTrialCount:      100,
		},
		ExperimentStatus: &api_pb.ExperimentStatus{
			Condition:      api_pb.ExperimentStatus_CREATED,
			StartTime:      "2019-02-03T04:05:06+09:00",
			CompletionTime: "",
		},
	}
	mockDB.EXPECT().GetExperiment(req.ExperimentName).Return(&testExp, nil)
	ret, err := s.GetExperiment(context.Background(), req)
	if ret.Experiment.Name != testExp.Name {
		t.Fatalf("GetExperiment Test fail expect experiment name %s got %s", ret.Experiment.Name, testExp.Name)
	}
	if err != nil {
		t.Fatalf("GetExperiment Error %v", err)
	}
}

func TestGetExperimentList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := &server{}
	mockDB := mockdb.NewMockKatibDBInterface(ctrl)
	dbIf = mockDB
	req := &api_pb.GetExperimentListRequest{}

	testExpList := []*api_pb.ExperimentSummary{
		&api_pb.ExperimentSummary{
			ExperimentName: "test1",
			Status: &api_pb.ExperimentStatus{
				Condition:      api_pb.ExperimentStatus_CREATED,
				StartTime:      "2019-02-03T04:05:06+09:00",
				CompletionTime: "",
			},
		},
		&api_pb.ExperimentSummary{
			ExperimentName: "test2",
			Status: &api_pb.ExperimentStatus{
				Condition:      api_pb.ExperimentStatus_SUCCEEDED,
				StartTime:      "2019-02-03T04:02:06+09:00",
				CompletionTime: "2019-02-03T04:03:06+09:00",
			},
		},
	}
	mockDB.EXPECT().GetExperimentList().Return(testExpList, nil)
	ret, err := s.GetExperimentList(context.Background(), req)
	if len(ret.ExperimentSummaries) != len(testExpList) {
		t.Fatalf("GetExperimentList Test fail expect experiment number %d got %d", len(testExpList), len(ret.ExperimentSummaries))
	}
	if err != nil {
		t.Fatalf("GetExperimentList Error %v", err)
	}
}

func TestUpdateExperimentStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := &server{}
	mockDB := mockdb.NewMockKatibDBInterface(ctrl)
	dbIf = mockDB

	req := &api_pb.UpdateExperimentStatusRequest{
		ExperimentName: "test1",
		NewStatus: &api_pb.ExperimentStatus{
			Condition:      api_pb.ExperimentStatus_RUNNING,
			StartTime:      "2019-02-03T04:05:06+09:00",
			CompletionTime: "",
		},
	}
	mockDB.EXPECT().UpdateExperimentStatus(req.ExperimentName, req.NewStatus).Return(nil)
	_, err := s.UpdateExperimentStatus(context.Background(), req)
	if err != nil {
		t.Fatalf("UpdateExperimentStatus Error %v", err)
	}
}

func TestUpdateAlgorithmExtraSettings(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := &server{}
	mockDB := mockdb.NewMockKatibDBInterface(ctrl)
	dbIf = mockDB

	req := &api_pb.UpdateAlgorithmExtraSettingsRequest{
		ExperimentName: "test1",
		ExtraAlgorithmSettings: []*api_pb.AlgorithmSetting{
			&api_pb.AlgorithmSetting{
				Name:  "set1",
				Value: "10",
			},
			&api_pb.AlgorithmSetting{
				Name:  "set2",
				Value: "0.5",
			},
		},
	}
	mockDB.EXPECT().UpdateAlgorithmExtraSettings(req.ExperimentName, req.ExtraAlgorithmSettings).Return(nil)
	_, err := s.UpdateAlgorithmExtraSettings(context.Background(), req)
	if err != nil {
		t.Fatalf("UpdateAlgorithmExtraSettings Error %v", err)
	}
}

func TestGetAlgorithmExtraSettings(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := &server{}
	mockDB := mockdb.NewMockKatibDBInterface(ctrl)
	dbIf = mockDB

	req := &api_pb.GetAlgorithmExtraSettingsRequest{
		ExperimentName: "test1",
	}
	extraAlgoSets := []*api_pb.AlgorithmSetting{
		&api_pb.AlgorithmSetting{
			Name:  "set1",
			Value: "10",
		},
		&api_pb.AlgorithmSetting{
			Name:  "set2",
			Value: "0.5",
		},
	}
	mockDB.EXPECT().GetAlgorithmExtraSettings(req.ExperimentName).Return(extraAlgoSets, nil)
	ret, err := s.GetAlgorithmExtraSettings(context.Background(), req)
	if len(ret.ExtraAlgorithmSettings) != len(extraAlgoSets) {
		t.Fatalf("GetAlgorithmExtraSettings Test fail expect experiment number %d got %d", len(extraAlgoSets), len(ret.ExtraAlgorithmSettings))
	}
	if err != nil {
		t.Fatalf("UpdateAlgorithmExtraSettings Error %v", err)
	}
}

func TestRegisterTrial(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := &server{}
	mockDB := mockdb.NewMockKatibDBInterface(ctrl)
	dbIf = mockDB

	req := &api_pb.RegisterTrialRequest{
		Trial: &api_pb.Trial{
			Name: "test1-trial1",
			Spec: &api_pb.TrialSpec{
				ExperimentName: "test1",
				RunSpec:        "",
				ParameterAssignments: &api_pb.TrialSpec_ParameterAssignments{
					Assignments: []*api_pb.ParameterAssignment{
						&api_pb.ParameterAssignment{
							Name:  "param1",
							Value: "10",
						},
						&api_pb.ParameterAssignment{
							Name:  "param2",
							Value: "0.1",
						},
					},
				},
			},
			Status: &api_pb.TrialStatus{
				Condition:      api_pb.TrialStatus_RUNNING,
				StartTime:      "2019-02-03T04:05:06+09:00",
				CompletionTime: "",
				Observation: &api_pb.Observation{
					Metrics: []*api_pb.Metric{
						&api_pb.Metric{
							Name:  "f1_score",
							Value: "88.95",
						},
						&api_pb.Metric{
							Name:  "loss",
							Value: "0.5",
						},
						&api_pb.Metric{
							Name:  "precision",
							Value: "88.7",
						},
						&api_pb.Metric{
							Name:  "recall",
							Value: "89.2",
						},
					},
				},
			},
		},
	}
	mockDB.EXPECT().RegisterTrial(req.Trial).Return(nil)
	_, err := s.RegisterTrial(context.Background(), req)
	if err != nil {
		t.Fatalf("RegisterTrial Error %v", err)
	}
}

func TestDeleteTrial(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := &server{}
	mockDB := mockdb.NewMockKatibDBInterface(ctrl)
	dbIf = mockDB

	req := &api_pb.DeleteTrialRequest{
		TrialName: "test1-trial1",
	}
	mockDB.EXPECT().DeleteTrial(req.TrialName).Return(nil)
	_, err := s.DeleteTrial(context.Background(), req)
	if err != nil {
		t.Fatalf("DeleteTrial Error %v", err)
	}
}

func TestGetTrialList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := &server{}
	mockDB := mockdb.NewMockKatibDBInterface(ctrl)
	dbIf = mockDB

	req := &api_pb.GetTrialListRequest{
		ExperimentName: "test1",
	}
	trialList := []*api_pb.Trial{
		&api_pb.Trial{
			Name: "test1-trial1",
			Spec: &api_pb.TrialSpec{
				ExperimentName: "test1",
				RunSpec:        "",
				ParameterAssignments: &api_pb.TrialSpec_ParameterAssignments{
					Assignments: []*api_pb.ParameterAssignment{
						&api_pb.ParameterAssignment{
							Name:  "param1",
							Value: "10",
						},
						&api_pb.ParameterAssignment{
							Name:  "param2",
							Value: "0.1",
						},
					},
				},
			},
			Status: &api_pb.TrialStatus{
				Condition:      api_pb.TrialStatus_RUNNING,
				StartTime:      "2019-02-03T04:05:06+09:00",
				CompletionTime: "",
				Observation: &api_pb.Observation{
					Metrics: []*api_pb.Metric{
						&api_pb.Metric{
							Name:  "f1_score",
							Value: "88.95",
						},
						&api_pb.Metric{
							Name:  "loss",
							Value: "0.5",
						},
						&api_pb.Metric{
							Name:  "precision",
							Value: "88.7",
						},
						&api_pb.Metric{
							Name:  "recall",
							Value: "89.2",
						},
					},
				},
			},
		},
		&api_pb.Trial{
			Name: "test1-trial2",
			Spec: &api_pb.TrialSpec{
				ExperimentName: "test1",
				RunSpec:        "",
				ParameterAssignments: &api_pb.TrialSpec_ParameterAssignments{
					Assignments: []*api_pb.ParameterAssignment{
						&api_pb.ParameterAssignment{
							Name:  "param1",
							Value: "20",
						},
						&api_pb.ParameterAssignment{
							Name:  "param2",
							Value: "0.5",
						},
					},
				},
			},
			Status: &api_pb.TrialStatus{
				Condition:      api_pb.TrialStatus_PENDING,
				StartTime:      "",
				CompletionTime: "",
				Observation: &api_pb.Observation{
					Metrics: []*api_pb.Metric{},
				},
			},
		},
	}
	mockDB.EXPECT().GetTrialList(req.ExperimentName).Return(trialList, nil)
	ret, err := s.GetTrialList(context.Background(), req)
	if len(ret.Trials) != len(trialList) {
		t.Fatalf("GetTrialList Test fail expect tiral number %d got %d", len(trialList), len(ret.Trials))
	}
	if err != nil {
		t.Fatalf("GetTrialList Error %v", err)
	}

}

func TestGetTrial(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := &server{}
	mockDB := mockdb.NewMockKatibDBInterface(ctrl)
	dbIf = mockDB
	req := &api_pb.GetTrialRequest{
		TrialName: "test1-trial1",
	}

	trial := &api_pb.Trial{
		Name: "test1-trial1",
		Spec: &api_pb.TrialSpec{
			ExperimentName: "test1",
			RunSpec:        "",
			ParameterAssignments: &api_pb.TrialSpec_ParameterAssignments{
				Assignments: []*api_pb.ParameterAssignment{
					&api_pb.ParameterAssignment{
						Name:  "param1",
						Value: "10",
					},
					&api_pb.ParameterAssignment{
						Name:  "param2",
						Value: "0.1",
					},
				},
			},
		},
		Status: &api_pb.TrialStatus{
			Condition:      api_pb.TrialStatus_RUNNING,
			StartTime:      "2019-02-03T04:05:06+09:00",
			CompletionTime: "",
			Observation: &api_pb.Observation{
				Metrics: []*api_pb.Metric{
					&api_pb.Metric{
						Name:  "f1_score",
						Value: "88.95",
					},
					&api_pb.Metric{
						Name:  "loss",
						Value: "0.5",
					},
					&api_pb.Metric{
						Name:  "precision",
						Value: "88.7",
					},
					&api_pb.Metric{
						Name:  "recall",
						Value: "89.2",
					},
				},
			},
		},
	}
	mockDB.EXPECT().GetTrial(req.TrialName).Return(trial, nil)
	ret, err := s.GetTrial(context.Background(), req)
	if ret.Trial.Name != trial.Name {
		t.Fatalf("GetTrial Test fail expect tiral %s got %s", trial.Name, ret.Trial.Name)
	}
	if err != nil {
		t.Fatalf("GetTrial Error %v", err)
	}
}

func TestUpdateTrialStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := &server{}
	mockDB := mockdb.NewMockKatibDBInterface(ctrl)
	dbIf = mockDB

	req := &api_pb.UpdateTrialStatusRequest{
		TrialName: "test1-trial1",
		NewStatus: &api_pb.TrialStatus{
			Condition:      api_pb.TrialStatus_COMPLETED,
			StartTime:      "2019-02-03T04:05:06+09:00",
			CompletionTime: "2019-02-03T05:05:06+09:00",
			Observation: &api_pb.Observation{
				Metrics: []*api_pb.Metric{
					&api_pb.Metric{
						Name:  "f1_score",
						Value: "88.95",
					},
					&api_pb.Metric{
						Name:  "loss",
						Value: "0.5",
					},
					&api_pb.Metric{
						Name:  "precision",
						Value: "88.7",
					},
					&api_pb.Metric{
						Name:  "recall",
						Value: "89.2",
					},
				},
			},
		},
	}
	mockDB.EXPECT().UpdateTrialStatus(req.TrialName, req.NewStatus).Return(nil)
	_, err := s.UpdateTrialStatus(context.Background(), req)
	if err != nil {
		t.Fatalf("UpdateTrialStatus Error %v", err)
	}
}

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
				&api_pb.MetricLog{
					TimeStamp: "2019-02-03T04:05:06+09:00",
					Metric: &api_pb.Metric{
						Name:  "f1_score",
						Value: "88.95",
					},
				},
				&api_pb.MetricLog{
					TimeStamp: "2019-02-03T04:05:06+09:00",
					Metric: &api_pb.Metric{
						Name:  "loss",
						Value: "0.5",
					},
				},
				&api_pb.MetricLog{
					TimeStamp: "2019-02-03T04:05:06+09:00",
					Metric: &api_pb.Metric{
						Name:  "precision",
						Value: "88.7",
					},
				},
				&api_pb.MetricLog{
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
			&api_pb.MetricLog{
				TimeStamp: "2019-02-03T04:05:06+09:00",
				Metric: &api_pb.Metric{
					Name:  "f1_score",
					Value: "88.95",
				},
			},
			&api_pb.MetricLog{
				TimeStamp: "2019-02-03T04:05:06+09:00",
				Metric: &api_pb.Metric{
					Name:  "loss",
					Value: "0.5",
				},
			},
			&api_pb.MetricLog{
				TimeStamp: "2019-02-03T04:05:06+09:00",
				Metric: &api_pb.Metric{
					Name:  "precision",
					Value: "88.7",
				},
			},
			&api_pb.MetricLog{
				TimeStamp: "2019-02-03T04:05:06+09:00",
				Metric: &api_pb.Metric{
					Name:  "recall",
					Value: "89.2",
				},
			},
		},
	}

	mockDB.EXPECT().GetObservationLog(req.TrialName, req.StartTime, req.EndTime).Return(obs, nil)
	ret, err := s.GetObservationLog(context.Background(), req)
	if err != nil {
		t.Fatalf("GetObservationLog Error %v", err)
	}
	if len(obs.MetricLogs) != len(ret.ObservationLog.MetricLogs) {
		t.Fatalf("GetObservationLog Test fail expect metrics number %d got %d", len(obs.MetricLogs), len(ret.ObservationLog.MetricLogs))
	}
}
func TestGetSuggestions(t *testing.T) {
}
func TestValidateAlgorithmSettings(t *testing.T) {
}
