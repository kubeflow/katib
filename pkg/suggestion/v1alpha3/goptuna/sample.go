package suggestion_goptuna_v1alpha3

import (
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
	nextTrial, err := study.Storage.GetTrial(nextTrialID)
	if err != nil {
		klog.Errorf("Failed to get a next trial: %s", err)
		return nil, err
	}

	var relativeSampleParams map[string]float64
	if study.RelativeSampler != nil {
		relativeSampleParams, err = study.RelativeSampler.SampleRelative(study, nextTrial, searchSpace)
		if err != nil {
			klog.Errorf("Failed to call SampleRelative: %s", err)
			return nil, err
		}
	}

	assignments := make([]*api_v1_alpha3.ParameterAssignment, 0, len(searchSpace))
	trial := goptuna.Trial{
		Study: study,
		ID:    nextTrialID,
	}

	for name := range searchSpace {
		switch distribution := searchSpace[name].(type) {
		case goptuna.UniformDistribution:
			var p float64
			if internalParam, ok := relativeSampleParams[name]; ok {
				p = internalParam
			} else {
				p, err = trial.SuggestUniform(name, distribution.Low, distribution.High)
				if err != nil {
					klog.Errorf("Failed to get suggested param: %s", err)
					return nil, err
				}
			}
			assignments = append(assignments, &api_v1_alpha3.ParameterAssignment{
				Name:  name,
				Value: strconv.FormatFloat(p, 'f', -1, 64),
			})
		case goptuna.IntUniformDistribution:
			var p int
			if internalParam, ok := relativeSampleParams[name]; ok {
				p = int(internalParam)
			} else {
				p, err = trial.SuggestInt(name, distribution.Low, distribution.High)
				if err != nil {
					klog.Errorf("Failed to get suggested param: %s", err)
					return nil, err
				}
			}
			assignments = append(assignments, &api_v1_alpha3.ParameterAssignment{
				Name:  name,
				Value: strconv.Itoa(p),
			})
		case goptuna.DiscreteUniformDistribution:
			var p float64
			if internalParam, ok := relativeSampleParams[name]; ok {
				p = internalParam
			} else {
				p, err = trial.SuggestDiscreteUniform(name, distribution.Low, distribution.High, distribution.Q)
				if err != nil {
					klog.Errorf("Failed to get suggested param: %s", err)
					return nil, err
				}
			}
			assignments = append(assignments, &api_v1_alpha3.ParameterAssignment{
				Name:  name,
				Value: strconv.FormatFloat(p, 'f', -1, 64),
			})
		case goptuna.CategoricalDistribution:
			var p string
			if internalParam, ok := relativeSampleParams[name]; ok {
				p = distribution.Choices[int(internalParam)]
			} else {
				p, err = trial.SuggestCategorical(name, distribution.Choices)
				if err != nil {
					klog.Errorf("Failed to get suggested param: %s", err)
					return nil, err
				}
			}
			assignments = append(assignments, &api_v1_alpha3.ParameterAssignment{
				Name:  name,
				Value: p,
			})
		}
	}
	return assignments, nil
}
