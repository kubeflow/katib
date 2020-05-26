package fake

import (
	"fmt"

	utilrand "k8s.io/apimachinery/pkg/util/rand"

	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/suggestion/suggestionclient"
)

type Fake struct {
}

func New() suggestionclient.SuggestionClient {
	return &Fake{}
}

func (f *Fake) SyncAssignments(
	instance *suggestionsv1beta1.Suggestion,
	e *experimentsv1beta1.Experiment,
	ts []trialsv1beta1.Trial) error {
	if int(instance.Status.SuggestionCount) != int(instance.Spec.Requests) {
		for i := 0; i < int(instance.Spec.Requests)-int(instance.Status.SuggestionCount); i++ {
			name := fmt.Sprintf("%s-%s", instance.Name, utilrand.String(8))
			instance.Status.Suggestions = append(instance.Status.Suggestions, suggestionsv1beta1.TrialAssignment{
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

func (f *Fake) ValidateAlgorithmSettings(instance *suggestionsv1beta1.Suggestion, e *experimentsv1beta1.Experiment) error {
	return nil
}
