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

echo "Generate open-api for the APIs"

CURRENT_DIR=$(dirname "${BASH_SOURCE[0]}")
KATIB_ROOT=$(realpath "${CURRENT_DIR}/..")
KATIB_PKG="github.com/kubeflow/katib"
CODEGEN_PKG=$(go list -m -mod=readonly -f "{{.Dir}}" k8s.io/code-generator)

cd "$CURRENT_DIR/.."

# shellcheck source=/dev/null
source "${CODEGEN_PKG}/kube_codegen.sh"

VERSION_LIST=(v1beta1)

for VERSION in "${VERSION_LIST[@]}"; do
  echo "Generating OpenAPI specification for ${VERSION} ..."

  kube::codegen::gen_openapi \
    --boilerplate "${KATIB_ROOT}/hack/boilerplate/boilerplate.go.txt" \
    --output-pkg "${KATIB_PKG}/pkg/apis/${VERSION}" \
    --output-dir "${KATIB_ROOT}/pkg/apis/${VERSION}" \
    --report-filename "${KATIB_ROOT}/hack/violation_exception_${VERSION}.list" \
    --update-report \
    "${KATIB_ROOT}/pkg/apis/controller"
done
