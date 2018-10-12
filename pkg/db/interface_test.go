package db

import (
	"database/sql/driver"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	api "github.com/kubeflow/katib/pkg/api"
)

var dbInterface VizierDBInterface
var mock sqlmock.Sqlmock

func TestMain(m *testing.M) {
	db, sm, err := sqlmock.New()
	mock = sm
	if err != nil {
		fmt.Printf("error opening db: %v\n", err)
		os.Exit(1)
	}
	//mock.ExpectBegin()
	dbInterface = NewWithSQLConn(db)
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS studies").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS study_permissions").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS trials").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS workers").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS worker_metrics").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS worker_lastlogs").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS suggestion_param").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS earlystop_param").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	dbInterface.DBInit()

	os.Exit(m.Run())
}

func TestGetStudyConfig(t *testing.T) {
	var in api.StudyConfig
	in.ParameterConfigs = new(api.StudyConfig_ParameterConfigs)
	//err := jsonpb.UnmarshalString("{}", &in)
	err := jsonpb.UnmarshalString(`{"configs": [{"name": "-abc"}]}`, in.ParameterConfigs)
	if err != nil {
		t.Errorf("err %v", err)
	}

	mock.ExpectExec("INSERT INTO studies VALUES").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	id, err := dbInterface.CreateStudy(&in)
	if err != nil {
		t.Errorf("CreateStudy error %v", err)
	}
	//	mock.ExpectExec("SELECT * FROM studies WHERE id").WithArgs(id).WillReturnRows(sqlmock.NewRows())
	mock.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{
			"id",
			"name",
			"owner",
			"optimization_type",
			"optimization_goal",
			"parameter_configs",
			"tags",
			"objective_value_name",
			"metrics",
			"job_id",
		}).
			AddRow("abc", "test", "admin", 1, 0.99, "{}", "", "", "", "test"))
	study, err := dbInterface.GetStudyConfig(id)
	if err != nil {
		t.Errorf("GetStudyConfig failed: %v", err)
	}
	fmt.Printf("%v", study)
	// TODO: check study data
}

func TestCreateStudyIdGeneration(t *testing.T) {
	var in api.StudyConfig
	in.ParameterConfigs = new(api.StudyConfig_ParameterConfigs)

	var ids []string
	for i := 0; i < 4; i++ {
		rand.Seed(int64(i))
		mock.ExpectExec("INSERT INTO studies VALUES").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
		id, err := dbInterface.CreateStudy(&in)
		if err != nil {
			t.Errorf("CreateStudy error %v", err)
		}
		ids = append(ids, id)
		t.Logf("id gen %d %s %v\n", i, id, err)
	}
	encountered := map[string]bool{}
	for i := 0; i < len(ids); i++ {
		if !encountered[ids[i]] {
			encountered[ids[i]] = true
		} else {
			t.Fatalf("Study ID duplicated %v", ids)
		}
	}
	for _, id := range ids {
		mock.ExpectExec("DELETE").WillReturnResult(sqlmock.NewResult(1, 1))
		err := dbInterface.DeleteStudy(id)
		if err != nil {
			t.Errorf("DeleteStudy error %v", err)
		}
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

var workerColumns = []string{"id",
	"study_id", "trial_id", "type",
	"status", "template_path", "tags"}

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
