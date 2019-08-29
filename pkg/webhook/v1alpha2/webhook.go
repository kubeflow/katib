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

package webhook

import (
	"os"

	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	"github.com/kubeflow/katib/pkg/webhook/v1alpha2/experiment"
	"github.com/kubeflow/katib/pkg/webhook/v1alpha2/pod"

	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/builder"
)

const (
	katibControllerName                   = "katib-controller"
	katibMetricsCollectorInjection        = "katib-metricscollector-injection"
	katibMetricsCollectorInjectionEnabled = "enabled"
)

func AddToManager(m manager.Manager) error {
	server, err := webhook.NewServer("katib-admission-server", m, webhook.ServerOptions{
		CertDir: "/tmp/cert",
		BootstrapOptions: &webhook.BootstrapOptions{
			Secret: &types.NamespacedName{
				Namespace: os.Getenv(experimentsv1alpha2.DefaultKatibNamespaceEnvName),
				Name:      katibControllerName,
			},
			Service: &webhook.Service{
				Namespace: os.Getenv(experimentsv1alpha2.DefaultKatibNamespaceEnvName),
				Name:      katibControllerName,
				Selectors: map[string]string{
					"app": katibControllerName,
				},
			},
			ValidatingWebhookConfigName: "katib-validating-webhook-config",
			MutatingWebhookConfigName:   "katib-mutating-webhook-config",
		},
	})
	if err != nil {
		return err
	}

	if err := register(m, server); err != nil {
		return err
	}

	return nil
}

func register(manager manager.Manager, server *webhook.Server) error {
	mutatingWebhook, err := builder.NewWebhookBuilder().
		Name("mutating.experiment.katib.kubeflow.org").
		Mutating().
		Operations(admissionregistrationv1beta1.Create, admissionregistrationv1beta1.Update).
		WithManager(manager).
		ForType(&experimentsv1alpha2.Experiment{}).
		Handlers(experiment.NewExperimentDefaulter(manager.GetClient())).
		Build()
	if err != nil {
		return err
	}
	validatingWebhook, err := builder.NewWebhookBuilder().
		Name("validating.experiment.katib.kubeflow.org").
		Validating().
		Operations(admissionregistrationv1beta1.Create, admissionregistrationv1beta1.Update).
		WithManager(manager).
		ForType(&experimentsv1alpha2.Experiment{}).
		Handlers(experiment.NewExperimentValidator(manager.GetClient())).
		Build()
	if err != nil {
		return err
	}
	nsSelector := &metav1.LabelSelector{
		MatchLabels: map[string]string{
			katibMetricsCollectorInjection: katibMetricsCollectorInjectionEnabled,
		},
	}
	injectWebhook, err := builder.NewWebhookBuilder().
		Name("mutating.pod.katib.kubeflow.org").
		NamespaceSelector(nsSelector).
		Mutating().
		Operations(admissionregistrationv1beta1.Create).
		WithManager(manager).
		ForType(&v1.Pod{}).
		Handlers(pod.NewSidecarInjector(manager.GetClient(), manager.GetConfig().Host)).
		Build()
	if err != nil {
		return err
	}
	return server.Register(mutatingWebhook, validatingWebhook, injectWebhook)
}
