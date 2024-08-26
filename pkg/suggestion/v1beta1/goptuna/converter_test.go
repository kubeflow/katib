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

package suggestion_goptuna_v1beta1

import (
	"testing"

	"github.com/c-bata/goptuna"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	api_v1_beta1 "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
)

func Test_toGoptunaDirection(t *testing.T) {
	for name, tc := range map[string]struct {
		objectiveType api_v1_beta1.ObjectiveType
		wantDirection goptuna.StudyDirection
	}{
		"minimize": {
			objectiveType: api_v1_beta1.ObjectiveType_MINIMIZE,
			wantDirection: goptuna.StudyDirectionMinimize,
		},
		"maximize": {
			objectiveType: api_v1_beta1.ObjectiveType_MAXIMIZE,
			wantDirection: goptuna.StudyDirectionMaximize,
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := toGoptunaDirection(tc.objectiveType)
			if diff := cmp.Diff(tc.wantDirection, got); len(diff) != 0 {
				t.Errorf("Unexpected direction from toGoptunaDirection (-want,+got):\n%s", diff)
			}
		})
	}
}

func Test_toGoptunaSearchSpace(t *testing.T) {
	cases := map[string]struct {
		parameters      []*api_v1_beta1.ParameterSpec
		wantSearchSpace map[string]interface{}
		wantError       error
	}{
		"Double parameter type": {
			parameters: []*api_v1_beta1.ParameterSpec{
				{
					Name:          "param-double",
					ParameterType: api_v1_beta1.ParameterType_DOUBLE,
					FeasibleSpace: &api_v1_beta1.FeasibleSpace{
						Max: "5.5",
						Min: "1.5",
					},
				},
			},
			wantSearchSpace: map[string]interface{}{
				"param-double": goptuna.UniformDistribution{
					High: 5.5,
					Low:  1.5,
				},
			},
		},
		"Double parameter type with step": {
			parameters: []*api_v1_beta1.ParameterSpec{
				{
					Name:          "param-double",
					ParameterType: api_v1_beta1.ParameterType_DOUBLE,
					FeasibleSpace: &api_v1_beta1.FeasibleSpace{
						Max:  "5.5",
						Min:  "1.5",
						Step: "0.5",
					},
				},
			},
			wantSearchSpace: map[string]interface{}{
				"param-double": goptuna.DiscreteUniformDistribution{
					High: 5.5,
					Low:  1.5,
					Q:    0.5,
				},
			},
		},
		"Int parameter type": {
			parameters: []*api_v1_beta1.ParameterSpec{
				{
					Name:          "param-int",
					ParameterType: api_v1_beta1.ParameterType_INT,
					FeasibleSpace: &api_v1_beta1.FeasibleSpace{
						Max: "5",
						Min: "1",
					},
				},
			},
			wantSearchSpace: map[string]interface{}{
				"param-int": goptuna.IntUniformDistribution{
					High: 5,
					Low:  1,
				},
			},
		},
		"Int parameter type with step": {
			parameters: []*api_v1_beta1.ParameterSpec{
				{
					Name:          "param-int",
					ParameterType: api_v1_beta1.ParameterType_INT,
					FeasibleSpace: &api_v1_beta1.FeasibleSpace{
						Max:  "5",
						Min:  "1",
						Step: "2",
					},
				},
			},
			wantSearchSpace: map[string]interface{}{
				"param-int": goptuna.StepIntUniformDistribution{
					High: 5,
					Low:  1,
					Step: 2,
				},
			},
		},
		"Discrete parameter type": {
			parameters: []*api_v1_beta1.ParameterSpec{
				{
					Name:          "param-discrete",
					ParameterType: api_v1_beta1.ParameterType_DISCRETE,
					FeasibleSpace: &api_v1_beta1.FeasibleSpace{
						List: []string{"3", "2", "6"},
					},
				},
			},
			wantSearchSpace: map[string]interface{}{
				"param-discrete": goptuna.CategoricalDistribution{
					Choices: []string{"3", "2", "6"},
				},
			},
		},
		"Categorical parameter type": {
			parameters: []*api_v1_beta1.ParameterSpec{
				{
					Name:          "param-categorical",
					ParameterType: api_v1_beta1.ParameterType_CATEGORICAL,
					FeasibleSpace: &api_v1_beta1.FeasibleSpace{
						List: []string{"cat1", "cat2", "cat3"},
					},
				},
			},
			wantSearchSpace: map[string]interface{}{
				"param-categorical": goptuna.CategoricalDistribution{
					Choices: []string{"cat1", "cat2", "cat3"},
				},
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := toGoptunaSearchSpace(tc.parameters)
			if diff := cmp.Diff(tc.wantError, err, cmpopts.EquateErrors()); len(diff) != 0 {
				t.Errorf("Unexpected error from toGoptunaSearchSpace (-want,+got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantSearchSpace, got); len(diff) != 0 {
				t.Errorf("Unexpected search space from toGoptunaSearchSpace (-want,+got):\n%s", diff)
			}
		})
	}
}
