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
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	commonapiv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	suggestionapi "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	suggestionapimock "github.com/kubeflow/katib/pkg/mock/v1beta1/api"
)

const (
	algorithmName              = "algorithm-name"
	earlyStoppingAlgorithmName = "early-stopping-name"
)

type k8sMatcher struct {
	x interface{}
}

func (k8s k8sMatcher) Matches(x interface{}) bool {
	switch ex := k8s.x.(type) {
	case proto.Message:
		return proto.Equal(ex, x.(proto.Message))
	default:
		return equality.Semantic.DeepEqual(k8s.x, x)
	}
}

func (k8s k8sMatcher) String() string {
	return fmt.Sprintf("is equal to %v", k8s.x)
}

func TestGetRPCClientSuggestion(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	fakeConn := &grpc.ClientConn{}
	actualClient := getRPCClientSuggestion(fakeConn)
	g.Expect(actualClient).To(gomega.Equal(suggestionapi.NewSuggestionClient(fakeConn)))
}

func TestGetRPCClientEarlyStopping(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	fakeConn := &grpc.ClientConn{}
	actualClient := getRPCClientEarlyStopping(fakeConn)
	g.Expect(actualClient).To(gomega.Equal(suggestionapi.NewEarlyStoppingClient(fakeConn)))
}

func TestSyncAssignments(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	rpcClientSuggestion := suggestionapimock.NewMockSuggestionClient(mockCtrl)
	rpcClientEarlyStopping := suggestionapimock.NewMockEarlyStoppingClient(mockCtrl)

	getRPCClientSuggestion = func(conn *grpc.ClientConn) suggestionapi.SuggestionClient {
		return rpcClientSuggestion
	}
	getRPCClientEarlyStopping = func(conn *grpc.ClientConn) suggestionapi.EarlyStoppingClient {
		return rpcClientEarlyStopping
	}

	suggestionClient := New()

	expectedRequestSuggestion := newFakeRequest()
	expectedRequestEarlyStopping := &suggestionapi.GetEarlyStoppingRulesRequest{
		Experiment:       newFakeRequest().Experiment,
		Trials:           newFakeRequest().Trials,
		DbManagerAddress: fmt.Sprintf("katib-db-manager.kubeflow:%v", consts.DefaultSuggestionPort),
	}

	getSuggestionReply := &suggestionapi.GetSuggestionsReply{
		ParameterAssignments: []*suggestionapi.GetSuggestionsReply_ParameterAssignments{
			{
				Assignments: []*suggestionapi.ParameterAssignment{
					{
						Name:  "param1-name",
						Value: "1",
					},
					{
						Name:  "param2-name",
						Value: "0.3",
					},
				},
			},
			{
				Assignments: []*suggestionapi.ParameterAssignment{
					{
						Name:  "param1-name",
						Value: "2",
					},
					{
						Name:  "param2-name",
						Value: "0.4",
					},
				},
			},
		},
		Algorithm: &suggestionapi.AlgorithmSpec{
			AlgorithmSettings: []*suggestionapi.AlgorithmSetting{
				{
					Name:  "overridden-name",
					Value: "suggestion-value",
				},
				{
					Name:  "new-suggestion-setting-name",
					Value: "value",
				},
			},
		},
	}

	getEarlyStoppingRulesReply := &suggestionapi.GetEarlyStoppingRulesReply{
		EarlyStoppingRules: []*suggestionapi.EarlyStoppingRule{
			{
				Name:       "accuracy",
				Value:      "0.7",
				Comparison: suggestionapi.ComparisonType_LESS,
				StartStep:  4,
			},
			{
				Name:       "epoch",
				Value:      "10",
				Comparison: suggestionapi.ComparisonType_EQUAL,
			},
		},
	}

	validRunGetSuggestions := rpcClientSuggestion.EXPECT().GetSuggestions(gomock.Any(), k8sMatcher{expectedRequestSuggestion}).Return(getSuggestionReply, nil)
	validRunGetEarlyStopRules := rpcClientEarlyStopping.EXPECT().GetEarlyStoppingRules(gomock.Any(), k8sMatcher{expectedRequestEarlyStopping}).Return(getEarlyStoppingRulesReply, nil)
	getSuggestionsFail := rpcClientSuggestion.EXPECT().GetSuggestions(gomock.Any(), gomock.Any()).Return(nil, errors.New("Suggestion service connection error"))

	invalidAssignmentsCount := rpcClientSuggestion.EXPECT().GetSuggestions(gomock.Any(), gomock.Any()).Return(
		&suggestionapi.GetSuggestionsReply{
			ParameterAssignments: []*suggestionapi.GetSuggestionsReply_ParameterAssignments{
				{
					Assignments: []*suggestionapi.ParameterAssignment{
						{
							Name:  "param1-name",
							Value: "1",
						},
					},
				},
			},
		}, nil)

	validRunGetSuggestions2 := rpcClientSuggestion.EXPECT().GetSuggestions(gomock.Any(), k8sMatcher{expectedRequestSuggestion}).Return(getSuggestionReply, nil)
	getEarlyStopRulesFail := rpcClientEarlyStopping.EXPECT().GetEarlyStoppingRules(gomock.Any(), gomock.Any()).Return(nil, errors.New("Suggestion service connection error"))

	gomock.InOrder(
		validRunGetSuggestions,
		validRunGetEarlyStopRules,
		getSuggestionsFail,
		invalidAssignmentsCount,
		validRunGetSuggestions2,
		getEarlyStopRulesFail,
	)

	tcs := []struct {
		experiment      *experimentsv1beta1.Experiment
		suggestion      *suggestionsv1beta1.Suggestion
		trials          []trialsv1beta1.Trial
		err             bool
		testDescription string
	}{
		// Experiment contains HP and NAS config just for the test purpose
		// validRunGetSuggestions + validRunGetEarlyStopRules case
		{
			experiment:      newFakeExperiment(),
			suggestion:      newFakeSuggestion(),
			trials:          newFakeTrials(),
			err:             false,
			testDescription: "SyncAssignments valid run",
		},
		{
			suggestion: func() *suggestionsv1beta1.Suggestion {
				s := newFakeSuggestion()
				s.Spec.Requests = 4
				s.Status.SuggestionCount = 6
				return s
			}(),
			err:             false,
			testDescription: "Negative request number",
		},
		// getSuggestionsFail case
		{
			experiment:      newFakeExperiment(),
			suggestion:      newFakeSuggestion(),
			trials:          newFakeTrials(),
			err:             true,
			testDescription: "Unable to execute GetSuggestions",
		},
		// invalidAssignmentsCount case
		{
			experiment:      newFakeExperiment(),
			suggestion:      newFakeSuggestion(),
			trials:          newFakeTrials(),
			err:             true,
			testDescription: "ParameterAssignments from response != request number",
		},
		// validRunGetSuggestions2 + getEarlyStopRulesFail case
		{
			experiment:      newFakeExperiment(),
			suggestion:      newFakeSuggestion(),
			trials:          newFakeTrials(),
			err:             true,
			testDescription: "Unable to execute GetEarlyStoppingRules",
		},
	}
	for _, tc := range tcs {
		err := suggestionClient.SyncAssignments(tc.suggestion, tc.experiment, tc.trials)
		if !tc.err && err != nil {
			t.Errorf("Case: %v failed. Expected nil, got %v", tc.testDescription, err)
		} else if tc.err && err == nil {
			t.Errorf("Case: %v failed. Expected err, got nil", tc.testDescription)
		}
	}
}

func TestValidateAlgorithmSettings(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	rpcClientSuggestion := suggestionapimock.NewMockSuggestionClient(mockCtrl)

	getRPCClientSuggestion = func(conn *grpc.ClientConn) suggestionapi.SuggestionClient {
		return rpcClientSuggestion
	}

	expectedRequest := &suggestionapi.ValidateAlgorithmSettingsRequest{
		Experiment: newFakeRequest().Experiment,
	}
	expectedRequest.Experiment.Spec.Algorithm.AlgorithmSettings = []*suggestionapi.AlgorithmSetting{
		{
			Name:  "overridden-name",
			Value: "value",
		},
	}

	validRun := rpcClientSuggestion.EXPECT().ValidateAlgorithmSettings(gomock.Any(), k8sMatcher{expectedRequest}, gomock.Any()).Return(nil, nil)

	invalidExperiment := rpcClientSuggestion.EXPECT().ValidateAlgorithmSettings(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil,
		status.Error(codes.InvalidArgument, "Invalid experiment parameter"))
	connectionError := rpcClientSuggestion.EXPECT().ValidateAlgorithmSettings(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil,
		status.Error(codes.Unavailable, "Unable to connect"))
	unimplementedMethod := rpcClientSuggestion.EXPECT().ValidateAlgorithmSettings(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil,
		status.Error(codes.Unimplemented, "Method not implemented"))

	suggestionClient := New()

	exp := newFakeExperiment()
	sug := newFakeSuggestion()

	gomock.InOrder(
		validRun,
		invalidExperiment,
		connectionError,
		unimplementedMethod)

	tcs := []struct {
		experiment      *experimentsv1beta1.Experiment
		suggestion      *suggestionsv1beta1.Suggestion
		err             bool
		testDescription string
	}{
		// validRun case
		{
			experiment:      exp,
			suggestion:      sug,
			err:             false,
			testDescription: "ValidateAlgorithmSettings valid run",
		},
		// invalidExperiment case
		{
			experiment:      exp,
			suggestion:      sug,
			err:             true,
			testDescription: "Invalid argument return in Experiment validation",
		},
		// connectionError case
		{
			experiment:      exp,
			suggestion:      sug,
			err:             true,
			testDescription: "Connection to suggestion service error",
		},
		// unimplementedMethod case
		{
			experiment:      exp,
			suggestion:      sug,
			err:             false,
			testDescription: "Unimplemented ValidateAlgorithmSettings method",
		},
	}
	for _, tc := range tcs {
		err := suggestionClient.ValidateAlgorithmSettings(tc.suggestion, tc.experiment)
		if !tc.err && err != nil {
			t.Errorf("Case: %v failed. Expected nil, got %v", tc.testDescription, err)
		} else if tc.err && err == nil {
			t.Errorf("Case: %v failed. Expected err, got nil", tc.testDescription)
		}
	}
}

func TestValidateEarlyStoppingSettings(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	rpcClientEarlyStopping := suggestionapimock.NewMockEarlyStoppingClient(mockCtrl)

	getRPCClientEarlyStopping = func(conn *grpc.ClientConn) suggestionapi.EarlyStoppingClient {
		return rpcClientEarlyStopping
	}

	expectedRequest := &suggestionapi.ValidateEarlyStoppingSettingsRequest{
		EarlyStopping: newFakeRequest().Experiment.Spec.EarlyStopping,
	}

	validRun := rpcClientEarlyStopping.EXPECT().ValidateEarlyStoppingSettings(gomock.Any(), k8sMatcher{expectedRequest}, gomock.Any()).Return(nil, nil)

	invalidExperiment := rpcClientEarlyStopping.EXPECT().ValidateEarlyStoppingSettings(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil,
		status.Error(codes.InvalidArgument, "Invalid experiment parameter"))
	connectionError := rpcClientEarlyStopping.EXPECT().ValidateEarlyStoppingSettings(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil,
		status.Error(codes.Unavailable, "Unable to connect"))
	unimplementedMethod := rpcClientEarlyStopping.EXPECT().ValidateEarlyStoppingSettings(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil,
		status.Error(codes.Unimplemented, "Method not implemented"))

	suggestionClient := New()

	exp := newFakeExperiment()
	sug := newFakeSuggestion()

	gomock.InOrder(
		validRun,
		invalidExperiment,
		connectionError,
		unimplementedMethod)

	tcs := []struct {
		experiment      *experimentsv1beta1.Experiment
		suggestion      *suggestionsv1beta1.Suggestion
		err             bool
		testDescription string
	}{
		// validRun case
		{
			experiment:      exp,
			suggestion:      sug,
			err:             false,
			testDescription: "ValidateEarlyStoppingSettings valid run",
		},
		// invalidExperiment case
		{
			experiment:      exp,
			suggestion:      sug,
			err:             true,
			testDescription: "Invalid argument return in Experiment validation",
		},
		// connectionError case
		{
			experiment:      exp,
			suggestion:      sug,
			err:             true,
			testDescription: "Connection to early stopping service error",
		},
		// unimplementedMethod case
		{
			experiment:      exp,
			suggestion:      sug,
			err:             false,
			testDescription: "Unimplemented ValidateEarlyStoppingSettings method",
		},
	}
	for _, tc := range tcs {
		err := suggestionClient.ValidateEarlyStoppingSettings(tc.suggestion, tc.experiment)
		if !tc.err && err != nil {
			t.Errorf("Case: %v failed. Expected nil, got %v", tc.testDescription, err)
		} else if tc.err && err == nil {
			t.Errorf("Case: %v failed. Expected err, got nil", tc.testDescription)
		}
	}
}

func TestConvertTrialConditionType(t *testing.T) {

	tcs := []struct {
		inCondition       trialsv1beta1.TrialConditionType
		expectedCondition suggestionapi.TrialStatus_TrialConditionType
		testDescription   string
	}{
		{
			inCondition:       trialsv1beta1.TrialCreated,
			expectedCondition: suggestionapi.TrialStatus_CREATED,
			testDescription:   "Convert created Trial condition",
		},
		{
			inCondition:       trialsv1beta1.TrialRunning,
			expectedCondition: suggestionapi.TrialStatus_RUNNING,
			testDescription:   "Convert running Trial condition",
		},
		{
			inCondition:       trialsv1beta1.TrialSucceeded,
			expectedCondition: suggestionapi.TrialStatus_SUCCEEDED,
			testDescription:   "Convert succeeded Trial condition",
		},
		{
			inCondition:       trialsv1beta1.TrialKilled,
			expectedCondition: suggestionapi.TrialStatus_KILLED,
			testDescription:   "Convert killed Trial condition",
		},
		{
			inCondition:       trialsv1beta1.TrialFailed,
			expectedCondition: suggestionapi.TrialStatus_FAILED,
			testDescription:   "Convert failed Trial condition",
		},
		{
			inCondition:       trialsv1beta1.TrialEarlyStopped,
			expectedCondition: suggestionapi.TrialStatus_EARLYSTOPPED,
			testDescription:   "Convert early stopped Trial condition",
		},
		{
			inCondition:       "Unknown",
			expectedCondition: suggestionapi.TrialStatus_UNKNOWN,
			testDescription:   "Convert unknown Trial condition",
		},
	}
	for _, tc := range tcs {
		actualCondition := convertTrialConditionType(tc.inCondition)
		if actualCondition != tc.expectedCondition {
			t.Errorf("Case: %v failed. Expected Trial condition %v, got %v", tc.testDescription, tc.expectedCondition, actualCondition)
		}
	}
}

func TestConvertObjectiveType(t *testing.T) {

	tcs := []struct {
		inType          commonapiv1beta1.ObjectiveType
		expectedType    suggestionapi.ObjectiveType
		testDescription string
	}{
		{
			inType:          commonv1beta1.ObjectiveTypeMaximize,
			expectedType:    suggestionapi.ObjectiveType_MAXIMIZE,
			testDescription: "Convert maximize objective type",
		},

		{
			inType:          commonv1beta1.ObjectiveTypeMinimize,
			expectedType:    suggestionapi.ObjectiveType_MINIMIZE,
			testDescription: "Convert minimize objective type",
		},
		{
			inType:          commonv1beta1.ObjectiveTypeUnknown,
			expectedType:    suggestionapi.ObjectiveType_UNKNOWN,
			testDescription: "Convert unknown objective type",
		},
	}
	for _, tc := range tcs {
		actualType := convertObjectiveType(tc.inType)
		if actualType != tc.expectedType {
			t.Errorf("Case: %v failed. Expected objective type %v, got %v", tc.testDescription, tc.expectedType, actualType)
		}
	}
}

func TestConvertParameterType(t *testing.T) {

	tcs := []struct {
		inType          experimentsv1beta1.ParameterType
		expectedType    suggestionapi.ParameterType
		testDescription string
	}{
		{
			inType:          experimentsv1beta1.ParameterTypeDiscrete,
			expectedType:    suggestionapi.ParameterType_DISCRETE,
			testDescription: "Convert discrete parameter type",
		},
		{
			inType:          experimentsv1beta1.ParameterTypeCategorical,
			expectedType:    suggestionapi.ParameterType_CATEGORICAL,
			testDescription: "Convert categorical parameter type",
		},
		{
			inType:          experimentsv1beta1.ParameterTypeDouble,
			expectedType:    suggestionapi.ParameterType_DOUBLE,
			testDescription: "Convert double parameter type",
		},
		{
			inType:          experimentsv1beta1.ParameterTypeInt,
			expectedType:    suggestionapi.ParameterType_INT,
			testDescription: "Convert int parameter type",
		},
		{
			inType:          experimentsv1beta1.ParameterTypeUnknown,
			expectedType:    suggestionapi.ParameterType_UNKNOWN_TYPE,
			testDescription: "Convert unknown parameter type",
		},
	}
	for _, tc := range tcs {
		actualType := convertParameterType(tc.inType)
		if actualType != tc.expectedType {
			t.Errorf("Case: %v failed. Expected parameter type %v, got %v", tc.testDescription, tc.expectedType, actualType)
		}
	}
}

func TestConvertFeasibleSpace(t *testing.T) {

	tcs := []struct {
		inFeasibleSpace       experimentsv1beta1.FeasibleSpace
		expectedFeasibleSpace *suggestionapi.FeasibleSpace
		testDescription       string
	}{
		{
			inFeasibleSpace: experimentsv1beta1.FeasibleSpace{
				Max:          "10",
				Min:          "1",
				List:         []string{"1", "2", "3"},
				Step:         "1",
				Distribution: experimentsv1beta1.DistributionUnknown,
			},
			expectedFeasibleSpace: &suggestionapi.FeasibleSpace{
				Max:  "10",
				Min:  "1",
				List: []string{"1", "2", "3"},
				Step: "1",
			},
			testDescription: "Convert feasible space with unknown distribution",
		},
		{
			inFeasibleSpace: experimentsv1beta1.FeasibleSpace{
				Max:          "100",
				Min:          "10",
				Step:         "10",
				Distribution: experimentsv1beta1.DistributionUniform,
			},
			expectedFeasibleSpace: &suggestionapi.FeasibleSpace{
				Max:          "100",
				Min:          "10",
				Step:         "10",
				Distribution: suggestionapi.Distribution_UNIFORM,
			},
			testDescription: "Convert feasible space with uniform distribution",
		},
	}

	for _, tc := range tcs {
		actualFeasibleSpace := convertFeasibleSpace(tc.inFeasibleSpace)
		if diff := cmp.Diff(tc.expectedFeasibleSpace, actualFeasibleSpace, cmpopts.IgnoreUnexported(suggestionapi.FeasibleSpace{})); diff != "" {
			t.Errorf("Case: %v failed. Unexpected difference (-want +got):\n%s", tc.testDescription, diff)
		}
	}
}

func TestConvertDistribution(t *testing.T) {

	tcs := []struct {
		inDistribution       experimentsv1beta1.Distribution
		expectedDistribution suggestionapi.Distribution
		testDescription      string
	}{
		{
			inDistribution:       experimentsv1beta1.DistributionUniform,
			expectedDistribution: suggestionapi.Distribution_UNIFORM,
			testDescription:      "Convert uniform distribution",
		},
		{
			inDistribution:       experimentsv1beta1.DistributionLogUniform,
			expectedDistribution: suggestionapi.Distribution_LOG_UNIFORM,
			testDescription:      "Convert log-uniform distribution",
		},
		{
			inDistribution:       experimentsv1beta1.DistributionNormal,
			expectedDistribution: suggestionapi.Distribution_NORMAL,
			testDescription:      "Convert normal distribution",
		},
		{
			inDistribution:       experimentsv1beta1.DistributionLogNormal,
			expectedDistribution: suggestionapi.Distribution_LOG_NORMAL,
			testDescription:      "Convert log-normal distribution",
		},
		{
			inDistribution:       experimentsv1beta1.DistributionUnknown,
			expectedDistribution: suggestionapi.Distribution_DISTRIBUTION_UNSPECIFIED,
			testDescription:      "Convert unknown distribution",
		},
	}

	for _, tc := range tcs {
		actualDistribution := convertDistribution(tc.inDistribution)
		if actualDistribution != tc.expectedDistribution {
			t.Errorf("Case: %v failed. Expected distribution %v, got %v", tc.testDescription, tc.expectedDistribution, actualDistribution)
		}
	}
}

func TestConvertTrialObservation(t *testing.T) {

	tcs := []struct {
		strategies          []commonv1beta1.MetricStrategy
		inObservation       *commonv1beta1.Observation
		expectedObservation *suggestionapi.Observation
		testDescription     string
	}{
		{
			strategies:          newFakeStrategies(),
			inObservation:       newFakeTrialObservation(),
			expectedObservation: newFakeRequestObservation(),
			testDescription:     "Run with min, max and latest metrics extract",
		},
		{
			strategies: newFakeStrategies(),
			inObservation: func() *commonapiv1beta1.Observation {
				obsIn := newFakeTrialObservation()
				obsIn.Metrics[0].Min = consts.UnavailableMetricValue
				return obsIn
			}(),
			expectedObservation: func() *suggestionapi.Observation {
				obsOut := newFakeRequestObservation()
				obsOut.Metrics[0].Value = "0.05"
				return obsOut
			}(),
			testDescription: "Observation doesn't have min metric, latest is assigned",
		},
		{
			strategies: newFakeStrategies(),
			inObservation: func() *commonapiv1beta1.Observation {
				obsIn := newFakeTrialObservation()
				obsIn.Metrics[1].Max = consts.UnavailableMetricValue
				return obsIn
			}(),
			expectedObservation: func() *suggestionapi.Observation {
				obsOut := newFakeRequestObservation()
				obsOut.Metrics[1].Value = "0.90"
				return obsOut
			}(),
			testDescription: "Observation doesn't have max metric, latest is assigned",
		},
	}
	for _, tc := range tcs {
		actualObservation := convertTrialObservation(tc.strategies, tc.inObservation)
		if !reflect.DeepEqual(actualObservation, tc.expectedObservation) {
			t.Errorf("Case: %v failed.\nExpected observation: %v \ngot: %v", tc.testDescription, tc.expectedObservation, actualObservation)
		}
	}
}

func TestConvertComparison(t *testing.T) {
	tcs := []struct {
		inComparison       suggestionapi.ComparisonType
		expectedComparison commonapiv1beta1.ComparisonType
		testDescription    string
	}{
		{
			inComparison:       suggestionapi.ComparisonType_EQUAL,
			expectedComparison: commonapiv1beta1.ComparisonTypeEqual,
			testDescription:    "Convert equal comparison type",
		},
		{
			inComparison:       suggestionapi.ComparisonType_LESS,
			expectedComparison: commonapiv1beta1.ComparisonTypeLess,
			testDescription:    "Convert less comparison type",
		},
		{
			inComparison:       suggestionapi.ComparisonType_GREATER,
			expectedComparison: commonapiv1beta1.ComparisonTypeGreater,
			testDescription:    "Convert greater comparison type",
		},
		{
			inComparison:       suggestionapi.ComparisonType_UNKNOWN_COMPARISON,
			expectedComparison: commonapiv1beta1.ComparisonTypeEqual,
			testDescription:    "Convert unknown comparison type",
		},
	}
	for _, tc := range tcs {
		actualComparison := convertComparison(tc.inComparison)
		if actualComparison != tc.expectedComparison {
			t.Errorf("Case: %v failed. Expected comparison type %v, got %v", tc.testDescription, tc.expectedComparison, actualComparison)
		}
	}
}

func newFakeStrategies() []commonv1beta1.MetricStrategy {
	return []commonv1beta1.MetricStrategy{
		{Name: "error", Value: commonv1beta1.ExtractByMin},
		{Name: "auc", Value: commonv1beta1.ExtractByMax},
		{Name: "accuracy", Value: commonv1beta1.ExtractByLatest},
	}
}

func newFakeTrialObservation() *commonv1beta1.Observation {
	return &commonv1beta1.Observation{
		Metrics: []commonv1beta1.Metric{
			{Name: "error", Min: "0.01", Max: "0.08", Latest: "0.05"},
			{Name: "auc", Min: "0.70", Max: "0.95", Latest: "0.90"},
			{Name: "accuracy", Min: "0.8", Max: "0.94", Latest: "0.93"},
		},
	}
}

func newFakeSuggestionTrialObservation() *commonv1beta1.Observation {
	return &commonv1beta1.Observation{
		Metrics: []commonv1beta1.Metric{
			{Name: "metric1-name", Min: "0.95", Max: "0.95", Latest: "0.95"},
		},
	}
}

func newFakeRequestObservation() *suggestionapi.Observation {
	return &suggestionapi.Observation{
		Metrics: []*suggestionapi.Metric{
			{
				Name:  "error",
				Value: "0.01",
			},
			{
				Name:  "auc",
				Value: "0.95",
			},
			{
				Name:  "accuracy",
				Value: "0.93",
			},
		},
	}
}

func newFakeObjective() *commonapiv1beta1.ObjectiveSpec {
	goal := 0.99

	return &commonv1beta1.ObjectiveSpec{
		Type:                  commonv1beta1.ObjectiveTypeMaximize,
		ObjectiveMetricName:   "metric1-name",
		AdditionalMetricNames: []string{"metric2-name"},
		Goal:                  &goal,
		MetricStrategies: []commonapiv1beta1.MetricStrategy{
			{Name: "metric1-name", Value: commonapiv1beta1.ExtractByLatest},
		},
	}
}

func newFakeTime() *metav1.Time {
	fakeTime, _ := time.Parse(timeFormat, "2020-01-01T15:04:05Z")
	return &metav1.Time{
		Time: fakeTime,
	}
}

func newFakeExperiment() *experimentsv1beta1.Experiment {
	var testInt int32 = 1

	fakeParameters := []experimentsv1beta1.ParameterSpec{
		{
			Name:          "param1-name",
			ParameterType: experimentsv1beta1.ParameterTypeInt,
			FeasibleSpace: experimentsv1beta1.FeasibleSpace{
				Max: "5",
				Min: "1",
			},
		},
		{
			Name:          "param2-name",
			ParameterType: experimentsv1beta1.ParameterTypeDouble,
			FeasibleSpace: experimentsv1beta1.FeasibleSpace{
				Max:  "0.4",
				Min:  "0.1",
				Step: "0.1",
			},
		},
	}

	return &experimentsv1beta1.Experiment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "experiment-name",
			Namespace: "namespace",
		},
		Spec: experimentsv1beta1.ExperimentSpec{
			ParallelTrialCount: &testInt,
			MaxTrialCount:      &testInt,
			Algorithm: &commonv1beta1.AlgorithmSpec{
				AlgorithmName: algorithmName,
				AlgorithmSettings: []commonv1beta1.AlgorithmSetting{
					{
						Name:  "overridden-name",
						Value: "value",
					},
				},
			},
			EarlyStopping: &commonapiv1beta1.EarlyStoppingSpec{
				AlgorithmName: earlyStoppingAlgorithmName,
				AlgorithmSettings: []commonapiv1beta1.EarlyStoppingSetting{
					{
						Name:  "setting-name",
						Value: "value",
					},
				},
			},
			Objective:  newFakeObjective(),
			Parameters: fakeParameters,
			NasConfig: &experimentsv1beta1.NasConfig{
				GraphConfig: experimentsv1beta1.GraphConfig{
					NumLayers:   &testInt,
					InputSizes:  []int32{32, 32, 3},
					OutputSizes: []int32{10},
				},
				Operations: []experimentsv1beta1.Operation{
					{
						OperationType: "operation-type",
						Parameters:    fakeParameters,
					},
				},
			},
		},
	}
}

func newFakeSuggestion() *suggestionsv1beta1.Suggestion {

	return &suggestionsv1beta1.Suggestion{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "suggestion-name",
			Namespace: "namespace",
		},
		Spec: suggestionsv1beta1.SuggestionSpec{
			Algorithm: &commonv1beta1.AlgorithmSpec{
				AlgorithmName: algorithmName,
				AlgorithmSettings: []commonv1beta1.AlgorithmSetting{
					{
						Name:  "overridden-name",
						Value: "value",
					},
				},
			},
			EarlyStopping: &commonapiv1beta1.EarlyStoppingSpec{
				AlgorithmName: earlyStoppingAlgorithmName,
				AlgorithmSettings: []commonapiv1beta1.EarlyStoppingSetting{
					{
						Name:  "setting-name",
						Value: "value",
					},
				},
			},
			Requests: 6,
		},
		Status: suggestionsv1beta1.SuggestionStatus{
			SuggestionCount: 4,
			AlgorithmSettings: []commonv1beta1.AlgorithmSetting{
				{
					Name:  "overridden-name",
					Value: "overridden-value",
				},
				{
					Name:  "new-setting",
					Value: "value",
				},
			},
		},
	}
}

func newFakeTrials() []trialsv1beta1.Trial {

	fakeConditions := []trialsv1beta1.TrialCondition{
		{
			Type: trialsv1beta1.TrialSucceeded,
		},
	}

	fakeEarlyStoppedConditions := []trialsv1beta1.TrialCondition{
		{
			Type:   trialsv1beta1.TrialEarlyStopped,
			Status: corev1.ConditionTrue,
		},
	}

	return []trialsv1beta1.Trial{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "trial1-name",
				Namespace: "namespace",
			},
			Spec: trialsv1beta1.TrialSpec{
				Objective: newFakeObjective(),
				ParameterAssignments: []commonapiv1beta1.ParameterAssignment{
					{
						Name:  "param1-name",
						Value: "1",
					},
					{
						Name:  "param2-name",
						Value: "0.1",
					},
				},
				Labels: map[string]string{},
			},
			Status: trialsv1beta1.TrialStatus{
				StartTime:      newFakeTime(),
				CompletionTime: newFakeTime(),
				Conditions:     fakeConditions,
				Observation:    newFakeSuggestionTrialObservation(),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "trial2-name",
				Namespace: "namespace",
			},
			Spec: trialsv1beta1.TrialSpec{
				Objective: newFakeObjective(),
				ParameterAssignments: []commonapiv1beta1.ParameterAssignment{
					{
						Name:  "param1-name",
						Value: "2",
					},
					{
						Name:  "param2-name",
						Value: "0.2",
					},
				},
				Labels: map[string]string{},
			},
			Status: trialsv1beta1.TrialStatus{
				Conditions:  fakeConditions,
				Observation: newFakeSuggestionTrialObservation(),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "trial3-name",
				Namespace: "namespace",
			},
			Status: trialsv1beta1.TrialStatus{
				Conditions: []trialsv1beta1.TrialCondition{
					{
						Type:    trialsv1beta1.TrialMetricsUnavailable,
						Status:  corev1.ConditionTrue,
						Message: "Metrics are not available",
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "trial4-name",
				Namespace: "namespace",
			},
			Spec: trialsv1beta1.TrialSpec{
				Objective: newFakeObjective(),
				ParameterAssignments: []commonapiv1beta1.ParameterAssignment{
					{
						Name:  "param1-name",
						Value: "4",
					},
					{
						Name:  "param2-name",
						Value: "0.4",
					},
				},
				Labels: map[string]string{},
			},
			Status: trialsv1beta1.TrialStatus{
				Conditions:  fakeEarlyStoppedConditions,
				Observation: nil,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "trial5-name",
				Namespace: "namespace",
			},
			Spec: trialsv1beta1.TrialSpec{
				Objective: newFakeObjective(),
				ParameterAssignments: []commonapiv1beta1.ParameterAssignment{
					{
						Name:  "param1-name",
						Value: "5",
					},
					{
						Name:  "param2-name",
						Value: "0.5",
					},
				},
				Labels: map[string]string{},
			},
			Status: trialsv1beta1.TrialStatus{
				Conditions:  fakeEarlyStoppedConditions,
				Observation: newFakeSuggestionTrialObservation(),
			},
		},
	}
}

func newFakeRequest() *suggestionapi.GetSuggestionsRequest {

	fakeParameters := []*suggestionapi.ParameterSpec{
		{
			Name:          "param1-name",
			ParameterType: suggestionapi.ParameterType_INT,
			FeasibleSpace: &suggestionapi.FeasibleSpace{
				Max: "5",
				Min: "1",
			},
		},
		{
			Name:          "param2-name",
			ParameterType: suggestionapi.ParameterType_DOUBLE,
			FeasibleSpace: &suggestionapi.FeasibleSpace{
				Max:  "0.4",
				Min:  "0.1",
				Step: "0.1",
			},
		},
	}

	fakeLabels := make(map[string]string)

	fakeObjective := &suggestionapi.ObjectiveSpec{
		Type:                  suggestionapi.ObjectiveType_MAXIMIZE,
		ObjectiveMetricName:   "metric1-name",
		AdditionalMetricNames: []string{"metric2-name"},
		Goal:                  0.99,
	}

	return &suggestionapi.GetSuggestionsRequest{
		Experiment: &suggestionapi.Experiment{
			Name: "experiment-name",
			Spec: &suggestionapi.ExperimentSpec{
				Algorithm: &suggestionapi.AlgorithmSpec{
					AlgorithmName: algorithmName,
					AlgorithmSettings: []*suggestionapi.AlgorithmSetting{
						{
							Name:  "overridden-name",
							Value: "overridden-value",
						},
						{
							Name:  "new-setting",
							Value: "value",
						},
					},
				},
				EarlyStopping: &suggestionapi.EarlyStoppingSpec{
					AlgorithmName: earlyStoppingAlgorithmName,
					AlgorithmSettings: []*suggestionapi.EarlyStoppingSetting{
						{
							Name:  "setting-name",
							Value: "value",
						},
					},
				},
				Objective:          fakeObjective,
				ParallelTrialCount: 1,
				MaxTrialCount:      1,
				ParameterSpecs: &suggestionapi.ExperimentSpec_ParameterSpecs{
					Parameters: fakeParameters,
				},
				NasConfig: &suggestionapi.NasConfig{
					GraphConfig: &suggestionapi.GraphConfig{
						NumLayers:   1,
						InputSizes:  []int32{32, 32, 3},
						OutputSizes: []int32{10},
					},
					Operations: &suggestionapi.NasConfig_Operations{
						Operation: []*suggestionapi.Operation{
							{
								OperationType: "operation-type",
								ParameterSpecs: &suggestionapi.Operation_ParameterSpecs{
									Parameters: fakeParameters,
								},
							},
						},
					},
				},
			},
		},
		Trials: []*suggestionapi.Trial{
			{
				Name: "trial1-name",
				Spec: &suggestionapi.TrialSpec{
					Objective: fakeObjective,
					ParameterAssignments: &suggestionapi.TrialSpec_ParameterAssignments{
						Assignments: []*suggestionapi.ParameterAssignment{
							{
								Name:  "param1-name",
								Value: "1",
							},
							{
								Name:  "param2-name",
								Value: "0.1",
							},
						},
					},
					Labels: fakeLabels,
				},
				Status: &suggestionapi.TrialStatus{
					StartTime:      newFakeTime().Format(timeFormat),
					CompletionTime: newFakeTime().Format(timeFormat),
					Condition:      suggestionapi.TrialStatus_SUCCEEDED,
					Observation: &suggestionapi.Observation{
						Metrics: []*suggestionapi.Metric{
							{
								Name:  "metric1-name",
								Value: "0.95",
							},
						},
					},
				},
			},
			{
				Name: "trial2-name",
				Spec: &suggestionapi.TrialSpec{
					Objective: fakeObjective,
					ParameterAssignments: &suggestionapi.TrialSpec_ParameterAssignments{
						Assignments: []*suggestionapi.ParameterAssignment{
							{
								Name:  "param1-name",
								Value: "2",
							},
							{
								Name:  "param2-name",
								Value: "0.2",
							},
						},
					},
					Labels: fakeLabels,
				},
				Status: &suggestionapi.TrialStatus{
					StartTime:      "",
					CompletionTime: "",
					Condition:      suggestionapi.TrialStatus_SUCCEEDED,
					Observation: &suggestionapi.Observation{
						Metrics: []*suggestionapi.Metric{
							{
								Name:  "metric1-name",
								Value: "0.95",
							},
						},
					},
				},
			},
			{
				Name: "trial5-name",
				Spec: &suggestionapi.TrialSpec{
					Objective: fakeObjective,
					ParameterAssignments: &suggestionapi.TrialSpec_ParameterAssignments{
						Assignments: []*suggestionapi.ParameterAssignment{
							{
								Name:  "param1-name",
								Value: "5",
							},
							{
								Name:  "param2-name",
								Value: "0.5",
							},
						},
					},
					Labels: fakeLabels,
				},
				Status: &suggestionapi.TrialStatus{
					StartTime:      "",
					CompletionTime: "",
					Condition:      suggestionapi.TrialStatus_EARLYSTOPPED,
					Observation: &suggestionapi.Observation{
						Metrics: []*suggestionapi.Metric{
							{
								Name:  "metric1-name",
								Value: "0.95",
							},
						},
					},
				},
			},
		},
		CurrentRequestNumber: 2,
		TotalRequestNumber:   6,
	}
}
