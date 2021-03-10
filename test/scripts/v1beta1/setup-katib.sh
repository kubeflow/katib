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
kubectl version
kubectl cluster-info

# Update images with current pull base sha.
echo "Updating Katib images with the current PR SHA: ${VERSION}"
FILE_PATH="manifests/v1beta1/installs/katib-standalone/kustomization.yaml"

# Change tag to all images in kustomization file.
sed -i -e "s@latest@${VERSION}@" ${FILE_PATH}

# Change Katib controller image.
sed -i -e "s@newName: docker.io/kubeflowkatib/katib-controller@newName: ${ECR_REGISTRY}/${REPO_NAME}/v1beta1/katib-controller@" ${FILE_PATH}

# Change Katib DB manager image.
sed -i -e "s@newName: docker.io/kubeflowkatib/katib-db-manager@newName: ${ECR_REGISTRY}/${REPO_NAME}/v1beta1/katib-db-manager@" ${FILE_PATH}

# Change Katib UI image.
sed -i -e "s@newName: docker.io/kubeflowkatib/katib-ui@newName: ${ECR_REGISTRY}/${REPO_NAME}/v1beta1/katib-ui@" ${FILE_PATH}

# Change Katib cert generator image.
sed -i -e "s@newName: docker.io/kubeflowkatib/cert-generator@newName: ${ECR_REGISTRY}/${REPO_NAME}/v1beta1/cert-generator@" ${FILE_PATH}

# Change Katib metrics collector images.
sed -i -e "s@docker.io/kubeflowkatib/file-metrics-collector@${ECR_REGISTRY}/${REPO_NAME}/v1beta1/file-metrics-collector@" ${FILE_PATH}
sed -i -e "s@docker.io/kubeflowkatib/tfevent-metrics-collector@${ECR_REGISTRY}/${REPO_NAME}/v1beta1/tfevent-metrics-collector@" ${FILE_PATH}

# Change Katib Suggestion images.
sed -i -e "s@docker.io/kubeflowkatib/suggestion-hyperopt@${ECR_REGISTRY}/${REPO_NAME}/v1beta1/suggestion-hyperopt@" ${FILE_PATH}
sed -i -e "s@docker.io/kubeflowkatib/suggestion-chocolate@${ECR_REGISTRY}/${REPO_NAME}/v1beta1/suggestion-chocolate@" ${FILE_PATH}
sed -i -e "s@docker.io/kubeflowkatib/suggestion-hyperband@${ECR_REGISTRY}/${REPO_NAME}/v1beta1/suggestion-hyperband@" ${FILE_PATH}
sed -i -e "s@docker.io/kubeflowkatib/suggestion-skopt@${ECR_REGISTRY}/${REPO_NAME}/v1beta1/suggestion-skopt@" ${FILE_PATH}
sed -i -e "s@docker.io/kubeflowkatib/suggestion-goptuna@${ECR_REGISTRY}/${REPO_NAME}/v1beta1/suggestion-goptuna@" ${FILE_PATH}
sed -i -e "s@docker.io/kubeflowkatib/suggestion-enas@${ECR_REGISTRY}/${REPO_NAME}/v1beta1/suggestion-enas@" ${FILE_PATH}
sed -i -e "s@docker.io/kubeflowkatib/suggestion-darts@${ECR_REGISTRY}/${REPO_NAME}/v1beta1/suggestion-darts@" ${FILE_PATH}

# Change Katib Early Stopping images.
sed -i -e "s@docker.io/kubeflowkatib/earlystopping-medianstop@${ECR_REGISTRY}/${REPO_NAME}/v1beta1/earlystopping-medianstop@" ${FILE_PATH}

echo "Katib images have been updated"
cat ${FILE_PATH}

# Update Trial template images in the examples.
./scripts/v1beta1/update-trial-images.sh -p "${ECR_REGISTRY}/${REPO_NAME}/v1beta1/trial-" -t ${VERSION}

echo "Creating Kubeflow namespace"
kubectl create namespace kubeflow

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
# 1 Pod for the cert-generator Job
PODNUM=$((PODNUM + 1))
until kubectl get pods -n kubeflow | grep -E 'Running|Completed' | [[ $(wc -l) -eq $PODNUM ]]; do
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
