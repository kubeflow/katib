package fake

import (
	experimentv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	api_pb "github.com/kubeflow/katib/pkg/api/v1alpha2"
	"github.com/kubeflow/katib/pkg/controller/v1alpha2/experiment/suggestion"
)

type Fake struct {
}

func New() suggestion.Suggestion {
	return &Fake{}
}

func (k *Fake) GetSuggestions(instance *experimentv1alpha2.Experiment, addCount int) ([]*api_pb.Trial, error) {
	return nil, nil
}
