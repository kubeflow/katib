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
	"strconv"
	"time"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/cmaes"
	"github.com/c-bata/goptuna/sobol"
	"github.com/c-bata/goptuna/tpe"
	api_v1_beta1 "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
)

func toGoptunaDirection(t api_v1_beta1.ObjectiveType) goptuna.StudyDirection {
	if t == api_v1_beta1.ObjectiveType_MINIMIZE {
		return goptuna.StudyDirectionMinimize
	}
	return goptuna.StudyDirectionMaximize
}

func toGoptunaSampler(algorithm *api_v1_beta1.AlgorithmSpec) (goptuna.Sampler, goptuna.RelativeSampler, error) {
	name := algorithm.GetAlgorithmName()
	if name == AlgorithmCMAES {
		opts := make([]cmaes.SamplerOption, 0, len(algorithm.GetAlgorithmSettings())+1)
		opts = append(opts, cmaes.SamplerOptionNStartupTrials(0))
		for _, s := range algorithm.GetAlgorithmSettings() {
			if s.Name == "random_state" {
				seed, err := strconv.Atoi(s.Value)
				if err != nil {
					return nil, nil, err
				}
				opts = append(opts, cmaes.SamplerOptionSeed(int64(seed)))
			} else if s.Name == "sigma" {
				sigma, err := strconv.ParseFloat(s.Value, 64)
				if err != nil {
					return nil, nil, err
				}
				opts = append(opts, cmaes.SamplerOptionInitialSigma(sigma))
			} else if s.Name == "restart_strategy" {
				if s.Value == "ipop" {
					// The argument is multiplier of population size before each restart and basically 2 is recommended.
					// According to the paper, it reveal similar performance for factors between 2 and 3.
					opts = append(opts, cmaes.SamplerOptionIPop(2))
				} else if s.Value == "bipop" {
					opts = append(opts, cmaes.SamplerOptionBIPop(2))
				} else if s.Value != "none" {
					return nil, nil, fmt.Errorf("invalid restart_strategy: '%s'", s.Value)
				}
			}
		}
		return nil, cmaes.NewSampler(opts...), nil
	} else if name == AlgorithmTPE {
		opts := make([]tpe.SamplerOption, 0, len(algorithm.GetAlgorithmSettings()))
		for _, s := range algorithm.GetAlgorithmSettings() {
			if s.Name == "random_state" {
				seed, err := strconv.Atoi(s.Value)
				if err != nil {
					return nil, nil, err
				}
				opts = append(opts, tpe.SamplerOptionSeed(int64(seed)))
			} else if s.Name == "n_startup_trials" {
				n, err := strconv.Atoi(s.Value)
				if err != nil {
					return nil, nil, err
				}
				opts = append(opts, tpe.SamplerOptionNumberOfStartupTrials(n))
			} else if s.Name == "n_ei_candidates" {
				n, err := strconv.Atoi(s.Value)
				if err != nil {
					return nil, nil, err
				}
				opts = append(opts, tpe.SamplerOptionNumberOfEICandidates(n))
			}
		}
		return tpe.NewSampler(opts...), nil, nil
	} else if name == AlgorithmSobol {
		return nil, sobol.NewSampler(), nil
	} else {
		opts := make([]goptuna.RandomSamplerOption, 0, len(algorithm.GetAlgorithmSettings()))
		for _, s := range algorithm.GetAlgorithmSettings() {
			if s.Name == "random_state" {
				seed, err := strconv.Atoi(s.Value)
				if err != nil {
					return nil, nil, err
				}
				opts = append(opts, goptuna.RandomSamplerOptionSeed(int64(seed)))
			}
		}
		return goptuna.NewRandomSampler(opts...), nil, nil
	}
}

func toGoptunaSearchSpace(parameters []*api_v1_beta1.ParameterSpec) (map[string]interface{}, error) {
	searchSpace := make(map[string]interface{}, len(parameters))
	for _, p := range parameters {
		if p.ParameterType == api_v1_beta1.ParameterType_DOUBLE {
			high, err := strconv.ParseFloat(p.GetFeasibleSpace().GetMax(), 64)
			if err != nil {
				return nil, err
			}
			low, err := strconv.ParseFloat(p.GetFeasibleSpace().GetMin(), 64)
			if err != nil {
				return nil, err
			}

			stepstr := p.GetFeasibleSpace().GetStep()
			if stepstr == "" {
				searchSpace[p.Name] = goptuna.UniformDistribution{
					High: high,
					Low:  low,
				}
			} else {
				step, err := strconv.ParseFloat(stepstr, 64)
				if err != nil {
					return nil, err
				}
				searchSpace[p.Name] = goptuna.DiscreteUniformDistribution{
					High: high,
					Low:  low,
					Q:    step,
				}
			}
		} else if p.ParameterType == api_v1_beta1.ParameterType_INT {
			high, err := strconv.Atoi(p.GetFeasibleSpace().GetMax())
			if err != nil {
				return nil, err
			}
			low, err := strconv.Atoi(p.GetFeasibleSpace().GetMin())
			if err != nil {
				return nil, err
			}
			stepstr := p.GetFeasibleSpace().GetStep()
			if stepstr == "" {
				searchSpace[p.Name] = goptuna.IntUniformDistribution{
					High: high,
					Low:  low,
				}
			} else {
				step, err := strconv.Atoi(stepstr)
				if err != nil {
					return nil, err
				}
				searchSpace[p.Name] = goptuna.StepIntUniformDistribution{
					High: high,
					Low:  low,
					Step: step,
				}
			}
		} else if p.ParameterType == api_v1_beta1.ParameterType_CATEGORICAL {
			choices := p.GetFeasibleSpace().GetList()
			searchSpace[p.Name] = goptuna.CategoricalDistribution{
				Choices: choices,
			}
		} else if p.ParameterType == api_v1_beta1.ParameterType_DISCRETE {
			// Use categorical distribution instead of goptuna.DiscreteUniformDistribution
			// because goptuna.DiscreteUniformDistributions needs to declare the parameter
			// space with minimum value, maximum value and interval.
			choices := p.GetFeasibleSpace().GetList()
			searchSpace[p.Name] = goptuna.CategoricalDistribution{
				Choices: choices,
			}
		} else {
			return nil, fmt.Errorf("Unsupported parameter type: %v", p.ParameterType)
		}
	}
	return searchSpace, nil
}

func toGoptunaState(condition api_v1_beta1.TrialStatus_TrialConditionType) (goptuna.TrialState, error) {
	if condition == api_v1_beta1.TrialStatus_CREATED {
		return goptuna.TrialStateRunning, nil
	} else if condition == api_v1_beta1.TrialStatus_RUNNING {
		return goptuna.TrialStateRunning, nil
	} else if condition == api_v1_beta1.TrialStatus_SUCCEEDED {
		return goptuna.TrialStateComplete, nil
	} else if condition == api_v1_beta1.TrialStatus_FAILED {
		return goptuna.TrialStateFail, nil
	} else if condition == api_v1_beta1.TrialStatus_EARLYSTOPPED {
		return goptuna.TrialStatePruned, nil
	}
	return goptuna.TrialStateFail, fmt.Errorf("Unexpected Trial condition: %v", condition)
}

func getFinalMetric(objectMetricName string, trial *api_v1_beta1.Trial) (float64, error) {
	metrics := trial.GetStatus().GetObservation().GetMetrics()
	for i := len(metrics) - 1; i >= 0; i-- {
		if metrics[i].GetName() != objectMetricName {
			continue
		}
		v, err := strconv.ParseFloat(metrics[i].GetValue(), 64)
		if err != nil {
			return 0, err
		}
		return v, nil
	}
	return 0, fmt.Errorf("No objective metric in Trial %v", trial)
}

func toGoptunaTrials(
	ktrials []*api_v1_beta1.Trial,
	objectMetricName string,
	study *goptuna.Study,
	searchSpace map[string]interface{},
) (map[string]goptuna.FrozenTrial, error) {
	gtrials := make(map[string]goptuna.FrozenTrial, len(ktrials))
	for i, kt := range ktrials {
		var err error
		var datetimeStart, datetimeComplete time.Time
		if kt.GetStatus().GetStartTime() != "" {
			datetimeStart, err = time.Parse(time.RFC3339Nano, kt.GetStatus().GetStartTime())
			if err != nil {
				return nil, err
			}
		}
		if kt.GetStatus().GetCompletionTime() != "" {
			datetimeComplete, err = time.Parse(time.RFC3339Nano, kt.GetStatus().GetCompletionTime())
			if err != nil {
				return nil, err
			}
		}
		state, err := toGoptunaState(kt.GetStatus().GetCondition())
		if err != nil {
			return nil, err
		}

		var finalValue float64
		if state == goptuna.TrialStateComplete {
			finalValue, err = getFinalMetric(objectMetricName, kt)
			if err != nil {
				return nil, err
			}
		}

		assignments := kt.GetSpec().GetParameterAssignments().GetAssignments()
		internalParams, externalParams, err := toGoptunaParams(assignments, searchSpace)
		if err != nil {
			return nil, err
		}

		gt := goptuna.FrozenTrial{
			ID:                 i, // dummy id
			StudyID:            study.ID,
			Number:             i, // dummy number
			State:              state,
			Value:              finalValue,
			IntermediateValues: nil,
			DatetimeStart:      datetimeStart,
			DatetimeComplete:   datetimeComplete,
			InternalParams:     internalParams,
			Params:             externalParams,
			Distributions:      searchSpace,
			UserAttrs:          nil,
			SystemAttrs:        nil,
		}
		gtrials[kt.GetName()] = gt
	}
	return gtrials, nil
}

func toGoptunaParams(
	assignments []*api_v1_beta1.ParameterAssignment,
	searchSpace map[string]interface{},
) (
	internalParams map[string]float64,
	externalParams map[string]interface{},
	err error,
) {
	internalParams = make(map[string]float64, len(assignments))
	externalParams = make(map[string]interface{}, len(assignments))

	for i := range assignments {
		name := assignments[i].GetName()
		valueStr := assignments[i].GetValue()

		switch d := searchSpace[name].(type) {
		case goptuna.UniformDistribution:
			p, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				return nil, nil, err
			}
			internalParams[name] = p
			externalParams[name] = d.ToExternalRepr(p)
		case goptuna.DiscreteUniformDistribution:
			p, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				return nil, nil, err
			}
			internalParams[name] = p
			externalParams[name] = d.ToExternalRepr(p)
		case goptuna.IntUniformDistribution:
			p, err := strconv.ParseInt(valueStr, 10, 64)
			if err != nil {
				return nil, nil, err
			}
			ir := float64(p)
			internalParams[name] = ir
			// externalParams[name] = p is prohibited because of reflect.DeepEqual() will be false
			// at findGoptunaTrialIDByParam() function.
			externalParams[name] = d.ToExternalRepr(ir)
		case goptuna.StepIntUniformDistribution:
			p, err := strconv.ParseInt(valueStr, 10, 64)
			if err != nil {
				return nil, nil, err
			}
			ir := float64(p)
			internalParams[name] = ir
			// externalParams[name] = p is prohibited because of reflect.DeepEqual() will be false
			// at findGoptunaTrialIDByParam() function.
			externalParams[name] = d.ToExternalRepr(ir)
		case goptuna.CategoricalDistribution:
			internalRepr := -1.0
			for i := range d.Choices {
				if d.Choices[i] == valueStr {
					internalRepr = float64(i)
					break
				}
			}
			if internalRepr == -1.0 {
				return nil, nil, fmt.Errorf("Invalid categorical value: %v", internalRepr)
			}
			internalParams[name] = internalRepr
			externalParams[name] = valueStr
		}
	}
	return internalParams, externalParams, nil
}

func createStudyAndSearchSpace(
	experiment *api_v1_beta1.Experiment,
) (*goptuna.Study, map[string]interface{}, error) {
	direction := toGoptunaDirection(experiment.GetSpec().GetObjective().GetType())
	independentSampler, relativeSampler, err := toGoptunaSampler(experiment.GetSpec().GetAlgorithm())
	if err != nil {
		return nil, nil, err
	}
	searchSpace, err := toGoptunaSearchSpace(experiment.GetSpec().GetParameterSpecs().GetParameters())
	if err != nil {
		return nil, nil, err
	}

	studyOpts := make([]goptuna.StudyOption, 0, 5)
	studyOpts = append(studyOpts, goptuna.StudyOptionDirection(direction))
	studyOpts = append(studyOpts, goptuna.StudyOptionDefineSearchSpace(searchSpace))
	studyOpts = append(studyOpts, goptuna.StudyOptionLogger(nil))
	if independentSampler != nil {
		studyOpts = append(studyOpts, goptuna.StudyOptionSampler(independentSampler))
	}
	if relativeSampler != nil {
		studyOpts = append(studyOpts, goptuna.StudyOptionRelativeSampler(relativeSampler))
	}

	study, err := goptuna.CreateStudy(defaultStudyName, studyOpts...)
	if err != nil {
		return nil, nil, err
	}

	return study, searchSpace, nil
}
