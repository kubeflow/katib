package v1alpha3

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (suggestion *Suggestion) IsExhausted() bool {
	return hasCondition(suggestion, SuggestionExhausted)
}

func (suggestion *Suggestion) IsCompleted() bool {
	return suggestion.IsSucceeded() || suggestion.IsFailed() || suggestion.IsExhausted()
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

func (suggestion *Suggestion) MarkSuggestionStatusRunning(reason, message string) {
	//suggestion.removeCondition(SuggestionRestarting)
	suggestion.setCondition(SuggestionRunning, v1.ConditionTrue, reason, message)
}

func (suggestion *Suggestion) MarkSuggestionStatusSucceeded(reason, message string) {
	currentCond := getCondition(suggestion, SuggestionRunning)
	if currentCond != nil {
		suggestion.setCondition(SuggestionRunning, v1.ConditionFalse, currentCond.Reason, currentCond.Message)
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

func (suggestion *Suggestion) MarkSuggestionStatusExhausted(reason, message string) {
	currentCond := getCondition(suggestion, SuggestionRunning)
	if currentCond != nil {
		suggestion.setCondition(SuggestionRunning, v1.ConditionFalse, currentCond.Reason, currentCond.Message)
	}
	suggestion.setCondition(SuggestionExhausted, v1.ConditionTrue, reason, message)
}

func (suggestion *Suggestion) MarkSuggestionStatusDeploymentReady(status v1.ConditionStatus, reason, message string) {
	suggestion.setCondition(SuggestionDeploymentReady, status, reason, message)
}
