package suggestion

import (
	"fmt"

	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	api_pb "github.com/kubeflow/katib/pkg/api/v1alpha2"
)

type Suggestion interface {
	GetSuggestions(instance *experimentsv1alpha2.Experiment, addCount int) ([]*api_pb.Trial, error)
}

type General struct {
}

func New() Suggestion {
	return &General{}
}

func (g *General) GetSuggestions(instance *experimentsv1alpha2.Experiment, addCount int) ([]*api_pb.Trial, error) {
	return nil, fmt.Errorf("Not implemented")
}
