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
	"github.com/spf13/viper"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/builder"

	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/consts"
	"github.com/kubeflow/katib/pkg/webhook/v1alpha3/common"
	"github.com/kubeflow/katib/pkg/webhook/v1alpha3/experiment"
	"github.com/kubeflow/katib/pkg/webhook/v1alpha3/pod"
)

func AddToManager(m manager.Manager, port int32, serviceName string) error {
	so := webhook.ServerOptions{
		CertDir: "/tmp/cert",
		BootstrapOptions: &webhook.BootstrapOptions{
			Service: &webhook.Service{
				Namespace: consts.DefaultKatibNamespace,
				Name:      serviceName,
				Selectors: map[string]string{
					"app": serviceName,
				},
			},
			ValidatingWebhookConfigName: "katib-validating-webhook-config",
			MutatingWebhookConfigName:   "katib-mutating-webhook-config",
		},
		Port: port,
	}

	// Decide if we should use local file system.
	// If not, we set a secret in BootstrapOptions.
	usingFS := viper.GetBool(consts.ConfigCertLocalFS)
	if !usingFS {
		so.BootstrapOptions.Secret = &types.NamespacedName{
			Namespace: consts.DefaultKatibNamespace,
			Name:      serviceName,
		}
	}
	server, err := webhook.NewServer("katib-admission-server", m, so)
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
		FailurePolicy(admissionregistrationv1beta1.Fail).
		Name("mutating.experiment.katib.kubeflow.org").
		Mutating().
		Operations(admissionregistrationv1beta1.Create, admissionregistrationv1beta1.Update).
		WithManager(manager).
		ForType(&experimentsv1alpha3.Experiment{}).
		Handlers(experiment.NewExperimentDefaulter(manager.GetClient())).
		Build()
	if err != nil {
		return err
	}
	validatingWebhook, err := builder.NewWebhookBuilder().
		FailurePolicy(admissionregistrationv1beta1.Fail).
		Name("validating.experiment.katib.kubeflow.org").
		Validating().
		Operations(admissionregistrationv1beta1.Create, admissionregistrationv1beta1.Update).
		WithManager(manager).
		ForType(&experimentsv1alpha3.Experiment{}).
		Handlers(experiment.NewExperimentValidator(manager.GetClient())).
		Build()
	if err != nil {
		return err
	}
	nsSelector := &metav1.LabelSelector{
		MatchLabels: map[string]string{
			common.KatibMetricsCollectorInjection: common.KatibMetricsCollectorInjectionEnabled,
		},
	}
	injectWebhook, err := builder.NewWebhookBuilder().
		FailurePolicy(admissionregistrationv1beta1.Fail).
		Name("mutating.pod.katib.kubeflow.org").
		NamespaceSelector(nsSelector).
		Mutating().
		Operations(admissionregistrationv1beta1.Create).
		WithManager(manager).
		ForType(&v1.Pod{}).
		Handlers(pod.NewSidecarInjector(manager.GetClient())).
		FailurePolicy(admissionregistrationv1beta1.Ignore).
		Build()
	if err != nil {
		return err
	}
	return server.Register(mutatingWebhook, validatingWebhook, injectWebhook)
}
