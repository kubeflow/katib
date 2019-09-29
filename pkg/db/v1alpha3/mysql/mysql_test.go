package mysql

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
	"github.com/kubeflow/katib/pkg/db/v1alpha3/common"
)

var dbInterface, mysqlInterface common.KatibDBInterface
var mock sqlmock.Sqlmock

var observationLogsColumns = []string{
	"trial_name",
	"id",
	"time",
	"metric_name",
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
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS observation_logs").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	dbInterface.DBInit()
	err = dbInterface.SelectOne()
	if err != nil {
		fmt.Printf("error `SELECT 1` probing: %v\n", err)
	}
	os.Exit(m.Run())
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
		"",
		"2016-12-31T21:01:05.123456Z",
		"",
	)
	if err != nil {
		t.Errorf("GetObservationLog failed %v", err)
	} else if len(obsLog.MetricLogs) != 2 {
		t.Errorf("GetObservationLog incorrect return %v", obsLog)
	}

}
