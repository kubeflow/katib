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
	"database/sql"

	commonv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/common/v1alpha2"
	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	api_pb "github.com/kubeflow/katib/pkg/api/v1alpha2"
)

func CreateExperimentInDB(instance *experimentsv1alpha2.Experiment) error {
	//TODO: Save experiment in to db
	// experiment := GetExperimentConf(instance)

	return nil
}

func UpdateExperimentStatusInDB(instance *experimentsv1alpha2.Experiment) error {

	return nil
}

func GetExperimentFromDB(instance *experimentsv1alpha2.Experiment) (*api_pb.GetExperimentReply, error) {
	return nil, sql.ErrNoRows
}

func GetExperimentConf(instance *experimentsv1alpha2.Experiment) *api_pb.Experiment {
	experiment := &api_pb.Experiment{
		ExperimentSpec: &api_pb.ExperimentSpec{
			Objective: &api_pb.ObjectiveSpec{
				AdditionalMetricNames: []string{},
			},
			Algorithm: &api_pb.AlgorithmSpec{
				AlgorithmSetting: []*api_pb.AlgorithmSetting{},
			},
		},
	}

	experiment.Name = instance.Name

	//Populate Objective
	switch instance.Spec.Objective.Type {
	case commonv1alpha2.ObjectiveTypeMaximize:
		experiment.ExperimentSpec.Objective.Type = api_pb.ObjectiveType_MAXIMIZE
	case commonv1alpha2.ObjectiveTypeMinimize:
		experiment.ExperimentSpec.Objective.Type = api_pb.ObjectiveType_MINIMIZE
	default:
		experiment.ExperimentSpec.Objective.Type = api_pb.ObjectiveType_UNKNOWN

	}
	experiment.ExperimentSpec.Objective.Goal = float32(*instance.Spec.Objective.Goal)
	experiment.ExperimentSpec.Objective.ObjectiveMetricName = instance.Spec.Objective.ObjectiveMetricName
	for _, m := range instance.Spec.Objective.AdditionalMetricNames {
		experiment.ExperimentSpec.Objective.AdditionalMetricNames = append(experiment.ExperimentSpec.Objective.AdditionalMetricNames, m)
	}

	//Populate Algorithm Spec
	experiment.ExperimentSpec.Algorithm.AlgorithmName = instance.Spec.Algorithm.AlgorithmName

	for _, as := range instance.Spec.Algorithm.AlgorithmSettings {
		experiment.ExperimentSpec.Algorithm.AlgorithmSetting = append(
			experiment.ExperimentSpec.Algorithm.AlgorithmSetting,
			&api_pb.AlgorithmSetting{
				Name:  as.Name,
				Value: as.Value,
			})
	}

	//Populate HP Experiment
	if instance.Spec.Parameters != nil {

		for _, p := range instance.Spec.Parameters {
			parameter := &api_pb.ParameterSpec{
				FeasibleSpace: &api_pb.FeasibleSpace{},
			}
			parameter.Name = p.Name
			parameter.FeasibleSpace.Min = p.FeasibleSpace.Min
			parameter.FeasibleSpace.Max = p.FeasibleSpace.Max
			parameter.FeasibleSpace.List = p.FeasibleSpace.List
			parameter.FeasibleSpace.Step = p.FeasibleSpace.Step

			switch p.ParameterType {
			case experimentsv1alpha2.ParameterTypeCategorical:
				parameter.ParameterType = api_pb.ParameterType_CATEGORICAL
			case experimentsv1alpha2.ParameterTypeDiscrete:
				parameter.ParameterType = api_pb.ParameterType_DISCRETE
			case experimentsv1alpha2.ParameterTypeDouble:
				parameter.ParameterType = api_pb.ParameterType_DOUBLE
			case experimentsv1alpha2.ParameterTypeInt:
				parameter.ParameterType = api_pb.ParameterType_INT
			case experimentsv1alpha2.ParameterTypeUnknown:
				parameter.ParameterType = api_pb.ParameterType_UNKNOWN_TYPE
			}
			experiment.ExperimentSpec.ParameterSpecs.Parameters = append(experiment.ExperimentSpec.ParameterSpecs.Parameters, parameter)
		}

	}

	//Populate NAS Experiment
	if instance.Spec.NasConfig != nil {

		nasConfig := &api_pb.NasConfig{
			GraphConfig: &api_pb.GraphConfig{},
			Operations: &api_pb.NasConfig_Operations{
				Operation: []*api_pb.Operation{},
			},
		}

		nasConfig.GraphConfig.NumLayers = *instance.Spec.NasConfig.GraphConfig.NumLayers

		for _, i := range instance.Spec.NasConfig.GraphConfig.InputSizes {
			nasConfig.GraphConfig.InputSizes = append(nasConfig.GraphConfig.InputSizes, i)
		}

		for _, o := range instance.Spec.NasConfig.GraphConfig.OutputSizes {
			nasConfig.GraphConfig.OutputSizes = append(nasConfig.GraphConfig.OutputSizes, o)
		}

		for _, op := range instance.Spec.NasConfig.Operations {
			operation := &api_pb.Operation{
				ParameterSpecs: &api_pb.Operation_ParameterSpecs{
					Parameters: []*api_pb.ParameterSpec{},
				},
			}

			operation.OperationType = op.OperationType

			for _, p := range op.Parameters {
				parameter := &api_pb.ParameterSpec{
					FeasibleSpace: &api_pb.FeasibleSpace{},
				}
				parameter.Name = p.Name
				parameter.FeasibleSpace.Min = p.FeasibleSpace.Min
				parameter.FeasibleSpace.Max = p.FeasibleSpace.Max
				parameter.FeasibleSpace.List = p.FeasibleSpace.List
				parameter.FeasibleSpace.Step = p.FeasibleSpace.Step

				switch p.ParameterType {
				case experimentsv1alpha2.ParameterTypeCategorical:
					parameter.ParameterType = api_pb.ParameterType_CATEGORICAL
				case experimentsv1alpha2.ParameterTypeDiscrete:
					parameter.ParameterType = api_pb.ParameterType_DISCRETE
				case experimentsv1alpha2.ParameterTypeDouble:
					parameter.ParameterType = api_pb.ParameterType_DOUBLE
				case experimentsv1alpha2.ParameterTypeInt:
					parameter.ParameterType = api_pb.ParameterType_INT
				case experimentsv1alpha2.ParameterTypeUnknown:
					parameter.ParameterType = api_pb.ParameterType_UNKNOWN_TYPE
				}
				operation.ParameterSpecs.Parameters = append(operation.ParameterSpecs.Parameters, parameter)
			}
			nasConfig.Operations.Operation = append(nasConfig.Operations.Operation, operation)
		}

		experiment.ExperimentSpec.NasConfig = nasConfig
	}

	return experiment

}
