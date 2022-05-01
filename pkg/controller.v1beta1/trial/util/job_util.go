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
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/util"
)

// TrialJobStatus is the internal representation for deployed Job status
type TrialJobStatus struct {
	// Condition describes the state of the Job at a certain point.
	// Can be either Running, Succeeded or Failed
	Condition ConditionType `json:"condition,omitempty"`

	// The reason received from Job's status, if it is possible
	Reason string `json:"reason,omitempty"`

	// The message received from Job's status, if it is possible
	Message string `json:"message,omitempty"`
}

// ConditionType describes the various conditions a Job can be in.
type ConditionType string

const (
	// JobRunning means that Job was deployed by Trial.
	// Job doesn't have succeeded or failed condition.
	JobRunning ConditionType = "Running"

	// JobSucceeded means that Job status satisfies Trial success condition
	JobSucceeded ConditionType = "Succeeded"

	// JobFailed means that Job status satisfies Trial failure condition
	JobFailed ConditionType = "Failed"
)

// GetDeployedJobStatus returns internal representation for deployed Job status.
func GetDeployedJobStatus(trial *trialsv1beta1.Trial, deployedJob *unstructured.Unstructured) (*TrialJobStatus, error) {

	trialJobStatus := &TrialJobStatus{}

	// Marshal unstructured Job to JSON
	// Deployed Job is valid JSON
	deployedJobJSON, _ := util.ConvertUnstructuredToString(deployedJob)

	// Try to get failure condition using spec.failureCondition expression
	failureJobCondition := gjson.Get(deployedJobJSON, trial.Spec.FailureCondition)

	// Condition exists if failureJobCondition is object or failureJobCondition is array with len > 0
	if failureJobCondition.IsObject() || (failureJobCondition.IsArray() && len(failureJobCondition.Array()) > 0) {
		strCondition := failureJobCondition.String()

		// If failureJobCondition is array we take first element to unmarshal in TrialJobStatus
		if failureJobCondition.IsArray() {
			strCondition = failureJobCondition.Array()[0].String()
		}

		// Unmarshal condition to Trial Job representation to get message and reason if it exists
		err := json.Unmarshal([]byte(strCondition), &trialJobStatus)
		if err != nil {
			return nil, fmt.Errorf("unmarshal failure condition to Trial Job status failed %v", err)
		}

		// Job condition is failed
		trialJobStatus.Condition = JobFailed
		return trialJobStatus, nil
	}

	// Try to get success condition using spec.successCondition expression
	successJobCondition := gjson.Get(deployedJobJSON, trial.Spec.SuccessCondition)

	// The same logic as for failure condition
	if successJobCondition.IsObject() || (successJobCondition.IsArray() && len(successJobCondition.Array()) > 0) {
		strCondition := successJobCondition.String()

		if successJobCondition.IsArray() {
			strCondition = successJobCondition.Array()[0].String()
		}

		err := json.Unmarshal([]byte(strCondition), &trialJobStatus)
		if err != nil {
			return nil, fmt.Errorf("unmarshal success condition to Trial Job status failed %v", err)
		}

		// Job condition is succeeded
		trialJobStatus.Condition = JobSucceeded
		return trialJobStatus, nil
	}

	// Set default Job condition is running when Job name is generated.
	// Check if Trial is not running
	if !trial.IsRunning() && deployedJob.GetName() != "" {
		trialJobStatus.Condition = JobRunning
		return trialJobStatus, nil
	}

	// Otherwise returns nil object and Trial status doesn't need to be updated
	return nil, nil
}
