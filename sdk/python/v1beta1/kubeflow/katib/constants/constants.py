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

import os

# Global CRD version
KATIB_VERSION = os.environ.get('EXPERIMENT_VERSION', 'v1beta1')

# Katib K8S constants
KUBEFLOW_GROUP = 'kubeflow.org'
EXPERIMENT_PLURAL = 'experiments'
SUGGESTION_PLURAL = 'suggestions'
TRIAL_PLURAL = 'trials'

# How long to wait in seconds for requests to the ApiServer
APISERVER_TIMEOUT = 120
