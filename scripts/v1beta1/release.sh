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

set -e

COMMIT=$(git rev-parse --short HEAD)
REGISTRY="docker.io/andreyvelichkevich"
VERSION="v1beta1"
TAG=${VERSION}-${COMMIT}

echo "Releasing images for Katib ${VERSION}..."
echo "Commit SHA: ${COMMIT}"
echo "Image registry: ${REGISTRY}"
echo -e "Image tag: ${TAG}\n"

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/../..
cd ${SCRIPT_ROOT}

# Building the images
make build REGISTRY=${REGISTRY} TAG=${TAG}

# Releasing the images
echo -e "\nAll Katib images have been successfully built\n"

# Katib core images
echo -e "\nPushing Katib controller image...\n"
docker push ${REGISTRY}/katib-controller:${TAG}

echo -e "\nPushing Katib DB manager image...\n"
docker push ${REGISTRY}/katib-db-manager:${TAG}

echo -e "\nPushing Katib UI image...\n"
docker push ${REGISTRY}/katib-ui:${TAG}

echo -e "\nPushing file metrics collector image...\n"
docker push ${REGISTRY}/file-metrics-collector:${TAG}

echo -e "\nPushing TF Event metrics collector image...\n"
docker push ${REGISTRY}/tfevent-metrics-collector:${TAG}

# Suggestion images
echo -e "\nPushing suggestion images..."

echo -e "\nPushing hyperopt suggestion...\n"
docker push ${REGISTRY}/suggestion-hyperopt:${TAG}

echo -e "\nPushing chocolate suggestion...\n"
docker push ${REGISTRY}/suggestion-chocolate:${TAG}

echo -e "\nPushing hyperband suggestion...\n"
docker push ${REGISTRY}/suggestion-hyperband:${TAG}

echo -e "\nPushing skopt suggestion...\n"
docker push ${REGISTRY}/suggestion-skopt:${TAG}

echo -e "\nPushing goptuna suggestion...\n"
docker push ${REGISTRY}/suggestion-goptuna:${TAG}

echo -e "\nPushing ENAS suggestion...\n"
docker push ${REGISTRY}/suggestion-enas:${TAG}

echo -e "\nPushing DARTS suggestion...\n"
docker push ${REGISTRY}/suggestion-darts:${TAG}

# Early stopping images
echo -e "\nPushing early stopping images...\n"

echo -e "\nPushing median stopping rule...\n"
docker push ${REGISTRY}/earlystopping-medianstop:${TAG}

# Training container images
echo -e "\nPushing training container images...\n"

echo -e "\nPushing mxnet mnist training container example...\n"
docker push ${REGISTRY}/mxnet-mnist:${TAG}

echo -e "\nPushing PyTorch mnist training container example...\n"
docker push ${REGISTRY}/pytorch-mnist:${TAG}

echo -e "\nPushing Keras CIFAR-10 CNN training container example for ENAS with GPU support...\n"
docker push ${REGISTRY}/enas-cnn-cifar10-gpu:${TAG}

echo -e "\nPushing Keras CIFAR-10 CNN training container example for ENAS with CPU support...\n"
docker push ${REGISTRY}/enas-cnn-cifar10-gpu:${TAG}

echo -e "\nPushing PyTorch CIFAR-10 CNN training container example for DARTS...\n"
docker push ${REGISTRY}/darts-cnn-cifar10:${TAG}

echo "Katib ${VERSION} for ${COMMIT} has been released successfully!"
