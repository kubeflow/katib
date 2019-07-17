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

package v1beta2

import (
	common "github.com/kubeflow/tf-operator/pkg/apis/common/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=pytorchjob

// PyTorchJob represents the configuration of PyTorchJob
type PyTorchJob struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Specification of the desired behavior of the PyTorchJob.
	Spec PyTorchJobSpec `json:"spec,omitempty"`

	// Most recently observed status of the PyTorchJob.
	// This data may not be up to date.
	// Populated by the system.
	// Read-only.
	Status common.JobStatus `json:"status,omitempty"`
}

// PyTorchJobSpec is a desired state description of the PyTorchJob.
type PyTorchJobSpec struct {
	// Specifies the duration in seconds relative to the startTime that the job may be active
	// before the system tries to terminate it; value must be positive integer.
	// This method applies only to pods with restartPolicy == OnFailure or Always.
	// +optional
	ActiveDeadlineSeconds *int64 `json:"activeDeadlineSeconds,omitempty"`

	// Optional number of retries before marking this job failed.
	// +optional
	BackoffLimit *int32 `json:"backoffLimit,omitempty"`

	// CleanPodPolicy defines the policy to kill pods after PyTorchJob is
	// succeeded.
	// Default to Running.
	CleanPodPolicy *common.CleanPodPolicy `json:"cleanPodPolicy,omitempty"`

	// TTLSecondsAfterFinished is the TTL to clean up pytorch-jobs (temporary
	// before kubernetes adds the cleanup controller).
	// It may take extra ReconcilePeriod seconds for the cleanup, since
	// reconcile gets called periodically.
	// Default to infinite.
	TTLSecondsAfterFinished *int32 `json:"ttlSecondsAfterFinished,omitempty"`

	// PyTorchReplicaSpecs is map of PyTorchReplicaType and PyTorchReplicaSpec
	// specifies the PyTorch replicas to run.
	// For example,
	//   {
	//     "Master": PyTorchReplicaSpec,
	//     "Worker": PyTorchReplicaSpec,
	//   }
	PyTorchReplicaSpecs map[PyTorchReplicaType]*common.ReplicaSpec `json:"pytorchReplicaSpecs"`
}

// PyTorchReplicaType is the type for PyTorchReplica.
type PyTorchReplicaType common.ReplicaType

const (
	// PyTorchReplicaTypeMaster is the type of Master of distributed PyTorch
	PyTorchReplicaTypeMaster PyTorchReplicaType = "Master"

	// PyTorchReplicaTypeWorker is the type for workers of distributed PyTorch.
	PyTorchReplicaTypeWorker PyTorchReplicaType = "Worker"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=pytorchjobs

// PyTorchJobList is a list of PyTorchJobs.
type PyTorchJobList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	metav1.ListMeta `json:"metadata,omitempty"`

	// List of PyTorchJobs.
	Items []PyTorchJob `json:"items"`
}
