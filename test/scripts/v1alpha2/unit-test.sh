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

GO_DIR=${GOPATH}/src/github.com/${REPO_OWNER}/${REPO_NAME}

echo "Copy source to GOPATH"
mkdir -p ${GO_DIR}
cp -r cmd ${GO_DIR}/cmd
cp -r pkg ${GO_DIR}/pkg
cp -r vendor ${GO_DIR}/vendor

echo "Run unit test cases"
cd ${GO_DIR}
os=$(go env GOOS)
arch=$(go env GOARCH)
# Install kubebuilder to setup kubernetes apiserver and etcd in the unit test.
curl -sL https://go.kubebuilder.io/dl/2.0.0/${os}/${arch} | tar -xz -C /tmp/kubebuilder
export PATH=$PATH:/tmp/kubebuilder/bin
go test ./...
cd - > /dev/null
