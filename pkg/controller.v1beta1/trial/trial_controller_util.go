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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	commonv1 "github.com/kubeflow/tf-operator/pkg/apis/common/v1"
)

const (
	cleanMetricsFinalizer = "clean-metrics-in-db"
)

func (r *ReconcileTrial) UpdateTrialStatusCondition(instance *trialsv1beta1.Trial, deployedJob *unstructured.Unstructured, jobCondition *commonv1.JobCondition) {
	if jobCondition == nil || instance == nil || deployedJob == nil {
		return
	}
	now := metav1.Now()
	jobConditionType := (*jobCondition).Type
	if jobConditionType == commonv1.JobSucceeded {
		if isTrialObservationAvailable(instance) {
			msg := "Trial has succeeded"
			instance.MarkTrialStatusSucceeded(corev1.ConditionTrue, TrialSucceededReason, msg)
			instance.Status.CompletionTime = &now

			eventMsg := fmt.Sprintf("Job %s has succeeded", deployedJob.GetName())
			r.recorder.Eventf(instance, corev1.EventTypeNormal, JobSucceededReason, eventMsg)
			r.collector.IncreaseTrialsSucceededCount(instance.Namespace)
		} else {
			msg := "Metrics are not available"
			instance.MarkTrialStatusSucceeded(corev1.ConditionFalse, TrialMetricsUnavailableReason, msg)

			eventMsg := fmt.Sprintf("Metrics are not available for Job %s", deployedJob.GetName())
			r.recorder.Eventf(instance, corev1.EventTypeWarning, JobMetricsUnavailableReason, eventMsg)
		}
	} else if jobConditionType == commonv1.JobFailed {
		msg := "Trial has failed"
		instance.MarkTrialStatusFailed(TrialFailedReason, msg)
		instance.Status.CompletionTime = &now

		jobConditionMessage := (*jobCondition).Message
		eventMsg := fmt.Sprintf("Job %s has failed: %s", deployedJob.GetName(), jobConditionMessage)
		r.recorder.Eventf(instance, corev1.EventTypeNormal, JobFailedReason, eventMsg)
		r.collector.IncreaseTrialsFailedCount(instance.Namespace)
	} else if jobConditionType == commonv1.JobRunning {
		msg := "Trial is running"
		instance.MarkTrialStatusRunning(TrialRunningReason, msg)
		jobConditionMessage := (*jobCondition).Message
		eventMsg := fmt.Sprintf("Job %s is running: %s",
			deployedJob.GetName(), jobConditionMessage)
		r.recorder.Eventf(instance, corev1.EventTypeNormal,
			JobRunningReason, eventMsg)
		// TODO(gaocegege): Should we maintain a TrialsRunningCount?
	}
	// else nothing to do
	return
}

func (r *ReconcileTrial) UpdateTrialStatusObservation(instance *trialsv1beta1.Trial, deployedJob *unstructured.Unstructured) error {
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
				instance.Status.Observation = &commonv1beta1.Observation{}
				metric := commonv1beta1.Metric{Name: objectiveMetricName, Value: *bestObjectiveValue}
				instance.Status.Observation.Metrics = []commonv1beta1.Metric{metric}
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

func (r *ReconcileTrial) updateFinalizers(instance *trialsv1beta1.Trial, finalizers []string) (reconcile.Result, error) {
	isDelete := true
	if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if _, err := r.DeleteTrialObservationLog(instance); err != nil {
			return reconcile.Result{}, err
		}
	} else {
		isDelete = false
	}
	instance.SetFinalizers(finalizers)
	if err := r.Update(context.TODO(), instance); err != nil {
		return reconcile.Result{}, err
	} else {
		if isDelete {
			r.collector.IncreaseTrialsDeletedCount(instance.Namespace)
		} else {
			r.collector.IncreaseTrialsCreatedCount(instance.Namespace)
		}
		// Need to requeue because finalizer update does not change metadata.generation
		return reconcile.Result{Requeue: true}, err
	}
}

func isTrialObservationAvailable(instance *trialsv1beta1.Trial) bool {
	if instance == nil {
		return false
	}
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

func isTrialComplete(instance *trialsv1beta1.Trial, jobCondition *commonv1.JobCondition) bool {
	if jobCondition == nil || instance == nil {
		return false
	}
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
	if jobCondition == nil {
		return false
	}
	jobConditionType := (*jobCondition).Type
	if jobConditionType == commonv1.JobSucceeded {
		return true
	}

	return false
}

func getBestObjectiveMetricValue(metricLogs []*api_pb.MetricLog, objectiveType commonv1beta1.ObjectiveType) *float64 {
	metricLogSize := len(metricLogs)
	if metricLogSize == 0 {
		return nil
	}

	bestObjectiveValue, _ := strconv.ParseFloat(metricLogs[0].Metric.Value, 64)
	for _, metricLog := range metricLogs[1:] {
		objectiveMetricValue, _ := strconv.ParseFloat(metricLog.Metric.Value, 64)
		if objectiveType == commonv1beta1.ObjectiveTypeMinimize {
			if objectiveMetricValue < bestObjectiveValue {
				bestObjectiveValue = objectiveMetricValue
			}
		} else if objectiveType == commonv1beta1.ObjectiveTypeMaximize {
			if objectiveMetricValue > bestObjectiveValue {
				bestObjectiveValue = objectiveMetricValue
			}
		}

	}
	return &bestObjectiveValue
}

func needUpdateFinalizers(trial *trialsv1beta1.Trial) (bool, []string) {
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
