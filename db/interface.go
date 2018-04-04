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

	api "github.com/mlkube/katib/api"

	_ "github.com/go-sql-driver/mysql"
)

const (
	db_driver      = "mysql"
	db_name        = "root:test@tcp(vizier-db:3306)/vizier"
	mysql_time_fmt = "2006-01-02 15:04:05.999999"
)

type GetTrialLogOpts struct {
}

type TrialLog struct {
	Time  string
	Value string
}

type VizierDBInterface interface {
	DB_Init()
	GetStudyConfig(string) (*api.StudyConfig, error)
	GetStudyList() ([]string, error)
	CreateStudy(*api.StudyConfig) (string, error)
	DeleteStudy(string) error

	GetTrial(string) (*api.Trial, error)
	GetTrialStatus(string) (api.TrialState, error)
	GetTrialList(string) ([]*api.Trial, error)
	CreateTrial(*api.Trial) error
	UpdateTrial(string, api.TrialState) error
	GetTrialLogs(string, *GetTrialLogOpts) ([]*TrialLog, error)
	GetTrialTimestamp(string) (*time.Time, error)
	StoreTrialLogs(string, []string) error
	DeleteTrial(string) error
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
	var dummy_id, configs, suggestion_parameters, tags, metrics, command, mconf string
	err := row.Scan(&dummy_id,
		&study.Name,
		&study.Owner,
		&study.OptimizationType,
		&study.OptimizationGoal,
		&configs,
		&study.SuggestAlgorithm,
		&study.AutostopAlgorithm,
		&study.StudyTaskName,
		&suggestion_parameters,
		&tags,
		&study.ObjectiveValueName,
		&metrics,
		&study.Image,
		&command,
		&study.Gpu,
		&study.Scheduler,
		&mconf,
		&study.PullSecret,
	)
	if err != nil {
		return nil, err
	}
	study.ParameterConfigs = new(api.StudyConfig_ParameterConfigs)
	err = jsonpb.UnmarshalString(configs, study.ParameterConfigs)
	if err != nil {
		return nil, err
	}

	var sp_array []string
	if len(suggestion_parameters) > 0 {
		sp_array = strings.Split(suggestion_parameters, ",\n")
	}
	study.SuggestionParameters = make([]*api.SuggestionParameter, len(sp_array))
	for i, j := range sp_array {
		sp := new(api.SuggestionParameter)
		err = jsonpb.UnmarshalString(j, sp)
		if err != nil {
			log.Printf("err unmarshal %s", j)
			return nil, err
		}
		study.SuggestionParameters[i] = sp
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

	study.Mount = new(api.MountConf)
	if mconf != "" {
		err = jsonpb.UnmarshalString(mconf, study.Mount)
		if err != nil {
			return nil, err
		}
	}

	study.Metrics = strings.Split(metrics, ",\n")
	study.Command = strings.Split(command, ",\n")
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
	configs, err := (&jsonpb.Marshaler{}).MarshalToString(in.ParameterConfigs)
	if err != nil {
		log.Fatalf("Error marshaling configs: %v", err)
	}

	suggestion_parameters := make([]string, len(in.SuggestionParameters))
	for i, elem := range in.SuggestionParameters {
		suggestion_parameters[i], err = (&jsonpb.Marshaler{}).MarshalToString(elem)
		if err != nil {
			log.Printf("Error marshalling %v: %v", elem, err)
		}
	}
	var mconf string = ""
	if in.Mount != nil {
		mconf, err = (&jsonpb.Marshaler{}).MarshalToString(in.Mount)
		if err != nil {
			log.Fatalf("Error marshaling mount configs: %v", err)
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

	var study_id string
	i := 3
	for true {
		study_id = generate_randid()
		_, err := d.db.Exec(
			"INSERT INTO studies VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			study_id,
			in.Name,
			in.Owner,
			in.OptimizationType,
			in.OptimizationGoal,
			configs,
			in.SuggestAlgorithm,
			in.AutostopAlgorithm,
			in.StudyTaskName,
			strings.Join(suggestion_parameters, ",\n"),
			strings.Join(tags, ",\n"),
			in.ObjectiveValueName,
			strings.Join(in.Metrics, ",\n"),
			in.Image,
			strings.Join(in.Command, ",\n"),
			in.Gpu,
			in.Scheduler,
			mconf,
			in.PullSecret,
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
			&trial.Status,
			&trial.ObjectiveValue,
			&tags,
		)
		if err != nil {
			return nil, err
		}
		// XXX need to unmarshall parameters & tags
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

func (d *db_conn) GetTrialStatus(id string) (api.TrialState, error) {
	status := api.TrialState_ERROR

	row := d.db.QueryRow("SELECT status FROM trials WHERE id = ?", id)
	err := row.Scan(&status)
	if err != nil {
		return status, err
	}
	return status, nil
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
		_, err = d.db.Exec("INSERT INTO trials VALUES (?, ?, ?, ?, ?, ?)",
			trial_id, trial.StudyId, strings.Join(params, ",\n"),
			trial.Status, trial.ObjectiveValue, strings.Join(tags, ",\n"))
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

func (d *db_conn) UpdateTrial(id string, newstatus api.TrialState) error {
	_, err := d.db.Exec("UPDATE trials SET status = ? WHERE id = ?", newstatus, id)
	return err
}

func (d *db_conn) GetTrialLogs(id string, opts *GetTrialLogOpts) ([]*TrialLog, error) {
	// TODO: opts not implemented
	rows, err := d.db.Query("SELECT (time, value) FROM trial_logs WHERE trial_id = ? ORDER BY time", id)
	if err != nil {
		return nil, err
	}

	var result []*TrialLog
	for rows.Next() {
		log1 := new(TrialLog)

		err := rows.Scan(&((*log1).Time), &((*log1).Value))
		if err != nil {
			log.Printf("Error scanning log: %v", err)
			continue
		}
		result = append(result, log1)
	}
	return result, nil
}

func (d *db_conn) GetTrialTimestamp(id string) (*time.Time, error) {
	var last_timestamp string

	row := d.db.QueryRow("SELECT time FROM trial_logs WHERE trial_id = ? ORDER BY time DESC LIMIT 1", id)
	err := row.Scan(&last_timestamp)
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

func (d *db_conn) StoreTrialLogs(trial_id string, logs []string) error {
	var lasterr error

	for _, logline := range logs {
		if logline == "" {
			continue
		}
		ls := strings.SplitN(logline, " ", 2)
		if len(ls) != 2 {
			log.Printf("Error parsing log: %s", logline)
			lasterr = errors.New("Error parsing log")
			continue
		}
		t, err := time.Parse(time.RFC3339Nano, ls[0])
		if err != nil {
			log.Printf("Error parsing time %s: %v", ls[0], err)
			lasterr = err
			continue
		}
		// use UTC as mysql DATETIME lacks timezone
		_, err = d.db.Exec("INSERT INTO trial_logs VALUES (?, ?, ?)",
			trial_id, t.UTC().Format("2006-01-02 15:04:05.999999"), ls[1])
		if err != nil {
			log.Printf("Error storing log %s: %v", logline, err)
			lasterr = err
		}
	}
	return lasterr
}

func (d *db_conn) DeleteTrial(id string) error {
	_, err := d.db.Exec("DELETE FROM trials WHERE id = ?", id)
	return err
}
