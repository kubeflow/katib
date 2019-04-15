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
	//v1 "k8s.io/api/core/v1"

	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	trialsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
)

var log = logf.Log.WithName("controller")

const (
	ExperimentCreatedReason   = "ExperimentCreated"
	ExperimentRunningReason   = "ExperimentRunning"
	ExperimentSucceededReason = "ExperimentSucceeded"
	ExperimentFailedReason    = "ExperimentFailed"
	ExperimentKilledReason    = "ExperimentKilled"
)

func UpdateExperimentStatus(instance *experimentsv1alpha2.Experiment, trials *trialsv1alpha2.TrialList) {

	isObjectiveGoalReached := updateTrialsSummary(instance, trials)

	updateExperimentStatusCondition(instance, isObjectiveGoalReached)

}

func updateTrialsSummary(instance *experimentsv1alpha2.Experiment, trials *trialsv1alpha2.TrialList) bool {

	var trialsPending, trialsRunning, trialsSucceeded, trialsFailed, trialsKilled int
	var bestTrialIndex int
	var bestTrialValue float64
	isObjectiveGoalReached := false
	objectiveMetricName := instance.Spec.Objective.ObjectiveMetricName
	objectiveValueGoal := instance.Spec.Objective.Goal
	objectiveType := instance.Spec.Objective.Type
	for index, trial := range trials.Items {
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

		for _, metric := range trial.Status.Observation.Metrics {
			if objectiveMetricName == metric.Name {
				if objectiveType == experimentsv1alpha2.ObjectiveTypeMinimize {
					if bestTrialValue < metric.Value {
						bestTrialValue = metric.Value
						bestTrialIndex = index
					}
					if bestTrialValue <= objectiveValueGoal {
						isObjectiveGoalReached = true
					}
				} else if objectiveType == experimentsv1alpha2.ObjectiveTypeMaximize {
					if bestTrialValue > metric.Value {
						bestTrialValue = metric.Value
						bestTrialIndex = index
					}
					if bestTrialValue >= objectiveValueGoal {
						isObjectiveGoalReached = true
					}
				}
			}
		}
	}
	if len(trials.Items) > 0 {
		instance.Status.TrialsPending = trialsPending
		instance.Status.TrialsRunning = trialsRunning
		instance.Status.TrialsSucceeded = trialsSucceeded
		instance.Status.TrialsFailed = trialsFailed
		instance.Status.TrialsKilled = trialsKilled

		bestTrial := trials.Items[bestTrialIndex]

		instance.Status.CurrentOptimalTrial.ParameterAssignments = []trialsv1alpha2.ParameterAssignment{}
		for _, parameterAssigment := range bestTrial.Spec.ParameterAssignments {
			instance.Status.CurrentOptimalTrial.ParameterAssignments = append(instance.Status.CurrentOptimalTrial.ParameterAssignments, parameterAssigment)
		}

		instance.Status.CurrentOptimalTrial.Observation.Metrics = []trialsv1alpha2.Metric{}
		for _, metric := range bestTrial.Status.Observation.Metrics {
			instance.Status.CurrentOptimalTrial.Observation.Metrics = append(instance.Status.CurrentOptimalTrial.Observation.Metrics, metric)
		}
	}
	return isObjectiveGoalReached
}

func updateExperimentStatusCondition(instance *experimentsv1alpha2.Experiment, isObjectiveGoalReached bool) {

	completedTrialsCount := instance.Status.TrialsSucceeded + instance.Status.TrialsFailed + instance.Status.TrialsKilled
	failedTrialsCount := instance.Status.TrialsFailed

	if (instance.Spec.MaxTrialCount != nil) && (completedTrialsCount >= *instance.Spec.MaxTrialCount) {
		msg := "Experiment has succeeded"
		instance.MarkExperimentStatusSucceeded(ExperimentSucceededReason, msg)
		return
	}

	if (instance.Spec.MaxFailedTrialCount != nil) && (failedTrialsCount >= *instance.Spec.MaxFailedTrialCount) {
		msg := "Experiment has failed"
		instance.MarkExperimentStatusFailed(ExperimentFailedReason, msg)
		return
	}

	if isObjectiveGoalReached {
		msg := "Experiment has succeeded"
		instance.MarkExperimentStatusSucceeded(ExperimentSucceededReason, msg)
		return
	}
	msg := "Experiment is running"
	instance.MarkExperimentStatusRunning(ExperimentRunningReason, msg)
}
