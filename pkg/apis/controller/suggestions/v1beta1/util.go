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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// SuggestionRestartReason is the reason for suggestion status when experiment is restarting
	SuggestionRestartReason = "Experiment is restarting"
)

func getCondition(suggestion *Suggestion, condType SuggestionConditionType) *SuggestionCondition {
	if suggestion.Status.Conditions != nil {
		for _, condition := range suggestion.Status.Conditions {
			if condition.Type == condType {
				return &condition
			}
		}
	}
	return nil
}

func hasCondition(suggestion *Suggestion, condType SuggestionConditionType) bool {
	cond := getCondition(suggestion, condType)
	if cond != nil && cond.Status == v1.ConditionTrue {
		return true
	}
	return false
}

func (suggestion *Suggestion) removeCondition(condType SuggestionConditionType) {
	var newConditions []SuggestionCondition
	for _, c := range suggestion.Status.Conditions {

		if c.Type == condType {
			continue
		}

		newConditions = append(newConditions, c)
	}
	suggestion.Status.Conditions = newConditions
}

func newCondition(conditionType SuggestionConditionType, status v1.ConditionStatus, reason, message string) SuggestionCondition {
	return SuggestionCondition{
		Type:               conditionType,
		Status:             status,
		LastUpdateTime:     metav1.Now(),
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            message,
	}
}

func (suggestion *Suggestion) IsCreated() bool {
	return hasCondition(suggestion, SuggestionCreated)
}

func (suggestion *Suggestion) IsFailed() bool {
	return hasCondition(suggestion, SuggestionFailed)
}

func (suggestion *Suggestion) IsSucceeded() bool {
	return hasCondition(suggestion, SuggestionSucceeded)
}

func (suggestion *Suggestion) IsRunning() bool {
	return hasCondition(suggestion, SuggestionRunning)
}

// IsRestarting returns true if suggestion running status is false and reason = SuggestionRestartReason
func (suggestion *Suggestion) IsRestarting() bool {
	cond := getCondition(suggestion, SuggestionRunning)
	if cond != nil && cond.Status == v1.ConditionFalse && cond.Reason == SuggestionRestartReason {
		return true
	}
	return false
}

func (suggestion *Suggestion) IsDeploymentReady() bool {
	return hasCondition(suggestion, SuggestionDeploymentReady)
}

func (suggestion *Suggestion) IsCompleted() bool {
	return suggestion.IsSucceeded() || suggestion.IsFailed()
}

func (suggestion *Suggestion) setCondition(conditionType SuggestionConditionType, status v1.ConditionStatus, reason, message string) {

	newCond := newCondition(conditionType, status, reason, message)
	currentCond := getCondition(suggestion, conditionType)
	// Do nothing if condition doesn't change
	if currentCond != nil && currentCond.Status == newCond.Status && currentCond.Reason == newCond.Reason {
		return
	}

	// Do not update lastTransitionTime if the status of the condition doesn't change.
	if currentCond != nil && currentCond.Status == newCond.Status {
		newCond.LastTransitionTime = currentCond.LastTransitionTime
	}

	suggestion.removeCondition(conditionType)
	suggestion.Status.Conditions = append(suggestion.Status.Conditions, newCond)
}

func (suggestion *Suggestion) MarkSuggestionStatusCreated(reason, message string) {
	suggestion.setCondition(SuggestionCreated, v1.ConditionTrue, reason, message)
}

// MarkSuggestionStatusRunning sets suggestion Running status.
func (suggestion *Suggestion) MarkSuggestionStatusRunning(status v1.ConditionStatus, reason, message string) {
	// When suggestion is restrating we need to remove succeeded status from suggestion.
	// That should happen only when ResumePolicy = FromVolume
	suggestion.removeCondition(SuggestionSucceeded)
	suggestion.setCondition(SuggestionRunning, status, reason, message)
}

// MarkSuggestionStatusSucceeded sets suggestion Succeeded status to true.
// Suggestion can be succeeded only if ResumeExperiment = Never or ResumeExperiment = FromVolume
func (suggestion *Suggestion) MarkSuggestionStatusSucceeded(reason, message string) {

	// When suggestion is Succeeded suggestion Running status is false
	runningCond := getCondition(suggestion, SuggestionRunning)
	succeededReason := "Suggestion is succeeded"
	if runningCond != nil {
		msg := "Suggestion is not running"
		suggestion.setCondition(SuggestionRunning, v1.ConditionFalse, succeededReason, msg)
	}

	// When suggestion is Succeeded suggestion DeploymentReady status is false
	deploymentReadyCond := getCondition(suggestion, SuggestionDeploymentReady)
	if deploymentReadyCond != nil {
		msg := "Deployment is not ready"
		suggestion.setCondition(SuggestionDeploymentReady, v1.ConditionFalse, succeededReason, msg)
	}

	suggestion.setCondition(SuggestionSucceeded, v1.ConditionTrue, reason, message)
}

func (suggestion *Suggestion) MarkSuggestionStatusFailed(reason, message string) {
	currentCond := getCondition(suggestion, SuggestionRunning)
	if currentCond != nil {
		suggestion.setCondition(SuggestionRunning, v1.ConditionFalse, currentCond.Reason, currentCond.Message)
	}
	suggestion.setCondition(SuggestionFailed, v1.ConditionTrue, reason, message)
}

func (suggestion *Suggestion) MarkSuggestionStatusDeploymentReady(status v1.ConditionStatus, reason, message string) {
	suggestion.setCondition(SuggestionDeploymentReady, status, reason, message)
}
