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

set -o errexit
set -o nounset
set -o pipefail

PREFIX="katib"
CMD_PREFIX="cmd"

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..

cd ${SCRIPT_ROOT}

echo "Building core image..."
docker build -t ${PREFIX}/vizier-core -f ${CMD_PREFIX}/manager/Dockerfile .

echo "Building suggestion images..."
docker build -t ${PREFIX}/suggestion-random -f ${CMD_PREFIX}/suggestion/random/Dockerfile .
docker build -t ${PREFIX}/suggestion-grid -f ${CMD_PREFIX}/suggestion/grid/Dockerfile .
#docker build -t ${PREFIX}/suggestion-hyperband -f ${CMD_PREFIX}/suggestion/hyperband/Dockerfile .
#docker build -t ${PREFIX}/suggestion-bayesianoptimization -f ${CMD_PREFIX}/suggestion/bayesianoptimization/Dockerfile .

echo "Building earlystopping images..."
docker build -t ${PREFIX}/earlystopping-medianstopping -f ${CMD_PREFIX}/earlystopping/medianstopping/Dockerfile .

echo "Building UI image..."
docker build -t ${PREFIX}/katib-frontend -f modeldb/Dockerfile .

#echo "Building CLI image..."
#docker build -t ${PREFIX}/katib-cli -f ${CMD_PREFIX}/cli/Dockerfile .
cd - > /dev/null
