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
	"context"
	"fmt"
	"strconv"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	commonv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/common/v1alpha3"
	trialsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1alpha3"
	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
	commonv1 "github.com/kubeflow/tf-operator/pkg/apis/common/v1"
)

const (
	cleanMetricsFinalizer = "clean-metrics-in-db"
)

func (r *ReconcileTrial) GetDeployedJobStatus(deployedJob *unstructured.Unstructured) (*commonv1.JobCondition, error) {
	jobCondition := commonv1.JobCondition{}
	jobCondition.Type = commonv1.JobRunning
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
			for _, cond := range jobStatus.Conditions {
				if cond.Type == batchv1.JobComplete && cond.Status == corev1.ConditionTrue {
					jobCondition.Type = commonv1.JobSucceeded
					//  JobConditions message not populated when succeeded for batchv1 Job
					break
				}
				if cond.Type == batchv1.JobFailed && cond.Status == corev1.ConditionTrue {
					jobCondition.Type = commonv1.JobFailed
					jobCondition.Message = cond.Message
					break
				}
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
				jobCondition.Type = lc.Type
				jobCondition.Message = lc.Message
			}
		}
	} else if unerr != nil {
		log.Error(unerr, "NestedFieldCopy unstructured to status error")
		return nil, unerr
	}
	return &jobCondition, nil
}

func (r *ReconcileTrial) UpdateTrialStatusCondition(instance *trialsv1alpha3.Trial, deployedJob *unstructured.Unstructured, jobCondition *commonv1.JobCondition) {
	now := metav1.Now()
	jobConditionType := (*jobCondition).Type
	if jobConditionType == commonv1.JobSucceeded {
		if isTrialObservationAvailable(instance) {
			msg := "Trial has succeeded"
			instance.MarkTrialStatusSucceeded(TrialSucceededReason, msg)
			instance.Status.CompletionTime = &now

			eventMsg := fmt.Sprintf("Job %s has succeeded", deployedJob.GetName())
			r.recorder.Eventf(instance, corev1.EventTypeNormal, JobSucceededReason, eventMsg)
		} else {
			eventMsg := fmt.Sprintf("Metrics are not available for Job %s", deployedJob.GetName())
			r.recorder.Eventf(instance, corev1.EventTypeWarning, MetricsUnavailableReason, eventMsg)
		}
	} else if jobConditionType == commonv1.JobFailed {
		msg := "Trial has failed"
		instance.MarkTrialStatusFailed(TrialFailedReason, msg)
		instance.Status.CompletionTime = &now

		jobConditionMessage := (*jobCondition).Message
		eventMsg := fmt.Sprintf("Job %s has failed: %s", deployedJob.GetName(), jobConditionMessage)
		r.recorder.Eventf(instance, corev1.EventTypeNormal, JobFailedReason, eventMsg)
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

func (r *ReconcileTrial) updateFinalizers(instance *trialsv1alpha3.Trial, finalizers []string) (reconcile.Result, error) {
	if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if _, err := r.DeleteTrialObservationLog(instance); err != nil {
			return reconcile.Result{}, err
		}
	}
	instance.SetFinalizers(finalizers)
	if err := r.Update(context.TODO(), instance); err != nil {
		return reconcile.Result{}, err
	} else {
		// Need to requeue because finalizer update does not change metadata.generation
		return reconcile.Result{Requeue: true}, err
	}
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

func isTrialComplete(instance *trialsv1alpha3.Trial, jobCondition *commonv1.JobCondition) bool {
	jobConditionType := (*jobCondition).Type
	if jobConditionType == commonv1.JobSucceeded && isTrialObservationAvailable(instance) {
		return true
	}
	if jobConditionType == commonv1.JobFailed {
		return true
	}

	return false
}

func isJobSucceeded(jobCondition *commonv1.JobCondition) bool {
	jobConditionType := (*jobCondition).Type
	if jobConditionType == commonv1.JobSucceeded {
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

func needUpdateFinalizers(trial *trialsv1alpha3.Trial) (bool, []string) {
	deleted := !trial.ObjectMeta.DeletionTimestamp.IsZero()
	pendingFinalizers := trial.GetFinalizers()
	contained := false
	for _, elem := range pendingFinalizers {
		if elem == cleanMetricsFinalizer {
			contained = true
			break
		}
	}

	if !deleted && !contained {
		finalizers := append(pendingFinalizers, cleanMetricsFinalizer)
		return true, finalizers
	}
	if deleted && contained {
		finalizers := []string{}
		for _, pendingFinalizer := range pendingFinalizers {
			if pendingFinalizer != cleanMetricsFinalizer {
				finalizers = append(finalizers, pendingFinalizer)
			}
		}
		return true, finalizers
	}
	return false, []string{}
}
