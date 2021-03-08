#!/bin/bash

# Copyright 2020 The Kubeflow Authors.
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

# This script is used to build all Katib images.
# It adds "<TAG>" and "latest" tag to them.
# Run ./scripts/v1beta1/build.sh -r <image-registry> -t <image-tag> to execute it.

set -e

usage() {
    echo "Usage: $0 [-r <REGISTRY>] [-t <TAG>]" 1>&2
    exit 1
}

while getopts ":t::r::p:" opt; do
    case $opt in
    r)
        REGISTRY=${OPTARG}
        ;;
    t)
        TAG=${OPTARG}
        ;;
    *)
        usage
        ;;
    esac
done

if [[ -z "$REGISTRY" || -z "$TAG" ]]; then
    echo "Image registry and tag must be set"
    echo "Usage: $0 [-r <REGISTRY>] [-t <TAG>]" 1>&2
    exit 1
fi

VERSION="v1beta1"
CMD_PREFIX="cmd"
MACHINE_ARCH=$(uname -m)

echo "Building images for Katib ${VERSION}..."
echo "Image registry: ${REGISTRY}"
echo "Image tag: ${TAG}"

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/../..
cd ${SCRIPT_ROOT}

# Katib core images
echo -e "\nBuilding Katib controller image...\n"
docker build -t ${REGISTRY}/katib-controller:${TAG} -t ${REGISTRY}/katib-controller:latest -f ${CMD_PREFIX}/katib-controller/${VERSION}/Dockerfile .

echo -e "\nBuilding Katib DB manager image...\n"
docker build -t ${REGISTRY}/katib-db-manager:${TAG} -t ${REGISTRY}/katib-db-manager:latest -f ${CMD_PREFIX}/db-manager/${VERSION}/Dockerfile .

echo -e "\nBuilding Katib UI image...\n"
docker build -t ${REGISTRY}/katib-ui:${TAG} -t ${REGISTRY}/katib-ui:latest -f ${CMD_PREFIX}/ui/${VERSION}/Dockerfile .

echo -e "\nBuilding file metrics collector image...\n"
docker build -t ${REGISTRY}/file-metrics-collector:${TAG} -t ${REGISTRY}/file-metrics-collector:latest -f ${CMD_PREFIX}/metricscollector/${VERSION}/file-metricscollector/Dockerfile .

echo -e "\nBuilding TF Event metrics collector image...\n"
if [ $MACHINE_ARCH == "aarch64" ]; then
    docker build -t ${REGISTRY}/tfevent-metrics-collector:${TAG} -t ${REGISTRY}/tfevent-metrics-collector:latest -f ${CMD_PREFIX}/metricscollector/${VERSION}/tfevent-metricscollector/Dockerfile.aarch64 .
elif [ $MACHINE_ARCH == "ppc64le" ]; then
    docker build -t ${REGISTRY}/tfevent-metrics-collector:${TAG} -t ${REGISTRY}/tfevent-metrics-collector:latest -f ${CMD_PREFIX}/metricscollector/${VERSION}/tfevent-metricscollector/Dockerfile.ppc64le .
else
    docker build -t ${REGISTRY}/tfevent-metrics-collector:${TAG} -t ${REGISTRY}/tfevent-metrics-collector:latest -f ${CMD_PREFIX}/metricscollector/${VERSION}/tfevent-metricscollector/Dockerfile .
fi

echo -e "\nBuilding Katib cert generator image...\n"
docker build -t ${REGISTRY}/cert-generator:${TAG} -t ${REGISTRY}/cert-generator:latest -f ${CMD_PREFIX}/cert-generator/${VERSION}/Dockerfile .

# Suggestion images
echo -e "\nBuilding suggestion images..."

echo -e "\nBuilding hyperopt suggestion...\n"
docker build -t ${REGISTRY}/suggestion-hyperopt:${TAG} -t ${REGISTRY}/suggestion-hyperopt:latest -f ${CMD_PREFIX}/suggestion/hyperopt/${VERSION}/Dockerfile .

echo -e "\nBuilding chocolate suggestion...\n"
docker build -t ${REGISTRY}/suggestion-chocolate:${TAG} -t ${REGISTRY}/suggestion-chocolate:latest -f ${CMD_PREFIX}/suggestion/chocolate/${VERSION}/Dockerfile .

echo -e "\nBuilding hyperband suggestion...\n"
docker build -t ${REGISTRY}/suggestion-hyperband:${TAG} -t ${REGISTRY}/suggestion-hyperband:latest -f ${CMD_PREFIX}/suggestion/hyperband/${VERSION}/Dockerfile .

echo -e "\nBuilding skopt suggestion...\n"
docker build -t ${REGISTRY}/suggestion-skopt:${TAG} -t ${REGISTRY}/suggestion-skopt:latest -f ${CMD_PREFIX}/suggestion/skopt/${VERSION}/Dockerfile .

echo -e "\nBuilding goptuna suggestion...\n"
docker build -t ${REGISTRY}/suggestion-goptuna:${TAG} -t ${REGISTRY}/suggestion-goptuna:latest -f ${CMD_PREFIX}/suggestion/goptuna/${VERSION}/Dockerfile .

echo -e "\nBuilding ENAS suggestion...\n"
if [ $MACHINE_ARCH == "aarch64" ]; then
    docker build -t ${REGISTRY}/suggestion-enas:${TAG} -t ${REGISTRY}/suggestion-enas:latest -f ${CMD_PREFIX}/suggestion/nas/enas/${VERSION}/Dockerfile.aarch64 .
else
    docker build -t ${REGISTRY}/suggestion-enas:${TAG} -t ${REGISTRY}/suggestion-enas:latest -f ${CMD_PREFIX}/suggestion/nas/enas/${VERSION}/Dockerfile .
fi

echo -e "\nBuilding DARTS suggestion...\n"
docker build -t ${REGISTRY}/suggestion-darts:${TAG} -t ${REGISTRY}/suggestion-darts:latest -f ${CMD_PREFIX}/suggestion/nas/darts/${VERSION}/Dockerfile .

# Early stopping images
echo -e "\nBuilding early stopping images...\n"

echo -e "\nBuilding median stopping rule...\n"
docker build -t ${REGISTRY}/earlystopping-medianstop:${TAG} -t ${REGISTRY}/earlystopping-medianstop:latest -f ${CMD_PREFIX}/earlystopping/medianstop/${VERSION}/Dockerfile .

# Training container images
echo -e "\nBuilding training container images..."

echo -e "\nBuilding mxnet mnist training container example...\n"
docker build -t ${REGISTRY}/mxnet-mnist:${TAG} -t ${REGISTRY}/mxnet-mnist:latest -f examples/${VERSION}/mxnet-mnist/Dockerfile .

echo -e "\nBuilding PyTorch mnist training container example...\n"
docker build -t ${REGISTRY}/pytorch-mnist:${TAG} -t ${REGISTRY}/pytorch-mnist:latest -f examples/${VERSION}/pytorch-mnist/Dockerfile .

echo -e "\nBuilding Keras CIFAR-10 CNN training container example for ENAS with GPU support...\n"
docker build -t ${REGISTRY}/enas-cnn-cifar10-gpu:${TAG} -t ${REGISTRY}/enas-cnn-cifar10-gpu:latest -f examples/${VERSION}/nas/enas-cnn-cifar10/Dockerfile.gpu .

echo -e "\nBuilding Keras CIFAR-10 CNN training container example for ENAS with CPU support...\n"
docker build -t ${REGISTRY}/enas-cnn-cifar10-cpu:${TAG} -t ${REGISTRY}/enas-cnn-cifar10-cpu:latest -f examples/${VERSION}/nas/enas-cnn-cifar10/Dockerfile.cpu .

echo -e "\nBuilding PyTorch CIFAR-10 CNN training container example for DARTS...\n"
docker build -t ${REGISTRY}/darts-cnn-cifar10:${TAG} -t ${REGISTRY}/darts-cnn-cifar10:latest -f examples/${VERSION}/nas/darts-cnn-cifar10/Dockerfile .
