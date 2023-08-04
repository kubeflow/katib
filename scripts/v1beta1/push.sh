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

# This script is used to push all Katib images.
# Run ./scripts/v1beta1/push.sh <IMAGE_REGISTRY> <TAG>

set -e

REGISTRY=$1
TAG=$2

if [[ -z "$REGISTRY" || -z "$TAG" ]]; then
  echo "Image registry and tag must be set"
  echo "Usage: $0 <image-registry> <image-tag>" 1>&2
  exit 1
fi

VERSION="v1beta1"

echo "Pushing images for Katib ${VERSION}..."
echo "Image registry: ${REGISTRY}"
echo "Image tag: ${TAG}"

# Katib core images
echo -e "\nPushing Katib controller image...\n"
docker push "${REGISTRY}/katib-controller:${TAG}"

echo -e "\nPushing Katib DB manager image...\n"
docker push "${REGISTRY}/katib-db-manager:${TAG}"

echo -e "\nPushing Katib UI image...\n"
docker push "${REGISTRY}/katib-ui:${TAG}"

echo -e "\nPushing file metrics collector image...\n"
docker push "${REGISTRY}/file-metrics-collector:${TAG}"

echo -e "\nPushing TF Event metrics collector image...\n"
docker push "${REGISTRY}/tfevent-metrics-collector:${TAG}"

# Suggestion images
echo -e "\nPushing suggestion images..."

echo -e "\nPushing hyperopt suggestion...\n"
docker push "${REGISTRY}/suggestion-hyperopt:${TAG}"

echo -e "\nPushing hyperband suggestion...\n"
docker push "${REGISTRY}/suggestion-hyperband:${TAG}"

echo -e "\nPushing skopt suggestion...\n"
docker push "${REGISTRY}/suggestion-skopt:${TAG}"

echo -e "\nPushing goptuna suggestion...\n"
docker push "${REGISTRY}/suggestion-goptuna:${TAG}"

echo -e "\nPushing optuna suggestion...\n"
docker push "${REGISTRY}/suggestion-optuna:${TAG}"

echo -e "\nPushing ENAS suggestion...\n"
docker push "${REGISTRY}/suggestion-enas:${TAG}"

echo -e "\nPushing DARTS suggestion...\n"
docker push "${REGISTRY}/suggestion-darts:${TAG}"

echo -e "\nPushing PBT suggestion...\n"
docker push "${REGISTRY}/suggestion-pbt:${TAG}"

# Early stopping images
echo -e "\nPushing early stopping images...\n"

echo -e "\nPushing median stopping rule...\n"
docker push "${REGISTRY}/earlystopping-medianstop:${TAG}"

# Training container images
echo -e "\nPushing training container images..."

echo -e "\nPushing mxnet mnist training container example...\n"
docker push "${REGISTRY}/mxnet-mnist:${TAG}"

echo -e "\nPushing Tensorflow with summaries mnist training container example...\n"
docker push "${REGISTRY}/tf-mnist-with-summaries:${TAG}"

echo -e "\nPushing PyTorch mnist training container example with CPU support...\n"
docker push "${REGISTRY}/pytorch-mnist-cpu:${TAG}"

echo -e "\nPushing PyTorch mnist training container example with GPU support...\n"
docker push "${REGISTRY}/pytorch-mnist-gpu:${TAG}"

echo -e "\nPushing Keras CIFAR-10 CNN training container example for ENAS with GPU support...\n"
docker push "${REGISTRY}/enas-cnn-cifar10-gpu:${TAG}"

echo -e "\nPushing Keras CIFAR-10 CNN training container example for ENAS with CPU support...\n"
docker push "${REGISTRY}/enas-cnn-cifar10-cpu:${TAG}"

echo -e "\nPushing PyTorch CIFAR-10 CNN training container example for DARTS with CPU support...\n"
docker push "${REGISTRY}/darts-cnn-cifar10-cpu:${TAG}"

echo -e "\nPushing PyTorch CIFAR-10 CNN training container example for DARTS with GPU support...\n"
docker push "${REGISTRY}/darts-cnn-cifar10-gpu:${TAG}"

echo -e "\nPushing dynamic learning rate training container example for PBT...\n"
docker push "${REGISTRY}/simple-pbt:${TAG}"

echo -e "\nAll Katib images with ${TAG} tag have been pushed successfully!\n"
