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

# This script is used to update images and tags in the Trial templates for all Katib examples and manifests.
# Run ./scripts/v1beta1/update-trial-template-tags.sh -p <image-prefix> -t <image-tag> to execute it.

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

echo "Updating Trial template images..."
echo "Image prefix: ${IMAGE_PREFIX}"
echo "Image tag: ${TAG}"

# Base prefix for the Trial template images.
BASE_IMAGE_PREFIX="docker.io/kubeflowkatib/"

# End of the each Trial template image.
MXNET_MNIST="mxnet-mnist"
PYTORCH_MNIST="pytorch-mnist"
ENAS_GPU="enas-cnn-cifar10-gpu"
ENAS_CPU="enas-cnn-cifar10-cpu"
DARTS="darts-cnn-cifar10"

# MXNet mnist.
# For MacOS we should set -i '' to avoid temp files from sed.
if [[ $(uname) == "Darwin" ]]; then
  find ./ -regex ".*\.yaml" -exec sed -i '' -e "s@${BASE_IMAGE_PREFIX}${MXNET_MNIST}:.*@${IMAGE_PREFIX}${MXNET_MNIST}:${TAG}@" {} \;
  find ./ -regex ".*\.yaml" -exec sed -i '' -e "s@${BASE_IMAGE_PREFIX}${PYTORCH_MNIST}:.*@${IMAGE_PREFIX}${PYTORCH_MNIST}:${TAG}@" {} \;
  find ./ -regex ".*\.yaml" -exec sed -i '' -e "s@${BASE_IMAGE_PREFIX}${ENAS_GPU}:.*@${IMAGE_PREFIX}${ENAS_GPU}:${TAG}@" {} \;
  find ./ -regex ".*\.yaml" -exec sed -i '' -e "s@${BASE_IMAGE_PREFIX}${ENAS_CPU}:.*@${IMAGE_PREFIX}${ENAS_CPU}:${TAG}@" {} \;
  find ./ -regex ".*\.yaml" -exec sed -i '' -e "s@${BASE_IMAGE_PREFIX}${DARTS}:.*@${IMAGE_PREFIX}${DARTS}:${TAG}@" {} \;
else
  find ./ -regex ".*\.yaml" -exec sed -i -e "s@${BASE_IMAGE_PREFIX}${MXNET_MNIST}:.*@${IMAGE_PREFIX}${MXNET_MNIST}:${TAG}@" {} \;
  find ./ -regex ".*\.yaml" -exec sed -i -e "s@${BASE_IMAGE_PREFIX}${PYTORCH_MNIST}:.*@${IMAGE_PREFIX}${PYTORCH_MNIST}:${TAG}@" {} \;
  find ./ -regex ".*\.yaml" -exec sed -i -e "s@${BASE_IMAGE_PREFIX}${ENAS_GPU}:.*@${IMAGE_PREFIX}${ENAS_GPU}:${TAG}@" {} \;
  find ./ -regex ".*\.yaml" -exec sed -i -e "s@${BASE_IMAGE_PREFIX}${ENAS_CPU}:.*@${IMAGE_PREFIX}${ENAS_CPU}:${TAG}@" {} \;
  find ./ -regex ".*\.yaml" -exec sed -i -e "s@${BASE_IMAGE_PREFIX}${DARTS}:.*@${IMAGE_PREFIX}${DARTS}:${TAG}@" {} \;
fi

echo "Trial template images has been updated"
