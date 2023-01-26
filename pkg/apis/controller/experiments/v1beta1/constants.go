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

package v1beta1

const (
	// DefaultTrialParallelCount is the default value of spec.parallelTrialCount.
	DefaultTrialParallelCount = 3

	// DefaultResumePolicy is the default value of spec.resumePolicy.
	DefaultResumePolicy = NeverResume

	// DefaultJobSuccessCondition is the default value of spec.trialTemplate.successCondition for Job.
	DefaultJobSuccessCondition = "status.conditions.#(type==\"Complete\")#|#(status==\"True\")#"

	// DefaultJobFailureCondition is the default value of spec.trialTemplate.failureCondition for Job.
	DefaultJobFailureCondition = "status.conditions.#(type==\"Failed\")#|#(status==\"True\")#"

	// DefaultKubeflowJobSuccessCondition is the default value of spec.trialTemplate.successCondition for Kubeflow Training Job.
	DefaultKubeflowJobSuccessCondition = "status.conditions.#(type==\"Succeeded\")#|#(status==\"True\")#"

	// DefaultKubeflowJobFailureCondition is the default value of spec.trialTemplate.failureCondition for Kubeflow Training Job.
	DefaultKubeflowJobFailureCondition = "status.conditions.#(type==\"Failed\")#|#(status==\"True\")#"
)

var (
	// DefaultKubeflowJobPrimaryPodLabels is the default value of spec.trialTemplate.primaryPodLabels for Kubeflow Training Job.
	DefaultKubeflowJobPrimaryPodLabels = map[string]string{"training.kubeflow.org/job-role": "master"}

	// KubeflowJobKinds is the list of Kubeflow Training Job kinds.
	KubeflowJobKinds = map[string]bool{
		"TFJob":      true,
		"PyTorchJob": true,
		"XGBoostJob": true,
		"MXJob":      true,
		"MPIJob":     true,
	}
)
