package util

import (
	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	suggestionsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/consts"
)

func SuggestionLabels(instance *suggestionsv1alpha3.Suggestion) map[string]string {
	return map[string]string{
		"deployment":               instance.Name,
		consts.LabelExperimentName: instance.Name,
		consts.LabelSuggestionName: instance.Name,
	}
}

func TrialLabels(instance *experimentsv1alpha3.Experiment) map[string]string {
	return map[string]string{
		consts.LabelExperimentName: instance.Name,
	}
}
