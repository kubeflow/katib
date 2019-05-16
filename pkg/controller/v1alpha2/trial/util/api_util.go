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
	commonv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/common/v1alpha2"
	trialsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
	api_pb "github.com/kubeflow/katib/pkg/api/v1alpha2"
	common "github.com/kubeflow/katib/pkg/common/v1alpha2"
)

func CreateTrialInDB(instance *trialsv1alpha2.Trial) error {
	trial := GetTrialConf(instance)
	request := &api_pb.RegisterTrialRequest{
		Trial: trial,
	}
	if _, err := common.RegisterTrial(request); err != nil {
		return err
	}
	return nil
}

func UpdateTrialStatusInDB(instance *trialsv1alpha2.Trial) error {

	return nil
}

func GetTrialObservation(instance *trialsv1alpha2.Trial) error {

	return nil
}

func GetTrialConf(instance *trialsv1alpha2.Trial) *api_pb.Trial {
	trial := &api_pb.Trial{
		Spec: &api_pb.TrialSpec{
			Objective: &api_pb.ObjectiveSpec{
				AdditionalMetricNames: []string{},
			},
			ParameterAssignments: &api_pb.TrialSpec_ParameterAssignments{
				Assignments: []*api_pb.ParameterAssignment{},
			},
		},
	}
	trial.Name = instance.Name

	trial.Spec.ExperimentName = instance.Labels["experiment"]

	//Populate Objective
	switch instance.Spec.Objective.Type {
	case commonv1alpha2.ObjectiveTypeMaximize:
		trial.Spec.Objective.Type = api_pb.ObjectiveType_MAXIMIZE
	case commonv1alpha2.ObjectiveTypeMinimize:
		trial.Spec.Objective.Type = api_pb.ObjectiveType_MINIMIZE
	default:
		trial.Spec.Objective.Type = api_pb.ObjectiveType_UNKNOWN

	}
	trial.Spec.Objective.Goal = float32(*instance.Spec.Objective.Goal)
	trial.Spec.Objective.ObjectiveMetricName = instance.Spec.Objective.ObjectiveMetricName
	for _, m := range instance.Spec.Objective.AdditionalMetricNames {
		trial.Spec.Objective.AdditionalMetricNames = append(trial.Spec.Objective.AdditionalMetricNames, m)
	}

	//Populate Parameter Assignments
	for _, p := range instance.Spec.ParameterAssignments {
		trial.Spec.ParameterAssignments.Assignments = append(
			trial.Spec.ParameterAssignments.Assignments,
			&api_pb.ParameterAssignment{
				Name:  p.Name,
				Value: p.Value,
			})
	}

	trial.Spec.RunSpec = instance.Spec.RunSpec

	trial.Spec.MetricsCollectorSpec = instance.Spec.MetricsCollectorSpec

	return trial
}
