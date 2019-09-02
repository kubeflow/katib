package suggestion

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	suggestionsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/suggestions/v1alpha2"
	api_pb "github.com/kubeflow/katib/pkg/api/v1alpha2"
	common "github.com/kubeflow/katib/pkg/common/v1alpha2"
)

// Suggestion is the interface for suggestions in Experiment controller.
type Suggestion interface {
	CreateSuggestion(instance *experimentsv1alpha2.Experiment) error
	GetSuggestions(instance *experimentsv1alpha2.Experiment, addCount int32) ([]*api_pb.Trial, error)
}

type General struct {
	scheme *runtime.Scheme
	client.Client
}

func New(scheme *runtime.Scheme, client client.Client) Suggestion {
	return &General{
		scheme: scheme,
		Client: client,
	}
}

func (g *General) CreateSuggestion(instance *experimentsv1alpha2.Experiment) error {
	suggestion := &suggestionsv1alpha2.Suggestion{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: suggestionsv1alpha2.SuggestionSpec{
			AlgorithmSpec: instance.Spec.Algorithm,
			Suggestions:   0,
		},
	}

	if err := controllerutil.SetControllerReference(instance, suggestion, g.scheme); err != nil {
		return err
	}

	if err := g.Create(context.TODO(), suggestion); err != nil {
		return err
	}
	return nil
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
