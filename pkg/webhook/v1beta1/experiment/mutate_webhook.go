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
	"context"
	"encoding/json"
	"net/http"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
)

// ExperimentDefaulter sets the Experiment default values.
type ExperimentDefaulter struct {
	client  client.Client
	decoder *admission.Decoder
}

// NewExperimentDefaulter returns a new Experiment defaulter with the given client.
func NewExperimentDefaulter(c client.Client, d *admission.Decoder) *ExperimentDefaulter {
	return &ExperimentDefaulter{
		client:  c,
		decoder: d,
	}
}

func (e *ExperimentDefaulter) Handle(ctx context.Context, req admission.Request) admission.Response {
	exp := &experimentsv1beta1.Experiment{}
	err := e.decoder.Decode(req, exp)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	expDefault := exp.DeepCopy()
	expDefault.SetDefault()

	marshaledExperiment, err := json.Marshal(expDefault)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.AdmissionRequest.Object.Raw, marshaledExperiment)
}
