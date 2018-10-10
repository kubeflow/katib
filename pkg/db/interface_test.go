package db

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

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
