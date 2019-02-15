// Copyright 2018 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testutil

import (
	"encoding/json"
	"strings"
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/cache"

	v1beta1 "github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1beta1"
	common "github.com/kubeflow/tf-operator/pkg/apis/common/v1beta1"
)

const (
	LabelGroupName      = "group_name"
	LabelPyTorchJobName = "pytorch_job_name"
)

var (
	// KeyFunc is the short name to DeletionHandlingMetaNamespaceKeyFunc.
	// IndexerInformer uses a delta queue, therefore for deletes we have to use this
	// key function but it should be just fine for non delete events.
	KeyFunc   = cache.DeletionHandlingMetaNamespaceKeyFunc
	GroupName = v1beta1.GroupName
)

func GenLabels(jobName string) map[string]string {
	return map[string]string{
		LabelGroupName:      GroupName,
		LabelPyTorchJobName: strings.Replace(jobName, "/", "-", -1),
	}
}

func GenOwnerReference(job *v1beta1.PyTorchJob) *metav1.OwnerReference {
	boolPtr := func(b bool) *bool { return &b }
	controllerRef := &metav1.OwnerReference{
		APIVersion:         v1beta1.SchemeGroupVersion.String(),
		Kind:               v1beta1.Kind,
		Name:               job.Name,
		UID:                job.UID,
		BlockOwnerDeletion: boolPtr(true),
		Controller:         boolPtr(true),
	}

	return controllerRef
}

// ConvertPyTorchJobToUnstructured uses JSON to convert PyTorchJob to Unstructured.
func ConvertPyTorchJobToUnstructured(job *v1beta1.PyTorchJob) (*unstructured.Unstructured, error) {
	var unstructured unstructured.Unstructured
	b, err := json.Marshal(job)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &unstructured); err != nil {
		return nil, err
	}
	return &unstructured, nil
}

func GetKey(job *v1beta1.PyTorchJob, t *testing.T) string {
	key, err := KeyFunc(job)
	if err != nil {
		t.Errorf("Unexpected error getting key for job %v: %v", job.Name, err)
		return ""
	}
	return key
}

func CheckCondition(job *v1beta1.PyTorchJob, condition common.JobConditionType, reason string) bool {
	for _, v := range job.Status.Conditions {
		if v.Type == condition && v.Status == v1.ConditionTrue && v.Reason == reason {
			return true
		}
	}
	return false
}
