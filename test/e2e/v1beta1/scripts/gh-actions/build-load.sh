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
set -o nounset

pushd .
cd "$(dirname "$0")/../../../../.."
trap popd EXIT

TRIAL_IMAGES=${1:-""}
EXPERIMENTS=${2:-""}
DEPLOY_KATIB_UI=${3:-false}

REGISTRY="docker.io/kubeflowkatib"
TAG="e2e-test"
VERSION="v1beta1"
CMD_PREFIX="cmd"
SPECIFIED_DEVICE_TYPE_IMAGES=("enas-cnn-cifar10-cpu" "darts-cnn-cifar10-cpu" "pytorch-mnist-cpu")

IFS="," read -r -a TRIAL_IMAGE_ARRAY <<< "$TRIAL_IMAGES"
IFS="," read -r -a EXPERIMENT_ARRAY <<< "$EXPERIMENTS"

_build_containers() {
  CONTAINER_NAME=${1:-"katib-controller"}
  DOCKERFILE=${2:-"$CMD_PREFIX/katib-controller/$VERSION/Dockerfile"}

  for image in "${SPECIFIED_DEVICE_TYPE_IMAGES[@]}"; do
    if [ "$image" = "$CONTAINER_NAME" ]; then
      DOCKERFILE="${DOCKERFILE//-cpu/}"
      DOCKERFILE="${DOCKERFILE}.cpu"
      break
    fi
  done

  echo -e "\nBuilding $CONTAINER_NAME image with $DOCKERFILE...\n"
  DOCKER_BUILDKIT=1 minikube image build --build-opt platform=linux/amd64 --all -t "$REGISTRY/$CONTAINER_NAME:$TAG" -f "$DOCKERFILE" .
}

_install_tools() {
  # install yq
  if [ -z "$(command -v yq)" ]; then
    wget -O /usr/local/bin/yq "https://github.com/mikefarah/yq/releases/download/v4.25.2/yq_$(uname -s)_$(uname -m)"
    chmod +x /usr/local/bin/yq
  fi
}

run() {
  CONTAINER_NAME=${1:-"katib-controller"}
  DOCKERFILE=${2:-"$CMD_PREFIX/katib-controller/$VERSION/Dockerfile"}

  _install_tools

  # CONTAINER_NAME is image for suggestion services
  if echo "$CONTAINER_NAME" | grep -q "^suggestion-"; then

    suggestions=()

    # Search for Suggestion Images required for Trial.
    for exp_name in "${EXPERIMENT_ARRAY[@]}"; do

      exp_path=$(find examples/v1beta1 -name "${exp_name}.yaml")
      algorithm_name="$(yq eval '.spec.algorithm.algorithmName' "$exp_path")"

      suggestion_image_name="$(algorithm_name=$algorithm_name yq eval '.runtime.suggestions.[] | select(.algorithmName == env(algorithm_name)) | .image' \
        manifests/v1beta1/installs/katib-standalone/katib-config.yaml | cut -d: -f1)"
      suggestion_name="$(basename "$suggestion_image_name")"

      suggestions+=("$suggestion_name")

    done

    for s in "${suggestions[@]}"; do
      if [ "$s" == "$CONTAINER_NAME" ]; then
        _build_containers "$CONTAINER_NAME" "$DOCKERFILE"
        break
      fi
    done

  # $CONTAINER_NAME is image for earlystopping services
  elif echo "$CONTAINER_NAME" | grep -q "^earlystopping-"; then

    earlystoppings=()

    # Search for EarlyStopping Images required for Trial.
    for exp_name in "${EXPERIMENT_ARRAY[@]}"; do

      exp_path=$(find examples/v1beta1 -name "${exp_name}.yaml")
      algorithm_name="$(yq eval '.spec.earlyStopping.algorithmName' "$exp_path")"

      earlystopping_image_name="$(algorithm_name=$algorithm_name yq eval '.runtime.earlyStoppings.[] | select(.algorithmName == env(algorithm_name)) | .image' \
        manifests/v1beta1/installs/katib-standalone/katib-config.yaml | cut -d: -f1)"
      earlystopping_name="$(basename "$earlystopping_image_name")"

      earlystoppings+=("$earlystopping_name")

    done

    for e in "${earlystoppings[@]}"; do
      if [ "$e" == "$CONTAINER_NAME" ]; then
        _build_containers "$CONTAINER_NAME" "$DOCKERFILE"
        break
      fi
    done

  # Others
  else
    _build_containers "$CONTAINER_NAME" "$DOCKERFILE"
  fi
}

echo "Building images for Katib ${VERSION}..."
echo "Image registry: ${REGISTRY}"
echo "Image tag: ${TAG}"

# Katib core images
run "katib-controller" "$CMD_PREFIX/katib-controller/$VERSION/Dockerfile"
run "katib-db-manager" "$CMD_PREFIX/db-manager/$VERSION/Dockerfile"

if "$DEPLOY_KATIB_UI"; then
  run "katib-ui" "${CMD_PREFIX}/ui/${VERSION}/Dockerfile"
fi

run "file-metrics-collector" "$CMD_PREFIX/metricscollector/$VERSION/file-metricscollector/Dockerfile"
run "tfevent-metrics-collector" "$CMD_PREFIX/metricscollector/$VERSION/tfevent-metricscollector/Dockerfile"

# Suggestion images
echo -e "\nBuilding suggestion images..."
run "suggestion-hyperopt" "$CMD_PREFIX/suggestion/hyperopt/$VERSION/Dockerfile"
run "suggestion-hyperband" "$CMD_PREFIX/suggestion/hyperband/$VERSION/Dockerfile"
run "suggestion-skopt" "$CMD_PREFIX/suggestion/skopt/$VERSION/Dockerfile"
run "suggestion-goptuna" "$CMD_PREFIX/suggestion/goptuna/$VERSION/Dockerfile"
run "suggestion-optuna" "$CMD_PREFIX/suggestion/optuna/$VERSION/Dockerfile"
run "suggestion-pbt" "$CMD_PREFIX/suggestion/pbt/$VERSION/Dockerfile"
run "suggestion-enas" "$CMD_PREFIX/suggestion/nas/enas/$VERSION/Dockerfile"
run "suggestion-darts" "$CMD_PREFIX/suggestion/nas/darts/$VERSION/Dockerfile"

# Early stopping images
echo -e "\nBuilding early stopping images...\n"
run "earlystopping-medianstop" "$CMD_PREFIX/earlystopping/medianstop/$VERSION/Dockerfile"

# Training container images
echo -e "\nBuilding training container images..."
for name in "${TRIAL_IMAGE_ARRAY[@]}"; do
  run "$name" "examples/$VERSION/trial-images/$name/Dockerfile"
done

echo -e "\nCleanup Build Cache...\n"
docker buildx prune -f

echo -e "\nAll Katib images with ${TAG} tag have been built successfully!\n"
