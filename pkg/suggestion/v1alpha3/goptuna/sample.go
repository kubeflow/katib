package suggestion_goptuna_v1alpha3

import (
	"errors"
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

func isSameTrialParam(ktrial, gtrial goptuna.FrozenTrial) bool {
	// Compare trial parameters by "internal representation".
	// In the internal representation, all parameters are represented by `float64` to store the storage
	// (because Goptuna supports not only in-memory but also RDB storage backend).
	// To represent categorical parameters, Goptuna holds an index of the list in the database.
	//
	// SearchSpace: map[string]interface{}{"x1": Uniform{Min: -10, Max: 10}, "x2": Categorical{Choices: []string{"param-1", "param-2"}}}
	// External representation: map[string]interface{}{"x1": 5.5, "x2": "param-2"}
	// Internal representation: map[string]float64{"x1": 5.5, "x2": 1.0}
	for name := range gtrial.InternalParams {
		gtrialParamValue := gtrial.InternalParams[name]
		ktrialParamValue, ok := ktrial.InternalParams[name]
		if !ok {
			// must not reach here
			klog.Errorf("Detect inconsistent internal parameters: %v and %v",
				ktrial.InternalParams, gtrial.InternalParams)
			return false
		}
		if gtrialParamValue != ktrialParamValue {
			return false
		}
	}
	return true
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

		if isSameTrialParam(ktrial, trials[i]) {
			return trials[i].ID, nil
		}
	}
	return -1, errors.New("same trial parameter is not found")
}
