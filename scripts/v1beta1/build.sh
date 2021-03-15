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
# It adds commit tag and release tag to the images.
# Commit tag format must be: v1beta1-<COMMIT-SHA>.
# Run ./scripts/v1beta1/build.sh <IMAGE_REGISTRY> <COMMIT_TAG> <RELEASE_TAG> to execute it.

set -e

REGISTRY=$1
COMMIT_TAG=$2
RELEASE_TAG=$3

if [[ -z "$REGISTRY" || -z "$COMMIT_TAG" || -z "$RELEASE_TAG" ]]; then
    echo "Image registry, commit tag and release tag must be set"
    echo "Usage: $0 <image-registry> <commit-tag> <release-tag>" 1>&2
    exit 1
fi

VERSION="v1beta1"
CMD_PREFIX="cmd"
MACHINE_ARCH=$(uname -m)

echo "Building images for Katib ${VERSION}..."
echo "Image registry: ${REGISTRY}"
echo "Image commit tag: ${COMMIT_TAG}"
echo "Image release tag: ${RELEASE_TAG}"

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/../..
cd ${SCRIPT_ROOT}

# Katib core images
echo -e "\nBuilding Katib controller image...\n"
docker build -t ${REGISTRY}/katib-controller:${COMMIT_TAG} -t ${REGISTRY}/katib-controller:${RELEASE_TAG} -f ${CMD_PREFIX}/katib-controller/${VERSION}/Dockerfile .

echo -e "\nBuilding Katib DB manager image...\n"
docker build -t ${REGISTRY}/katib-db-manager:${COMMIT_TAG} -t ${REGISTRY}/katib-db-manager:${RELEASE_TAG} -f ${CMD_PREFIX}/db-manager/${VERSION}/Dockerfile .

echo -e "\nBuilding Katib UI image...\n"
docker build -t ${REGISTRY}/katib-ui:${COMMIT_TAG} -t ${REGISTRY}/katib-ui:${RELEASE_TAG} -f ${CMD_PREFIX}/ui/${VERSION}/Dockerfile .

echo -e "\nBuilding Katib cert generator image...\n"
docker build -t ${REGISTRY}/cert-generator:${COMMIT_TAG} -t ${REGISTRY}/cert-generator:${RELEASE_TAG} -f ${CMD_PREFIX}/cert-generator/${VERSION}/Dockerfile .

echo -e "\nBuilding file metrics collector image...\n"
docker build -t ${REGISTRY}/file-metrics-collector:${COMMIT_TAG} -t ${REGISTRY}/file-metrics-collector:${RELEASE_TAG} -f ${CMD_PREFIX}/metricscollector/${VERSION}/file-metricscollector/Dockerfile .

echo -e "\nBuilding TF Event metrics collector image...\n"
if [ $MACHINE_ARCH == "aarch64" ]; then
    docker build -t ${REGISTRY}/tfevent-metrics-collector:${COMMIT_TAG} -t ${REGISTRY}/tfevent-metrics-collector:${RELEASE_TAG} -f ${CMD_PREFIX}/metricscollector/${VERSION}/tfevent-metricscollector/Dockerfile.aarch64 .
elif [ $MACHINE_ARCH == "ppc64le" ]; then
    docker build -t ${REGISTRY}/tfevent-metrics-collector:${COMMIT_TAG} -t ${REGISTRY}/tfevent-metrics-collector:${RELEASE_TAG} -f ${CMD_PREFIX}/metricscollector/${VERSION}/tfevent-metricscollector/Dockerfile.ppc64le .
else
    docker build -t ${REGISTRY}/tfevent-metrics-collector:${COMMIT_TAG} -t ${REGISTRY}/tfevent-metrics-collector:${RELEASE_TAG} -f ${CMD_PREFIX}/metricscollector/${VERSION}/tfevent-metricscollector/Dockerfile .
fi

# Suggestion images
echo -e "\nBuilding suggestion images..."

echo -e "\nBuilding hyperopt suggestion...\n"
docker build -t ${REGISTRY}/suggestion-hyperopt:${COMMIT_TAG} -t ${REGISTRY}/suggestion-hyperopt:${RELEASE_TAG} -f ${CMD_PREFIX}/suggestion/hyperopt/${VERSION}/Dockerfile .

echo -e "\nBuilding chocolate suggestion...\n"
docker build -t ${REGISTRY}/suggestion-chocolate:${COMMIT_TAG} -t ${REGISTRY}/suggestion-chocolate:${RELEASE_TAG} -f ${CMD_PREFIX}/suggestion/chocolate/${VERSION}/Dockerfile .

echo -e "\nBuilding hyperband suggestion...\n"
docker build -t ${REGISTRY}/suggestion-hyperband:${COMMIT_TAG} -t ${REGISTRY}/suggestion-hyperband:${RELEASE_TAG} -f ${CMD_PREFIX}/suggestion/hyperband/${VERSION}/Dockerfile .

echo -e "\nBuilding skopt suggestion...\n"
docker build -t ${REGISTRY}/suggestion-skopt:${COMMIT_TAG} -t ${REGISTRY}/suggestion-skopt:${RELEASE_TAG} -f ${CMD_PREFIX}/suggestion/skopt/${VERSION}/Dockerfile .

echo -e "\nBuilding goptuna suggestion...\n"
docker build -t ${REGISTRY}/suggestion-goptuna:${COMMIT_TAG} -t ${REGISTRY}/suggestion-goptuna:${RELEASE_TAG} -f ${CMD_PREFIX}/suggestion/goptuna/${VERSION}/Dockerfile .

echo -e "\nBuilding ENAS suggestion...\n"
if [ $MACHINE_ARCH == "aarch64" ]; then
    docker build -t ${REGISTRY}/suggestion-enas:${COMMIT_TAG} -t ${REGISTRY}/suggestion-enas:${RELEASE_TAG} -f ${CMD_PREFIX}/suggestion/nas/enas/${VERSION}/Dockerfile.aarch64 .
else
    docker build -t ${REGISTRY}/suggestion-enas:${COMMIT_TAG} -t ${REGISTRY}/suggestion-enas:${RELEASE_TAG} -f ${CMD_PREFIX}/suggestion/nas/enas/${VERSION}/Dockerfile .
fi

echo -e "\nBuilding DARTS suggestion...\n"
docker build -t ${REGISTRY}/suggestion-darts:${COMMIT_TAG} -t ${REGISTRY}/suggestion-darts:${RELEASE_TAG} -f ${CMD_PREFIX}/suggestion/nas/darts/${VERSION}/Dockerfile .

# Early stopping images
echo -e "\nBuilding early stopping images...\n"

echo -e "\nBuilding median stopping rule...\n"
docker build -t ${REGISTRY}/earlystopping-medianstop:${COMMIT_TAG} -t ${REGISTRY}/earlystopping-medianstop:${RELEASE_TAG} -f ${CMD_PREFIX}/earlystopping/medianstop/${VERSION}/Dockerfile .

# Training container images
echo -e "\nBuilding training container images..."

echo -e "\nBuilding mxnet mnist training container example...\n"
docker build -t ${REGISTRY}/mxnet-mnist:${COMMIT_TAG} -t ${REGISTRY}/mxnet-mnist:${RELEASE_TAG} -f examples/${VERSION}/mxnet-mnist/Dockerfile .

echo -e "\nBuilding PyTorch mnist training container example...\n"
docker build -t ${REGISTRY}/pytorch-mnist:${COMMIT_TAG} -t ${REGISTRY}/pytorch-mnist:${RELEASE_TAG} -f examples/${VERSION}/pytorch-mnist/Dockerfile .

echo -e "\nBuilding Keras CIFAR-10 CNN training container example for ENAS with GPU support...\n"
docker build -t ${REGISTRY}/enas-cnn-cifar10-gpu:${COMMIT_TAG} -t ${REGISTRY}/enas-cnn-cifar10-gpu:${RELEASE_TAG} -f examples/${VERSION}/nas/enas-cnn-cifar10/Dockerfile.gpu .

echo -e "\nBuilding Keras CIFAR-10 CNN training container example for ENAS with CPU support...\n"
docker build -t ${REGISTRY}/enas-cnn-cifar10-cpu:${COMMIT_TAG} -t ${REGISTRY}/enas-cnn-cifar10-cpu:${RELEASE_TAG} -f examples/${VERSION}/nas/enas-cnn-cifar10/Dockerfile.cpu .

echo -e "\nBuilding PyTorch CIFAR-10 CNN training container example for DARTS...\n"
docker build -t ${REGISTRY}/darts-cnn-cifar10:${COMMIT_TAG} -t ${REGISTRY}/darts-cnn-cifar10:${RELEASE_TAG} -f examples/${VERSION}/nas/darts-cnn-cifar10/Dockerfile .

echo -e "\nAll Katib images have been built successfully!\n"
