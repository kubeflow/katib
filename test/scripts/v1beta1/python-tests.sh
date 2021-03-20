#!/bin/bash

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

# This shell script is used to run the python tests in the argo workflow

set -o errexit
set -o nounset
set -o pipefail

export PYTHONPATH=$(pwd):$(pwd)/pkg/apis/manager/v1beta1/python:$(pwd)/pkg/apis/manager/health/python
pip install -r test/suggestion/v1beta1/test_requirements.txt
pip install -r cmd/suggestion/chocolate/v1beta1/requirements.txt
pip install -r cmd/suggestion/hyperopt/v1beta1/requirements.txt
pip install -r cmd/suggestion/skopt/v1beta1/requirements.txt
pip install -r cmd/suggestion/nas/enas/v1beta1/requirements.txt
pip install -r cmd/suggestion/hyperband/v1beta1/requirements.txt
pip install -r cmd/suggestion/nas/darts/v1beta1/requirements.txt
pytest -s ./test/suggestion/v1beta1

pip install -r cmd/earlystopping/medianstop/v1beta1/requirements.txt
pytest -s ./test/earlystopping/v1beta1
