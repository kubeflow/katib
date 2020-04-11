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

	minDiff := math.MaxFloat64
	estimatedTrialID := -1
	for i := range trials {
		if trials[i].State != goptuna.TrialStateRunning {
			continue
		}

		// skip the trial id which already mapped by other katib trial name
		if existInMapping(trials[i].ID) {
			continue
		}

		var diff float64
		for name := range trials[i].InternalParams {
			gtrialParamValue := trials[i].InternalParams[name]
			ktrialParamValue, ok := ktrial.InternalParams[name]
			if !ok {
				return -1, errors.New("must not reach here")
			}
			diff += math.Abs(gtrialParamValue - ktrialParamValue)
		}

		if diff < minDiff {
			minDiff = diff
			estimatedTrialID = trials[i].ID
		}
	}
	if estimatedTrialID == -1 {
		return -1, errors.New("goptuna trial is not found")
	}
	return estimatedTrialID, nil
}
