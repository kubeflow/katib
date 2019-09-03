package fake

import (
	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/experiment/suggestion"
)

type Fake struct {
}

func New() suggestion.Suggestion {
	return &Fake{}
}

func (k *Fake) GetSuggestions(instance *experimentsv1alpha3.Experiment, addCount int32) ([]*api_pb.Trial, error) {
	return nil, nil
}
