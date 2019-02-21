package db

import (
	crand "crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"

	api "github.com/kubeflow/katib/pkg/api"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbDriver     = "mysql"
	dbNameTmpl   = "root:%s@tcp(vizier-db:3306)/vizier?timeout=5s"
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

type VizierDBInterface interface {
	DBInit()
	SelectOne() error

	GetStudy(string) (*api.StudyConfig, error)
	GetStudyList() ([]string, error)

	CreateStudy(*api.StudyConfig) (string, error)
	UpdateStudy(string, *api.StudyConfig) error
	DeleteStudy(string) error

	GetTrial(string) (*api.Trial, error)
	GetTrialList(string) ([]*api.Trial, error)
	CreateTrial(*api.Trial) error
	UpdateTrial(*api.Trial) error
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
	GetWorkerFullInfo(string, string, string, bool) (*api.GetWorkerFullInfoReply, error)

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

func getDbName() string {
	dbPass := os.Getenv("MYSQL_ROOT_PASSWORD")
	if dbPass == "" {
		log.Printf("WARN: Env var MYSQL_ROOT_PASSWORD is empty. Falling back to \"test\".")

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

func NewWithSQLConn(db *sql.DB) (VizierDBInterface, error) {
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

func (d *dbConn) GetStudyMetrics(id string) ([]string, error) {
	row := d.db.QueryRow("SELECT metrics FROM studies WHERE id = ?", id)
	var metrics string
	err := row.Scan(&metrics)
	if err != nil {
		return nil, err
	}
	retMetrics := strings.Split(metrics, ",\n")
	return retMetrics, nil
}

func (d *dbConn) GetStudy(StudyID string) (*api.StudyConfig, error) {
	row := d.db.QueryRow("SELECT * FROM studies WHERE id = ?", StudyID)
	study := new(api.StudyConfig)
	var dummyID, nasConfig, parameters, tags, metrics string
	err := row.Scan(&dummyID,
		&study.Name,
		&study.Owner,
		&study.OptimizationType,
		&study.OptimizationGoal,
		&parameters,
		&tags,
		&study.ObjectiveValueName,
		&metrics,
		&nasConfig,
		&study.JobId,
		&study.JobType,
	)
	if err != nil {
		return nil, err
	}
	if parameters != "" {
		study.ParameterConfigs = new(api.StudyConfig_ParameterConfigs)
		err = jsonpb.UnmarshalString(parameters, study.ParameterConfigs)
		if err != nil {
			return nil, err
		}
	}
	if nasConfig != "" {
		study.NasConfig = new(api.NasConfig)
		err = jsonpb.UnmarshalString(nasConfig, study.NasConfig)
		if err != nil {
			log.Printf("Failed to unmarshal NasConfig")
			return nil, err
		}
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

func New() (VizierDBInterface, error) {
	db, err := openSQLConn(dbDriver, getDbName(), connectInterval, connectTimeout)
	if err != nil {
		return nil, fmt.Errorf("DB open failed: %v", err)
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

func isDBDuplicateError(err error) bool {
	errmsg := strings.ToLower(err.Error())
	if strings.Contains(errmsg, "unique") || strings.Contains(errmsg, "duplicate") {
		return true
	}
	return false
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
	if in.JobId != "" {
		var temporaryId string
		err := d.db.QueryRow("SELECT id FROM studies WHERE job_id = ?", in.JobId).Scan(&temporaryId)
		if err == nil {
			return "", fmt.Errorf("Study %s in Job %s already exist.", in.Name, in.JobId)
		}
	}

	var nasConfig string
	var configs string
	var err error
	if in.NasConfig != nil {
		nasConfig, err = (&jsonpb.Marshaler{}).MarshalToString(in.NasConfig)
		if err != nil {
			log.Fatalf("Error marshaling nasConfig: %v", err)
		}
	}

	if in.ParameterConfigs != nil {
		configs, err = (&jsonpb.Marshaler{}).MarshalToString(in.ParameterConfigs)
		if err != nil {
			log.Fatalf("Error marshaling configs: %v", err)
		}
	}

	tags := make([]string, len(in.Tags))
	for i, elem := range in.Tags {
		tags[i], err = (&jsonpb.Marshaler{}).MarshalToString(elem)
		if err != nil {
			log.Printf("Error marshalling %v: %v", elem, err)
			continue
		}
	}

	/* WE PROBABLY DON'T NEED METRICS AND ALSO THIS LOGIC IS KIND OF CONFUSING */
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
			"INSERT INTO studies VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			studyID,
			in.Name,
			in.Owner,
			in.OptimizationType,
			in.OptimizationGoal,
			configs,
			strings.Join(tags, ",\n"),
			in.ObjectiveValueName,
			strings.Join(in.Metrics, ",\n"),
			nasConfig,
			in.JobId,
			in.JobType,
		)
		if err == nil {
			break
		} else if isDBDuplicateError(err) {
			i--
			if i > 0 {
				continue
			}
		}
		return "", err
	}

	return studyID, nil
}

// UpdateStudy updates the corresponding row in the DB.
// It only updates name, owner, tags and job_id.
// Other columns are silently ignored.
func (d *dbConn) UpdateStudy(studyID string, in *api.StudyConfig) error {

	/* THINK ABOUT TRIALS */
	var err error

	tags := make([]string, len(in.Tags))
	for i, elem := range in.Tags {
		tags[i], err = (&jsonpb.Marshaler{}).MarshalToString(elem)
		if err != nil {
			log.Printf("Error marshalling %v: %v", elem, err)
			continue
		}
	}
	_, err = d.db.Exec(`UPDATE studies SET name = ?, owner = ?, tags = ?,
                job_id = ? WHERE id = ?`,
		in.Name,
		in.Owner,
		strings.Join(tags, ",\n"),
		in.JobId,
		studyID)
	return err
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
		var timeStamp string
		err := rows.Scan(&trial.TrialId,
			&trial.StudyId,
			&parameters,
			&trial.ObjectiveValue,
			&tags,
			&timeStamp,
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

func marshalTrial(trial *api.Trial) ([]string, []string, error) {
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
	return params, tags, lastErr
}

// CreateTrial stores into the trials DB table.
// As a side-effect, it generates and sets trial.TrialId.
// Users should not overwrite TrialId.

/* TODO FIX CREATE_TRIAL & DELETE_TRIAL IN CASE OF NAS, SINCE WE WANT TRIALS */
func (d *dbConn) CreateTrial(trial *api.Trial) error {
	params, tags, lastErr := marshalTrial(trial)

	var trialID string
	i := 3
	for true {
		trialID = generateRandid()
		timeString := time.Now().UTC().Format(mysqlTimeFmt)
		_, err := d.db.Exec("INSERT INTO trials VALUES (?, ?, ?, ?, ?, ?)",
			trialID, trial.StudyId, strings.Join(params, ",\n"),
			trial.ObjectiveValue, strings.Join(tags, ",\n"), timeString)
		if err == nil {
			trial.TrialId = trialID
			break
		} else if isDBDuplicateError(err) {
			i--
			if i > 0 {
				continue
			}
		}
		return err
	}
	return lastErr
}

// UpdateTrial updates the corresponding row in the DB.
// It only updates parameters and tags. Other columns are silently ignored.
func (d *dbConn) UpdateTrial(trial *api.Trial) error {
	params, tags, lastErr := marshalTrial(trial)
	_, err := d.db.Exec(`UPDATE trials SET parameters = ?, tags = ?,
		WHERE id = ?`,
		strings.Join(params, ",\n"), strings.Join(tags, ",\n"),
		trial.TrialId)
	if err != nil {
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
		} else if isDBDuplicateError(err) {
			i--
			if i > 0 {
				continue
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

func (d *dbConn) GetWorkerFullInfo(studyId string, trialId string, workerId string, OnlyLatestLog bool) (*api.GetWorkerFullInfoReply, error) {
	ret := &api.GetWorkerFullInfoReply{}
	var err error
	ws := []*api.Worker{}

	if workerId != "" {
		w, err := d.GetWorker(workerId)
		ws = append(ws, w)
		if err != nil {
			return ret, err
		}
	} else {
		ws, err = d.GetWorkerList(studyId, trialId)
		if err != nil {
			return ret, err
		}
	}
	ts, err := d.GetTrialList(studyId)
	if err != nil {
		return ret, err
	}
	// Actually no need to get full config now
	metrics, err := d.GetStudyMetrics(studyId)
	if err != nil {
		return ret, err
	}

	plist := make(map[string][]*api.Parameter)
	for _, t := range ts {
		plist[t.TrialId] = t.ParameterSet
	}

	wfilist := make([]*api.WorkerFullInfo, len(ws))
	var qstr, id string
	if OnlyLatestLog {
		qstr = `
		SELECT 
			WM.worker_id, WM.time, WM.name, WM.value 
		FROM (
			SELECT 
				Master.worker_id, Master.time,  Master.name,  Master.value
			FROM (
				SELECT 
					worker_id, name, 
					MAX(id) AS MaxID
				FROM 
					worker_metrics 
				GROUP BY 
					worker_id, name
				) AS LATEST
				JOIN worker_metrics AS Master
				ON Master.id = LATEST.MaxID
		) AS WM 
		JOIN workers AS WS 
		ON WM.worker_id = WS.id 
		AND`
	} else {
		qstr = `
		SELECT 
			WM.worker_id, WM.time, WM.name, WM.value 
		FROM 
			worker_metrics AS WM 
		JOIN workers AS WS 
		ON WM.worker_id = WS.id 
		AND`
	}
	if workerId != "" {
		if OnlyLatestLog {
			qstr = `
			SELECT 
			WM.worker_id, WM.time, WM.name, WM.value 
			FROM (
				SELECT 
					Master.worker_id, Master.time,  Master.name,  Master.value
				FROM (
					SELECT 
						worker_id, name, 
						MAX(id) AS MaxID
					FROM 
						worker_metrics 
					GROUP BY 
						worker_id, name
				) AS LATEST
				JOIN worker_metrics AS Master
				ON Master.id = LATEST.MaxID
				AND Master.worker_id = ?
			) AS WM`
		} else {
			qstr = "SELECT worker_id, time, name, value FROM worker_metrics WHERE worker_id = ?"
		}
		id = workerId
	} else if trialId != "" {
		qstr += " WS.trial_id = ? "
		id = trialId
	} else if studyId != "" {
		qstr += " WS.study_id = ? "
		id = studyId
	}
	rows, err := d.db.Query(qstr+" ORDER BY time", id)
	if err != nil {
		log.Printf("SQL query: %v", err)
		return ret, err
	}
	metricslist := make(map[string]map[string][]*api.MetricsValueTime, len(ws))
	for rows.Next() {
		var name, value, timeStr, wid string
		err := rows.Scan(&wid, &timeStr, &name, &value)
		if err != nil {
			log.Printf("Error scanning log: %v", err)
			continue
		}
		ptime, err := time.Parse(mysqlTimeFmt, timeStr)
		if err != nil {
			log.Printf("Error parsing time %s: %v", timeStr, err)
			continue
		}
		if _, ok := metricslist[wid]; ok {
			metricslist[wid][name] = append(metricslist[wid][name], &api.MetricsValueTime{
				Value: value,
				Time:  ptime.UTC().Format(time.RFC3339Nano),
			})
		} else {
			metricslist[wid] = make(map[string][]*api.MetricsValueTime, len(metrics))
			metricslist[wid][name] = append(metricslist[wid][name], &api.MetricsValueTime{
				Value: value,
				Time:  ptime.UTC().Format(time.RFC3339Nano),
			})
		}
	}
	for i, w := range ws {
		wfilist[i] = &api.WorkerFullInfo{
			Worker:       w,
			ParameterSet: plist[w.TrialId],
		}
		for _, m := range metrics {
			if v, ok := metricslist[w.WorkerId][m]; ok {
				wfilist[i].MetricsLogs = append(wfilist[i].MetricsLogs, &api.MetricsLog{
					Name:   m,
					Values: v,
				},
				)
			}
		}
	}
	ret.WorkerFullInfos = wfilist
	return ret, nil
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
		} else if !isDBDuplicateError(err) {
			return "", err
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
	rows, err = d.db.Query("SELECT id, suggestion_algo, parameters FROM suggestion_param WHERE study_id = ?", studyID)
	if err != nil {
		return nil, err
	}
	var result []*api.SuggestionParameterSet
	for rows.Next() {
		var id string
		var algorithm string
		var params string
		err := rows.Scan(&id, &algorithm, &params)
		if err != nil {
			return nil, err
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
		paramID = generateRandid()
		_, err = d.db.Exec("INSERT INTO earlystopping_param VALUES (?,?, ?, ?)",
			paramID, algorithm, studyID, strings.Join(ps, ",\n"))
		if err == nil {
			break
		} else if !isDBDuplicateError(err) {
			return "", err
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
	rows, err = d.db.Query("SELECT id, earlystop_algo, parameters FROM earlystopping_param WHERE study_id = ?", studyID)
	if err != nil {
		return nil, err
	}
	var result []*api.EarlyStoppingParameterSet
	for rows.Next() {
		var id string
		var algorithm string
		var params string
		err := rows.Scan(&id, &algorithm, &params)
		if err != nil {
			return nil, err
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
