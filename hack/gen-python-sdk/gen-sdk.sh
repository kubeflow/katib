#!/usr/bin/env bash

# Copyright 2021 The Kubeflow Authors.
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

SWAGGER_JAR_URL="https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/4.3.1/openapi-generator-cli-4.3.1.jar"
SWAGGER_CODEGEN_JAR="hack/gen-python-sdk/openapi-generator-cli.jar"

SWAGGER_CODEGEN_CONF="hack/gen-python-sdk/swagger_config.json"
SWAGGER_CODEGEN_FILE="pkg/apis/KATIB_VERSION/swagger.json"

TMP_CODEGEN_PATH="sdk/tmp/KATIB_VERSION"
SDK_OUTPUT_PATH="sdk/python"
POST_GEN_PYTHON_HANDLER="hack/gen-python-sdk/post_gen.py"
KATIB_VERSIONS=(v1beta1)

# Download JAR package if file doesn't exist.
if ! test -f ${SWAGGER_CODEGEN_JAR}; then
    echo "Downloading the openapi generator JAR package ..."
    wget -O ${SWAGGER_CODEGEN_JAR} ${SWAGGER_JAR_URL}
fi

for VERSION in ${KATIB_VERSIONS[@]}; do
    echo "Generating Python SDK for Kubeflow Katib ${VERSION} ..."
    SWAGGER_FILE=${SWAGGER_CODEGEN_FILE/KATIB_VERSION/$VERSION}
    TMP_PATH=${TMP_CODEGEN_PATH/KATIB_VERSION/$VERSION}
    java -jar ${SWAGGER_CODEGEN_JAR} generate -i ${SWAGGER_FILE} -g python -o ${TMP_PATH} -c ${SWAGGER_CODEGEN_CONF}

    # Run post gen script.
    python ${POST_GEN_PYTHON_HANDLER} ${TMP_PATH} ${SDK_OUTPUT_PATH}/${VERSION}
done
