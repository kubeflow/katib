package db

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/jsonpb"

	api "github.com/kubeflow/katib/pkg/api"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var dbInterface, mysqlInterface VizierDBInterface
var mock sqlmock.Sqlmock

var studyColumns = []string{
	"id", "name", "owner", "optimization_type", "optimization_goal",
	"parameter_configs", "tags", "objective_value_name",
	"metrics", "nas_config", "job_id", "job_type"}
var trialColumns = []string{
	"id", "study_id", "parameters", "objective_value", "tags", "time"}
var workerColumns = []string{"id",
	"study_id", "trial_id", "type",
	"status", "template_path", "tags"}

func TestMain(m *testing.M) {
	db, sm, err := sqlmock.New()
	mock = sm
	if err != nil {
		fmt.Printf("error opening db: %v\n", err)
		os.Exit(1)
	}
	//mock.ExpectBegin()
	dbInterface, err = NewWithSQLConn(db)
	if err != nil {
		fmt.Printf("error NewWithSQLConn: %v\n", err)
	}
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS studies").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS study_permissions").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS trials").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS workers").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS worker_metrics").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS worker_lastlogs").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS suggestion_param").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS earlystop_param").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	dbInterface.DBInit()
	err = dbInterface.SelectOne()
	if err != nil {
		fmt.Printf("error `SELECT 1` probing: %v\n", err)
	}

	mysqlAddr := os.Getenv("TEST_MYSQL")
	if mysqlAddr != "" {
		mysql, err := sql.Open("mysql", "root:test123@tcp("+mysqlAddr+")/vizier")
		if err != nil {
			fmt.Printf("error opening db: %v\n", err)
			os.Exit(1)
		}

		mysqlInterface, err = NewWithSQLConn(mysql)
		if err != nil {
			fmt.Printf("error initializing db interface: %v\n", err)
			os.Exit(1)
		}

		mysqlInterface.DBInit()
	}

	os.Exit(m.Run())
}

func TestOpenSQLConn(t *testing.T) {
	_, _, err := sqlmock.New()
	if err != nil {
		fmt.Printf("error opening db: %v\n", err)
		os.Exit(1)
	}
	mysqlAddr := os.Getenv("TEST_MYSQL")
	if mysqlAddr != "" {
		_, err := openSQLConn("mysql", "root:test123@tcp("+mysqlAddr+")/vizier", time.Second, 3*time.Second)
		if err != nil {
			t.Errorf("openSQLConn error: %v", err)
		}
	}
	_, err = openSQLConn("mysql", "root:test123@tcp(dummy)/vizier", time.Second, 3*time.Second)
	if err.Error() != "Timeout waiting for DB conn successfully opened." {
		t.Errorf("openSQLConn should timeout but got error: %v", err)
	}
}

func TestCreateStudy(t *testing.T) {
	var in api.StudyConfig
	in.ParameterConfigs = new(api.StudyConfig_ParameterConfigs)
	//err := jsonpb.UnmarshalString("{}", &in)
	err := jsonpb.UnmarshalString(`{"configs": [{"name": "-abc"}]}`, in.ParameterConfigs)
	if err != nil {
		t.Errorf("err %v", err)
	}

	mock.ExpectExec("INSERT INTO studies VALUES").WithArgs().WillReturnError(errors.New("sql: Duplicated key"))
	mock.ExpectExec("INSERT INTO studies VALUES").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	id, err := dbInterface.CreateStudy(&in)
	if err != nil {
		t.Errorf("CreateStudy error %v", err)
	} else if len(id) != 16 {
		t.Errorf("CreateStudy returned incorrect ID %s", id)
	}
}

func TestGetStudyConfig(t *testing.T) {
	id := generateRandid()
	mock.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows(studyColumns).AddRow(
			"abc", "test", "admin", 1, 0.99, "{}", "", "", "", "", "test", "hp"))
	study, err := dbInterface.GetStudy(id)
	if err != nil {
		t.Errorf("GetStudyConfig failed: %v", err)
	} else if study.Name != "test" || study.Owner != "admin" {
		t.Errorf("GetStudyConfig incorrect return %v", study)
	}
}

func TestGetStudyList(t *testing.T) {
	ids := []string{"abcde1234567890f", "bcde1234567890fa"}
	mock.ExpectQuery("SELECT id FROM studies").WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(ids[0]).AddRow(ids[1]))
	r, err := dbInterface.GetStudyList()
	if err != nil {
		t.Errorf("GetStudyList error %v", err)
	}
	if len(r) != len(ids) {
		t.Errorf("GetStudyList returned incorrect number of ids %d != %d",
			len(r), len(ids))
	}
	for i, id := range r {
		if ids[i] != id {
			t.Errorf("GetStudyList returned incorrect ID %s != %s",
				id, ids[i])
		}
	}
}

func TestUpdateStudy(t *testing.T) {
	studyID := generateRandid()
	var in api.StudyConfig
	in.Name = "hoge"
	in.Owner = "joe"
	in.JobId = "foobar123"

	mock.ExpectExec(`UPDATE studies SET name = \?, owner = \?, tags = \?,
                job_id = \? WHERE id = \?`,
	).WithArgs(in.Name, in.Owner, "", in.JobId, studyID).WillReturnResult(sqlmock.NewResult(1, 1))
	err := dbInterface.UpdateStudy(studyID, &in)
	if err != nil {
		t.Errorf("UpdateStudy error %v", err)
	}
}

func TestDeleteStudy(t *testing.T) {
	studyID := generateRandid()
	mock.ExpectExec(`DELETE FROM studies WHERE id = \?`).WithArgs(studyID).WillReturnResult(sqlmock.NewResult(1, 1))
	err := dbInterface.DeleteStudy(studyID)
	if err != nil {
		t.Errorf("DeleteStudy error %v", err)
	}
}

func TestCreateStudyIdGeneration(t *testing.T) {
	if mysqlInterface == nil {
		t.Skip("TEST_MYSQL is not defined.")
	}
	var in api.StudyConfig
	in.ParameterConfigs = new(api.StudyConfig_ParameterConfigs)

	seed := rand.Int63()
	encountered := map[string]bool{}
	for i := 0; i < 4; i++ {
		// Repeadedly use the seed to force the same ID generation
		rand.Seed(seed)
		id, err := mysqlInterface.CreateStudy(&in)
		if i == 3 {
			if err == nil || !isDBDuplicateError(err) {
				t.Errorf("Expected an duplicate error but got %v",
					err)
			} else {
				break
			}
		}
		if err != nil {
			t.Errorf("CreateStudy error %v", err)
		} else if !encountered[id] {
			encountered[id] = true
		} else {
			t.Fatalf("Study ID duplicated %s", id)
		}
		t.Logf("id gen %d %s %v\n", i, id, err)
	}
	for id, _ := range encountered {
		err := mysqlInterface.DeleteStudy(id)
		if err != nil {
			t.Errorf("DeleteStudy error %v", err)
		}
	}
}

func TestGetTrial(t *testing.T) {
	id := generateRandid()
	mock.ExpectQuery(`SELECT \* FROM trials WHERE id = \?`).WillReturnRows(
		sqlmock.NewRows(trialColumns).AddRow(
			id, "s1234567890abcde",
			"{\"name\": \"1\"},\n{}", "obj_val",
			"{\"name\": \"foo\"},\n{}",
			""))
	trial, err := dbInterface.GetTrial(id)
	if err != nil {
		t.Errorf("GetTrial error %v", err)
	} else if len((*trial).Tags) != 2 {
		t.Errorf("GetTrial returned incorrect Tag %v", (*trial).Tags)
	}
}

func TestGetTrialList(t *testing.T) {
	studyID := generateRandid()
	var ids = []string{"abcdef1234567890", "bcdef1234567890a"}
	rows := sqlmock.NewRows(trialColumns)
	for _, id := range ids {
		rows.AddRow(id, studyID, "", "obj_val", "", "")
	}
	mock.ExpectQuery(`SELECT \* FROM trials WHERE study_id = \?`).WithArgs(studyID).WillReturnRows(rows)
	trials, err := dbInterface.GetTrialList(studyID)
	if err != nil {
		t.Errorf("GetTrialList error %v", err)
	} else if len(trials) != len(ids) {
		t.Errorf("GetTrialList returned incorrect number of trials %d != %d",
			len(trials), len(ids))
	}
}

// don't know how to fix since time is automatically generate
// func TestCreateTrial(t *testing.T) {
// 	var trial api.Trial
// 	trial.StudyId = generateRandid()
// 	mock.ExpectExec(`INSERT INTO trials VALUES \(`).WithArgs(sqlmock.AnyArg(),
// 		trial.StudyId, "", "", "", "").WillReturnResult(sqlmock.NewResult(1, 1))
// 	err := dbInterface.CreateTrial(&trial)
// 	if err != nil {
// 		t.Errorf("CreateTrial error %v", err)
// 	}
// }

func TestUpdateTrial(t *testing.T) {
	var trial api.Trial
	trial.TrialId = generateRandid()
	trial.StudyId = generateRandid()
	trial.ParameterSet = make([]*api.Parameter, 1)
	trial.ParameterSet[0] = &api.Parameter{Name: "abc"}
	mock.ExpectExec(`UPDATE trials SET parameters = \?, tags = \?,
		WHERE id = \?`,
	).WithArgs("{\"name\":\"abc\"}", "", trial.TrialId).WillReturnResult(sqlmock.NewResult(1, 1))
	err := dbInterface.UpdateTrial(&trial)
	if err != nil {
		t.Errorf("UpdateTrial error %v", err)
	}
}

func TestDeleteTrial(t *testing.T) {
	id := generateRandid()
	mock.ExpectExec(`DELETE FROM trials WHERE id = \?`).WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))
	err := dbInterface.DeleteTrial(id)
	if err != nil {
		t.Errorf("DeleteTrial error %v", err)
	}
}

func TestCreateWorker(t *testing.T) {
	var w api.Worker
	w.StudyId = generateRandid()
	w.TrialId = generateRandid()

	mock.ExpectExec("INSERT INTO workers VALUES").WithArgs(sqlmock.AnyArg(), w.StudyId, w.TrialId, w.Type, api.State_PENDING, w.TemplatePath, "").WillReturnResult(sqlmock.NewResult(1, 1))
	worker_id, err := dbInterface.CreateWorker(&w)

	if err != nil {
		t.Errorf("CreateWorker error %v", err)
	} else if worker_id != w.WorkerId {
		t.Errorf("Worker ID doesn't match %s != %s",
			worker_id, w.WorkerId)
	}
}

const defaultWorkerID = "w123456789abcdef"
const objValueName = "obj_value"

func TestGetWorker(t *testing.T) {
	mock.ExpectQuery(`SELECT \* FROM workers WHERE id = \?`).WithArgs(defaultWorkerID).WillReturnRows(sqlmock.NewRows(workerColumns).AddRow(defaultWorkerID, 1, 1, 1, 1, 1, ""))
	worker, err := dbInterface.GetWorker(defaultWorkerID)
	if err != nil {
		t.Errorf("GetWorker error %v", err)
	} else if worker.WorkerId != defaultWorkerID {
		t.Errorf("GetWorker returned incorrect ID %s", worker.WorkerId)
	}
}

func TestGetWorkerStatus(t *testing.T) {
	mock.ExpectQuery(`SELECT status FROM workers WHERE id = \?`).WithArgs(defaultWorkerID).WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow(api.State_RUNNING))
	status, err := dbInterface.GetWorkerStatus(defaultWorkerID)
	if err != nil {
		t.Errorf("GetWorker error %v", err)
	} else if *status != api.State_RUNNING {
		t.Errorf("GetWorkerStatus returned incorrect %s", *status)
	}
}

func TestGetWorkerList(t *testing.T) {
	var trial_id = generateRandid()
	mock.ExpectQuery(`SELECT \* FROM workers WHERE trial_id = \?`).WithArgs(trial_id).WillReturnRows(sqlmock.NewRows(workerColumns).AddRow(1, 1, 1, 1, 1, 1, ""))
	out_workers, err := dbInterface.GetWorkerList("", trial_id)
	if err != nil {
		t.Errorf("GetWorkerList error %v", err)
	} else if len(out_workers) != 1 {
		t.Errorf("GetWorkerList returned incorrect number of workers %d",
			len(out_workers))
	}

	var study_id = generateRandid()
	mock.ExpectQuery(`SELECT \* FROM workers WHERE study_id = \?`).WithArgs(study_id).WillReturnRows(sqlmock.NewRows(workerColumns).AddRow(1, 1, 1, 1, 1, 1, "").AddRow(1, 1, 1, 1, 1, 1, ""))
	out_workers, err = dbInterface.GetWorkerList(study_id, "")
	if err != nil {
		t.Errorf("GetWorkerList error %v", err)
	} else if len(out_workers) != 2 {
		t.Errorf("GetWorkerList returned incorrect number of workers %d",
			len(out_workers))
	}
}

func TestUpdateWorker(t *testing.T) {
	mock.ExpectExec(`UPDATE workers SET status = \? WHERE id = \?`).WithArgs(api.State_COMPLETED, defaultWorkerID).WillReturnResult(sqlmock.NewResult(1, 1))
	err := dbInterface.UpdateWorker(defaultWorkerID, api.State_COMPLETED)
	if err != nil {
		t.Errorf("UpdateWorker error %v", err)
	}
}
func TestDeleteWorker(t *testing.T) {
	mock.ExpectExec(`DELETE FROM workers WHERE id = \?`).WithArgs(defaultWorkerID).WillReturnResult(sqlmock.NewResult(1, 1))
	err := dbInterface.DeleteWorker(defaultWorkerID)
	if err != nil {
		t.Errorf("DeleteWorker error %v", err)
	}

}

func TestGetWorkerFullInfo(t *testing.T) {
	studyID := generateRandid()
	wRows := sqlmock.NewRows(workerColumns)
	wRows.AddRow("w1134567890abcde", studyID, "", "", "1", "", "")
	wRows.AddRow("w2234567890abcde", studyID, "", "", "2", "", "")
	mock.ExpectQuery(`SELECT \* FROM workers WHERE study_id = \?`).WithArgs(studyID).WillReturnRows(wRows)
	mock.ExpectQuery(`SELECT \* FROM trials WHERE study_id = \?`).WithArgs(studyID).WillReturnRows(
		sqlmock.NewRows(trialColumns))
	mock.ExpectQuery(`SELECT metrics FROM studies WHERE id = \?`).WithArgs(studyID).WillReturnRows(sqlmock.NewRows([]string{"metrics"}).AddRow("foo,\nbar"))

	WMRows := sqlmock.NewRows([]string{"WM.worker_id", "WM.time", "WM.name", "WM.value"})
	WMRows.AddRow("w1134567890abcde", "2012-01-01 09:54:32", "foo", "1")
	WMRows.AddRow("w1134567890abcde", "2012-01-01 09:54:32", "bar", "1")
	mock.ExpectQuery(`SELECT WM.worker_id, WM.time, WM.name, WM.value FROM .* MaxID .* ON WM.worker_id`).WithArgs(studyID).WillReturnRows(WMRows)

	fi, err := dbInterface.GetWorkerFullInfo(studyID, "", "", true)
	if err != nil {
		t.Errorf("GetWorkerFullInfo error %v", err)
	} else if len(fi.WorkerFullInfos) != 2 ||
		len(fi.WorkerFullInfos[0].MetricsLogs) != 2 {
		t.Errorf("GetWorkerFullInfo incorrect return  %v", fi)
	}
}

type MetricsLogData struct {
	stored bool
	name   string
	time   string
}

func newMetricsLog(ms []MetricsLogData) []*api.MetricsLog {
	mlog := make([]*api.MetricsLog, len(ms))
	for i, m := range ms {
		value := fmt.Sprintf("%d", i)
		mlog[i] = &api.MetricsLog{
			Name: m.name, Values: []*api.MetricsValueTime{
				{Time: m.time, Value: value}}}
		if m.stored {
			t, _ := time.Parse(time.RFC3339Nano, m.time)
			timeStr := t.UTC().Format(mysqlTimeFmt)
			var isObj int64
			if m.name == objValueName {
				isObj = 1
			}
			ex := mock.ExpectExec(
				`INSERT INTO worker_metrics \(worker_id, time, name, value, is_objective\)`)
			ex.WithArgs(defaultWorkerID, timeStr, m.name, value, isObj).WillReturnResult(sqlmock.NewResult(1, 1))
		}
	}
	return mlog
}

func TestStoreWorkerLogs(t *testing.T) {
	var tests = []struct {
		lastMetrics  [][]interface{}
		newMetrics   []MetricsLogData
		newTimestamp string
	}{
		{
			[][]interface{}{
				{"2012-01-02 19:54:31.999999", "foo", "2"},
				{"2012-01-02 19:54:31.999999", "baz", "4"}},
			[]MetricsLogData{
				{false, "foo", "2012-01-02T09:54:31.995555Z"},
				{true, "bar", "2012-01-02T19:54:31.999999Z"},
				{false, "foo", "2012-01-02T19:54:31.999999Z"},
				{true, "bar", "2012-01-02T19:54:31.999999Z"},
				{false, "baz", "2012-01-02T19:54:31.999999Z"},
				{true, "obj_value", "2012-01-02T21:54:34.1234+02:00"},
				{true, "hoge", "2012-01-02T19:54:33Z"}},
			"2012-01-02 19:54:34.1234",
		},
		{
			[][]interface{}{
				{"2012-01-02 20:54:31.999999", nil, nil}},
			[]MetricsLogData{
				{false, "foo", "2012-01-02T09:54:31.995555Z"},
				{true, "bar", "2012-01-02T20:54:31.99999901Z"},
				{true, "foo", "2012-01-02T20:54:31.99999902Z"},
				{true, "baz", "2012-01-02T20:54:31.99999903Z"},
				{true, "obj_value", "2012-01-02T20:54:32Z"},
			},
			"2012-01-02 20:54:32",
		},
		{
			[][]interface{}{},
			[]MetricsLogData{
				{true, "foo", "2012-01-02T09:54:31.995555Z"},
				{true, "baz", "2012-01-02T20:54:31.99999903Z"},
				{true, "obj_value", "2012-01-02T20:54:32Z"},
			},
			"2012-01-02 20:54:32",
		},
	}

	for i, test := range tests {
		rows := sqlmock.NewRows([]string{"worker_lastlogs.time", "name", "value"})
		for _, r := range test.lastMetrics {
			rows.AddRow(r[0], r[1], r[2])
		}

		mock.ExpectQuery(`SELECT .* FROM worker_lastlogs
                        LEFT JOIN worker_metrics ON .* WHERE worker_lastlogs.worker_id = \?`).WithArgs(defaultWorkerID).WillReturnRows(rows)
		mock.ExpectQuery(`SELECT objective_value_name FROM workers
			JOIN \(studies\) ON \(workers.study_id = studies.id\) WHERE
			workers.id = \?`).WithArgs(defaultWorkerID).WillReturnRows(
			sqlmock.NewRows([]string{"objective_value_name"}).AddRow(
				objValueName))

		mlogs := newMetricsLog(test.newMetrics)

		mock.ExpectExec("REPLACE INTO worker_lastlogs").WithArgs(defaultWorkerID, test.newTimestamp).WillReturnResult(sqlmock.NewResult(1, 1))
		err := dbInterface.StoreWorkerLogs(defaultWorkerID, mlogs)
		if err != nil {
			t.Errorf("StoreWorkerLogs test %d error %v", i, err)
		}
	}
}

func TestGetWorkerTimestamp(t *testing.T) {
	timeStr := "2012-01-02 20:54:32.123456"
	timeVal, _ := time.Parse(mysqlTimeFmt, timeStr)

	mock.ExpectQuery(`SELECT time FROM worker_lastlogs WHERE worker_id = \?`).WithArgs(defaultWorkerID).WillReturnRows(
		sqlmock.NewRows([]string{"time"}).AddRow(timeStr))
	tm, err := dbInterface.GetWorkerTimestamp(defaultWorkerID)
	if err != nil {
		t.Errorf("GetWorkerTimestamp error %v", err)
	} else if *tm != timeVal {
		t.Errorf("GetWorkerTimestamp incorrect time %v", *tm)
	}

	mock.ExpectQuery(`SELECT time FROM worker_lastlogs WHERE worker_id = \?`).WithArgs(defaultWorkerID).WillReturnRows(
		sqlmock.NewRows([]string{"time"}))
	tm, err = dbInterface.GetWorkerTimestamp(defaultWorkerID)
	if tm != nil {
		t.Errorf("GetWorkerTimestamp expected nil return %v", *tm)
	}
	if err != nil {
		t.Errorf("GetWorkerTimestamp error %v", err)
	}
}

func TestGetWorkerLogs(t *testing.T) {
	var tests = []struct {
		opts  *GetWorkerLogOpts
		query string
		args  []driver.Value
	}{
		{nil, " ORDER BY time", []driver.Value{}},
		{
			&GetWorkerLogOpts{
				Name: "foo",
			},
			` AND name = \? ORDER BY time`,
			[]driver.Value{"foo"},
		},
	}

	for i, test := range tests {
		args := append([]driver.Value{defaultWorkerID}, test.args...)
		mock.ExpectQuery(`SELECT time, name, value
                        FROM worker_metrics WHERE worker_id = \?` + test.query).WithArgs(args...).WillReturnRows(
			sqlmock.NewRows([]string{"time", "name", "value"}).AddRow("2012-01-02 12:34:56.789", "foo", "3.14159"))
		logs, err := dbInterface.GetWorkerLogs(defaultWorkerID, test.opts)
		if err != nil {
			t.Errorf("GetWorkerLogs test %d error %v", i, err)
		} else if len(logs) != 1 {
			t.Errorf("GetWorkerLogs test %d incorrect result %v", i, logs)
		}
	}
}

func TestSetSuggestionParam(t *testing.T) {
	sp := make([]*api.SuggestionParameter, 1)
	sp[0] = &api.SuggestionParameter{Name: "DefaultGrid", Value: "1"}
	studyID := generateRandid()
	mock.ExpectExec("INSERT INTO suggestion_param VALUES").WithArgs(
		sqlmock.AnyArg(), "grid", studyID,
		`{"name":"DefaultGrid","value":"1"}`).WillReturnResult(sqlmock.NewResult(1, 1))
	id, err := dbInterface.SetSuggestionParam("grid", studyID, sp)
	if err != nil {
		t.Errorf("SetSuggestionParam error %v", err)
	} else if len(id) != 16 {
		t.Errorf("SetSuggestionParam returned incorrect ID %s", id)
	}
}

func TestUpdateSuggestionParam(t *testing.T) {
	sp := make([]*api.SuggestionParameter, 1)
	sp[0] = &api.SuggestionParameter{Name: "DefaultGrid", Value: "12"}
	id := generateRandid()
	mock.ExpectExec(`UPDATE suggestion_param SET parameters = \? WHERE id = \?`).WithArgs(
		`{"name":"DefaultGrid","value":"12"}`, id).WillReturnResult(sqlmock.NewResult(1, 1))
	err := dbInterface.UpdateSuggestionParam(id, sp)
	if err != nil {
		t.Errorf("UpdateSuggestionParam error %v", err)
	}
}

func TestGetSuggestionParam(t *testing.T) {
	id := generateRandid()
	mock.ExpectQuery(`SELECT parameters FROM suggestion_param WHERE id = \?`).WithArgs(id).WillReturnRows(
		sqlmock.NewRows([]string{"parameters"}).AddRow(
			`{"name":"DefaultGrid","value":"12"}`))
	sp, err := dbInterface.GetSuggestionParam(id)
	if err != nil {
		t.Errorf("GetSuggestionParam error %v", err)
	} else if len(sp) != 1 {
		t.Errorf("GetSuggestionParam returned incorrect number of data %v", sp)
	}
}

func TestGetSuggestionParamList(t *testing.T) {
	studyID := generateRandid()
	mock.ExpectQuery(`SELECT id, suggestion_algo, parameters FROM suggestion_param WHERE study_id = \?`).WithArgs(studyID).WillReturnRows(
		sqlmock.NewRows([]string{"id", "suggestion_algo", "parameters"}).AddRow(
			generateRandid(), "random", "{}"))

	sp, err := dbInterface.GetSuggestionParamList(studyID)
	if err != nil {
		t.Errorf("GetSuggestionParamList error %v", err)
	} else if len(sp) != 1 {
		t.Errorf("GetSuggestionParamList returned incorrect number of data %v", sp)
	}
}

func TestSetEarlyStopParam(t *testing.T) {
	ep := make([]*api.EarlyStoppingParameter, 1)
	ep[0] = &api.EarlyStoppingParameter{Name: "LeastStep", Value: "1"}
	studyID := generateRandid()
	mock.ExpectExec("INSERT INTO earlystopping_param VALUES").WithArgs(
		sqlmock.AnyArg(), "medianstopping", studyID,
		`{"name":"LeastStep","value":"1"}`).WillReturnResult(sqlmock.NewResult(1, 1))
	id, err := dbInterface.SetEarlyStopParam("medianstopping", studyID, ep)
	if err != nil {
		t.Errorf("SetEarlyStopParam error %v", err)
	} else if len(id) != 16 {
		t.Errorf("SetEarlyStopParam returned incorrect ID %s", id)
	}
}

func TestUpdateEarlyStopParam(t *testing.T) {
	ep := make([]*api.EarlyStoppingParameter, 1)
	ep[0] = &api.EarlyStoppingParameter{Name: "LeastStep", Value: "12"}
	id := generateRandid()
	mock.ExpectExec(`UPDATE earlystopping_param SET parameters = \? WHERE id = \?`).WithArgs(
		`{"name":"LeastStep","value":"12"}`, id).WillReturnResult(sqlmock.NewResult(1, 1))
	err := dbInterface.UpdateEarlyStopParam(id, ep)
	if err != nil {
		t.Errorf("UpdateEarlyStopParamerror %v", err)
	}
}

func TestGetEarlyStopParam(t *testing.T) {
	id := generateRandid()
	mock.ExpectQuery(`SELECT parameters FROM earlystopping_param WHERE id = \?`).WithArgs(id).WillReturnRows(
		sqlmock.NewRows([]string{"parameters"}).AddRow(
			`{"name":"LeastStep","value":"12"}`))
	ep, err := dbInterface.GetEarlyStopParam(id)
	if err != nil {
		t.Errorf("GetEarlyStopParam error %v", err)
	} else if len(ep) != 1 {
		t.Errorf("GetEarlyStopParam returned incorrect number of data %v", ep)
	}
}

func TestGetEarlyStopParamList(t *testing.T) {
	studyID := generateRandid()
	mock.ExpectQuery(`SELECT id, earlystop_algo, parameters FROM earlystopping_param WHERE study_id = \?`).WithArgs(studyID).WillReturnRows(
		sqlmock.NewRows([]string{"id", "earlystop_algo", "parameters"}).AddRow(
			generateRandid(), "medianstopping", "{}"))

	ep, err := dbInterface.GetEarlyStopParamList(studyID)
	if err != nil {
		t.Errorf("GetEarlyStopParamList error %v", err)
	} else if len(ep) != 1 {
		t.Errorf("GetEarlyStopParamList returned incorrect number of data %v", ep)
	}
}
