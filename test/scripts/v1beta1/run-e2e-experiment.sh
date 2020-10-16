#!/bin/bash

# Copyright 2020 The Kubeflow Authors.
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

# This shell script is used to run Katib Experiment.
# Input parameter - path to Experiment yaml.

set -o errexit
set -o nounset
set -o pipefail

CLUSTER_NAME="${CLUSTER_NAME}"
AWS_REGION="${AWS_REGION}"
EXPERIMENT_FILE=$1

# GO_DIR=${GOPATH}/src/github.com/${REPO_OWNER}/${REPO_NAME}

echo "Configuring kubeconfig.."
aws eks update-kubeconfig --region=${REGION} --name=${CLUSTER_NAME}

echo "Katib deployments"
kubectl -n kubeflow get deploy
echo "Katib services"
kubectl -n kubeflow get svc
echo "Katib pods"
kubectl -n kubeflow get pod
echo "Katib persistent volume claims "
kubectl get pvc -n kubeflow

# cd ${GO_DIR}/test/e2e/v1beta1

echo "Running Experiment from ${EXPERIMENT_FILE} file"
go run ./test/scripts/v1beta1/run-e2e-experiment ${EXPERIMENT_FILE}

exit 0
