package suggestion

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	suggestionsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/suggestions/v1alpha2"
)

var log = logf.Log.WithName("experiment-suggestionc-client")

// Suggestion is the interface for suggestions in Experiment controller.
type Suggestion interface {
	CreateSuggestion(instance *experimentsv1alpha2.Experiment) error
	RequestSuggestions(s *suggestionsv1alpha2.Suggestion, addcount int32) error
	UpdateSuggestion(s *suggestionsv1alpha2.Suggestion, os *suggestionsv1alpha2.Suggestion) error
	GetSuggestions(s *suggestionsv1alpha2.Suggestion) ([]suggestionsv1alpha2.TrialAssignment, error)
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

func (g *General) RequestSuggestions(s *suggestionsv1alpha2.Suggestion, addcount int32) error {
	if int(s.Spec.Suggestions) > len(s.Status.Assignments) {
		return nil
	}
	s.Spec.Suggestions += addcount
	if err := g.Update(context.TODO(), s); err != nil {
		return err
	}
	return nil
}

func (g *General) CreateSuggestion(instance *experimentsv1alpha2.Experiment) error {
	logger := log.WithValues("experiment", types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace})
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

	found := &suggestionsv1alpha2.Suggestion{}
	err := g.Get(context.TODO(), types.NamespacedName{Name: suggestion.Name, Namespace: suggestion.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating Suggestion", "namespace", suggestion.Namespace, "name", suggestion.Name)
		err = g.Create(context.TODO(), suggestion)
		return err
	} else if err != nil {
		return err
	}
	return nil
}

func (g *General) UpdateSuggestion(s *suggestionsv1alpha2.Suggestion, os *suggestionsv1alpha2.Suggestion) error {
	if !equality.Semantic.DeepEqual(s.Status.Assignments, os.Status.Assignments) {
		if err := g.Status().Update(context.TODO(), s); err != nil {
			return err
		}
	}
	return nil
}

func (g *General) GetSuggestions(
	s *suggestionsv1alpha2.Suggestion) ([]suggestionsv1alpha2.TrialAssignment, error) {
	tas := make([]suggestionsv1alpha2.TrialAssignment, 0)
	for i := range s.Status.Assignments {
		if s.Status.Assignments[i].Name == nil {
			name := fmt.Sprintf("%s-%s", s.Name, utilrand.String(8))
			// Set the name for suggestions.
			s.Status.Assignments[i].Name = &name
			tas = append(tas, s.Status.Assignments[i])
		}
	}
	return tas, nil
}
