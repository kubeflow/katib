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

# Delete CR first
studyjobs=`kubectl get studyjobs --all-namespaces | awk 'NR>1' | awk '{print $1"/"$2}'`
for s in $studyjobs
do
  kubectl delete studyjobs $s --grace-period=0 --force;
done

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/../..

cd ${SCRIPT_ROOT}
kubectl delete -f manifests/v1alpha1/pv
kubectl delete -f manifests/v1alpha1/vizier/db
kubectl delete -f manifests/v1alpha1/vizier/core
kubectl delete -f manifests/v1alpha1/vizier/core-rest
kubectl delete -f manifests/v1alpha1/vizier/ui
kubectl delete -f manifests/v1alpha1/vizier/suggestion/random
kubectl delete -f manifests/v1alpha1/vizier/suggestion/grid
kubectl delete -f manifests/v1alpha1/vizier/suggestion/hyperband
kubectl delete -f manifests/v1alpha1/vizier/suggestion/bayesianoptimization
kubectl delete -f manifests/v1alpha1/vizier/suggestion/nasrl
kubectl delete -f manifests/v1alpha1/vizier/suggestion/nasenvelopenet
kubectl delete -f manifests/v1alpha1/vizier/earlystopping/medianstopping
kubectl delete -f manifests/v1alpha1/studyjobcontroller
kubectl delete -f manifests/v1alpha1
cd - > /dev/null
