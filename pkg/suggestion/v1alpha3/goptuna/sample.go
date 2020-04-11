package suggestion_goptuna_v1alpha3

import (
	"errors"
	"math"
	"strconv"

	"github.com/c-bata/goptuna"
	api_v1_alpha3 "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
	"k8s.io/klog"
)

func sampleNextParam(study *goptuna.Study, searchSpace map[string]interface{}) ([]*api_v1_alpha3.ParameterAssignment, error) {
	nextTrialID, err := study.Storage.CreateNewTrial(study.ID)
	if err != nil {
		klog.Errorf("Failed to create a new trial: %s", err)
		return nil, err
	}
	trial := goptuna.Trial{
		Study: study,
		ID:    nextTrialID,
	}

	err = trial.CallRelativeSampler()
	if err != nil {
		klog.Errorf("Failed to sample relative parameters: %s", err)
		return nil, err
	}

	assignments := make([]*api_v1_alpha3.ParameterAssignment, 0, len(searchSpace))
	for name := range searchSpace {
		switch distribution := searchSpace[name].(type) {
		case goptuna.UniformDistribution:
			p, err := trial.SuggestFloat(name, distribution.Low, distribution.High)
			if err != nil {
				klog.Errorf("Failed to get suggested param: %s", err)
				return nil, err
			}
			assignments = append(assignments, &api_v1_alpha3.ParameterAssignment{
				Name:  name,
				Value: strconv.FormatFloat(p, 'f', -1, 64),
			})
		case goptuna.DiscreteUniformDistribution:
			p, err := trial.SuggestDiscreteFloat(name, distribution.Low, distribution.High, distribution.Q)
			if err != nil {
				klog.Errorf("Failed to get suggested param: %s", err)
				return nil, err
			}
			assignments = append(assignments, &api_v1_alpha3.ParameterAssignment{
				Name:  name,
				Value: strconv.FormatFloat(p, 'f', -1, 64),
			})
		case goptuna.IntUniformDistribution:
			p, err := trial.SuggestInt(name, distribution.Low, distribution.High)
			if err != nil {
				klog.Errorf("Failed to get suggested param: %s", err)
				return nil, err
			}
			assignments = append(assignments, &api_v1_alpha3.ParameterAssignment{
				Name:  name,
				Value: strconv.Itoa(p),
			})
		case goptuna.StepIntUniformDistribution:
			p, err := trial.SuggestStepInt(name, distribution.Low, distribution.High, distribution.Step)
			if err != nil {
				klog.Errorf("Failed to get suggested param: %s", err)
				return nil, err
			}
			assignments = append(assignments, &api_v1_alpha3.ParameterAssignment{
				Name:  name,
				Value: strconv.Itoa(p),
			})
		case goptuna.CategoricalDistribution:
			p, err := trial.SuggestCategorical(name, distribution.Choices)
			if err != nil {
				klog.Errorf("Failed to get suggested param: %s", err)
				return nil, err
			}
			assignments = append(assignments, &api_v1_alpha3.ParameterAssignment{
				Name:  name,
				Value: p,
			})
		}
	}

	if t, err := study.Storage.GetTrial(trial.ID); err == nil {
		klog.Infof("Success to sample new trial: trialID=%d, params=%v", t.ID, t.Params)
	}
	return assignments, nil
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

	minManhattan := math.MaxFloat64
	estimatedTrialID := -1
	for i := range trials {
		if trials[i].State != goptuna.TrialStateRunning {
			continue
		}

		// skip the trial id which already mapped by other katib trial name
		if existInMapping(trials[i].ID) {
			continue
		}

		// To understand how this function estimate the Goptuna trial ID from the parameters,
		// you need to understand the 'internal representation' in Goptuna.
		// Goptuna trials holds the parameters in two types of representations.
		// To explain the representation format, please see the following example.
		//
		// * Search space: map[string]interface{}{"x1": Uniform{Min: -10, Max: 10}, "x2": Categorical{Choices: []string{"param-1", "param-2", "param-3"}}}
		// * External representation: map[string]interface{}{"x1": 5.5, "x2": "param-2"}
		// * Internal representation: map[string]float64{"x1": 5.5, "x2": 1.0}
		//
		// In the internal representation, all parameters are represented by `float64` to store the storage
		// (because Goptuna supports not only in-memory but also RDB storage backend).
		// To represent categorical parameters, Goptuna holds an index of the list in the database.
		//
		// This function calculates Manhattan distance of internal representation parameters.
		// Then returns trialID which has the most 'similar' parameters.
		var manhattan float64
		for name := range trials[i].InternalParams {
			gtrialParamValue := trials[i].InternalParams[name]
			ktrialParamValue, ok := ktrial.InternalParams[name]
			if !ok {
				return -1, errors.New("must not reach here")
			}
			manhattan += math.Abs(gtrialParamValue - ktrialParamValue)
		}

		if manhattan < minManhattan {
			minManhattan = manhattan
			estimatedTrialID = trials[i].ID
		}
	}
	if estimatedTrialID == -1 {
		return -1, errors.New("goptuna trial is not found")
	}
	return estimatedTrialID, nil
}
