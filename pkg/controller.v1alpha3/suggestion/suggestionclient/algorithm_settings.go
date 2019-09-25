package suggestionclient

import (
	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1alpha3"
	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	suggestionsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1alpha3"
	suggestionapi "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
)

// appendAlgorithmSettingsFromSuggestion appends the algorithm settings
// in suggestion to Experiment.
// Algorithm settings in suggestion will overwrite the settings in experiment.
func appendAlgorithmSettingsFromSuggestion(experiment *experimentsv1alpha3.Experiment, algoSettingsInSuggestion *common.AlgorithmSpec) {
	algoSettingsInExperiment := experiment.Spec.Algorithm
	for _, setting := range algoSettingsInSuggestion.AlgorithmSettings {
		if index, found := contains(algoSettingsInExperiment, setting.Name); found {
			// If the setting is found in Experiment, update it.
			algoSettingsInExperiment.AlgorithmSettings[index].Value = setting.Value
		} else {
			// If not found, append it.
			algoSettingsInExperiment.AlgorithmSettings = append(
				algoSettingsInExperiment.AlgorithmSettings, setting)
		}
	}
}

func updateAlgorithmSettings(suggestion *suggestionsv1alpha3.Suggestion, algorithm *suggestionapi.AlgorithmSpec) {
	if suggestion.Status.Algorithm == nil {
		suggestion.Status.Algorithm = &common.AlgorithmSpec{}
	}
	algoSettingsInSuggestion := suggestion.Status.Algorithm
	for _, setting := range algorithm.AlgorithmSetting {
		if setting != nil {
			if index, found := contains(algoSettingsInSuggestion, setting.Name); found {
				// If the setting is found in Suggestion, update it.
				algoSettingsInSuggestion.AlgorithmSettings[index].Value = setting.Value
			} else {
				// If not found, append it.
				algoSettingsInSuggestion.AlgorithmSettings = append(algoSettingsInSuggestion.AlgorithmSettings, common.AlgorithmSetting{
					Name:  setting.Name,
					Value: setting.Value,
				})
			}
		}
	}
}

func contains(algorithmSettings *common.AlgorithmSpec,
	name string) (int, bool) {
	for i, s := range algorithmSettings.AlgorithmSettings {
		if s.Name == name {
			return i, true
		}
	}
	return -1, false
}
