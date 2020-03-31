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
TAG="latest"
PREFIX="katib/v1alpha3"
CMD_PREFIX="cmd"
MACHINE_ARCH=`uname -m`

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/../..

cd ${SCRIPT_ROOT}

usage() { echo "Usage: $0 [-t <tag>] [-r <registry>] [-p <prefix>]" 1>&2; exit 1; }

while getopts ":t::r::p:" opt; do
    case $opt in
        t)
            TAG=${OPTARG}
            ;;
        r)
            REGISTRY=${OPTARG}
            ;;
        p)
            PREFIX=${OPTARG}
            ;;
        *)
            usage
            ;;
    esac
done
echo "Registry: ${REGISTRY}, tag: ${TAG}, prefix: ${PREFIX}"

echo "Building core image..."
docker build -t ${REGISTRY}/${PREFIX}/katib-controller:${TAG} -f ${CMD_PREFIX}/katib-controller/v1alpha3/Dockerfile .
docker build -t ${REGISTRY}/${PREFIX}/katib-db-manager:${TAG} -f ${CMD_PREFIX}/db-manager/v1alpha3/Dockerfile .

echo "Building UI image..."
docker build -t ${REGISTRY}/${PREFIX}/katib-ui:${TAG} -f ${CMD_PREFIX}/ui/v1alpha3/Dockerfile .

echo "Building file metrics collector image..."
docker build -t ${REGISTRY}/${PREFIX}/file-metrics-collector:${TAG} -f ${CMD_PREFIX}/metricscollector/v1alpha3/file-metricscollector/Dockerfile .

echo "Building TF Event metrics collector image..."
if [ $MACHINE_ARCH == "aarch64" ]; then
        docker build -t ${REGISTRY}/${PREFIX}/tfevent-metrics-collector:${TAG} -f ${CMD_PREFIX}/metricscollector/v1alpha3/tfevent-metricscollector/Dockerfile.aarch64 .
elif [ $MACHINE_ARCH == "ppc64le" ]; then
	docker build -t ${REGISTRY}/${PREFIX}/tfevent-metrics-collector:${TAG} -f ${CMD_PREFIX}/metricscollector/v1alpha3/tfevent-metricscollector/Dockerfile.ppc64le .
else
        docker build -t ${REGISTRY}/${PREFIX}/tfevent-metrics-collector:${TAG} -f ${CMD_PREFIX}/metricscollector/v1alpha3/tfevent-metricscollector/Dockerfile .
fi

echo "Building suggestion images..."
docker build -t ${REGISTRY}/${PREFIX}/suggestion-hyperopt:${TAG} -f ${CMD_PREFIX}/suggestion/hyperopt/v1alpha3/Dockerfile .
docker build -t ${REGISTRY}/${PREFIX}/suggestion-skopt:${TAG} -f ${CMD_PREFIX}/suggestion/skopt/v1alpha3/Dockerfile .
docker build -t ${REGISTRY}/${PREFIX}/suggestion-chocolate:${TAG} -f ${CMD_PREFIX}/suggestion/chocolate/v1alpha3/Dockerfile .
if [ $MACHINE_ARCH == "aarch64" ]; then
	docker build -t ${REGISTRY}/${PREFIX}/suggestion-nasrl:${TAG} -f ${CMD_PREFIX}/suggestion/nasrl/v1alpha3/Dockerfile.aarch64 .
else
	docker build -t ${REGISTRY}/${PREFIX}/suggestion-nasrl:${TAG} -f ${CMD_PREFIX}/suggestion/nasrl/v1alpha3/Dockerfile .
fi
docker build -t ${REGISTRY}/${PREFIX}/suggestion-hyperband:${TAG} -f ${CMD_PREFIX}/suggestion/hyperband/v1alpha3/Dockerfile .
docker build -t ${REGISTRY}/${PREFIX}/suggestion-goptuna:${TAG} -f ${CMD_PREFIX}/suggestion/goptuna/v1alpha3/Dockerfile .
