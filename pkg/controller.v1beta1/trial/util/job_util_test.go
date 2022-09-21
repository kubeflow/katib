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
	"reflect"
	"testing"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	trialsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/util"
)

const (
	testMessage = "test-message"
	testReason  = "test-reason"
)

func TestGetDeployedJobStatus(t *testing.T) {

	successCondition := "status.conditions.#(type==\"Complete\")#|#(status==\"True\")#"
	failureCondition := "status.conditions.#(type==\"Failed\")#|#(status==\"True\")#"

	tcs := []struct {
		trial                  *trialsv1beta1.Trial
		deployedJob            *unstructured.Unstructured
		expectedTrialJobStatus *TrialJobStatus
		err                    bool
		testDescription        string
	}{
		{
			trial: newFakeTrial(successCondition, failureCondition),
			deployedJob: func() *unstructured.Unstructured {
				job := newFakeJob()
				job.Status.Conditions[0].Status = corev1.ConditionFalse
				job.Status.Conditions[1].Status = corev1.ConditionFalse
				return newFakeDeployedJob(job)
			}(),
			expectedTrialJobStatus: func() *TrialJobStatus {
				return &TrialJobStatus{
					Condition: JobRunning,
				}
			}(),
			err:             false,
			testDescription: "Job status is running",
		},
		{
			trial:       newFakeTrial(successCondition, failureCondition),
			deployedJob: newFakeDeployedJob(newFakeJob()),
			expectedTrialJobStatus: func() *TrialJobStatus {
				return &TrialJobStatus{
					Condition: JobSucceeded,
					Message:   testMessage,
					Reason:    testReason,
				}
			}(),
			err:             false,
			testDescription: "Job status is succeeded, reason and message must be returned",
		},
		{
			trial: newFakeTrial(successCondition, failureCondition),
			deployedJob: func() *unstructured.Unstructured {
				job := newFakeJob()
				job.Status.Conditions[0].Status = corev1.ConditionTrue
				job.Status.Conditions[1].Status = corev1.ConditionFalse
				return newFakeDeployedJob(job)
			}(),
			expectedTrialJobStatus: func() *TrialJobStatus {
				return &TrialJobStatus{
					Condition: JobFailed,
					Message:   testMessage,
					Reason:    testReason,
				}
			}(),
			err:             false,
			testDescription: "Job status is failed, reason and message must be returned",
		},
		{
			trial:       newFakeTrial("status.[@this].#(succeeded==1)", failureCondition),
			deployedJob: newFakeDeployedJob(newFakeJob()),
			expectedTrialJobStatus: func() *TrialJobStatus {
				return &TrialJobStatus{
					Condition: JobSucceeded,
				}
			}(),
			err:             false,
			testDescription: "Job status is succeeded because status.succeeded = 1",
		},
	}

	for _, tc := range tcs {
		actualTrialJobStatus, err := GetDeployedJobStatus(tc.trial, tc.deployedJob)

		if tc.err && err == nil {
			t.Errorf("Case: %v failed. Expected err, got nil", tc.testDescription)
		} else if !tc.err {
			if err != nil {
				t.Errorf("Case: %v failed. Expected nil, got %v", tc.testDescription, err)
			} else if !reflect.DeepEqual(tc.expectedTrialJobStatus, actualTrialJobStatus) {
				t.Errorf("Case: %v failed. Expected %v\n got %v", tc.testDescription, tc.expectedTrialJobStatus, actualTrialJobStatus)
			}
		}
	}
}

func newFakeTrial(successCondition, failureCondition string) *trialsv1beta1.Trial {
	return &trialsv1beta1.Trial{
		Spec: trialsv1beta1.TrialSpec{
			SuccessCondition: successCondition,
			FailureCondition: failureCondition,
		},
	}
}

func newFakeJob() *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-job",
		},
		Status: batchv1.JobStatus{
			Conditions: []batchv1.JobCondition{
				{
					Type:    batchv1.JobFailed,
					Status:  corev1.ConditionFalse,
					Reason:  testReason,
					Message: testMessage,
				},
				{
					Type:    batchv1.JobComplete,
					Status:  corev1.ConditionTrue,
					Reason:  testReason,
					Message: testMessage,
				},
			},
			Succeeded: 1,
		},
	}
}
func newFakeDeployedJob(job interface{}) *unstructured.Unstructured {

	jobUnstructured, _ := util.ConvertObjectToUnstructured(job)
	return jobUnstructured
}
