/*
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

package studyjob

import (
	"context"
	"net/http"

	katibv1alpha1 "github.com/kubeflow/katib/pkg/api/operators/apis/studyjob/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

// studyJobValidator validates Pods
type studyJobValidator struct {
	client  client.Client
	decoder types.Decoder
}

// Implement admission.Handler so the controller can handle admission request.
var _ admission.Handler = &studyJobValidator{}

func (v *studyJobValidator) Handle(ctx context.Context, req types.Request) types.Response {
	inst := &katibv1alpha1.StudyJob{}
	err := v.decoder.Decode(req, inst)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	err = validateStudy(inst)
	if err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)
	}
	return admission.ValidationResponse(true, "")
}

// studyJobValidator implements inject.Client.
// A client will be automatically injected.
var _ inject.Client = &studyJobValidator{}

// InjectClient injects the client.
func (v *studyJobValidator) InjectClient(c client.Client) error {
	v.client = c
	return nil
}

// studyJobValidator implements inject.Decoder.
// A decoder will be automatically injected.
var _ inject.Decoder = &studyJobValidator{}

// InjectDecoder injects the decoder.
func (v *studyJobValidator) InjectDecoder(d types.Decoder) error {
	v.decoder = d
	return nil
}
