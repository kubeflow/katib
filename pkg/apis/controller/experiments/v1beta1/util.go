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

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getCondition(exp *Experiment, condType ExperimentConditionType) *ExperimentCondition {
	if exp.Status.Conditions != nil {
		for _, condition := range exp.Status.Conditions {
			if condition.Type == condType {
				return &condition
			}
		}
	}
	return nil
}

func hasCondition(exp *Experiment, condType ExperimentConditionType) bool {
	cond := getCondition(exp, condType)
	if cond != nil && cond.Status == v1.ConditionTrue {
		return true
	}
	return false
}

func (exp *Experiment) removeCondition(condType ExperimentConditionType) {
	var newConditions []ExperimentCondition
	for _, c := range exp.Status.Conditions {

		if c.Type == condType {
			continue
		}

		newConditions = append(newConditions, c)
	}
	exp.Status.Conditions = newConditions
}

func newCondition(conditionType ExperimentConditionType, status v1.ConditionStatus, reason, message string) ExperimentCondition {
	return ExperimentCondition{
		Type:               conditionType,
		Status:             status,
		LastUpdateTime:     metav1.Now(),
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            message,
	}
}

func (exp *Experiment) IsCreated() bool {
	return hasCondition(exp, ExperimentCreated)
}

func (exp *Experiment) IsSucceeded() bool {
	return hasCondition(exp, ExperimentSucceeded)
}

func (exp *Experiment) IsFailed() bool {
	return hasCondition(exp, ExperimentFailed)
}

func (exp *Experiment) IsRunning() bool {
	return hasCondition(exp, ExperimentRunning)
}

func (exp *Experiment) IsRestarting() bool {
	return hasCondition(exp, ExperimentRestarting)
}

func (exp *Experiment) IsCompleted() bool {
	return exp.IsSucceeded() || exp.IsFailed()
}

func (exp *Experiment) IsCompletedReason(reason string) bool {
	cond := getCondition(exp, ExperimentSucceeded)
	if cond != nil && cond.Status == v1.ConditionTrue && cond.Reason == reason {
		return true
	}
	return false
}

func (exp *Experiment) HasRunningTrials() bool {
	return exp.Status.TrialsRunning != 0
}

func (exp *Experiment) GetLastConditionType() (ExperimentConditionType, error) {
	if len(exp.Status.Conditions) > 0 {
		return exp.Status.Conditions[len(exp.Status.Conditions)-1].Type, nil
	}
	return "", errors.New("Experiment doesn't have any condition")
}

func (exp *Experiment) setCondition(conditionType ExperimentConditionType, status v1.ConditionStatus, reason, message string) {

	newCond := newCondition(conditionType, status, reason, message)
	currentCond := getCondition(exp, conditionType)
	// Do nothing if condition doesn't change
	if currentCond != nil && currentCond.Status == newCond.Status && currentCond.Reason == newCond.Reason {
		return
	}

	// Do not update lastTransitionTime if the status of the condition doesn't change.
	if currentCond != nil && currentCond.Status == newCond.Status {
		newCond.LastTransitionTime = currentCond.LastTransitionTime
	}

	exp.removeCondition(conditionType)
	exp.Status.Conditions = append(exp.Status.Conditions, newCond)
}

func (exp *Experiment) MarkExperimentStatusCreated(reason, message string) {
	exp.setCondition(ExperimentCreated, v1.ConditionTrue, reason, message)
}

func (exp *Experiment) MarkExperimentStatusRunning(reason, message string) {
	//exp.removeCondition(ExperimentRestarting)
	exp.setCondition(ExperimentRunning, v1.ConditionTrue, reason, message)
}

func (exp *Experiment) MarkExperimentStatusRestarting(reason, message string) {
	exp.removeCondition(ExperimentSucceeded)
	exp.removeCondition(ExperimentFailed)
	exp.setCondition(ExperimentRestarting, v1.ConditionTrue, reason, message)
}

func (exp *Experiment) MarkExperimentStatusSucceeded(reason, message string) {
	currentCond := getCondition(exp, ExperimentRunning)
	if currentCond != nil {
		exp.setCondition(ExperimentRunning, v1.ConditionFalse, currentCond.Reason, currentCond.Message)
	}
	exp.setCondition(ExperimentSucceeded, v1.ConditionTrue, reason, message)

}

func (exp *Experiment) MarkExperimentStatusFailed(reason, message string) {
	currentCond := getCondition(exp, ExperimentRunning)
	if currentCond != nil {
		exp.setCondition(ExperimentRunning, v1.ConditionFalse, currentCond.Reason, currentCond.Message)
	}
	exp.setCondition(ExperimentFailed, v1.ConditionTrue, reason, message)
}
