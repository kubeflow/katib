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
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

type certificates struct {
	certPem []byte
	keyPem  []byte
	cert    *x509.Certificate
	key     *rsa.PrivateKey
}

func (c *certificates) encode(rawKey, rawDer []byte) error {
	key := &bytes.Buffer{}
	if err := pem.Encode(key, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: rawKey}); err != nil {
		return err
	}
	c.keyPem = key.Bytes()

	cert := &bytes.Buffer{}
	if err := pem.Encode(cert, &pem.Block{Type: "CERTIFICATE", Bytes: rawDer}); err != nil {
		return err
	}
	c.certPem = cert.Bytes()
	return nil
}
