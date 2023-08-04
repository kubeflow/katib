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

package webhook

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/kubeflow/katib/pkg/webhook/v1beta1/experiment"
	"github.com/kubeflow/katib/pkg/webhook/v1beta1/pod"
)

func AddToManager(mgr manager.Manager, hookServer webhook.Server) error {
	if err := mgr.Add(hookServer); err != nil {
		return fmt.Errorf("Add webhook server to the manager failed: %v", err)
	}

	decoder := admission.NewDecoder(mgr.GetScheme())
	experimentValidator := experiment.NewExperimentValidator(mgr.GetClient(), decoder)
	experimentDefaulter := experiment.NewExperimentDefaulter(mgr.GetClient(), decoder)
	sidecarInjector := pod.NewSidecarInjector(mgr.GetClient(), decoder)

	hookServer.Register("/validate-experiment", &webhook.Admission{Handler: experimentValidator})
	hookServer.Register("/mutate-experiment", &webhook.Admission{Handler: experimentDefaulter})
	hookServer.Register("/mutate-pod", &webhook.Admission{Handler: sidecarInjector})
	return nil
}
