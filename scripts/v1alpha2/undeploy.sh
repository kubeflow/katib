#!/bin/bash

# Copyright 2019 The Kubeflow Authors.
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

# Delete CR first
experiments=`kubectl get experiments --all-namespaces | awk 'NR>1' | awk '{print $1"/"$2}'`
for s in $experiments
do
  kubectl delete experiments $s --grace-period=0 --force;
done

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/../..

cd ${SCRIPT_ROOT}
kubectl delete -f manifests/v1alpha2/katib-controller
kubectl delete -f manifests/v1alpha2/katib/manager
kubectl delete -f manifests/v1alpha2/katib/manager-rest
kubectl delete -f manifests/v1alpha2/katib/db
kubectl delete -f manifests/v1alpha2/katib/ui
kubectl delete -f manifests/v1alpha2/katib/pv
kubectl delete -f manifests/v1alpha2/katib/suggestion/random
kubectl delete -f manifests/v1alpha2/katib/suggestion/bayesianoptimization
kubectl delete -f manifests/v1alpha2
cd - > /dev/null
