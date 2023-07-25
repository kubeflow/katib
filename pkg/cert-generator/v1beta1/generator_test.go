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

package v1beta1

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	admissionregistration "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	configv1beta1 "github.com/kubeflow/katib/pkg/apis/config/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

func TestGenerate(t *testing.T) {
	const testNamespace = "test"

	controllerDeployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "katib-controller",
			Namespace: testNamespace,
			UID:       "test",
		},
	}
	emptyVWebhookConfig := &admissionregistration.ValidatingWebhookConfiguration{
		TypeMeta: metav1.TypeMeta{
			APIVersion: admissionregistration.SchemeGroupVersion.String(),
			Kind:       "ValidatingWebhookConfiguration",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: Webhook,
		},
		Webhooks: []admissionregistration.ValidatingWebhook{
			{
				Name:         strings.Join([]string{"validator.experiment", Webhook}, "."),
				ClientConfig: admissionregistration.WebhookClientConfig{},
			},
		},
	}
	emptyMWebhookConfig := &admissionregistration.MutatingWebhookConfiguration{
		TypeMeta: metav1.TypeMeta{
			APIVersion: admissionregistration.SchemeGroupVersion.String(),
			Kind:       "MutatingWebhookConfiguration",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: Webhook,
		},
		Webhooks: []admissionregistration.MutatingWebhook{
			{
				Name:         strings.Join([]string{"defaulter.experiment", Webhook}, "."),
				ClientConfig: admissionregistration.WebhookClientConfig{},
			},
			{
				Name:         strings.Join([]string{"mutator.pod", Webhook}, "."),
				ClientConfig: admissionregistration.WebhookClientConfig{},
			},
		},
	}
	controllerSecret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Secret,
			Namespace: testNamespace,
		},
	}
	controllerService := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      configv1beta1.DefaultWebhookServiceName,
			Namespace: testNamespace,
		},
	}

	tests := map[string]struct {
		objects   []client.Object
		opts      *CertGenerator
		wantError error
	}{
		"Generate successfully": {
			opts: &CertGenerator{
				namespace:      testNamespace,
				serviceName:    "katib-controller",
				controllerName: "katib-controller",
			},
			objects: []client.Object{
				controllerDeployment,
				emptyVWebhookConfig,
				emptyMWebhookConfig,
				controllerService,
			},
		},
		"There is an old Secret, katib-webhook-cert": {
			opts: &CertGenerator{
				namespace:      testNamespace,
				serviceName:    "katib-controller",
				controllerName: "katib-controller",
			},
			objects: []client.Object{
				controllerDeployment,
				emptyVWebhookConfig,
				emptyMWebhookConfig,
				controllerService,
				controllerSecret,
			},
		},
		"There is not Deployment, katib-controller": {
			opts: &CertGenerator{
				namespace:      testNamespace,
				serviceName:    "katib-controller",
				controllerName: "katib-controller",
			},
			objects: []client.Object{
				emptyVWebhookConfig,
				emptyMWebhookConfig,
				controllerService,
			},
			wantError: errCreateCertSecretFail,
		},
		"There is not ValidatingWebhookConfiguration": {
			opts: &CertGenerator{
				namespace:      testNamespace,
				serviceName:    "katib-controller",
				controllerName: "katib-controller",
			},
			objects: []client.Object{
				controllerDeployment,
				emptyMWebhookConfig,
				controllerService,
			},
			wantError: errInjectCertError,
		},
		"There is not MutatingWebhookConfiguration": {
			opts: &CertGenerator{
				namespace:      testNamespace,
				serviceName:    "katib-controller",
				controllerName: "katib-controller",
			},
			objects: []client.Object{
				controllerDeployment,
				emptyVWebhookConfig,
				controllerService,
			},
			wantError: errInjectCertError,
		},
		"There is no Service katib-controller": {
			opts: &CertGenerator{
				namespace:      testNamespace,
				serviceName:    "katib-controller",
				controllerName: "katib-controller",
			},
			objects: []client.Object{
				controllerDeployment,
				emptyVWebhookConfig,
				emptyMWebhookConfig,
			},
			wantError: errServiceNotFound,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if err := os.RemoveAll(consts.CertDir); err != nil {
				t.Fatalf("Failed to clean up cert dir: %v", err)
			}

			kc := buildFakeClient(tc.objects)
			tc.opts.kubeClient = kc
			err := tc.opts.generate(context.Background())
			if diff := cmp.Diff(tc.wantError, err, cmpopts.EquateErrors()); len(diff) != 0 {
				t.Errorf("Unexpected error from generate() (-want,+got):\n%s", diff)
			}

			if tc.wantError == nil {
				secret := &corev1.Secret{}
				if err = kc.Get(context.Background(), client.ObjectKey{Name: Secret, Namespace: testNamespace}, secret); err != nil {
					t.Fatalf("Failed to get a controllerSecret: %v", err)
				}
				if !metav1.IsControlledBy(secret, controllerDeployment) {
					t.Errorf("Unexpected owner for the secret: %v", secret.OwnerReferences)
				}
				if len(secret.Data[serverKeyName]) == 0 {
					t.Errorf("Unexpected tls.key embedded in secret: %v", secret.Data)
				}
				if len(secret.Data[serverCertName]) == 0 {
					t.Errorf("Unexpected tls.crt embedded in secret: %v", secret.Data)
				}

				if _, err = os.Stat(filepath.Join(consts.CertDir, serverKeyName)); err != nil {
					t.Errorf("Failed to find tls.key: %v", err)
				}
				if _, err = os.Stat(filepath.Join(consts.CertDir, serverCertName)); err != nil {
					t.Errorf("Failed to find tls.crt: %v", err)
				}

				vConfig := &admissionregistration.ValidatingWebhookConfiguration{}
				if err = kc.Get(context.Background(), client.ObjectKey{Name: Webhook}, vConfig); err != nil {
					t.Fatalf("Failed to get a ValidatingWebhookConfiguration: %v", err)
				}
				if len(vConfig.Webhooks[0].ClientConfig.CABundle) == 0 {
					t.Errorf("Unexpected tls.crt embedded in ValidatingWebhookConfiguration: %v", vConfig.Webhooks)
				}

				mConfig := &admissionregistration.MutatingWebhookConfiguration{}
				if err = kc.Get(context.Background(), client.ObjectKey{Name: Webhook}, mConfig); err != nil {
					t.Fatalf("Failed to get a MutatingWebhookConfiguration: %v", err)
				}
				if len(mConfig.Webhooks[0].ClientConfig.CABundle) == 0 || len(mConfig.Webhooks[1].ClientConfig.CABundle) == 0 {
					t.Errorf("Unexpected tls.crt embedded in MutatingWebhookConfiguration: %v", mConfig.Webhooks)
				}
			}
		})
	}
}

func buildFakeClient(kubeResources []client.Object) client.Client {
	fakeClientBuilder := fake.NewClientBuilder().WithScheme(scheme.Scheme)
	if len(kubeResources) > 0 {
		fakeClientBuilder.WithObjects(kubeResources...)
	}
	return fakeClientBuilder.Build()
}
