#!/bin/bash

# Copyright 2021 The Kubeflow Authors.
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

# This script is used to update images and tags which are hosted by Katib community.
#
# The image names postfix must be equal to Katib format.
# Please check the images format here: https://github.com/kubeflow/katib/blob/master/docs/images-location.md
# For example, for Katib controller the image postfix is "katib-controller".
#
# This script updates the following images:
# 1. Katib Core Components.
# 2. Katib Metrics Collectors
# 3. Katib Suggestions
# 4. Katib Early Stopping
# 5. Katib Trial training containers
#
# Run ./scripts/v1beta1/update-images.sh -p <image-prefix> -t <image-tag> to execute it.
# For example, to update images with registry "docker.io/private/" and tag "v0.12.0", run:
# ./scripts/v1beta1/update-images.sh -p docker.io/private/ -t v0.12.0

usage() {
  echo "Usage: $0 [-p <IMAGE_PREFIX> -t <TAG>]" 1>&2
  exit 1
}

while getopts ":p:t:" opt; do
  case $opt in
  p)
    IMAGE_PREFIX=${OPTARG}
    ;;
  t)
    TAG=${OPTARG}
    ;;
  *)
    usage
    ;;
  esac
done

if [[ -z "$IMAGE_PREFIX" || -z "$TAG" ]]; then
  echo "Image prefix and tag must be set"
  echo "Usage: $0 [-p <IMAGE_PREFIX> -t <TAG>]" 1>&2
  exit 1
fi

# This function edits YAML files data for a given path.
# $1 argument - path for files to search.
# $2 argument - old string regex to be replaced.
# $3 argument - new string.
update_yaml_files() {
  # For MacOS we should set -i '' to avoid temp files from sed.
  if [[ $(uname) == "Darwin" ]]; then
    find $1 -regex ".*\.yaml" -exec sed -i '' -e "s@$2@$3@" {} \;
  else
    find $1 -regex ".*\.yaml" -exec sed -i -e "s@$2@$3@" {} \;
  fi
}

# Base prefix for the Katib images.
BASE_PREFIX="docker.io/kubeflowkatib/"

echo "Updating Katib images..."
echo "Image prefix: ${IMAGE_PREFIX}"
echo -e "Image tag: ${TAG}\n"

# Katib Core images.
# echo -e "Updating Katib Core images\n"
# update_yaml_files "manifests/v1beta1/installs/" "newName: ${BASE_PREFIX}" "newName: ${IMAGE_PREFIX}"
# update_yaml_files "manifests/v1beta1/installs/" "newTag: .*" "newTag: ${TAG}"

# # Katib Config images.
# CONFIG_PATH="manifests/v1beta1/components/controller/katib-config.yaml"

# echo -e "Update Katib Metrics Collectors, Suggestion and EarlyStopping images\n"
# update_yaml_files "${CONFIG_PATH}" "${BASE_PREFIX}" "${IMAGE_PREFIX}"
# update_yaml_files "${CONFIG_PATH}" ":[^[:space:]].*\"" ":${TAG}\""

# Katib Trial training container images.

# Postfix for the each Trial image.
MXNET_MNIST="mxnet-mnist"
PYTORCH_MNIST="pytorch-mnist"
ENAS_GPU="enas-cnn-cifar10-gpu"
ENAS_CPU="enas-cnn-cifar10-cpu"
DARTS="darts-cnn-cifar10"

echo -e "Update Katib Trial training container images\n"
update_yaml_files "examples/v1beta1/" "${BASE_PREFIX}${MXNET_MNIST}:.*" "${IMAGE_PREFIX}${MXNET_MNIST}:${TAG}"
update_yaml_files "examples/v1beta1/" "${BASE_PREFIX}${PYTORCH_MNIST}:.*" "${IMAGE_PREFIX}${PYTORCH_MNIST}:${TAG}"
update_yaml_files "examples/v1beta1/" "${BASE_PREFIX}${ENAS_GPU}:.*" "${IMAGE_PREFIX}${ENAS_GPU}:${TAG}"
update_yaml_files "examples/v1beta1/" "${BASE_PREFIX}${ENAS_CPU}:.*" "${IMAGE_PREFIX}${ENAS_CPU}:${TAG}"
update_yaml_files "examples/v1beta1/" "${BASE_PREFIX}${DARTS}:.*" "${IMAGE_PREFIX}${DARTS}:${TAG}"

echo -e "Katib images have been updated\n"
