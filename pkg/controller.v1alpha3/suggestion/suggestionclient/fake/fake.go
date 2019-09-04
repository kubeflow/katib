package fake

import (
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
	if len(instance.Status.Assignments) != int(*instance.Spec.Suggestions) {
		for i := 0; i < int(*instance.Spec.Suggestions)-len(instance.Status.Assignments); i++ {
			instance.Status.Assignments = append(instance.Status.Assignments, suggestionsv1alpha3.TrialAssignment{
				Assignments: []common.ParameterAssignment{
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
