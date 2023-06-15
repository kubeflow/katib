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

set -o errexit
set -o nounset
set -o pipefail

echo "Generate deepcopy, clientset, listers, informers for the APIs"

if [[ -z "${GOPATH:-}" ]]; then
    GOPATH=$(go env GOPATH)
    export GOPATH
fi

# Grab code-generator version from go.mod
CODEGEN_VERSION=$(cd ../.. && grep 'k8s.io/code-generator' go.mod | awk '{print $2}')
CODEGEN_PKG="$GOPATH/pkg/mod/k8s.io/code-generator@${CODEGEN_VERSION}"

if [[ ! -d "${CODEGEN_PKG}" ]]; then
    echo "${CODEGEN_PKG} is missing. Please run 'go mod download'."
    exit 0
fi

echo ">> Using ${CODEGEN_PKG} for the code generator"

# Ensure we can execute.
chmod +x "${CODEGEN_PKG}/generate-groups.sh"

PROJECT_ROOT=${GOPATH}/src/github.com/kubeflow/katib

modules=(experiments suggestions trials common)
versions=(v1beta1)
versionStr=$(printf ",%s" "${versions[@]}")
GROUP_VERSIONS=$(printf "%s:${versionStr:1} " "${modules[@]}")

echo "Generating clients for ${GROUP_VERSIONS} ..."
"${CODEGEN_PKG}/generate-groups.sh" \
    all \
    github.com/kubeflow/katib/pkg/client/controller \
    github.com/kubeflow/katib/pkg/apis/controller \
    "${GROUP_VERSIONS}" \
    --go-header-file "${PROJECT_ROOT}/hack/boilerplate/boilerplate.go.txt"

echo "Generating deepcopy for config.kubeflow.org ..."
"${PROJECT_ROOT}/bin/controller-gen" \
    object:headerFile="${PROJECT_ROOT}/hack/boilerplate/boilerplate.go.txt" \
    paths="${PROJECT_ROOT}/pkg/apis/config/..."
