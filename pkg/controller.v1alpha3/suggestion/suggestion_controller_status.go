package suggestion

import (
	"context"

	suggestionsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1alpha3"
	"k8s.io/apimachinery/pkg/api/equality"
)

func (r *ReconcileSuggestion) updateStatus(s *suggestionsv1alpha3.Suggestion, oldS *suggestionsv1alpha3.Suggestion) error {
	if !equality.Semantic.DeepEqual(s.Status, oldS.Status) {
		if err := r.Status().Update(context.TODO(), s); err != nil {
			return err
		}
	}
	return nil
}
