package suggestionclient

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	commonapiv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	suggestionapi "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	suggestionapimock "github.com/kubeflow/katib/pkg/mock/v1beta1/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type k8sMatcher struct {
	x interface{}
}

func (k8s k8sMatcher) Matches(x interface{}) bool {
	return equality.Semantic.DeepEqual(k8s.x, x)
}

func (k8s k8sMatcher) String() string {
	return fmt.Sprintf("is equal to %v", k8s.x)
}

func TestGetRPCClient(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	fakeConn := &grpc.ClientConn{}
	actualClient := getRPCClient(fakeConn)
	g.Expect(actualClient).To(gomega.Equal(suggestionapi.NewSuggestionClient(fakeConn)))
}

func TestSyncAssignments(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	rpcClient := suggestionapimock.NewMockSuggestionClient(mockCtrl)

	getRPCClient = func(conn *grpc.ClientConn) suggestionapi.SuggestionClient {
		return rpcClient
	}
	suggestionClient := New()

	experiment := newFakeExperiment()
	suggestion := newFakeSuggestion()
	trials := newFakeTrials()

	expectedRequest := newFakeRequest()

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
					Name:  "overridded-name",
					Value: "suggestion-value",
				},
				{
					Name:  "new-suggestion-setting-name",
					Value: "value",
				},
			},
		},
	}

	validRun := rpcClient.EXPECT().GetSuggestions(gomock.Any(), k8sMatcher{expectedRequest}).Return(getSuggestionReply, nil)

	getSuggestionsFail := rpcClient.EXPECT().GetSuggestions(gomock.Any(), gomock.Any()).Return(nil, errors.New("Suggestion service connection error"))

	invalidAssignmentsCount := rpcClient.EXPECT().GetSuggestions(gomock.Any(), gomock.Any()).Return(
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

	gomock.InOrder(
		validRun,
		getSuggestionsFail,
		invalidAssignmentsCount,
	)

	tcs := []struct {
		Experiment      *experimentsv1beta1.Experiment
		Suggestion      *suggestionsv1beta1.Suggestion
		Trials          []trialsv1beta1.Trial
		Err             bool
		TestDescription string
	}{
		// Experiment contains HP and NAS config just for the test purpose
		// validRun case
		{
			Experiment:      experiment,
			Suggestion:      suggestion,
			Trials:          trials,
			Err:             false,
			TestDescription: "SyncAssignments valid run",
		},
		{
			Suggestion: func() *suggestionsv1beta1.Suggestion {
				s := newFakeSuggestion()
				s.Spec.Requests = 4
				s.Status.SuggestionCount = 6
				return s
			}(),
			Err:             false,
			TestDescription: "Negative request number",
		},
		// getSuggestionsFail case
		{
			Experiment:      experiment,
			Suggestion:      suggestion,
			Trials:          trials,
			Err:             true,
			TestDescription: "Unable to execute GetSuggestions",
		},
		// invalidAssignmentsCount case
		{
			Experiment:      experiment,
			Suggestion:      suggestion,
			Trials:          trials,
			Err:             true,
			TestDescription: "ParameterAssignments from response != request number",
		},
	}
	for _, tc := range tcs {
		err := suggestionClient.SyncAssignments(tc.Suggestion, tc.Experiment, tc.Trials)
		if !tc.Err && err != nil {
			t.Errorf("Case: %v failed. Expected nil, got %v", tc.TestDescription, err)
		} else if tc.Err && err == nil {
			t.Errorf("Case: %v failed. Expected err, got nil", tc.TestDescription)
		}
	}
}

func TestValidateAlgorithmSettings(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	rpcClient := suggestionapimock.NewMockSuggestionClient(mockCtrl)

	getRPCClient = func(conn *grpc.ClientConn) suggestionapi.SuggestionClient {
		return rpcClient
	}

	expectedRequest := &suggestionapi.ValidateAlgorithmSettingsRequest{
		Experiment: newFakeRequest().Experiment,
	}
	expectedRequest.Experiment.Spec.Algorithm.AlgorithmSettings = []*suggestionapi.AlgorithmSetting{
		{
			Name:  "overridded-name",
			Value: "value",
		},
	}

	validRun := rpcClient.EXPECT().ValidateAlgorithmSettings(gomock.Any(), k8sMatcher{expectedRequest}, gomock.Any()).Return(nil, nil)

	invalidExperiment := rpcClient.EXPECT().ValidateAlgorithmSettings(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil,
		status.Error(codes.InvalidArgument, "Invalid experiment parameter"))
	connectionError := rpcClient.EXPECT().ValidateAlgorithmSettings(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil,
		status.Error(codes.Unavailable, "Unable to connect"))
	unimplementedMethod := rpcClient.EXPECT().ValidateAlgorithmSettings(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil,
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
		Experiment      *experimentsv1beta1.Experiment
		Suggestion      *suggestionsv1beta1.Suggestion
		Err             bool
		TestDescription string
	}{
		// validRun case
		{
			Experiment:      exp,
			Suggestion:      sug,
			Err:             false,
			TestDescription: "ValidateAlgorithmSettings valid run",
		},
		// invalidExperiment case
		{
			Experiment:      exp,
			Suggestion:      sug,
			Err:             true,
			TestDescription: "Invalid argument return in Experiment validation",
		},
		// connectionError case
		{
			Experiment:      exp,
			Suggestion:      sug,
			Err:             true,
			TestDescription: "Connection to suggestion service error",
		},
		// unimplementedMethod case
		{
			Experiment:      exp,
			Suggestion:      sug,
			Err:             false,
			TestDescription: "Unimplemented ValidateAlgorithmSettings method",
		},
	}
	for _, tc := range tcs {
		err := suggestionClient.ValidateAlgorithmSettings(tc.Suggestion, tc.Experiment)
		if !tc.Err && err != nil {
			t.Errorf("Case: %v failed. Expected nil, got %v", tc.TestDescription, err)
		} else if tc.Err && err == nil {
			t.Errorf("Case: %v failed. Expected err, got nil", tc.TestDescription)
		}
	}
}

func TestConvertTrialConditionType(t *testing.T) {

	tcs := []struct {
		InCondition       trialsv1beta1.TrialConditionType
		ExpectedCondition suggestionapi.TrialStatus_TrialConditionType
		TestDescription   string
	}{
		{
			InCondition:       trialsv1beta1.TrialCreated,
			ExpectedCondition: suggestionapi.TrialStatus_CREATED,
			TestDescription:   "Convert created Trial condition",
		},
		{
			InCondition:       trialsv1beta1.TrialRunning,
			ExpectedCondition: suggestionapi.TrialStatus_RUNNING,
			TestDescription:   "Convert running Trial condition",
		},
		{
			InCondition:       trialsv1beta1.TrialSucceeded,
			ExpectedCondition: suggestionapi.TrialStatus_SUCCEEDED,
			TestDescription:   "Convert succeeded Trial condition",
		},
		{
			InCondition:       trialsv1beta1.TrialKilled,
			ExpectedCondition: suggestionapi.TrialStatus_KILLED,
			TestDescription:   "Convert killed Trial condition",
		},
		{
			InCondition:       trialsv1beta1.TrialFailed,
			ExpectedCondition: suggestionapi.TrialStatus_FAILED,
			TestDescription:   "Convert failed Trial condition",
		},
		{
			InCondition:       "Unknown",
			ExpectedCondition: suggestionapi.TrialStatus_UNKNOWN,
			TestDescription:   "Convert unknown Trial condition",
		},
	}
	for _, tc := range tcs {
		actualCondition := convertTrialConditionType(tc.InCondition)
		if actualCondition != tc.ExpectedCondition {
			t.Errorf("Case: %v failed. Expected Trial condition %v, got %v", tc.TestDescription, tc.ExpectedCondition, actualCondition)
		}
	}
}

func TestConvertObjectiveType(t *testing.T) {

	tcs := []struct {
		InType          commonapiv1beta1.ObjectiveType
		ExpectedType    suggestionapi.ObjectiveType
		TestDescription string
	}{
		{
			InType:          commonv1beta1.ObjectiveTypeMaximize,
			ExpectedType:    suggestionapi.ObjectiveType_MAXIMIZE,
			TestDescription: "Convert maximize objective type",
		},

		{
			InType:          commonv1beta1.ObjectiveTypeMinimize,
			ExpectedType:    suggestionapi.ObjectiveType_MINIMIZE,
			TestDescription: "Convert minimize objective type",
		},
		{
			InType:          commonv1beta1.ObjectiveTypeUnknown,
			ExpectedType:    suggestionapi.ObjectiveType_UNKNOWN,
			TestDescription: "Convert unknown objective type",
		},
	}
	for _, tc := range tcs {
		actualType := convertObjectiveType(tc.InType)
		if actualType != tc.ExpectedType {
			t.Errorf("Case: %v failed. Expected objective type %v, got %v", tc.TestDescription, tc.ExpectedType, actualType)
		}
	}
}

func TestConvertParameterType(t *testing.T) {

	tcs := []struct {
		InType          experimentsv1beta1.ParameterType
		ExpectedType    suggestionapi.ParameterType
		TestDescription string
	}{
		{
			InType:          experimentsv1beta1.ParameterTypeDiscrete,
			ExpectedType:    suggestionapi.ParameterType_DISCRETE,
			TestDescription: "Convert discrete parameter type",
		},
		{
			InType:          experimentsv1beta1.ParameterTypeCategorical,
			ExpectedType:    suggestionapi.ParameterType_CATEGORICAL,
			TestDescription: "Convert categorical parameter type",
		},
		{
			InType:          experimentsv1beta1.ParameterTypeDouble,
			ExpectedType:    suggestionapi.ParameterType_DOUBLE,
			TestDescription: "Convert double parameter type",
		},
		{
			InType:          experimentsv1beta1.ParameterTypeInt,
			ExpectedType:    suggestionapi.ParameterType_INT,
			TestDescription: "Convert int parameter type",
		},
		{
			InType:          experimentsv1beta1.ParameterTypeUnknown,
			ExpectedType:    suggestionapi.ParameterType_UNKNOWN_TYPE,
			TestDescription: "Convert unknown parameter type",
		},
	}
	for _, tc := range tcs {
		actualType := convertParameterType(tc.InType)
		if actualType != tc.ExpectedType {
			t.Errorf("Case: %v failed. Expected parameter type %v, got %v", tc.TestDescription, tc.ExpectedType, actualType)
		}
	}
}

func TestConvertTrialObservation(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	var strategies = []commonv1beta1.MetricStrategy{
		{Name: "error", Value: commonv1beta1.ExtractByMin},
		{Name: "auc", Value: commonv1beta1.ExtractByMax},
		{Name: "accuracy", Value: commonv1beta1.ExtractByLatest},
	}
	var observation = &commonv1beta1.Observation{
		Metrics: []commonv1beta1.Metric{
			{Name: "error", Min: 0.01, Max: 0.08, Latest: "0.05"},
			{Name: "auc", Min: 0.70, Max: 0.95, Latest: "0.90"},
			{Name: "accuracy", Min: 0.8, Max: 0.94, Latest: "0.93"},
		},
	}
	obsPb := convertTrialObservation(strategies, observation)
	g.Expect(obsPb.Metrics[0].Name).To(gomega.Equal("error"))
	value, _ := strconv.ParseFloat(obsPb.Metrics[0].Value, 64)
	g.Expect(value).To(gomega.Equal(0.01))
	g.Expect(obsPb.Metrics[1].Name).To(gomega.Equal("auc"))
	value, _ = strconv.ParseFloat(obsPb.Metrics[1].Value, 64)
	g.Expect(value).To(gomega.Equal(0.95))
	g.Expect(obsPb.Metrics[2].Name).To(gomega.Equal("accuracy"))
	value, _ = strconv.ParseFloat(obsPb.Metrics[2].Value, 64)
	g.Expect(value).To(gomega.Equal(0.93))
}

func newFakeObjective() *commonapiv1beta1.ObjectiveSpec {
	goal := 0.99

	return &commonv1beta1.ObjectiveSpec{
		Type:                  commonv1beta1.ObjectiveTypeMaximize,
		ObjectiveMetricName:   "metric1-name",
		AdditionalMetricNames: []string{"metric2-name"},
		Goal:                  &goal,
	}
}

func newFakeTime() *metav1.Time {
	fakeTime, _ := time.Parse(timeFormat, "2020-01-01T15:04:05Z")
	return &metav1.Time{
		Time: fakeTime,
	}
}

func newFakeExperiment() *experimentsv1beta1.Experiment {
	var testInt int32
	testInt = 1

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
				AlgorithmName: "algorithm-name",
				AlgorithmSettings: []commonv1beta1.AlgorithmSetting{
					{
						Name:  "overridded-name",
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
			AlgorithmName: "algorithm-name",
			Requests:      6,
		},
		Status: suggestionsv1beta1.SuggestionStatus{
			SuggestionCount: 4,
			AlgorithmSettings: []commonv1beta1.AlgorithmSetting{
				{
					Name:  "overridded-name",
					Value: "overridded-value",
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
			},
			Status: trialsv1beta1.TrialStatus{
				StartTime:      newFakeTime(),
				CompletionTime: newFakeTime(),
				Conditions:     fakeConditions,
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
			},
			Status: trialsv1beta1.TrialStatus{
				Conditions: fakeConditions,
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
					AlgorithmName: "algorithm-name",
					AlgorithmSettings: []*suggestionapi.AlgorithmSetting{
						{
							Name:  "overridded-name",
							Value: "overridded-value",
						},
						{
							Name:  "new-setting",
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
				},
				Status: &suggestionapi.TrialStatus{
					StartTime:      newFakeTime().Format(timeFormat),
					CompletionTime: newFakeTime().Format(timeFormat),
					Condition:      suggestionapi.TrialStatus_SUCCEEDED,
					Observation:    &suggestionapi.Observation{},
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
				},
				Status: &suggestionapi.TrialStatus{
					StartTime:      "",
					CompletionTime: "",
					Condition:      suggestionapi.TrialStatus_SUCCEEDED,
					Observation:    &suggestionapi.Observation{},
				},
			},
		},
		RequestNumber: 2,
	}
}
