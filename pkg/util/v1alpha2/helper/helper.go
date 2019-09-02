package helper

import (
	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	suggestionsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/suggestions/v1alpha2"
	"github.com/kubeflow/katib/pkg/controller/v1alpha2/consts"
)

func SuggestionLabels(instance *suggestionsv1alpha2.Suggestion) map[string]string {
	return map[string]string{
		"deployment":               instance.Name,
		consts.LabelExperimentName: instance.Name,
		consts.LabelSuggestionName: instance.Name,
	}
}

func TrialLabels(instance *experimentsv1alpha2.Experiment) map[string]string {
	return map[string]string{
		consts.LabelExperimentName: instance.Name,
	}
}
