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

import inspect
import json
import os
import textwrap
from typing import Any, Callable, Dict, List, Optional, Tuple, Union

from kubeflow.katib import models
from kubeflow.katib.constants import constants


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


def validate_metrics_value(value: Any):
    """Validate if the metrics value can be converted to type `float`."""
    try:
        float(value)
    except Exception:
        raise ValueError(
            f"Invalid value {value} for metrics value. "
            "The metrics value should have or can be converted to type `float`. "
        )


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


def get_command_using_train_func(
    train_func: Optional[Callable],
    train_func_parameters: Optional[Dict[str, Any]] = None,
    packages_to_install: Optional[List[str]] = None,
    pip_index_url: str = "https://pypi.org/simple",
) -> Tuple[List[str], List[str]]:
    """
    Get container args and command from the given training function and parameters.
    """
    # Check if function is callable.
    if not callable(train_func):
        raise ValueError(
            f"Training function must be callable, got function type: {type(train_func)}"
        )

    # Extract function implementation.
    func_code = inspect.getsource(train_func)

    # Function might be defined in some indented scope (e.g. in another function).
    # We need to dedent the function code.
    func_code = textwrap.dedent(func_code)

    # Wrap function code to execute it from the file. For example:
    # def train(parameters):
    #     print('Start Training...')
    # train({'lr': 0.01})
    if train_func_parameters is None:
        func_code = f"{func_code}\n{train_func.__name__}()\n"
    else:
        func_code = f"{func_code}\n{train_func.__name__}({train_func_parameters})\n"

    # Prepare execute script template.
    exec_script = textwrap.dedent(
        """
                program_path=$(mktemp -d)
                read -r -d '' SCRIPT << EOM\n
                {func_code}
                EOM
                printf "%s" \"$SCRIPT\" > \"$program_path/ephemeral_script.py\"
                python3 -u \"$program_path/ephemeral_script.py\""""
    )

    # Add function code to the execute script.
    exec_script = exec_script.format(func_code=func_code)

    # Install Python packages if that is required.
    if packages_to_install is not None:
        exec_script = (
            get_script_for_python_packages(packages_to_install, pip_index_url) + exec_script
        )

    # Return container command and args to execute training function.
    return ["bash", "-c"], [exec_script]


def get_container_spec(
    name: str,
    base_image: str,
    train_func: Optional[Callable] = None,
    train_func_parameters: Optional[Dict[str, Any]] = None,
    packages_to_install: Optional[List[str]] = None,
    pip_index_url: str = "https://pypi.org/simple",
    args: Optional[List[str]] = None,
    resources: Union[dict, models.V1ResourceRequirements, None] = None,
    volume_mounts: Optional[List[models.V1VolumeMount]] = None,
    env: Optional[List[models.V1EnvVar]] = None,
    env_from: Optional[List[models.V1EnvFromSource]] = None,
) -> models.V1Container:
    """
    Get container spec for the given parameters.
    """

    if name is None or base_image is None:
        raise ValueError("Container name or base image cannot be none")

    # Create initial container spec.
    container_spec = models.V1Container(
        name=name, image=base_image, args=args, volume_mounts=volume_mounts
    )

    # If training function is set, override container command and args to execute the function.
    if train_func is not None:
        container_spec.command, container_spec.args = get_command_using_train_func(
            train_func=train_func,
            train_func_parameters=train_func_parameters,
            packages_to_install=packages_to_install,
            pip_index_url=pip_index_url,
        )

    # Convert dict to the Kubernetes container resources if that is required.
    if isinstance(resources, dict):
        # Convert all keys in resources to lowercase.
        resources = {k.lower(): v for k, v in resources.items()}
        if "gpu" in resources:
            resources["nvidia.com/gpu"] = resources.pop("gpu")

        resources = models.V1ResourceRequirements(
            requests=resources,
            limits=resources,
        )

    # Add resources to the container spec.
    container_spec.resources = resources

    # Add environment variables to the container spec.
    if env:
        container_spec.env = env
    if env_from:
        container_spec.env_from = env_from

    return container_spec


def get_pod_template_spec(
    containers: List[models.V1Container],
    init_containers: Optional[List[models.V1Container]] = None,
    volumes: Optional[List[models.V1Volume]] = None,
    restart_policy: Optional[str] = None,
) -> models.V1PodTemplateSpec:
    """
    Get Pod template spec for the given parameters.
    """

    # Create Pod template spec. If the value is None, Pod doesn't have that parameter
    pod_template_spec = models.V1PodTemplateSpec(
        metadata=models.V1ObjectMeta(annotations={"sidecar.istio.io/inject": "false"}),
        spec=models.V1PodSpec(
            init_containers=init_containers,
            containers=containers,
            volumes=volumes,
            restart_policy=restart_policy,
        ),
    )

    return pod_template_spec


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
        metadata={"name": pvc_name, "namespace": namespace},
        spec=models.V1PersistentVolumeClaimSpec(
            access_modes=storage_config["access_modes"],
            resources=models.V1ResourceRequirements(requests={"storage": storage_config["size"]}),
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
