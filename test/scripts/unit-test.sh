#!/bin/bash

# Copyright 2018 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This shell script is used to build a cluster and create a namespace from our
# argo workflow

set -o errexit
set -o nounset
set -o pipefail

export PATH=${GOPATH}/bin:/usr/local/go/bin:${PATH}
REGISTRY="${GCP_REGISTRY}"
CLUSTER_NAME="${CLUSTER_NAME}"
ZONE="${GCP_ZONE}"
PROJECT="${GCP_PROJECT}"
GO_DIR=${GOPATH}/src/github.com/${REPO_OWNER}/${REPO_NAME}
VERSION=$(git describe --tags --always --dirty)

echo "Activating service-account"
gcloud auth activate-service-account --key-file=${GOOGLE_APPLICATION_CREDENTIALS}

echo "Configuring kubectl"

gcloud container clusters describe ${CLUSTER_NAME} \
  --zone ${ZONE} \
  --format 'value(masterAuth.clusterCaCertificate)'|  base64 -d > ca.pem

gcloud container clusters describe ${CLUSTER_NAME} \
  --zone ${ZONE} \
  --format 'value(masterAuth.clientCertificate)'  |  base64 -d > client.pem

gcloud container clusters describe ${CLUSTER_NAME} \
  --zone ${ZONE} \
  --format 'value(masterAuth.clientKey)' |  base64 -d > key.rsa

kubectl config set-credentials temp-admin --username=admin --client-certificate=./client.pem --client-key=./key.rsa
kubectl config set-context temp-context --cluster=$(kubectl config get-clusters | grep ${CLUSTER_NAME}) --user=temp-admin
kubectl config use-context temp-context

kubectl delete pod mysql-ut --ignore-not-found
kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  labels:
    instance: mysql-ut
  name: mysql-ut
spec:
  containers:
  - name: mysql
    image: mysql
    env:
    - name: MYSQL_ROOT_PASSWORD
      value: test123
    - name: MYSQL_DATABASE
      value: vizier
EOF

TIMEOUT=120
until kubectl logs -l instance=mysql-ut |grep 'init process done'; do
    sleep 10
    TIMEOUT=$(( TIMEOUT - 1 ))
    if [ "$TIMEOUT" -eq 0 ]; then
	echo "DB failed to start"
	kubectl get pod -l instance=mysql-ut
	exit 1
    fi
done

kubectl port-forward $(kubectl get pod -l instance=mysql-ut -o=name) 3306:3306&
export TEST_MYSQL=localhost

echo "Run unit test cases"
cd ${GO_DIR}
go test ./...
cd - > /dev/null
