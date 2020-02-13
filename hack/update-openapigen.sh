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
VERSION="v1alpha3"
SWAGGER_CODEGEN_FILE="pkg/apis/v1alpha3/swagger.json"
SWAGGER_VERSION="0.1"

echo "Generating OpenAPI specification ..."
go run ${CODEGEN_PKG}/cmd/openapi-gen/main.go \
  --go-header-file ${PROJECT_ROOT}/hack/boilerplate.go.txt \
  --input-dirs github.com/kubeflow/katib/pkg/apis/controller/common/${VERSION},github.com/kubeflow/katib/pkg/apis/controller/experiments/${VERSION},github.com/kubeflow/katib/pkg/apis/controller/suggestions/${VERSION},github.com/kubeflow/katib/pkg/apis/controller/trials/${VERSION} \
  --output-package github.com/kubeflow/katib/pkg/apis/${VERSION} \
  $@

echo "Generating swagger file ..."
go run hack/swagger/main.go ${SWAGGER_VERSION} > ${SWAGGER_CODEGEN_FILE}
