package suggestionclient

import (
	suggestionsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/suggestions/v1alpha2"
)

type SuggestionClient interface {
	SyncAssignments(instance *suggestionsv1alpha2.Suggestion) error
}

type General struct {
}

func New() SuggestionClient {
	return &General{}
}

func (g *General) SyncAssignments(instance *suggestionsv1alpha2.Suggestion) error {
	return nil
}
