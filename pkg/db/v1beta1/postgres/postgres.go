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
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"

	v1beta1 "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	"github.com/kubeflow/katib/pkg/db/v1beta1/common"
	"github.com/kubeflow/katib/pkg/util/v1beta1/env"
	"k8s.io/klog"
)

const dbDriver = "postgres"

type dbConn struct {
	db *sql.DB
}

func getDbName() string {
	dbPassEnvName := common.DBPasswordEnvName
	dbPass := os.Getenv(dbPassEnvName)

	dbUser := env.GetEnvOrDefault(
		common.DBUserEnvName, common.DefaultPostgreSQLUser)
	dbHost := env.GetEnvOrDefault(
		common.PostgreSQLDBHostEnvName, common.DefaultPostgreSQLHost)
	dbPort := env.GetEnvOrDefault(
		common.PostgreSQLDBPortEnvName, common.DefaultPostgreSQLPort)
	dbName := env.GetEnvOrDefault(common.PostgreSQLDatabase,
		common.DefaultPostgreSQLDatabase)
	sslMode := env.GetEnvOrDefault(common.PostgreSSLMode,
		common.DefaultPostgreSSLMode)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPass, dbName, sslMode)

	return psqlInfo
}

func NewDBInterface(connectTimeout time.Duration) (common.KatibDBInterface, error) {
	db, err := common.OpenSQLConn(dbDriver, getDbName(), common.ConnectInterval, connectTimeout)
	if err != nil {
		return nil, fmt.Errorf("DB open failed: %v", err)
	}
	return &dbConn{db: db}, nil
}

func (d *dbConn) RegisterObservationLog(trialName string, observationLog *v1beta1.ObservationLog) error {
	statement := "INSERT INTO observation_logs (trial_name, time, metric_name, value) VALUES "
	values := []interface{}{}

	index_of_qparam := 1
	for _, mlog := range observationLog.MetricLogs {
		if mlog.TimeStamp == "" {
			continue
		}
		t, err := time.Parse(time.RFC3339Nano, mlog.TimeStamp)
		if err != nil {
			return fmt.Errorf("Error parsing start time %s: %v", mlog.TimeStamp, err)
		}
		sqlTimeStr := t.UTC().Format(time.RFC3339Nano)

		statement += fmt.Sprintf("($%d, $%d, $%d, $%d),",
			index_of_qparam, index_of_qparam+1, index_of_qparam+2, index_of_qparam+3,
		)
		values = append(values, trialName, sqlTimeStr, mlog.Metric.Name, mlog.Metric.Value)
		index_of_qparam += 4
	}

	statement = statement[:len(statement)-1]

	// Prepare the statement
	stmt, err := d.db.Prepare(statement)
	if err != nil {
		return fmt.Errorf("Prepare SQL statement failed: %v", err)
	}

	// Defer Close the statement
	defer stmt.Close()

	// Execute INSERT
	_, err = stmt.Exec(values...)
	if err != nil {
		return fmt.Errorf("Execute SQL INSERT failed: %v", err)
	}

	return nil
}

func (d *dbConn) GetObservationLog(trialName string, metricName string, startTime string, endTime string) (*v1beta1.ObservationLog, error) {
	qfield := []interface{}{trialName}
	qstr := ""
	index_of_qparam := 1

	base_stmt := fmt.Sprintf("SELECT time, metric_name, value FROM observation_logs WHERE trial_name = $%d", index_of_qparam)
	index_of_qparam += 1

	if metricName != "" {
		qstr += fmt.Sprintf(" AND metric_name = $%d", index_of_qparam)
		qfield = append(qfield, metricName)
		index_of_qparam += 1
	}

	if startTime != "" {
		s_time, err := time.Parse(time.RFC3339Nano, startTime)
		if err != nil {
			return nil, fmt.Errorf("Error parsing start time %s: %v", startTime, err)
		}
		formattedStartTime := s_time.UTC().Format(time.RFC3339Nano)
		qstr += fmt.Sprintf(" AND time >= $%d", index_of_qparam)
		qfield = append(qfield, formattedStartTime)
		index_of_qparam += 1
	}
	if endTime != "" {
		e_time, err := time.Parse(time.RFC3339Nano, endTime)
		if err != nil {
			return nil, fmt.Errorf("Error parsing completion time %s: %v", endTime, err)
		}
		formattedEndTime := e_time.UTC().Format(time.RFC3339Nano)
		qstr += fmt.Sprintf(" AND time <= $%d", index_of_qparam)
		qfield = append(qfield, formattedEndTime)
		// index_of_qparam += 1  // if any other filters are added, this should be incremented
	}

	rows, err := d.db.Query(base_stmt+qstr+" ORDER BY time", qfield...)
	if err != nil {
		return nil, fmt.Errorf("Failed to get ObservationLogs %v", err)
	}

	// Defer Close the rows
	defer rows.Close()

	result := &v1beta1.ObservationLog{
		MetricLogs: []*v1beta1.MetricLog{},
	}
	for rows.Next() {
		var mname, mvalue, sqlTimeStr string
		err := rows.Scan(&sqlTimeStr, &mname, &mvalue)
		if err != nil {
			klog.Errorf("Error scanning log: %v", err)
			continue
		}
		ptime, err := time.Parse(time.RFC3339Nano, sqlTimeStr)
		if err != nil {
			klog.Errorf("Error parsing time %s: %v", sqlTimeStr, err)
			continue
		}
		timeStamp := ptime.UTC().Format(time.RFC3339Nano)
		result.MetricLogs = append(result.MetricLogs, &v1beta1.MetricLog{
			TimeStamp: timeStamp,
			Metric: &v1beta1.Metric{
				Name:  mname,
				Value: mvalue,
			},
		})
	}

	return result, nil
}

func (d *dbConn) DeleteObservationLog(trialName string) error {
	_, err := d.db.Exec("DELETE FROM observation_logs WHERE trial_name = $1", trialName)

	return err
}
