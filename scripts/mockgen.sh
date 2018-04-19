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

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..

cd ${SCRIPT_ROOT}
echo "Generating ManagerClient..."
mockgen -package mock -destination mock/api/manager.go github.com/kubeflow/hp-tuning/api ManagerClient
echo "Generating SuggestionClient..."
mockgen -package mock -destination mock/api/suggestion.go github.com/kubeflow/hp-tuning/api SuggestionClient
echo "Generating VizierDBInterface..."
mockgen -package mock -destination mock/db/db.go github.com/kubeflow/hp-tuning/db VizierDBInterface
echo "Generating WorkerInterface..."
mockgen -package mock -destination mock/worker/worker.go  github.com/kubeflow/hp-tuning/manager/worker_interface WorkerInterface
echo "Generating ModelStore..."
mockgen -package mock -destination mock/modelstore/modelstore.go  github.com/kubeflow/hp-tuning/manager/modelstore ModelStore
cd - > /dev/null
