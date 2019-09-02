package fake

import (
	common "github.com/kubeflow/katib/pkg/api/operators/apis/common/v1alpha2"
	suggestionsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/suggestions/v1alpha2"
	"github.com/kubeflow/katib/pkg/controller/v1alpha2/suggestion/suggestionclient"
)

type Fake struct {
}

func New() suggestionclient.SuggestionClient {
	return &Fake{}
}

func (f *Fake) SyncAssignments(instance *suggestionsv1alpha2.Suggestion) error {
	if len(instance.Status.Assignments) != int(instance.Spec.Suggestions) {
		for i := 0; i < int(instance.Spec.Suggestions)-len(instance.Status.Assignments); i++ {
			instance.Status.Assignments = append(instance.Status.Assignments, suggestionsv1alpha2.TrialAssignment{
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
