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
	"fmt"
	"reflect"
	"strconv"

	"github.com/c-bata/goptuna"
	api_v1_beta1 "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
)

func sampleNextParam(study *goptuna.Study, searchSpace map[string]interface{}) (int, []*api_v1_beta1.ParameterAssignment, error) {
	nextTrialID, err := study.Storage.CreateNewTrial(study.ID)
	if err != nil {
		return -1, nil, err
	}
	trial := goptuna.Trial{
		Study: study,
		ID:    nextTrialID,
	}

	// Sample parameters and stored in Goptuna storage.
	err = trial.CallRelativeSampler()
	if err != nil {
		return nextTrialID, nil, err
	}

	assignments := make([]*api_v1_beta1.ParameterAssignment, 0, len(searchSpace))
	for name := range searchSpace {
		switch distribution := searchSpace[name].(type) {
		case goptuna.UniformDistribution:
			p, err := trial.SuggestFloat(name, distribution.Low, distribution.High)
			if err != nil {
				return nextTrialID, nil, err
			}
			assignments = append(assignments, &api_v1_beta1.ParameterAssignment{
				Name:  name,
				Value: strconv.FormatFloat(p, 'f', -1, 64),
			})
		case goptuna.DiscreteUniformDistribution:
			p, err := trial.SuggestDiscreteFloat(name, distribution.Low, distribution.High, distribution.Q)
			if err != nil {
				return nextTrialID, nil, err
			}
			assignments = append(assignments, &api_v1_beta1.ParameterAssignment{
				Name:  name,
				Value: strconv.FormatFloat(p, 'f', -1, 64),
			})
		case goptuna.IntUniformDistribution:
			p, err := trial.SuggestInt(name, distribution.Low, distribution.High)
			if err != nil {
				return nextTrialID, nil, err
			}
			assignments = append(assignments, &api_v1_beta1.ParameterAssignment{
				Name:  name,
				Value: strconv.Itoa(p),
			})
		case goptuna.StepIntUniformDistribution:
			p, err := trial.SuggestStepInt(name, distribution.Low, distribution.High, distribution.Step)
			if err != nil {
				return nextTrialID, nil, err
			}
			assignments = append(assignments, &api_v1_beta1.ParameterAssignment{
				Name:  name,
				Value: strconv.Itoa(p),
			})
		case goptuna.CategoricalDistribution:
			p, err := trial.SuggestCategorical(name, distribution.Choices)
			if err != nil {
				return nextTrialID, nil, err
			}
			assignments = append(assignments, &api_v1_beta1.ParameterAssignment{
				Name:  name,
				Value: p,
			})
		}
	}
	return nextTrialID, assignments, nil
}

func findGoptunaTrialIDByParam(study *goptuna.Study, trialMapping map[string]int, ktrial goptuna.FrozenTrial) (int, error) {
	trials, err := study.GetTrials()
	if err != nil {
		return -1, err
	}

	existInMapping := func(trialID int) bool {
		for j := range trialMapping {
			if trialMapping[j] == trialID {
				return true
			}
		}
		return false
	}

	for i := len(trials) - 1; i >= 0; i-- {
		if trials[i].State != goptuna.TrialStateRunning {
			continue
		}

		// skip the trial id which already mapped by other katib trial name
		if existInMapping(trials[i].ID) {
			continue
		}

		if reflect.DeepEqual(ktrial.Params, trials[i].Params) {
			return trials[i].ID, nil
		}
	}
	return -1, fmt.Errorf("Same parameter is not found for Trial: %v", ktrial)
}
