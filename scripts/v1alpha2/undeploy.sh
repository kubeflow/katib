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
experiments=`kubectl get experiments --all-namespaces | awk '{if (NR>1) {print $1"/"$2}}'`
for s in $experiments
do
  ns=`echo $s|cut -d "/" -f 1`
  exp=`echo $s|cut -d "/" -f 2`
  kubectl delete experiments $exp -n $ns
done

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/../..

cd ${SCRIPT_ROOT}
kubectl delete -f manifests/v1alpha2/katib-controller
kubectl delete -f manifests/v1alpha2/manager
kubectl delete -f manifests/v1alpha2/manager-rest
kubectl delete -f manifests/v1alpha2/db
kubectl delete -f manifests/v1alpha2/ui
kubectl delete -f manifests/v1alpha2/pv
kubectl delete -f manifests/v1alpha2/suggestion/random
kubectl delete -f manifests/v1alpha2/suggestion/bayesianoptimization
kubectl delete -f manifests/v1alpha2/suggestion/nasrl
kubectl delete -f manifests/v1alpha2/suggestion/grid
kubectl delete -f manifests/v1alpha2
cd - > /dev/null
