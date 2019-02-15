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
	"time"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1alpha2 "github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1alpha2"
)

func NewPyTorchJobWithCleanPolicy(master, worker int, policy v1alpha2.CleanPodPolicy) *v1alpha2.PyTorchJob {
	if master == 1 {
		job := NewPyTorchJobWithMaster(worker)
		job.Spec.CleanPodPolicy = &policy
		return job
	}
	job := NewPyTorchJob(worker)
	job.Spec.CleanPodPolicy = &policy
	return job
}

func NewPyTorchJobWithCleanupJobDelay(master, worker int, ttl *int32) *v1alpha2.PyTorchJob {
	if master == 1 {
		job := NewPyTorchJobWithMaster(worker)
		job.Spec.TTLSecondsAfterFinished = ttl
		policy := v1alpha2.CleanPodPolicyNone
		job.Spec.CleanPodPolicy = &policy
		return job
	}
	job := NewPyTorchJob(worker)
	job.Spec.TTLSecondsAfterFinished = ttl
	policy := v1alpha2.CleanPodPolicyNone
	job.Spec.CleanPodPolicy = &policy
	return job
}

func NewPyTorchJobWithMaster(worker int) *v1alpha2.PyTorchJob {
	job := NewPyTorchJob(worker)
	job.Spec.PyTorchReplicaSpecs[v1alpha2.PyTorchReplicaTypeMaster] = &v1alpha2.PyTorchReplicaSpec{
		Template: NewPyTorchReplicaSpecTemplate(),
	}
	return job
}

func NewPyTorchJob(worker int) *v1alpha2.PyTorchJob {
	job := &v1alpha2.PyTorchJob{
		TypeMeta: metav1.TypeMeta{
			Kind: v1alpha2.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      TestPyTorchJobName,
			Namespace: metav1.NamespaceDefault,
		},
		Spec: v1alpha2.PyTorchJobSpec{
			PyTorchReplicaSpecs: make(map[v1alpha2.PyTorchReplicaType]*v1alpha2.PyTorchReplicaSpec),
		},
	}

	if worker > 0 {
		worker := int32(worker)
		workerReplicaSpec := &v1alpha2.PyTorchReplicaSpec{
			Replicas: &worker,
			Template: NewPyTorchReplicaSpecTemplate(),
		}
		job.Spec.PyTorchReplicaSpecs[v1alpha2.PyTorchReplicaTypeWorker] = workerReplicaSpec
	}

	return job
}

func NewPyTorchReplicaSpecTemplate() v1.PodTemplateSpec {
	return v1.PodTemplateSpec{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				v1.Container{
					Name:  v1alpha2.DefaultContainerName,
					Image: TestImageName,
					Args:  []string{"Fake", "Fake"},
					Ports: []v1.ContainerPort{
						v1.ContainerPort{
							Name:          v1alpha2.DefaultPortName,
							ContainerPort: v1alpha2.DefaultPort,
						},
					},
				},
			},
		},
	}
}

func SetPyTorchJobCompletionTime(job *v1alpha2.PyTorchJob) {
	now := metav1.Time{Time: time.Now()}
	job.Status.CompletionTime = &now
}
