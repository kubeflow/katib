// This test assumes mysql listening on localhost:3306, which can be
// prepared by the following:
// docker run -e MYSQL_ROOT_PASSWORD=test123 -e MYSQL_DATABASE=vizier -p 3306:3306 mysql

package db

import (
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"math/rand"
	"os"
	"testing"

	api "github.com/mlkube/katib/api"

	_ "github.com/go-sql-driver/mysql"
)

var db_interface VizierDBInterface

func TestMain(m *testing.M) {
	db, err := sql.Open("mysql", "root:test123@tcp(localhost:3306)/vizier")
	if err != nil {
		fmt.Printf("error opening db: %v\n", err)
		os.Exit(1)
	}
	db_interface = NewWithSqlConn(db)
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

	id, err := db_interface.CreateStudy(&in)
	if err != nil {
		t.Fatalf("CreateStudy error %v", err)
	}
	study, err := db_interface.GetStudyConfig(id)
	if err != nil {
		t.Fatalf("GetStudyConfig failed: %v", err)
	}
	fmt.Printf("%v", study)
	// TODO: check study data
}

func TestCreateStudyIdGeneration(t *testing.T) {
	var in api.StudyConfig
	in.ParameterConfigs = new(api.StudyConfig_ParameterConfigs)

	var ids []string
	for i := 0; i < 4; i++ {
		rand.Seed(1)
		id, err := db_interface.CreateStudy(&in)
		if i < 3 {
			if err != nil {
				t.Errorf("CreateStudy error %v", err)
			}
			ids = append(ids, id)
		} else if err == nil {
			t.Fatal("Expected error but succeeded")
		}
		t.Logf("id gen %d %s %v\n", i, id, err)
	}
	for _, id := range ids {
		err := db_interface.DeleteStudy(id)
		if err != nil {
			t.Errorf("DeleteStudy error %v", err)
		}
	}
}
