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
kubectl apply -f manifests/v1alpha1/0-namespace.yaml
kubectl apply -f manifests/v1alpha1/pv
kubectl apply -f manifests/v1alpha1/vizier/db
kubectl apply -f manifests/v1alpha1/vizier/core
kubectl apply -f manifests/v1alpha1/vizier/core-rest
kubectl apply -f manifests/v1alpha1/vizier/ui
kubectl apply -f manifests/v1alpha1/vizier/suggestion/random
kubectl apply -f manifests/v1alpha1/vizier/suggestion/grid
kubectl apply -f manifests/v1alpha1/vizier/suggestion/hyperband
kubectl apply -f manifests/v1alpha1/vizier/suggestion/bayesianoptimization
kubectl apply -f manifests/v1alpha1/vizier/suggestion/nasrl
kubectl apply -f manifests/v1alpha1/vizier/earlystopping/medianstopping
kubectl apply -f manifests/v1alpha1/studyjobcontroller/crd.yaml
kubectl apply -f manifests/v1alpha1/studyjobcontroller/tfjobsCrd.yaml
kubectl apply -f manifests/v1alpha1/studyjobcontroller/pytorchjobsCrd.yaml
kubectl apply -f manifests/v1alpha1/studyjobcontroller/rbac.yaml
kubectl apply -f manifests/v1alpha1/studyjobcontroller/mcrbac.yaml
kubectl apply -f manifests/v1alpha1/studyjobcontroller/workerConfigMap.yaml
kubectl apply -f manifests/v1alpha1/studyjobcontroller/metricsControllerConfigMap.yaml
kubectl apply -f manifests/v1alpha1/studyjobcontroller/studyjobcontroller.yaml
cd - > /dev/null
