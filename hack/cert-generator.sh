#!/bin/bash

# Copyright 2021 The Kubeflow Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

service="katib-controller"
namespace=${KATIB_CORE_NAMESPACE}
full_service_domain="${service}.${namespace}.svc"

job_name="katib-cert-generator"
secret="katib-webhook-cert"
webhook="katib.kubeflow.org"

# Fully qualified name of the CSR object.
csr="certificatesigningrequests.v1beta1.certificates.k8s.io"
# CSR name.
csr_name=${service}.${namespace}

if [ ! -x "$(command -v openssl)" ]; then
  echo "ERROR: openssl not found"
  exit 1
fi

tmpdir=$(mktemp -d)
echo "INFO: Creating certs in tmpdir ${tmpdir} "

cat <<EOF >>"${tmpdir}/csr.conf"
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1 = ${service}
DNS.2 = ${service}.${namespace}
DNS.3 = ${full_service_domain}
EOF

openssl genrsa -out "${tmpdir}/server-key.pem" 2048
openssl req -new -key "${tmpdir}/server-key.pem" -subj "/CN=system:node:${full_service_domain}/O=system:nodes" -out "${tmpdir}/server.csr" -config "${tmpdir}/csr.conf"

# Clean-up any previously created CSR for our service.
set +e
if kubectl get "${csr}/${csr_name}"; then
  if kubectl delete "${csr}/${csr_name}"; then
    echo "WARN: Previous CSR was found and removed."
  fi
fi

# Create server cert/key CSR and send it to k8s api.
set -e
echo "INFO: Creating CSR: ${csr_name}"

# signerName is not supported in Kubernetes <= 1.17
# See: https://github.com/kubeflow/katib/issues/1500
cat <<EOF | kubectl create --validate=false -f -
apiVersion: certificates.k8s.io/v1beta1
kind: CertificateSigningRequest
metadata:
  name: ${csr_name}
spec:
  groups:
  - system:authenticated
  request: $(base64 <"${tmpdir}/server.csr" | tr -d '\n')
  signerName: kubernetes.io/kubelet-serving
  usages:
  - digital signature
  - key encipherment
  - server auth
EOF

# Verify that CSR has been created.
set +e
while true; do
  if kubectl get "${csr}/${csr_name}"; then
    break
  fi
done

# Approve and fetch the signed certificate.
set -e
kubectl certificate approve "${csr}/${csr_name}"

# Verify that certificate has been signed.
set +e
i=1
while [ "$i" -ne 20 ]; do
  server_cert=$(kubectl get "${csr}/${csr_name}" -o jsonpath='{.status.certificate}')
  if [ "${server_cert}" != '' ]; then
    break
  fi
  sleep 3
  i=$((i + 1))
done

set -e
if [ "${server_cert}" = '' ]; then
  echo "ERROR: After approving csr ${csr_name}, the signed certificate did not appear on the resource. Giving up after 1 minute."
  exit 1
fi

echo "${server_cert}" | openssl base64 -d -A -out "${tmpdir}/server-cert.pem"

# Clean-up any previously created secret.
set +e
if kubectl get secret ${secret} -n ${namespace}; then
  if kubectl delete secret ${secret} -n ${namespace}; then
    echo "WARN: Previous secret was found and removed."
  fi
fi

# Get cert generator Job UID.
set -e
job_uid=$(kubectl get job ${job_name} -n ${namespace} -o jsonpath='{.metadata.uid}')

# Create secret with CA cert and server cert/key.
# Add ownerReferences to clean-up secret with cert generator Job.
echo "INFO: Creating Secret: ${secret}"
cat <<EOF | kubectl create -f -
apiVersion: v1
kind: Secret
metadata:
  name: ${secret}
  namespace: ${namespace}
  ownerReferences:
    - apiVersion: batch/v1
      kind: Job
      controller: true
      name: ${job_name}
      uid: ${job_uid}
type: kubernetes.io/tls
data:
  tls.key: $(base64 <"${tmpdir}/server-key.pem" | tr -d '\n')
  tls.crt: $(base64 <"${tmpdir}/server-cert.pem" | tr -d '\n')
EOF

patch_webhook() {
  set +e
  kind=$1
  webhook=$2
  path=$3
  caBundle=$4

  i=0
  while [[ $i -ne 5 ]]; do
    echo "INFO: Trying to patch ${kind} adding the caBundle."
    if kubectl patch ${kind} ${webhook} --type='json' -p "[{'op': 'replace', 'path': '${path}', 'value':'${caBundle}'}]"; then
      break
    fi
    echo "WARNING: Webhook are not patched. Retrying in 5s..."
    sleep 5
    i=$((i + 1))
  done
  if [[ $i -eq 5 ]]; then
    echo "ERROR: Unable to patch ${kind} ${webhook}"
    exit 1
  fi
}

caBundle=$(base64 </run/secrets/kubernetes.io/serviceaccount/ca.crt | tr -d '\n')
echo "INFO: Encoded CA:"
echo -e "${caBundle} \n"

# Patch the webhook to add the caBundle.
patch_webhook "ValidatingWebhookConfiguration" ${webhook} "/webhooks/0/clientConfig/caBundle" ${caBundle}
patch_webhook "MutatingWebhookConfiguration" ${webhook} "/webhooks/0/clientConfig/caBundle" ${caBundle}
patch_webhook "MutatingWebhookConfiguration" ${webhook} "/webhooks/1/clientConfig/caBundle" ${caBundle}
