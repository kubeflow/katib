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

import json
import os
import textwrap
from typing import Callable
import inspect

from kubeflow.katib import models
from kubeflow.katib.constants import constants
from typing import Any, Callable, Dict, List, Optional, Union
import transformers


def is_running_in_k8s():
    return os.path.isdir("/var/run/secrets/kubernetes.io/")


def get_current_k8s_namespace():
    with open("/var/run/secrets/kubernetes.io/serviceaccount/namespace", "r") as f:
        return f.readline()


def get_default_target_namespace():
    if not is_running_in_k8s():
        return "default"
    return get_current_k8s_namespace()


def set_katib_namespace(katib):
    katib_namespace = katib.metadata.namespace
    namespace = katib_namespace or get_default_target_namespace()
    return namespace


def has_condition(conditions, condition_type):
    """Verify if the condition list has the required condition.
    Condition should be valid object with `type` and `status`.
    """

    for c in conditions:
        if c.type == condition_type and c.status == constants.CONDITION_STATUS_TRUE:
            return True
    return False


def print_experiment_status(experiment: models.V1beta1Experiment):
    if experiment.status:
        print(
            "Experiment Trials status: {} Trials, {} Pending Trials, "
            "{} Running Trials, {} Succeeded Trials, {} Failed Trials, "
            "{} EarlyStopped Trials, {} MetricsUnavailable Trials".format(
                experiment.status.trials or 0,
                experiment.status.trials_pending or 0,
                experiment.status.trials_running or 0,
                experiment.status.trials_succeeded or 0,
                experiment.status.trials_failed or 0,
                experiment.status.trials_early_stopped or 0,
                experiment.status.trial_metrics_unavailable or 0,
            )
        )
        print(f"Current Optimal Trial:\n {experiment.status.current_optimal_trial}")
        print(f"Experiment conditions:\n {experiment.status.conditions}")


def validate_objective_function(objective: Callable):

    # Check if objective function is callable.
    if not callable(objective):
        raise ValueError(
            f"Objective function must be callable, got function type: {type(objective)}"
        )

    # Verify the objective function arguments.
    objective_signature = inspect.signature(objective)
    try:
        objective_signature.bind({})
    except Exception:
        raise ValueError(
            "Invalid args in the Objective function. "
            "The function args must have only 'parameters' dictionary. "
            f"Current Objective arguments: {objective_signature}"
        )


def get_script_for_python_packages(packages_to_install, pip_index_url):
    packages_str = " ".join([str(package) for package in packages_to_install])

    script_for_python_packages = textwrap.dedent(
        f"""
        if ! [ -x "$(command -v pip)" ]; then
            python3 -m ensurepip || python3 -m ensurepip --user || apt-get install python3-pip
        fi

        PIP_DISABLE_PIP_VERSION_CHECK=1 python3 -m pip install --prefer-binary --quiet \
        --no-warn-script-location --index-url {pip_index_url} {packages_str}
        """
    )

    return script_for_python_packages


class FakeResponse:
    """Fake object of RESTResponse to deserialize
    Ref) https://github.com/kubeflow/katib/pull/1630#discussion_r697877815
    Ref) https://github.com/kubernetes-client/python/issues/977#issuecomment-592030030
    """

    def __init__(self, obj):
        self.data = json.dumps(obj)


def get_pvc_spec(
    pvc_name: str,
    namespace: str,
    storage_config: Dict[str, Optional[Union[str, List[str]]]],
):
    if pvc_name is None or namespace is None:
        raise ValueError("One of the required storage config argument is None")

    if "size" not in storage_config:
        storage_config["size"] = constants.PVC_DEFAULT_SIZE

    if "access_modes" not in storage_config:
        storage_config["access_modes"] = constants.PVC_DEFAULT_ACCESS_MODES

    pvc_spec = models.V1PersistentVolumeClaim(
        api_version="v1",
        kind="PersistentVolumeClaim",
        metadata={"name": pvc_name, "namepsace": namespace},
        spec=models.V1PersistentVolumeClaimSpec(
            access_modes=storage_config["access_modes"],
            resources=models.V1ResourceRequirements(
                requests={"storage": storage_config["size"]}
            ),
        ),
    )

    if "storage_class" in storage_config:
        pvc_spec.spec.storage_class_name = storage_config["storage_class"]

    return pvc_spec


class SetEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, set):
            return list(obj)
        if isinstance(obj, type):
            return obj.__name__
        return json.JSONEncoder.default(self, obj)