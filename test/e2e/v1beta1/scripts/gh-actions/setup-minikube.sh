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

# This shell script is used to setup Katib deployment.

set -o errexit
set -o pipefail
set -o nounset
cd "$(dirname "$0")"

DEPLOY_KATIB_UI=${1:-false}
TRIAL_IMAGES=${2:-""}
EXPERIMENTS=${3:-""}

echo "Start to setup Minikube Kubernetes Cluster"
kubectl version
kubectl cluster-info
kubectl get nodes

echo "Build and Load container images"
./build-load.sh "$TRIAL_IMAGES" "$EXPERIMENTS" "$DEPLOY_KATIB_UI"
