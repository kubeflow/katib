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

package postgres

import (
	"fmt"

	"k8s.io/klog"

	"github.com/kubeflow/katib/pkg/db/v1beta1/common"
	"github.com/kubeflow/katib/pkg/util/v1beta1/env"
)

func (d *dbConn) DBInit() {
	db := d.db
	skipDbInitialization := env.GetEnvOrDefault(common.SkipDbInitializationEnvName, "false")

	if skipDbInitialization == "false" {
		klog.Info("Initializing v1beta1 DB schema")

		_, err := db.Exec(`CREATE TABLE IF NOT EXISTS observation_logs
		(trial_name VARCHAR(255) NOT NULL,
		id serial PRIMARY KEY,
		time TIMESTAMP(6),
		metric_name VARCHAR(255) NOT NULL,
		value TEXT NOT NULL)`)
		if err != nil {
			klog.Fatalf("Error creating observation_logs table: %v", err)
		}
	} else {
		klog.Info("Skipping v1beta1 DB schema initialization.")

		_, err := db.Query(`SELECT trial_name, id, time, metric_name, value FROM observation_logs LIMIT 1`)
		if err != nil {
			klog.Fatalf("Error validating observation_logs table: %v", err)
		}
	}
}

func (d *dbConn) SelectOne() error {
	db := d.db
	_, err := db.Exec(`SELECT 1`)
	if err != nil {
		return fmt.Errorf("Error `SELECT 1` probing: %v", err)
	}
	return nil
}
