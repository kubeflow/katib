package db

import (
	"errors"

	"github.com/kubeflow/katib/pkg/db/v1beta1/common"
	"github.com/kubeflow/katib/pkg/db/v1beta1/mysql"
)

func NewKatibDBInterface(dbName string) (common.KatibDBInterface, error) {

	if dbName == common.MySqlDBNameEnvValue {
		return mysql.NewDBInterface()
	}
	return nil, errors.New("Invalid DB Name")
}
