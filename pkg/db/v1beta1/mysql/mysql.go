package mysql

import (
	crand "crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"k8s.io/klog"

	v1beta1 "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	"github.com/kubeflow/katib/pkg/db/v1beta1/common"
	"github.com/kubeflow/katib/pkg/util/v1beta1/env"
)

const (
	dbDriver     = "mysql"
	mysqlTimeFmt = "2006-01-02 15:04:05.999999"

	connectInterval = 5 * time.Second
	connectTimeout  = 60 * time.Second
)

type dbConn struct {
	db *sql.DB
}

func CreateMySQLConfig(user, password string, mysqlServiceHost string,
	mysqlServicePort string, dbName string, mysqlExtraParams map[string]string) *mysql.Config {

	params := map[string]string{}

	for k, v := range mysqlExtraParams {
		params[k] = v
	}

	return &mysql.Config{
		User:                 user,
		Passwd:               password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", mysqlServiceHost, mysqlServicePort),
		Params:               params,
		DBName:               dbName,
		AllowNativePasswords: true,
		Timeout:              connectInterval,
	}
}

func initDB() (string, error) {
	dbPass := os.Getenv(common.DBPasswordEnvName)
	dbUser := env.GetEnvOrDefault(
		common.DBUserEnvName, common.DefaultMySQLUser)
	dbHost := env.GetEnvOrDefault(
		common.DBHostEnvName, common.DefaultMySQLHost)
	dbPort := env.GetEnvOrDefault(
		common.DBPortEnvName, common.DefaultMySQLPort)
	dbName := env.GetEnvOrDefault(common.DBNameEnvName,
		common.DefaultMySQLDatabase)
	dbExtraParamsStr := env.GetEnvOrDefault(common.DBExtraParamsEnvName,
		common.DefaultMySQLExtraParams)

	dbExtraParams := make(map[string]string)
	err := json.Unmarshal([]byte(dbExtraParamsStr), &dbExtraParams)
	if err != nil {
		klog.Errorf("Failed to get dbExtraParams. Error: %v", err)
	}

	if dbPass != "" {
		mysqlConfig := CreateMySQLConfig(dbUser, dbPass, dbHost, dbPort, "", dbExtraParams)

		ticker := time.NewTicker(connectInterval)
		defer ticker.Stop()

		timeoutC := time.After(connectTimeout)
		for {
			select {
			case <-ticker.C:
				if db, err := sql.Open(dbDriver, mysqlConfig.FormatDSN()); err == nil {
					if _, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName)); err == nil {
						mysqlConfig.DBName = dbName
						return mysqlConfig.FormatDSN(), nil
					}
					return "", fmt.Errorf("Creation of Katib db '%s' failed: %v", dbName, err)
				} else {
					klog.Errorf("Open sql connection failed: %v", err)
				}
			case <-timeoutC:
				return "", fmt.Errorf("Timeout waiting for DB conn successfully opened")
			}
		}
	} else {
		mysqlConfig := CreateMySQLConfig(dbUser, dbPass, dbHost, dbPort, dbName, dbExtraParams)
		return mysqlConfig.FormatDSN(), nil
	}
}

func openSQLConn(driverName string, dataSourceName string, interval time.Duration,
	timeout time.Duration) (*sql.DB, error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	timeoutC := time.After(timeout)
	for {
		select {
		case <-ticker.C:
			if db, err := sql.Open(driverName, dataSourceName); err == nil {
				if err = db.Ping(); err == nil {
					return db, nil
				}
				klog.Errorf("Ping to Katib db failed: %v", err)
			} else {
				klog.Errorf("Open sql connection failed: %v", err)
			}
		case <-timeoutC:
			return nil, fmt.Errorf("Timeout waiting for DB conn successfully opened.")
		}
	}
}

func NewWithSQLConn(db *sql.DB) (common.KatibDBInterface, error) {
	d := new(dbConn)
	d.db = db
	seed, err := crand.Int(crand.Reader, big.NewInt(1<<63-1))
	if err != nil {
		return nil, fmt.Errorf("RNG initialization failed: %v", err)
	}
	// We can do the following instead, but it creates a locking issue
	//d.rng = rand.New(rand.NewSource(seed.Int64()))
	rand.Seed(seed.Int64())

	return d, nil
}

func NewDBInterface() (common.KatibDBInterface, error) {
	dataSourceName, err := initDB()
	if err != nil {
		return nil, fmt.Errorf("DB init failed: %v", err)
	}
	db, err := openSQLConn(dbDriver, dataSourceName, connectInterval, connectTimeout)
	if err != nil {
		return nil, fmt.Errorf("DB open failed: %v", err)
	}
	return NewWithSQLConn(db)
}

func (d *dbConn) RegisterObservationLog(trialName string, observationLog *v1beta1.ObservationLog) error {
	sqlQuery := "INSERT INTO observation_logs (trial_name, time, metric_name, value) VALUES "
	values := []interface{}{}

	for _, mlog := range observationLog.MetricLogs {
		if mlog.TimeStamp == "" {
			continue
		}
		t, err := time.Parse(time.RFC3339Nano, mlog.TimeStamp)
		if err != nil {
			return fmt.Errorf("Error parsing start time %s: %v", mlog.TimeStamp, err)
		}
		sqlTimeStr := t.UTC().Format(mysqlTimeFmt)

		sqlQuery += "(?, ?, ?, ?),"
		values = append(values, trialName, sqlTimeStr, mlog.Metric.Name, mlog.Metric.Value)
	}
	sqlQuery = sqlQuery[0 : len(sqlQuery)-1]

	// Prepare the statement
	stmt, err := d.db.Prepare(sqlQuery)
	if err != nil {
		return fmt.Errorf("Pepare SQL statement failed: %v", err)
	}

	// Execute INSERT
	_, err = stmt.Exec(values...)
	if err != nil {
		return fmt.Errorf("Execute SQL INSERT failed: %v", err)
	}

	return nil
}

func (d *dbConn) DeleteObservationLog(trialName string) error {
	_, err := d.db.Exec("DELETE FROM observation_logs WHERE trial_name = ?", trialName)
	return err
}

func (d *dbConn) GetObservationLog(trialName string, metricName string, startTime string, endTime string) (*v1beta1.ObservationLog, error) {
	qfield := []interface{}{trialName}
	qstr := ""
	if metricName != "" {
		qstr += " AND metric_name = ?"
		qfield = append(qfield, metricName)
	}
	if startTime != "" {
		s_time, err := time.Parse(time.RFC3339Nano, startTime)
		if err != nil {
			return nil, fmt.Errorf("Error parsing start time %s: %v", startTime, err)
		}
		formattedStartTime := s_time.UTC().Format(mysqlTimeFmt)
		qstr += " AND time >= ?"
		qfield = append(qfield, formattedStartTime)
	}
	if endTime != "" {
		e_time, err := time.Parse(time.RFC3339Nano, endTime)
		if err != nil {
			return nil, fmt.Errorf("Error parsing completion time %s: %v", endTime, err)
		}
		formattedEndTime := e_time.UTC().Format(mysqlTimeFmt)
		qstr += " AND time <= ?"
		qfield = append(qfield, formattedEndTime)
	}
	rows, err := d.db.Query("SELECT time, metric_name, value FROM observation_logs WHERE trial_name = ?"+qstr+" ORDER BY time",
		qfield...)
	if err != nil {
		return nil, fmt.Errorf("Failed to get ObservationLogs %v", err)
	}
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
		ptime, err := time.Parse(mysqlTimeFmt, sqlTimeStr)
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
