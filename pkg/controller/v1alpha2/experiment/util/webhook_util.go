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
	"bytes"
	"fmt"

	ep_v1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

func DefaultExperiment(instance *ep_v1alpha2.Experiment) error {
	instance.SetDefault()
	return nil
}

func ValidateExperiment(instance *ep_v1alpha2.Experiment) error {
	if err := validateObjective(instance.Spec.Objective); err != nil {
		return err
	}
	if err := validateAlgorithm(instance.Spec.Algorithm); err != nil {
		return err
	}

	if err := validateTrialTemplate(instance); err != nil {
		return err
	}

	if len(instance.Spec.Parameters) == 0 && instance.Spec.NasConfig == nil {
		return fmt.Errorf("spec.parameters or spec.nasConfig must be specified.")
	}

	if len(instance.Spec.Parameters) > 0 && instance.Spec.NasConfig != nil {
		return fmt.Errorf("Only one of spec.parameters and spec.nasConfig can be specified.")
	}

	if err := validateAlgorithmSettings(instance); err !=nil {
		return err
	}
	return nil
}

func validateAlgorithmSettings(inst *ep_v1alpha2.Experiment) error {
	// TODO: it need call ValidateAlgorithmSettings API of vizier-core manager, implement it when vizier-core done
	return nil
}

func validateObjective(obj *ep_v1alpha2.ObjectiveSpec) error {
	if obj == nil {
		return fmt.Errorf("No spec.objective specified.")
	}
	if obj.Type != ep_v1alpha2.OptimizationTypeMinimize && obj.Type != ep_v1alpha2.OptimizationTypeMaximize  {
		return fmt.Errorf("spec.objective.type must be %s or %s.", ep_v1alpha2.OptimizationTypeMinimize, ep_v1alpha2.OptimizationTypeMaximize)
	}
	if obj.ObjectiveMetricName == "" {
		return fmt.Errorf("No spec.objective.objectiveMetricName specified.")
	}
	return nil
}

func validateAlgorithm(ag *ep_v1alpha2.AlgorithmSpec) error {
	if ag == nil {
		return fmt.Errorf("No spec.algorithm specified.")
	}
	if ag.AlgorithmName == "" {
		return fmt.Errorf("No spec.algorithm.name specified.")
	}

	return nil
}

func validateTrialTemplate(instance *ep_v1alpha2.Experiment) error {
	trialName := fmt.Sprintf("%s-trial", instance.GetName())
	trialParams := TrialTemplateParams{
		Experiment: instance.GetName(),
		Trial:      trialName,
		NameSpace:  instance.GetNamespace(),
	}
	runSpec, err := GetRunSpec(instance, trialParams)
	if err != nil {
		return fmt.Errorf("Invalid spec.trialTemplate: %v.", err)
	}

	bufSize := 1024
	buf := bytes.NewBufferString(runSpec)

	job := &unstructured.Unstructured{}
	if err := k8syaml.NewYAMLOrJSONDecoder(buf, bufSize).Decode(job); err != nil {
		return fmt.Errorf("Invalid spec.trialTemplate: %v.", err)
	}

	if err := validateSupportedJob(job); err != nil {
		return fmt.Errorf("Invalid spec.trialTemplate: %v.", err)
	}

	if job.GetNamespace() != instance.GetNamespace() {
		return fmt.Errorf("Invalid spec.trialTemplate: metadata.namespace should be %s or {{.NameSpace}}", instance.GetNamespace())
	}
	if job.GetName() != trialName {
		return fmt.Errorf("Invalid spec.trialTemplate: metadata.name should be {{.Trial}}")
	}
	return nil
}

func validateSupportedJob(job *unstructured.Unstructured) error {
	gvk := job.GroupVersionKind()
	supportedJobs := GetSupportdJobList()
	for _, sJob := range supportedJobs {
		if gvk == sJob {
			return nil
		}
	}
	return fmt.Errorf("Cannot support to run job: %v", gvk)
}
