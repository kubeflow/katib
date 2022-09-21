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
	"reflect"
	"testing"

	"github.com/c-bata/goptuna"
	api_v1_beta1 "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
)

func Test_toGoptunaDirection(t *testing.T) {
	for _, tt := range []struct {
		name          string
		objectiveType api_v1_beta1.ObjectiveType
		expected      goptuna.StudyDirection
	}{
		{
			name:          "minimize",
			objectiveType: api_v1_beta1.ObjectiveType_MINIMIZE,
			expected:      goptuna.StudyDirectionMinimize,
		},
		{
			name:          "maximize",
			objectiveType: api_v1_beta1.ObjectiveType_MAXIMIZE,
			expected:      goptuna.StudyDirectionMaximize,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := toGoptunaDirection(tt.objectiveType)
			if got != tt.expected {
				t.Errorf("toGoptunaDirection() got = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_toGoptunaSearchSpace(t *testing.T) {
	tests := []struct {
		name       string
		parameters []*api_v1_beta1.ParameterSpec
		want       map[string]interface{}
		wantErr    bool
	}{
		{
			name: "Double parameter type",
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
			want: map[string]interface{}{
				"param-double": goptuna.UniformDistribution{
					High: 5.5,
					Low:  1.5,
				},
			},
			wantErr: false,
		},
		{
			name: "Double parameter type with step",
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
			want: map[string]interface{}{
				"param-double": goptuna.DiscreteUniformDistribution{
					High: 5.5,
					Low:  1.5,
					Q:    0.5,
				},
			},
			wantErr: false,
		},
		{
			name: "Int parameter type",
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
			want: map[string]interface{}{
				"param-int": goptuna.IntUniformDistribution{
					High: 5,
					Low:  1,
				},
			},
			wantErr: false,
		},
		{
			name: "Int parameter type with step",
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
			want: map[string]interface{}{
				"param-int": goptuna.StepIntUniformDistribution{
					High: 5,
					Low:  1,
					Step: 2,
				},
			},
			wantErr: false,
		},
		{
			name: "Discrete parameter type",
			parameters: []*api_v1_beta1.ParameterSpec{
				{
					Name:          "param-discrete",
					ParameterType: api_v1_beta1.ParameterType_DISCRETE,
					FeasibleSpace: &api_v1_beta1.FeasibleSpace{
						List: []string{"3", "2", "6"},
					},
				},
			},
			want: map[string]interface{}{
				"param-discrete": goptuna.CategoricalDistribution{
					Choices: []string{"3", "2", "6"},
				},
			},
			wantErr: false,
		},
		{
			name: "Categorical parameter type",
			parameters: []*api_v1_beta1.ParameterSpec{
				{
					Name:          "param-categorical",
					ParameterType: api_v1_beta1.ParameterType_CATEGORICAL,
					FeasibleSpace: &api_v1_beta1.FeasibleSpace{
						List: []string{"cat1", "cat2", "cat3"},
					},
				},
			},
			want: map[string]interface{}{
				"param-categorical": goptuna.CategoricalDistribution{
					Choices: []string{"cat1", "cat2", "cat3"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toGoptunaSearchSpace(tt.parameters)
			if (err != nil) != tt.wantErr {
				t.Errorf("toGoptunaSearchSpace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toGoptunaSearchSpace() got = %v, want %v", got, tt.want)
			}
		})
	}
}
