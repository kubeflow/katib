package suggestionclient

import (
	"context"
	"fmt"
	"time"

	"github.com/kubeflow/katib/pkg/controller.v1alpha3/consts"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
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
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	log        = logf.Log.WithName("suggestion-client")
	timeout    = 60 * time.Second
	timeFormat = "2006-01-02T15:04:05Z"
)

// SuggestionClient is the interface to communicate with algorithm services.
type SuggestionClient interface {
	SyncAssignments(instance *suggestionsv1alpha3.Suggestion, e *experimentsv1alpha3.Experiment,
		ts []trialsv1alpha3.Trial) error

	ValidateAlgorithmSettings(instance *suggestionsv1alpha3.Suggestion, e *experimentsv1alpha3.Experiment) error
}

// General is the implementation for SuggestionClient.
type General struct {
}

// New creates a new SuggestionClient.
func New() SuggestionClient {
	return &General{}
}

// SyncAssignments syncs assignments from algorithm services.
func (g *General) SyncAssignments(
	instance *suggestionsv1alpha3.Suggestion,
	e *experimentsv1alpha3.Experiment,
	ts []trialsv1alpha3.Trial) error {
	logger := log.WithValues("Suggestion", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	requestNum := int(instance.Spec.Requests) - int(instance.Status.SuggestionCount)
	if requestNum <= 0 {
		return nil
	}

	endpoint := util.GetAlgorithmEndpoint(instance)
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := suggestionapi.NewSuggestionClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Algorithm settings in suggestion will overwrite the settings in experiment.
	filledE := e.DeepCopy()
	appendAlgorithmSettingsFromSuggestion(filledE,
		instance.Status.AlgorithmSettings)

	request := &suggestionapi.GetSuggestionsRequest{
		Experiment:    g.ConvertExperiment(filledE),
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
	instance.Status.SuggestionCount = int32(len(instance.Status.Suggestions))

	if response.Algorithm != nil {
		updateAlgorithmSettings(instance, response.Algorithm)
	}
	return nil
}

// ValidateAlgorithmSettings validates if the algorithm specific configurations are valid.
func (g *General) ValidateAlgorithmSettings(instance *suggestionsv1alpha3.Suggestion, e *experimentsv1alpha3.Experiment) error {
	logger := log.WithValues("Suggestion", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	endpoint := util.GetAlgorithmEndpoint(instance)

	callOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(consts.DefaultGRPCRetryPeriod)),
		grpc_retry.WithMax(consts.DefaultGRPCRetryAttempts),
	}
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure(),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(callOpts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(callOpts...)),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := suggestionapi.NewSuggestionClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	request := &suggestionapi.ValidateAlgorithmSettingsRequest{
		Experiment: g.ConvertExperiment(e),
	}

	// See https://github.com/grpc/grpc-go/issues/2636
	// See https://github.com/grpc/grpc-go/pull/2503
	_, err = client.ValidateAlgorithmSettings(ctx, request, grpc.WaitForReady(true))
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
			ObjectiveMetricName: e.Spec.Objective.ObjectiveMetricName,
		},
		ParameterSpecs: &suggestionapi.ExperimentSpec_ParameterSpecs{
			Parameters: convertParameters(e.Spec.Parameters),
		},
	}
	// Set Goal if user defines it in Objective
	if e.Spec.Objective.Goal != nil {
		res.Spec.Objective.Goal = *e.Spec.Objective.Goal
	}
	// Set NasConfig if the user defines it in Spec.
	if e.Spec.NasConfig != nil {
		res.Spec.NasConfig = convertNasConfig(e.Spec.NasConfig)
	}
	if e.Spec.ParallelTrialCount != nil {
		res.Spec.ParallelTrialCount = *e.Spec.ParallelTrialCount
	}
	if e.Spec.MaxTrialCount != nil {
		res.Spec.MaxTrialCount = *e.Spec.MaxTrialCount
	}
	return res
}

// ConvertTrials converts CRD to the GRPC definition.
func (g *General) ConvertTrials(ts []trialsv1alpha3.Trial) []*suggestionapi.Trial {
	trialsRes := make([]*suggestionapi.Trial, 0)
	for _, t := range ts {
		trial := &suggestionapi.Trial{
			Name: t.Name,
			Spec: &suggestionapi.TrialSpec{
				Objective: &suggestionapi.ObjectiveSpec{
					Type:                  convertObjectiveType(t.Spec.Objective.Type),
					ObjectiveMetricName:   t.Spec.Objective.ObjectiveMetricName,
					AdditionalMetricNames: t.Spec.Objective.AdditionalMetricNames,
				},
				ParameterAssignments: convertTrialParameterAssignments(
					t.Spec.ParameterAssignments),
			},
			Status: &suggestionapi.TrialStatus{
				StartTime:      convertTrialStatusTime(t.Status.StartTime),
				CompletionTime: convertTrialStatusTime(t.Status.CompletionTime),
				Observation: convertTrialObservation(
					t.Status.Observation),
			},
		}
		if t.Spec.Objective.Goal != nil {
			trial.Spec.Objective.Goal = *t.Spec.Objective.Goal
		}
		if len(t.Status.Conditions) > 0 {
			// We send only the latest condition of the Trial!
			trial.Status.Condition = convertTrialConditionType(
				t.Status.Conditions[len(t.Status.Conditions)-1].Type)
		}
		trialsRes = append(trialsRes, trial)
	}

	return trialsRes
}

// convertTrialParameterAssignments convert ParameterAssignments CRD to the GRPC definition
func convertTrialParameterAssignments(pas []commonapiv1alpha3.ParameterAssignment) *suggestionapi.TrialSpec_ParameterAssignments {
	tsPas := &suggestionapi.TrialSpec_ParameterAssignments{
		Assignments: make([]*suggestionapi.ParameterAssignment, 0),
	}
	for _, pa := range pas {
		tsPas.Assignments = append(tsPas.Assignments, &suggestionapi.ParameterAssignment{
			Name:  pa.Name,
			Value: pa.Value,
		})
	}

	return tsPas
}

// convertTrialConditionType convert Trial Status Condition Type CRD to the GRPC definition
func convertTrialConditionType(conditionType trialsv1alpha3.TrialConditionType) suggestionapi.TrialStatus_TrialConditionType {
	switch conditionType {
	case trialsv1alpha3.TrialCreated:
		return suggestionapi.TrialStatus_CREATED
	case trialsv1alpha3.TrialRunning:
		return suggestionapi.TrialStatus_RUNNING
	case trialsv1alpha3.TrialSucceeded:
		return suggestionapi.TrialStatus_SUCCEEDED
	case trialsv1alpha3.TrialKilled:
		return suggestionapi.TrialStatus_KILLED
	case trialsv1alpha3.TrialFailed:
		return suggestionapi.TrialStatus_FAILED
	default:
		return suggestionapi.TrialStatus_UNKNOWN
	}
}

// convertTrialObservation convert Trial Observation Metrics CRD to the GRPC definition
func convertTrialObservation(observation *commonapiv1alpha3.Observation) *suggestionapi.Observation {
	resObservation := &suggestionapi.Observation{
		Metrics: make([]*suggestionapi.Metric, 0),
	}
	if observation != nil && observation.Metrics != nil {
		for _, m := range observation.Metrics {
			resObservation.Metrics = append(resObservation.Metrics, &suggestionapi.Metric{
				Name:  m.Name,
				Value: fmt.Sprintf("%f", m.Value),
			})
		}
	}
	return resObservation

}

// convertTrialStatusTime convert Trial Status Time CRD to the GRPC definition
func convertTrialStatusTime(time *metav1.Time) string {
	if time != nil {
		return time.Format(timeFormat)
	}
	return ""
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
