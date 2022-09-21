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
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/kubeflow/katib/pkg/cert-generator/v1beta1/consts"
	"github.com/spf13/cobra"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	"math/big"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"time"
)

// generateOptions contains values for all certificates.
type generateOptions struct {
	namespace         string
	serviceName       string
	jobName           string
	fullServiceDomain string
}

// NewGenerateCmd sets up `generate` subcommand.
func NewGenerateCmd(kubeClient client.Client) *cobra.Command {
	o := &generateOptions{}
	cmd := &cobra.Command{
		Use:          "generate",
		Short:        "generate server cert for webhook",
		Long:         "generate server cert for webhook",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.run(context.TODO(), kubeClient); err != nil {
				return err
			}
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVarP(&o.namespace, "namespace", "n", "kubeflow", "set namespace")
	f.StringVarP(&o.jobName, "jobName", "j", consts.JobName, "set job name")
	f.StringVarP(&o.serviceName, "serviceName", "s", consts.Service, "set service name")
	return cmd
}

// run is main function for `generate` subcommand.
func (o *generateOptions) run(ctx context.Context, kubeClient client.Client) error {
	controllerService := &corev1.Service{}
	if err := kubeClient.Get(ctx, client.ObjectKey{Namespace: o.namespace, Name: o.serviceName}, controllerService); err != nil {
		klog.Errorf("Unable to locate controller service: %s", o.serviceName)
		return err
	}

	o.fullServiceDomain = strings.Join([]string{o.serviceName, o.namespace, "svc"}, ".")

	caKeyPair, err := o.createCACert()
	if err != nil {
		return err
	}
	keyPair, err := o.createCert(caKeyPair)
	if err != nil {
		return err
	}

	if err = o.createWebhookCertSecret(ctx, kubeClient, caKeyPair, keyPair); err != nil {
		return err
	}
	if err = o.injectCert(ctx, kubeClient, caKeyPair); err != nil {
		return err
	}

	return nil
}

// createCACert creates the self-signed CA certificate and private key.
func (o *generateOptions) createCACert() (*certificates, error) {
	now := time.Now()
	template := &x509.Certificate{
		SerialNumber: big.NewInt(0),
		Subject: pkix.Name{
			CommonName:   consts.CAName,
			Organization: []string{consts.Katib},
		},
		DNSNames: []string{
			consts.CAName,
		},
		NotBefore:             now,
		NotAfter:              now.Add(24 * time.Hour * 365 * 10),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	klog.Info("Generating the self-signed CA certificate and private key.")
	rawKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, rawKey.Public(), rawKey)
	if err != nil {
		return nil, err
	}

	return encode(rawKey, der)
}

// createCert creates public certificate and private key signed with self-signed CA certificate and private key.
func (o *generateOptions) createCert(caKeyPair *certificates) (*certificates, error) {
	now := time.Now()
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: o.fullServiceDomain,
		},
		DNSNames: []string{
			o.serviceName,
			strings.Join([]string{o.serviceName, o.namespace}, "."),
			o.fullServiceDomain,
		},
		NotBefore:             now,
		NotAfter:              now.Add(24 * time.Hour * 365 * 10),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: false,
	}

	klog.Info("Generating public certificate and private key signed with self-singed CA cert and private key.")
	rawKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	der, err := x509.CreateCertificate(rand.Reader, template, caKeyPair.cert, rawKey.Public(), caKeyPair.key)
	if err != nil {
		return nil, err
	}

	return encode(rawKey, der)
}

// createWebhookCertSecret creates Secret embedded ca.key, ca.crt, tls.key and tls.crt.
func (o *generateOptions) createWebhookCertSecret(ctx context.Context, kubeClient client.Client, caKeyPair *certificates, keyPair *certificates) error {

	certGeneratorJob := &batchv1.Job{}
	if err := kubeClient.Get(ctx, client.ObjectKey{Namespace: o.namespace, Name: o.jobName}, certGeneratorJob); err != nil {
		return err
	}

	// Create secret with CA cert and server cert/key.
	// Add ownerReferences to clean-up secret with cert generator Job.
	isController := true
	jobUID := certGeneratorJob.UID
	webhookCertSecret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      consts.Secret,
			Namespace: o.namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "batch/v1",
					Kind:       "Job",
					Controller: &isController,
					Name:       o.jobName,
					UID:        jobUID,
				},
			},
		},
		Type: corev1.SecretTypeTLS,
		Data: map[string][]byte{
			"ca.key":  caKeyPair.keyPem,
			"ca.crt":  caKeyPair.certPem,
			"tls.key": keyPair.keyPem,
			"tls.crt": keyPair.certPem,
		},
	}

	oldSecret := &corev1.Secret{}
	err := kubeClient.Get(ctx, client.ObjectKey{Namespace: o.namespace, Name: consts.Secret}, oldSecret)
	switch {
	case err != nil && !k8serrors.IsNotFound(err):
		return err
	case err == nil:
		klog.Warning("Previous secret was found and removed.")
		if err = kubeClient.Delete(ctx, oldSecret); err != nil {
			return err
		}
	}

	klog.Infof("Creating Secret: %s", consts.Secret)
	if err = kubeClient.Create(ctx, webhookCertSecret); err != nil {
		return err
	}
	return nil
}

// injectCert applies patch to ValidatingWebhookConfiguration and MutatingWebhookConfiguration.
func (o *generateOptions) injectCert(ctx context.Context, kubeClient client.Client, caKeypair *certificates) error {
	validatingConf := &admissionregistrationv1.ValidatingWebhookConfiguration{}
	if err := kubeClient.Get(ctx, client.ObjectKey{Name: consts.Webhook}, validatingConf); err != nil {
		return err
	}
	newValidatingConf := validatingConf.DeepCopy()
	newValidatingConf.Webhooks[0].ClientConfig.CABundle = caKeypair.certPem

	klog.Info("Trying to patch ValidatingWebhookConfiguration adding the caBundle.")
	if err := kubeClient.Patch(ctx, newValidatingConf, client.MergeFrom(validatingConf)); err != nil {
		klog.Errorf("Unable to patch ValidatingWebhookConfiguration %s", consts.Webhook)
		return err
	}

	mutatingConf := &admissionregistrationv1.MutatingWebhookConfiguration{}
	if err := kubeClient.Get(ctx, client.ObjectKey{Name: consts.Webhook}, mutatingConf); err != nil {
		return err
	}
	newMutatingConf := mutatingConf.DeepCopy()
	newMutatingConf.Webhooks[0].ClientConfig.CABundle = caKeypair.certPem
	newMutatingConf.Webhooks[1].ClientConfig.CABundle = caKeypair.certPem

	klog.Info("Trying to patch MutatingWebhookConfiguration adding the caBundle.")
	if err := kubeClient.Patch(ctx, newMutatingConf, client.MergeFrom(mutatingConf)); err != nil {
		klog.Errorf("Unable to patch MutatingWebhookConfiguration %s", consts.Webhook)
		return err
	}

	return nil
}
