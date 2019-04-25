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

        experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	"github.com/kubeflow/katib/pkg/controller/v1alpha2/experiment/util"
        "sigs.k8s.io/controller-runtime/pkg/client"
        "sigs.k8s.io/controller-runtime/pkg/runtime/inject"
        "sigs.k8s.io/controller-runtime/pkg/webhook/admission"
        "sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

// experimentDefaulter that sets default fields in experiment
type experimentDefaulter struct {
	client  client.Client
	decoder types.Decoder
}

var _ admission.Handler = &experimentDefaulter{}

func (e *experimentDefaulter) Handle(ctx context.Context, req types.Request) types.Response {
        inst := &experimentsv1alpha2.Experiment{}
        err := e.decoder.Decode(req, inst)
        if err != nil {
                return admission.ErrorResponse(http.StatusBadRequest, err)
        }

	copy := inst.DeepCopy()
        err = util.DefaultExperiment(copy)
        if err != nil {
                return admission.ErrorResponse(http.StatusInternalServerError, err)
        }

	return admission.PatchResponse(inst, copy)
}

var _ inject.Client = &experimentDefaulter{}

func (e *experimentDefaulter) InjectClient(c client.Client) error {
	e.client = c
	return nil
}

var _ inject.Decoder = &experimentDefaulter{}

func (e *experimentDefaulter) InjectDecoder(d types.Decoder) error {
	e.decoder = d
	return nil
}
