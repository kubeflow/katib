/*
Copyright 2019 The Kubernetes Authors.

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

package util

import (
	"fmt"

	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
)

func DefaultExperiment(instance *experimentsv1alpha2.Experiment) error {
	instance.Default()
	return nil
}

func ValidateExperiment(instance *experimentsv1alpha2.Experiment) error {
	if err := validateObjective(instance.Spec.Objective); err != nil {
		return nil
	}

	if err := validateAlgorithm(instance.Spec.Algorithm); err != nil {
		return nil
	}

	if err := validateTrialTemplate(instance, instance.Spec.TrialTemplate); err != nil {
		return nil
	}

	if len(instance.Spec.Parameters) == 0 && instance.Spec.NasConfig == nil {
		return fmt.Errorf("spec.Parameters or spec.NasConfig must be specified.")
	}

	if len(instance.Spec.Parameters) > 0 && instance.Spec.NasConfig != nil {
		return fmt.Errorf("Only one of spec.Parameters and spec.NasConfig can be specified.")
	}

	if len(instance.Spec.Parameters) > 0 {
		if err := validateParameters(instance.Spec.Parameters); err != nil {
			return nil;
		}
	}

	if instance.Spec.NasConfig != nil {
		if err := validateNasConfig(instance.Spec.NasConfig); err != nil {
			return nil;
		}
	}

	return nil
}

func validateObjective(oj *experimentsv1alpha2.ObjectiveSpec) error {
	return nil
}

func validateAlgorithm(ag *experimentsv1alpha2.AlgorithmSpec) error {
	return nil
}

func validateTrialTemplate(instance *experimentsv1alpha2.Experiment, t *experimentsv1alpha2.TrialTemplate) error {
	if t == nil {
		return fmt.Errorf("No spec.trialTemplate specified.")
	}

	return nil
}

func validateParameters(ps []experimentsv1alpha2.ParameterSpec) error {
	return nil
}

func validateNasConfig(nas *experimentsv1alpha2.NasConfig) error {
	return nil
}
