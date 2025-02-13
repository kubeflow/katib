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

import copy
import inspect
import json
import logging
import os
import textwrap
from typing import Any, Callable, Dict, List, Optional, Union

from kubeflow.katib import models
from kubeflow.katib.constants import constants
from kubeflow.katib.types import types
from kubeflow.training import models as training_models
from kubeflow.training.constants.constants import (
    API_VERSION,
    JOB_PARAMETERS,
    PYTORCHJOB_KIND,
)
from kubernetes import client

logger = logging.getLogger(__name__)


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


class SetEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, set):
            return list(obj)
        if isinstance(obj, type):
            return obj.__name__
        return json.JSONEncoder.default(self, obj)


def get_trial_substitutions_from_dict(
    parameters: Dict[str, Any],
    experiment_parameters: List[models.V1beta1ParameterSpec],
    trial_parameters: List[models.V1beta1TrialParameterSpec],
) -> Dict[str, str]:
    for p_name, p_value in parameters.items():
        # If input parameter value is Katib Experiment parameter sample.
        if isinstance(p_value, models.V1beta1ParameterSpec):
            # Wrap value for the function input.
            parameters[p_name] = f"${{trialParameters.{p_name}}}"

            # Add value to the Katib Experiment parameters.
            p_value.name = p_name
            experiment_parameters.append(p_value)

            # Add value to the Katib Experiment's Trial parameters.
            trial_parameters.append(
                models.V1beta1TrialParameterSpec(name=p_name, reference=p_name)
            )
        else:
            # Otherwise, add value to the function input.
            parameters[p_name] = p_value

    return parameters


def get_trial_substitutions_from_trainer(
    parameters: Union["TrainingArguments", "LoraConfig"],  # noqa: F821
    experiment_params: List[models.V1beta1ParameterSpec],
    trial_params: List[models.V1beta1TrialParameterSpec],
) -> Dict[str, str]:
    from peft import LoraConfig  # noqa: F401
    from transformers import TrainingArguments  # noqa: F401

    if isinstance(parameters, TrainingArguments):
        parameters_dict = parameters.to_dict()
    else:
        parameters_dict = parameters.__dict__

    for p_name, p_value in parameters_dict.items():
        if not hasattr(parameters, p_name):
            logger.warning(f"Training parameter {p_name} is not supported.")
            continue

        if isinstance(p_value, models.V1beta1ParameterSpec):
            old_attr = getattr(parameters, p_name, None)
            if old_attr is not None:
                value = f"${{trialParameters.{p_name}}}"
            setattr(parameters, p_name, value)
            p_value.name = p_name
            experiment_params.append(p_value)
            trial_params.append(
                models.V1beta1TrialParameterSpec(name=p_name, reference=p_name)
            )
        elif p_value is not None:
            old_attr = getattr(parameters, p_name, None)
            if old_attr is not None:
                if isinstance(p_value, dict):
                    # Update the existing dictionary without nesting
                    value = copy.deepcopy(p_value)
                else:
                    value = type(old_attr)(p_value)
            setattr(parameters, p_name, value)

    if isinstance(parameters, TrainingArguments):
        parameters = json.dumps(parameters.to_dict())
    else:
        parameters = json.dumps(parameters.__dict__, cls=SetEncoder)

    return parameters


def get_exec_script_from_objective(
    objective: Callable,
    entrypoint: str,
    input_params: Dict[str, Any],
    packages_to_install: Optional[List[str]] = None,
    pip_index_url: str = "https://pypi.org/simple",
) -> str:
    """
    Get executable script for container args from the given objective function and parameters.
    """
    # Validate objective function.
    validate_objective_function(objective)

    # Extract objective function implementation.
    objective_code = inspect.getsource(objective)

    # Objective function might be defined in some indented scope
    # (e.g. in another function). We need to dedent the function code.
    objective_code = textwrap.dedent(objective_code)

    # Wrap objective function to execute it from the file. For example:
    # def objective(parameters):
    #     print(f'Parameters are {parameters}')
    # objective({
    #     'lr': '${trialParameters.lr}',
    #     'epochs': '${trialParameters.epochs}',
    #     'is_dist': False
    # })
    objective_code = f"{objective_code}\n{objective.__name__}({input_params})\n"

    # Prepare execute script template.
    exec_script = textwrap.dedent(
        """
                program_path=$(mktemp -d)
                read -r -d '' SCRIPT << EOM\n
                {objective_code}
                EOM
                printf "%s" \"$SCRIPT\" > \"$program_path/ephemeral_script.py\"
                {entrypoint} \"$program_path/ephemeral_script.py\""""
    )

    # Add objective code to the execute script.
    exec_script = exec_script.format(
        objective_code=objective_code, entrypoint=entrypoint
    )

    # Install Python packages if that is required.
    if packages_to_install is not None:
        exec_script = (
            get_script_for_python_packages(packages_to_install, pip_index_url)
            + exec_script
        )

    # Return executable script to execute objective function.
    return exec_script


def get_trial_template_with_job(
    retain_trials: bool,
    trial_parameters: List[models.V1beta1TrialParameterSpec],
    pod_template_spec: client.V1PodTemplateSpec,
) -> models.V1beta1TrialTemplate:
    """
    Get Trial template with Job as a Trial's Worker
    """

    # Restart policy must be set for the Job.
    pod_template_spec.spec.restart_policy = "Never"  # type: ignore

    # Use Job as a Trial spec.
    job = client.V1Job(
        api_version="batch/v1",
        kind="Job",
        spec=client.V1JobSpec(template=pod_template_spec),
    )

    trial_template = models.V1beta1TrialTemplate(
        primary_container_name=constants.DEFAULT_PRIMARY_CONTAINER_NAME,
        retain=retain_trials,
        trial_parameters=trial_parameters,
        trial_spec=job,
    )
    return trial_template


def get_trial_template_with_pytorchjob(
    retain_trials: bool,
    trial_parameters: List[models.V1beta1TrialParameterSpec],
    resources_per_trial: types.TrainerResources,
    master_pod_template_spec: models.V1PodTemplateSpec,
    worker_pod_template_spec: models.V1PodTemplateSpec,
) -> models.V1beta1TrialTemplate:
    """
    Get Trial template with PyTorchJob as a Trial's Worker
    """

    # Use PyTorchJob as a Trial spec.
    pytorchjob = training_models.KubeflowOrgV1PyTorchJob(
        api_version=API_VERSION,
        kind=PYTORCHJOB_KIND,
        spec=training_models.KubeflowOrgV1PyTorchJobSpec(
            run_policy=training_models.KubeflowOrgV1RunPolicy(clean_pod_policy=None),
            nproc_per_node=str(resources_per_trial.num_procs_per_worker),
            pytorch_replica_specs={
                "Master": training_models.KubeflowOrgV1ReplicaSpec(
                    replicas=1,
                    template=master_pod_template_spec,
                )
            },
        ),
    )

    # Add Worker replica if number of workers > 1
    if resources_per_trial.num_workers > 1:
        pytorchjob.spec.pytorch_replica_specs["Worker"] = (
            training_models.KubeflowOrgV1ReplicaSpec(
                replicas=resources_per_trial.num_workers - 1,
                template=worker_pod_template_spec,
            )
        )

    trial_template = models.V1beta1TrialTemplate(
        primary_container_name=JOB_PARAMETERS[PYTORCHJOB_KIND]["container"],
        retain=retain_trials,
        trial_parameters=trial_parameters,
        trial_spec=pytorchjob,
    )
    return trial_template
