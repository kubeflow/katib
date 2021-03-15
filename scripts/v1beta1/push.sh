#!/bin/bash

# Copyright 2021 The Kubeflow Authors.
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

# This script is used to build and push all Katib images to the registry.
# It adds commit tag and release tag to the images.
# Commit tag format must be: v1beta1-<COMMIT-SHA>.
# Run ./scripts/v1beta1/push.sh <IMAGE_REGISTRY> <COMMIT_TAG> <RELEASE_TAG> to execute it.

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

# Building the images
make build REGISTRY=${REGISTRY} COMMIT_TAG=${COMMIT_TAG} RELEASE_TAG=${RELEASE_TAG}

echo "Pushing images for Katib ${VERSION}..."
echo "Image registry: ${REGISTRY}"
echo "Image commit tag: ${COMMIT_TAG}"
echo "Image release tag: ${RELEASE_TAG}"

# Katib core images
echo -e "\nPushing Katib controller image...\n"
docker push ${REGISTRY}/katib-controller:${COMMIT_TAG}
docker push ${REGISTRY}/katib-controller:${RELEASE_TAG}

echo -e "\nPushing Katib DB manager image...\n"
docker push ${REGISTRY}/katib-db-manager:${COMMIT_TAG}
docker push ${REGISTRY}/katib-db-manager:${RELEASE_TAG}

echo -e "\nPushing Katib UI image...\n"
docker push ${REGISTRY}/katib-ui:${COMMIT_TAG}
docker push ${REGISTRY}/katib-ui:${RELEASE_TAG}

echo -e "\nPushing Katib cert generator image...\n"
docker push ${REGISTRY}/cert-generator:${COMMIT_TAG}
docker push ${REGISTRY}/cert-generator:${RELEASE_TAG}

echo -e "\nPushing file metrics collector image...\n"
docker push ${REGISTRY}/file-metrics-collector:${COMMIT_TAG}
docker push ${REGISTRY}/file-metrics-collector:${RELEASE_TAG}

echo -e "\nPushing TF Event metrics collector image...\n"
docker push ${REGISTRY}/tfevent-metrics-collector:${COMMIT_TAG}
docker push ${REGISTRY}/tfevent-metrics-collector:${RELEASE_TAG}

# Suggestion images
echo -e "\nPushing suggestion images..."

echo -e "\nPushing hyperopt suggestion...\n"
docker push ${REGISTRY}/suggestion-hyperopt:${COMMIT_TAG}
docker push ${REGISTRY}/suggestion-hyperopt:${RELEASE_TAG}

echo -e "\nPushing chocolate suggestion...\n"
docker push ${REGISTRY}/suggestion-chocolate:${COMMIT_TAG}
docker push ${REGISTRY}/suggestion-chocolate:${RELEASE_TAG}

echo -e "\nPushing hyperband suggestion...\n"
docker push ${REGISTRY}/suggestion-hyperband:${COMMIT_TAG}
docker push ${REGISTRY}/suggestion-hyperband:${RELEASE_TAG}

echo -e "\nPushing skopt suggestion...\n"
docker push ${REGISTRY}/suggestion-skopt:${COMMIT_TAG}
docker push ${REGISTRY}/suggestion-skopt:${RELEASE_TAG}

echo -e "\nPushing goptuna suggestion...\n"
docker push ${REGISTRY}/suggestion-goptuna:${COMMIT_TAG}
docker push ${REGISTRY}/suggestion-goptuna:${RELEASE_TAG}

echo -e "\nPushing ENAS suggestion...\n"
docker push ${REGISTRY}/suggestion-enas:${COMMIT_TAG}
docker push ${REGISTRY}/suggestion-enas:${RELEASE_TAG}

echo -e "\nPushing DARTS suggestion...\n"
docker push ${REGISTRY}/suggestion-darts:${COMMIT_TAG}
docker push ${REGISTRY}/suggestion-darts:${RELEASE_TAG}

# Early stopping images
echo -e "\nPushing early stopping images...\n"

echo -e "\nPushing median stopping rule...\n"
docker push ${REGISTRY}/earlystopping-medianstop:${COMMIT_TAG}
docker push ${REGISTRY}/earlystopping-medianstop:${RELEASE_TAG}

# Training container images
echo -e "\nPushing training container images..."

echo -e "\nPushing mxnet mnist training container example...\n"
docker push ${REGISTRY}/mxnet-mnist:${COMMIT_TAG}
docker push ${REGISTRY}/mxnet-mnist:${RELEASE_TAG}

echo -e "\nPushing PyTorch mnist training container example...\n"
docker push ${REGISTRY}/pytorch-mnist:${COMMIT_TAG}
docker push ${REGISTRY}/pytorch-mnist:${RELEASE_TAG}

echo -e "\nPushing Keras CIFAR-10 CNN training container example for ENAS with GPU support...\n"
docker push ${REGISTRY}/enas-cnn-cifar10-gpu:${COMMIT_TAG}
docker push ${REGISTRY}/enas-cnn-cifar10-gpu:${RELEASE_TAG}

echo -e "\nPushing Keras CIFAR-10 CNN training container example for ENAS with CPU support...\n"
docker push ${REGISTRY}/enas-cnn-cifar10-cpu:${COMMIT_TAG}
docker push ${REGISTRY}/enas-cnn-cifar10-cpu:${RELEASE_TAG}

echo -e "\nPushing PyTorch CIFAR-10 CNN training container example for DARTS...\n"
docker push ${REGISTRY}/darts-cnn-cifar10:${COMMIT_TAG}
docker push ${REGISTRY}/darts-cnn-cifar10:${RELEASE_TAG}

echo -e "\nAll Katib images have been pushed successfully!\n"
