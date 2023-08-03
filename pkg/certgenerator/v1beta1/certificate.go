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
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

// certificates contains all certificates for katib-webhook-cert.
type certificates struct {
	certPem []byte
	keyPem  []byte
	cert    *x509.Certificate
	key     *rsa.PrivateKey
}

// encode creates PEM key and convert DER to CRT.
func encode(rawKey *rsa.PrivateKey, der []byte) (*certificates, error) {
	keyPem := &bytes.Buffer{}
	if err := pem.Encode(keyPem, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(rawKey)}); err != nil {
		return nil, err
	}

	certPem := &bytes.Buffer{}
	if err := pem.Encode(certPem, &pem.Block{Type: "CERTIFICATE", Bytes: der}); err != nil {
		return nil, err
	}

	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, err
	}

	return &certificates{
		certPem: certPem.Bytes(),
		keyPem:  keyPem.Bytes(),
		cert:    cert,
		key:     rawKey,
	}, nil
}
