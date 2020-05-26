package util

import (
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

// SuggestionLabels returns the expected suggestion labels.
func SuggestionLabels(instance *suggestionsv1beta1.Suggestion) map[string]string {
	res := make(map[string]string)
	for k, v := range instance.Labels {
		res[k] = v
	}
	res[consts.LabelDeploymentName] = GetAlgorithmDeploymentName(instance)
	res[consts.LabelExperimentName] = instance.Name
	res[consts.LabelSuggestionName] = instance.Name

	return res
}

// TrialLabels returns the expected trial labels.
func TrialLabels(instance *experimentsv1beta1.Experiment) map[string]string {
	res := make(map[string]string)
	for k, v := range instance.Labels {
		res[k] = v
	}
	res[consts.LabelExperimentName] = instance.Name

	return res
}
