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
	ExperimentCreatedReason   = "ExperimentCreated"
	ExperimentRunningReason   = "ExperimentRunning"
	ExperimentSucceededReason = "ExperimentSucceeded"
	ExperimentFailedReason    = "ExperimentFailed"
	ExperimentKilledReason    = "ExperimentKilled"
)

func UpdateExperimentStatus(collector *ExperimentsCollector, instance *experimentsv1alpha3.Experiment, trials *trialsv1alpha3.TrialList) error {

	isObjectiveGoalReached := updateTrialsSummary(instance, trials)

	if !instance.IsCompleted() {
		UpdateExperimentStatusCondition(collector, instance, isObjectiveGoalReached, false)
	}
	return nil

}

func updateTrialsSummary(instance *experimentsv1alpha3.Experiment, trials *trialsv1alpha3.TrialList) bool {

	var totalTrials, trialsPending, trialsRunning, trialsSucceeded, trialsFailed, trialsKilled int32
	var bestTrialValue float64
	bestTrialIndex := -1
	isObjectiveGoalReached := false
	objectiveValueGoal := *instance.Spec.Objective.Goal
	objectiveType := instance.Spec.Objective.Type
	objectiveMetricName := instance.Spec.Objective.ObjectiveMetricName

	for index, trial := range trials.Items {
		totalTrials++
		if trial.IsKilled() {
			trialsKilled++
		} else if trial.IsFailed() {
			trialsFailed++
		} else if trial.IsSucceeded() {
			trialsSucceeded++
		} else if trial.IsRunning() {
			trialsRunning++
		} else {
			trialsPending++
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

	instance.Status.Trials = totalTrials
	instance.Status.TrialsPending = trialsPending
	instance.Status.TrialsRunning = trialsRunning
	instance.Status.TrialsSucceeded = trialsSucceeded
	instance.Status.TrialsFailed = trialsFailed
	instance.Status.TrialsKilled = trialsKilled

	// if best trial is set
	if bestTrialIndex != -1 {
		bestTrial := trials.Items[bestTrialIndex]

		instance.Status.CurrentOptimalTrial.ParameterAssignments = []commonv1alpha3.ParameterAssignment{}
		for _, parameterAssigment := range bestTrial.Spec.ParameterAssignments {
			instance.Status.CurrentOptimalTrial.ParameterAssignments = append(instance.Status.CurrentOptimalTrial.ParameterAssignments, parameterAssigment)
		}

		instance.Status.CurrentOptimalTrial.Observation.Metrics = []commonv1alpha3.Metric{}
		for _, metric := range bestTrial.Status.Observation.Metrics {
			instance.Status.CurrentOptimalTrial.Observation.Metrics = append(instance.Status.CurrentOptimalTrial.Observation.Metrics, metric)
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

func UpdateExperimentStatusCondition(collector *ExperimentsCollector, instance *experimentsv1alpha3.Experiment, isObjectiveGoalReached bool, getSuggestionDone bool) {

	completedTrialsCount := instance.Status.TrialsSucceeded + instance.Status.TrialsFailed + instance.Status.TrialsKilled
	failedTrialsCount := instance.Status.TrialsFailed
	now := metav1.Now()

	if isObjectiveGoalReached {
		msg := "Experiment has succeeded because Objective goal has reached"
		instance.MarkExperimentStatusSucceeded(ExperimentSucceededReason, msg)
		instance.Status.CompletionTime = &now
		collector.IncreaseExperimentsSucceededCount(instance.Namespace)
		return
	}

	if (instance.Spec.MaxTrialCount != nil) && (completedTrialsCount >= *instance.Spec.MaxTrialCount) {
		msg := "Experiment has succeeded because max trial count has reached"
		instance.MarkExperimentStatusSucceeded(ExperimentSucceededReason, msg)
		instance.Status.CompletionTime = &now
		collector.IncreaseExperimentsSucceededCount(instance.Namespace)
		return
	}

	if getSuggestionDone && (instance.Status.TrialsPending+instance.Status.TrialsRunning) == 0 {
		msg := "Experiment has succeeded because suggestion service has reached the end"
		instance.MarkExperimentStatusSucceeded(ExperimentSucceededReason, msg)
		instance.Status.CompletionTime = &now
		collector.IncreaseExperimentsSucceededCount(instance.Namespace)
		return
	}

	if (instance.Spec.MaxFailedTrialCount != nil) && (failedTrialsCount >= *instance.Spec.MaxFailedTrialCount) {
		msg := "Experiment has failed because max failed count has reached"
		instance.MarkExperimentStatusFailed(ExperimentFailedReason, msg)
		instance.Status.CompletionTime = &now
		collector.IncreaseExperimentsFailedCount(instance.Namespace)
		return
	}

	msg := "Experiment is running"
	instance.MarkExperimentStatusRunning(ExperimentRunningReason, msg)
}
