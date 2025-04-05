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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1beta1

import (
	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
)

// TrialAssignmentApplyConfiguration represents a declarative configuration of the TrialAssignment type for use
// with apply.
type TrialAssignmentApplyConfiguration struct {
	ParameterAssignments []commonv1beta1.ParameterAssignment `json:"parameterAssignments,omitempty"`
	Name                 *string                             `json:"name,omitempty"`
	EarlyStoppingRules   []commonv1beta1.EarlyStoppingRule   `json:"earlyStoppingRules,omitempty"`
	Labels               map[string]string                   `json:"labels,omitempty"`
}

// TrialAssignmentApplyConfiguration constructs a declarative configuration of the TrialAssignment type for use with
// apply.
func TrialAssignment() *TrialAssignmentApplyConfiguration {
	return &TrialAssignmentApplyConfiguration{}
}

// WithParameterAssignments adds the given value to the ParameterAssignments field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the ParameterAssignments field.
func (b *TrialAssignmentApplyConfiguration) WithParameterAssignments(values ...commonv1beta1.ParameterAssignment) *TrialAssignmentApplyConfiguration {
	for i := range values {
		b.ParameterAssignments = append(b.ParameterAssignments, values[i])
	}
	return b
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *TrialAssignmentApplyConfiguration) WithName(value string) *TrialAssignmentApplyConfiguration {
	b.Name = &value
	return b
}

// WithEarlyStoppingRules adds the given value to the EarlyStoppingRules field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the EarlyStoppingRules field.
func (b *TrialAssignmentApplyConfiguration) WithEarlyStoppingRules(values ...commonv1beta1.EarlyStoppingRule) *TrialAssignmentApplyConfiguration {
	for i := range values {
		b.EarlyStoppingRules = append(b.EarlyStoppingRules, values[i])
	}
	return b
}

// WithLabels puts the entries into the Labels field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the Labels field,
// overwriting an existing map entries in Labels field with the same key.
func (b *TrialAssignmentApplyConfiguration) WithLabels(entries map[string]string) *TrialAssignmentApplyConfiguration {
	if b.Labels == nil && len(entries) > 0 {
		b.Labels = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.Labels[k] = v
	}
	return b
}
