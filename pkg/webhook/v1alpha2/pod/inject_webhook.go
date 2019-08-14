/*
Copyright 2019 The Kubernetes Authors.

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

package pod

import (
	"context"
	"net/http"

	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"

	v1 "k8s.io/api/core/v1"
	// "github.com/kubeflow/katib/pkg/webhook/v1alpha2/experiment/injector"
)

// For debug
var log = logf.Log.WithName("injector-webhook")

// sidecarInjector that inject metrics collect sidecar into master pod
type sidecarInjector struct {
	client  client.Client
	decoder types.Decoder
	// injector.Injector
}

var _ admission.Handler = &sidecarInjector{}

func (s *sidecarInjector) Handle(ctx context.Context, req types.Request) types.Response {
	pod := &v1.Pod{}
	err := s.decoder.Decode(req, pod)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	// Check whether the pod need to be mutated
	if !s.MutationRequired(pod) {
		log.Info("Skipping mutation for " + pod.Name + " due to policy check")
		return admission.ValidationResponse(true, "")
	}

	// Do mutation
	mutatedPod := s.Mutate(pod)
	// if err != nil {
	// 	return admission.ErrorResponse(http.StatusBadRequest, err)
	// }

	return admission.PatchResponse(pod, mutatedPod)
}

var _ inject.Client = &sidecarInjector{}

func (s *sidecarInjector) InjectClient(c client.Client) error {
	s.client = c
	return nil
}

var _ inject.Decoder = &sidecarInjector{}

func (s *sidecarInjector) InjectDecoder(d types.Decoder) error {
	s.decoder = d
	return nil
}

func NewSidecarInjector(c client.Client) *sidecarInjector {
	return &sidecarInjector{
		client: c,
	}
}

func (s *sidecarInjector) MutationRequired(pod *v1.Pod) bool {
	labels := pod.Labels
	for k, v := range labels {
		log.Info("FOR TEST: pod label {" + k + ": " + v + "}")
	}

	return true
}

func (s *sidecarInjector) Mutate(pod *v1.Pod) *v1.Pod {
	mutatedPod := pod.DeepCopy()

	// Hard code container, just for test
	injectContainer := v1.Container{
		Name:  "sidecar-nginxk",
		Image: "nginx:1.12.2",
	}
	mutatedPod.Spec.Containers = append(mutatedPod.Spec.Containers, injectContainer)

	return mutatedPod
}
