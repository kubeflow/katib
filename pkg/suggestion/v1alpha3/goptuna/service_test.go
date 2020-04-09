package suggestion_goptuna_v1alpha3_test

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	api_v1_alpha3 "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
	suggestion_goptuna_v1alpha3 "github.com/kubeflow/katib/pkg/suggestion/v1alpha3/goptuna"
)

func TestSuggestionService_GetSuggestions(t *testing.T) {
	ctx := context.TODO()
	for _, tt := range []struct {
		name         string
		req          *api_v1_alpha3.GetSuggestionsRequest
		expectedCode codes.Code
	}{
		{
			name:         "empty request",
			req:          nil,
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "cmaes request",
			req: &api_v1_alpha3.GetSuggestionsRequest{
				Experiment: &api_v1_alpha3.Experiment{
					Name: "test",
					Spec: &api_v1_alpha3.ExperimentSpec{
						Algorithm: &api_v1_alpha3.AlgorithmSpec{
							AlgorithmName: "cmaes",
							AlgorithmSetting: []*api_v1_alpha3.AlgorithmSetting{
								{
									Name:  "random_state",
									Value: "10",
								},
							},
							EarlyStoppingSpec: nil,
						},
						Objective: &api_v1_alpha3.ObjectiveSpec{
							Type:                  api_v1_alpha3.ObjectiveType_MINIMIZE,
							Goal:                  0.1,
							ObjectiveMetricName:   "metric-1",
							AdditionalMetricNames: nil,
						},
						ParameterSpecs: &api_v1_alpha3.ExperimentSpec_ParameterSpecs{
							Parameters: []*api_v1_alpha3.ParameterSpec{
								{
									Name:          "param-1",
									ParameterType: api_v1_alpha3.ParameterType_INT,
									FeasibleSpace: &api_v1_alpha3.FeasibleSpace{
										Max:  "10",
										Min:  "-10",
										List: nil,
										Step: "",
									},
								},
								{
									Name:          "param-2",
									ParameterType: api_v1_alpha3.ParameterType_CATEGORICAL,
									FeasibleSpace: &api_v1_alpha3.FeasibleSpace{
										List: []string{"cat1", "cat2", "cat3"},
									},
								},
								{
									Name:          "param-3",
									ParameterType: api_v1_alpha3.ParameterType_DISCRETE,
									FeasibleSpace: &api_v1_alpha3.FeasibleSpace{
										List: []string{"3", "2", "6"},
									},
								},
								{
									Name:          "param-4",
									ParameterType: api_v1_alpha3.ParameterType_DOUBLE,
									FeasibleSpace: &api_v1_alpha3.FeasibleSpace{
										Max: "5.5",
										Min: "-1.5",
									},
								},
							},
						},
						TrialTemplate:        "",
						MetricsCollectorSpec: "",
						ParallelTrialCount:   0,
						MaxTrialCount:        0,
						NasConfig:            nil,
					},
				},
				Trials: []*api_v1_alpha3.Trial{
					{
						Name: "test-asfjh",
						Spec: &api_v1_alpha3.TrialSpec{
							ExperimentName: "",
							Objective: &api_v1_alpha3.ObjectiveSpec{
								Type:                  api_v1_alpha3.ObjectiveType_MAXIMIZE,
								Goal:                  0.9,
								ObjectiveMetricName:   "metric-2",
								AdditionalMetricNames: nil,
							},
							ParameterAssignments: &api_v1_alpha3.TrialSpec_ParameterAssignments{
								Assignments: []*api_v1_alpha3.ParameterAssignment{
									{
										Name:  "param-1",
										Value: "2",
									},
									{
										Name:  "param-2",
										Value: "cat1",
									},
									{
										Name:  "param-3",
										Value: "2",
									},
									{
										Name:  "param-4",
										Value: "3.44",
									},
								},
							},
							RunSpec:              "",
							MetricsCollectorSpec: "",
						},
						Status: &api_v1_alpha3.TrialStatus{
							StartTime:      time.Now().Format(time.RFC3339Nano),
							CompletionTime: time.Now().Format(time.RFC3339Nano),
							Condition:      api_v1_alpha3.TrialStatus_SUCCEEDED,
							Observation: &api_v1_alpha3.Observation{
								Metrics: []*api_v1_alpha3.Metric{
									{
										Name:  "metric-1",
										Value: "435",
									},
									{
										Name:  "metric-2",
										Value: "5643",
									},
								},
							},
						},
					},
					{
						Name: "test-234hs",
						Spec: &api_v1_alpha3.TrialSpec{
							ExperimentName: "",
							Objective: &api_v1_alpha3.ObjectiveSpec{
								Type:                  api_v1_alpha3.ObjectiveType_MAXIMIZE,
								Goal:                  0.9,
								ObjectiveMetricName:   "metric-2",
								AdditionalMetricNames: nil,
							},
							ParameterAssignments: &api_v1_alpha3.TrialSpec_ParameterAssignments{
								Assignments: []*api_v1_alpha3.ParameterAssignment{
									{
										Name:  "param-1",
										Value: "3",
									},
									{
										Name:  "param-2",
										Value: "cat2",
									},
									{
										Name:  "param-3",
										Value: "6",
									},
									{
										Name:  "param-4",
										Value: "4.44",
									},
								},
							},
							RunSpec:              "",
							MetricsCollectorSpec: "",
						},
						Status: &api_v1_alpha3.TrialStatus{
							StartTime:      time.Now().Format(time.RFC3339Nano),
							CompletionTime: time.Now().Format(time.RFC3339Nano),
							Condition:      api_v1_alpha3.TrialStatus_SUCCEEDED,
							Observation: &api_v1_alpha3.Observation{
								Metrics: []*api_v1_alpha3.Metric{
									{
										Name:  "metric-1",
										Value: "123",
									},
									{
										Name:  "metric-2",
										Value: "3028",
									},
								},
							},
						},
					},
				},
				RequestNumber: 2,
			},
			expectedCode: codes.OK,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s := &suggestion_goptuna_v1alpha3.SuggestionService{}
			reply, err := s.GetSuggestions(ctx, tt.req)

			c, ok := status.FromError(err)
			if !ok {
				t.Errorf("GeteSuggestion() returns non-gRPC error")
			}
			if tt.expectedCode != c.Code() {
				t.Errorf("GetSuggestions() should return = %v, but got %v", tt.expectedCode, c.Code())
			}

			if tt.expectedCode != codes.OK {
				return
			}

			if len(reply.ParameterAssignments) != int(tt.req.RequestNumber) {
				t.Errorf("GetSuggestions() should return %d suggestions, but got %#v", tt.req.RequestNumber, reply.ParameterAssignments)
				return
			}

			params := tt.req.GetExperiment().GetSpec().GetParameterSpecs().GetParameters()
			for i := range reply.ParameterAssignments {
				assignments := reply.ParameterAssignments[i].Assignments
				if len(assignments) != len(params) {
					t.Errorf("Each assignments should holds %d parameters, but got %#v", len(params), assignments)
					return
				}
			}
		})
	}
}
