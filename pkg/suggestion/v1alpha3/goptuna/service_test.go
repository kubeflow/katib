package suggestion_goptuna_v1alpha3_test

import (
	"context"
	"testing"

	api_v1_alpha3 "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
	suggestion_goptuna_v1alpha3 "github.com/kubeflow/katib/pkg/suggestion/v1alpha3/goptuna"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestSuggestionService_GetSuggestions(t *testing.T) {
	ctx := context.TODO()
	parameterSpecs := &api_v1_alpha3.ExperimentSpec_ParameterSpecs{
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
	}

	for _, tt := range []struct {
		name         string
		req          *api_v1_alpha3.GetSuggestionsRequest
		expectedCode codes.Code
	}{
		{
			name: "CMA-ES request",
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
						},
						Objective: &api_v1_alpha3.ObjectiveSpec{
							Type:                  api_v1_alpha3.ObjectiveType_MINIMIZE,
							Goal:                  0.1,
							ObjectiveMetricName:   "metric-1",
							AdditionalMetricNames: nil,
						},
						ParameterSpecs: parameterSpecs,
					},
				},
				RequestNumber: 2,
			},
			expectedCode: codes.OK,
		},
		{
			name: "TPE request",
			req: &api_v1_alpha3.GetSuggestionsRequest{
				Experiment: &api_v1_alpha3.Experiment{
					Name: "test",
					Spec: &api_v1_alpha3.ExperimentSpec{
						Algorithm: &api_v1_alpha3.AlgorithmSpec{
							AlgorithmName: "tpe",
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
						ParameterSpecs: parameterSpecs,
					},
				},
				RequestNumber: 2,
			},
		},
		{
			name: "Random request",
			req: &api_v1_alpha3.GetSuggestionsRequest{
				Experiment: &api_v1_alpha3.Experiment{
					Name: "test",
					Spec: &api_v1_alpha3.ExperimentSpec{
						Algorithm: &api_v1_alpha3.AlgorithmSpec{
							AlgorithmName: "random",
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
						ParameterSpecs: parameterSpecs,
					},
				},
				RequestNumber: 2,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s := &suggestion_goptuna_v1alpha3.SuggestionService{}
			reply, err := s.GetSuggestions(ctx, tt.req)

			c, ok := status.FromError(err)
			if !ok {
				t.Errorf("GetSuggestion() returns non-gRPC error")
			}
			if tt.expectedCode != c.Code() {
				t.Errorf("GetSuggestions() should return = %v, but got %v", tt.expectedCode, c.Code())
			}

			if c.Code() != codes.OK {
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
