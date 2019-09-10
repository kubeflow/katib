package suggestion

import (
	"context"
	//"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	//utilrand "k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	suggestionsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1alpha3"
)

var log = logf.Log.WithName("experiment-suggestion-client")

type Suggestion interface {
	CreateSuggestion(instance *experimentsv1alpha3.Experiment, suggestionRequests int32) error
	UpdateSuggestion(suggestion *suggestionsv1alpha3.Suggestion, suggestionRequests int32) error
	GetSuggestions(suggestion *suggestionsv1alpha3.Suggestion) []suggestionsv1alpha3.TrialAssignment
}

type General struct {
	scheme *runtime.Scheme
	client.Client
}

func New(scheme *runtime.Scheme, client client.Client) Suggestion {
	return &General{scheme: scheme, Client: client}
}

func (g *General) GetSuggestions(suggestion *suggestionsv1alpha3.Suggestion) []suggestionsv1alpha3.TrialAssignment {
	return suggestion.Status.Suggestions
}

func (g *General) CreateSuggestion(instance *experimentsv1alpha3.Experiment, suggestionRequests int32) error {
	logger := log.WithValues("experiment", types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace})
	suggestion := &suggestionsv1alpha3.Suggestion{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: suggestionsv1alpha3.SuggestionSpec{
			AlgorithmSpec: instance.Spec.Algorithm,
			Requests:      suggestionRequests,
		},
	}

	if err := controllerutil.SetControllerReference(instance, suggestion, g.scheme); err != nil {
		logger.Error(err, "Error in setting controller reference")
		return err
	}
	logger.Info("Creating Suggestion", "namespace", instance.Namespace, "name", instance.Name)
	if err := g.Create(context.TODO(), suggestion); err != nil {
		return err
	}
	return nil
}

func (g *General) UpdateSuggestion(suggestion *suggestionsv1alpha3.Suggestion, suggestionRequests int32) error {
	if err := g.Update(context.TODO(), suggestion); err != nil {
		return err
	}
	return nil
}
