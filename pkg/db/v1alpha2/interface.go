package db

import (
	crand "crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/jsonpb"
	"k8s.io/klog"

	v1alpha2 "github.com/kubeflow/katib/pkg/api/v1alpha2"
)

const (
	dbDriver     = "mysql"
	dbNameTmpl   = "root:%s@tcp(katib-db:3306)/katib?timeout=5s"
	mysqlTimeFmt = "2006-01-02 15:04:05.999999"

	connectInterval = 5 * time.Second
	connectTimeout  = 60 * time.Second
)

type GetWorkerLogOpts struct {
	Name       string
	SinceTime  *time.Time
	Descending bool
	Limit      int32
	Objective  bool
}

type WorkerLog struct {
	Time  time.Time
	Name  string
	Value string
}

type KatibDBInterface interface {
	DBInit()
	SelectOne() error

	RegisterExperiment(experiment *v1alpha2.Experiment) error
	PreCheckRegisterExperiment(experiment *v1alpha2.Experiment) (bool, error)
	DeleteExperiment(experimentName string) error
	GetExperiment(experimentName string) (*v1alpha2.Experiment, error)
	GetExperimentList() ([]*v1alpha2.ExperimentSummary, error)
	UpdateExperimentStatus(experimentName string, newStatus *v1alpha2.ExperimentStatus) error

	UpdateAlgorithmExtraSettings(experimentName string, extraAlgorithmSetting []*v1alpha2.AlgorithmSetting) error
	GetAlgorithmExtraSettings(experimentName string) ([]*v1alpha2.AlgorithmSetting, error)

	RegisterTrial(trial *v1alpha2.Trial) error
	GetTrialList(experimentName string, filter string) ([]*v1alpha2.Trial, error)
	GetTrial(trialName string) (*v1alpha2.Trial, error)
	UpdateTrialStatus(trialName string, newStatus *v1alpha2.TrialStatus) error
	DeleteTrial(trialName string) error

	RegisterObservationLog(trialName string, observationLog *v1alpha2.ObservationLog) error
	GetObservationLog(trialName string, metricName string, startTime string, endTime string) (*v1alpha2.ObservationLog, error)
}

type dbConn struct {
	db *sql.DB
}

var rs1Letters = []rune("abcdefghijklmnopqrstuvwxyz")

func getDbName() string {
	dbPass := os.Getenv("MYSQL_ROOT_PASSWORD")
	if dbPass == "" {
		klog.Info("WARN: Env var MYSQL_ROOT_PASSWORD is empty. Falling back to \"test\".")

		// For backward compatibility, e.g. in case that all but vizier-core
		// is older ones so we do not have Secret nor upgraded vizier-db.
		dbPass = "test"
	}

	return fmt.Sprintf(dbNameTmpl, dbPass)
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
			}
		case <-timeoutC:
			return nil, fmt.Errorf("Timeout waiting for DB conn successfully opened.")
		}
	}
}

func NewWithSQLConn(db *sql.DB) (KatibDBInterface, error) {
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

func New() (KatibDBInterface, error) {
	db, err := openSQLConn(dbDriver, getDbName(), connectInterval, connectTimeout)
	if err != nil {
		return nil, fmt.Errorf("DB open failed: %v", err)
	}
	return NewWithSQLConn(db)
}

func (d *dbConn) RegisterExperiment(experiment *v1alpha2.Experiment) error {
	var paramSpecs string
	var objSpec string
	var algoSpec string
	var nasConfig string
	var start_time string
	var completion_time string
	var err error
	if experiment.Spec != nil {
		if experiment.Spec.ParameterSpecs != nil {
			paramSpecs, err = (&jsonpb.Marshaler{}).MarshalToString(experiment.Spec.ParameterSpecs)
			if err != nil {
				return fmt.Errorf("Error marshaling Parameters: %v", err)
			}
		}
		if experiment.Spec.Objective != nil {
			objSpec, err = (&jsonpb.Marshaler{}).MarshalToString(experiment.Spec.Objective)
			if err != nil {
				return fmt.Errorf("Error marshaling Objective: %v", err)
			}
		}
		if experiment.Spec.Algorithm != nil {
			algoSpec, err = (&jsonpb.Marshaler{}).MarshalToString(experiment.Spec.Algorithm)
			if err != nil {
				return fmt.Errorf("Error marshaling Algorithm: %v", err)
			}
		}
		if experiment.Spec.NasConfig != nil {
			nasConfig, err = (&jsonpb.Marshaler{}).MarshalToString(experiment.Spec.NasConfig)
			if err != nil {
				return fmt.Errorf("Error marshaling NasConfig: %v", err)
			}
		}
	} else {
		return fmt.Errorf("Invalid experiment: spec is nil.")
	}

	start_time = time.Now().UTC().Format(mysqlTimeFmt)
	completion_time = time.Time{}.UTC().Format(mysqlTimeFmt)
	if experiment.Status != nil {
		if experiment.Status.StartTime != "" {
			s_time, err := time.Parse(time.RFC3339Nano, experiment.Status.StartTime)
			if err != nil {
				return fmt.Errorf("Error parsing start time %s: %v", experiment.Status.StartTime, err)
			}
			start_time = s_time.UTC().Format(mysqlTimeFmt)
		}
		if experiment.Status.CompletionTime != "" {
			c_time, err := time.Parse(time.RFC3339Nano, experiment.Status.CompletionTime)
			if err != nil {
				return fmt.Errorf("Error parsing completion time %s: %v", experiment.Status.CompletionTime, err)
			}
			completion_time = c_time.UTC().Format(mysqlTimeFmt)
		}
	}

	_, err = d.db.Exec(
		`INSERT INTO experiments (
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
			nas_config) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		experiment.Name,
		paramSpecs,
		objSpec,
		algoSpec,
		experiment.Spec.TrialTemplate,
		experiment.Spec.MetricsCollectorSpec,
		experiment.Spec.ParallelTrialCount,
		experiment.Spec.MaxTrialCount,
		experiment.Status.Condition,
		start_time,
		completion_time,
		nasConfig,
	)
	return err
}

func (d *dbConn) DeleteExperiment(experimentName string) error {
	_, err := d.db.Exec("DELETE FROM experiments WHERE name = ?", experimentName)
	return err
}

func (d *dbConn) PreCheckRegisterExperiment(experiment *v1alpha2.Experiment) (bool, error) {
	_, err := d.GetExperiment(experiment.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, err
	} else {
		return false, nil
	}
}

func (d *dbConn) GetExperiment(experimentName string) (*v1alpha2.Experiment, error) {
	var id string
	var paramSpecs string
	var objSpec string
	var algoSpec string
	var nasConfig string
	var start_time string
	var completion_time string

	experiment := &v1alpha2.Experiment{
		Spec:   &v1alpha2.ExperimentSpec{},
		Status: &v1alpha2.ExperimentStatus{},
	}
	row := d.db.QueryRow("SELECT * FROM experiments WHERE name = ?", experimentName)
	err := row.Scan(
		&id,
		&experiment.Name,
		&paramSpecs,
		&objSpec,
		&algoSpec,
		&experiment.Spec.TrialTemplate,
		&experiment.Spec.MetricsCollectorSpec,
		&experiment.Spec.ParallelTrialCount,
		&experiment.Spec.MaxTrialCount,
		&experiment.Status.Condition,
		&start_time,
		&completion_time,
		&nasConfig,
	)
	if err != nil {
		return nil, err
	}
	if paramSpecs != "" {
		experiment.Spec.ParameterSpecs = new(v1alpha2.ExperimentSpec_ParameterSpecs)
		err = jsonpb.UnmarshalString(paramSpecs, experiment.Spec.ParameterSpecs)
		if err != nil {
			return nil, err
		}
	}
	if objSpec != "" {
		experiment.Spec.Objective = new(v1alpha2.ObjectiveSpec)
		err = jsonpb.UnmarshalString(objSpec, experiment.Spec.Objective)
		if err != nil {
			return nil, err
		}
	}
	if algoSpec != "" {
		experiment.Spec.Algorithm = new(v1alpha2.AlgorithmSpec)
		err = jsonpb.UnmarshalString(algoSpec, experiment.Spec.Algorithm)
		if err != nil {
			return nil, err
		}
	}
	if nasConfig != "" {
		experiment.Spec.NasConfig = new(v1alpha2.NasConfig)
		err = jsonpb.UnmarshalString(nasConfig, experiment.Spec.NasConfig)
		if err != nil {
			return nil, err
		}
	}
	if start_time != "" {
		start_timeMysql, err := time.Parse(mysqlTimeFmt, start_time)
		if err != nil {
			return nil, fmt.Errorf("Error parsing Trial start time %s to mysqlFormat: %v", start_time, err)
		}
		experiment.Status.StartTime = start_timeMysql.UTC().Format(time.RFC3339Nano)
	}
	if completion_time != "" {
		completion_timeMysql, err := time.Parse(mysqlTimeFmt, completion_time)
		if err != nil {
			return nil, fmt.Errorf("Error parsing Trial completion time %s to mysqlFormat: %v", completion_time, err)
		}
		experiment.Status.CompletionTime = completion_timeMysql.UTC().Format(time.RFC3339Nano)
	}
	return experiment, nil
}

func (d *dbConn) GetExperimentList() ([]*v1alpha2.ExperimentSummary, error) {
	rows, err := d.db.Query(`SELECT name, status, start_time, completion_time FROM experiments`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*v1alpha2.ExperimentSummary
	var start_time string
	var completion_time string
	for rows.Next() {
		experiment_sum := v1alpha2.ExperimentSummary{
			ExperimentName: "",
			Status:         &v1alpha2.ExperimentStatus{},
		}
		err = rows.Scan(
			&experiment_sum.ExperimentName,
			&experiment_sum.Status.Condition,
			&start_time,
			&completion_time,
		)
		if err != nil {
			return nil, fmt.Errorf("Fail to get Experiment from DB. %v", err)
		}
		if start_time != "" {
			start_timeMysql, err := time.Parse(mysqlTimeFmt, start_time)
			if err != nil {
				return nil, fmt.Errorf("Error parsing Trial start time %s to mysqlFormat: %v", start_time, err)
			}
			experiment_sum.Status.StartTime = start_timeMysql.UTC().Format(time.RFC3339Nano)
		}
		if completion_time != "" {
			completion_timeMysql, err := time.Parse(mysqlTimeFmt, completion_time)
			if err != nil {
				return nil, fmt.Errorf("Error parsing Trial completion time %s to mysqlFormat: %v", completion_time, err)
			}
			experiment_sum.Status.CompletionTime = completion_timeMysql.UTC().Format(time.RFC3339Nano)
		}
		result = append(result, &experiment_sum)
	}
	return result, nil
}

func (d *dbConn) UpdateExperimentStatus(experimentName string, newStatus *v1alpha2.ExperimentStatus) error {
	start_time := ""
	completion_time := ""
	var err error
	if newStatus.StartTime != "" {
		s_time, err := time.Parse(time.RFC3339Nano, newStatus.StartTime)
		if err != nil {
			return fmt.Errorf("Error parsing start time %s: %v", newStatus.StartTime, err)
		}
		start_time = s_time.UTC().Format(mysqlTimeFmt)
	}
	if newStatus.CompletionTime != "" {
		c_time, err := time.Parse(time.RFC3339Nano, newStatus.CompletionTime)
		if err != nil {
			return fmt.Errorf("Error parsing completion time %s: %v", newStatus.CompletionTime, err)
		}
		completion_time = c_time.UTC().Format(mysqlTimeFmt)
	}
	_, err = d.db.Exec(`UPDATE experiments SET status = ?,
		start_time = ?,
		completion_time = ? WHERE name = ?`,
		newStatus.Condition,
		start_time,
		completion_time,
		experimentName)
	return err
}

func (d *dbConn) UpdateAlgorithmExtraSettings(experimentName string, extraAlgorithmSetting []*v1alpha2.AlgorithmSetting) error {
	aesList, err := d.GetAlgorithmExtraSettings(experimentName)
	if err != nil {
		return fmt.Errorf("Failed to get current state %v", err)
	}
	for _, neas := range extraAlgorithmSetting {
		isin := false
		for _, ceas := range aesList {
			if ceas.Name == neas.Name {
				_, err = d.db.Exec(`UPDATE extra_algorithm_settings SET value = ? ,
						WHERE experiment_name = ? AND setting_name = ?`,
					neas.Value, experimentName, ceas.Name)
				if err != nil {
					return fmt.Errorf("Failed to update state %v", err)
				}
				isin = true
				break
			}
		}
		if !isin {
			_, err = d.db.Exec(
				`INSERT INTO extra_algorithm_settings (
			experiment_name,
			setting_name,
			value) VALUES (?, ?, ?)`,
				experimentName,
				neas.Name,
				neas.Value,
			)
			if err != nil {
				return fmt.Errorf("Failed to update state %v", err)
			}
		}
	}
	return nil
}

func (d *dbConn) GetAlgorithmExtraSettings(experimentName string) ([]*v1alpha2.AlgorithmSetting, error) {
	rows, err := d.db.Query("SELECT setting_name, value FROM extra_algorithm_settings WHERE experiment_name = ?", experimentName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*v1alpha2.AlgorithmSetting
	for rows.Next() {
		as := new(v1alpha2.AlgorithmSetting)
		err := rows.Scan(
			&as.Name,
			&as.Value,
		)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan ExtraSetting %v", err)
		}
		result = append(result, as)
	}
	return result, nil
}

func (d *dbConn) RegisterTrial(trial *v1alpha2.Trial) error {
	var objSpec string
	var paramAssignment string
	var start_time string
	var completion_time string
	var observation string
	var err error
	if trial.Spec != nil {
		if trial.Spec.Objective != nil {
			objSpec, err = (&jsonpb.Marshaler{}).MarshalToString(trial.Spec.Objective)
			if err != nil {
				return fmt.Errorf("Error marshaling Objective: %v", err)
			}
		}
		if trial.Spec.ParameterAssignments != nil {
			paramAssignment, err = (&jsonpb.Marshaler{}).MarshalToString(trial.Spec.ParameterAssignments)
			if err != nil {
				return fmt.Errorf("Error marshaling Parameters: %v", err)
			}
		}
	} else {
		return fmt.Errorf("Invalid trial: spec is nil.")
	}

	start_time = time.Now().UTC().Format(mysqlTimeFmt)
	completion_time = time.Time{}.UTC().Format(mysqlTimeFmt)

	if trial.Status != nil {
		if trial.Status.Observation != nil {
			observation, err = (&jsonpb.Marshaler{}).MarshalToString(trial.Status.Observation)
			if err != nil {
				return fmt.Errorf("Error marshaling Objective: %v", err)
			}
		}
		if trial.Status.StartTime != "" {
			s_time, err := time.Parse(time.RFC3339Nano, trial.Status.StartTime)
			if err != nil {
				return fmt.Errorf("Error parsing start time %s: %v", trial.Status.StartTime, err)
			}
			start_time = s_time.UTC().Format(mysqlTimeFmt)
		}
		if trial.Status.CompletionTime != "" {
			c_time, err := time.Parse(time.RFC3339Nano, trial.Status.CompletionTime)
			if err != nil {
				return fmt.Errorf("Error parsing completion time %s: %v", trial.Status.CompletionTime, err)
			}
			completion_time = c_time.UTC().Format(mysqlTimeFmt)
		}
	}
	_, err = d.db.Exec(
		`INSERT INTO trials (
			name, 
			experiment_name,
			objective,
			parameter_assignments,
			run_spec,
			metrics_collector_spec,
			observation,
			status,
			start_time,
			completion_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		trial.Name,
		trial.Spec.ExperimentName,
		objSpec,
		paramAssignment,
		trial.Spec.RunSpec,
		trial.Spec.MetricsCollectorSpec,
		observation,
		trial.Status.Condition,
		start_time,
		completion_time,
	)
	return err
}

func (d *dbConn) GetTrialList(experimentName string, filter string) ([]*v1alpha2.Trial, error) {
	var id string
	var objSpec string
	var paramAssignment string
	var start_time string
	var completion_time string
	var observation string
	var qstr = "SELECT * FROM trials WHERE experiment_name = ?"
	var qfield = []interface{}{experimentName}
	if filter != "" {
		//Currently only support filter by name.
		//TODO support other type of fiter
		//e.g.
		//* filter:name=foo
		//* filter:start_time>x
		//*filter:end_time<=y
		qstr += " AND name LIKE '%?%'"
		qfield = append(qfield, filter)
	}
	rows, err := d.db.Query(qstr, qfield...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*v1alpha2.Trial
	for rows.Next() {
		trial := &v1alpha2.Trial{
			Spec:   &v1alpha2.TrialSpec{},
			Status: &v1alpha2.TrialStatus{},
		}
		err := rows.Scan(
			&id,
			&trial.Name,
			&trial.Spec.ExperimentName,
			&objSpec,
			&paramAssignment,
			&trial.Spec.RunSpec,
			&trial.Spec.MetricsCollectorSpec,
			&observation,
			&trial.Status.Condition,
			&start_time,
			&completion_time,
		)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan trial %v", err)
		}
		if objSpec != "" {
			trial.Spec.Objective = new(v1alpha2.ObjectiveSpec)
			err = jsonpb.UnmarshalString(objSpec, trial.Spec.Objective)
			if err != nil {
				return nil, err
			}
		}
		if paramAssignment != "" {
			trial.Spec.ParameterAssignments = new(v1alpha2.TrialSpec_ParameterAssignments)
			err = jsonpb.UnmarshalString(paramAssignment, trial.Spec.ParameterAssignments)
			if err != nil {
				return nil, err
			}
		}
		if observation != "" {
			trial.Status.Observation = new(v1alpha2.Observation)
			err = jsonpb.UnmarshalString(observation, trial.Status.Observation)
			if err != nil {
				return nil, err
			}
		}
		if start_time != "" {
			start_timeMysql, err := time.Parse(mysqlTimeFmt, start_time)
			if err != nil {
				return nil, fmt.Errorf("Error parsing Trial start time %s to mysqlFormat: %v", start_time, err)
			}
			trial.Status.StartTime = start_timeMysql.UTC().Format(time.RFC3339Nano)
		}
		if completion_time != "" {
			completion_timeMysql, err := time.Parse(mysqlTimeFmt, completion_time)
			if err != nil {
				return nil, fmt.Errorf("Error parsing Trial completion time %s to mysqlFormat: %v", completion_time, err)
			}
			trial.Status.CompletionTime = completion_timeMysql.UTC().Format(time.RFC3339Nano)
		}
		result = append(result, trial)
	}
	return result, nil
}

func (d *dbConn) GetTrial(trialName string) (*v1alpha2.Trial, error) {
	var id string
	var objSpec string
	var paramAssignment string
	var start_time string
	var completion_time string
	var observation string
	trial := &v1alpha2.Trial{
		Spec:   &v1alpha2.TrialSpec{},
		Status: &v1alpha2.TrialStatus{},
	}
	row := d.db.QueryRow("SELECT * FROM trials WHERE name = ?", trialName)
	err := row.Scan(
		&id,
		&trial.Name,
		&trial.Spec.ExperimentName,
		&objSpec,
		&paramAssignment,
		&trial.Spec.RunSpec,
		&trial.Spec.MetricsCollectorSpec,
		&observation,
		&trial.Status.Condition,
		&start_time,
		&completion_time,
	)
	if objSpec != "" {
		trial.Spec.Objective = new(v1alpha2.ObjectiveSpec)
		err = jsonpb.UnmarshalString(objSpec, trial.Spec.Objective)
		if err != nil {
			return nil, err
		}
	}
	if paramAssignment != "" {
		trial.Spec.ParameterAssignments = new(v1alpha2.TrialSpec_ParameterAssignments)
		err = jsonpb.UnmarshalString(paramAssignment, trial.Spec.ParameterAssignments)
		if err != nil {
			return nil, err
		}
	}
	if observation != "" {
		trial.Status.Observation = new(v1alpha2.Observation)
		err = jsonpb.UnmarshalString(observation, trial.Status.Observation)
		if err != nil {
			return nil, err
		}
	}
	if start_time != "" {
		start_timeMysql, err := time.Parse(mysqlTimeFmt, start_time)
		if err != nil {
			return nil, fmt.Errorf("Error parsing Trial start time %s to mysqlFormat: %v", start_time, err)
		}
		trial.Status.StartTime = start_timeMysql.UTC().Format(time.RFC3339Nano)
	}
	if completion_time != "" {
		completion_timeMysql, err := time.Parse(mysqlTimeFmt, completion_time)
		if err != nil {
			return nil, fmt.Errorf("Error parsing Trial completion time %s to mysqlFormat: %v", completion_time, err)
		}
		trial.Status.CompletionTime = completion_timeMysql.UTC().Format(time.RFC3339Nano)
	}

	return trial, nil
}

func (d *dbConn) UpdateTrialStatus(trialName string, newStatus *v1alpha2.TrialStatus) error {
	var observation string = ""
	var formattedStartTime, formattedCompletionTime string = "", ""
	var err error
	if newStatus.Observation != nil {
		observation, err = (&jsonpb.Marshaler{}).MarshalToString(newStatus.Observation)
		if err != nil {
			return fmt.Errorf("Error marshaling Objective: %v", err)
		}
	}

	if newStatus.StartTime != "" {
		start_time, err := time.Parse(time.RFC3339Nano, newStatus.StartTime)
		if err != nil {
			return fmt.Errorf("Error parsing start time %s: %v", newStatus.StartTime, err)
		}
		formattedStartTime = start_time.UTC().Format(mysqlTimeFmt)
	}
	if newStatus.CompletionTime != "" {
		completion_time, err := time.Parse(time.RFC3339Nano, newStatus.CompletionTime)
		if err != nil {
			return fmt.Errorf("Error parsing completion time %s: %v", newStatus.CompletionTime, err)
		}
		formattedCompletionTime = completion_time.UTC().Format(mysqlTimeFmt)
	}
	_, err = d.db.Exec(`UPDATE trials SET status = ?,
		start_time = ?,
		completion_time = ?,
		observation = ? WHERE name = ?`,
		newStatus.Condition,
		formattedStartTime,
		formattedCompletionTime,
		observation,
		trialName)
	return err
}
func (d *dbConn) DeleteTrial(trialName string) error {
	_, err := d.db.Exec("DELETE FROM trials WHERE name = ?", trialName)
	return err
}

func (d *dbConn) RegisterObservationLog(trialName string, observationLog *v1alpha2.ObservationLog) error {
	var mname, mvalue string
	for _, mlog := range observationLog.MetricLogs {
		mname = mlog.Metric.Name
		mvalue = mlog.Metric.Value
		if mlog.TimeStamp == "" {
			continue
		}
		t, err := time.Parse(time.RFC3339Nano, mlog.TimeStamp)
		if err != nil {
			return fmt.Errorf("Error parsing start time %s: %v", mlog.TimeStamp, err)
		}
		sqlTimeStr := t.UTC().Format(mysqlTimeFmt)
		_, err = d.db.Exec(
			`INSERT INTO observation_logs (
				trial_name,
				time,
				metric_name,
				value
			) VALUES (?, ?, ?, ?)`,
			trialName,
			sqlTimeStr,
			mname,
			mvalue,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
func (d *dbConn) GetObservationLog(trialName string, metricName string, startTime string, endTime string) (*v1alpha2.ObservationLog, error) {
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
	result := &v1alpha2.ObservationLog{
		MetricLogs: []*v1alpha2.MetricLog{},
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
		result.MetricLogs = append(result.MetricLogs, &v1alpha2.MetricLog{
			TimeStamp: timeStamp,
			Metric: &v1alpha2.Metric{
				Name:  mname,
				Value: mvalue,
			},
		})
	}
	return result, nil
}
