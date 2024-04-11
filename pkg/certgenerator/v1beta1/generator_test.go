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
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	admissionregistration "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"

	configv1beta1 "github.com/kubeflow/katib/pkg/apis/config/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

func TestGenerate(t *testing.T) {
	const testNamespace = "test"

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
	webhookSecret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "katib-test-secret",
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
				namespace:          testNamespace,
				webhookServiceName: "katib-controller",
				webhookSecretName:  "katib-test-secret",
			},
			objects: []client.Object{
				emptyVWebhookConfig,
				emptyMWebhookConfig,
				controllerService,
				webhookSecret,
			},
		},
		"There is not ValidatingWebhookConfiguration": {
			opts: &CertGenerator{
				namespace:          testNamespace,
				webhookServiceName: "katib-controller",
				webhookSecretName:  "katib-test-secret",
			},
			objects: []client.Object{
				emptyMWebhookConfig,
				controllerService,
				webhookSecret,
			},
			wantError: errInjectCertError,
		},
		"There is not MutatingWebhookConfiguration": {
			opts: &CertGenerator{
				namespace:          testNamespace,
				webhookServiceName: "katib-controller",
				webhookSecretName:  "katib-test-secret",
			},
			objects: []client.Object{
				emptyVWebhookConfig,
				controllerService,
				webhookSecret,
			},
			wantError: errInjectCertError,
		},
		"There is not Service katib-controller": {
			opts: &CertGenerator{
				namespace:          testNamespace,
				webhookServiceName: "katib-controller",
				webhookSecretName:  "katib-test-secret",
			},
			objects: []client.Object{
				emptyVWebhookConfig,
				emptyMWebhookConfig,
				webhookSecret,
			},
			wantError: errServiceNotFound,
		},
		"There is not Secret katib-webhook-cert": {
			opts: &CertGenerator{
				namespace:          testNamespace,
				webhookServiceName: "katib-controller",
				webhookSecretName:  "katib-test-secret",
			},
			objects: []client.Object{
				emptyVWebhookConfig,
				emptyMWebhookConfig,
				controllerService,
			},
			wantError: errCertCheckFail,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			kc := buildFakeClient(tc.objects)
			tc.opts.kubeClient = kc
			err := tc.opts.generate(context.Background())
			if diff := cmp.Diff(tc.wantError, err, cmpopts.EquateErrors()); len(diff) != 0 {
				t.Errorf("Unexpected error from generate() (-want,+got):\n%s", diff)
			}

			if tc.wantError == nil {
				secret := &corev1.Secret{}
				if err = kc.Get(context.Background(), client.ObjectKey{Name: tc.opts.webhookSecretName, Namespace: testNamespace}, secret); err != nil {
					t.Fatalf("Failed to get a webhookSecret: %v", err)
				}
				if len(secret.Data[serverKeyName]) == 0 {
					t.Errorf("Unexpected tls.key embedded in secret: %v", secret.Data)
				}
				if len(secret.Data[serverCertName]) == 0 {
					t.Errorf("Unexpected tls.crt embedded in secret: %v", secret.Data)
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
	fakeClientBuilder := fake.NewClientBuilder().
		WithScheme(scheme.Scheme).
		WithInterceptorFuncs(interceptor.Funcs{Patch: ssaAsStrategicMergePatchFunc})
	if len(kubeResources) > 0 {
		fakeClientBuilder.WithObjects(kubeResources...)
	}
	return fakeClientBuilder.Build()
}

type ssaPatchAsStrategicMerge struct {
	client.Patch
}

func (*ssaPatchAsStrategicMerge) Type() types.PatchType {
	return types.StrategicMergePatchType
}

func wrapSSAPatch(patch client.Patch) client.Patch {
	if patch.Type() == types.ApplyPatchType {
		return &ssaPatchAsStrategicMerge{Patch: patch}
	}
	return patch
}

// ssaAsStrategicMergePatchFunc returns the patch interceptor.
// TODO (tenzen-y): Once the fake client supports server-side apply, we should remove these interceptor.
// REF: https://github.com/kubernetes/kubernetes/issues/115598
func ssaAsStrategicMergePatchFunc(
	ctx context.Context,
	cli client.WithWatch,
	obj client.Object,
	patch client.Patch,
	opts ...client.PatchOption,
) error {
	return cli.Patch(ctx, obj, wrapSSAPatch(patch), opts...)
}

func TestEnsureCertMounted(t *testing.T) {
	tests := map[string]struct {
		keyExist  bool
		certExist bool
		wantExist bool
	}{
		"key and cert exist": {
			keyExist:  true,
			certExist: true,
			wantExist: true,
		},
		"key doesn't exist": {
			keyExist:  false,
			certExist: true,
			wantExist: false,
		},
		"cert doesn't exist": {
			keyExist:  true,
			certExist: false,
			wantExist: false,
		},
		"all files doesn't exist": {
			keyExist:  false,
			certExist: false,
			wantExist: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.keyExist || tc.certExist {
				if err := os.MkdirAll(consts.CertDir, 0760); err != nil {
					t.Fatalf("Failed to set up directory: %v", err)
				}
				defer func() {
					if err := os.RemoveAll(consts.CertDir); err != nil {
						t.Fatalf("Failed to clean up directory: %v", err)
					}
				}()
			}
			if tc.keyExist {
				if _, err := os.Create(filepath.Join(consts.CertDir, serverKeyName)); err != nil {
					t.Fatalf("Failed to create tls.key: %v", err)
				}
			}
			if tc.certExist {
				if _, err := os.Create(filepath.Join(consts.CertDir, serverCertName)); err != nil {
					t.Fatalf("Failed to create tls.crt: %v", err)
				}
			}
			ensureFunc := ensureCertMounted(time.Now())
			got, _ := ensureFunc(context.Background())
			if tc.wantExist != got {
				t.Errorf("Unexpected value from ensureCertMounted: \n(want: %v, got: %v)\n", tc.wantExist, got)
			}
		})
	}
}
