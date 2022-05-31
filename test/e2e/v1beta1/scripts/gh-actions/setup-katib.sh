#!/usr/bin/env bash

# Copyright 2022 The Kubeflow Authors.
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

E2E_TEST_IMAGE_TAG="e2e-test"
TRAINING_OPERATOR_VERSION="v1.4.0"
DEPLOY_TRAINING_OPERATOR=$1
cd "$(dirname "$0")"

echo "Start to install Katib"
kubectl cluster-info
kubectl get nodes

# Update Katib images with `e2e-test`.
cd ../../../../../ && make update-images OLD_PREFIX="docker.io/kubeflowkatib/" NEW_PREFIX="docker.io/kubeflowkatib/" TAG="$E2E_TEST_IMAGE_TAG" && cd -

echo -e "\n The Katib will be deployed with the following configs"
cat "manifests/v1beta1/installs/katib-standalone/kustomization.yaml"
cat "manifests/v1beta1/components/controller/katib-config.yaml"

echo "Creating Kubeflow namespace"
kubectl create namespace kubeflow

if "$DEPLOY_TRAINING_OPERATOR"; then
  echo "Deploying Training Operator v1.4.0"
  kubectl apply -k "github.com/kubeflow/training-operator/manifests/overlays/standalone?ref=$TRAINING_OPERATOR_VERSION"
fi

echo "Deploying Katib"
cd ../../../../../ && make deploy && cd -

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
kubectl apply -f ../testdata/valid-experiment.yaml
kubectl delete -f ../testdata/valid-experiment.yaml

set +o errexit
kubectl apply -f ../testdata/invalid-experiment.yaml
if [ $? -ne 1 ]; then
  echo "Failed to create invalid-experiment: return code $?"
  exit 1
fi
set -o errexit

# Build the binary for e2e test
echo "Building run-e2e-experiment for e2e test cases"
mkdir -p ../bin
go build -o ../bin/run-e2e-experiment ../hack/gh-actions/run-e2e-experiment.go
chmod +x ../bin/run-e2e-experiment

exit 0
