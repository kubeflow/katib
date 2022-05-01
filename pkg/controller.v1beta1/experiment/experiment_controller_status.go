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

package experiment

import (
	"context"

	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
)

type updateStatusFunc func(instance *experimentsv1beta1.Experiment) error

func (r *ReconcileExperiment) updateStatus(instance *experimentsv1beta1.Experiment) error {
	err := r.Status().Update(context.TODO(), instance)
	if err != nil {
		return err
	}
	return nil
}
