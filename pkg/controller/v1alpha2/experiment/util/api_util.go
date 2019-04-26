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

	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	katibapiv1alpha2 "github.com/kubeflow/katib/pkg/api/v1alpha2"
)

func CreateExperimentInDB(experiment *katibapiv1alpha2.Experiment) error {

	return nil
}

func UpdateExperimentStatusInDB(instance *experimentsv1alpha2.Experiment) error {

	return nil
}

func GetSuggestions(instance *experimentsv1alpha2.Experiment, addCount int) ([]*katibapiv1alpha2.Trial, error) {

	return nil, nil
}

func GetExperimentConf(instance *experimentsv1alpha2.Experiment) *katibapiv1alpha2.Experiment {
	jobType := getJobType(instance)
	if jobType == jobTypeNAS {
		return populateNASExperiment(instance)
	}
	return populateHPExperiment(instance)
}

func getJobType(instance *experimentsv1alpha2.Experiment) string {
	if instance.Spec.NasConfig != nil {
		return jobTypeNAS
	}
	return jobTypeHP
}

func populateHPExperiment(instance *experimentsv1alpha2.Experiment) *katibapiv1alpha2.Experiment {
	experiment := &katibapiv1alpha2.Experiment{
		ExperimentSpec: &katibapiv1alpha2.ExperimentSpec{
			Objective: &katibapiv1alpha2.ObjectiveSpec{
				AdditionalMetricsNames: []string{},
			},
			Algorithm: &katibapiv1alpha2.AlgorithmSpec{
				AlgorithmSetting: []*katibapiv1alpha2.AlgorithmSetting{},
			},
		},
	}

	populateCommonExperimentFields(instance, experiment)

	for _, p := range instance.Spec.Parameters {
		parameter := &katibapiv1alpha2.ParameterSpec{
			FeasibleSpace: &katibapiv1alpha2.FeasibleSpace{},
		}
		parameter.Name = p.Name
		parameter.FeasibleSpace.Min = p.FeasibleSpace.Min
		parameter.FeasibleSpace.Max = p.FeasibleSpace.Max
		parameter.FeasibleSpace.List = p.FeasibleSpace.List
		parameter.FeasibleSpace.Step = p.FeasibleSpace.Step

		switch p.ParameterType {
		case experimentsv1alpha2.ParameterTypeCategorical:
			parameter.ParameterType = katibapiv1alpha2.ParameterType_CATEGORICAL
		case experimentsv1alpha2.ParameterTypeDiscrete:
			parameter.ParameterType = katibapiv1alpha2.ParameterType_DISCRETE
		case experimentsv1alpha2.ParameterTypeDouble:
			parameter.ParameterType = katibapiv1alpha2.ParameterType_DOUBLE
		case experimentsv1alpha2.ParameterTypeInt:
			parameter.ParameterType = katibapiv1alpha2.ParameterType_INT
		case experimentsv1alpha2.ParameterTypeUnknown:
			parameter.ParameterType = katibapiv1alpha2.ParameterType_UNKNOWN_TYPE
		}
		experiment.ExperimentSpec.ParameterSpecs.Parameters = append(experiment.ExperimentSpec.ParameterSpecs.Parameters, parameter)
	}

	return experiment
}

func populateNASExperiment(instance *experimentsv1alpha2.Experiment) *katibapiv1alpha2.Experiment {
	experiment := &katibapiv1alpha2.Experiment{
		ExperimentSpec: &katibapiv1alpha2.ExperimentSpec{
			Objective: &katibapiv1alpha2.ObjectiveSpec{
				AdditionalMetricsNames: []string{},
			},
			Algorithm: &katibapiv1alpha2.AlgorithmSpec{
				AlgorithmSetting: []*katibapiv1alpha2.AlgorithmSetting{},
			},
		},
	}
	populateCommonExperimentFields(instance, experiment)

	nasConfig := &katibapiv1alpha2.NasConfig{
		GraphConfig: &katibapiv1alpha2.GraphConfig{},
		Operations: &katibapiv1alpha2.NasConfig_Operations{
			Operation: []*katibapiv1alpha2.Operation{},
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
		operation := &katibapiv1alpha2.Operation{
			ParameterSpecs: &katibapiv1alpha2.Operation_ParameterSpecs{
				Parameters: []*katibapiv1alpha2.ParameterSpec{},
			},
		}

		operation.OperationType = op.OperationType
		for _, p := range op.Parameters {
			parameter := &katibapiv1alpha2.ParameterSpec{
				FeasibleSpace: &katibapiv1alpha2.FeasibleSpace{},
			}
			parameter.Name = p.Name
			parameter.FeasibleSpace.Min = p.FeasibleSpace.Min
			parameter.FeasibleSpace.Max = p.FeasibleSpace.Max
			parameter.FeasibleSpace.List = p.FeasibleSpace.List
			parameter.FeasibleSpace.Step = p.FeasibleSpace.Step

			switch p.ParameterType {
			case experimentsv1alpha2.ParameterTypeCategorical:
				parameter.ParameterType = katibapiv1alpha2.ParameterType_CATEGORICAL
			case experimentsv1alpha2.ParameterTypeDiscrete:
				parameter.ParameterType = katibapiv1alpha2.ParameterType_DISCRETE
			case experimentsv1alpha2.ParameterTypeDouble:
				parameter.ParameterType = katibapiv1alpha2.ParameterType_DOUBLE
			case experimentsv1alpha2.ParameterTypeInt:
				parameter.ParameterType = katibapiv1alpha2.ParameterType_INT
			case experimentsv1alpha2.ParameterTypeUnknown:
				parameter.ParameterType = katibapiv1alpha2.ParameterType_UNKNOWN_TYPE
			}
			operation.ParameterSpecs.Parameters = append(operation.ParameterSpecs.Parameters, parameter)
		}
		nasConfig.Operations.Operation = append(nasConfig.Operations.Operation, operation)
	}

	experiment.ExperimentSpec.NasConfig = nasConfig
	return experiment
}

func populateCommonExperimentFields(instance *experimentsv1alpha2.Experiment, experiment *katibapiv1alpha2.Experiment) {
	experiment.Name = instance.ObjectMeta.Name

	//Populate Objective
	switch instance.Spec.Objective.Type {
	case experimentsv1alpha2.ObjectiveTypeMaximize:
		experiment.ExperimentSpec.Objective.Type = katibapiv1alpha2.ObjectiveType_MAXIMIZE
	case experimentsv1alpha2.ObjectiveTypeMinimize:
		experiment.ExperimentSpec.Objective.Type = katibapiv1alpha2.ObjectiveType_MINIMIZE
	default:
		experiment.ExperimentSpec.Objective.Type = katibapiv1alpha2.ObjectiveType_UNKNOWN

	}
	experiment.ExperimentSpec.Objective.Goal = float32(*instance.Spec.Objective.Goal)
	experiment.ExperimentSpec.Objective.ObjectiveMetricName = instance.Spec.Objective.ObjectiveMetricName
	for _, m := range instance.Spec.Objective.AdditionalMetricsNames {
		experiment.ExperimentSpec.Objective.AdditionalMetricsNames = append(experiment.ExperimentSpec.Objective.AdditionalMetricsNames, m)
	}

	//Populate Algorithm Spec
	experiment.ExperimentSpec.Algorithm.AlgorithmName = instance.Spec.Algorithm.AlgorithmName

	for _, as := range instance.Spec.Algorithm.AlgorithmSettings {
		experiment.ExperimentSpec.Algorithm.AlgorithmSetting = append(
			experiment.ExperimentSpec.Algorithm.AlgorithmSetting,
			&katibapiv1alpha2.AlgorithmSetting{
				Name:  as.Name,
				Value: as.Value,
			})
	}

}
