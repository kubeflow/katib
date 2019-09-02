package suggestionclient

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/types"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	commonapiv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/common/v1alpha2"
	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	suggestionsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/suggestions/v1alpha2"
	trialsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
	suggestionapi "github.com/kubeflow/katib/pkg/api/v1alpha2"
	"github.com/kubeflow/katib/pkg/controller/v1alpha2/consts"
)

var log = logf.Log.WithName("suggestion-client")

type SuggestionClient interface {
	SyncAssignments(
		instance *suggestionsv1alpha2.Suggestion,
		e *experimentsv1alpha2.Experiment,
		ts []trialsv1alpha2.Trial) error
}

type General struct {
}

func New() SuggestionClient {
	return &General{}
}

func (g *General) SyncAssignments(
	instance *suggestionsv1alpha2.Suggestion,
	e *experimentsv1alpha2.Experiment,
	ts []trialsv1alpha2.Trial) error {
	logger := log.WithValues("Suggestion", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	requestNum := int(instance.Spec.Suggestions) - len(instance.Status.Assignments)
	if requestNum <= 0 {
		return nil
	}

	endpoint := fmt.Sprintf("%s:%d", instance.Name, consts.DefaultSuggestionPort)
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := suggestionapi.NewSuggestionClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := &suggestionapi.GetSuggestionsRequest{
		Experiment:    g.ConvertExperiment(e),
		Trials:        g.ConvertTrials(ts),
		RequestNumber: int32(requestNum),
	}
	response, err := client.GetSuggestions(ctx, request)
	logger.V(0).Info("Getting suggestions", "endpoint", endpoint, "response", response, "request", request)
	if err != nil {
		return err
	}
	if len(response.Trials) == 0 {
		return fmt.Errorf("The response contains 0 trials")
	}
	return nil
}

// ConvertExperiment converts CRD to the GRPC definition.
func (g *General) ConvertExperiment(e *experimentsv1alpha2.Experiment) *suggestionapi.Experiment {
	res := &suggestionapi.Experiment{}
	res.Name = e.Name
	res.Spec = &suggestionapi.ExperimentSpec{
		Algorithm: &suggestionapi.AlgorithmSpec{
			AlgorithmName:    e.Spec.Algorithm.AlgorithmName,
			AlgorithmSetting: convertAlgorithmSettings(e.Spec.Algorithm.AlgorithmSettings),
		},
		Objective: &suggestionapi.ObjectiveSpec{
			Type:                convertObjectiveType(e.Spec.Objective.Type),
			Goal:                *e.Spec.Objective.Goal,
			ObjectiveMetricName: e.Spec.Objective.ObjectiveMetricName,
		},
		ParameterSpecs: &suggestionapi.ExperimentSpec_ParameterSpecs{
			Parameters: convertParameters(e.Spec.Parameters),
		},
	}
	return res
}

// ConvertTrials converts CRD to the GRPC definition.
func (g *General) ConvertTrials(
	t []trialsv1alpha2.Trial) []*suggestionapi.Trial {
	res := make([]*suggestionapi.Trial, 0)
	return res
}

// ComposeTrialsTemplate composes trials with raw template from the GRPC response.
func (g *General) ComposeTrialsTemplate(ts []*suggestionapi.Trial) []trialsv1alpha2.Trial {
	res := make([]trialsv1alpha2.Trial, 0)
	for _, t := range ts {
		res = append(res, trialsv1alpha2.Trial{
			Spec: trialsv1alpha2.TrialSpec{
				ParameterAssignments: composeParameterAssignments(
					t.Spec.ParameterAssignments.Assignments),
			},
		})
	}
	return res
}

func composeParameterAssignments(pas []*suggestionapi.ParameterAssignment) []commonapiv1alpha2.ParameterAssignment {
	res := make([]commonapiv1alpha2.ParameterAssignment, 0)
	for _, pa := range pas {
		res = append(res, commonapiv1alpha2.ParameterAssignment{
			Name:  pa.Name,
			Value: pa.Value,
		})
	}
	return res
}

func convertObjectiveType(typ commonapiv1alpha2.ObjectiveType) suggestionapi.ObjectiveType {
	switch typ {
	case commonapiv1alpha2.ObjectiveTypeMaximize:
		return suggestionapi.ObjectiveType_MAXIMIZE
	case commonapiv1alpha2.ObjectiveTypeMinimize:
		return suggestionapi.ObjectiveType_MINIMIZE
	default:
		return suggestionapi.ObjectiveType_UNKNOWN
	}
}

func convertAlgorithmSettings(as []commonapiv1alpha2.AlgorithmSetting) []*suggestionapi.AlgorithmSetting {
	res := make([]*suggestionapi.AlgorithmSetting, 0)
	for _, s := range as {
		res = append(res, &suggestionapi.AlgorithmSetting{
			Name:  s.Name,
			Value: s.Value,
		})
	}
	return res
}

func convertParameters(ps []experimentsv1alpha2.ParameterSpec) []*suggestionapi.ParameterSpec {
	res := make([]*suggestionapi.ParameterSpec, 0)
	for _, p := range ps {
		res = append(res, &suggestionapi.ParameterSpec{
			Name:          p.Name,
			ParameterType: convertParameterType(p.ParameterType),
			FeasibleSpace: convertFeasibleSpace(p.FeasibleSpace),
		})
	}
	return res
}

func convertParameterType(typ experimentsv1alpha2.ParameterType) suggestionapi.ParameterType {
	switch typ {
	case experimentsv1alpha2.ParameterTypeDiscrete:
		return suggestionapi.ParameterType_DISCRETE
	case experimentsv1alpha2.ParameterTypeCategorical:
		return suggestionapi.ParameterType_CATEGORICAL
	case experimentsv1alpha2.ParameterTypeDouble:
		return suggestionapi.ParameterType_DOUBLE
	case experimentsv1alpha2.ParameterTypeInt:
		return suggestionapi.ParameterType_INT
	default:
		return suggestionapi.ParameterType_UNKNOWN_TYPE
	}
}

func convertFeasibleSpace(fs experimentsv1alpha2.FeasibleSpace) *suggestionapi.FeasibleSpace {
	res := &suggestionapi.FeasibleSpace{
		Max:  fs.Max,
		Min:  fs.Min,
		List: fs.List,
		Step: fs.Step,
	}
	return res
}
