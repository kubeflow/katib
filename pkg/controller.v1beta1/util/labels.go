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

package util

import (
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

// SuggestionLabels returns the expected suggestion labels.
func SuggestionLabels(instance *suggestionsv1beta1.Suggestion) map[string]string {
	res := make(map[string]string)
	for k, v := range instance.Labels {
		res[k] = v
	}
	res[consts.LabelDeploymentName] = GetSuggestionDeploymentName(instance)
	res[consts.LabelExperimentName] = instance.Name
	res[consts.LabelSuggestionName] = instance.Name

	return res
}

// TrialLabels returns the expected trial labels.
func TrialLabels(instance *experimentsv1beta1.Experiment) map[string]string {
	res := make(map[string]string)
	for k, v := range instance.Labels {
		res[k] = v
	}
	res[consts.LabelExperimentName] = instance.Name

	return res
}
