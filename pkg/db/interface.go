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
	db_driver      = "mysql"
	db_name        = "root:test@tcp(vizier-db:3306)/vizier"
	mysql_time_fmt = "2006-01-02 15:04:05.999999"
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
	DB_Init()
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

type db_conn struct {
	db *sql.DB
}

var rs1Letters = []rune("abcdefghijklmnopqrstuvwxyz")

func NewWithSqlConn(db *sql.DB) VizierDBInterface {
	d := new(db_conn)
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
	db, err := sql.Open(db_driver, db_name)
	if err != nil {
		log.Fatalf("DB open failed: %v", err)
	}
	return NewWithSqlConn(db)
}

func generate_randid() string {
	// UUID isn't quite handy in the Go world
	id_ := make([]byte, 8)
	_, err := rand.Read(id_)
	if err != nil {
		log.Printf("Error reading random: %v", err)
		return ""
	}
	return string(rs1Letters[rand.Intn(len(rs1Letters))]) + fmt.Sprintf("%016x", id_)[1:]
}

func (d *db_conn) GetStudyConfig(id string) (*api.StudyConfig, error) {
	row := d.db.QueryRow("SELECT * FROM studies WHERE id = ?", id)

	study := new(api.StudyConfig)
	var dummy_id, configs, tags, metrics string
	err := row.Scan(&dummy_id,
		&study.Name,
		&study.Owner,
		&study.OptimizationType,
		&study.OptimizationGoal,
		&configs,
		&tags,
		&study.ObjectiveValueName,
		&metrics,
	)
	if err != nil {
		return nil, err
	}
	study.ParameterConfigs = new(api.StudyConfig_ParameterConfigs)
	err = jsonpb.UnmarshalString(configs, study.ParameterConfigs)
	if err != nil {
		return nil, err
	}

	var tags_array []string
	if len(tags) > 0 {
		tags_array = strings.Split(tags, ",\n")
	}
	study.Tags = make([]*api.Tag, len(tags_array))
	for i, j := range tags_array {
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

func (d *db_conn) GetStudyList() ([]string, error) {
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

func (d *db_conn) CreateStudy(in *api.StudyConfig) (string, error) {
	if in.ParameterConfigs == nil {
		return "", errors.New("ParameterConfigs must be set")

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

	var study_id string
	i := 3
	for true {
		study_id = generate_randid()
		_, err := d.db.Exec(
			"INSERT INTO studies VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
			study_id,
			in.Name,
			in.Owner,
			in.OptimizationType,
			in.OptimizationGoal,
			configs,
			strings.Join(tags, ",\n"),
			in.ObjectiveValueName,
			strings.Join(in.Metrics, ",\n"),
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
			study_id, perm)
		if err != nil {
			log.Printf("Error storing permission (%s, %s): %v",
				study_id, perm, err)
		}
	}

	return study_id, nil
}

func (d *db_conn) DeleteStudy(id string) error {
	_, err := d.db.Exec("DELETE FROM studies WHERE id = ?", id)
	return err
}

func (d *db_conn) getTrials(trial_id string, study_id string) ([]*api.Trial, error) {
	var rows *sql.Rows
	var err error

	if trial_id != "" {
		rows, err = d.db.Query("SELECT * FROM trials WHERE id = ?", trial_id)
	} else if study_id != "" {
		rows, err = d.db.Query("SELECT * FROM trials WHERE study_id = ?", study_id)
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

func (d *db_conn) GetTrial(id string) (*api.Trial, error) {
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

func (d *db_conn) GetTrialList(id string) ([]*api.Trial, error) {
	trials, err := d.getTrials("", id)

	return trials, err
}

func (d *db_conn) CreateTrial(trial *api.Trial) error {
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

	var trial_id string
	i := 3
	for true {
		trial_id = generate_randid()
		_, err = d.db.Exec("INSERT INTO trials VALUES (?, ?, ?, ?, ?)",
			trial_id, trial.StudyId, strings.Join(params, ",\n"),
			trial.ObjectiveValue, strings.Join(tags, ",\n"))
		if err == nil {
			trial.TrialId = trial_id
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

func (d *db_conn) DeleteTrial(id string) error {
	_, err := d.db.Exec("DELETE FROM trials WHERE id = ?", id)
	return err
}

func (d *db_conn) GetWorkerLogs(id string, opts *GetWorkerLogOpts) ([]*WorkerLog, error) {
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
		var time_str string

		err := rows.Scan(&time_str, &((*log1).Name), &((*log1).Value))
		if err != nil {
			log.Printf("Error scanning log: %v", err)
			continue
		}
		log1.Time, err = time.Parse(mysql_time_fmt, time_str)
		if err != nil {
			log.Printf("Error parsing time %s: %v", time_str, err)
			continue
		}
		result = append(result, log1)
	}
	return result, nil
}

func (d *db_conn) getWorkerLastlog(id string, value *string) (*time.Time, error) {
	var last_timestamp string
	var err error

	if value != nil {
		row := d.db.QueryRow("SELECT time, value FROM worker_lastlogs WHERE worker_id = ?", id)
		err = row.Scan(&last_timestamp, value)
	} else {
		row := d.db.QueryRow("SELECT time FROM worker_lastlogs WHERE worker_id = ?", id)
		err = row.Scan(&last_timestamp)
	}

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		mt, err := time.Parse(mysql_time_fmt, last_timestamp)
		if err != nil {
			log.Printf("Error parsing time in log %s: %v",
				last_timestamp, err)
			return nil, err
		}
		return &mt, nil
	}
}

func (d *db_conn) GetWorkerTimestamp(id string) (*time.Time, error) {
	return d.getWorkerLastlog(id, nil)
}

func (d *db_conn) storeWorkerLog(worker_id string, time string, metrics_name string, metrics_value string, objective_value_name string) error {
	is_objective := 0
	if metrics_name == objective_value_name {
		is_objective = 1
	}
	_, err := d.db.Exec("INSERT INTO worker_metrics (worker_id, time, name, value, is_objective) VALUES (?, ?, ?, ?, ?)",
		worker_id, time, metrics_name, metrics_value, is_objective)
	if err != nil {
		return err
	}
	return nil
}

func (d *db_conn) StoreWorkerLogs(worker_id string, logs []*api.MetricsLog) error {
	var lasterr error
	var last_value string

	db_t, err := d.getWorkerLastlog(worker_id, &last_value)
	if err != nil {
		log.Printf("Error getting last log timestamp: %v", err)
	}

	row := d.db.QueryRow("SELECT objective_value_name FROM workers "+
		"JOIN (studies) ON (workers.study_id = studies.id) WHERE "+
		"workers.id = ?", worker_id)
	var objective_value_name string
	err = row.Scan(&objective_value_name)
	if err != nil {
		log.Printf("Cannot get objective_value_name or metrics: %v", err)
		return err
	}

	var formatted_time string
	var ls []string
	for _, mlog := range logs {
		metrics_name := mlog.Name
		for _, mv := range mlog.Values {
			t, err := time.Parse(time.RFC3339Nano, mv.Time)
			if err != nil {
				log.Printf("Error parsing time %s: %v", mv.Time, err)
				lasterr = err
				continue
			}
			if db_t != nil && !t.After(*db_t) {
				// db_t is from mysql and has microsec precision.
				// This code assumes nanosec fractions are rounded down.
				continue
			}
			// use UTC as mysql DATETIME lacks timezone
			formatted_time = t.UTC().Format(mysql_time_fmt)
			if db_t != nil {
				// Parse again to get rounding effect
				//reparsed_time, err := time.Parse(mysql_time_fmt, formatted_time)
				//if reparsed_time == *db_t {
				//	if mv.Value == last_value {
				//	 stored_logs are already in DB
				//	 This assignment ensures the remaining
				//	 logs will be stored in DB.
				//		db_t = nil
				//		continue
				//	}
				//	// We don't know this is necessary or not yet.
				//	stored_logs = append(stored_logs, &mv.Value)
				//	continue
				//}
				// (reparsed_time > *db_t) can be assumed
				err = d.storeWorkerLog(worker_id,
					db_t.UTC().Format(mysql_time_fmt),
					metrics_name, mv.Value,
					objective_value_name)
				if err != nil {
					log.Printf("Error storing log %s: %v", mv.Value, err)
					lasterr = err
				}
				db_t = nil
			} else {
				err = d.storeWorkerLog(worker_id,
					formatted_time,
					metrics_name, mv.Value,
					objective_value_name)
				if err != nil {
					log.Printf("Error storing log %s: %v", mv.Value, err)
					lasterr = err
				}
			}
		}
	}
	if lasterr != nil {
		// If lastlog were updated, logs that couldn't be saved
		// would be lost.
		return lasterr
	}
	if len(ls) == 2 {
		_, err = d.db.Exec("REPLACE INTO worker_lastlogs VALUES (?, ?, ?)",
			worker_id, formatted_time, ls[1])
	}
	return err
}

func (d *db_conn) getWorkers(worker_id string, trial_id string, study_id string) ([]*api.Worker, error) {
	var rows *sql.Rows
	var err error

	if worker_id != "" {
		rows, err = d.db.Query("SELECT * FROM workers WHERE id = ?", worker_id)
	} else if trial_id != "" {
		rows, err = d.db.Query("SELECT * FROM workers WHERE trial_id = ?", trial_id)
	} else if study_id != "" {
		rows, err = d.db.Query("SELECT * FROM workers WHERE study_id = ?", study_id)
	} else {
		return nil, errors.New("worker_id, trial_id or study_id must be set")
	}

	if err != nil {
		return nil, err
	}

	var result []*api.Worker
	for rows.Next() {
		worker := new(api.Worker)

		var config, tags string
		err := rows.Scan(
			&worker.WorkerId,
			&worker.StudyId,
			&worker.TrialId,
			&worker.Type,
			&worker.Status,
			&config,
			&tags,
		)
		if err != nil {
			return nil, err
		}
		worker.Config = new(api.WorkerConfig)
		err = jsonpb.UnmarshalString(config, worker.Config)
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

func (d *db_conn) GetWorker(id string) (*api.Worker, error) {
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

func (d *db_conn) GetWorkerStatus(id string) (*api.State, error) {
	status := api.State_ERROR
	row := d.db.QueryRow("SELECT status FROM workers WHERE id = ?", id)
	err := row.Scan(&status)
	if err != nil {
		return &status, err
	}
	return &status, nil
}

func (d *db_conn) GetWorkerList(sid string, tid string) ([]*api.Worker, error) {
	workers, err := d.getWorkers("", tid, sid)
	return workers, err
}

func (d *db_conn) CreateWorker(worker *api.Worker) (string, error) {
	// Users should not overwrite worker.id
	var err, lastErr error
	config, err := (&jsonpb.Marshaler{}).MarshalToString(worker.Config)
	if err != nil {
		log.Fatalf("Error marshaling configs: %v", err)
		lastErr = err
	}

	tags := make([]string, len(worker.Tags))
	for i := range tags {
		tags[i], err = (&jsonpb.Marshaler{}).MarshalToString(worker.Tags[i])
		if err != nil {
			log.Printf("Error marshalling worker.Tags %v: %v",
				worker.Tags[i], err)
			lastErr = err
		}
	}

	var worker_id string
	i := 3
	for true {
		worker_id = generate_randid()
		_, err = d.db.Exec("INSERT INTO workers VALUES (?, ?, ?, ?, ?, ?, ?)",
			worker_id, worker.StudyId, worker.TrialId, worker.Type,
			api.State_PENDING, config, strings.Join(tags, ",\n"))
		if err == nil {
			worker.WorkerId = worker_id
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

func (d *db_conn) UpdateWorker(id string, newstatus api.State) error {
	_, err := d.db.Exec("UPDATE workers SET status = ? WHERE id = ?", newstatus, id)
	return err
}

func (d *db_conn) DeleteWorker(id string) error {
	_, err := d.db.Exec("DELETE FROM workers WHERE id = ?", id)
	return err
}

func (d *db_conn) SetSuggestionParam(algorithm string, studyId string, params []*api.SuggestionParameter) (string, error) {
	var err error
	ps := make([]string, len(params))
	for i, elem := range params {
		ps[i], err = (&jsonpb.Marshaler{}).MarshalToString(elem)
		if err != nil {
			log.Printf("Error marshalling %v: %v", elem, err)
			return "", err
		}
	}
	var paramId string
	for true {
		paramId = generate_randid()
		_, err = d.db.Exec("INSERT INTO suggestion_param VALUES (?, ?, ?, ?)",
			paramId, algorithm, studyId, strings.Join(ps, ",\n"))
		if err == nil {
			break
		}
	}
	return paramId, err
}

func (d *db_conn) UpdateSuggestionParam(paramId string, params []*api.SuggestionParameter) error {
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
		strings.Join(ps, ",\n"), paramId)
	return err
}

func (d *db_conn) GetSuggestionParam(paramId string) ([]*api.SuggestionParameter, error) {
	var params string
	row := d.db.QueryRow("SELECT parameters FROM suggestion_param WHERE id = ?", paramId)
	err := row.Scan(&params)
	if err != nil {
		return nil, err
	}
	var p_array []string
	if len(params) > 0 {
		p_array = strings.Split(params, ",\n")
	} else {
		return nil, nil
	}
	ret := make([]*api.SuggestionParameter, len(p_array))
	for i, j := range p_array {
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

func (d *db_conn) GetSuggestionParamList(studyId string) ([]*api.SuggestionParameterSet, error) {
	var rows *sql.Rows
	var err error
	rows, err = d.db.Query("SELECT * FROM suggestion_param WHERE study_id = ?", studyId)
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
		if studyId != sID {
			continue
		}
		var p_array []string
		if len(params) > 0 {
			p_array = strings.Split(params, ",\n")
		} else {
			return nil, nil
		}
		suggestparams := make([]*api.SuggestionParameter, len(p_array))
		for i, j := range p_array {
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

func (d *db_conn) SetEarlyStopParam(algorithm string, studyId string, params []*api.EarlyStoppingParameter) (string, error) {
	ps := make([]string, len(params))
	var err error
	for i, elem := range params {
		ps[i], err = (&jsonpb.Marshaler{}).MarshalToString(elem)
		if err != nil {
			log.Printf("Error marshalling %v: %v", elem, err)
			return "", err
		}
	}
	var paramId string
	for true {
		paramId := generate_randid()
		_, err = d.db.Exec("INSERT INTO earlystopping_param VALUES (?,?, ?, ?)",
			paramId, algorithm, studyId, strings.Join(ps, ",\n"))
		if err == nil {
			break
		}
	}
	return paramId, nil
}

func (d *db_conn) UpdateEarlyStopParam(paramId string, params []*api.EarlyStoppingParameter) error {
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
		strings.Join(ps, ",\n"), paramId)
	return err
}

func (d *db_conn) GetEarlyStopParam(paramId string) ([]*api.EarlyStoppingParameter, error) {
	var params string
	row := d.db.QueryRow("SELECT parameters FROM earlystopping_param WHERE id = ?", paramId)
	err := row.Scan(&params)
	if err != nil {
		return nil, err
	}
	var p_array []string
	if len(params) > 0 {
		p_array = strings.Split(params, ",\n")
	} else {
		return nil, nil
	}
	ret := make([]*api.EarlyStoppingParameter, len(p_array))
	for i, j := range p_array {
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

func (d *db_conn) GetEarlyStopParamList(studyId string) ([]*api.EarlyStoppingParameterSet, error) {
	var rows *sql.Rows
	var err error
	rows, err = d.db.Query("SELECT * FROM earlystopping_param WHERE study_id = ?", studyId)
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
		if studyId != sID {
			continue
		}
		var p_array []string
		if len(params) > 0 {
			p_array = strings.Split(params, ",\n")
		} else {
			return nil, nil
		}
		esparams := make([]*api.EarlyStoppingParameter, len(p_array))
		for i, j := range p_array {
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
