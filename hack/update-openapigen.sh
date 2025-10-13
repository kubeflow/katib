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

cd "$CURRENT_DIR/.."

# Get the kube-openapi binary.
OPENAPI_PKG=$(go list -m -mod=readonly -f "{{.Dir}}" k8s.io/kube-openapi)
echo ">> Using ${OPENAPI_PKG}"

VERSION_LIST=(v1beta1)
SWAGGER_VERSION="0.1"

for VERSION in "${VERSION_LIST[@]}"; do
  SWAGGER_CODEGEN_FILE=${KATIB_ROOT}/pkg/apis/${VERSION}/swagger.json

  echo "Generating OpenAPI specification for ${VERSION} ..."
  
  go run "${OPENAPI_PKG}/cmd/openapi-gen" \
    --go-header-file "${KATIB_ROOT}/hack/boilerplate/boilerplate.go.txt" \
    --output-pkg "${KATIB_PKG}/pkg/apis/${VERSION}" \
    --output-dir "${KATIB_ROOT}/pkg/apis/${VERSION}" \
    --output-file "zz_generated.openapi.go" \
    --report-filename "${KATIB_ROOT}/hack/violation_exception_${VERSION}.list" \
    "${KATIB_ROOT}/pkg/apis/controller/common/${VERSION}" \
    "${KATIB_ROOT}/pkg/apis/controller/experiments/${VERSION}" \
    "${KATIB_ROOT}/pkg/apis/controller/suggestions/${VERSION}" \
    "${KATIB_ROOT}/pkg/apis/controller/trials/${VERSION}"
  
  echo "Generating OpenAPI Swagger for ${VERSION} ..."
  go run "${KATIB_ROOT}/hack/swagger/main.go" "${VERSION}-${SWAGGER_VERSION}" "${VERSION}" >"${SWAGGER_CODEGEN_FILE}"
done


# This list needs to cover all of the types used transitively from the Kubeflow Katib APIs.
# Update this list if Kubeflow Katib depends on new external APIs.
EXTRA_PACKAGES=(
  k8s.io/apimachinery/pkg/apis/meta/v1
  k8s.io/apimachinery/pkg/apis/meta/v1/unstructured
  k8s.io/apimachinery/pkg/api/resource
  k8s.io/apimachinery/pkg/runtime
  k8s.io/apimachinery/pkg/util/intstr
  k8s.io/api/core/v1
)

# Generating OpenAPI Swagger for Kubeflow Katib.
# TODO (andreyvelich): This is needed for Kubeflow SDK for now.
# TODO (andreyvelich): temporally, we store zz_generated.openapi.go in pkg/apis/kubeflowsdk dir
echo "Generate OpenAPI Swagger for Kubeflow SDK"
go run "${OPENAPI_PKG}/cmd/openapi-gen" \
  --go-header-file "${KATIB_ROOT}/hack/boilerplate/boilerplate.go.txt" \
  --output-pkg "${KATIB_PKG}/pkg/apis/kubeflowsdk" \
  --output-dir "${KATIB_ROOT}/pkg/apis/kubeflowsdk" \
  --output-file "zz_generated.openapi.go" \
  --report-filename "${KATIB_ROOT}/hack/violation_exception_v1beta1.list" \
  "${EXTRA_PACKAGES[@]}" \
  "${KATIB_ROOT}/pkg/apis/controller/common/v1beta1" \
  "${KATIB_ROOT}/pkg/apis/controller/experiments/v1beta1" \
  "${KATIB_ROOT}/pkg/apis/controller/suggestions/v1beta1" \
  "${KATIB_ROOT}/pkg/apis/controller/trials/v1beta1"

go run hack/swagger_kubeflow_sdk/main.go >api/openapi-spec/swagger.json
