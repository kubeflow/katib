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

CURRENT_DIR=$(dirname "${BASH_SOURCE[0]}")
KATIB_ROOT=$(realpath "${CURRENT_DIR}/..")
KATIB_PKG="github.com/kubeflow/katib"
CODEGEN_PKG=$(go list -m -mod=readonly -f "{{.Dir}}" k8s.io/code-generator)

cd "$CURRENT_DIR/.."

echo "${CODEGEN_PKG}"
# shellcheck source=/dev/null
source "${CODEGEN_PKG}/kube_codegen.sh"

echo "Generating conversion and defaults functions for config.kubeflow.org ..."
kube::codegen::gen_helpers \
    --boilerplate "${KATIB_ROOT}/hack/boilerplate/boilerplate.go.txt" \
    "${KATIB_ROOT}/pkg/apis/config"

echo "Generating clients for config.kubeflow.org ..."
kube::codegen::gen_client \
    --boilerplate "${KATIB_ROOT}/hack/boilerplate/boilerplate.go.txt" \
    --output-dir "${KATIB_ROOT}/pkg/client/controller" \
    --output-pkg "${KATIB_PKG}/pkg/client/controller" \
    --with-watch \
    --with-applyconfig \
    "${KATIB_ROOT}/pkg/apis/controller"
