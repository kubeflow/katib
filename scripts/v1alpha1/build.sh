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

echo "Building core image..."
docker build -t ${PREFIX}/v1alpha1/vizier-core -f ${CMD_PREFIX}/manager/v1alpha1/Dockerfile .
docker build -t ${PREFIX}/v1alpha1/studyjob-controller -f ${CMD_PREFIX}/katib-controller/v1alpha1/Dockerfile .
docker build -t ${PREFIX}/v1alpha1/metrics-collector -f ${CMD_PREFIX}/metricscollector/v1alpha1/Dockerfile .
docker build -t ${PREFIX}/v1alpha1/tfevent-metrics-collector -f ${CMD_PREFIX}/tfevent-metricscollector/v1alpha1/Dockerfile .

echo "Building REST API for core image..."
docker build -t ${PREFIX}/v1alpha1/vizier-core-rest -f ${CMD_PREFIX}/manager-rest/v1alpha1/Dockerfile .

echo "Building suggestion images..."
docker build -t ${PREFIX}/v1alpha1/suggestion-random -f ${CMD_PREFIX}/suggestion/random/v1alpha1/Dockerfile .
docker build -t ${PREFIX}/v1alpha1/suggestion-grid -f ${CMD_PREFIX}/suggestion/grid/v1alpha1/Dockerfile .
docker build -t ${PREFIX}/v1alpha1/suggestion-hyperband -f ${CMD_PREFIX}/suggestion/hyperband/v1alpha1/Dockerfile .
docker build -t ${PREFIX}/v1alpha1/suggestion-bayesianoptimization -f ${CMD_PREFIX}/suggestion/bayesianoptimization/v1alpha1/Dockerfile .
docker build -t ${PREFIX}/v1alpha1/suggestion-nasrl -f ${CMD_PREFIX}/suggestion/nasrl/v1alpha1/Dockerfile .

echo "Building earlystopping images..."
docker build -t ${PREFIX}/v1alpha1/earlystopping-medianstopping -f ${CMD_PREFIX}/earlystopping/medianstopping/v1alpha1/Dockerfile .

echo "Building UI image..."
docker build -t ${PREFIX}/v1alpha1/katib-ui -f ${CMD_PREFIX}/ui/v1alpha1/Dockerfile .

