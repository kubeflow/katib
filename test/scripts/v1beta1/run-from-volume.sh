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
GO_DIR=${GOPATH}/src/github.com/${REPO_OWNER}/${REPO_NAME}

# Activate gcloud service account
source test/scripts/v1beta1/utils.sh
_activate_service_account

echo "Configuring kubectl"

echo "CLUSTER_NAME: ${CLUSTER_NAME}"
echo "ZONE: ${GCP_ZONE}"
echo "PROJECT: ${GCP_PROJECT}"

gcloud --project ${PROJECT} container clusters get-credentials ${CLUSTER_NAME} \
  --zone ${ZONE}
kubectl config set-context $(kubectl config current-context) --namespace=default
USER=$(gcloud config get-value account)

echo "All Katib components are running."
kubectl version
kubectl cluster-info
echo "Katib deployments"
kubectl -n kubeflow get deploy
echo "Katib services"
kubectl -n kubeflow get svc
echo "Katib pods"
kubectl -n kubeflow get pod

cd ${GO_DIR}/test/e2e/v1beta1

# Set number of epochs to 2 for faster execution
sed -i -e "s@--batch-size=64@--num-epochs=2@" ../../../examples/v1beta1/resume-experiment/from-volume-resume.yaml

echo "Running e2e test for resume from volume experiment"
export KUBECONFIG=$HOME/.kube/config
./run-e2e-experiment ../../../examples/v1beta1/resume-experiment/from-volume-resume.yaml

kubectl -n kubeflow describe suggestion from-volume-resume
kubectl -n kubeflow describe experiment from-volume-resume

echo "Available volumes"
kubectl get pvc -n kubeflow
kubectl get pv

echo "Resuming the completed experiment with resume from volume"
./resume-e2e-experiment ../../../examples/v1beta1/resume-experiment/from-volume-resume.yaml

kubectl -n kubeflow describe suggestion from-volume-resume
kubectl -n kubeflow describe experiment from-volume-resume

kubectl -n kubeflow delete experiment from-volume-resume

exit 0
