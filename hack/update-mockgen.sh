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

# This script generates files using mockgen.
# Usage: `hack/update-mockgen.sh`.

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT="$(dirname "${BASH_SOURCE[0]}")/.."

cd "${SCRIPT_ROOT}"

# Grab mockgen version from go.mod
MOCKGEN_VERSION=$(grep 'github.com/golang/mock' go.mod | awk '{print $2}')

MOCKGEN="./bin/mockgen-${MOCKGEN_VERSION}"
if [ ! -f "${MOCKGEN}" ]; then
  echo "Please run 'make mockgen' first"
  exit 1
fi

echo "Generating v1beta1 Suggestion RPC Client..."
"${MOCKGEN}" -package mock -destination pkg/mock/v1beta1/api/suggestion.go github.com/kubeflow/katib/pkg/apis/manager/v1beta1 SuggestionClient
echo "Generating v1beta1 EarlyStopping RPC Client..."
"${MOCKGEN}" -package mock -destination pkg/mock/v1beta1/api/earlystopping.go github.com/kubeflow/katib/pkg/apis/manager/v1beta1 EarlyStoppingClient
echo "Generating v1beta1 KatibDBInterface..."
"${MOCKGEN}" -package mock -destination pkg/mock/v1beta1/db/db.go github.com/kubeflow/katib/pkg/db/v1beta1/common KatibDBInterface
echo "Generating v1beta1 Generator..."
"${MOCKGEN}" -package mock -destination pkg/mock/v1beta1/experiment/manifest/generator.go github.com/kubeflow/katib/pkg/controller.v1beta1/experiment/manifest Generator
echo "Generating v1beta1 KatibClient..."
"${MOCKGEN}" -package mock -destination pkg/mock/v1beta1/util/katibclient/katibclient.go github.com/kubeflow/katib/pkg/util/v1beta1/katibclient Client
echo "Generating v1beta1 ManagerClient in Trial Controller..."
"${MOCKGEN}" -package mock -destination pkg/mock/v1beta1/trial/managerclient/katibmanager.go github.com/kubeflow/katib/pkg/controller.v1beta1/trial/managerclient ManagerClient
echo "Generating v1beta1 Suggestion in Experiment Controller..."
"${MOCKGEN}" -package mock -destination pkg/mock/v1beta1/experiment/suggestion/suggestion.go github.com/kubeflow/katib/pkg/controller.v1beta1/experiment/suggestion Suggestion
echo "Generating v1beta1 SuggestionClient in Suggestion Controller..."
"${MOCKGEN}" -package mock -destination pkg/mock/v1beta1/suggestion/suggestionclient/suggestionclient.go github.com/kubeflow/katib/pkg/controller.v1beta1/suggestion/suggestionclient SuggestionClient
