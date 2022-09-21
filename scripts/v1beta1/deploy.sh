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

set -o xtrace
set -o errexit

SCRIPT_ROOT="$(dirname "${BASH_SOURCE[0]}")/../.."

cd "${SCRIPT_ROOT}"

WITH_DATABASE_TYPE=${1:-mysql}

# if mysql, use below kustomize, else use postgres
if [ "$WITH_DATABASE_TYPE" == "mysql" ]; then
    kustomize build manifests/v1beta1/installs/katib-standalone | kubectl apply -f -
elif [ "$WITH_DATABASE_TYPE" == "postgres" ]; then
    kustomize build manifests/v1beta1/installs/katib-standalone-postgres | kubectl apply -f -
else
    echo "Unknown database type: $WITH_DATABASE_TYPE"
    exit 1
fi
