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

package trial

import (
	"context"

	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
)

const (
	// For Trials
	TrialCreatedReason            = "TrialCreated"
	TrialRunningReason            = "TrialRunning"
	TrialSucceededReason          = "TrialSucceeded"
	TrialMetricsUnavailableReason = "MetricsUnavailable"
	TrialFailedReason             = "TrialFailed"

	// For Jobs
	JobCreatedReason            = "JobCreated"
	JobDeletedReason            = "JobDeleted"
	JobSucceededReason          = "JobSucceeded"
	JobMetricsUnavailableReason = "MetricsUnavailable"
	JobFailedReason             = "JobFailed"
	JobRunningReason            = "JobRunning"
)

type updateStatusFunc func(instance *trialsv1beta1.Trial) error

func (r *ReconcileTrial) updateStatus(instance *trialsv1beta1.Trial) error {
	err := r.Status().Update(context.TODO(), instance)
	if err != nil {
		return err
	}
	return nil
}
