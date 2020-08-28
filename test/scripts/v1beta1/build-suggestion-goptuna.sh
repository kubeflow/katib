#!/bin/bash

# Copyright 2018 The Kubeflow Authors.
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

# This shell script is used to build an image from our argo workflow

set -o errexit
set -o nounset
set -o pipefail

export PATH=${GOPATH}/bin:/usr/local/go/bin:${PATH}
REGISTRY="${GCP_REGISTRY}"
PROJECT="${GCP_PROJECT}"
GO_DIR=${GOPATH}/src/github.com/${REPO_OWNER}/${REPO_NAME}-suggestion-goptuna
VERSION=$(git describe --tags --always --dirty)

# Activate gcloud service account
source utils.sh
_activate_service_account

echo "Copy source to GOPATH"
mkdir -p ${GO_DIR}
cp -r cmd ${GO_DIR}/cmd
cp -r pkg ${GO_DIR}/pkg
cp -r vendor ${GO_DIR}/vendor

cd ${GO_DIR}

cp cmd/suggestion/goptuna/v1beta1/Dockerfile .
gcloud builds submit . --tag=${REGISTRY}/${REPO_NAME}/v1beta1/suggestion-goptuna:${VERSION} --project=${PROJECT}
gcloud container images add-tag --quiet ${REGISTRY}/${REPO_NAME}/v1beta1/suggestion-goptuna:${VERSION} ${REGISTRY}/${REPO_NAME}/v1beta1/suggestion-goptuna:latest --verbosity=info
