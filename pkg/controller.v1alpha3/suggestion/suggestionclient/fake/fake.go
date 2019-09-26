package fake

import (
	"fmt"

	utilrand "k8s.io/apimachinery/pkg/util/rand"

	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1alpha3"
	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	suggestionsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1alpha3"
	trialsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/suggestion/suggestionclient"
)

type Fake struct {
}

func New() suggestionclient.SuggestionClient {
	return &Fake{}
}

func (f *Fake) SyncAssignments(
	instance *suggestionsv1alpha3.Suggestion,
	e *experimentsv1alpha3.Experiment,
	ts []trialsv1alpha3.Trial) error {
	if int(instance.Status.SuggestionCount) != int(instance.Spec.Requests) {
		for i := 0; i < int(instance.Spec.Requests)-int(instance.Status.SuggestionCount); i++ {
			name := fmt.Sprintf("%s-%s", instance.Name, utilrand.String(8))
			instance.Status.Suggestions = append(instance.Status.Suggestions, suggestionsv1alpha3.TrialAssignment{
				Name: name,
				ParameterAssignments: []common.ParameterAssignment{
					{
						Name:  "--lr",
						Value: "0.03",
					},
					{
						Name:  "--num-layers",
						Value: "4",
					},
					{
						Name:  "--optimizer",
						Value: "adam",
					},
				},
			})
		}
	}
	return nil
}

func (f *Fake) ValidateAlgorithmSettings(instance *suggestionsv1alpha3.Suggestion, e *experimentsv1alpha3.Experiment) error {
	return nil
}
