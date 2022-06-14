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

package generate

import (
	"github.com/kubeflow/katib/pkg/cert-generator/v1beta1/consts"
	admissionregistration "k8s.io/api/admissionregistration/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"log"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"strings"
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
			Name: consts.Webhook,
		},
		Webhooks: []admissionregistration.ValidatingWebhook{
			{
				Name: strings.Join([]string{"validator.experiment", consts.Webhook}, "."),
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
			Name: consts.Webhook,
		},
		Webhooks: []admissionregistration.MutatingWebhook{
			{
				Name: strings.Join([]string{"defaulter.experiment", consts.Webhook}, "."),
				ClientConfig: admissionregistration.WebhookClientConfig{
					CABundle: []byte("CG=="),
				},
			},
			{
				Name: strings.Join([]string{"mutator.pod", consts.Webhook}, "."),
				ClientConfig: admissionregistration.WebhookClientConfig{
					CABundle: []byte("CG=="),
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
	testControllerService := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      consts.Service,
			Namespace: testNamespace,
		},
	}

	tests := []struct {
		testDescription string
		err             bool
		objects         []client.Object
	}{
		{
			testDescription: "Generate successfully",
			err:             false,
			objects: []client.Object{
				testGeneratorJob,
				testValidatingWebhook,
				testMutatingWebhook,
				testControllerService,
			},
		},
		{
			testDescription: "There is old Secret, katib-webhook-cert",
			err:             false,
			objects: []client.Object{
				testGeneratorJob,
				testValidatingWebhook,
				testMutatingWebhook,
				oldWebhookCertSecret,
				testControllerService,
			},
		},
		{
			testDescription: "There is not Job, katib-cert-generator",
			err:             true,
			objects: []client.Object{
				testValidatingWebhook,
				testMutatingWebhook,
				testControllerService,
			},
		},
		{
			testDescription: "There is not ValidatingWebhookConfiguration",
			err:             true,
			objects: []client.Object{
				testGeneratorJob,
				testMutatingWebhook,
				testControllerService,
			},
		},
		{
			testDescription: "There is not MutatingWebhookConfiguration",
			err:             true,
			objects: []client.Object{
				testGeneratorJob,
				testValidatingWebhook,
				testControllerService,
			},
		},
		{
			testDescription: "There is no Service katib-controller",
			err:             true,
			objects: []client.Object{
				testGeneratorJob,
				testMutatingWebhook,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testDescription, func(t *testing.T) {
			if err := executeGeneratorCommand(test.objects, testNamespace); (err != nil) != test.err {
				t.Errorf("expected error: %v, got: '%v'\n", test.err, err)
			}
		})
	}

}

func executeGeneratorCommand(kubeResources []client.Object, namespace string) error {

	fakeClientBuilder := fake.NewClientBuilder().WithScheme(scheme.Scheme)
	if len(kubeResources) > 0 {
		for _, r := range kubeResources {
			fakeClientBuilder.WithObjects(r)
		}
	}
	cmd := NewGenerateCmd(fakeClientBuilder.Build())
	if err := cmd.Flags().Set("namespace", namespace); err != nil {
		log.Fatal(err)
	}

	return cmd.Execute()
}
