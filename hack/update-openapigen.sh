#!/bin/bash

# Copyright 2019 The Kubeflow Authors.
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

set -o errexit
set -o nounset
set -o pipefail

if [[ -z "${GOPATH:-}" ]]; then
  export GOPATH=$(go env GOPATH)
fi

PROJECT_ROOT=$(pwd)
CODEGEN_PKG=$(pwd)/vendor/k8s.io/code-generator

KATIB_VERSION_V1ALPHA3="v1alpha3"
SWAGGER_CODEGEN_FILE_V1ALPHA3="pkg/apis/v1alpha3/swagger.json"

KATIB_VERSION_V1BETA1="v1beta1"
SWAGGER_CODEGEN_FILE_V1BETA1="pkg/apis/v1beta1/swagger.json"

SWAGGER_VERSION="0.1"

# OpenAPI for v1alpha3 Katib version
echo "Generating OpenAPI specification for v1alpha3..."
go run ${CODEGEN_PKG}/cmd/openapi-gen/main.go \
  --go-header-file ${PROJECT_ROOT}/hack/boilerplate.go.txt \
  --input-dirs github.com/kubeflow/katib/pkg/apis/controller/common/${KATIB_VERSION_V1ALPHA3},github.com/kubeflow/katib/pkg/apis/controller/experiments/${KATIB_VERSION_V1ALPHA3},github.com/kubeflow/katib/pkg/apis/controller/suggestions/${KATIB_VERSION_V1ALPHA3},github.com/kubeflow/katib/pkg/apis/controller/trials/${KATIB_VERSION_V1ALPHA3} \
  --output-package github.com/kubeflow/katib/pkg/apis/${KATIB_VERSION_V1ALPHA3} \
  $@

# Swagger file for v1alpha3 Katib version
echo "Generating Swagger file for v1alpha3..."
go run hack/swagger/main.go ${SWAGGER_VERSION} ${KATIB_VERSION_V1ALPHA3} >${SWAGGER_CODEGEN_FILE_V1ALPHA3}

# OpenAPI for v1beta1 Katib version
echo "Generating OpenAPI specification for v1beta1..."
go run ${CODEGEN_PKG}/cmd/openapi-gen/main.go \
  --go-header-file ${PROJECT_ROOT}/hack/boilerplate.go.txt \
  --input-dirs github.com/kubeflow/katib/pkg/apis/controller/common/${KATIB_VERSION_V1BETA1},github.com/kubeflow/katib/pkg/apis/controller/experiments/${KATIB_VERSION_V1BETA1},github.com/kubeflow/katib/pkg/apis/controller/suggestions/${KATIB_VERSION_V1BETA1},github.com/kubeflow/katib/pkg/apis/controller/trials/${KATIB_VERSION_V1BETA1} \
  --output-package github.com/kubeflow/katib/pkg/apis/${KATIB_VERSION_V1BETA1} \
  $@

# Swagger file for v1beta1 Katib version
echo "Generating Swagger file for v1beta1..."
go run hack/swagger/main.go ${SWAGGER_VERSION} ${KATIB_VERSION_V1BETA1} >${SWAGGER_CODEGEN_FILE_V1BETA1}
