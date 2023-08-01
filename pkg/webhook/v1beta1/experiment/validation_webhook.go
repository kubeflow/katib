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

package experiment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	commonv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/common/v1beta1"
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/experiment/manifest"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/util"
	"github.com/kubeflow/katib/pkg/webhook/v1beta1/common"
	"github.com/kubeflow/katib/pkg/webhook/v1beta1/experiment/validator"
)

// ExperimentValidator validates Experiments.
type ExperimentValidator struct {
	client  client.Client
	decoder *admission.Decoder
	validator.Validator
}

// NewExperimentValidator returns a new Experiment validator with the given client.
func NewExperimentValidator(c client.Client, d *admission.Decoder) *ExperimentValidator {
	p := manifest.New(c)
	return &ExperimentValidator{
		client:    c,
		Validator: validator.New(p),
		decoder:   d,
	}
}

func (v *ExperimentValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	inst := &experimentsv1beta1.Experiment{}
	var oldInst *experimentsv1beta1.Experiment
	err := v.decoder.Decode(req, inst)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	if len(req.AdmissionRequest.OldObject.Raw) > 0 {
		oldDecoder := json.NewDecoder(bytes.NewBuffer(req.AdmissionRequest.OldObject.Raw))
		oldInst = &experimentsv1beta1.Experiment{}
		if err := oldDecoder.Decode(&oldInst); err != nil {
			return admission.Errored(http.StatusBadRequest, fmt.Errorf("Cannot decode incoming old object: %v", err))
		}
	}

	// After metrics collector sidecar injection in Job level done, delete validation for namespace labels
	ns := &v1.Namespace{}
	if err := v.client.Get(context.TODO(), types.NamespacedName{Name: req.AdmissionRequest.Namespace}, ns); err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
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
		return admission.Errored(http.StatusBadRequest, err)
	}

	err = v.ValidateExperiment(inst, oldInst)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	// By default we create PV with cluster local host path for each experiment with ResumePolicy = FromVolume.
	// User should not submit new experiment if previous PV for the same experiment was not deleted.
	// We unable to watch for the PV events in controller.
	// Webhook forbids experiment creation until coresponding PV will be deleted.
	if inst.Spec.ResumePolicy == experimentsv1beta1.FromVolume && oldInst == nil {
		// Create suggestion with name, namespace and algorithm name to get appropriate PV
		suggestion := &suggestionsv1beta1.Suggestion{
			ObjectMeta: metav1.ObjectMeta{
				Name:      inst.Name,
				Namespace: inst.Namespace,
			},
			Spec: suggestionsv1beta1.SuggestionSpec{
				Algorithm: &commonv1beta1.AlgorithmSpec{
					AlgorithmName: inst.Spec.Algorithm.AlgorithmName,
				},
			},
		}

		// Get PV name from Suggestion
		PVName := util.GetSuggestionPersistentVolumeName(suggestion)
		err := v.client.Get(context.TODO(), types.NamespacedName{Name: PVName}, &v1.PersistentVolume{})
		if !errors.IsNotFound(err) {
			returnError := fmt.Errorf("Cannot create the Experiment: %v in namespace: %v, PV: %v is not deleted", inst.Name, inst.Namespace, PVName)
			if err != nil {
				returnError = fmt.Errorf("Cannot create the Experiment: %v in namespace: %v, error: %v", inst.Name, inst.Namespace, err)
			}
			return admission.Errored(http.StatusBadRequest, returnError)
		}
	}

	return admission.ValidationResponse(true, "")
}
