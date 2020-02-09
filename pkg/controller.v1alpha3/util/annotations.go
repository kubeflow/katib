package util

import (
	suggestionsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/consts"
)

// SuggestionAnnotations returns the expected suggestion annotations.
func SuggestionAnnotations(instance *suggestionsv1alpha3.Suggestion) map[string]string {
	res := make(map[string]string)
	for k, v := range instance.Annotations {
		res[k] = v
	}
	res[consts.AnnotationIstioSidecarInjectName] = consts.AnnotationIstioSidecarInjectValue

	return res
}
