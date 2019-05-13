package db

import (
	//	"database/sql"
	//	"database/sql/driver"
	//	"errors"
	"fmt"
	"os"
	"testing"
	//"time"

	_ "github.com/go-sql-driver/mysql"
	//	"github.com/golang/protobuf/jsonpb"

	api_pb "github.com/kubeflow/katib/pkg/api/v1alpha2"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var dbInterface, mysqlInterface KatibDBInterface
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
	dbInterface, err = NewWithSQLConn(db)
	if err != nil {
		fmt.Printf("error NewWithSQLConn: %v\n", err)
	}
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS experiments").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS trials").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS observation_logs").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS extra_algorithm_settings").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	dbInterface.DBInit()
	err = dbInterface.SelectOne()
	if err != nil {
		fmt.Printf("error `SELECT 1` probing: %v\n", err)
	}
	os.Exit(m.Run())
}

func TestRegisterExperiment(t *testing.T) {
	experiment := &api_pb.Experiment{
		Name: "test1",
		ExperimentSpec: &api_pb.ExperimentSpec{
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
		ExperimentStatus: &api_pb.ExperimentStatus{
			Condition:      api_pb.ExperimentStatus_CREATED,
			StartTime:      "2016-12-31T20:02:05.123456Z",
			CompletionTime: "",
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
		experiment.ExperimentSpec.TrialTemplate,
		experiment.ExperimentSpec.MetricsCollectorSpec,
		experiment.ExperimentSpec.ParallelTrialCount,
		experiment.ExperimentSpec.MaxTrialCount,
		experiment.ExperimentStatus.Condition,
		"2016-12-31 20:02:05.123456",
		"",
		"",
	).WillReturnResult(sqlmock.NewResult(1, 1))
	err := dbInterface.RegisterExperiment(experiment)
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
			"",
			"",
		),
	)
	experiment, err := dbInterface.GetExperiment("test1")
	if err != nil {
		t.Errorf("GetExperiment failed %v", err)
	} else if experiment.Name != "test1" {
		t.Errorf("GetExperiment incorrect return %v", experiment)
	}
}

func TestGetExperimentList(t *testing.T) {
	mock.ExpectQuery("SELECT  name, status, start_time, completion_time FROM experiments").WillReturnRows(
		sqlmock.NewRows([]string{"name", "status", "start_time", "completion_time"}).AddRow(
			"test1",
			api_pb.ExperimentStatus_CREATED,
			"2016-12-31 20:02:05.123456",
			"",
		).AddRow(
			"test2",
			api_pb.ExperimentStatus_SUCCEEDED,
			"2016-12-31 20:02:05.123456",
			"2016-12-31 20:05:05.123456",
		),
	)
	experiments, err := dbInterface.GetExperimentList()
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
	completion_time := ""

	mock.ExpectExec(`UPDATE experiments SET status = \?,
	start_time = \?,
	completion_time = \? WHERE name = \?`,
	).WithArgs(condition, start_time, completion_time, exp_name).WillReturnResult(sqlmock.NewResult(1, 1))

	err := dbInterface.UpdateExperimentStatus(exp_name,
		&api_pb.ExperimentStatus{
			Condition:      condition,
			StartTime:      "2016-12-31T20:02:05.123456Z",
			CompletionTime: "",
		},
	)
	if err != nil {
		t.Errorf("UpdateExperiment failed %v", err)
	}
}

func TestUpdateAlgorithmExtraSettings(t *testing.T) {
	//TestUpdateAlgorithmExtraSettings inclueds TestGetAlgorithmExtraSettings
	//settin1 is already stored and setting2 is not exist in DB.
	exp_name := "test1"
	exAlgoSet := []*api_pb.AlgorithmSetting{
		&api_pb.AlgorithmSetting{
			Name:  "setting1",
			Value: "100",
		},
		&api_pb.AlgorithmSetting{
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

	mock.ExpectExec(`UPDATE extra_algorithm_settings SET value = \? ,
	WHERE experiment_name = \? AND setting_name = \?`,
	).WithArgs(exAlgoSet[0].Value, exp_name, exAlgoSet[0].Name).WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`INSERT INTO extra_algorithm_settings \(
		experiment_name,
		setting_name,
		value\)`,
	).WithArgs(exp_name, exAlgoSet[1].Name, exAlgoSet[1].Value).WillReturnResult(sqlmock.NewResult(1, 1))
	err := dbInterface.UpdateAlgorithmExtraSettings(exp_name, exAlgoSet)
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
					&api_pb.ParameterAssignment{
						Name:  "param1",
						Value: "0.9",
					},
					&api_pb.ParameterAssignment{
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
			StartTime:      "2016-12-31T20:02:05.123456Z",
			CompletionTime: "",
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
		"",
	).WillReturnResult(sqlmock.NewResult(1, 1))
	err := dbInterface.RegisterTrial(trial)
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
			"",
		).AddRow(
			2,
			"test1_trial2",
			"test1",
			"{\"goal\":0.99,\"objectiveMetricName\":\"f1_score\",\"additionalMetricNames\":[\"loss\",\"precision\",\"recall\"]}",
			"{\"assignments\":[{\"name\":\"param1\",\"value\":\"0.8\"},{\"name\":\"param2\",\"value\":\"20\"}]}",
			"",
			"",
			"{\"metrics\":[{\"name\":\"f1_score\",\"value\":\"88.5\"},{\"name\":\"loss\",\"value\":\"0.8\"},{\"name\":\"precision\",\"value\":\"88.2\"},{\"name\":\"recall\",\"value\":\"89.0\"}]}",
			api_pb.TrialStatus_COMPLETED,
			"2016-12-31 20:02:05.123456",
			"",
		),
	)
	trials, err := dbInterface.GetTrialList("test1", "trial")
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
			"",
		),
	)
	trial, err := dbInterface.GetTrial("test1_trial1")
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
	completion_time := ""

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

	err := dbInterface.UpdateTrialStatus(
		trial_name,
		&api_pb.TrialStatus{
			Condition:      condition,
			StartTime:      "2016-12-31T20:02:05.123456Z",
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
	)
	if err != nil {
		t.Errorf("UpdateTrial failed %v", err)
	}
}

func TestRegisterObservationLog(t *testing.T) {
	obsLog := &api_pb.ObservationLog{
		MetricLogs: []*api_pb.MetricLog{
			&api_pb.MetricLog{
				TimeStamp: "2016-12-31T20:02:05.123456Z",
				Metric: &api_pb.Metric{
					Name:  "f1_score",
					Value: "88.95",
				},
			},
			&api_pb.MetricLog{
				TimeStamp: "2016-12-31T20:02:05.123456Z",
				Metric: &api_pb.Metric{
					Name:  "loss",
					Value: "0.5",
				},
			},
			&api_pb.MetricLog{
				TimeStamp: "2016-12-31T20:02:05.123456Z",
				Metric: &api_pb.Metric{
					Name:  "precision",
					Value: "88.7",
				},
			},
			&api_pb.MetricLog{
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
	err := dbInterface.RegisterObservationLog("test1_trial1", obsLog)
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
	obsLog, err := dbInterface.GetObservationLog(
		"test1_trial1",
		"2016-12-31T21:01:05.123456Z",
		"",
	)
	if err != nil {
		t.Errorf("GetObservationLog failed %v", err)
	} else if len(obsLog.MetricLogs) != 2 {
		t.Errorf("GetObservationLog incorrect return %v", obsLog)
	}

}
