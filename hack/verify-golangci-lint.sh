#!/usr/bin/env bash

# Copyright 2021 The Kubeflow Authors.
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

set -o errexit
set -o pipefail

cd $(dirname "$0")/..

if ! which golangci-lint >/dev/null; then
	echo 'Can not find golangci-lint, install with:'
	echo 'go get -u github.com/golangci/golangci-lint/cmd/golangci-lint'
	exit 1
fi

echo 'running golangci-lint'
golangci-lint run --timeout 5m
