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

# This shell script is used to run Katib Experiment.
# Input parameter - path to Experiment yaml.

set -o errexit
set -o nounset
set -o pipefail

cd "$(dirname "$0")"
EXPERIMENT_FILES=${1:-""}
IFS="," read -r -a EXPERIMENT_FILE_ARRAY <<< "$EXPERIMENT_FILES"

echo "Katib deployments"
kubectl -n kubeflow get deploy
echo "Katib services"
kubectl -n kubeflow get svc
echo "Katib pods"
kubectl -n kubeflow get pod
echo "Katib persistent volume claims"
kubectl get pvc -n kubeflow
echo "Available CRDs"
kubectl get crd

if [ -z "$EXPERIMENT_FILES" ]; then
  echo "Skip Test for Experiment"
  exit 0
fi

for exp_name in "${EXPERIMENT_FILE_ARRAY[@]}"; do
  echo "Running Experiment from $exp_name file"
  exp_path=$(find ../../../../../examples/v1beta1 -name "${exp_name}.yaml")
  python run-e2e-experiment.py --experiment-path "${exp_path}" --namespace default \
  --verbose || (kubectl get pods -n kubeflow && exit 1)
done

exit 0
