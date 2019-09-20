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

REGISTRY="gcr.io/kubeflow-images-public"
PREFIX="katib"
CMD_PREFIX="cmd"

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/../..

cd ${SCRIPT_ROOT}

echo "Building suggestion images..."
docker build -t ${REGISTRY}/${PREFIX}/v1alpha3/suggestion-hyperopt -f ${CMD_PREFIX}/suggestion/hyperopt/v1alpha3/Dockerfile .

echo "Building core image..."
docker build -t ${REGISTRY}/${PREFIX}/v1alpha3/katib-controller -f ${CMD_PREFIX}/katib-controller/v1alpha3/Dockerfile .
docker build -t ${REGISTRY}/${PREFIX}/v1alpha3/katib-manager -f ${CMD_PREFIX}/manager/v1alpha3/Dockerfile .
docker build -t ${REGISTRY}/${PREFIX}/v1alpha3/katib-manager-rest -f ${CMD_PREFIX}/manager-rest/v1alpha3/Dockerfile .

echo "Building UI image..."
docker build -t ${REGISTRY}/${PREFIX}/v1alpha3/katib-ui -f ${CMD_PREFIX}/ui/v1alpha3/Dockerfile .

echo "Building file metrics collector image..."
docker build -t ${REGISTRY}/${PREFIX}/v1alpha3/file-metrics-collector -f ${CMD_PREFIX}/metricscollector/v1alpha3/file-metricscollector/Dockerfile .

echo "Building TF Event metrics collector image..."
docker build -t ${REGISTRY}/${PREFIX}/v1alpha3/tfevent-metrics-collector -f ${CMD_PREFIX}/metricscollector/v1alpha3/tfevent-metricscollector/Dockerfile .
