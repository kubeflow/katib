package suggestionclient

import (
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	suggestionsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1alpha3"
	trialsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1alpha3"
)

var log = logf.Log.WithName("suggestion-client")

type SuggestionClient interface {
	SyncAssignments(
		instance *suggestionsv1alpha3.Suggestion,
		e *experimentsv1alpha3.Experiment,
		ts []trialsv1alpha3.Trial) error
}
