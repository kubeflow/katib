package suggestion

import (
	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
	common "github.com/kubeflow/katib/pkg/common/v1alpha3"
)

type Suggestion interface {
	GetSuggestions(instance *experimentsv1alpha3.Experiment, addCount int32) ([]*api_pb.Trial, error)
}

type General struct {
}

func New() Suggestion {
	return &General{}
}

func (g *General) GetSuggestions(instance *experimentsv1alpha3.Experiment, addCount int32) ([]*api_pb.Trial, error) {
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
