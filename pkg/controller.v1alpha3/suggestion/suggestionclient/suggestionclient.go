package suggestionclient

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/types"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	commonapiv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/common/v1alpha3"
	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	suggestionsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1alpha3"
	trialsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1alpha3"
	suggestionapi "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/consts"
)

var log = logf.Log.WithName("suggestion-client")

type SuggestionClient interface {
	SyncAssignments(instance *suggestionsv1alpha3.Suggestion, e *experimentsv1alpha3.Experiment,
		ts []trialsv1alpha3.Trial) error

	ValidateAlgorithmSettings(instance *suggestionsv1alpha3.Suggestion, e *experimentsv1alpha3.Experiment) error
}

type General struct {
}

func New() SuggestionClient {
	return &General{}
}

func (g *General) SyncAssignments(
	instance *suggestionsv1alpha3.Suggestion,
	e *experimentsv1alpha3.Experiment,
	ts []trialsv1alpha3.Trial) error {
	logger := log.WithValues("Suggestion", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	requestNum := int(instance.Spec.Requests) - len(instance.Status.Suggestions)
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
	if err != nil {
		return err
	}
	logger.V(0).Info("Getting suggestions", "endpoint", endpoint, "response", response, "request", request)
	if len(response.ParameterAssignments) != requestNum {
		err := fmt.Errorf("The response contains unexpected trials")
		logger.Error(err, "The response contains unexpected trials", "requestNum", requestNum, "response", response)
		return err
	}
	for _, t := range response.ParameterAssignments {
		instance.Status.Suggestions = append(instance.Status.Suggestions,
			suggestionsv1alpha3.TrialAssignment{
				Name:                 fmt.Sprintf("%s-%s", instance.Name, utilrand.String(8)),
				ParameterAssignments: composeParameterAssignments(t.Assignments),
			})
	}

	// TODO(gaocegege): Set algorithm settings
	return nil
}

func (g *General) ValidateAlgorithmSettings(instance *suggestionsv1alpha3.Suggestion, e *experimentsv1alpha3.Experiment) error {
	logger := log.WithValues("Suggestion", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	endpoint := fmt.Sprintf("%s:%d", instance.Name, consts.DefaultSuggestionPort)
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := suggestionapi.NewSuggestionClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := &suggestionapi.ValidateAlgorithmSettingsRequest{
		Experiment: g.ConvertExperiment(e),
	}
	_, err = client.ValidateAlgorithmSettings(ctx, request)
	statusCode, _ := status.FromError(err)

	// validation error
	if statusCode.Code() == codes.InvalidArgument || statusCode.Code() == codes.Unknown {
		logger.Error(err, "ValidateAlgorithmSettings error")
		return fmt.Errorf("ValidateAlgorithmSettings Error: %v", statusCode.Message())
	}

	// Connection error
	if statusCode.Code() == codes.Unavailable {
		logger.Error(err, "Connection to Suggestion algorithm service currently unavailable")
		return err
	}

	// Validate to true as function is not implemented
	if statusCode.Code() == codes.Unimplemented {
		logger.Info("Method ValidateAlgorithmSettings not found", "Suggestion service", e.Spec.Algorithm.AlgorithmName)
		return nil
	}
	logger.Info("Algorithm settings validated")
	return nil
}

// ConvertExperiment converts CRD to the GRPC definition.
func (g *General) ConvertExperiment(e *experimentsv1alpha3.Experiment) *suggestionapi.Experiment {
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
	t []trialsv1alpha3.Trial) []*suggestionapi.Trial {
	res := make([]*suggestionapi.Trial, 0)
	return res
}

// ComposeTrialsTemplate composes trials with raw template from the GRPC response.
func (g *General) ComposeTrialsTemplate(ts []*suggestionapi.Trial) []trialsv1alpha3.Trial {
	res := make([]trialsv1alpha3.Trial, 0)
	for _, t := range ts {
		res = append(res, trialsv1alpha3.Trial{
			Spec: trialsv1alpha3.TrialSpec{
				ParameterAssignments: composeParameterAssignments(
					t.Spec.ParameterAssignments.Assignments),
			},
		})
	}
	return res
}

func composeParameterAssignments(pas []*suggestionapi.ParameterAssignment) []commonapiv1alpha3.ParameterAssignment {
	res := make([]commonapiv1alpha3.ParameterAssignment, 0)
	for _, pa := range pas {
		res = append(res, commonapiv1alpha3.ParameterAssignment{
			Name:  pa.Name,
			Value: pa.Value,
		})
	}
	return res
}

func convertObjectiveType(typ commonapiv1alpha3.ObjectiveType) suggestionapi.ObjectiveType {
	switch typ {
	case commonapiv1alpha3.ObjectiveTypeMaximize:
		return suggestionapi.ObjectiveType_MAXIMIZE
	case commonapiv1alpha3.ObjectiveTypeMinimize:
		return suggestionapi.ObjectiveType_MINIMIZE
	default:
		return suggestionapi.ObjectiveType_UNKNOWN
	}
}

func convertAlgorithmSettings(as []commonapiv1alpha3.AlgorithmSetting) []*suggestionapi.AlgorithmSetting {
	res := make([]*suggestionapi.AlgorithmSetting, 0)
	for _, s := range as {
		res = append(res, &suggestionapi.AlgorithmSetting{
			Name:  s.Name,
			Value: s.Value,
		})
	}
	return res
}

func convertParameters(ps []experimentsv1alpha3.ParameterSpec) []*suggestionapi.ParameterSpec {
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

func convertParameterType(typ experimentsv1alpha3.ParameterType) suggestionapi.ParameterType {
	switch typ {
	case experimentsv1alpha3.ParameterTypeDiscrete:
		return suggestionapi.ParameterType_DISCRETE
	case experimentsv1alpha3.ParameterTypeCategorical:
		return suggestionapi.ParameterType_CATEGORICAL
	case experimentsv1alpha3.ParameterTypeDouble:
		return suggestionapi.ParameterType_DOUBLE
	case experimentsv1alpha3.ParameterTypeInt:
		return suggestionapi.ParameterType_INT
	default:
		return suggestionapi.ParameterType_UNKNOWN_TYPE
	}
}

func convertFeasibleSpace(fs experimentsv1alpha3.FeasibleSpace) *suggestionapi.FeasibleSpace {
	res := &suggestionapi.FeasibleSpace{
		Max:  fs.Max,
		Min:  fs.Min,
		List: fs.List,
		Step: fs.Step,
	}
	return res
}
