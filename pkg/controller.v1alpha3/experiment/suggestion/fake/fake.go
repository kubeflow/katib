package fake

import (
	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	suggestionsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/experiment/suggestion"
)

type Fake struct {
}

func New() suggestion.Suggestion {
	return &Fake{}
}

func (f *Fake) CreateSuggestion(instance *experimentsv1alpha3.Experiment, suggestionRequests int32) error {
	return nil
}

func (f *Fake) UpdateSuggestion(suggestion *suggestionsv1alpha3.Suggestion, suggestionRequests int32) error {
	return nil
}

func (f *Fake) GetSuggestions(suggestion *suggestionsv1alpha3.Suggestion) []suggestionsv1alpha3.TrialAssignment {
	return nil
}
