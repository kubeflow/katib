#!/bin/bash

# Copyright 2018 The Kubeflow Authors.
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

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/../..

cd ${SCRIPT_ROOT}
kubectl apply -f manifests/v1alpha2
kubectl apply -f manifests/v1alpha2/katib-controller
kubectl apply -f manifests/v1alpha2/manager
kubectl apply -f manifests/v1alpha2/manager-rest
kubectl apply -f manifests/v1alpha2/pv
kubectl apply -f manifests/v1alpha2/db
kubectl apply -f manifests/v1alpha2/ui
kubectl apply -f manifests/v1alpha2/suggestion/random
kubectl apply -f manifests/v1alpha2/suggestion/bayesianoptimization
kubectl apply -f manifests/v1alpha2/suggestion/nasrl
kubectl apply -f manifests/v1alpha2/suggestion/grid
cd - > /dev/null
