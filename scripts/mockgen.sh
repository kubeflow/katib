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
echo "Generating v1alpha1 ManagerClient..."
mockgen -package mock -destination pkg/mock/v1alpha1/api/manager.go github.com/kubeflow/katib/pkg/api/v1alpha1 ManagerClient
echo "Generating v1alpha1 SuggestionClient..."
mockgen -package mock -destination pkg/mock/v1alpha1/api/suggestion.go github.com/kubeflow/katib/pkg/api/v1alpha1 SuggestionClient
echo "Generating v1alpha1 VizierDBInterface..."
mockgen -package mock -destination pkg/mock/v1alpha1/db/db.go github.com/kubeflow/katib/pkg/db/v1alpha1 VizierDBInterface
echo "Generating ModelStore..."
mockgen -package mock -destination pkg/mock/modelstore/modelstore.go  github.com/kubeflow/katib/pkg/manager/modelstore ModelStore

echo "Generating v1alpha2 ManagerClient..."
mockgen -package mock -destination pkg/mock/v1alpha2/api/manager.go github.com/kubeflow/katib/pkg/api/v1alpha2 ManagerClient
echo "Generating v1alpha2 SuggestionClient..."
mockgen -package mock -destination pkg/mock/v1alpha2/api/suggestion.go github.com/kubeflow/katib/pkg/api/v1alpha2 SuggestionClient
echo "Generating v1alpha2 KatibDBInterface..."
mockgen -package mock -destination pkg/mock/v1alpha2/db/db.go github.com/kubeflow/katib/pkg/db/v1alpha2 KatibDBInterface
