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

package suggestionclient

import (
	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	suggestionapi "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
)

// appendAlgorithmSettingsFromSuggestion appends the algorithm settings from Suggestion status to Experiment.
// Algorithm settings in Suggestion status will overwrite the settings in Experiment.
func appendAlgorithmSettingsFromSuggestion(experiment *experimentsv1beta1.Experiment, algoSettingsInSuggestion []common.AlgorithmSetting) {
	algoSettingsInExperiment := experiment.Spec.Algorithm
	for _, setting := range algoSettingsInSuggestion {
		if index, found := contains(
			algoSettingsInExperiment.AlgorithmSettings, setting.Name); found {
			// If the setting is found in Experiment, update it.
			algoSettingsInExperiment.AlgorithmSettings[index].Value = setting.Value
		} else {
			// If not found, append it.
			algoSettingsInExperiment.AlgorithmSettings = append(
				algoSettingsInExperiment.AlgorithmSettings, setting)
		}
	}
}

func updateAlgorithmSettings(suggestion *suggestionsv1beta1.Suggestion, algorithm *suggestionapi.AlgorithmSpec) {
	for _, setting := range algorithm.AlgorithmSettings {
		if setting != nil {
			if index, found := contains(suggestion.Status.AlgorithmSettings, setting.Name); found {
				// If the setting is found in Suggestion, update it.
				suggestion.Status.AlgorithmSettings[index].Value = setting.Value
			} else {
				// If not found, append it.
				suggestion.Status.AlgorithmSettings = append(suggestion.Status.AlgorithmSettings, common.AlgorithmSetting{
					Name:  setting.Name,
					Value: setting.Value,
				})
			}
		}
	}
}

func contains(algorithmSettings []common.AlgorithmSetting,
	name string) (int, bool) {
	for i, s := range algorithmSettings {
		if s.Name == name {
			return i, true
		}
	}
	return -1, false
}
