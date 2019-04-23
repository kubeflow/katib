/*

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

package util

import (
	//v1 "k8s.io/api/core/v1"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	trialsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
	commonv1beta1 "github.com/kubeflow/tf-operator/pkg/apis/common/v1beta1"
)

var log = logf.Log.WithName("trial-status-util")

const (
	DefaultJobKind       = "Job"
	TrialCreatedReason   = "TrialCreated"
	TrialRunningReason   = "TrialRunning"
	TrialSucceededReason = "TrialSucceeded"
	TrialFailedReason    = "TrialFailed"
	TrialKilledReason    = "TrialKilled"
)

func UpdateTrialStatusCondition(instance *trialsv1alpha2.Trial, deployedJob *unstructured.Unstructured) error {

	kind := deployedJob.GetKind()
	status, ok, unerr := unstructured.NestedFieldCopy(deployedJob.Object, "status")

	if ok {
		statusMap := status.(map[string]interface{})
		switch kind {

		case DefaultJobKind:
			jobStatus := batchv1.JobStatus{}
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(statusMap, &jobStatus)
			if err != nil {
				log.Error(err, "Convert unstructured to status error")
				return err
			}
			if jobStatus.Active == 0 && jobStatus.Succeeded > 0 {
				msg := "Trial has succeeded"
				instance.MarkTrialStatusSucceeded(TrialSucceededReason, msg)
			} else if jobStatus.Failed > 0 {
				msg := "Trial has failed"
				instance.MarkTrialStatusFailed(TrialFailedReason, msg)
			}
		default:
			jobStatus := commonv1beta1.JobStatus{}
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(statusMap, &jobStatus)

			if err != nil {
				log.Error(err, "Convert unstructured to status error")
				return err
			}
			if len(jobStatus.Conditions) > 0 {
				lc := jobStatus.Conditions[len(jobStatus.Conditions)-1]
				if lc.Type == commonv1beta1.JobSucceeded {
					msg := "Trial has succeeded"
					instance.MarkTrialStatusSucceeded(TrialSucceededReason, msg)
				} else if lc.Type == commonv1beta1.JobFailed {
					msg := "Trial has failed"
					instance.MarkTrialStatusFailed(TrialFailedReason, msg)
				}
			}
		}
	} else if unerr != nil {
		log.Error(unerr, "NestedFieldCopy unstructured to status error")
		return unerr
	}
	return nil
}

func UpdateTrialStatusObservation(instance *trialsv1alpha2.Trial, deployedJob *unstructured.Unstructured) error {

	// read GetObservationLog call and update observation field
	return nil
}
