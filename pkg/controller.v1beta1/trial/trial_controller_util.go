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
	"fmt"
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	trialutil "github.com/kubeflow/katib/pkg/controller.v1beta1/trial/util"
)

const (
	cleanMetricsFinalizer = "clean-metrics-in-db"
)

// UpdateTrialStatusCondition updates Trial status from current deployed Job status
func (r *ReconcileTrial) UpdateTrialStatusCondition(instance *trialsv1beta1.Trial, deployedJobName string, jobStatus *trialutil.TrialJobStatus) {
	logger := log.WithValues("Trial", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})

	timeNow := metav1.Now()

	if jobStatus.Condition == trialutil.JobSucceeded {
		if instance.IsObservationAvailable() && !instance.IsSucceeded() {
			if !instance.IsEarlyStopped() {
				msg := "Trial has succeeded"
				reason := TrialSucceededReason

				// Get message and reason from deployed job
				if jobStatus.Message != "" {
					msg = fmt.Sprintf("%v. Job message: %v", msg, jobStatus.Message)
				}
				if jobStatus.Reason != "" {
					reason = fmt.Sprintf("%v. Job reason: %v", reason, jobStatus.Reason)
				}

				logger.Info("Trial status changed to Succeeded")
				instance.MarkTrialStatusSucceeded(corev1.ConditionTrue, reason, msg)
				instance.Status.CompletionTime = &timeNow

				eventMsg := fmt.Sprintf("Job %v has succeeded", deployedJobName)
				r.recorder.Eventf(instance, corev1.EventTypeNormal, JobSucceededReason, eventMsg)
				r.collector.IncreaseTrialsSucceededCount(instance.Namespace)
			}
		} else if !instance.IsMetricsUnavailable() {
			msg := "Metrics are not available"
			reason := TrialMetricsUnavailableReason

			// Get message and reason from deployed job
			if jobStatus.Message != "" {
				msg = fmt.Sprintf("%v. Job message: %v", msg, jobStatus.Message)
			}
			if jobStatus.Reason != "" {
				reason = fmt.Sprintf("%v. Job reason: %v", reason, jobStatus.Reason)
			}

			logger.Info("Trial status changed to Metrics Unavailable")
			instance.MarkTrialStatusMetricsUnavailable(reason, msg)
			instance.Status.CompletionTime = &timeNow

			eventMsg := fmt.Sprintf("Metrics are not available for Job %v", deployedJobName)
			r.recorder.Eventf(instance, corev1.EventTypeWarning, JobMetricsUnavailableReason, eventMsg)
			r.collector.IncreaseTrialsMetricsUnavailableCount(instance.Namespace)
		}
	} else if jobStatus.Condition == trialutil.JobFailed && !instance.IsFailed() && !instance.IsEarlyStopped() {
		msg := "Trial has failed"
		reason := TrialFailedReason

		// Get message and reason from deployed job
		if jobStatus.Message != "" {
			msg = fmt.Sprintf("%v. Job message: %v", msg, jobStatus.Message)
		}
		if jobStatus.Reason != "" {
			reason = fmt.Sprintf("%v. Job reason: %v", reason, jobStatus.Reason)
		}

		instance.MarkTrialStatusFailed(reason, msg)
		instance.Status.CompletionTime = &timeNow

		eventMsg := fmt.Sprintf("Job %v has failed", deployedJobName)
		if jobStatus.Message != "" || jobStatus.Reason != "" {
			eventMsg = fmt.Sprintf("%v. %v %v", eventMsg, jobStatus.Message, jobStatus.Reason)
		}

		r.recorder.Eventf(instance, corev1.EventTypeNormal, JobFailedReason, eventMsg)
		r.collector.IncreaseTrialsFailedCount(instance.Namespace)
		logger.Info("Trial status changed to Failed")
	} else if jobStatus.Condition == trialutil.JobRunning && !instance.IsRunning() && !instance.IsEarlyStopped() {
		msg := "Trial is running"
		instance.MarkTrialStatusRunning(TrialRunningReason, msg)

		eventMsg := fmt.Sprintf("Job %v is running", deployedJobName)
		r.recorder.Eventf(instance, corev1.EventTypeNormal, JobRunningReason, eventMsg)
		logger.Info("Trial status changed to Running")
		// TODO(gaocegege): Should we maintain a TrialsRunningCount?
	}
	// else nothing to do
}

func (r *ReconcileTrial) UpdateTrialStatusObservation(instance *trialsv1beta1.Trial) error {
	reply, err := r.GetTrialObservationLog(instance)
	if err != nil {
		log.Error(err, "Get trial observation log error")
		return err
	}
	metricStrategies := instance.Spec.Objective.MetricStrategies
	if len(reply.ObservationLog.MetricLogs) != 0 {
		observation, err := getMetrics(reply.ObservationLog.MetricLogs, metricStrategies)
		if err != nil {
			log.Error(err, "Get metrics from logs error")
			return err
		}
		instance.Status.Observation = observation
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

func getMetrics(metricLogs []*api_pb.MetricLog, strategies []commonv1beta1.MetricStrategy) (*commonv1beta1.Observation, error) {
	metrics := make(map[string]*commonv1beta1.Metric)
	timestamps := make(map[string]*time.Time)
	for _, strategy := range strategies {
		timestamps[strategy.Name] = nil
		metrics[strategy.Name] = &commonv1beta1.Metric{
			Name:   strategy.Name,
			Min:    consts.UnavailableMetricValue,
			Max:    consts.UnavailableMetricValue,
			Latest: consts.UnavailableMetricValue,
		}
	}

	for _, metricLog := range metricLogs {
		metric, ok := metrics[metricLog.Metric.Name]
		if !ok {
			continue
		}
		strValue := metricLog.Metric.Value
		floatValue, err := strconv.ParseFloat(strValue, 64)
		if err == nil {
			if metric.Min == consts.UnavailableMetricValue {
				metric.Min = strValue
				metric.Max = strValue
			} else {
				// We can't get error here, because we parsed this value before
				minMetric, _ := strconv.ParseFloat(metric.Min, 64)
				maxMetric, _ := strconv.ParseFloat(metric.Max, 64)
				if floatValue < minMetric {
					metric.Min = strValue
				} else if floatValue > maxMetric {
					metric.Max = strValue
				}
			}
		}
		currentTime, err := time.Parse(time.RFC3339Nano, metricLog.TimeStamp)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamps %s: %e", metricLog.TimeStamp, err)
		}
		timestamp := timestamps[metricLog.Metric.Name]
		if timestamp == nil || !timestamp.After(currentTime) {
			timestamps[metricLog.Metric.Name] = &currentTime
			metric.Latest = strValue
		}
	}

	observation := &commonv1beta1.Observation{}
	for _, metric := range metrics {
		observation.Metrics = append(observation.Metrics, *metric)
	}

	return observation, nil
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
