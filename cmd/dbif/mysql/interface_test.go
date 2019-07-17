package main

import (
	//	"database/sql"
	//	"database/sql/driver"
	//	"errors"
	"testing"
	"fmt"
	"os"
	//"time"
	"context"

	_ "github.com/go-sql-driver/mysql"
	//	"github.com/golang/protobuf/jsonpb"
	api_pb "github.com/kubeflow/katib/pkg/api/v1alpha2/dbif"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var dbInterface *dbConn
var mock sqlmock.Sqlmock

var experimentColums = []string{
	"id",
	"name",
	"parameters",
	"objective",
	"algorithm",
	"trial_template",
	"metrics_collector_spec",
	"parallel_trial_count",
	"max_trial_count",
	"status",
	"start_time",
	"completion_time",
	"nas_config",
}
var trialColumns = []string{
	"id",
	"name",
	"experiment_name",
	"objective",
	"parameter_assignments",
	"run_spec",
	"metrics_collector_spec",
	"observation",
	"status",
	"start_time",
	"completion_time",
}

var observationLogsColumns = []string{
	"trial_name",
	"id",
	"time",
	"metric_name",
	"value",
}

var extraAlgorithmSettingsColumns = []string{
	"experiment_name",
	"id",
	"setting_name",
	"value",
}

func TestMain(m *testing.M) {
	db, sm, err := sqlmock.New()
	mock = sm
	if err != nil {
		fmt.Printf("error opening db: %v\n", err)
		os.Exit(1)
	}
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS experiments").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS trials").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS observation_logs").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS extra_algorithm_settings").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	dbInterface, err = NewWithSQLConn(db)
	if err != nil {
		fmt.Printf("error NewWithSQLConn: %v\n", err)
	}
	_, err = dbInterface.SelectOne(context.Background(), &api_pb.SelectOneRequest{})
	if err != nil {
		fmt.Printf("error `SELECT 1` probing: %v\n", err)
	}
	os.Exit(m.Run())
}

func TestRegisterExperiment(t *testing.T) {
	experiment := &api_pb.Experiment{
		Name: "test1",
		Spec: &api_pb.ExperimentSpec{
			ParameterSpecs: &api_pb.ExperimentSpec_ParameterSpecs{
				Parameters: []*api_pb.ParameterSpec{},
			},
			Objective: &api_pb.ObjectiveSpec{
				Type:                  api_pb.ObjectiveType_UNKNOWN,
				Goal:                  0.99,
				ObjectiveMetricName:   "f1_score",
				AdditionalMetricNames: []string{"loss", "precision", "recall"},
			},
			Algorithm:          &api_pb.AlgorithmSpec{},
			TrialTemplate:      "",
			ParallelTrialCount: 10,
			MaxTrialCount:      100,
		},
		Status: &api_pb.ExperimentStatus{
			Condition:      api_pb.ExperimentStatus_CREATED,
			StartTime:      "2016-12-31T20:02:05.123456Z",
			CompletionTime: "2016-12-31T20:02:06.123456Z",
		},
	}
	mock.ExpectExec(
		`INSERT INTO experiments \(
			name, 
			parameters, 
			objective, 
			algorithm, 
			trial_template,
			metrics_collector_spec,
			parallel_trial_count, 
			max_trial_count,
			status,
			start_time,
			completion_time,
			nas_config\)`,
	).WithArgs(
		"test1",
		"{\"parameters\":[]}",
		"{\"goal\":0.99,\"objectiveMetricName\":\"f1_score\",\"additionalMetricNames\":[\"loss\",\"precision\",\"recall\"]}",
		"{}",
		experiment.Spec.TrialTemplate,
		experiment.Spec.MetricsCollectorSpec,
		experiment.Spec.ParallelTrialCount,
		experiment.Spec.MaxTrialCount,
		experiment.Status.Condition,
		"2016-12-31 20:02:05.123456",
		"2016-12-31 20:02:06.123456",
		"",
	).WillReturnResult(sqlmock.NewResult(1, 1))
	_, err := dbInterface.RegisterExperiment(context.Background(), &api_pb.RegisterExperimentRequest{Experiment: experiment})
	if err != nil {
		t.Errorf("RegisterExperiment failed: %v", err)
	}
}

func TestGetExperiment(t *testing.T) {
	mock.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows(experimentColums).AddRow(
			1,
			"test1",
			"{\"parameters\":[]}",
			"{\"goal\":0.99,\"objectiveMetricName\":\"f1_score\",\"additionalMetricNames\":[\"loss\",\"precision\",\"recall\"]}",
			"{}",
			"",
			"",
			10,
			100,
			api_pb.ExperimentStatus_CREATED,
			"2016-12-31 20:02:05.123456",
			"2016-12-31 20:02:06.123456",
			"",
		),
	)
	response, err := dbInterface.GetExperiment(context.Background(), &api_pb.GetExperimentRequest{ExperimentName: "test1"})
	if err != nil {
		t.Errorf("GetExperiment failed %v", err)
	} else 
	{
		experiment := response.Experiment
		if experiment.Name != "test1" {
		t.Errorf("GetExperiment incorrect return %v", experiment)
		}
	}
}

func TestGetExperimentList(t *testing.T) {
	mock.ExpectQuery("SELECT  name, status, start_time, completion_time FROM experiments").WillReturnRows(
		sqlmock.NewRows([]string{"name", "status", "start_time", "completion_time"}).AddRow(
			"test1",
			api_pb.ExperimentStatus_CREATED,
			"2016-12-31 20:02:05.123456",
			"2016-12-31 20:02:06.123456",
		).AddRow(
			"test2",
			api_pb.ExperimentStatus_SUCCEEDED,
			"2016-12-31 20:02:05.123456",
			"2016-12-31 20:05:05.123456",
		),
	)
	response, err := dbInterface.GetExperimentList(context.Background(), &api_pb.GetExperimentListRequest{})
	experiments := response.ExperimentSummaries
	if err != nil {
		t.Errorf("GetExperimentList failed %v", err)
	} else if len(experiments) != 2 {
		t.Errorf("Wrong Experiment number %d", len(experiments))
	} else if experiments[0].ExperimentName != "test1" || experiments[1].ExperimentName != "test2" {
		t.Errorf("GetExperimentList incorrect return %v", experiments)
	}
}

func TestUpdateExperimentStatus(t *testing.T) {
	condition := api_pb.ExperimentStatus_RUNNING
	exp_name := "test1"
	start_time := "2016-12-31 20:02:05.123456"
	completion_time := "2016-12-31 20:02:06.123456"

	mock.ExpectExec(`UPDATE experiments SET status = \?,
	start_time = \?,
	completion_time = \? WHERE name = \?`,
	).WithArgs(condition, start_time, completion_time, exp_name).WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := dbInterface.UpdateExperimentStatus(context.Background(), 
	&api_pb.UpdateExperimentStatusRequest{ExperimentName: exp_name,
		NewStatus: &api_pb.ExperimentStatus{
			Condition:      condition,
			StartTime:      "2016-12-31T20:02:05.123456Z",
			CompletionTime: "2016-12-31T20:02:06.123456Z",
		},
	})
	if err != nil {
		t.Errorf("UpdateExperiment failed %v", err)
	}
}

func TestUpdateAlgorithmExtraSettings(t *testing.T) {
	//TestUpdateAlgorithmExtraSettings inclueds TestGetAlgorithmExtraSettings
	//settin1 is already stored and setting2 is not exist in DB.
	exp_name := "test1"
	exAlgoSet := []*api_pb.AlgorithmSetting{
		{
			Name:  "setting1",
			Value: "100",
		},
		{
			Name:  "setting2",
			Value: "0.2",
		},
	}
	mock.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"setting_name", "value"}).AddRow(
			"setting1",
			"50",
		),
	)

	mock.ExpectExec(`UPDATE extra_algorithm_settings SET value = \?
	WHERE experiment_name = \? AND setting_name = \?`,
	).WithArgs(exAlgoSet[0].Value, exp_name, exAlgoSet[0].Name).WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`INSERT INTO extra_algorithm_settings \(
		experiment_name,
		setting_name,
		value\)`,
	).WithArgs(exp_name, exAlgoSet[1].Name, exAlgoSet[1].Value).WillReturnResult(sqlmock.NewResult(1, 1))
	_, err := dbInterface.UpdateAlgorithmExtraSettings(context.Background(), 
	&api_pb.UpdateAlgorithmExtraSettingsRequest{
		ExperimentName: exp_name, 
		ExtraAlgorithmSettings: exAlgoSet,
	})
	if err != nil {
		t.Errorf("UpdateAlgorithmExtraSettings failed %v", err)
	}

}

func TestRegisterTrial(t *testing.T) {
	trial := &api_pb.Trial{
		Name: "test1_trial1",
		Spec: &api_pb.TrialSpec{
			ExperimentName: "test1",
			Objective: &api_pb.ObjectiveSpec{
				Type:                  api_pb.ObjectiveType_UNKNOWN,
				Goal:                  0.99,
				ObjectiveMetricName:   "f1_score",
				AdditionalMetricNames: []string{"loss", "precision", "recall"},
			},
			ParameterAssignments: &api_pb.TrialSpec_ParameterAssignments{
				Assignments: []*api_pb.ParameterAssignment{
					{
						Name:  "param1",
						Value: "0.9",
					},
					{
						Name:  "param2",
						Value: "10",
					},
				},
			},
			RunSpec: "",
		},
		Status: &api_pb.TrialStatus{
			Condition: api_pb.TrialStatus_RUNNING,
			Observation: &api_pb.Observation{
				Metrics: []*api_pb.Metric{
					{
						Name:  "f1_score",
						Value: "88.95",
					},
					{
						Name:  "loss",
						Value: "0.5",
					},
					{
						Name:  "precision",
						Value: "88.7",
					},
					{
						Name:  "recall",
						Value: "89.2",
					},
				},
			},
			StartTime:      "2016-12-31T20:02:05.123456Z",
			CompletionTime: "2016-12-31T20:02:06.123456Z",
		},
	}
	mock.ExpectExec(
		`INSERT INTO trials \(
			name, 
			experiment_name,
			objective,
			parameter_assignments,
			run_spec,
			metrics_collector_spec,
			observation,
			status,
			start_time,
			completion_time\)`,
	).WithArgs(
		"test1_trial1",
		"test1",
		"{\"goal\":0.99,\"objectiveMetricName\":\"f1_score\",\"additionalMetricNames\":[\"loss\",\"precision\",\"recall\"]}",
		"{\"assignments\":[{\"name\":\"param1\",\"value\":\"0.9\"},{\"name\":\"param2\",\"value\":\"10\"}]}",
		"",
		"",
		"{\"metrics\":[{\"name\":\"f1_score\",\"value\":\"88.95\"},{\"name\":\"loss\",\"value\":\"0.5\"},{\"name\":\"precision\",\"value\":\"88.7\"},{\"name\":\"recall\",\"value\":\"89.2\"}]}",
		trial.Status.Condition,
		"2016-12-31 20:02:05.123456",
		"2016-12-31 20:02:06.123456",
	).WillReturnResult(sqlmock.NewResult(1, 1))
	_, err := dbInterface.RegisterTrial(context.Background(), &api_pb.RegisterTrialRequest{Trial: trial})
	if err != nil {
		t.Errorf("RegisterTrial failed: %v", err)
	}

}

func TestGetTrialList(t *testing.T) {
	mock.ExpectQuery(`SELECT`).WillReturnRows(
		sqlmock.NewRows(trialColumns).AddRow(
			1,
			"test1_trial1",
			"test1",
			"{\"goal\":0.99,\"objectiveMetricName\":\"f1_score\",\"additionalMetricNames\":[\"loss\",\"precision\",\"recall\"]}",
			"{\"assignments\":[{\"name\":\"param1\",\"value\":\"0.9\"},{\"name\":\"param2\",\"value\":\"10\"}]}",
			"",
			"",
			"{\"metrics\":[{\"name\":\"f1_score\",\"value\":\"88.95\"},{\"name\":\"loss\",\"value\":\"0.5\"},{\"name\":\"precision\",\"value\":\"88.7\"},{\"name\":\"recall\",\"value\":\"89.2\"}]}",
			api_pb.TrialStatus_RUNNING,
			"2016-12-31 20:02:05.123456",
			"2016-12-31 20:02:06.123456",
		).AddRow(
			2,
			"test1_trial2",
			"test1",
			"{\"goal\":0.99,\"objectiveMetricName\":\"f1_score\",\"additionalMetricNames\":[\"loss\",\"precision\",\"recall\"]}",
			"{\"assignments\":[{\"name\":\"param1\",\"value\":\"0.8\"},{\"name\":\"param2\",\"value\":\"20\"}]}",
			"",
			"",
			"{\"metrics\":[{\"name\":\"f1_score\",\"value\":\"88.5\"},{\"name\":\"loss\",\"value\":\"0.8\"},{\"name\":\"precision\",\"value\":\"88.2\"},{\"name\":\"recall\",\"value\":\"89.0\"}]}",
			api_pb.TrialStatus_SUCCEEDED,
			"2016-12-31 20:02:05.123456",
			"2016-12-31 20:02:06.123456",
		),
	)
	response, err := dbInterface.GetTrialList(context.Background(), &api_pb.GetTrialListRequest{
		ExperimentName: "test1", 
		Filter: "trial",
	})
	trials := response.Trials
	if err != nil {
		t.Errorf("GetTrialList failed %v", err)
	} else if len(trials) != 2 {
		t.Errorf("Wrong trial number %d", len(trials))
	} else if trials[0].Name != "test1_trial1" || trials[1].Name != "test1_trial2" {
		t.Errorf("GetTrialList incorrect return %v", trials)
	}
}

func TestGetTrial(t *testing.T) {
	mock.ExpectQuery(`SELECT \* FROM trials WHERE name = \?`).WillReturnRows(
		sqlmock.NewRows(trialColumns).AddRow(
			1,
			"test1_trial1",
			"test1",
			"{\"goal\":0.99,\"objectiveMetricName\":\"f1_score\",\"additionalMetricNames\":[\"loss\",\"precision\",\"recall\"]}",
			"{\"assignments\":[{\"name\":\"param1\",\"value\":\"0.9\"},{\"name\":\"param2\",\"value\":\"10\"}]}",
			"",
			"",
			"{\"metrics\":[{\"name\":\"f1_score\",\"value\":\"88.95\"},{\"name\":\"loss\",\"value\":\"0.5\"},{\"name\":\"precision\",\"value\":\"88.7\"},{\"name\":\"recall\",\"value\":\"89.2\"}]}",
			api_pb.TrialStatus_RUNNING,
			"2016-12-31 20:02:05.123456",
			"2016-12-31 20:02:06.123456",
		),
	)
	response, err := dbInterface.GetTrial(context.Background(), &api_pb.GetTrialRequest{TrialName: "test1_trial1"})
	trial := response.Trial
	if err != nil {
		t.Errorf("GetTria failed %v", err)
	} else if trial.Name != "test1_trial1" {
		t.Errorf("GetTrial incorrect return %v", trial)
	}
}

func TestUpdateTrialStatus(t *testing.T) {
	condition := api_pb.TrialStatus_RUNNING
	trial_name := "test1_trial1"
	start_time := "2016-12-31 20:02:05.123456"
	completion_time := "2016-12-31 20:02:06.123456"

	mock.ExpectExec(`UPDATE trials SET status = \?,
	start_time = \?,
	completion_time = \?,
	observation = \? WHERE name = \?`,
	).WithArgs(
		condition,
		start_time,
		completion_time,
		"{\"metrics\":[{\"name\":\"f1_score\",\"value\":\"88.95\"},{\"name\":\"loss\",\"value\":\"0.5\"},{\"name\":\"precision\",\"value\":\"88.7\"},{\"name\":\"recall\",\"value\":\"89.2\"}]}",
		trial_name).WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := dbInterface.UpdateTrialStatus(context.Background(), &api_pb.UpdateTrialStatusRequest{
		TrialName: trial_name,
		NewStatus: &api_pb.TrialStatus{
			Condition:      condition,
			StartTime:      "2016-12-31T20:02:05.123456Z",
			CompletionTime: "2016-12-31T20:02:06.123456Z",
			Observation: &api_pb.Observation{
				Metrics: []*api_pb.Metric{
					{
						Name:  "f1_score",
						Value: "88.95",
					},
					{
						Name:  "loss",
						Value: "0.5",
					},
					{
						Name:  "precision",
						Value: "88.7",
					},
					{
						Name:  "recall",
						Value: "89.2",
					},
				},
			},
		},
	})
	if err != nil {
		t.Errorf("UpdateTrial failed %v", err)
	}
}

func TestRegisterObservationLog(t *testing.T) {
	obsLog := &api_pb.ObservationLog{
		MetricLogs: []*api_pb.MetricLog{
			{
				TimeStamp: "2016-12-31T20:02:05.123456Z",
				Metric: &api_pb.Metric{
					Name:  "f1_score",
					Value: "88.95",
				},
			},
			{
				TimeStamp: "2016-12-31T20:02:05.123456Z",
				Metric: &api_pb.Metric{
					Name:  "loss",
					Value: "0.5",
				},
			},
			{
				TimeStamp: "2016-12-31T20:02:05.123456Z",
				Metric: &api_pb.Metric{
					Name:  "precision",
					Value: "88.7",
				},
			},
			{
				TimeStamp: "2016-12-31T20:02:05.123456Z",
				Metric: &api_pb.Metric{
					Name:  "recall",
					Value: "89.2",
				},
			},
		},
	}
	for _, m := range obsLog.MetricLogs {
		mock.ExpectExec(
			`INSERT INTO observation_logs \(
				trial_name,
				time,
				metric_name,
				value
			\)`,
		).WithArgs(
			"test1_trial1",
			"2016-12-31 20:02:05.123456",
			m.Metric.Name,
			m.Metric.Value,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	}
	_, err := dbInterface.ReportObservationLog(context.Background(), 
			&api_pb.ReportObservationLogRequest{TrialName: "test1_trial1", ObservationLog: obsLog})
	if err != nil {
		t.Errorf("RegisterExperiment failed: %v", err)
	}

}

func TestGetObservationLog(t *testing.T) {
	mock.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"time", "metric_name", "value"}).AddRow(
			"2016-12-31 21:02:05.123456",
			"loss",
			"0.9",
		).AddRow(
			"2016-12-31 22:02:05.123456",
			"loss",
			"0.9",
		),
	)
	response, err := dbInterface.GetObservationLog(context.Background(), &api_pb.GetObservationLogRequest{
		TrialName: "test1_trial1",
		MetricName: "",
		StartTime: "2016-12-31T21:01:05.123456Z",
		EndTime: "",
	})
	obsLog := response.ObservationLog
	if err != nil {
		t.Errorf("GetObservationLog failed %v", err)
	} else if len(obsLog.MetricLogs) != 2 {
		t.Errorf("GetObservationLog incorrect return %v", obsLog)
	}

}
