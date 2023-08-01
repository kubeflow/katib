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
# Run ./scripts/v1beta1/update-images.sh <OLD_PREFIX> <NEW_PREFIX> <TAG> to execute it.
# For example, to update images from: docker.io/kubeflowkatib/ to: docker.io/private/ registry with tag: v0.12.0, run:
# ./scripts/v1beta1/update-images.sh docker.io/kubeflowkatib/ docker.io/private/ v0.12.0

set -o errexit
set -o pipefail
set -o nounset

OLD_PREFIX=${1:-""}
NEW_PREFIX=${2:-""}
TAG=${3:-""}

if [[ -z "$OLD_PREFIX" || -z "$NEW_PREFIX" || -z "$TAG" ]]; then
  echo "Image old prefix, new prefix, and tag must be set"
  echo -e "Usage: $0 <OLD_PREFIX> <NEW_PREFIX> <TAG>\n" 1>&2
  echo "For example, to update images from: docker.io/kubeflowkatib/ to: docker.io/private/ registry with tag: v0.12.0, run:"
  echo "$0 docker.io/kubeflowkatib/ docker.io/private/ v0.12.0"
  exit 1
fi

# This function edits YAML files data for a given path.
# $1 argument - path for files to search.
# $2 argument - old string regex to be replaced.
# $3 argument - new string.
update_yaml_files() {
  # For MacOS we should set -i '' to avoid temp files from sed.
  if [[ $(uname) == "Darwin" ]]; then
    find "$1" -regex ".*\.yaml" -exec sed -i '' -e "s@$2@$3@" {} \;
  else
    find "$1" -regex ".*\.yaml" -exec sed -i -e "s@$2@$3@" {} \;
  fi
}

echo "Updating Katib images..."
echo "Image old prefix: ${OLD_PREFIX}"
echo "Image new prefix: ${NEW_PREFIX}"
echo -e "Image tag: ${TAG}\n"

# Katib Core images.
INSTALLS_PATH="manifests/v1beta1/installs/"

echo -e "Updating Katib Core images\n"
update_yaml_files "${INSTALLS_PATH}" "newName: ${OLD_PREFIX}" "newName: ${NEW_PREFIX}"
update_yaml_files "${INSTALLS_PATH}" "newTag: .*" "newTag: ${TAG}"

# Katib Config images.
echo -e "Update Katib Metrics Collectors, Suggestions and EarlyStopping images\n"
for config in manifests/v1beta1/installs/**/katib-config.yaml; do
  update_yaml_files "${config}" "${OLD_PREFIX}" "${NEW_PREFIX}"
  update_yaml_files "${config}" ":[^[:space:]].*" ":${TAG}"
done

# Katib Trial training container images.

# Postfixes for the each Trial image.
MXNET_MNIST="mxnet-mnist"
PYTORCH_MNIST_CPU="pytorch-mnist-cpu"
PYTORCH_MNIST_GPU="pytorch-mnist-gpu"
TF_MNIST_WITH_SUMMARIES="tf-mnist-with-summaries"
ENAS_GPU="enas-cnn-cifar10-gpu"
ENAS_CPU="enas-cnn-cifar10-cpu"
DARTS_GPU="darts-cnn-cifar10-gpu"
DARTS_CPU="darts-cnn-cifar10-cpu"
SIMPLE_PBT="simple-pbt"

echo -e "Update Katib Trial training container images\n"
update_yaml_files "./" "${OLD_PREFIX}${MXNET_MNIST}:.*" "${NEW_PREFIX}${MXNET_MNIST}:${TAG}"
update_yaml_files "./" "${OLD_PREFIX}${PYTORCH_MNIST_CPU}:.*" "${NEW_PREFIX}${PYTORCH_MNIST_CPU}:${TAG}"
update_yaml_files "./" "${OLD_PREFIX}${PYTORCH_MNIST_GPU}:.*" "${NEW_PREFIX}${PYTORCH_MNIST_GPU}:${TAG}"
update_yaml_files "./" "${OLD_PREFIX}${TF_MNIST_WITH_SUMMARIES}:.*" "${NEW_PREFIX}${TF_MNIST_WITH_SUMMARIES}:${TAG}"
update_yaml_files "./" "${OLD_PREFIX}${ENAS_GPU}:.*" "${NEW_PREFIX}${ENAS_GPU}:${TAG}"
update_yaml_files "./" "${OLD_PREFIX}${ENAS_CPU}:.*" "${NEW_PREFIX}${ENAS_CPU}:${TAG}"
update_yaml_files "./" "${OLD_PREFIX}${DARTS_GPU}:.*" "${NEW_PREFIX}${DARTS_GPU}:${TAG}"
update_yaml_files "./" "${OLD_PREFIX}${DARTS_CPU}:.*" "${NEW_PREFIX}${DARTS_CPU}:${TAG}"
update_yaml_files "./" "${OLD_PREFIX}${SIMPLE_PBT}:.*" "${NEW_PREFIX}${SIMPLE_PBT}:${TAG}"

echo "Katib images have been updated"
