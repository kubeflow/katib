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
	"fmt"
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"

	v1alpha2 "github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1alpha2"
)

func NewBaseService(name string, job *v1alpha2.PyTorchJob, t *testing.T) *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Labels:          GenLabels(job.Name),
			Namespace:       job.Namespace,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(job, controllerKind)},
		},
	}
}

func NewService(job *v1alpha2.PyTorchJob, typ string, index int, t *testing.T) *v1.Service {
	service := NewBaseService(fmt.Sprintf("%s-%d", typ, index), job, t)
	service.Labels[replicaTypeLabel] = typ
	service.Labels[replicaIndexLabel] = fmt.Sprintf("%d", index)
	return service
}

// NewServiceList creates count pods with the given phase for the given job
func NewServiceList(count int32, job *v1alpha2.PyTorchJob, typ string, t *testing.T) []*v1.Service {
	services := []*v1.Service{}
	for i := int32(0); i < count; i++ {
		newService := NewService(job, typ, int(i), t)
		services = append(services, newService)
	}
	return services
}

func SetServices(serviceIndexer cache.Indexer, job *v1alpha2.PyTorchJob, typ string, activeWorkerServices int32, t *testing.T) {
	for _, service := range NewServiceList(activeWorkerServices, job, typ, t) {
		if err := serviceIndexer.Add(service); err != nil {
			t.Errorf("unexpected error when adding service %v", err)
		}
	}
}
