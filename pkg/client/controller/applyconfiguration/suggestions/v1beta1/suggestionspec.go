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
	v1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
)

// SuggestionSpecApplyConfiguration represents a declarative configuration of the SuggestionSpec type for use
// with apply.
type SuggestionSpecApplyConfiguration struct {
	Algorithm     *v1beta1.AlgorithmSpec               `json:"algorithm,omitempty"`
	EarlyStopping *v1beta1.EarlyStoppingSpec           `json:"earlyStopping,omitempty"`
	Requests      *int32                               `json:"requests,omitempty"`
	ResumePolicy  *experimentsv1beta1.ResumePolicyType `json:"resumePolicy,omitempty"`
}

// SuggestionSpecApplyConfiguration constructs a declarative configuration of the SuggestionSpec type for use with
// apply.
func SuggestionSpec() *SuggestionSpecApplyConfiguration {
	return &SuggestionSpecApplyConfiguration{}
}

// WithAlgorithm sets the Algorithm field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Algorithm field is set to the value of the last call.
func (b *SuggestionSpecApplyConfiguration) WithAlgorithm(value v1beta1.AlgorithmSpec) *SuggestionSpecApplyConfiguration {
	b.Algorithm = &value
	return b
}

// WithEarlyStopping sets the EarlyStopping field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the EarlyStopping field is set to the value of the last call.
func (b *SuggestionSpecApplyConfiguration) WithEarlyStopping(value v1beta1.EarlyStoppingSpec) *SuggestionSpecApplyConfiguration {
	b.EarlyStopping = &value
	return b
}

// WithRequests sets the Requests field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Requests field is set to the value of the last call.
func (b *SuggestionSpecApplyConfiguration) WithRequests(value int32) *SuggestionSpecApplyConfiguration {
	b.Requests = &value
	return b
}

// WithResumePolicy sets the ResumePolicy field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ResumePolicy field is set to the value of the last call.
func (b *SuggestionSpecApplyConfiguration) WithResumePolicy(value experimentsv1beta1.ResumePolicyType) *SuggestionSpecApplyConfiguration {
	b.ResumePolicy = &value
	return b
}
