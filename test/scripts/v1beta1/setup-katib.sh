#!/bin/bash

# Copyright 2020 The Kubeflow Authors.
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

# This shell script is used to setup Katib deployment.

set -o errexit
set -o nounset
set -o pipefail

CLUSTER_NAME="${CLUSTER_NAME}"
AWS_REGION="${AWS_REGION}"
ECR_REGISTRY="${ECR_REGISTRY}"
VERSION="${PULL_BASE_SHA}"

echo "Start to install Katib"
echo "CLUSTER_NAME: ${CLUSTER_NAME}"
echo "AWS_REGION: ${AWS_REGION}"
echo "ECR_REGISTRY: ${ECR_REGISTRY}"
echo "VERSION: ${PULL_BASE_SHA}"

echo "Configuring kubeconfig.."
aws eks update-kubeconfig --region=${AWS_REGION} --name=${CLUSTER_NAME}

# Update images with current pull base sha.
echo "Updating Katib images with current PR SHA: ${VERSION}"
# Katib controller
sed -i -e "s@image: gcr.io\/kubeflow-images-public\/katib\/v1beta1\/katib-controller@image: ${ECR_REGISTRY}\/${REPO_NAME}\/v1beta1\/katib-controller:${VERSION}@" manifests/v1beta1/katib-controller/katib-controller.yaml

# Metrics collector
sed -i -e "s@gcr.io\/kubeflow-images-public\/katib\/v1beta1\/file-metrics-collector@${ECR_REGISTRY}\/${REPO_NAME}\/v1beta1\/file-metrics-collector:${VERSION}@" manifests/v1beta1/katib-controller/katib-config.yaml
sed -i -e "s@gcr.io\/kubeflow-images-public\/katib\/v1beta1\/tfevent-metrics-collector@${ECR_REGISTRY}\/${REPO_NAME}\/v1beta1\/tfevent-metrics-collector:${VERSION}@" manifests/v1beta1/katib-controller/katib-config.yaml

# Katib DB manager
sed -i -e "s@image: gcr.io\/kubeflow-images-public\/katib\/v1beta1\/katib-db-manager@image: ${ECR_REGISTRY}\/${REPO_NAME}\/v1beta1\/katib-db-manager:${VERSION}@" manifests/v1beta1/db-manager/deployment.yaml

# UI
sed -i -e "s@image: gcr.io\/kubeflow-images-public\/katib\/v1beta1\/katib-ui@image: ${ECR_REGISTRY}\/${REPO_NAME}\/v1beta1\/katib-ui:${VERSION}@" manifests/v1beta1/ui/deployment.yaml

# Suggestion algorithms
sed -i -e "s@gcr.io\/kubeflow-images-public\/katib\/v1beta1\/suggestion-enas@${ECR_REGISTRY}\/${REPO_NAME}\/v1beta1\/suggestion-enas:${VERSION}@" manifests/v1beta1/katib-controller/katib-config.yaml
sed -i -e "s@gcr.io\/kubeflow-images-public\/katib\/v1beta1\/suggestion-hyperband@${ECR_REGISTRY}\/${REPO_NAME}\/v1beta1\/suggestion-hyperband:${VERSION}@" manifests/v1beta1/katib-controller/katib-config.yaml
sed -i -e "s@gcr.io\/kubeflow-images-public\/katib\/v1beta1\/suggestion-chocolate@${ECR_REGISTRY}\/${REPO_NAME}\/v1beta1\/suggestion-chocolate:${VERSION}@" manifests/v1beta1/katib-controller/katib-config.yaml
sed -i -e "s@gcr.io\/kubeflow-images-public\/katib\/v1beta1\/suggestion-hyperopt@${ECR_REGISTRY}\/${REPO_NAME}\/v1beta1\/suggestion-hyperopt:${VERSION}@" manifests/v1beta1/katib-controller/katib-config.yaml
sed -i -e "s@gcr.io\/kubeflow-images-public\/katib\/v1beta1\/suggestion-skopt@${ECR_REGISTRY}\/${REPO_NAME}\/v1beta1\/suggestion-skopt:${VERSION}@" manifests/v1beta1/katib-controller/katib-config.yaml
sed -i -e "s@gcr.io\/kubeflow-images-public\/katib\/v1beta1\/suggestion-goptuna@${ECR_REGISTRY}\/${REPO_NAME}\/v1beta1\/suggestion-goptuna:${VERSION}@" manifests/v1beta1/katib-controller/katib-config.yaml
sed -i -e "s@gcr.io\/kubeflow-images-public\/katib\/v1beta1\/suggestion-darts@${ECR_REGISTRY}\/${REPO_NAME}\/v1beta1\/suggestion-darts:${VERSION}@" manifests/v1beta1/katib-controller/katib-config.yaml

cat manifests/v1beta1/katib-controller/katib-config.yaml

echo "Deploying tf-operator from kubeflow/manifests master"
cd "${MANIFESTS_DIR}/tf-training/tf-job-crds/base"
kustomize build . | kubectl apply -f -
cd "${MANIFESTS_DIR}/tf-training/tf-job-operator/base"
kustomize build . | kubectl apply -f -

echo "Deploying pytorch-operator from kubeflow/manifests master"
cd "${MANIFESTS_DIR}/pytorch-job/pytorch-job-crds/base"
kustomize build . | kubectl apply -f -
cd "${MANIFESTS_DIR}/pytorch-job/pytorch-operator/base/"
kustomize build . | kubectl apply -f -

echo "Deploying Katib"
cd "${GOPATH}/src/github.com/kubeflow/katib"
make deploy

# Wait until all Katib pods is running.
TIMEOUT=120
PODNUM=$(kubectl get deploy -n kubeflow | grep -v NAME | wc -l)
until kubectl get pods -n kubeflow | grep Running | [[ $(wc -l) -eq $PODNUM ]]; do
  echo Pod Status $(kubectl get pods -n kubeflow | grep "1/1" | wc -l)/$PODNUM

  sleep 10
  TIMEOUT=$((TIMEOUT - 1))
  if [[ $TIMEOUT -eq 0 ]]; then
    echo "NG"
    kubectl get pods -n kubeflow
    exit 1
  fi
done

echo "All Katib components are running."
kubectl version
kubectl cluster-info
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

exit 0
