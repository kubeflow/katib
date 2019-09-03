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

# echo "Building core image..."
docker build -t ${PREFIX}/v1alpha2/katib-controller -f ${CMD_PREFIX}/katib-controller/v1alpha2/Dockerfile .
# docker build -t ${PREFIX}/v1alpha2/katib-manager -f ${CMD_PREFIX}/manager/v1alpha2/Dockerfile .
# docker build -t ${PREFIX}/v1alpha2/katib-manager-rest -f ${CMD_PREFIX}/manager-rest/v1alpha2/Dockerfile .
# docker build -t ${PREFIX}/v1alpha2/metrics-collector -f ${CMD_PREFIX}/metricscollector/v1alpha2/Dockerfile .

# echo "Building UI image..."
# docker build -t ${PREFIX}/v1alpha2/katib-ui -f ${CMD_PREFIX}/ui/v1alpha2/Dockerfile .

# echo "Building sidecar metrics collector image..."
# docker build -t ${PREFIX}/v1alpha2/sidecar-metrics-collector -f ${CMD_PREFIX}/sidecar-metricscollector/v1alpha2/Dockerfile .

# echo "Building TF Event metrics collector image..."
# docker build -t ${PREFIX}/v1alpha2/tfevent-metrics-collector -f ${CMD_PREFIX}/tfevent-metricscollector/v1alpha2/Dockerfile .

echo "Building suggestion images..."
# docker build -t ${PREFIX}/v1alpha2/suggestion-random -f ${CMD_PREFIX}/suggestion/random/v1alpha2/Dockerfile .
# docker build -t ${PREFIX}/v1alpha2/suggestion-bayesianoptimization -f ${CMD_PREFIX}/suggestion/bayesianoptimization/v1alpha2/Dockerfile .
# docker build -t ${PREFIX}/v1alpha2/suggestion-grid -f ${CMD_PREFIX}/suggestion/grid/v1alpha2/Dockerfile .
# docker build -t ${PREFIX}/v1alpha2/suggestion-hyperband -f ${CMD_PREFIX}/suggestion/hyperband/v1alpha2/Dockerfile .
# docker build -t ${PREFIX}/v1alpha2/suggestion-nasrl -f ${CMD_PREFIX}/suggestion/nasrl/v1alpha2/Dockerfile .
