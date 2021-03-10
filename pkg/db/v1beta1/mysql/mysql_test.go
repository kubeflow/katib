package mysql

import (
	"fmt"
	"os"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"

	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	"github.com/kubeflow/katib/pkg/db/v1beta1/common"
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
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS observation_logs").WillReturnResult(sqlmock.NewResult(1, 1))
	dbInterface.DBInit()
	mock.ExpectExec("SELECT 1").
		WillReturnResult(sqlmock.NewResult(1, 1))
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
		},
	}
	mock.ExpectPrepare("INSERT")
	mock.ExpectExec(
		"INSERT",
	).WithArgs(
		"test1_trial1",
		"2016-12-31 20:02:05.123456",
		"f1_score",
		"88.95",
		"test1_trial1",
		"2016-12-31 20:02:05.123456",
		"loss",
		"0.5",
	).WillReturnResult(sqlmock.NewResult(1, 1))

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
		"loss",
		"2016-12-31T21:01:05.123456Z",
		"2016-12-31T22:10:20.123456Z",
	)
	if err != nil {
		t.Errorf("GetObservationLog failed %v", err)
	} else if len(obsLog.MetricLogs) != 2 {
		t.Errorf("GetObservationLog incorrect return %v", obsLog)
	}

}

func TestDeleteObservationLog(t *testing.T) {
	trialName := "test1_trial1"

	mock.ExpectExec(
		"DELETE FROM observation_logs",
	).WithArgs(trialName).WillReturnResult(sqlmock.NewResult(1, 1))

	err := dbInterface.DeleteObservationLog(trialName)
	if err != nil {
		t.Errorf("DeleteObservationLog failed: %v", err)
	}
}

func TestGetDbName(t *testing.T) {
	dsn := "root@tcp(katib-mysql:3306)/katib?timeout=5s&maxAllowedPacket=0"
	dataSourceName, _ := initDB()
	if dataSourceName != dsn {
		t.Errorf("initDB returns wrong value %v", dataSourceName)
	}

}
