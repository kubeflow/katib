package suggestion

import (
	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	api_pb "github.com/kubeflow/katib/pkg/api/v1alpha2"
	common "github.com/kubeflow/katib/pkg/common/v1alpha2"
)

type Suggestion interface {
	GetSuggestions(instance *experimentsv1alpha2.Experiment, addCount int32) ([]*api_pb.Trial, error)
}

type General struct {
}

func New() Suggestion {
	return &General{}
}

func (g *General) GetSuggestions(instance *experimentsv1alpha2.Experiment, addCount int32) ([]*api_pb.Trial, error) {
	request := &api_pb.GetSuggestionsRequest{
		ExperimentName: instance.Name,
		AlgorithmName:  instance.Spec.Algorithm.AlgorithmName,
		RequestNumber:  addCount,
	}
	if reply, err := common.GetSuggestions(request); err != nil {
		return nil, err
	} else {
		return reply.Trials, nil
	}
}
