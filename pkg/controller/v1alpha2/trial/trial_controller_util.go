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

package trial

import (
	"strconv"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	commonv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/common/v1alpha2"
	trialsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
	api_pb "github.com/kubeflow/katib/pkg/api/v1alpha2"
	commonv1 "github.com/kubeflow/tf-operator/pkg/apis/common/v1"
)

func (r *ReconcileTrial) UpdateTrialStatusCondition(instance *trialsv1alpha2.Trial, deployedJob *unstructured.Unstructured) error {

	kind := deployedJob.GetKind()
	status, ok, unerr := unstructured.NestedFieldCopy(deployedJob.Object, "status")
	now := metav1.Now()

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
				instance.Status.CompletionTime = &now
			} else if jobStatus.Failed > 0 {
				msg := "Trial has failed"
				instance.MarkTrialStatusFailed(TrialFailedReason, msg)
				instance.Status.CompletionTime = &now
			}
		default:
			jobStatus := commonv1.JobStatus{}
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(statusMap, &jobStatus)

			if err != nil {
				log.Error(err, "Convert unstructured to status error")
				return err
			}
			if len(jobStatus.Conditions) > 0 {
				lc := jobStatus.Conditions[len(jobStatus.Conditions)-1]
				if lc.Type == commonv1.JobSucceeded {
					msg := "Trial has succeeded"
					instance.MarkTrialStatusSucceeded(TrialSucceededReason, msg)
					instance.Status.CompletionTime = &now
				} else if lc.Type == commonv1.JobFailed {
					msg := "Trial has failed"
					instance.MarkTrialStatusFailed(TrialFailedReason, msg)
					instance.Status.CompletionTime = &now
				}
			}
		}
	} else if unerr != nil {
		log.Error(unerr, "NestedFieldCopy unstructured to status error")
		return unerr
	}
	return nil
}

func (r *ReconcileTrial) UpdateTrialStatusObservation(instance *trialsv1alpha2.Trial, deployedJob *unstructured.Unstructured) error {
	objectiveMetricName := instance.Spec.Objective.ObjectiveMetricName
	reply, err := r.GetTrialObservationLog(instance)
	if err != nil {
		log.Error(err, "Get trial observation log error")
		return err
	}
	if reply.ObservationLog != nil {
		bestObjectiveValue := getBestObjectiveMetricValue(reply.ObservationLog.MetricLogs, instance.Spec.Objective.Type)
		if bestObjectiveValue != nil {
			if instance.Status.Observation == nil {
				instance.Status.Observation = &commonv1alpha2.Observation{}
				metric := commonv1alpha2.Metric{Name: objectiveMetricName, Value: *bestObjectiveValue}
				instance.Status.Observation.Metrics = []commonv1alpha2.Metric{metric}
			} else {
				for index, metric := range instance.Status.Observation.Metrics {
					if metric.Name == objectiveMetricName {
						instance.Status.Observation.Metrics[index].Value = *bestObjectiveValue
					}
				}
			}
		}
	}
	return nil
}

func isTrialObservationAvailable(instance *trialsv1alpha2.Trial) bool {
	objectiveMetricName := instance.Spec.Objective.ObjectiveMetricName
	if instance.Status.Observation != nil && instance.Status.Observation.Metrics != nil {
		for _, metric := range instance.Status.Observation.Metrics {
			if metric.Name == objectiveMetricName {
				return true
			}
		}
	}
	return false
}

func getBestObjectiveMetricValue(metricLogs []*api_pb.MetricLog, objectiveType commonv1alpha2.ObjectiveType) *float64 {
	metricLogSize := len(metricLogs)
	if metricLogSize == 0 {
		return nil
	}

	bestObjectiveValue, _ := strconv.ParseFloat(metricLogs[0].Metric.Value, 32)
	for _, metricLog := range metricLogs[1:] {
		objectiveMetricValue, _ := strconv.ParseFloat(metricLog.Metric.Value, 32)
		if objectiveType == commonv1alpha2.ObjectiveTypeMinimize {
			if objectiveMetricValue < bestObjectiveValue {
				bestObjectiveValue = objectiveMetricValue
			}
		} else if objectiveType == commonv1alpha2.ObjectiveTypeMaximize {
			if objectiveMetricValue > bestObjectiveValue {
				bestObjectiveValue = objectiveMetricValue
			}
		}

	}
	return &bestObjectiveValue
}
