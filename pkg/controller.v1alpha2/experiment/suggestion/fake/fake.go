package fake

import (
	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha2"
	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1alpha2"
	"github.com/kubeflow/katib/pkg/controller.v1alpha2/experiment/suggestion"
)

type Fake struct {
}

func New() suggestion.Suggestion {
	return &Fake{}
}

func (k *Fake) GetSuggestions(instance *experimentsv1alpha2.Experiment, addCount int32) ([]*api_pb.Trial, error) {
	return nil, nil
}
