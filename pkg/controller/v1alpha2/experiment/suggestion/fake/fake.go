package fake

import (
	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	suggestionsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/suggestions/v1alpha2"
	"github.com/kubeflow/katib/pkg/controller/v1alpha2/experiment/suggestion"
)

type Fake struct {
}

func New() suggestion.Suggestion {
	return &Fake{}
}

func (f *Fake) RequestSuggestions(s *suggestionsv1alpha2.Suggestion, addcount int32) error {
	return nil
}

func (f *Fake) CreateSuggestion(instance *experimentsv1alpha2.Experiment) error {
	return nil
}

func (f *Fake) UpdateSuggestion(s *suggestionsv1alpha2.Suggestion, os *suggestionsv1alpha2.Suggestion) error {
	return nil
}

func (f *Fake) GetSuggestions(
	s *suggestionsv1alpha2.Suggestion) ([]suggestionsv1alpha2.TrialAssignment, error) {
	return nil, nil
}
