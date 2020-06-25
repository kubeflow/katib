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
echo "Generating v1alpha3 ManagerClient..."
mockgen -package mock -destination pkg/mock/v1alpha3/api/manager.go github.com/kubeflow/katib/pkg/apis/manager/v1alpha3 ManagerClient
echo "Generating v1alpha3 SuggestionClient..."
mockgen -package mock -destination pkg/mock/v1alpha3/api/suggestion.go github.com/kubeflow/katib/pkg/apis/manager/v1alpha3 SuggestionClient
echo "Generating v1alpha3 KatibDBInterface..."
mockgen -package mock -destination pkg/mock/v1alpha3/db/db.go github.com/kubeflow/katib/pkg/db/v1alpha3/common KatibDBInterface
echo "Generating v1alpha3 Generator..."
mockgen -package mock -destination pkg/mock/v1alpha3/experiment/manifest/generator.go github.com/kubeflow/katib/pkg/controller.v1alpha3/experiment/manifest Generator
echo "Generating v1alpha3 KatibClient..."
mockgen -package mock -destination pkg/mock/v1alpha3/util/katibclient/katibclient.go github.com/kubeflow/katib/pkg/util/v1alpha3/katibclient Client
echo "Generating v1alpha3 ManagerClient in Trial Controller..."
mockgen -package mock -destination pkg/mock/v1alpha3/trial/managerclient/katibmanager.go github.com/kubeflow/katib/pkg/controller.v1alpha3/trial/managerclient ManagerClient
echo "Generating v1alpha3 Suggestion in Experiment Controller..."
mockgen -package mock -destination pkg/mock/v1alpha3/experiment/suggestion/suggestion.go github.com/kubeflow/katib/pkg/controller.v1alpha3/experiment/suggestion Suggestion

echo "Generating v1beta1 Suggestion RPC Client..."
mockgen -package mock -destination pkg/mock/v1beta1/api/suggestion.go github.com/kubeflow/katib/pkg/apis/manager/v1beta1 SuggestionClient
echo "Generating v1beta1 KatibDBInterface..."
mockgen -package mock -destination pkg/mock/v1beta1/db/db.go github.com/kubeflow/katib/pkg/db/v1beta1/common KatibDBInterface
echo "Generating v1beta1 Generator..."
mockgen -package mock -destination pkg/mock/v1beta1/experiment/manifest/generator.go github.com/kubeflow/katib/pkg/controller.v1beta1/experiment/manifest Generator
echo "Generating v1beta1 KatibClient..."
mockgen -package mock -destination pkg/mock/v1beta1/util/katibclient/katibclient.go github.com/kubeflow/katib/pkg/util/v1beta1/katibclient Client
echo "Generating v1beta1 ManagerClient in Trial Controller..."
mockgen -package mock -destination pkg/mock/v1beta1/trial/managerclient/katibmanager.go github.com/kubeflow/katib/pkg/controller.v1beta1/trial/managerclient ManagerClient
echo "Generating v1beta1 Suggestion in Experiment Controller..."
mockgen -package mock -destination pkg/mock/v1beta1/experiment/suggestion/suggestion.go github.com/kubeflow/katib/pkg/controller.v1beta1/experiment/suggestion Suggestion
