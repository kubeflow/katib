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

package postgres

import (
	"fmt"
	"os"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	_ "github.com/lib/pq"

	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	"github.com/kubeflow/katib/pkg/db/v1beta1/common"
)

var dbInterface common.KatibDBInterface
var mock sqlmock.Sqlmock

func TestMain(m *testing.M) {
	db, sm, err := sqlmock.New()
	mock = sm
	if err != nil {
		fmt.Printf("error opening db: %v\n", err)
		os.Exit(1)
	}
	dbInterface = &dbConn{db: db}
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS observation_logs").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	dbInterface.DBInit()
	mock.ExpectExec("SELECT 1").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	err = dbInterface.SelectOne()
	if err != nil {
		fmt.Printf("error `SELECT 1` probing: %v\n", err)
	}
	m.Run()
}

func TestRegisterObservationLog(t *testing.T) {
	obsLog := &api_pb.ObservationLog{
		MetricLogs: []*api_pb.MetricLog{
			{
				TimeStamp: "2016-12-31T20:01:05.123456Z",
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
		"2016-12-31T20:01:05.123456Z",
		"f1_score",
		"88.95",
		"test1_trial1",
		"2016-12-31T20:02:05.123456Z",
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
			"2016-12-31T20:01:05.123456Z",
			"loss",
			"0.9",
		).AddRow(
			"2016-12-31T20:02:05.123456Z",
			"loss",
			"0.9",
		),
	)
	obsLog, err := dbInterface.GetObservationLog(
		"test1_trial1",
		"loss",
		"2016-12-31T20:01:05.123456Z",
		"2016-12-31T20:02:05.123456Z",
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
	cases := map[string]struct {
		updateEnvs map[string]string
		wantName   string
	}{
		"All parameters are default": {
			wantName: "host=katib-postgres port=5432 user=katib password= dbname=katib sslmode=disable",
		},
		"Set DB_USER": {
			updateEnvs: map[string]string{
				common.DBUserEnvName: "testUser",
			},
			wantName: "host=katib-postgres port=5432 user=testUser password= dbname=katib sslmode=disable",
		},
		"Set KATIB_POSTGRESQL_DB_HOST": {
			updateEnvs: map[string]string{
				common.PostgreSQLDBHostEnvName: "testHost",
			},
			wantName: "host=testHost port=5432 user=katib password= dbname=katib sslmode=disable",
		},
		"Set KATIB_POSTGRESQL_DB_PORT": {
			updateEnvs: map[string]string{
				common.PostgreSQLDBPortEnvName: "1234",
			},
			wantName: "host=katib-postgres port=1234 user=katib password= dbname=katib sslmode=disable",
		},
		"Set KATIB_POSTGRESQL_DB_DATABASE": {
			updateEnvs: map[string]string{
				common.PostgreSQLDatabase: "testDB",
			},
			wantName: "host=katib-postgres port=5432 user=katib password= dbname=testDB sslmode=disable",
		},
		"Set DB_PASSWORD": {
			updateEnvs: map[string]string{
				common.DBPasswordEnvName: "testPassword",
			},
			wantName: "host=katib-postgres port=5432 user=katib password=testPassword dbname=katib sslmode=disable",
		},
		"Set KATIB_POSTGRESQL_SSL_MODE": {
			updateEnvs: map[string]string{
				common.PostgreSSLMode: "require",
			},
			wantName: "host=katib-postgres port=5432 user=katib password= dbname=katib sslmode=require",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			for k, v := range tc.updateEnvs {
				t.Setenv(k, v)
			}
			gotName := getDbName()
			if diff := cmp.Diff(tc.wantName, gotName); len(diff) != 0 {
				t.Errorf("Unexpected DBName (-want,+got):\n%s", diff)
			}
		})
	}
}
