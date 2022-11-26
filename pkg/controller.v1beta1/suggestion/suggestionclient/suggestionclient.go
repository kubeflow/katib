/*
Copyright 2022 The Kubeflow Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package suggestionclient

import (
	"context"
	"fmt"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	commonapiv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	suggestionapi "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	katibmanagerv1beta1 "github.com/kubeflow/katib/pkg/common/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/util"
)

var (
	log        = logf.Log.WithName("suggestion-client")
	timeout    = 60 * time.Second
	timeFormat = "2006-01-02T15:04:05Z"

	getRPCClientSuggestion = func(conn *grpc.ClientConn) suggestionapi.SuggestionClient {
		return suggestionapi.NewSuggestionClient(conn)
	}

	getRPCClientEarlyStopping = func(conn *grpc.ClientConn) suggestionapi.EarlyStoppingClient {
		return suggestionapi.NewEarlyStoppingClient(conn)
	}

	callValidatorOpts = []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(consts.DefaultGRPCRetryPeriod)),
		grpc_retry.WithMax(consts.DefaultGRPCRetryAttempts),
	}
)

// SuggestionClient is the interface to communicate with algorithm services.
type SuggestionClient interface {
	SyncAssignments(instance *suggestionsv1beta1.Suggestion, e *experimentsv1beta1.Experiment,
		ts []trialsv1beta1.Trial) error

	ValidateAlgorithmSettings(instance *suggestionsv1beta1.Suggestion, e *experimentsv1beta1.Experiment) error
	ValidateEarlyStoppingSettings(instance *suggestionsv1beta1.Suggestion, e *experimentsv1beta1.Experiment) error
}

// General is the implementation for SuggestionClient.
type General struct {
}

// New creates a new SuggestionClient.
func New() SuggestionClient {
	return &General{}
}

// SyncAssignments syncs assignments from Suggestion and EarlyStopping service.
// If early stopping is set, we call GetEarlyStoppingRules after GetSuggestions
func (g *General) SyncAssignments(
	instance *suggestionsv1beta1.Suggestion,
	e *experimentsv1beta1.Experiment,
	ts []trialsv1beta1.Trial) error {
	logger := log.WithValues("Suggestion", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	currentRequestNum := int(instance.Spec.Requests) - int(instance.Status.SuggestionCount)
	if currentRequestNum <= 0 {
		return nil
	}

	endpoint := util.GetAlgorithmEndpoint(instance)
	connSuggestion, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer connSuggestion.Close()

	// Create client for Suggestion service
	rpcClientSuggestion := getRPCClientSuggestion(connSuggestion)
	ctx, cancelSuggestion := context.WithTimeout(context.Background(), timeout)
	defer cancelSuggestion()

	// Overwrite Experiment algorithm settings from the Suggestion status before the gRPC request.
	// Original algorithm settings are located in Suggestion.Spec.Algorithm.AlgorithmSettings.
	filledE := e.DeepCopy()
	appendAlgorithmSettingsFromSuggestion(filledE,
		instance.Status.AlgorithmSettings)

	requestSuggestion := &suggestionapi.GetSuggestionsRequest{
		Experiment:           g.ConvertExperiment(filledE),
		Trials:               g.ConvertTrials(ts),
		CurrentRequestNumber: int32(currentRequestNum),
		TotalRequestNumber:   int32(instance.Spec.Requests),
	}

	// Get new suggestions
	responseSuggestion, err := rpcClientSuggestion.GetSuggestions(ctx, requestSuggestion)
	if err != nil {
		return err
	}
	logger.Info("Getting suggestions", "endpoint", endpoint, "Number of current request parameters", currentRequestNum, "Number of response parameters", len(responseSuggestion.ParameterAssignments))
	if len(responseSuggestion.ParameterAssignments) != currentRequestNum {
		err := fmt.Errorf("The response contains unexpected trials")
		logger.Error(err, "The response contains unexpected trials")
		return err
	}

	earlyStoppingRules := []commonapiv1beta1.EarlyStoppingRule{}
	// If early stopping is set, call GetEarlyStoppingRules after GetSuggestions.
	if instance.Spec.EarlyStopping != nil {
		endpoint = util.GetEarlyStoppingEndpoint(instance)
		connEarlyStopping, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
		defer connEarlyStopping.Close()

		// Create client for EarlyStopping service
		rpcClientEarlyStopping := getRPCClientEarlyStopping(connEarlyStopping)
		ctx, cancelEarlyStopping := context.WithTimeout(context.Background(), timeout)
		defer cancelEarlyStopping()

		requestEarlyStopping := &suggestionapi.GetEarlyStoppingRulesRequest{
			Experiment:       g.ConvertExperiment(filledE),
			Trials:           g.ConvertTrials(ts),
			DbManagerAddress: katibmanagerv1beta1.GetDBManagerAddr(),
		}

		// Get new early stopping rules
		responseEarlyStopping, err := rpcClientEarlyStopping.GetEarlyStoppingRules(ctx, requestEarlyStopping)
		if err != nil {
			return err
		}

		logger.Info("Getting early stopping rules", "endpoint", endpoint, "response", responseEarlyStopping)

		for _, rule := range responseEarlyStopping.EarlyStoppingRules {
			earlyStoppingRules = append(earlyStoppingRules,
				commonapiv1beta1.EarlyStoppingRule{
					Name:       rule.Name,
					Value:      rule.Value,
					Comparison: convertComparison(rule.Comparison),
					StartStep:  int(rule.StartStep),
				},
			)
		}
	}

	trialAssignments := []suggestionsv1beta1.TrialAssignment{}
	for _, t := range responseSuggestion.ParameterAssignments {
		var trialName string
		if t.TrialName != "" {
			trialName = t.TrialName
		} else {
			trialName = fmt.Sprintf("%s-%s", instance.Name, utilrand.String(8))
		}

		assignment := suggestionsv1beta1.TrialAssignment{
			Name:                 trialName,
			ParameterAssignments: composeParameterAssignments(t.Assignments),
			EarlyStoppingRules:   earlyStoppingRules,
		}
		if t.Labels != nil {
			assignment.Labels = t.Labels
		}
		trialAssignments = append(trialAssignments, assignment)
	}

	instance.Status.Suggestions = append(instance.Status.Suggestions, trialAssignments...)
	instance.Status.SuggestionCount = int32(len(instance.Status.Suggestions))

	if responseSuggestion.Algorithm != nil {
		updateAlgorithmSettings(instance, responseSuggestion.Algorithm)
	}
	return nil
}

// ValidateAlgorithmSettings validates if the algorithm specific configurations are valid.
func (g *General) ValidateAlgorithmSettings(instance *suggestionsv1beta1.Suggestion, e *experimentsv1beta1.Experiment) error {
	logger := log.WithValues("Suggestion", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	endpoint := util.GetAlgorithmEndpoint(instance)

	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(callValidatorOpts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(callValidatorOpts...)),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	rpcClient := getRPCClientSuggestion(conn)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	request := &suggestionapi.ValidateAlgorithmSettingsRequest{
		Experiment: g.ConvertExperiment(e),
	}

	// See https://github.com/grpc/grpc-go/issues/2636
	// See https://github.com/grpc/grpc-go/pull/2503
	_, err = rpcClient.ValidateAlgorithmSettings(ctx, request, grpc.WaitForReady(true))
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
	logger.Info("Algorithm settings are validated")
	return nil
}

// ValidateEarlyStoppingSettings validates if the algorithm specific configurations for early stopping are valid.
func (g *General) ValidateEarlyStoppingSettings(instance *suggestionsv1beta1.Suggestion, e *experimentsv1beta1.Experiment) error {
	logger := log.WithValues("EarlyStopping", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	endpoint := util.GetEarlyStoppingEndpoint(instance)

	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(callValidatorOpts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(callValidatorOpts...)),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	rpcClient := getRPCClientEarlyStopping(conn)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	request := &suggestionapi.ValidateEarlyStoppingSettingsRequest{
		EarlyStopping: g.ConvertExperiment(e).Spec.EarlyStopping,
	}

	// See https://github.com/grpc/grpc-go/issues/2636
	// See https://github.com/grpc/grpc-go/pull/2503
	_, err = rpcClient.ValidateEarlyStoppingSettings(ctx, request, grpc.WaitForReady(true))
	statusCode, _ := status.FromError(err)

	// validation error
	if statusCode.Code() == codes.InvalidArgument || statusCode.Code() == codes.Unknown {
		logger.Error(err, "ValidateEarlyStoppingSettings error")
		return fmt.Errorf("ValidateEarlyStoppingSettings Error: %v", statusCode.Message())
	}

	// Connection error
	if statusCode.Code() == codes.Unavailable {
		logger.Error(err, "Connection to EarlyStopping algorithm service currently unavailable")
		return err
	}

	// Validate to true as function is not implemented
	if statusCode.Code() == codes.Unimplemented {
		logger.Info("Method ValidateEarlyStoppingSettings not found", "EarlyStopping service", e.Spec.EarlyStopping.AlgorithmName)
		return nil
	}
	logger.Info("EarlyStopping settings are validated")
	return nil
}

// ConvertExperiment converts CRD to the GRPC definition.
func (g *General) ConvertExperiment(e *experimentsv1beta1.Experiment) *suggestionapi.Experiment {
	res := &suggestionapi.Experiment{}
	res.Name = e.Name
	res.Spec = &suggestionapi.ExperimentSpec{
		Algorithm: &suggestionapi.AlgorithmSpec{
			AlgorithmName:     e.Spec.Algorithm.AlgorithmName,
			AlgorithmSettings: convertAlgorithmSettings(e.Spec.Algorithm.AlgorithmSettings),
		},
		Objective: &suggestionapi.ObjectiveSpec{
			Type:                  convertObjectiveType(e.Spec.Objective.Type),
			ObjectiveMetricName:   e.Spec.Objective.ObjectiveMetricName,
			AdditionalMetricNames: e.Spec.Objective.AdditionalMetricNames,
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
	// Set early stopping if it is needed
	if e.Spec.EarlyStopping != nil {
		res.Spec.EarlyStopping = &suggestionapi.EarlyStoppingSpec{
			AlgorithmName:     e.Spec.EarlyStopping.AlgorithmName,
			AlgorithmSettings: convertEarlyStoppingSettings(e.Spec.EarlyStopping.AlgorithmSettings),
		}
	}
	return res
}

// ConvertTrials converts CRD to the GRPC definition.
func (g *General) ConvertTrials(ts []trialsv1beta1.Trial) []*suggestionapi.Trial {
	trialsRes := make([]*suggestionapi.Trial, 0)
	for _, t := range ts {
		// Skip Trials with unavailable metrics.
		if t.IsMetricsUnavailable() {
			continue
		}
		if !t.IsObservationAvailable() && t.IsEarlyStopped() {
			continue
		}
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
				Observation:    convertTrialObservation(t.Spec.Objective.MetricStrategies, t.Status.Observation),
			},
		}
		if t.Spec.Labels != nil {
			trial.Spec.Labels = t.Spec.Labels
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
func convertTrialParameterAssignments(pas []commonapiv1beta1.ParameterAssignment) *suggestionapi.TrialSpec_ParameterAssignments {
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
func convertTrialConditionType(conditionType trialsv1beta1.TrialConditionType) suggestionapi.TrialStatus_TrialConditionType {
	switch conditionType {
	case trialsv1beta1.TrialCreated:
		return suggestionapi.TrialStatus_CREATED
	case trialsv1beta1.TrialRunning:
		return suggestionapi.TrialStatus_RUNNING
	case trialsv1beta1.TrialSucceeded:
		return suggestionapi.TrialStatus_SUCCEEDED
	case trialsv1beta1.TrialKilled:
		return suggestionapi.TrialStatus_KILLED
	case trialsv1beta1.TrialFailed:
		return suggestionapi.TrialStatus_FAILED
	case trialsv1beta1.TrialEarlyStopped:
		return suggestionapi.TrialStatus_EARLYSTOPPED
	default:
		return suggestionapi.TrialStatus_UNKNOWN
	}
}

// convertTrialObservation convert Trial Observation Metrics CRD to the GRPC definition
func convertTrialObservation(strategies []commonapiv1beta1.MetricStrategy, observation *commonapiv1beta1.Observation) *suggestionapi.Observation {
	resObservation := &suggestionapi.Observation{
		Metrics: make([]*suggestionapi.Metric, 0),
	}
	strategyMap := make(map[string]commonapiv1beta1.MetricStrategyType)
	for _, strategy := range strategies {
		strategyMap[strategy.Name] = strategy.Value
	}
	if observation != nil && observation.Metrics != nil {
		for _, m := range observation.Metrics {
			var value string
			switch strategy := strategyMap[m.Name]; strategy {
			case commonapiv1beta1.ExtractByMin:
				if m.Min == consts.UnavailableMetricValue {
					value = m.Latest
				} else {
					value = m.Min
				}
			case commonapiv1beta1.ExtractByMax:
				if m.Max == consts.UnavailableMetricValue {
					value = m.Latest
				} else {
					value = m.Max
				}
			case commonapiv1beta1.ExtractByLatest:
				value = m.Latest
			}
			resObservation.Metrics = append(resObservation.Metrics, &suggestionapi.Metric{
				Name:  m.Name,
				Value: value,
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

func composeParameterAssignments(pas []*suggestionapi.ParameterAssignment) []commonapiv1beta1.ParameterAssignment {
	res := make([]commonapiv1beta1.ParameterAssignment, 0)
	for _, pa := range pas {
		res = append(res, commonapiv1beta1.ParameterAssignment{
			Name:  pa.Name,
			Value: pa.Value,
		})
	}
	return res
}

func convertObjectiveType(typ commonapiv1beta1.ObjectiveType) suggestionapi.ObjectiveType {
	switch typ {
	case commonapiv1beta1.ObjectiveTypeMaximize:
		return suggestionapi.ObjectiveType_MAXIMIZE
	case commonapiv1beta1.ObjectiveTypeMinimize:
		return suggestionapi.ObjectiveType_MINIMIZE
	default:
		return suggestionapi.ObjectiveType_UNKNOWN
	}
}

func convertAlgorithmSettings(as []commonapiv1beta1.AlgorithmSetting) []*suggestionapi.AlgorithmSetting {
	res := make([]*suggestionapi.AlgorithmSetting, 0)
	for _, s := range as {
		res = append(res, &suggestionapi.AlgorithmSetting{
			Name:  s.Name,
			Value: s.Value,
		})
	}
	return res
}

func convertEarlyStoppingSettings(es []commonapiv1beta1.EarlyStoppingSetting) []*suggestionapi.EarlyStoppingSetting {
	res := make([]*suggestionapi.EarlyStoppingSetting, 0)
	for _, s := range es {
		res = append(res, &suggestionapi.EarlyStoppingSetting{
			Name:  s.Name,
			Value: s.Value,
		})
	}
	return res
}

func convertParameters(ps []experimentsv1beta1.ParameterSpec) []*suggestionapi.ParameterSpec {
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

func convertParameterType(typ experimentsv1beta1.ParameterType) suggestionapi.ParameterType {
	switch typ {
	case experimentsv1beta1.ParameterTypeDiscrete:
		return suggestionapi.ParameterType_DISCRETE
	case experimentsv1beta1.ParameterTypeCategorical:
		return suggestionapi.ParameterType_CATEGORICAL
	case experimentsv1beta1.ParameterTypeDouble:
		return suggestionapi.ParameterType_DOUBLE
	case experimentsv1beta1.ParameterTypeInt:
		return suggestionapi.ParameterType_INT
	default:
		return suggestionapi.ParameterType_UNKNOWN_TYPE
	}
}

func convertFeasibleSpace(fs experimentsv1beta1.FeasibleSpace) *suggestionapi.FeasibleSpace {
	res := &suggestionapi.FeasibleSpace{
		Max:  fs.Max,
		Min:  fs.Min,
		List: fs.List,
		Step: fs.Step,
	}
	return res
}

func convertComparison(comparison suggestionapi.ComparisonType) commonapiv1beta1.ComparisonType {
	switch comparison {
	case suggestionapi.ComparisonType_EQUAL:
		return commonapiv1beta1.ComparisonTypeEqual
	case suggestionapi.ComparisonType_LESS:
		return commonapiv1beta1.ComparisonTypeLess
	case suggestionapi.ComparisonType_GREATER:
		return commonapiv1beta1.ComparisonTypeGreater
	default:
		return commonapiv1beta1.ComparisonTypeEqual
	}
}
