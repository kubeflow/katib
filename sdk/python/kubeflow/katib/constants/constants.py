# Copyright 2019 kubeflow.org.
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

# Katib K8S constants
EXPERIMENT_GROUP = 'kubeflow.org'
EXPERIMENT_KIND = 'experiment'
EXPERIMENT_PLURAL = 'experiments'
EXPERIMENT_VERSION = os.environ.get('EXPERIMENT_VERSION', 'v1alpha3')
EXPERIMENT_LOGLEVEL = os.environ.get('EXPERIMENT_LOGLEVEL', 'INFO').upper()

TRIAL_KIND = 'trial'
TRIAL_PLURAL = 'trials'
TRIAL_VERSION = os.environ.get('EXPERIMENT_VERSION', 'v1alpha3')

# How long to wait in seconds for requests to the ApiServer
APISERVER_TIMEOUT = 120
