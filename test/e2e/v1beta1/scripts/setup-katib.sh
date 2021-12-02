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

# This shell script is used to setup Katib deployment.

set -o errexit
set -o nounset
set -o pipefail

echo "Start to install Katib"
echo "CLUSTER_NAME: ${CLUSTER_NAME}"
echo "AWS_REGION: ${AWS_REGION}"
echo "ECR_REGISTRY: ${ECR_REGISTRY}"
echo "VERSION: ${PULL_PULL_SHA}"

echo "Configuring kubeconfig.."
aws eks update-kubeconfig --region=${AWS_REGION} --name=${CLUSTER_NAME}
kubectl version
kubectl cluster-info

# Update Katib images with the current PULL SHA.
make update-images PREFIX="${ECR_REGISTRY}/${REPO_NAME}/v1beta1/" TAG="${PULL_PULL_SHA}"

echo -e "\n The Katib will be deployed with the following images"
cat "manifests/v1beta1/installs/katib-standalone/kustomization.yaml"
cat "manifests/v1beta1/components/controller/katib-config.yaml"

# Update Trial training container images.
make update-images PREFIX=""
./scripts/v1beta1/update-trial-images.sh -p "${ECR_REGISTRY}/${REPO_NAME}/v1beta1/trial-" -t ${VERSION}

echo "Creating Kubeflow namespace"
kubectl create namespace kubeflow

echo "Deploying training-operator from kubeflow/manifests v1.4 branch"
cd "${MANIFESTS_DIR}/apps/training-operator/upstream/overlays/kubeflow"
kustomize build . | kubectl apply -f -

echo "Deploying Katib"
cd "${GOPATH}/src/github.com/kubeflow/katib"
make deploy

# Wait until all Katib pods is running.
TIMEOUT=120s
kubectl wait --for=condition=complete --timeout=${TIMEOUT} -l katib.kubeflow.org/component=cert-generator -n kubeflow job
kubectl wait --for=condition=ready --timeout=${TIMEOUT} -l "katib.kubeflow.org/component in (controller,db-manager,mysql,ui)" -n kubeflow pod

echo "All Katib components are running."
echo "Katib deployments"
kubectl -n kubeflow get deploy
echo "Katib services"
kubectl -n kubeflow get svc
echo "Katib pods"
kubectl -n kubeflow get pod

# Check that Katib is working with 2 Experiments.
kubectl apply -f test/e2e/v1beta1/valid-experiment.yaml
kubectl delete -f test/e2e/v1beta1/valid-experiment.yaml

set +o errexit
kubectl apply -f test/e2e/v1beta1/invalid-experiment.yaml
if [ $? -ne 1 ]; then
  echo "Failed to create invalid-experiment: return code $?"
  exit 1
fi

# Build the binary for e2e test
echo "Building run-e2e-experiment for e2e test cases"
go build -o run-e2e-experiment test/e2e/v1beta1/run-e2e-experiment.go

exit 0
