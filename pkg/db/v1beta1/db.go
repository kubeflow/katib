/*
Copyright 2022 The Kubeflow Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package db

import (
	"errors"
	"time"

	"github.com/kubeflow/katib/pkg/db/v1beta1/common"
	"github.com/kubeflow/katib/pkg/db/v1beta1/mysql"
	"github.com/kubeflow/katib/pkg/db/v1beta1/postgres"
	"k8s.io/klog"
)

func NewKatibDBInterface(dbName string, connectTimeout time.Duration) (common.KatibDBInterface, error) {

	if dbName == common.MySqlDBNameEnvValue {
		klog.Info("Using MySQL")
		return mysql.NewDBInterface(connectTimeout)
	} else if dbName == common.PostgresSQLDBNameEnvValue {
		klog.Info("Using Postgres")
		return postgres.NewDBInterface(connectTimeout)
	}
	return nil, errors.New("Invalid DB Name")
}
