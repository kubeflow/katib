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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	v1 "k8s.io/api/core/v1"
	ktypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"

	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/experiment/manifest"
	"github.com/kubeflow/katib/pkg/webhook/v1alpha3/common"
	"github.com/kubeflow/katib/pkg/webhook/v1alpha3/experiment/validator"
)

// experimentValidator validates Pods
type experimentValidator struct {
	admission.Handler
	client  client.Client
	decoder types.Decoder
	validator.Validator
}

func NewExperimentValidator(c client.Client) *experimentValidator {
	p := manifest.New(c)
	return &experimentValidator{
		Validator: validator.New(p),
	}
}

func (v *experimentValidator) Handle(ctx context.Context, req types.Request) types.Response {
	inst := &experimentsv1alpha3.Experiment{}
	var oldInst *experimentsv1alpha3.Experiment
	err := v.decoder.Decode(req, inst)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}
	if len(req.AdmissionRequest.OldObject.Raw) > 0 {
		oldDecoder := json.NewDecoder(bytes.NewBuffer(req.AdmissionRequest.OldObject.Raw))
		oldInst = &experimentsv1alpha3.Experiment{}
		if err := oldDecoder.Decode(&oldInst); err != nil {
			return admission.ErrorResponse(http.StatusBadRequest, fmt.Errorf("cannot decode incoming old object: %v", err))
		}
	}

	// After metrics collector sidecar injection in Job level done, delete validation for namespace labels
	ns := &v1.Namespace{}
	if err := v.client.Get(context.TODO(), ktypes.NamespacedName{Name: req.AdmissionRequest.Namespace}, ns); err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)
	}
	validNS := true
	if ns.Labels == nil {
		validNS = false
	} else {
		if v, ok := ns.Labels[common.KatibMetricsCollectorInjection]; !ok || v != common.KatibMetricsCollectorInjectionEnabled {
			validNS = false
		}
	}
	if !validNS {
		err = fmt.Errorf("Cannot create the Experiment %q in namespace %q: the namespace lacks label \"%s: %s\"",
			inst.Name, req.AdmissionRequest.Namespace, common.KatibMetricsCollectorInjection, common.KatibMetricsCollectorInjectionEnabled)
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	err = v.ValidateExperiment(inst, oldInst)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}
	return admission.ValidationResponse(true, "")
}

// experimentValidator implements inject.Client.
// A client will be automatically injected.
var _ inject.Client = &experimentValidator{}

// InjectClient injects the client.
func (v *experimentValidator) InjectClient(c client.Client) error {
	v.client = c
	v.Validator.InjectClient(c)
	return nil
}

// experimentValidator implements inject.Decoder.
// A decoder will be automatically injected.
var _ inject.Decoder = &experimentValidator{}

// InjectDecoder injects the decoder.
func (v *experimentValidator) InjectDecoder(d types.Decoder) error {
	v.decoder = d
	return nil
}
