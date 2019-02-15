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

# This shell script is used to run a pytorch job with custom cleanpod policies


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
APP_NAME=test-app
KUBEFLOW_VERSION=master
KF_ENV=pytorch

echo "Activating service-account"
gcloud auth activate-service-account --key-file=${GOOGLE_APPLICATION_CREDENTIALS}
echo "Configuring kubectl"
gcloud --project ${PROJECT} container clusters get-credentials ${CLUSTER_NAME} \
    --zone ${ZONE}


cd ${GO_DIR}

echo "Running smoke test"
SENDRECV_TEST_IMAGE_TAG="pytorch-dist-sendrecv-test:1.0"
go run ./test/e2e/v1alpha2/cleanpolicy_all.go --namespace=${NAMESPACE} --image=${REGISTRY}/${SENDRECV_TEST_IMAGE_TAG} --name=sendrecvjob-cleanall

echo "Running mnist test"
MNIST_TEST_IMAGE_TAG="pytorch-dist-mnist_test:1.0"
go run ./test/e2e/v1alpha2/cleanpolicy_all.go --namespace=${NAMESPACE} --image=${REGISTRY}/${MNIST_TEST_IMAGE_TAG} --name=mnistjob-cleanall

