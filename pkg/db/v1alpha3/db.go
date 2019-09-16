package db

import (
	"errors"

	"github.com/kubeflow/katib/pkg/db/v1alpha3/common"
	"github.com/kubeflow/katib/pkg/db/v1alpha3/mysql"
)

func NewKatibDBInterface(dbName string) (common.KatibDBInterface, error) {

	if dbName == "mysql" {
		return mysql.NewDBInterface()
	}
	return nil, errors.New("Invalid DB Name")
}
