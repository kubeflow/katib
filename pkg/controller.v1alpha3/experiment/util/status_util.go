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

package util

import (
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	commonv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/common/v1alpha3"
	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	trialsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var log = logf.Log.WithName("experiment-status-util")

const (
	ExperimentCreatedReason              = "ExperimentCreated"
	ExperimentRunningReason              = "ExperimentRunning"
	ExperimentRestartingReason           = "ExperimentRestarting"
	ExperimentGoalReachedReason          = "ExperimentGoalReached"
	ExperimentMaxTrialsReachedReason     = "ExperimentMaxTrialsReached"
	ExperimentSuggestionEndReachedReason = "ExperimentSuggestionEndReached"
	ExperimentFailedReason               = "ExperimentFailed"
	ExperimentKilledReason               = "ExperimentKilled"
)

func UpdateExperimentStatus(collector *ExperimentsCollector, instance *experimentsv1alpha3.Experiment, trials *trialsv1alpha3.TrialList) error {

	isObjectiveGoalReached := updateTrialsSummary(instance, trials)

	if !instance.IsCompleted() {
		UpdateExperimentStatusCondition(collector, instance, isObjectiveGoalReached, false)
	}
	return nil

}

func updateTrialsSummary(instance *experimentsv1alpha3.Experiment, trials *trialsv1alpha3.TrialList) bool {

	var totalTrials int32
	var bestTrialValue float64
	sts := &instance.Status
	sts.RunningTrials, sts.PendingTrials, sts.FailedTrials, sts.SucceededTrials, sts.KilledTrials = nil, nil, nil, nil, nil
	bestTrialIndex := -1
	isObjectiveGoalReached := false
	objectiveValueGoal := *instance.Spec.Objective.Goal
	objectiveType := instance.Spec.Objective.Type
	objectiveMetricName := instance.Spec.Objective.ObjectiveMetricName

	for index, trial := range trials.Items {
		totalTrials++
		if trial.IsKilled() {
			sts.KilledTrials = append(sts.KilledTrials, trial.Name)
		} else if trial.IsFailed() {
			sts.KilledTrials = append(sts.KilledTrials, trial.Name)
		} else if trial.IsSucceeded() {
			sts.SucceededTrials = append(sts.SucceededTrials, trial.Name)
		} else if trial.IsRunning() {
			sts.RunningTrials = append(sts.RunningTrials, trial.Name)
		} else {
			sts.PendingTrials = append(sts.PendingTrials, trial.Name)
		}

		objectiveMetricValue := getObjectiveMetricValue(trial, objectiveMetricName)
		if objectiveMetricValue == nil {
			continue
		}

		//intialize vars to objective metric value of the first trial
		if bestTrialIndex == -1 {
			bestTrialValue = *objectiveMetricValue
			bestTrialIndex = index
		}

		if objectiveType == commonv1alpha3.ObjectiveTypeMinimize {
			if *objectiveMetricValue < bestTrialValue {
				bestTrialValue = *objectiveMetricValue
				bestTrialIndex = index
			}
			if bestTrialValue <= objectiveValueGoal {
				isObjectiveGoalReached = true
			}
		} else if objectiveType == commonv1alpha3.ObjectiveTypeMaximize {
			if *objectiveMetricValue > bestTrialValue {
				bestTrialValue = *objectiveMetricValue
				bestTrialIndex = index
			}
			if bestTrialValue >= objectiveValueGoal {
				isObjectiveGoalReached = true
			}
		}
	}

	sts.Trials = totalTrials

	// if best trial is set
	if bestTrialIndex != -1 {
		bestTrial := trials.Items[bestTrialIndex]

		sts.CurrentOptimalTrial.BestTrialName = bestTrial.Name

		sts.CurrentOptimalTrial.ParameterAssignments = []commonv1alpha3.ParameterAssignment{}
		for _, parameterAssigment := range bestTrial.Spec.ParameterAssignments {
			sts.CurrentOptimalTrial.ParameterAssignments = append(sts.CurrentOptimalTrial.ParameterAssignments, parameterAssigment)
		}

		sts.CurrentOptimalTrial.Observation.Metrics = []commonv1alpha3.Metric{}
		for _, metric := range bestTrial.Status.Observation.Metrics {
			sts.CurrentOptimalTrial.Observation.Metrics = append(sts.CurrentOptimalTrial.Observation.Metrics, metric)
		}
	}
	return isObjectiveGoalReached
}

func getObjectiveMetricValue(trial trialsv1alpha3.Trial, objectiveMetricName string) *float64 {
	if trial.Status.Observation == nil {
		return nil
	}
	for _, metric := range trial.Status.Observation.Metrics {
		if objectiveMetricName == metric.Name {
			return &metric.Value
		}
	}
	return nil
}

// UpdateExperimentStatusCondition updates the experiment status.
func UpdateExperimentStatusCondition(collector *ExperimentsCollector, instance *experimentsv1alpha3.Experiment, isObjectiveGoalReached bool, getSuggestionDone bool) {
	completedCount := int32(len(instance.Status.SucceededTrials) + len(instance.Status.FailedTrials) + len(instance.Status.KilledTrials))
	failedCount := int32(len(instance.Status.FailedTrials))
	activeCount := int32(len(instance.Status.PendingTrials) + len(instance.Status.RunningTrials))
	now := metav1.Now()

	if isObjectiveGoalReached {
		msg := "Experiment has succeeded because Objective goal has reached"
		instance.MarkExperimentStatusSucceeded(ExperimentGoalReachedReason, msg)
		instance.Status.CompletionTime = &now
		collector.IncreaseExperimentsSucceededCount(instance.Namespace)
		return
	}

	// First check if MaxFailedTrialCount is reached.
	if (instance.Spec.MaxFailedTrialCount != nil) && (failedCount > *instance.Spec.MaxFailedTrialCount) {
		msg := "Experiment has failed because max failed count has reached"
		instance.MarkExperimentStatusFailed(ExperimentFailedReason, msg)
		instance.Status.CompletionTime = &now
		collector.IncreaseExperimentsFailedCount(instance.Namespace)
		return
	}

	// Then Check if MaxTrialCount is reached.
	if (instance.Spec.MaxTrialCount != nil) && (completedCount >= *instance.Spec.MaxTrialCount) {
		msg := "Experiment has succeeded because max trial count has reached"
		instance.MarkExperimentStatusSucceeded(ExperimentMaxTrialsReachedReason, msg)
		instance.Status.CompletionTime = &now
		collector.IncreaseExperimentsSucceededCount(instance.Namespace)
		return
	}

	if getSuggestionDone && activeCount == 0 {
		msg := "Experiment has succeeded because suggestion service has reached the end"
		instance.MarkExperimentStatusSucceeded(ExperimentSuggestionEndReachedReason, msg)
		instance.Status.CompletionTime = &now
		collector.IncreaseExperimentsSucceededCount(instance.Namespace)
		return
	}

	msg := "Experiment is running"
	instance.MarkExperimentStatusRunning(ExperimentRunningReason, msg)
}

func IsCompletedExperimentRestartable(instance *experimentsv1alpha3.Experiment) bool {
	if instance.IsSucceeded() && instance.IsCompletedReason(ExperimentMaxTrialsReachedReason) {
		return true
	}
	return false
}
