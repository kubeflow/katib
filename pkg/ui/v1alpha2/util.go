package ui

import (
	experimentv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
)

func getExperimentCurrentCondition(experiment *experimentv1alpha2.Experiment) experimentv1alpha2.ExperimentConditionType {
	var condition string

	if experiment.IsFailed() {
		return experimentv1alpha2.ExperimentFailed
	}
	if experiment.IsSucceeded() {
		return experimentv1alpha2.ExperimentSucceeded
	}
	//TODO: Add logic here or in experiments api util
	// if experiment.IsRunning() {
	// 	return experimentv1alpha2.ExperimentRunning
	// }
	return experimentv1alpha2.ExperimentRunning
}
