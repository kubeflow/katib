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

	commonv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/common/v1alpha3"
	trialsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1alpha3"
	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
	commonv1 "github.com/kubeflow/tf-operator/pkg/apis/common/v1"
)

func (r *ReconcileTrial) GetDeployedJobStatus(deployedJob *unstructured.Unstructured) (*commonv1.JobConditionType, error) {
	jobConditionType := commonv1.JobRunning
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
				return nil, err
			}
			if jobStatus.Active == 0 && jobStatus.Succeeded > 0 {
				jobConditionType = commonv1.JobSucceeded
			} else if jobStatus.Failed > 0 {
				jobConditionType = commonv1.JobFailed
			}
		default:
			jobStatus := commonv1.JobStatus{}
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(statusMap, &jobStatus)

			if err != nil {
				log.Error(err, "Convert unstructured to status error")
				return nil, err
			}
			if len(jobStatus.Conditions) > 0 {
				lc := jobStatus.Conditions[len(jobStatus.Conditions)-1]
				jobConditionType = lc.Type
			}
		}
	} else if unerr != nil {
		log.Error(unerr, "NestedFieldCopy unstructured to status error")
		return nil, unerr
	}
	return &jobConditionType, nil
}

func (r *ReconcileTrial) UpdateTrialStatusCondition(instance *trialsv1alpha3.Trial, jobCondition commonv1.JobConditionType) {
	now := metav1.Now()
	if jobCondition == commonv1.JobSucceeded {
		msg := "Trial has succeeded"
		instance.MarkTrialStatusSucceeded(TrialSucceededReason, msg)
		instance.Status.CompletionTime = &now
	} else if jobCondition == commonv1.JobFailed {
		msg := "Trial has failed"
		instance.MarkTrialStatusFailed(TrialFailedReason, msg)
		instance.Status.CompletionTime = &now
	}
	//else nothing to do
	return
}

func (r *ReconcileTrial) UpdateTrialStatusObservation(instance *trialsv1alpha3.Trial, deployedJob *unstructured.Unstructured) error {
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
				instance.Status.Observation = &commonv1alpha3.Observation{}
				metric := commonv1alpha3.Metric{Name: objectiveMetricName, Value: *bestObjectiveValue}
				instance.Status.Observation.Metrics = []commonv1alpha3.Metric{metric}
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

func isTrialObservationAvailable(instance *trialsv1alpha3.Trial) bool {
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

func isTrialComplete(instance *trialsv1alpha3.Trial, jobConditionType commonv1.JobConditionType) bool {
	if jobConditionType == commonv1.JobSucceeded && isTrialObservationAvailable(instance) {
		return true
	}
	if jobConditionType == commonv1.JobFailed {
		return true
	}

	return false
}

func getBestObjectiveMetricValue(metricLogs []*api_pb.MetricLog, objectiveType commonv1alpha3.ObjectiveType) *float64 {
	metricLogSize := len(metricLogs)
	if metricLogSize == 0 {
		return nil
	}

	bestObjectiveValue, _ := strconv.ParseFloat(metricLogs[0].Metric.Value, 32)
	for _, metricLog := range metricLogs[1:] {
		objectiveMetricValue, _ := strconv.ParseFloat(metricLog.Metric.Value, 32)
		if objectiveType == commonv1alpha3.ObjectiveTypeMinimize {
			if objectiveMetricValue < bestObjectiveValue {
				bestObjectiveValue = objectiveMetricValue
			}
		} else if objectiveType == commonv1alpha3.ObjectiveTypeMaximize {
			if objectiveMetricValue > bestObjectiveValue {
				bestObjectiveValue = objectiveMetricValue
			}
		}

	}
	return &bestObjectiveValue
}
