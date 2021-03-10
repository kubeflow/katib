package db

import (
	"github.com/kubeflow/katib/pkg/db/v1beta1/common"
	"github.com/kubeflow/katib/pkg/db/v1beta1/mysql"
)

func NewKatibDBInterface() (common.KatibDBInterface, error) {
	return mysql.NewDBInterface()
}
