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

CLUSTER_NAME="${CLUSTER_NAME}"
ZONE="${GCP_ZONE}"
PROJECT="${GCP_PROJECT}"
NAMESPACE="${DEPLOY_NAMESPACE}"
REGISTRY="${GCP_REGISTRY}"
VERSION=$(git describe --tags --always --dirty)
GO_DIR=${GOPATH}/src/github.com/${REPO_OWNER}/${REPO_NAME}

echo "Activating service-account"
gcloud auth activate-service-account --key-file=${GOOGLE_APPLICATION_CREDENTIALS}

echo "Configuring kubectl"
gcloud --project ${PROJECT} container clusters get-credentials ${CLUSTER_NAME} \
    --zone ${ZONE}


# Create serviceaccount to overcome the RBAC issue in GKE.
# Please see https://github.com/coreos/prometheus-operator/pull/360/files?short_path=a141d34#diff-a141d34d5164b9b18482d4e57c88ec97
echo "Create clusterrolebining"
kubectl apply -f manifests/0-namespace.yaml
kubectl create serviceaccount katib -n katib
kubectl create clusterrolebinding katib-cluster-admin-binding --clusterrole=cluster-admin --serviceaccount=katib:katib

echo "Install Katib "
sed -i -e "s@image: katib\/vizier-core@image: ${REGISTRY}\/${REPO_NAME}\/vizier-core:${VERSION}@" manifests/vizier/core/deployment.yaml
sed -i -e "s@image: katib\/suggestion-random@image: ${REGISTRY}\/${REPO_NAME}\/suggestion-random:${VERSION}@" manifests/vizier/suggestion/random/deployment.yaml
sed -i -e "s@image: katib\/suggestion-grid@image: ${REGISTRY}\/${REPO_NAME}\/suggestion-grid:${VERSION}@" manifests/vizier/suggestion/grid/deployment.yaml
sed -i -e "s@image: katib\/suggestion-hyperband@image: ${REGISTRY}\/${REPO_NAME}\/suggestion-hyperband:${VERSION}@" manifests/vizier/suggestion/hyperband/deployment.yaml
sed -i -e "s@image: katib\/dlk-manager@image: ${REGISTRY}\/${REPO_NAME}\/dlk-manager:${VERSION}@" manifests/dlk/deployment.yaml
sed -i -e "s@image: katib\/katib-frontend@image: ${REGISTRY}\/${REPO_NAME}\/katib-frontend:${VERSION}@" manifests/modeldb/frontend/deployment.yaml
cd ${GO_DIR}
./deploy.sh

TIMEOUT=60
until kubectl get pods -n katib | grep 1/1 | [[ $(wc -l) -eq 7 ]]; do
    sleep 10
    TIMEOUT=$(( TIMEOUT - 1 ))
    if [[ $TIMEOUT -eq 0 ]];then
        echo "NG"
        exit 1
    fi
done
echo "OK"

#echo "Run go tests"
#cd ${GO_DIR}
#go run ./test/e2e/main.go --namespace=${NAMESPACE}
