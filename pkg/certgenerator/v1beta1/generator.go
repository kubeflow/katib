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
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	configv1beta1 "github.com/kubeflow/katib/pkg/apis/config/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

var (
	errServiceNotFound      = errors.New("unable to locate controller service")
	errCertCheckFail        = errors.New("failed to check if certs already exist")
	errCreateCertFail       = errors.New("failed to create certs")
	errCreateCertSecretFail = errors.New("failed to create secret embedded certs")
	errInjectCertError      = errors.New("failed to inject certs into WebhookConfigurations")
)

// CertGenerator is the manager to generate certs.
type CertGenerator struct {
	namespace   string
	serviceName string
	secretName  string
	kubeClient  client.Client
	certsReady  chan struct{}

	certs             *certificates
	fullServiceDomain string
}

var _ manager.Runnable = &CertGenerator{}
var _ manager.LeaderElectionRunnable = &CertGenerator{}

func (c *CertGenerator) Start(ctx context.Context) error {
	if err := c.generate(ctx); err != nil {
		return err
	}
	// Close a certsReady means start to register controllers to the manager.
	close(c.certsReady)
	return nil
}

func (c *CertGenerator) NeedLeaderElection() bool {
	return true
}

// AddToManager adds the cert-generator to the manager.
func AddToManager(mgr manager.Manager, config configv1beta1.CertGeneratorConfig, certsReady chan struct{}) error {
	return mgr.Add(&CertGenerator{
		namespace:   consts.DefaultKatibNamespace,
		serviceName: config.WebhookServiceName,
		secretName:  config.WebhookSecretName,
		kubeClient:  mgr.GetClient(),
		certsReady:  certsReady,
	})
}

// generate generates certificates for the admission webhooks.
func (c *CertGenerator) generate(ctx context.Context) error {
	controllerService := &corev1.Service{}
	if err := c.kubeClient.Get(ctx, client.ObjectKey{Name: c.serviceName, Namespace: c.namespace}, controllerService); err != nil {
		return fmt.Errorf("%w: %v", errServiceNotFound, err)
	}

	certExist, err := c.isCertExist(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", errCertCheckFail, err)
	}
	if !certExist {
		c.fullServiceDomain = strings.Join([]string{c.serviceName, c.namespace, "svc"}, ".")

		if err = c.createCert(); err != nil {
			return fmt.Errorf("%w: %v", errCreateCertFail, err)
		}
		if err = c.updateCertSecret(ctx); err != nil {
			return fmt.Errorf("%w: %v", errCreateCertSecretFail, err)
		}
	}
	if err = c.injectCert(ctx); err != nil {
		return fmt.Errorf("%w: %v", errInjectCertError, err)
	}
	return nil
}

// isCertExist checks if a secret embedded certs already exists.
// For example, it will return true if the katib-controller is created with enabled leader-election
// since another controller pod will create the secret.
func (c *CertGenerator) isCertExist(ctx context.Context) (bool, error) {
	secret := &corev1.Secret{}
	if err := c.kubeClient.Get(ctx, client.ObjectKey{Name: c.secretName, Namespace: c.namespace}, secret); err != nil {
		return false, err
	}
	key := secret.Data[serverKeyName]
	cert := secret.Data[serverCertName]
	if len(key) != 0 && len(cert) != 0 {
		c.certs = &certificates{
			keyPem:  key,
			certPem: cert,
		}
		return true, nil
	}
	return false, nil
}

// createCert creates the self-signed certificate and private key.
func (c *CertGenerator) createCert() error {
	now := time.Now()
	template := &x509.Certificate{
		SerialNumber: big.NewInt(0),
		Subject: pkix.Name{
			CommonName: c.fullServiceDomain,
		},
		DNSNames: []string{
			c.fullServiceDomain,
		},
		NotBefore:   now,
		NotAfter:    now.Add(24 * time.Hour * 365 * 10),
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	klog.Info("Generating self-signed public certificate and private key.")
	rawKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, rawKey.Public(), rawKey)
	if err != nil {
		return err
	}
	if c.certs, err = encode(rawKey, der); err != nil {
		return err
	}
	return nil
}

// updateCertSecret updates Secret embedded tls.key and tls.crt.
func (c *CertGenerator) updateCertSecret(ctx context.Context) error {
	secret := &corev1.Secret{}
	if err := c.kubeClient.Get(ctx, client.ObjectKey{Name: c.secretName, Namespace: c.namespace}, secret); err != nil {
		return err
	}
	newSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			UID:        secret.UID,
			Name:       secret.Name,
			Namespace:  secret.Namespace,
			Generation: secret.Generation,
		},
		TypeMeta: secret.TypeMeta,
	}
	if len(newSecret.APIVersion) == 0 {
		newSecret.APIVersion = corev1.SchemeGroupVersion.String()
	}
	if len(newSecret.Kind) == 0 {
		newSecret.Kind = "Secret"
	}
	newSecret.Data = map[string][]byte{
		serverKeyName:  c.certs.keyPem,
		serverCertName: c.certs.certPem,
	}
	return c.kubeClient.Patch(ctx, newSecret, client.Apply, client.FieldOwner(ssaFieldOwnerName), client.ForceOwnership)
}

// injectCert applies patch to ValidatingWebhookConfiguration and MutatingWebhookConfiguration.
func (c *CertGenerator) injectCert(ctx context.Context) error {
	vWebhookConfig := &admissionregistrationv1.ValidatingWebhookConfiguration{}
	if err := c.kubeClient.Get(ctx, client.ObjectKey{Name: Webhook}, vWebhookConfig); err != nil {
		return err
	}
	if !bytes.Equal(vWebhookConfig.Webhooks[0].ClientConfig.CABundle, c.certs.certPem) {
		newVWebhookConfig := &admissionregistrationv1.ValidatingWebhookConfiguration{
			ObjectMeta: metav1.ObjectMeta{
				UID:        vWebhookConfig.UID,
				Name:       vWebhookConfig.Name,
				Namespace:  vWebhookConfig.Namespace,
				Generation: vWebhookConfig.Generation,
			},
			TypeMeta: vWebhookConfig.TypeMeta,
		}
		if len(newVWebhookConfig.APIVersion) == 0 {
			newVWebhookConfig.APIVersion = admissionregistrationv1.SchemeGroupVersion.String()
		}
		if len(newVWebhookConfig.Kind) == 0 {
			newVWebhookConfig.Kind = "ValidatingWebhookConfiguration"
		}
		newVWebhookConfig.Webhooks = vWebhookConfig.Webhooks
		newVWebhookConfig.Webhooks[0].ClientConfig.CABundle = c.certs.certPem

		klog.Info("Trying to patch ValidatingWebhookConfiguration adding the caBundle.")
		err := c.kubeClient.Patch(ctx, newVWebhookConfig, client.Apply, client.FieldOwner(ssaFieldOwnerName), client.ForceOwnership)
		if err != nil {
			klog.Errorf("Unable to patch ValidatingWebhookConfiguration %q", Webhook)
			return err
		}
	}

	mWebhookConfig := &admissionregistrationv1.MutatingWebhookConfiguration{}
	if err := c.kubeClient.Get(ctx, client.ObjectKey{Name: Webhook}, mWebhookConfig); err != nil {
		return err
	}
	if !bytes.Equal(mWebhookConfig.Webhooks[0].ClientConfig.CABundle, c.certs.certPem) ||
		!bytes.Equal(mWebhookConfig.Webhooks[1].ClientConfig.CABundle, c.certs.certPem) {
		newMWebhookConfig := &admissionregistrationv1.MutatingWebhookConfiguration{
			ObjectMeta: metav1.ObjectMeta{
				UID:        mWebhookConfig.UID,
				Name:       mWebhookConfig.Name,
				Namespace:  mWebhookConfig.Namespace,
				Generation: mWebhookConfig.Generation,
			},
			TypeMeta: mWebhookConfig.TypeMeta,
		}
		if len(newMWebhookConfig.APIVersion) == 0 {
			newMWebhookConfig.APIVersion = admissionregistrationv1.SchemeGroupVersion.String()
		}
		if len(newMWebhookConfig.Kind) == 0 {
			newMWebhookConfig.Kind = "MutatingWebhookConfiguration"
		}
		newMWebhookConfig.Webhooks = mWebhookConfig.Webhooks
		newMWebhookConfig.Webhooks[0].ClientConfig.CABundle = c.certs.certPem
		newMWebhookConfig.Webhooks[1].ClientConfig.CABundle = c.certs.certPem

		klog.Info("Trying to patch MutatingWebhookConfiguration adding the caBundle.")
		err := c.kubeClient.Patch(ctx, newMWebhookConfig, client.Apply, client.FieldOwner(ssaFieldOwnerName), client.ForceOwnership)
		if err != nil {
			klog.Errorf("Unable to patch MutatingWebhookConfiguration %q", Webhook)
			return err
		}
	}
	return nil
}
