package suggestion

import (
	"context"

	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	"k8s.io/apimachinery/pkg/api/equality"
)

const (
	SuggestionCreatedReason      = "SuggestionCreated"
	SuggestionDeploymentReady    = "DeploymentReady"
	SuggestionDeploymentNotReady = "DeploymentNotReady"
	SuggestionRunningReason      = "SuggestionRunning"
	SuggestionSucceededReason    = "SuggestionSucceeded"
	SuggestionFailedReason       = "SuggestionFailed"
	SuggestionKilledReason       = "SuggestionKilled"
)

func (r *ReconcileSuggestion) updateStatus(s *suggestionsv1beta1.Suggestion, oldS *suggestionsv1beta1.Suggestion) error {
	if !equality.Semantic.DeepEqual(s.Status, oldS.Status) {
		if err := r.Status().Update(context.TODO(), s); err != nil {
			return err
		}
	}
	return nil
}

func (r *ReconcileSuggestion) updateStatusCondition(s *suggestionsv1beta1.Suggestion, oldS *suggestionsv1beta1.Suggestion) error {
	if !equality.Semantic.DeepEqual(s.Status.Conditions, oldS.Status.Conditions) {
		newConditions := s.Status.Conditions
		s.Status = oldS.Status
		s.Status.Conditions = newConditions
		if err := r.Status().Update(context.TODO(), s); err != nil {
			return err
		}
	}
	return nil
}
