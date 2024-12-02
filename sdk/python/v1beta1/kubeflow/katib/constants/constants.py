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

# How long to wait in seconds for requests to the Kubernetes or gRPC API Server.
DEFAULT_TIMEOUT = 120

# RFC3339 time format
RFC3339_FORMAT = "%Y-%m-%dT%H:%M:%SZ"

# Global CRD version
KATIB_VERSION = os.environ.get("EXPERIMENT_VERSION", "v1beta1")

# Katib K8S constants
KUBEFLOW_GROUP = "kubeflow.org"
API_VERSION = f"{KUBEFLOW_GROUP}/{KATIB_VERSION}"
EXPERIMENT_KIND = "Experiment"
EXPERIMENT_PLURAL = "experiments"
SUGGESTION_PLURAL = "suggestions"
TRIAL_PLURAL = "trials"


DEFAULT_PRIMARY_CONTAINER_NAME = "training-container"

# Label to identify Experiment's resources.
EXPERIMENT_LABEL = "katib.kubeflow.org/experiment"

# True means that Katib CR is in this condition.
CONDITION_STATUS_TRUE = "True"

# Experiment conditions.
# TODO (andreyvelich): Use API enums when Katib SDK supports it.
# Ref: https://github.com/kubeflow/katib/issues/1969.
EXPERIMENT_CONDITION_CREATED = "Created"
EXPERIMENT_CONDITION_RUNNING = "Running"
EXPERIMENT_CONDITION_RESTARTING = "Restarting"
EXPERIMENT_CONDITION_SUCCEEDED = "Succeeded"
EXPERIMENT_CONDITION_FAILED = "Failed"

# Trial conditions.
TRIAL_CONDITION_SUCCEEDED = "Succeeded"

# Supported base images for the Katib Trials.
# TODO (andreyvelich): Implement list_base_images function to get each image description.
BASE_IMAGE_TENSORFLOW = "docker.io/tensorflow/tensorflow:2.13.0"
BASE_IMAGE_TENSORFLOW_GPU = "docker.io/tensorflow/tensorflow:2.13.0-gpu"
BASE_IMAGE_PYTORCH = "docker.io/pytorch/pytorch:2.2.1-cuda12.1-cudnn8-runtime"

DEFAULT_DB_MANAGER_ADDRESS = "katib-db-manager.kubeflow:6789"

# The default value for dataset and model storage PVC.
PVC_DEFAULT_SIZE = "10Gi"
# The default value for PVC access modes.
PVC_DEFAULT_ACCESS_MODES = ["ReadWriteOnce", "ReadOnlyMany"]
