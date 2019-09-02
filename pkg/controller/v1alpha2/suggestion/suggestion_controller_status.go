package suggestion

import (
	"context"

	suggestionsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/suggestions/v1alpha2"
	"k8s.io/apimachinery/pkg/api/equality"
)

func (r *ReconcileSuggestion) updateStatus(s *suggestionsv1alpha2.Suggestion, oldS *suggestionsv1alpha2.Suggestion) error {
	if !equality.Semantic.DeepEqual(s.Status, oldS.Status) {
		if err := r.Status().Update(context.TODO(), s); err != nil {
			return err
		}
	}
	return nil
}
