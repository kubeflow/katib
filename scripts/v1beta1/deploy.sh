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
kubectl get validatingwebhookconfigurations katib-validating-webhook-config && kubectl delete validatingwebhookconfigurations katib-validating-webhook-config
kubectl get mutatingwebhookconfigurations katib-mutating-webhook-config && kubectl delete mutatingwebhookconfigurations katib-mutating-webhook-config
kubectl apply -f manifests/v1beta1
kubectl apply -f manifests/v1beta1/katib-controller
kubectl apply -f manifests/v1beta1/db-manager
kubectl apply -f manifests/v1beta1/pv
kubectl apply -f manifests/v1beta1/mysql-db
kubectl apply -f manifests/v1beta1/ui
cd - > /dev/null
