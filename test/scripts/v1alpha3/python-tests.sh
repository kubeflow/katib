#!/bin/bash

# Copyright 2018 The Kubernetes Authors.
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

# This shell script is used to run the python tests in the argo workflow

set -o errexit
set -o nounset
set -o pipefail

export PYTHONPATH=$(pwd):$(pwd)/pkg/apis/manager/v1alpha3/python:$(pwd)/pkg/apis/manager/health/python
pip install -r test/suggestion/v1alpha3/test_requirements.txt
pip install -r cmd/suggestion/chocolate/v1alpha3/requirements.txt
pip install -r cmd/suggestion/hyperopt/v1alpha3/requirements.txt
pip install -r cmd/suggestion/skopt/v1alpha3/requirements.txt
pytest -s ./test
