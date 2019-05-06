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

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/../..

cd ${SCRIPT_ROOT}

docker build -t ${PREFIX}/katib-controller -f ${CMD_PREFIX}/katib-controller/Dockerfile .

echo "Building suggestion images..."
docker build -t ${PREFIX}/suggestion-hyperband -f ${CMD_PREFIX}/suggestion/v1alpha2/hyperband/Dockerfile .
docker build -t ${PREFIX}/suggestion -f ${CMD_PREFIX}/suggestion/v1alpha2/katib-suggestion/Dockerfile .
docker build -t ${PREFIX}/suggestion-nasrl -f ${CMD_PREFIX}/suggestion/v1alpha2/nasrl/Dockerfile .
