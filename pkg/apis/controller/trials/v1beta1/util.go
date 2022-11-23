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

package v1beta1

import (
	"errors"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getCondition(trial *Trial, condType TrialConditionType) *TrialCondition {
	for _, condition := range trial.Status.Conditions {
		if condition.Type == condType {
			return &condition
		}
	}
	return nil
}

func hasCondition(trial *Trial, condType TrialConditionType) bool {
	cond := getCondition(trial, condType)
	if cond != nil && cond.Status == v1.ConditionTrue {
		return true
	}
	return false
}

func (trial *Trial) removeCondition(condType TrialConditionType) {
	var newConditions []TrialCondition
	for _, c := range trial.Status.Conditions {

		if c.Type == condType {
			continue
		}

		newConditions = append(newConditions, c)
	}
	trial.Status.Conditions = newConditions
}

func newCondition(conditionType TrialConditionType, status v1.ConditionStatus, reason, message string) TrialCondition {
	return TrialCondition{
		Type:               conditionType,
		Status:             status,
		LastUpdateTime:     metav1.Now(),
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            message,
	}
}

func (trial *Trial) IsCreated() bool {
	return hasCondition(trial, TrialCreated)
}

func (trial *Trial) IsRunning() bool {
	return hasCondition(trial, TrialRunning)
}

func (trial *Trial) IsSucceeded() bool {
	return hasCondition(trial, TrialSucceeded)
}

func (trial *Trial) IsFailed() bool {
	return hasCondition(trial, TrialFailed)
}

func (trial *Trial) IsKilled() bool {
	return hasCondition(trial, TrialKilled)
}

// IsMetricsUnavailable returns true if Trial metrics are not available
func (trial *Trial) IsMetricsUnavailable() bool {
	return hasCondition(trial, TrialMetricsUnavailable)
}

// IsObservationAvailable return ture if the Trial has valid observations updated
func (trial *Trial) IsObservationAvailable() bool {
	if trial.Spec.Objective == nil {
		return false
	}
	objectiveMetricName := trial.Spec.Objective.ObjectiveMetricName
	if trial.Status.Observation != nil && trial.Status.Observation.Metrics != nil {
		for _, metric := range trial.Status.Observation.Metrics {
			if metric.Name == objectiveMetricName && metric.Latest != consts.UnavailableMetricValue {
				return true
			}
		}
	}
	return false
}

func (trial *Trial) IsCompleted() bool {
	return trial.IsSucceeded() || trial.IsFailed() || trial.IsKilled() || trial.IsEarlyStopped() || trial.IsMetricsUnavailable()
}

func (trial *Trial) IsEarlyStopped() bool {
	return hasCondition(trial, TrialEarlyStopped)
}

func (trial *Trial) GetLastConditionType() (TrialConditionType, error) {
	if len(trial.Status.Conditions) > 0 {
		return trial.Status.Conditions[len(trial.Status.Conditions)-1].Type, nil
	}
	return "", errors.New("Trial doesn't have any condition")
}

func (trial *Trial) setCondition(conditionType TrialConditionType, status v1.ConditionStatus, reason, message string) {

	newCond := newCondition(conditionType, status, reason, message)
	currentCond := getCondition(trial, conditionType)
	// Do nothing if condition doesn't change
	if currentCond != nil && currentCond.Status == newCond.Status && currentCond.Reason == newCond.Reason {
		return
	}

	// Do not update lastTransitionTime if the status of the condition doesn't change.
	if currentCond != nil && currentCond.Status == newCond.Status {
		newCond.LastTransitionTime = currentCond.LastTransitionTime
	}

	trial.removeCondition(conditionType)
	trial.Status.Conditions = append(trial.Status.Conditions, newCond)
}

func (trial *Trial) MarkTrialStatusCreated(reason, message string) {
	trial.setCondition(TrialCreated, v1.ConditionTrue, reason, message)
}

func (trial *Trial) MarkTrialStatusRunning(reason, message string) {
	trial.setCondition(TrialRunning, v1.ConditionTrue, reason, message)
}

func (trial *Trial) MarkTrialStatusSucceeded(status v1.ConditionStatus, reason, message string) {
	currentCond := getCondition(trial, TrialRunning)
	if currentCond != nil {
		trial.setCondition(TrialRunning, v1.ConditionFalse, currentCond.Reason, currentCond.Message)
	}
	trial.setCondition(TrialSucceeded, status, reason, message)

}

func (trial *Trial) MarkTrialStatusFailed(reason, message string) {
	currentCond := getCondition(trial, TrialRunning)
	if currentCond != nil {
		trial.setCondition(TrialRunning, v1.ConditionFalse, currentCond.Reason, currentCond.Message)
	}
	trial.setCondition(TrialFailed, v1.ConditionTrue, reason, message)
}

func (trial *Trial) MarkTrialStatusKilled(reason, message string) {
	currentCond := getCondition(trial, TrialRunning)
	if currentCond != nil {
		trial.setCondition(TrialRunning, v1.ConditionFalse, currentCond.Reason, currentCond.Message)
	}
	trial.setCondition(TrialKilled, v1.ConditionTrue, reason, message)
}

func (trial *Trial) MarkTrialStatusMetricsUnavailable(reason, message string) {
	currentCond := getCondition(trial, TrialRunning)
	if currentCond != nil {
		trial.setCondition(TrialRunning, v1.ConditionFalse, currentCond.Reason, currentCond.Message)
	}
	trial.setCondition(TrialMetricsUnavailable, v1.ConditionTrue, reason, message)
}
