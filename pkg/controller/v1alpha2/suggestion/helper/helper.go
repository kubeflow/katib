package helper

import (
	suggestionsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/suggestions/v1alpha2"
	"github.com/kubeflow/katib/pkg/controller/v1alpha2/consts"
)

func Labels(instance *suggestionsv1alpha2.Suggestion) map[string]string {
	return map[string]string{
		"deployment":               instance.Name,
		consts.LabelExperimentName: instance.Name,
		consts.LabelSuggestionName: instance.Name,
	}
}
