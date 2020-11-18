#!/usr/bin/env bash

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

SWAGGER_JAR_URL="http://search.maven.org/maven2/io/swagger/swagger-codegen-cli/2.4.6/swagger-codegen-cli-2.4.6.jar"
SWAGGER_CODEGEN_JAR="hack/gen-python-sdk/swagger-codegen-cli.jar"
SWAGGER_CODEGEN_CONF="hack/gen-python-sdk/swagger_config.json"
SWAGGER_CODEGEN_FILE="pkg/apis/KATIB_VERSION/swagger.json"
TMP_CODEGEN_PATH="sdk/tmp/KATIB_VERSION"
SDK_OUTPUT_PATH="sdk/python"
POST_GEN_PYTHON_HANDLER="hack/gen-python-sdk/post_gen.py"
KATIB_VERSIONS=(v1beta1)

echo "Downloading the swagger-codegen JAR package ..."
wget -O ${SWAGGER_CODEGEN_JAR} ${SWAGGER_JAR_URL}

for VERSION in ${KATIB_VERSIONS[@]}; do
    echo "Generating Python SDK for Kubeflow Katib ${VERSION} ..."
    SWAGGER_FILE=${SWAGGER_CODEGEN_FILE/KATIB_VERSION/$VERSION}
    TMP_PATH=${TMP_CODEGEN_PATH/KATIB_VERSION/$VERSION}
    java -jar ${SWAGGER_CODEGEN_JAR} generate -i ${SWAGGER_FILE} -l python -o ${TMP_PATH} -c ${SWAGGER_CODEGEN_CONF}

    python ${POST_GEN_PYTHON_HANDLER} ${TMP_PATH} ${SDK_OUTPUT_PATH}/${VERSION}
done

rm ${SWAGGER_CODEGEN_JAR}
