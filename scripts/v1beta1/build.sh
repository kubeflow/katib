#!/usr/bin/env bash

# Copyright 2022 The Kubeflow Authors.
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

# This script is used to build all Katib images.
# Run ./scripts/v1beta1/build.sh <IMAGE_REGISTRY> <TAG> to execute it.

set -e

REGISTRY=$1
TAG=$2
ARCH=$3

if [[ -z "${REGISTRY}" || -z "${TAG}" || -z "${ARCH}" ]]; then
  echo "Image registry, tag and cpu-architecture must be set"
  echo "Usage: $0 <image-registry> <image-tag> <cpu-architecture>" 1>&2
  exit 1
fi

SUPPORTED_CPU_ARCHS=(amd64 arm64 ppc64le)
function check_specified_cpu_arch() {
  for SUPPORTED_ARCH in "${SUPPORTED_CPU_ARCHS[@]}"; do
    if [ "${ARCH}" = "${SUPPORTED_ARCH}" ]; then
      return 0
    fi
  done
  echo "CPU architecture '${ARCH}' is not supported"
  echo "You can use '${SUPPORTED_CPU_ARCHS[*]}'"
  echo "To get machine architecture run: uname -m"
  return 1
}
check_specified_cpu_arch

VERSION="v1beta1"
CMD_PREFIX="cmd"

echo "Building images for Katib ${VERSION}..."
echo "Image registry: ${REGISTRY}"
echo "Image tag: ${TAG}"
echo "CPU architecture: ${ARCH}"

SCRIPT_ROOT=$(dirname "$0")/../..
cd "${SCRIPT_ROOT}"

# Katib core images
echo -e "\nBuilding Katib controller image...\n"
docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/katib-controller:${TAG}" -f ${CMD_PREFIX}/katib-controller/${VERSION}/Dockerfile .

echo -e "\nBuilding Katib DB manager image...\n"
docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/katib-db-manager:${TAG}" -f ${CMD_PREFIX}/db-manager/${VERSION}/Dockerfile .

echo -e "\nBuilding Katib UI image...\n"
docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/katib-ui:${TAG}" -f ${CMD_PREFIX}/ui/${VERSION}/Dockerfile .

echo -e "\nBuilding file metrics collector image...\n"
docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/file-metrics-collector:${TAG}" -f ${CMD_PREFIX}/metricscollector/${VERSION}/file-metricscollector/Dockerfile .

echo -e "\nBuilding TF Event metrics collector image...\n"
if [ "${ARCH}" == "ppc64le" ]; then
  docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/tfevent-metrics-collector:${TAG}" -f ${CMD_PREFIX}/metricscollector/${VERSION}/tfevent-metricscollector/Dockerfile.ppc64le .
else
  docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/tfevent-metrics-collector:${TAG}" -f ${CMD_PREFIX}/metricscollector/${VERSION}/tfevent-metricscollector/Dockerfile .
fi

# Suggestion images
echo -e "\nBuilding suggestion images..."

echo -e "\nBuilding hyperopt suggestion...\n"
docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/suggestion-hyperopt:${TAG}" -f ${CMD_PREFIX}/suggestion/hyperopt/${VERSION}/Dockerfile .

echo -e "\nBuilding hyperband suggestion...\n"
docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/suggestion-hyperband:${TAG}" -f ${CMD_PREFIX}/suggestion/hyperband/${VERSION}/Dockerfile .

echo -e "\nBuilding skopt suggestion...\n"
docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/suggestion-skopt:${TAG}" -f ${CMD_PREFIX}/suggestion/skopt/${VERSION}/Dockerfile .

echo -e "\nBuilding goptuna suggestion...\n"
docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/suggestion-goptuna:${TAG}" -f ${CMD_PREFIX}/suggestion/goptuna/${VERSION}/Dockerfile .

echo -e "\nBuilding optuna suggestion...\n"
docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/suggestion-optuna:${TAG}" -f ${CMD_PREFIX}/suggestion/optuna/${VERSION}/Dockerfile .

echo -e "\nBuilding ENAS suggestion...\n"
docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/suggestion-enas:${TAG}" -f ${CMD_PREFIX}/suggestion/nas/enas/${VERSION}/Dockerfile .

echo -e "\nBuilding DARTS suggestion...\n"
docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/suggestion-darts:${TAG}" -f ${CMD_PREFIX}/suggestion/nas/darts/${VERSION}/Dockerfile .

echo -e "\nBuilding PBT suggestion...\n"
docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/suggestion-pbt:${TAG}" -f ${CMD_PREFIX}/suggestion/pbt/${VERSION}/Dockerfile .

# Early stopping images
echo -e "\nBuilding early stopping images...\n"

echo -e "\nBuilding median stopping rule...\n"
docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/earlystopping-medianstop:${TAG}" -f ${CMD_PREFIX}/earlystopping/medianstop/${VERSION}/Dockerfile .

# Training container images
echo -e "\nBuilding training container images..."

echo -e "\nBuilding dynamic learning rate training container example for PBT...\n"
docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/simple-pbt:${TAG}" -f examples/${VERSION}/trial-images/simple-pbt/Dockerfile .

echo -e "\nBuilding PyTorch CIFAR-10 CNN training container example for DARTS with CPU support...\n"
docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/darts-cnn-cifar10-cpu:${TAG}" -f examples/${VERSION}/trial-images/darts-cnn-cifar10/Dockerfile.cpu .

echo -e "\nBuilding PyTorch CIFAR-10 CNN training container example for DARTS with GPU support...\n"
docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/darts-cnn-cifar10-gpu:${TAG}" -f examples/${VERSION}/trial-images/darts-cnn-cifar10/Dockerfile.gpu .

echo -e "\nBuilding mxnet mnist training container example...\n"
docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/mxnet-mnist:${TAG}" -f examples/${VERSION}/trial-images/mxnet-mnist/Dockerfile .

echo -e "\nBuilding PyTorch mnist training container example with CPU support...\n"
docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/pytorch-mnist-cpu:${TAG}" -f examples/${VERSION}/trial-images/pytorch-mnist/Dockerfile.cpu .

echo -e "\nBuilding PyTorch mnist training container example with GPU support...\n"
docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/pytorch-mnist-gpu:${TAG}" -f examples/${VERSION}/trial-images/pytorch-mnist/Dockerfile.gpu .

echo -e "\nBuilding Tensorflow with summaries mnist training container example...\n"
docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/tf-mnist-with-summaries:${TAG}" -f examples/${VERSION}/trial-images/tf-mnist-with-summaries/Dockerfile .

echo -e "\nBuilding Keras CIFAR-10 CNN training container example for ENAS with CPU support...\n"
docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/enas-cnn-cifar10-cpu:${TAG}" -f examples/${VERSION}/trial-images/enas-cnn-cifar10/Dockerfile.cpu .

echo -e "\nBuilding Keras CIFAR-10 CNN training container example for ENAS with GPU support...\n"
docker buildx build --platform "linux/${ARCH}" -t "${REGISTRY}/enas-cnn-cifar10-gpu:${TAG}" -f examples/${VERSION}/trial-images/enas-cnn-cifar10/Dockerfile.gpu .

echo -e "\nAll Katib images with ${TAG} tag have been built successfully!\n"
