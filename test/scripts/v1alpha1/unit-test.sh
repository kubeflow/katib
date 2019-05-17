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

gcloud --project ${PROJECT} container clusters get-credentials ${CLUSTER_NAME} \
  --zone ${ZONE}
kubectl config set-context $(kubectl config current-context) --namespace=default

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

echo "Copy source to GOPATH"
mkdir -p ${GO_DIR}
cp -r cmd ${GO_DIR}/cmd
cp -r pkg ${GO_DIR}/pkg
cp -r vendor ${GO_DIR}/vendor

echo "Run unit test cases"
cd ${GO_DIR}

os=$(go env GOOS)
arch=$(go env GOARCH)
version=1.0.2
echo "os: ${os}, arch: ${arch}"

# download the release
curl -L -O "https://github.com/kubernetes-sigs/kubebuilder/releases/download/v${version}/kubebuilder_${version}_${os}_${arch}.tar.gz"

# extract the archive
tar -zxvf kubebuilder_${version}_${os}_${arch}.tar.gz
mv kubebuilder_${version}_${os}_${arch} /usr/local/bin/kubebuilder
export PATH=$PATH:/usr/local/bin

go test ./...
cd - > /dev/null
