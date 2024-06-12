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

package certgenerator

import (
	"fmt"

	cert "github.com/open-policy-agent/cert-controller/pkg/rotator"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	configv1beta1 "github.com/kubeflow/katib/pkg/apis/config/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

const Webhook = "katib.kubeflow.org"

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;update
// +kubebuilder:rbac:groups="admissionregistration.k8s.io",resources=validatingwebhookconfigurations,verbs=get;list;watch;update
// +kubebuilder:rbac:groups="admissionregistration.k8s.io",resources=mutatingwebhookconfigurations,verbs=get;list;watch;update

// AddToManager adds the cert-generator to the manager.
func AddToManager(mgr manager.Manager, cfg configv1beta1.CertGeneratorConfig, certsReady chan struct{}) error {
	return cert.AddRotator(mgr, &cert.CertRotator{
		SecretKey: types.NamespacedName{
			Namespace: consts.DefaultKatibNamespace,
			Name:      cfg.WebhookSecretName,
		},
		CertDir:        consts.CertDir,
		CAName:         "katib-ca",
		CAOrganization: "katib",
		DNSName:        fmt.Sprintf("%s.%s.svc", cfg.WebhookServiceName, consts.DefaultKatibNamespace),
		IsReady:        certsReady,
		Webhooks: []cert.WebhookInfo{
			{Name: Webhook, Type: cert.Validating},
			{Name: Webhook, Type: cert.Mutating},
		},
		FieldOwner: "cert-generator",
		// When the Katib is running in the leader election mode,
		// we expect webhook server will run in primary and secondary instance
		RequireLeaderElection: false,
	})
}
