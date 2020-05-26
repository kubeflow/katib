package fake

import (
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/experiment/suggestion"
)

type Fake struct {
}

func New() suggestion.Suggestion {
	return &Fake{}
}

func (f *Fake) GetOrCreateSuggestion(instance *experimentsv1beta1.Experiment, suggestionRequests int32) (*suggestionsv1beta1.Suggestion, error) {
	return nil, nil
}

func (f *Fake) UpdateSuggestion(suggestion *suggestionsv1beta1.Suggestion) error {
	return nil
}

func (f *Fake) UpdateSuggestionStatus(suggestion *suggestionsv1beta1.Suggestion) error {
	return nil
}
