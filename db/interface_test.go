// This test assumes mysql listening on localhost:3306, which can be
// prepared by the following:
// docker run -e MYSQL_ROOT_PASSWORD=test123 -e MYSQL_DATABASE=vizier -p 3306:3306 mysql

package db

import (
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"math/rand"
	"os"
	"testing"

	api "github.com/kubeflow/hp-tuning/api"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var db_interface VizierDBInterface
var mock sqlmock.Sqlmock

func TestMain(m *testing.M) {
	//	db, err := sql.Open("mysql", "root:test123@tcp(localhost:3306)/vizier")
	db, sm, err := sqlmock.New()
	mock = sm
	if err != nil {
		fmt.Printf("error opening db: %v\n", err)
		os.Exit(1)
	}
	//mock.ExpectBegin()
	db_interface = NewWithSqlConn(db)
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS studies").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS study_permissions").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS trials").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS trial_logs").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS workers").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	db_interface.DB_Init()

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
	id, err := db_interface.CreateStudy(&in)
	if err != nil {
		t.Errorf("CreateStudy error %v", err)
	}
	//	mock.ExpectExec("SELECT * FROM studies WHERE id").WithArgs(id).WillReturnRows(sqlmock.NewRows())
	mock.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"id",
			"name",
			"owner",
			"optimization_type",
			"optimization_goal",
			"parameter_configs",
			"suggest_algo",
			"early_stop_algo",
			"study_task_name",
			"suggestion_parameters",
			"early_stopping_parameters",
			"tags",
			"objective_value_name",
			"metrics",
			"image",
			"command",
			"gpu",
			"scheduler",
			"mount",
			"pull_secret",
		}).
			AddRow("abc", "test", "admin", 1, 0.99, "{}", "random", "test", "", "", "", "", "", "", "", "", 1, "", "", ""))
	study, err := db_interface.GetStudyConfig(id)
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
		id, err := db_interface.CreateStudy(&in)
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
		err := db_interface.DeleteStudy(id)
		if err != nil {
			t.Errorf("DeleteStudy error %v", err)
		}
	}
}
