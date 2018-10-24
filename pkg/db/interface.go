package db

import (
	crand "crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"log"
	"math/big"
	"math/rand"
	"strings"
	"time"

	api "github.com/kubeflow/katib/pkg/api"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbDriver     = "mysql"
	dbName       = "root:test@tcp(vizier-db:3306)/vizier"
	mysqlTimeFmt = "2006-01-02 15:04:05.999999"
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

type VizierDBInterface interface {
	DBInit()
	GetStudyConfig(string) (*api.StudyConfig, error)
	GetStudyList() ([]string, error)
	CreateStudy(*api.StudyConfig) (string, error)
	DeleteStudy(string) error

	GetTrial(string) (*api.Trial, error)
	GetTrialList(string) ([]*api.Trial, error)
	CreateTrial(*api.Trial) error
	DeleteTrial(string) error

	GetWorker(string) (*api.Worker, error)
	GetWorkerStatus(string) (*api.State, error)
	GetWorkerList(string, string) ([]*api.Worker, error)
	GetWorkerLogs(string, *GetWorkerLogOpts) ([]*WorkerLog, error)
	GetWorkerTimestamp(string) (*time.Time, error)
	StoreWorkerLogs(string, []*api.MetricsLog) error
	CreateWorker(*api.Worker) (string, error)
	UpdateWorker(string, api.State) error
	DeleteWorker(string) error

	SetSuggestionParam(string, string, []*api.SuggestionParameter) (string, error)
	UpdateSuggestionParam(string, []*api.SuggestionParameter) error
	GetSuggestionParam(string) ([]*api.SuggestionParameter, error)
	GetSuggestionParamList(string) ([]*api.SuggestionParameterSet, error)
	SetEarlyStopParam(string, string, []*api.EarlyStoppingParameter) (string, error)
	UpdateEarlyStopParam(string, []*api.EarlyStoppingParameter) error
	GetEarlyStopParam(string) ([]*api.EarlyStoppingParameter, error)
	GetEarlyStopParamList(string) ([]*api.EarlyStoppingParameterSet, error)
}

type dbConn struct {
	db *sql.DB
}

var rs1Letters = []rune("abcdefghijklmnopqrstuvwxyz")

func NewWithSQLConn(db *sql.DB) VizierDBInterface {
	d := new(dbConn)
	d.db = db
	seed, err := crand.Int(crand.Reader, big.NewInt(1<<63-1))
	if err != nil {
		log.Fatalf("RNG initialization failed: %v", err)
	}
	// We can do the following instead, but it creates a locking issue
	//d.rng = rand.New(rand.NewSource(seed.Int64()))
	rand.Seed(seed.Int64())

	return d
}

func New() VizierDBInterface {
	db, err := sql.Open(dbDriver, dbName)
	if err != nil {
		log.Fatalf("DB open failed: %v", err)
	}
	return NewWithSQLConn(db)
}

func generateRandid() string {
	// UUID isn't quite handy in the Go world
	id := make([]byte, 8)
	_, err := rand.Read(id)
	if err != nil {
		log.Printf("Error reading random: %v", err)
		return ""
	}
	return string(rs1Letters[rand.Intn(len(rs1Letters))]) + fmt.Sprintf("%016x", id)[1:]
}

func (d *dbConn) GetStudyConfig(id string) (*api.StudyConfig, error) {
	row := d.db.QueryRow("SELECT * FROM studies WHERE id = ?", id)

	study := new(api.StudyConfig)
	var dummyID, configs, tags, metrics string
	err := row.Scan(&dummyID,
		&study.Name,
		&study.Owner,
		&study.OptimizationType,
		&study.OptimizationGoal,
		&configs,
		&tags,
		&study.ObjectiveValueName,
		&metrics,
		&study.JobId,
	)
	if err != nil {
		return nil, err
	}
	study.ParameterConfigs = new(api.StudyConfig_ParameterConfigs)
	err = jsonpb.UnmarshalString(configs, study.ParameterConfigs)
	if err != nil {
		return nil, err
	}

	var tagsArray []string
	if len(tags) > 0 {
		tagsArray = strings.Split(tags, ",\n")
	}
	study.Tags = make([]*api.Tag, len(tagsArray))
	for i, j := range tagsArray {
		tag := new(api.Tag)
		err = jsonpb.UnmarshalString(j, tag)
		if err != nil {
			log.Printf("err unmarshal %s", j)
			return nil, err
		}
		study.Tags[i] = tag
	}
	study.Metrics = strings.Split(metrics, ",\n")
	return study, nil
}

func (d *dbConn) GetStudyList() ([]string, error) {
	rows, err := d.db.Query("SELECT id FROM studies")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var result []string
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			log.Printf("err scanning studies.id: %v", err)
			continue
		}
		result = append(result, id)
	}
	return result, nil
}

func (d *dbConn) CreateStudy(in *api.StudyConfig) (string, error) {
	if in.ParameterConfigs == nil {
		return "", errors.New("ParameterConfigs must be set")

	}
	if in.JobId != "" {
		row := d.db.QueryRow("SELECT * FROM studies WHERE job_id = ?", in.JobId)
		dummyStudy := new(api.StudyConfig)
		var dummyID, dummyConfigs, dummyTags, dummyMetrics, dummyJobID string
		err := row.Scan(&dummyID,
			&dummyStudy.Name,
			&dummyStudy.Owner,
			&dummyStudy.OptimizationType,
			&dummyStudy.OptimizationGoal,
			&dummyConfigs,
			&dummyTags,
			&dummyStudy.ObjectiveValueName,
			&dummyMetrics,
			&dummyJobID,
		)
		if err == nil {
			return "", fmt.Errorf("Study %s in Job %s already exist.", in.Name, in.JobId)
		}
	}

	configs, err := (&jsonpb.Marshaler{}).MarshalToString(in.ParameterConfigs)
	if err != nil {
		log.Fatalf("Error marshaling configs: %v", err)
	}

	tags := make([]string, len(in.Tags))
	for i, elem := range in.Tags {
		tags[i], err = (&jsonpb.Marshaler{}).MarshalToString(elem)
		if err != nil {
			log.Printf("Error marshalling %v: %v", elem, err)
			continue
		}
	}

	var isin bool = false
	for _, m := range in.Metrics {
		if m == in.ObjectiveValueName {
			isin = true
		}
	}
	if !isin {
		in.Metrics = append(in.Metrics, in.ObjectiveValueName)
	}

	var studyID string
	i := 3
	for true {
		studyID = generateRandid()
		_, err := d.db.Exec(
			"INSERT INTO studies VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			studyID,
			in.Name,
			in.Owner,
			in.OptimizationType,
			in.OptimizationGoal,
			configs,
			strings.Join(tags, ",\n"),
			in.ObjectiveValueName,
			strings.Join(in.Metrics, ",\n"),
			in.JobId,
		)
		if err == nil {
			break
		} else {
			errmsg := strings.ToLower(err.Error())
			if strings.Contains(errmsg, "unique") || strings.Contains(errmsg, "duplicate") {
				i--
				if i > 0 {
					continue
				}
			}
		}
		return "", err
	}
	for _, perm := range in.AccessPermissions {
		_, err := d.db.Exec(
			"INSERT INTO study_permissions (study_id, access_permission) "+
				"VALUES (?, ?)",
			studyID, perm)
		if err != nil {
			log.Printf("Error storing permission (%s, %s): %v",
				studyID, perm, err)
		}
	}

	return studyID, nil
}

func (d *dbConn) DeleteStudy(id string) error {
	_, err := d.db.Exec("DELETE FROM studies WHERE id = ?", id)
	return err
}

func (d *dbConn) getTrials(trialID string, studyID string) ([]*api.Trial, error) {
	var rows *sql.Rows
	var err error

	if trialID != "" {
		rows, err = d.db.Query("SELECT * FROM trials WHERE id = ?", trialID)
	} else if studyID != "" {
		rows, err = d.db.Query("SELECT * FROM trials WHERE study_id = ?", studyID)
	} else {
		return nil, errors.New("trial_id or study_id must be set")
	}

	if err != nil {
		return nil, err
	}

	var result []*api.Trial
	for rows.Next() {
		trial := new(api.Trial)

		var parameters, tags string
		err := rows.Scan(&trial.TrialId,
			&trial.StudyId,
			&parameters,
			&trial.ObjectiveValue,
			&tags,
		)
		if err != nil {
			return nil, err
		}
		params := strings.Split(parameters, ",\n")
		p := make([]*api.Parameter, len(params))
		for i, pstr := range params {
			if pstr == "" {
				continue
			}
			p[i] = &api.Parameter{}
			err := jsonpb.UnmarshalString(pstr, p[i])
			if err != nil {
				return nil, err
			}
		}
		trial.ParameterSet = p
		taglist := strings.Split(tags, ",\n")
		t := make([]*api.Tag, len(taglist))
		for i, tstr := range taglist {
			t[i] = &api.Tag{}
			if tstr == "" {
				continue
			}
			err := jsonpb.UnmarshalString(tstr, t[i])
			if err != nil {
				return nil, err
			}
		}
		trial.Tags = t
		result = append(result, trial)
	}

	return result, nil
}

func (d *dbConn) GetTrial(id string) (*api.Trial, error) {
	trials, err := d.getTrials(id, "")
	if err != nil {
		return nil, err
	}

	if len(trials) > 1 {
		return nil, errors.New("multiple trials found")
	} else if len(trials) == 0 {
		return nil, errors.New("trials not found")
	}

	return trials[0], nil
}

func (d *dbConn) GetTrialList(id string) ([]*api.Trial, error) {
	trials, err := d.getTrials("", id)

	return trials, err
}

func (d *dbConn) CreateTrial(trial *api.Trial) error {
	// This function sets trial.id, unlike old dbInsertTrials().
	// Users should not overwrite trial.id
	var err, lastErr error

	params := make([]string, len(trial.ParameterSet))
	for i, elem := range trial.ParameterSet {
		params[i], err = (&jsonpb.Marshaler{}).MarshalToString(elem)
		if err != nil {
			log.Printf("Error marshalling trial.ParameterSet %v: %v",
				elem, err)
			lastErr = err
		}
	}
	tags := make([]string, len(trial.Tags))
	for i := range tags {
		tags[i], err = (&jsonpb.Marshaler{}).MarshalToString(trial.Tags[i])
		if err != nil {
			log.Printf("Error marshalling trial.Tags %v: %v",
				trial.Tags[i], err)
			lastErr = err
		}
	}

	var trialID string
	i := 3
	for true {
		trialID = generateRandid()
		_, err = d.db.Exec("INSERT INTO trials VALUES (?, ?, ?, ?, ?)",
			trialID, trial.StudyId, strings.Join(params, ",\n"),
			trial.ObjectiveValue, strings.Join(tags, ",\n"))
		if err == nil {
			trial.TrialId = trialID
			break
		} else {
			errmsg := strings.ToLower(err.Error())
			if strings.Contains(errmsg, "unique") || strings.Contains(errmsg, "duplicate") {
				i--
				if i > 0 {
					continue
				}
			}
		}
		return err
	}
	return lastErr
}

func (d *dbConn) DeleteTrial(id string) error {
	_, err := d.db.Exec("DELETE FROM trials WHERE id = ?", id)
	return err
}

func (d *dbConn) GetWorkerLogs(id string, opts *GetWorkerLogOpts) ([]*WorkerLog, error) {
	qstr := ""
	qfield := []interface{}{id}
	order := ""
	if opts != nil {
		if opts.SinceTime != nil {
			qstr += " AND time >= ?"
			qfield = append(qfield, opts.SinceTime)
		}
		if opts.Name != "" {
			qstr += " AND name = ?"
			qfield = append(qfield, opts.Name)
		}
		if opts.Objective {
			qstr += " AND is_objective = 1"
		}
		if opts.Descending {
			order = " DESC"
		}
		if opts.Limit > 0 {
			order += fmt.Sprintf(" LIMIT %d", opts.Limit)
		}
	}

	rows, err := d.db.Query("SELECT time, name, value FROM worker_metrics WHERE worker_id = ?"+
		qstr+" ORDER BY time"+order, qfield...)
	if err != nil {
		return nil, err
	}

	var result []*WorkerLog
	for rows.Next() {
		log1 := new(WorkerLog)
		var timeStr string

		err := rows.Scan(&timeStr, &((*log1).Name), &((*log1).Value))
		if err != nil {
			log.Printf("Error scanning log: %v", err)
			continue
		}
		log1.Time, err = time.Parse(mysqlTimeFmt, timeStr)
		if err != nil {
			log.Printf("Error parsing time %s: %v", timeStr, err)
			continue
		}
		result = append(result, log1)
	}
	return result, nil
}

func (d *dbConn) getWorkerLastlogs(id string) (time.Time, []*WorkerLog, error) {
	var timeStr string
	var timeVal time.Time
	var err error

	// Use LEFT JOIN to ensure a result even if there's no matching
	// in worker_metrics.
	rows, err := d.db.Query(
		`SELECT worker_lastlogs.time, name, value FROM worker_lastlogs
                 LEFT JOIN worker_metrics
                 ON (worker_lastlogs.worker_id = worker_metrics.worker_id AND worker_lastlogs.time = worker_metrics.time)
                 WHERE worker_lastlogs.worker_id = ?`, id)
	if err != nil {
		return timeVal, nil, err
	}

	var result []*WorkerLog
	for rows.Next() {
		log1 := new(WorkerLog)
		var thisTime string
		var name, value sql.NullString

		err := rows.Scan(&thisTime, &name, &value)
		if err != nil {
			log.Printf("Error scanning log: %v", err)
			continue
		}
		if timeStr == "" {
			timeStr = thisTime
			timeVal, err = time.Parse(mysqlTimeFmt, timeStr)
			if err != nil {
				log.Printf("Error parsing time %s: %v", timeStr, err)
				return timeVal, nil, err
			}
		} else if timeStr != thisTime {
			log.Printf("Unexpected query result %s != %s",
				timeStr, thisTime)
		}
		log1.Time = timeVal
		if !name.Valid {
			continue
		}
		(*log1).Name = name.String
		(*log1).Value = value.String
		result = append(result, log1)
	}
	return timeVal, result, nil
}

func (d *dbConn) GetWorkerTimestamp(id string) (*time.Time, error) {
	var lastTimestamp string

	row := d.db.QueryRow("SELECT time FROM worker_lastlogs WHERE worker_id = ?", id)
	err := row.Scan(&lastTimestamp)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		mt, err := time.Parse(mysqlTimeFmt, lastTimestamp)
		if err != nil {
			log.Printf("Error parsing time in log %s: %v",
				lastTimestamp, err)
			return nil, err
		}
		return &mt, nil
	}
}

func (d *dbConn) storeWorkerLog(workerID string, time string, metricsName string, metricsValue string, objectiveValueName string) error {
	isObjective := 0
	if metricsName == objectiveValueName {
		isObjective = 1
	}
	_, err := d.db.Exec("INSERT INTO worker_metrics (worker_id, time, name, value, is_objective) VALUES (?, ?, ?, ?, ?)",
		workerID, time, metricsName, metricsValue, isObjective)
	if err != nil {
		return err
	}
	return nil
}

func (d *dbConn) StoreWorkerLogs(workerID string, logs []*api.MetricsLog) error {
	var lasterr error

	dbT, lastLogs, err := d.getWorkerLastlogs(workerID)
	if err != nil {
		log.Printf("Error getting last log timestamp: %v", err)
	}

	row := d.db.QueryRow("SELECT objective_value_name FROM workers "+
		"JOIN (studies) ON (workers.study_id = studies.id) WHERE "+
		"workers.id = ?", workerID)
	var objectiveValueName string
	err = row.Scan(&objectiveValueName)
	if err != nil {
		log.Printf("Cannot get objective_value_name or metrics: %v", err)
		return err
	}

	// Store logs when
	//   1. a log is newer than dbT, or,
	//   2. a log is not yet in the DB when the timestamps are equal
	var formattedTime string
	var lastTime time.Time
	for _, mlog := range logs {
		metricsName := mlog.Name
	logLoop:
		for _, mv := range mlog.Values {
			t, err := time.Parse(time.RFC3339Nano, mv.Time)
			if err != nil {
				log.Printf("Error parsing time %s: %v", mv.Time, err)
				lasterr = err
				continue
			}
			if t.Before(dbT) {
				// dbT is from mysql and has microsec precision.
				// This code assumes nanosec fractions are rounded down.
				continue
			}
			// use UTC as mysql DATETIME lacks timezone
			formattedTime = t.UTC().Format(mysqlTimeFmt)
			if !dbT.IsZero() {
				// Parse again to get rounding effect, otherwise
				// the next comparison will be almost always false.
				reparsed_time, err := time.Parse(mysqlTimeFmt, formattedTime)
				if err != nil {
					log.Printf("Error parsing time %s: %v", formattedTime, err)
					lasterr = err
					continue
				}
				if reparsed_time == dbT {
					for _, l := range lastLogs {
						if l.Name == metricsName && l.Value == mv.Value {
							continue logLoop
						}
					}
				}
			}
			err = d.storeWorkerLog(workerID,
				formattedTime,
				metricsName, mv.Value,
				objectiveValueName)
			if err != nil {
				log.Printf("Error storing log %s: %v", mv.Value, err)
				lasterr = err
			} else if t.After(lastTime) {
				lastTime = t
			}
		}
	}
	if lasterr != nil {
		// If lastlog were updated, logs that couldn't be saved
		// would be lost.
		return lasterr
	}
	if !lastTime.IsZero() {
		formattedTime = lastTime.UTC().Format(mysqlTimeFmt)
		_, err = d.db.Exec("REPLACE INTO worker_lastlogs VALUES (?, ?)",
			workerID, formattedTime)
	}
	return err
}

func (d *dbConn) getWorkers(workerID string, trialID string, studyID string) ([]*api.Worker, error) {
	var rows *sql.Rows
	var err error

	if workerID != "" {
		rows, err = d.db.Query("SELECT * FROM workers WHERE id = ?", workerID)
	} else if trialID != "" {
		rows, err = d.db.Query("SELECT * FROM workers WHERE trial_id = ?", trialID)
	} else if studyID != "" {
		rows, err = d.db.Query("SELECT * FROM workers WHERE study_id = ?", studyID)
	} else {
		return nil, errors.New("worker_id, trial_id or study_id must be set")
	}

	if err != nil {
		return nil, err
	}

	var result []*api.Worker
	for rows.Next() {
		worker := new(api.Worker)

		var tags string
		err := rows.Scan(
			&worker.WorkerId,
			&worker.StudyId,
			&worker.TrialId,
			&worker.Type,
			&worker.Status,
			&worker.TemplatePath,
			&tags,
		)
		if err != nil {
			return nil, err
		}

		taglist := strings.Split(tags, ",\n")
		t := make([]*api.Tag, len(taglist))
		for i, tstr := range taglist {
			t[i] = &api.Tag{}
			if tstr == "" {
				continue
			}
			err := jsonpb.UnmarshalString(tstr, t[i])
			if err != nil {
				return nil, err
			}
		}
		worker.Tags = t
		result = append(result, worker)
	}
	return result, nil
}

func (d *dbConn) GetWorker(id string) (*api.Worker, error) {
	workers, err := d.getWorkers(id, "", "")
	if err != nil {
		return nil, err
	}
	if len(workers) > 1 {
		return nil, errors.New("multiple workers found")
	} else if len(workers) == 0 {
		return nil, errors.New("worker not found")
	}

	return workers[0], nil

}

func (d *dbConn) GetWorkerStatus(id string) (*api.State, error) {
	status := api.State_ERROR
	row := d.db.QueryRow("SELECT status FROM workers WHERE id = ?", id)
	err := row.Scan(&status)
	if err != nil {
		return &status, err
	}
	return &status, nil
}

func (d *dbConn) GetWorkerList(sid string, tid string) ([]*api.Worker, error) {
	workers, err := d.getWorkers("", tid, sid)
	return workers, err
}

func (d *dbConn) CreateWorker(worker *api.Worker) (string, error) {
	// Users should not overwrite worker.id
	var err, lastErr error
	tags := make([]string, len(worker.Tags))
	for i := range tags {
		tags[i], err = (&jsonpb.Marshaler{}).MarshalToString(worker.Tags[i])
		if err != nil {
			log.Printf("Error marshalling worker.Tags %v: %v",
				worker.Tags[i], err)
			lastErr = err
		}
	}

	var workerID string
	i := 3
	for true {
		workerID = generateRandid()
		_, err = d.db.Exec("INSERT INTO workers VALUES (?, ?, ?, ?, ?, ?, ?)",
			workerID, worker.StudyId, worker.TrialId, worker.Type,
			api.State_PENDING, worker.TemplatePath, strings.Join(tags, ",\n"))
		if err == nil {
			worker.WorkerId = workerID
			break
		} else {
			errmsg := strings.ToLower(err.Error())
			if strings.Contains(errmsg, "unique") || strings.Contains(errmsg, "duplicate") {
				i--
				if i > 0 {
					continue
				}
			}
		}
		return "", err
	}
	return worker.WorkerId, lastErr

}

func (d *dbConn) UpdateWorker(id string, newstatus api.State) error {
	_, err := d.db.Exec("UPDATE workers SET status = ? WHERE id = ?", newstatus, id)
	return err
}

func (d *dbConn) DeleteWorker(id string) error {
	_, err := d.db.Exec("DELETE FROM workers WHERE id = ?", id)
	return err
}

func (d *dbConn) SetSuggestionParam(algorithm string, studyID string, params []*api.SuggestionParameter) (string, error) {
	var err error
	ps := make([]string, len(params))
	for i, elem := range params {
		ps[i], err = (&jsonpb.Marshaler{}).MarshalToString(elem)
		if err != nil {
			log.Printf("Error marshalling %v: %v", elem, err)
			return "", err
		}
	}
	var paramID string
	for true {
		paramID = generateRandid()
		_, err = d.db.Exec("INSERT INTO suggestion_param VALUES (?, ?, ?, ?)",
			paramID, algorithm, studyID, strings.Join(ps, ",\n"))
		if err == nil {
			break
		}
	}
	return paramID, err
}

func (d *dbConn) UpdateSuggestionParam(paramID string, params []*api.SuggestionParameter) error {
	var err error
	ps := make([]string, len(params))
	for i, elem := range params {
		ps[i], err = (&jsonpb.Marshaler{}).MarshalToString(elem)
		if err != nil {
			log.Printf("Error marshalling %v: %v", elem, err)
			return err
		}
	}
	_, err = d.db.Exec("UPDATE suggestion_param SET parameters = ? WHERE id = ?",
		strings.Join(ps, ",\n"), paramID)
	return err
}

func (d *dbConn) GetSuggestionParam(paramID string) ([]*api.SuggestionParameter, error) {
	var params string
	row := d.db.QueryRow("SELECT parameters FROM suggestion_param WHERE id = ?", paramID)
	err := row.Scan(&params)
	if err != nil {
		return nil, err
	}
	var pArray []string
	if len(params) > 0 {
		pArray = strings.Split(params, ",\n")
	} else {
		return nil, nil
	}
	ret := make([]*api.SuggestionParameter, len(pArray))
	for i, j := range pArray {
		p := new(api.SuggestionParameter)
		err = jsonpb.UnmarshalString(j, p)
		if err != nil {
			log.Printf("err unmarshal %s", j)
			return nil, err
		}
		ret[i] = p
	}
	return ret, nil
}

func (d *dbConn) GetSuggestionParamList(studyID string) ([]*api.SuggestionParameterSet, error) {
	var rows *sql.Rows
	var err error
	rows, err = d.db.Query("SELECT * FROM suggestion_param WHERE study_id = ?", studyID)
	if err != nil {
		return nil, err
	}
	var result []*api.SuggestionParameterSet
	for rows.Next() {
		var id string
		var algorithm string
		var params string
		var sID string
		err := rows.Scan(&id, &sID, &algorithm, &params)
		if err != nil {
			return nil, err
		}
		if studyID != sID {
			continue
		}
		var pArray []string
		if len(params) > 0 {
			pArray = strings.Split(params, ",\n")
		} else {
			return nil, nil
		}
		suggestparams := make([]*api.SuggestionParameter, len(pArray))
		for i, j := range pArray {
			p := new(api.SuggestionParameter)
			err = jsonpb.UnmarshalString(j, p)
			if err != nil {
				log.Printf("err unmarshal %s", j)
				return nil, err
			}
			suggestparams[i] = p
		}
		result = append(result, &api.SuggestionParameterSet{
			ParamId:              id,
			SuggestionAlgorithm:  algorithm,
			SuggestionParameters: suggestparams,
		})
	}
	return result, nil
}

func (d *dbConn) SetEarlyStopParam(algorithm string, studyID string, params []*api.EarlyStoppingParameter) (string, error) {
	ps := make([]string, len(params))
	var err error
	for i, elem := range params {
		ps[i], err = (&jsonpb.Marshaler{}).MarshalToString(elem)
		if err != nil {
			log.Printf("Error marshalling %v: %v", elem, err)
			return "", err
		}
	}
	var paramID string
	for true {
		paramID := generateRandid()
		_, err = d.db.Exec("INSERT INTO earlystopping_param VALUES (?,?, ?, ?)",
			paramID, algorithm, studyID, strings.Join(ps, ",\n"))
		if err == nil {
			break
		}
	}
	return paramID, nil
}

func (d *dbConn) UpdateEarlyStopParam(paramID string, params []*api.EarlyStoppingParameter) error {
	ps := make([]string, len(params))
	var err error
	for i, elem := range params {
		ps[i], err = (&jsonpb.Marshaler{}).MarshalToString(elem)
		if err != nil {
			log.Printf("Error marshalling %v: %v", elem, err)
			return err
		}
	}
	_, err = d.db.Exec("UPDATE earlystopping_param SET parameters = ? WHERE id = ?",
		strings.Join(ps, ",\n"), paramID)
	return err
}

func (d *dbConn) GetEarlyStopParam(paramID string) ([]*api.EarlyStoppingParameter, error) {
	var params string
	row := d.db.QueryRow("SELECT parameters FROM earlystopping_param WHERE id = ?", paramID)
	err := row.Scan(&params)
	if err != nil {
		return nil, err
	}
	var pArray []string
	if len(params) > 0 {
		pArray = strings.Split(params, ",\n")
	} else {
		return nil, nil
	}
	ret := make([]*api.EarlyStoppingParameter, len(pArray))
	for i, j := range pArray {
		p := new(api.EarlyStoppingParameter)
		err = jsonpb.UnmarshalString(j, p)
		if err != nil {
			log.Printf("err unmarshal %s", j)
			return nil, err
		}
		ret[i] = p
	}
	return ret, nil
}

func (d *dbConn) GetEarlyStopParamList(studyID string) ([]*api.EarlyStoppingParameterSet, error) {
	var rows *sql.Rows
	var err error
	rows, err = d.db.Query("SELECT * FROM earlystopping_param WHERE study_id = ?", studyID)
	if err != nil {
		return nil, err
	}
	var result []*api.EarlyStoppingParameterSet
	for rows.Next() {
		var id string
		var algorithm string
		var params string
		var sID string
		err := rows.Scan(&id, &sID, &algorithm, &params)
		if err != nil {
			return nil, err
		}
		if studyID != sID {
			continue
		}
		var pArray []string
		if len(params) > 0 {
			pArray = strings.Split(params, ",\n")
		} else {
			return nil, nil
		}
		esparams := make([]*api.EarlyStoppingParameter, len(pArray))
		for i, j := range pArray {
			p := new(api.EarlyStoppingParameter)
			err = jsonpb.UnmarshalString(j, p)
			if err != nil {
				log.Printf("err unmarshal %s", j)
				return nil, err
			}
			esparams[i] = p
		}
		result = append(result, &api.EarlyStoppingParameterSet{
			ParamId:                 id,
			EarlyStoppingAlgorithm:  algorithm,
			EarlyStoppingParameters: esparams,
		})
	}
	return result, nil
}
