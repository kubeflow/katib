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

package experiment

import (
	"context"
	"net/http"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"

	v1 "k8s.io/api/core/v1"

	"github.com/kubeflow/katib/pkg/webhook/v1alpha2/experiment/injector"
)

// experimentDefaulter that sets default fields in experiment
type sidecarInjector struct {
	client  client.Client
	decoder types.Decoder
	injector.Injector
}

var _ admission.Handler = &sidecarInjector{}

func (e *sidecarInjector) Handle(ctx context.Context, req types.Request) types.Response {
	inst := &v1.Pod{}
	err := e.decoder.Decode(req, inst)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	// TODO:

	return admission.PatchResponse(inst, copy)
}

var _ inject.Client = &sidecarInjector{}

func (e *sidecarInjector) InjectClient(c client.Client) error {
	e.client = c
	return nil
}

var _ inject.Decoder = &sidecarInjector{}

func (e *sidecarInjector) InjectDecoder(d types.Decoder) error {
	e.decoder = d
	return nil
}

func NewSidecarInjector(c client.Client) *sidecarInjector {
	return &sidecarInjector{
		client: c,
	}
}
