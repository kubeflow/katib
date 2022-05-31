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

set -o errexit
set -o pipefail
cd "$(dirname "$0")"

REGISTRY=docker.io/kubeflowkatib
TAG="e2e-test"
VERSION="v1beta1"
CMD_PREFIX="cmd"
# shellcheck disable=SC2206
TRIAL_IMAGES=(${1//,/ })

_build_containers() {
  CONTAINER_NAME=$1
  DOCKERFILE=$2
  echo -e "\nBuilding $CONTAINER_NAME image...\n"
  docker build -t "$REGISTRY/$CONTAINER_NAME:$TAG" -f "../../../../../$DOCKERFILE" ../../../../../
}

_load_minikube_cluster() {
  CONTAINER_NAME=$1
  echo -e "\nLoading $CONTAINER_NAME image...\n"
  minikube image load "$REGISTRY/$CONTAINER_NAME:$TAG"
}

cleanup_build_cache() {
  echo -e "\nCleanup Build Cache...\n"
  echo y | docker builder prune
}

run() {
  CONTAINER_NAME=$1
  DOCKERFILE=$2
  _build_containers "$CONTAINER_NAME" "$DOCKERFILE"
  _load_minikube_cluster "$CONTAINER_NAME"
}

echo "Building images for Katib ${VERSION}..."
echo "Image registry: ${REGISTRY}"
echo "Image tag: ${TAG}"

# Katib core images
run "katib-controller" "$CMD_PREFIX/katib-controller/$VERSION/Dockerfile"
run "katib-db-manager" "$CMD_PREFIX/db-manager/$VERSION/Dockerfile"
run "cert-generator" "$CMD_PREFIX/cert-generator/$VERSION/Dockerfile"
run "file-metrics-collector" "$CMD_PREFIX/metricscollector/$VERSION/file-metricscollector/Dockerfile"
run "tfevent-metrics-collector" "$CMD_PREFIX/metricscollector/$VERSION/tfevent-metricscollector/Dockerfile"
cleanup_build_cache

# Suggestion images
echo -e "\nBuilding suggestion images..."
run "suggestion-hyperopt" "$CMD_PREFIX/suggestion/hyperopt/$VERSION/Dockerfile"
run "suggestion-chocolate" "$CMD_PREFIX/suggestion/chocolate/$VERSION/Dockerfile"
run "suggestion-hyperband" "$CMD_PREFIX/suggestion/hyperband/$VERSION/Dockerfile"
run "suggestion-skopt" "$CMD_PREFIX/suggestion/skopt/$VERSION/Dockerfile"
run "suggestion-goptuna" "$CMD_PREFIX/suggestion/goptuna/$VERSION/Dockerfile"
run "suggestion-optuna" "$CMD_PREFIX/suggestion/optuna/$VERSION/Dockerfile"
run "suggestion-enas" "$CMD_PREFIX/suggestion/nas/enas/$VERSION/Dockerfile"
run "suggestion-darts" "$CMD_PREFIX/suggestion/nas/darts/$VERSION/Dockerfile"
cleanup_build_cache

# Early stopping images
echo -e "\nBuilding early stopping images...\n"
run "earlystopping-medianstop" "$CMD_PREFIX/earlystopping/medianstop/$VERSION/Dockerfile"
cleanup_build_cache

# Training container images
echo -e "\nBuilding training container images..."
for name in "${TRIAL_IMAGES[@]}"; do
  run "$name" "examples/$VERSION/trial-images/$name/Dockerfile"
done
cleanup_build_cache

echo -e "\nAll Katib images with ${TAG} tag have been built successfully!\n"
