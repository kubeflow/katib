#!/bin/bash

set -e

service="katib-controller"
webhook="katib.kubeflow.org"
secret="katib-webhook-cert"
namespace="kubeflow"
fullServiceDomain="${service}.${namespace}.svc"

# Fully qualified name of the CSR object.
csr="certificatesigningrequests.v1beta1.certificates.k8s.io"

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
DNS.3 = ${fullServiceDomain}
EOF

openssl genrsa -out "${tmpdir}/server-key.pem" 2048
openssl req -new -key "${tmpdir}/server-key.pem" -subj "/CN=${fullServiceDomain}" -out "${tmpdir}/server.csr" -config "${tmpdir}/csr.conf"

csrName=${service}.${namespace}
echo "INFO: Creating CSR: ${csrName} "

# Clean-up any previously created CSR for our service.
# Ignore errors.
set +e
if kubectl get "${csr}/${csrName}"; then
  if kubectl delete "${csr}/${csrName}"; then
    echo "WARN: Previous CSR was found and removed."
  fi
fi

# Create server cert/key CSR and send it to k8s api.
set -e
cat <<EOF | kubectl create --validate=false -f -
apiVersion: certificates.k8s.io/v1beta1
kind: CertificateSigningRequest
metadata:
  name: ${csrName}
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
# Ignore errors.
set +e
while true; do
  if kubectl get "${csr}/${csrName}"; then
    break
  fi
done

# Approve and fetch the signed certificate.
set -e
kubectl certificate approve "${csr}/${csrName}"

# Verify that certificate has been signed.
set +e
i=1
while [ "$i" -ne 20 ]; do
  serverCert=$(kubectl get "${csr}/${csrName}" -o jsonpath='{.status.certificate}')
  if [ "${serverCert}" != '' ]; then
    break
  fi
  sleep 3
  i=$((i + 1))
done

set -e
if [ "${serverCert}" = '' ]; then
  echo "ERROR: After approving csr ${csrName}, the signed certificate did not appear on the resource. Giving up after 1 minute." >&2
  exit 1
fi

echo "${serverCert}" | openssl base64 -d -A -out "${tmpdir}/server-cert.pem"

# Clean-up any previously created secret.
# Ignore errors.
set +e
if kubectl get secret ${secret} -n ${namespace}; then
  if kubectl delete secret ${secret} -n ${namespace}; then
    echo "WARN: Previous secret was found and removed."
  fi
fi
# Create the secret with CA cert and server cert/key.
kubectl create secret tls "${secret}" \
  --key="${tmpdir}/server-key.pem" \
  --cert="${tmpdir}/server-cert.pem" \
  --dry-run -o yaml |
  kubectl -n "${namespace}" apply -f -

caBundle=$(base64 </run/secrets/kubernetes.io/serviceaccount/ca.crt | tr -d '\n')
echo "INFO: Encoded CA:"
echo -e "${caBundle} \n"

patch_webhook() {
  set +e
  kind=$1
  webhook=$2
  path=$3
  caBundle=$4

  i=0
  while [[ $i -lt 5 ]]; do
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

# Patch the webhook to add the caBundle.
patch_webhook "ValidatingWebhookConfiguration" ${webhook} "/webhooks/0/clientConfig/caBundle" ${caBundle}
patch_webhook "MutatingWebhookConfiguration" ${webhook} "/webhooks/0/clientConfig/caBundle" ${caBundle}
patch_webhook "MutatingWebhookConfiguration" ${webhook} "/webhooks/1/clientConfig/caBundle" ${caBundle}
