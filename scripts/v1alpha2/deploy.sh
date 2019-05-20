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
kubectl apply -f manifests/v1alpha2/dbif/mysql/
kubectl apply -f manifests/v1alpha2/katib/manager
kubectl apply -f manifests/v1alpha2/katib/manager-rest
kubectl apply -f manifests/v1alpha2/katib/pv
kubectl apply -f manifests/v1alpha2/katib/db
kubectl apply -f manifests/v1alpha2/katib/suggestion/random
cd - > /dev/null
