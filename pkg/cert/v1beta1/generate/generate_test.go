/*
Copyright 2021 The Kubeflow Authors.

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

package generate

import (
	"github.com/kubeflow/katib/pkg/cert/v1beta1/consts"
	"github.com/kubeflow/katib/pkg/cert/v1beta1/kube"
	admissionregistration "k8s.io/api/admissionregistration/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"log"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func TestGenerate(t *testing.T) {

	const testNamespace = "test"

	testGeneratorJob := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      consts.JobName,
			Namespace: testNamespace,
			UID:       "test",
		},
	}
	testValidatingWebhook := &admissionregistration.ValidatingWebhookConfiguration{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "admissionregistration.k8s.io/v1",
			Kind:       "ValidatingWebhookConfiguration",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "katib.kubeflow.org",
		},
		Webhooks: []admissionregistration.ValidatingWebhook{
			{
				Name: "validator.experiment.katib.kubeflow.org",
				ClientConfig: admissionregistration.WebhookClientConfig{
					CABundle: []byte("CG=="),
				},
			},
		},
	}
	testMutatingWebhook := &admissionregistration.MutatingWebhookConfiguration{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "admissionregistration.k8s.io/v1",
			Kind:       "MutatingWebhookConfiguration",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "katib.kubeflow.org",
		},
		Webhooks: []admissionregistration.MutatingWebhook{
			{
				Name: "defaulter.experiment.katib.kubeflow.org",
				ClientConfig: admissionregistration.WebhookClientConfig{
					CABundle: []byte("CG=="),
				},
			},
			{
				Name: "mutator.pod.katib.kubeflow.org",
				ClientConfig: admissionregistration.WebhookClientConfig{
					CABundle: []byte("CG"),
				},
			},
		},
	}

	oldWebhookCertSecret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      consts.Secret,
			Namespace: testNamespace,
		},
	}

	tests := []struct {
		name      string
		wantError bool
		objects   []client.Object
	}{
		{
			name:      "generate successfully",
			wantError: false,
			objects: []client.Object{
				testGeneratorJob,
				testValidatingWebhook,
				testMutatingWebhook,
			},
		},
		{
			name:      "old secret exists",
			wantError: false,
			objects: []client.Object{
				testGeneratorJob,
				testValidatingWebhook,
				testMutatingWebhook,
				oldWebhookCertSecret,
			},
		},
		{
			name:      "missing katib-cert-generator job",
			wantError: true,
			objects: []client.Object{
				testValidatingWebhook,
				testMutatingWebhook,
			},
		},
		{
			name:      "missing validatingWebhookConfiguration",
			wantError: true,
			objects: []client.Object{
				testGeneratorJob,
				testMutatingWebhook,
			},
		},
		{
			name:      "missing mutatingWebhookConfiguration",
			wantError: true,
			objects: []client.Object{
				testGeneratorJob,
				testValidatingWebhook,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := executeGeneratorCommand(test.objects, testNamespace); (err != nil) != test.wantError {
				t.Errorf("expected error, got '%v'\n", err)
			}
		})
	}

}

func executeGeneratorCommand(kubeResources []client.Object, namespace string) error {
	scm := runtime.NewScheme()
	if err := scheme.AddToScheme(scm); err != nil {
		log.Fatal(err)
	}
	fakeClientBuilder := fake.NewClientBuilder().WithScheme(scm)
	if len(kubeResources) > 0 {
		for _, r := range kubeResources {
			fakeClientBuilder.WithObjects(r)
		}
	}

	c := &kube.Client{KubeClient: fakeClientBuilder.Build()}
	cmd := NewGenerateCmd(c)
	if err := cmd.Flags().Set("namespace", namespace); err != nil {
		log.Fatal(err)
	}

	return cmd.Execute()
}
