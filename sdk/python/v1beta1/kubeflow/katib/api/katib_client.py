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
import multiprocessing
import textwrap
import time
from typing import Any, Callable, Dict, List, Optional, Union
import json
import logging
logger = logging.getLogger(__name__)

import grpc
import kubeflow.katib.katib_api_pb2 as katib_api_pb2
from kubeflow.katib import models
from kubeflow.katib.api_client import ApiClient
from kubeflow.katib.constants import constants
from kubeflow.katib.utils import utils
from kubernetes import client, config


class KatibClient(object):
    def __init__(
        self,
        config_file: str = None,
        context: str = None,
        client_configuration: client.Configuration = None,
        namespace: str = utils.get_default_target_namespace(),
    ):
        """KatibClient constructor.

        Args:
            config_file: Path to the kube-config file. Defaults to ~/.kube/config.
            context: Set the active context. Defaults to current_context from the kube-config.
            client_configuration: Client configuration for cluster authentication.
                You have to provide valid configuration with Bearer token or
                with username and password.
                You can find an example here: https://github.com/kubernetes-client/python/blob/67f9c7a97081b4526470cad53576bc3b71fa6fcc/examples/remote_cluster.py#L31
            namespace: Target Kubernetes namespace. Can be overridden during method invocations.
        """

        self.in_cluster = False
        # If client configuration is not set, use kube-config to access Kubernetes APIs.
        if client_configuration is None:
            # Load kube-config or in-cluster config.
            if config_file or not utils.is_running_in_k8s():
                config.load_kube_config(config_file=config_file, context=context)
            else:
                config.load_incluster_config()
                self.in_cluster = True

        k8s_client = client.ApiClient(client_configuration)
        self.custom_api = client.CustomObjectsApi(k8s_client)
        self.core_api = client.CoreV1Api(k8s_client)
        self.api_client = ApiClient()
        self.namespace = namespace

    def _is_ipython(self):
        """Returns whether we are running in notebook."""

        try:
            import IPython

            ipy = IPython.get_ipython()
            if ipy is None:
                return False
        except ImportError:
            return False
        return True

    def create_experiment(
        self,
        experiment: models.V1beta1Experiment,
        namespace: Optional[str] = None,
    ):
        """Create the Katib Experiment.

        Args:
            experiment: Experiment object of type V1beta1Experiment.
            namespace: Namespace for the Experiment.

        Raises:
            TimeoutError: Timeout to create Katib Experiment.
            RuntimeError: Failed to create Katib Experiment.
        """

        namespace = namespace or self.namespace

        experiment_name = None
        if type(experiment) == models.V1beta1Experiment:
            if experiment.metadata.name is not None:
                experiment_name = experiment.metadata.name
            elif experiment.metadata.generate_name is not None:
                experiment_name = experiment.metadata.generate_name
        elif "name" in experiment["metadata"]:
            experiment_name = experiment["metadata"]["name"]
        elif "generate_name" in experiment["metadata"]:
            experiment_name = experiment["metadata"]["generate_name"]

        if experiment_name is None:
            raise ValueError("Experiment must have a name or generateName")

        try:
            outputs = self.custom_api.create_namespaced_custom_object(
                constants.KUBEFLOW_GROUP,
                constants.KATIB_VERSION,
                namespace,
                constants.EXPERIMENT_PLURAL,
                experiment,
            )
            experiment_name = outputs["metadata"][
                "name"
            ]  # if "generate_name" is used, "name" gets a prefix from server
        except multiprocessing.TimeoutError:
            raise TimeoutError(
                f"Timeout to create Katib Experiment: {namespace}/{experiment_name}"
            )
        except Exception as e:
            if hasattr(e, "status") and e.status == 409:
                raise Exception(
                    f"A Katib Experiment with the name {namespace}/{experiment_name} already exists."
                )
            raise RuntimeError(
                f"Failed to create Katib Experiment: {namespace}/{experiment_name}"
            )

        # TODO (andreyvelich): Use proper logger.
        print(f"Experiment {namespace}/{experiment_name} has been created")

        if self._is_ipython():
            if self.in_cluster:
                import IPython

                IPython.display.display(
                    IPython.display.HTML(
                        "Katib Experiment {} "
                        'link <a href="/_/katib/#/katib/hp_monitor/{}/{}" target="_blank">here</a>'.format(
                            experiment_name,
                            namespace,
                            experiment_name,
                        )
                    )
                )

    def tune(
        self,
        # TODO (andreyvelich): How to be consistent with other APIs (name) ?
        name: str,
        model_provider_parameters: Optional[Any] = None,
		dataset_provider_parameters: Optional[Any] = None,
		storage_config: Optional[Dict[str, Optional[Union[str, List[str]]]]] = {
            "size": constants.PVC_DEFAULT_SIZE,
            "storage_class": None,
            "access_modes": constants.PVC_DEFAULT_ACCESS_MODES,
        },
        objective: Optional[Callable] = None,
        base_image: Optional[str] = constants.BASE_IMAGE_TENSORFLOW,
        trainer_parameters = None,
        namespace: Optional[str] = None,
        env_per_trial: Optional[
            Union[Dict[str, str], List[Union[client.V1EnvVar, client.V1EnvFromSource]]]
        ] = None,
        algorithm_name: str = "random",
        algorithm_settings: Union[
            dict, List[models.V1beta1AlgorithmSetting], None
        ] = None,
        objective_metric_name: str = None,
        additional_metric_names: List[str] = [],
        objective_type: str = "maximize",
        objective_goal: float = None,
        max_trial_count: int = None,
        parallel_trial_count: int = None,
        max_failed_trial_count: int = None,
        resources_per_trial: Union[dict, client.V1ResourceRequirements, None] = None,
        retain_trials: bool = False,
        packages_to_install: List[str] = None,
        pip_index_url: str = "https://pypi.org/simple",
    ):
        """Create HyperParameter Tuning Katib Experiment using one of the following options:
        - External models and datasets: Specify both `model_provider_parameters` and `dataset_provider_parameters` to download models and datasets from external platforms (currently supports HuggingFace and Amazon S3) using the Storage Initializer. The `trainer_parameters` should be of type `HuggingFaceTrainerParams` to set the hyperparameters search space. This API will automatically define the "Trainer" class in HuggingFace with the provided parameters and utilize `Trainer.train()` from HuggingFace to obtain the metrics for optimizing hyperparameters. 
        - Custom objective function: Specify the `objective` parameter to define your own objective function. The `base_image` parameter will be used to execute the objective function. `trainer_parameters` should be a dictionary to define the search space for these parameters.

        Args:
            name: Name for the Experiment.
            model_provider_parameters: Parameters for the model provider in the Storage Initializer.
                For example, HuggingFace model name and Transformer type for that model, like: AutoModelForSequenceClassification. This argument must be the type of `kubeflow.storage_initializer.hugging_face.HuggingFaceModelParams`.
            dataset_provider_parameters: Parameters for the dataset provider in the Storage Initializer.
                For example, name of the HuggingFace dataset or AWS S3 configuration. This argument must be the type of `kubeflow.storage_initializer.hugging_face.HuggingFaceDatasetParams` or `kubeflow.storage_initializer.s3.S3DatasetParams`
            storage_config: Configuration for Storage Initializer PVC to download pre-trained model and dataset.
                You can configure PVC size and storage class name in this argument.
            objective: Objective function that Katib uses to train the model.
                This function must be Callable and it must have only one dict argument.
                Katib uses this argument to send HyperParameters to the function.
                The function should not use any code declared outside of the function
                definition. Import statements must be added inside the function.
            base_image: Image to use when executing the objective function.
            trainer_parameters: Parameters for configuring the training process, including settings for the hyperparameters search space.
                You should use the Katib SDK to define the search space for these parameters.
                If you choose to use external models and datasets, it should be of type `HuggingFaceTrainerParams`. For example:
                ```
                trainer_parameters = HuggingFaceTrainerParams(
                    training_parameters = transformers.TrainingArguments(
                        learning_rate = katib.search.double(min=0.1, max=0.2),
                    ),
                ),
                ```
                If you choose a custom objective function, it should be a dictionary. For example:
                ```
                trainer_parameters = {"lr": katib.search.double(min=0.1, max=0.2)}
                ```
                Also, you can use these parameters to define input for training the external models or your custom objective function.
            namespace: Namespace for the Experiment.
            env_per_trial: Environment variable(s) to be attached to each trial container.
                You can specify a dictionary as a mapping object representing the environment
                variables. Otherwise, you can specify a list, in which the element can either
                be a kubernetes.client.models.V1EnvVar (documented here:
                https://github.com/kubernetes-client/python/blob/master/kubernetes/docs/V1EnvVar.md)
                or a kubernetes.client.models.V1EnvFromSource (documented here:
                https://github.com/kubernetes-client/python/blob/master/kubernetes/docs/V1EnvFromSource.md)
            algorithm_name: Search algorithm for the HyperParameter tuning.
            algorithm_settings: Settings for the search algorithm given.
                For available fields, check this doc: https://www.kubeflow.org/docs/components/katib/experiment/#search-algorithms-in-detail.
            objective_metric_name: Objective metric that Katib optimizes.
            additional_metric_names: List of metrics that Katib collects from the
                objective function in addition to objective metric.
            objective_type: Type for the Experiment optimization for the objective metric.
                Must be one of `minimize` or `maximize`.
            objective_goal: Objective goal that Experiment should reach to be Succeeded.
            max_trial_count: Maximum number of Trials to run. For the default
                values check this doc: https://www.kubeflow.org/docs/components/katib/experiment/#configuration-spec.
            parallel_trial_count: Number of Trials that Experiment runs in parallel.
            max_failed_trial_count: Maximum number of Trials allowed to fail.
            resources_per_trial: A parameter that lets you specify how much
            resources each trial container should have. You can either specify a
            kubernetes.client.V1ResourceRequirements object (documented here:
            https://github.com/kubernetes-client/python/blob/master/kubernetes/docs/V1ResourceRequirements.md)
            or a dictionary that includes one or more of the following keys:
            `cpu`, `memory`, or `gpu` (other keys will be ignored). Appropriate
            values for these keys are documented here:
            https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/.
            For example:
                {
                    "cpu": "1",
                    "gpu": "1",
                    "memory": "2Gi",
                }
            Please note, `gpu` specifies a resource request with a key of
            `nvidia.com/gpu`, i.e. an NVIDIA GPU. If you need a different type
            of GPU, pass in a V1ResourceRequirement instance instead, since it's
            more flexible. This parameter is optional and defaults to None.
            retain_trials: Whether Trials' resources (e.g. pods) are deleted after Succeeded state.
            packages_to_install: List of Python packages to install in addition
                to the base image packages. These packages are installed before
                executing the objective function.
            pip_index_url: The PyPI url from which to install Python packages.

        Raises:
            ValueError: Function arguments have incorrect type or value.
            TimeoutError: Timeout to create Katib Experiment.
            RuntimeError: Failed to create Katib Experiment.
        """

        print(
            "Thank you for using `tune` API for LLMs hyperparameters optimization. This feature is in alpha stage Kubeflow community is looking for your feedback. Please share your experience via #kubeflow-katib Slack channel or Kubeflow Katib GitHub."
        )
        
        if (
            ((model_provider_parameters is not None) and (dataset_provider_parameters is not None)) == (objective is not None)
        ):
            raise ValueError(
                "Invalid configuration for creating a Katib Experiment for hyperparameter optimization. "
                "You should only specify one of the following options: 1) `model_provider_parameters` and `dataset_provider_parameters`; 2) `objective`."
            )
        
        if (
            not name
            or not trainer_parameters
        ):
            raise ValueError("One of the required parameters is None")

        namespace = namespace or self.namespace

        # Create Katib Experiment template.
        experiment = models.V1beta1Experiment(
            api_version=f"{constants.KUBEFLOW_GROUP}/{constants.KATIB_VERSION}",
            kind=constants.EXPERIMENT_KIND,
            metadata=models.V1ObjectMeta(name=name, namespace=namespace),
            spec=models.V1beta1ExperimentSpec(),
        )

        # Add Objective to the Katib Experiment.
        experiment.spec.objective = models.V1beta1ObjectiveSpec(
            type=objective_type,
            objective_metric_name=objective_metric_name,
            additional_metric_names=additional_metric_names,
        )
        if objective_goal is not None:
            experiment.spec.objective.goal = objective_goal

        # Add Algorithm to the Katib Experiment.
        if isinstance(algorithm_settings, dict):
            algorithm_settings = [
                models.V1beta1AlgorithmSetting(name=str(k), value=str(v))
                for k, v in algorithm_settings.items()
            ]

        experiment.spec.algorithm = models.V1beta1AlgorithmSpec(
            algorithm_name=algorithm_name,
            algorithm_settings=algorithm_settings,
        )

        # Add Trial budget to the Katib Experiment.
        if max_trial_count is not None:
            experiment.spec.max_trial_count = max_trial_count
        if parallel_trial_count is not None:
            experiment.spec.parallel_trial_count = parallel_trial_count
        if max_failed_trial_count is not None:
            experiment.spec.max_failed_trial_count = max_failed_trial_count
        
        # Add resources to the Katib Experiment.
        if isinstance(resources_per_trial, dict):
            if "gpu" in resources_per_trial:
                resources_per_trial["nvidia.com/gpu"] = resources_per_trial.pop("gpu")

            resources_per_trial = client.V1ResourceRequirements(
                requests=resources_per_trial,
                limits=resources_per_trial,
            )

        # Add environment variables to the Katib Experiment.
        env = []
        env_from = []
        if isinstance(env_per_trial, dict):
            env = [
                client.V1EnvVar(name=str(k), value=str(v))
                for k, v in env_per_trial.items()
            ]
        elif env_per_trial:
            for x in env_per_trial:
                if isinstance(x, client.V1EnvVar):
                    env.append(x)
                elif isinstance(x, client.V1EnvFromSource):
                    env_from.append(x)
                else:
                    raise ValueError(
                        f"Incorrect value for env_per_trial: {env_per_trial}"
                    )

        # Create Container and Pod specifications.
        # If users choose to use a custom objective function.
        if objective is not None:
            # Validate objective function.
            utils.validate_objective_function(objective)

            # Extract objective function implementation.
            objective_code = inspect.getsource(objective)

            # Objective function might be defined in some indented scope
            # (e.g. in another function). We need to dedent the function code.
            objective_code = textwrap.dedent(objective_code)

            # Iterate over input parameters.
            input_params = {}
            experiment_params = []
            trial_params = []
            for p_name, p_value in trainer_parameters.items():
                # If input parameter value is Katib Experiment parameter sample.
                if isinstance(p_value, models.V1beta1ParameterSpec):
                    # Wrap value for the function input.
                    input_params[p_name] = f"${{trialParameters.{p_name}}}"

                    # Add value to the Katib Experiment parameters.
                    p_value.name = p_name
                    experiment_params.append(p_value)

                    # Add value to the Katib Experiment's Trial parameters.
                    trial_params.append(
                        models.V1beta1TrialParameterSpec(name=p_name, reference=p_name)
                    )
                else:
                    # Otherwise, add value to the function input.
                    input_params[p_name] = p_value

            # Wrap objective function to execute it from the file. For example
            # def objective(parameters):
            #     print(f'Parameters are {parameters}')
            # objective({'lr': '${trialParameters.lr}', 'epochs': '${trialParameters.epochs}', 'is_dist': False})
            objective_code = f"{objective_code}\n{objective.__name__}({input_params})\n"

            # Prepare execute script template.
            exec_script = textwrap.dedent(
                """
                program_path=$(mktemp -d)
                read -r -d '' SCRIPT << EOM\n
                {objective_code}
                EOM
                printf "%s" "$SCRIPT" > $program_path/ephemeral_objective.py
                python3 -u $program_path/ephemeral_objective.py"""
            )

            # Add objective code to the execute script.
            exec_script = exec_script.format(objective_code=objective_code)

            # Install Python packages if that is required.
            if packages_to_install is not None:
                exec_script = (
                    utils.get_script_for_python_packages(packages_to_install, pip_index_url)
                    + exec_script
                )
            
            # create app container spec
            container_spec = client.V1Container(
                                name=constants.DEFAULT_PRIMARY_CONTAINER_NAME,
                                image=base_image,
                                command=["bash", "-c"],
                                args=[exec_script],
                                env=env if env else None,
                                env_from=env_from if env_from else None,
                                resources=resources_per_trial,
                            )
            
            pod_spec = client.V1PodTemplateSpec(
                metadata=models.V1ObjectMeta(
                    annotations={"sidecar.istio.io/inject": "false"}
                ),
                spec=client.V1PodSpec(
                    restart_policy="Never",
                    containers=[container_spec],
                ),
            )

        # If users choose to use external models and datasets.
        else:
            try:
                import peft
                import transformers
                from kubeflow.storage_initializer.s3 import S3DatasetParams
                from kubeflow.storage_initializer.hugging_face import (
                    HuggingFaceModelParams,
                    HuggingFaceDatasetParams,
                    HuggingFaceTrainerParams,
                )
            except ImportError:
                raise ImportError(
                    "Tune API dependencies not installed. "
                    + "Run: pip install -U 'kubeflow-training[huggingface]' "
                )
            
            # Create PVC for the Storage Initializer.
            try:
                self.core_api.create_namespaced_persistent_volume_claim(
                    namespace=namespace,
                    body=utils.get_pvc_spec(
                        pvc_name=constants.STORAGE_INITIALIZER,
                        namespace=namespace,
                        storage_config=storage_config,
                    ),
                )
            except Exception as e:
                pvc_list = self.core_api.list_namespaced_persistent_volume_claim(namespace)
                # Check if the PVC with the specified name exists.
                for pvc in pvc_list.items:
                    if pvc.metadata.name == constants.STORAGE_INITIALIZER:
                        print(
                            f"PVC '{constants.STORAGE_INITIALIZER}' already exists in namespace "
                            f"{namespace}."
                        )
                        break
                else:
                    raise RuntimeError(f"failed to create PVC. Error: {e}")
            
            if isinstance(model_provider_parameters, HuggingFaceModelParams):
                mp = "hf"
            else:
                raise ValueError("Model provider parameters must be an instance of HuggingFaceModelParams.")
            
            if isinstance(dataset_provider_parameters, S3DatasetParams):
                dp = "s3"
            elif isinstance(dataset_provider_parameters, HuggingFaceDatasetParams):
                dp = "hf"
            else:
                raise ValueError("Dataset provider parameters must be an instance of S3DatasetParams or HuggingFaceDatasetParams.")
            
            # Iterate over input parameters.
            experiment_params = []
            trial_params = []

            training_args = trainer_parameters.training_parameters
            for p_name, p_value in trainer_parameters.training_parameters.to_dict().items():
                if not hasattr(training_args, p_name):
                    logger.warning(f"Training parameter {p_name} is not supported by the current transformer.")
                    continue
                if isinstance(p_value, models.V1beta1ParameterSpec):
                    old_attr = getattr(training_args, p_name, None)
                    if old_attr is not None:
                        value = f"${{trialParameters.{p_name}}}"
                    setattr(training_args, p_name, value)
                    p_value.name = p_name
                    experiment_params.append(p_value)
                    trial_params.append(models.V1beta1TrialParameterSpec(name=p_name, reference=p_name))
                elif p_value is not None:
                    old_attr = getattr(training_args, p_name, None)
                    if old_attr is not None:
                        value = type(old_attr)(p_value)
                    setattr(training_args, p_name, value)

            lora_config = trainer_parameters.lora_config
            for p_name, p_value in trainer_parameters.lora_config.__dict__.items():
                if not hasattr(lora_config, p_name):
                    logger.warning(f"Training parameter {p_name} is not supported by the current peft.")
                    continue
                if isinstance(p_value, models.V1beta1ParameterSpec):
                    old_attr = getattr(lora_config, p_name, None)
                    if old_attr is not None:
                        value = f"${{trialParameters.{p_name}}}"
                    setattr(lora_config, p_name, value)
                    p_value.name = p_name
                    experiment_params.append(p_value)
                    trial_params.append(models.V1beta1TrialParameterSpec(name=p_name, reference=p_name))
                elif p_value is not None:
                    old_attr = getattr(lora_config, p_name, None)
                    if old_attr is not None:
                        value = type(old_attr)(p_value)
                    setattr(lora_config, p_name, value)

            # create init container spec.
            init_container_spec = client.V1Container(
                name=constants.STORAGE_INITIALIZER,
                image=constants.STORAGE_INITIALIZER_IMAGE,
                args=[
                    "--model_provider",
                    mp,
                    "--model_provider_parameters",
                    json.dumps(model_provider_parameters.__dict__, cls=utils.SetEncoder), 
                    "--dataset_provider",
                    dp,
                    "--dataset_provider_parameters",
                    json.dumps(dataset_provider_parameters.__dict__),
                ],
                volume_mounts=[constants.STORAGE_INITIALIZER_VOLUME_MOUNT],
            )

            from kubeflow.storage_initializer.constants import (
                VOLUME_PATH_DATASET,
                VOLUME_PATH_MODEL,
            )   

            lora_config = json.dumps(lora_config.__dict__, cls=utils.SetEncoder)
            training_args = json.dumps(training_args.to_dict())
            # create app container spec.
            container_spec = client.V1Container(
                name=constants.DEFAULT_PRIMARY_CONTAINER_NAME,
                image=constants.TRAINER_TRANSFORMER_IMAGE,
                args=[
                    "--model_uri",
                    model_provider_parameters.model_uri,
                    "--transformer_type",
                    model_provider_parameters.transformer_type.__name__,
                    "--model_dir",
                    VOLUME_PATH_MODEL,
                    "--dataset_dir",
                    VOLUME_PATH_DATASET,
                    "--lora_config",
                    f"'{lora_config}'",
                    "--training_parameters",
                    f"'{training_args}'",
                ],
                volume_mounts=[constants.STORAGE_INITIALIZER_VOLUME_MOUNT],
                env=env if env else None,
                env_from=env_from if env_from else None,
                resources=resources_per_trial,
            )

            pod_spec = client.V1PodTemplateSpec(
                metadata=models.V1ObjectMeta(
                    annotations={"sidecar.istio.io/inject": "false"}
                ),
                spec=client.V1PodSpec(
                    restart_policy="Never",
                    containers=[container_spec],
                    init_containers=[init_container_spec],
                    volumes=[constants.STORAGE_INITIALIZER_VOLUME],
                ),
            )
        
        # Create Trial specification.
        trial_spec = client.V1Job(
            api_version="batch/v1",
            kind="Job",
            spec=client.V1JobSpec(
                template=pod_spec,
            ),
        )

        # Create Trial template.
        trial_template = models.V1beta1TrialTemplate(
            primary_container_name=constants.DEFAULT_PRIMARY_CONTAINER_NAME,
            retain=retain_trials,
            trial_parameters=trial_params,
            trial_spec=trial_spec,
        )

        # Add parameters to the Katib Experiment.
        experiment.spec.parameters = experiment_params

        # Add Trial template to the Katib Experiment.
        experiment.spec.trial_template = trial_template

        # Create the Katib Experiment.
        self.create_experiment(experiment, namespace)

    def get_experiment(
        self,
        name: str,
        namespace: Optional[str] = None,
        timeout: int = constants.DEFAULT_TIMEOUT,
    ):
        """Get the Katib Experiment.

        Args:
            name: Name for the Experiment.
            namespace: Namespace for the Experiment.
            timeout: Optional, Kubernetes API server timeout in seconds
                to execute the request.

        Returns:
            V1beta1Experiment: Katib Experiment object.

        Raises:
            TimeoutError: Timeout to get Katib Experiment.
            RuntimeError: Failed to get Katib Experiment.
        """

        namespace = namespace or self.namespace

        try:
            thread = self.custom_api.get_namespaced_custom_object(
                constants.KUBEFLOW_GROUP,
                constants.KATIB_VERSION,
                namespace,
                constants.EXPERIMENT_PLURAL,
                name,
                async_req=True,
            )
            response = utils.FakeResponse(thread.get(timeout))
            experiment = self.api_client.deserialize(response, models.V1beta1Experiment)
            return experiment

        except multiprocessing.TimeoutError:
            raise TimeoutError(f"Timeout to get Katib Experiment: {namespace}/{name}")
        except Exception:
            raise RuntimeError(f"Failed to get Katib Experiment: {namespace}/{name}")

    def list_experiments(
        self,
        namespace: Optional[str] = None,
        timeout: int = constants.DEFAULT_TIMEOUT,
    ):
        """List of all Katib Experiments in namespace.

        Args:
            namespace: Namespace to list the Experiments.
            timeout: Optional, Kubernetes API server timeout in seconds
                to execute the request.

        Returns:
            list[V1beta1Experiment]: List of Katib Experiment objects. It returns
            empty list if Experiments cannot be found.

        Raises:
            TimeoutError: Timeout to list Katib Experiments.
            RuntimeError: Failed to list Katib Experiments.
        """

        namespace = namespace or self.namespace

        result = []
        try:
            thread = self.custom_api.list_namespaced_custom_object(
                constants.KUBEFLOW_GROUP,
                constants.KATIB_VERSION,
                namespace=namespace,
                plural=constants.EXPERIMENT_PLURAL,
                async_req=True,
            )
            response = thread.get(timeout)
            result = [
                self.api_client.deserialize(
                    utils.FakeResponse(item), models.V1beta1Experiment
                )
                for item in response.get("items")
            ]
        except multiprocessing.TimeoutError:
            raise TimeoutError(
                f"Timeout to list Katib Experiments in namespace: {namespace}"
            )
        except Exception:
            raise RuntimeError(
                f"Failed to list Katib Experiments in namespace: {namespace}"
            )
        return result

    def get_experiment_conditions(
        self,
        name: str,
        namespace: Optional[str] = None,
        experiment: models.V1beta1Experiment = None,
        timeout: int = constants.DEFAULT_TIMEOUT,
    ):
        """Get the Experiment conditions. Experiment is in the condition when
        `status` is True for the appropriate condition `type`.

        Args:
            name: Name for the Experiment.
            namespace: Namespace for the Experiment.
            experiment: Optionally, Experiment object can be set to get the conditions.
            timeout: Optional, Kubernetes API server timeout in seconds
                to execute the request.

        Returns:
            list[V1beta1ExperimentCondition]: List of Experiment conditions with
                last transition time, last update time, message, reason, type, and
                status. It returns empty list if Experiment does not have any
                conditions yet.

        Raises:
            TimeoutError: Timeout to get Katib Experiment.
            RuntimeError: Failed to get Katib Experiment.
        """

        namespace = namespace or self.namespace

        if experiment is None:
            experiment = self.get_experiment(name, namespace, timeout)

        if (
            experiment.status
            and experiment.status.conditions
            and len(experiment.status.conditions) > 0
        ):
            return experiment.status.conditions

        return []

    def is_experiment_created(
        self,
        name: str,
        namespace: Optional[str] = None,
        experiment: models.V1beta1Experiment = None,
        timeout: int = constants.DEFAULT_TIMEOUT,
    ):
        """Check if Experiment is Created.

        Args:
            name: Name for the Experiment.
            namespace: Namespace for the Experiment.
            experiment: Optionally, Experiment object can be set to check the status.
            timeout: Optional, Kubernetes API server timeout in seconds
                to execute the request.

        Returns:
            bool: True is Experiment is Created, else False.

        Raises:
            TimeoutError: Timeout to get Katib Experiment.
            RuntimeError: Failed to get Katib Experiment.
        """

        namespace = namespace or self.namespace

        return utils.has_condition(
            self.get_experiment_conditions(name, namespace, experiment, timeout),
            constants.EXPERIMENT_CONDITION_CREATED,
        )

    def is_experiment_running(
        self,
        name: str,
        namespace: Optional[str] = None,
        experiment: models.V1beta1Experiment = None,
        timeout: int = constants.DEFAULT_TIMEOUT,
    ):
        """Check if Experiment is Running.

        Args:
            name: Name for the Experiment.
            namespace: Namespace for the Experiment.
            experiment: Optionally, Experiment object can be set to check the status.
            timeout: Optional, Kubernetes API server timeout in seconds
                to execute the request.

        Returns:
            bool: True is Experiment is Running, else False.

        Raises:
            TimeoutError: Timeout to get Katib Experiment.
            RuntimeError: Failed to get Katib Experiment.
        """

        namespace = namespace or self.namespace

        return utils.has_condition(
            self.get_experiment_conditions(name, namespace, experiment, timeout),
            constants.EXPERIMENT_CONDITION_RUNNING,
        )

    def is_experiment_restarting(
        self,
        name: str,
        namespace: Optional[str] = None,
        experiment: models.V1beta1Experiment = None,
        timeout: int = constants.DEFAULT_TIMEOUT,
    ):
        """Check if Experiment is Restarting.
        Args:
            name: Name for the Experiment.
            namespace: Namespace for the Experiment.
            experiment: Optionally, Experiment object can be set to check the status.
            timeout: Optional, Kubernetes API server timeout in seconds
                to execute the request.

        Returns:
            bool: True is Experiment is Resting, else False.

        Raises:
            TimeoutError: Timeout to get Katib Experiment.
            RuntimeError: Failed to get Katib Experiment.
        """

        namespace = namespace or self.namespace

        return utils.has_condition(
            self.get_experiment_conditions(name, namespace, experiment, timeout),
            constants.EXPERIMENT_CONDITION_RESTARTING,
        )

    def is_experiment_succeeded(
        self,
        name: str,
        namespace: Optional[str] = None,
        experiment: models.V1beta1Experiment = None,
        timeout: int = constants.DEFAULT_TIMEOUT,
    ):
        """Check if Experiment is Succeeded.
        Args:
            name: Name for the Experiment.
            namespace: Namespace for the Experiment.
            experiment: Optionally, Experiment object can be set to check the status.
            timeout: Optional, Kubernetes API server timeout in seconds
                to execute the request.

        Returns:
            bool: True is Experiment is Succeeded, else False.

        Raises:
            TimeoutError: Timeout to get Katib Experiment.
            RuntimeError: Failed to get Katib Experiment.
        """

        namespace = namespace or self.namespace

        return utils.has_condition(
            self.get_experiment_conditions(name, namespace, experiment, timeout),
            constants.EXPERIMENT_CONDITION_SUCCEEDED,
        )

    def is_experiment_failed(
        self,
        name: str,
        namespace: Optional[str] = None,
        experiment: models.V1beta1Experiment = None,
        timeout: int = constants.DEFAULT_TIMEOUT,
    ):
        """Check if Experiment is Failed.
        Args:
            name: Name for the Experiment.
            namespace: Namespace for the Experiment.
            experiment: Optionally, Experiment object can be set to check the status.
            timeout: Optional, Kubernetes API server timeout in seconds
                to execute the request.

        Returns:
            bool: True is Experiment is Failed, else False.

        Raises:
            TimeoutError: Timeout to get Katib Experiment.
            RuntimeError: Failed to get Katib Experiment.
        """

        namespace = namespace or self.namespace

        return utils.has_condition(
            self.get_experiment_conditions(name, namespace, experiment, timeout),
            constants.EXPERIMENT_CONDITION_FAILED,
        )

    def wait_for_experiment_condition(
        self,
        name: str,
        namespace: Optional[str] = None,
        expected_condition: str = constants.EXPERIMENT_CONDITION_SUCCEEDED,
        timeout: int = 600,
        polling_interval: int = 15,
        apiserver_timeout: int = constants.DEFAULT_TIMEOUT,
    ):
        """Wait until Experiment reaches specific condition. By default it waits
        for the Succeeded condition.

        Args:
            name: Name for the Experiment.
            namespace: Namespace for the Experiment.
            expected_condition: Which condition Experiment should reach.
            timeout: How many seconds to wait until Experiment reaches condition.
            polling_interval: The polling interval in seconds to get Experiment status.
            apiserver_timeout: Optional, Kubernetes API server timeout in seconds
                to execute the request.

        Returns:
            V1beta1Experiment: Katib Experiment object.

        Raises:
            RuntimeError: Failed to get Katib Experiment or Experiment reaches
                Failed state if it does not wait for this condition.
            TimeoutError: Timeout waiting for Experiment to reach required condition
                or timeout to get Katib Experiment.
        """

        namespace = namespace or self.namespace

        for _ in range(round(timeout / polling_interval)):
            # We should get Experiment only once per cycle and check the statuses.
            experiment = self.get_experiment(name, namespace, apiserver_timeout)

            # Wait for Failed condition.
            if (
                expected_condition == constants.EXPERIMENT_CONDITION_FAILED
                and self.is_experiment_failed(
                    name, namespace, experiment, apiserver_timeout
                )
            ):
                utils.print_experiment_status(experiment)
                print(f"Experiment: {namespace}/{name} is {expected_condition}\n\n\n")
                return experiment

            # Raise exception if Experiment is Failed.
            elif self.is_experiment_failed(
                name, namespace, experiment, apiserver_timeout
            ):
                raise RuntimeError(
                    f"Experiment: {namespace}/{name} is Failed. "
                    f"Experiment conditions: {experiment.status.conditions}"
                )

            # Check if Experiment reaches Created condition.
            elif (
                expected_condition == constants.EXPERIMENT_CONDITION_CREATED
                and self.is_experiment_created(
                    name, namespace, experiment, apiserver_timeout
                )
            ):
                utils.print_experiment_status(experiment)
                print(f"Experiment: {namespace}/{name} is {expected_condition}\n\n\n")
                return experiment

            # Check if Experiment reaches Running condition.
            elif (
                expected_condition == constants.EXPERIMENT_CONDITION_RUNNING
                and self.is_experiment_running(
                    name, namespace, experiment, apiserver_timeout
                )
            ):
                utils.print_experiment_status(experiment)
                print(f"Experiment: {namespace}/{name} is {expected_condition}\n\n\n")
                return experiment

            # Check if Experiment reaches Restarting condition.
            elif (
                expected_condition == constants.EXPERIMENT_CONDITION_RESTARTING
                and self.is_experiment_restarting(
                    name, namespace, experiment, apiserver_timeout
                )
            ):
                utils.print_experiment_status(experiment)
                print(f"Experiment: {namespace}/{name} is {expected_condition}\n\n\n")
                return experiment

            # Check if Experiment reaches Succeeded condition.
            elif (
                expected_condition == constants.EXPERIMENT_CONDITION_SUCCEEDED
                and self.is_experiment_succeeded(
                    name, namespace, experiment, apiserver_timeout
                )
            ):
                utils.print_experiment_status(experiment)
                print(f"Experiment: {namespace}/{name} is {expected_condition}\n\n\n")
                return experiment

            # Otherwise, print the current Experiment results and sleep for the pooling interval.
            utils.print_experiment_status(experiment)
            print(
                f"Waiting for Experiment: {namespace}/{name} to reach {expected_condition} condition\n\n\n"
            )
            time.sleep(polling_interval)

        raise TimeoutError(
            f"Timeout waiting for Experiment: {namespace}/{name} to reach {expected_condition} state"
        )

    def edit_experiment_budget(
        self,
        name: str,
        namespace: Optional[str] = None,
        max_trial_count: int = None,
        parallel_trial_count: int = None,
        max_failed_trial_count: int = None,
        timeout: int = constants.DEFAULT_TIMEOUT,
    ):
        """Update Experiment budget for the running Trials. You can modify Trial
        budget to resume Succeeded Experiments with `LongRunning` and `FromVolume`
        resume policies.

        Learn about resuming Experiments here: https://www.kubeflow.org/docs/components/katib/resume-experiment/

        Args:
            name: Name for the Experiment.
            namespace: Namespace for the Experiment.
            max_trial_count: The new maximum number of Trials.
            parallel_trial_count: The new number of Trials that Experiment runs in parallel.
            max_failed_trial_count: The new maximum number of Trials allowed to fail.
            timeout: Optional, Kubernetes API server timeout in seconds
                to execute the request.

        Raises:
            ValueError: The new Trial budget is not set.
            TimeoutError: Timeout to edit/get Katib Experiment or timeout to wait
                until Experiment reaches Restarting condition.
            RuntimeError: Failed to edit/get Katib Experiment or Experiment
                reaches Failed condition.
        """

        namespace = namespace or self.namespace

        # The new Trial budget must be set.
        if (
            max_trial_count is None
            and parallel_trial_count is None
            and max_failed_trial_count is None
        ):
            raise ValueError(
                "Invalid input arguments. "
                "You have to set max_trial_count, parallel_trial_count, or max_failed_trial_count "
                "to modify Experiment Trial budget."
            )

        # Modify the Experiment Trial budget.
        experiment = self.get_experiment(name, namespace, timeout)
        if max_trial_count is not None:
            experiment.spec.max_trial_count = max_trial_count
        if parallel_trial_count is not None:
            experiment.spec.parallel_trial_count = parallel_trial_count
        if max_failed_trial_count is not None:
            experiment.spec.max_failed_trial_count = max_failed_trial_count

        # Update Experiment with the new Trial budget.
        try:
            self.custom_api.patch_namespaced_custom_object(
                constants.KUBEFLOW_GROUP,
                constants.KATIB_VERSION,
                namespace,
                constants.EXPERIMENT_PLURAL,
                name,
                experiment,
            )
        except multiprocessing.TimeoutError:
            raise TimeoutError(f"Timeout to edit Katib Experiment: {namespace}/{name}")
        except Exception:
            raise RuntimeError(f"Failed to edit Katib Experiment: {namespace}/{name}")

        print(f"Experiment {namespace}/{name} has been updated")

    def delete_experiment(
        self,
        name: str,
        namespace: Optional[str] = None,
        delete_options: client.V1DeleteOptions = None,
    ):
        """Delete the Katib Experiment.

        Args:
            name: Name for the Experiment.
            namespace: Namespace for the Experiment.
            delete_options: Optional, V1DeleteOptions to set while deleting
                Katib Experiment. For example, grace period seconds.

        Raises:
            TimeoutError: Timeout to delete Katib Experiment.
            RuntimeError: Failed to delete Katib Experiment.
        """

        namespace = namespace or self.namespace

        try:
            self.custom_api.delete_namespaced_custom_object(
                constants.KUBEFLOW_GROUP,
                constants.KATIB_VERSION,
                namespace,
                constants.EXPERIMENT_PLURAL,
                name,
                body=delete_options,
            )
        except multiprocessing.TimeoutError:
            raise TimeoutError(
                f"Timeout to delete Katib Experiment: {namespace}/{name}"
            )
        except Exception:
            raise RuntimeError(f"Failed to delete Katib Experiment: {namespace}/{name}")

        # TODO (andreyvelich): Use proper logger.
        print(f"Experiment {namespace}/{name} has been deleted")

    def get_suggestion(
        self,
        name: str,
        namespace: Optional[str] = None,
        timeout: int = constants.DEFAULT_TIMEOUT,
    ):
        """Get the Katib Suggestion.

        Args:
            name: Name for the Suggestion.
            namespace: Namespace for the Suggestion.
            timeout: Optional, Kubernetes API server timeout in seconds
                to execute the request.

        Returns:
            V1beta1Suggestion: Katib Suggestion object.

        Raises:
            TimeoutError: Timeout to get Katib Suggestion.
            RuntimeError: Failed to get Katib Suggestion.
        """

        namespace = namespace or self.namespace

        try:
            thread = self.custom_api.get_namespaced_custom_object(
                constants.KUBEFLOW_GROUP,
                constants.KATIB_VERSION,
                namespace,
                constants.SUGGESTION_PLURAL,
                name,
                async_req=True,
            )
            response = utils.FakeResponse(thread.get(timeout))
            suggestion = self.api_client.deserialize(response, models.V1beta1Suggestion)
            return suggestion

        except multiprocessing.TimeoutError:
            raise TimeoutError(f"Timeout to get Katib Suggestion: {namespace}/{name}")
        except Exception:
            raise RuntimeError(f"Failed to get Katib Suggestion: {namespace}/{name}")

    def list_suggestions(
        self,
        namespace: Optional[str] = None,
        timeout: int = constants.DEFAULT_TIMEOUT,
    ):
        """List of all Katib Suggestion in namespace.

        Args:
            namespace: Namespace to list the Suggestions.
            timeout: Optional, Kubernetes API server timeout in seconds
                to execute the request.

        Returns:
            list[V1beta1Suggestion]: List of Katib Suggestions objects. It returns
            empty list if Suggestions cannot be found.

        Raises:
            TimeoutError: Timeout to list Katib Suggestions.
            RuntimeError: Failed to list Katib Suggestions.
        """

        namespace = namespace or self.namespace

        result = []
        try:
            thread = self.custom_api.list_namespaced_custom_object(
                constants.KUBEFLOW_GROUP,
                constants.KATIB_VERSION,
                namespace=namespace,
                plural=constants.EXPERIMENT_PLURAL,
                async_req=True,
            )
            response = thread.get(timeout)
            result = [
                self.api_client.deserialize(
                    utils.FakeResponse(item), models.V1beta1Suggestion
                )
                for item in response.get("items")
            ]
        except multiprocessing.TimeoutError:
            raise TimeoutError(
                f"Timeout to list Katib Suggestions in namespace: {namespace}"
            )
        except Exception:
            raise RuntimeError(
                f"Failed to list Katib Suggestions in namespace: {namespace}"
            )
        return result

    def get_trial(
        self,
        name: str,
        namespace: Optional[str] = None,
        timeout: int = constants.DEFAULT_TIMEOUT,
    ):
        """Get the Katib Trial.

        Args:
            name: Name for the Trial.
            namespace: Namespace for the Trial.
            timeout: Optional, Kubernetes API server timeout in seconds
                to execute the request.

        Returns:
            V1beta1Trial: Katib Trial object.

        Raises:
            TimeoutError: Timeout to get Katib Trial.
            RuntimeError: Failed to get Katib Trial.
        """

        namespace = namespace or self.namespace

        try:
            thread = self.custom_api.get_namespaced_custom_object(
                constants.KUBEFLOW_GROUP,
                constants.KATIB_VERSION,
                namespace,
                constants.TRIAL_PLURAL,
                name,
                async_req=True,
            )
            response = utils.FakeResponse(thread.get(timeout))
            trial = self.api_client.deserialize(response, models.V1beta1Trial)
            return trial

        except multiprocessing.TimeoutError:
            raise TimeoutError(f"Timeout to get Katib Trial: {namespace}/{name}")
        except Exception:
            raise RuntimeError(f"Failed to get Katib Trial: {namespace}/{name}")

    def list_trials(
        self,
        experiment_name: str = None,
        namespace: Optional[str] = None,
        timeout: int = constants.DEFAULT_TIMEOUT,
    ):
        """List of all Trials in namespace. If Experiment name is set,
        it returns all Trials belong to the Experiment.

        Args:
            experiment_name: Optional name for the Experiment.
            namespace: Namespace to list the Trials.
            timeout: Optional, Kubernetes API server timeout in seconds
                to execute the request.

        Returns:
            list[V1beta1Trial]: List of Katib Trial objects. It returns
            empty list if Trials cannot be found.

        Raises:
            TimeoutError: Timeout to list Katib Trials.
            RuntimeError: Failed to list Katib Trials.
        """

        namespace = namespace or self.namespace

        result = []
        try:
            if experiment_name is None:
                thread = self.custom_api.list_namespaced_custom_object(
                    constants.KUBEFLOW_GROUP,
                    constants.KATIB_VERSION,
                    namespace=namespace,
                    plural=constants.TRIAL_PLURAL,
                    async_req=True,
                )
            else:
                thread = self.custom_api.list_namespaced_custom_object(
                    constants.KUBEFLOW_GROUP,
                    constants.KATIB_VERSION,
                    namespace=namespace,
                    plural=constants.TRIAL_PLURAL,
                    label_selector=f"{constants.EXPERIMENT_LABEL}={experiment_name}",
                    async_req=True,
                )
            response = thread.get(timeout)
            result = [
                self.api_client.deserialize(
                    utils.FakeResponse(item), models.V1beta1Trial
                )
                for item in response.get("items")
            ]
        except multiprocessing.TimeoutError:
            raise TimeoutError(
                f"Timeout to list Katib Trials in namespace: {namespace}"
            )
        except Exception:
            raise RuntimeError(f"Failed to list Katib Trials in namespace: {namespace}")
        return result

    def get_success_trial_details(
        self,
        experiment_name: str = None,
        namespace: Optional[str] = None,
        timeout: int = constants.DEFAULT_TIMEOUT,
    ):
        """Get the Succeeded Trial details. If Experiment name is set,
        it returns Succeeded Trials details belong to the Experiment.

        Args:
            experiment_name: Optional name for the Experiment.
            namespace: Namespace to list the Trials.
            timeout: Optional, Kubernetes API server timeout in seconds
                to execute the request.

        Returns:
            list[dict]: Trial names with hyperparameters and metrics.
            It returns empty list if Succeeded Trials cannot be found.

        Raises:
            TimeoutError: Timeout to list Katib Trials.
            RuntimeError: Failed to list Katib Trials.
        """

        namespace = namespace or self.namespace

        result = []
        try:
            if experiment_name is None:
                thread = self.custom_api.list_namespaced_custom_object(
                    constants.KUBEFLOW_GROUP,
                    constants.KATIB_VERSION,
                    namespace=namespace,
                    plural=constants.TRIAL_PLURAL,
                    async_req=True,
                )
            else:
                thread = self.custom_api.list_namespaced_custom_object(
                    constants.KUBEFLOW_GROUP,
                    constants.KATIB_VERSION,
                    namespace=namespace,
                    plural=constants.TRIAL_PLURAL,
                    label_selector=f"{constants.EXPERIMENT_LABEL}={experiment_name}",
                    async_req=True,
                )
            response = thread.get(timeout)
            for item in response.get("items"):
                trial = self.api_client.deserialize(
                    utils.FakeResponse(item), models.V1beta1Trial
                )
                if (
                    trial.status
                    and trial.status.conditions
                    and len(trial.status.conditions) > 0
                ):
                    if utils.has_condition(
                        trial.status.conditions, constants.TRIAL_CONDITION_SUCCEEDED
                    ):
                        output = {}
                        output["name"] = trial.metadata.name
                        output["parameter_assignments"] = (
                            trial.spec.parameter_assignments
                        )
                        output["metrics"] = trial.status.observation.metrics
                        result.append(output)
        except multiprocessing.TimeoutError:
            raise TimeoutError(
                f"Timeout to list Katib Trials in namespace: {namespace}"
            )
        except Exception:
            raise RuntimeError(f"Failed to list Katib Trials in namespace: {namespace}")
        return result

    def get_optimal_hyperparameters(
        self,
        name: str,
        namespace: Optional[str] = None,
        timeout: int = constants.DEFAULT_TIMEOUT,
    ):
        """Get the current optimal Trial from the Experiment.

        Args:
            name: Name for the Experiment.
            namespace: Namespace for the Experiment.
            timeout: Optional, Kubernetes API server timeout in seconds
                to execute the request.

        Returns:
            V1beta1OptimalTrial: The most optimal Trial for the Experiment.
            It returns `None` if Experiment does not have optimal Trial yet.

        Raises:
            TimeoutError: Timeout to get Katib Experiment.
            RuntimeError: Failed to get Katib Experiment.
        """

        namespace = namespace or self.namespace

        experiment = self.get_experiment(name, namespace, timeout)
        if (
            experiment.status
            and experiment.status.current_optimal_trial
            and experiment.status.current_optimal_trial.observation.metrics
        ):
            return experiment.status.current_optimal_trial
        else:
            return None

    def get_trial_metrics(
        self,
        name: str,
        namespace: Optional[str] = None,
        db_manager_address: str = constants.DEFAULT_DB_MANAGER_ADDRESS,
        timeout: str = constants.DEFAULT_TIMEOUT,
    ):
        """Get the Trial Metric Results from the Katib DB.
        Katib DB Manager service should be accessible while calling this API.

        If you run this API in-cluster (e.g. from the Kubeflow Notebook) you can
        use the default Katib DB Manager address: `katib-db-manager.kubeflow:6789`.

        If you run this API outside the cluster, you have to port-forward the
        Katib DB Manager before getting the Trial metrics: `kubectl port-forward svc/katib-db-manager -n kubeflow 6789`.
        In that case, you can use this Katib DB Manager address: `localhost:6789`.

        You can use `curl` to verify that Katib DB Manager is reachable: `curl <db-manager-address>`.

        Args:
            name: Name for the Trial.
            namespace: Namespace for the Trial.
            db-manager-address: Address for the Katib DB Manager in this format: `ip-address:port`.
            timeout: Optional, gRPC API Server timeout in seconds to get metrics.

        Returns:
            List of MetricLog objects
            (https://github.com/kubeflow/katib/blob/4a2db414d85f29f17bc8ec6ff3462beef29585da/pkg/apis/manager/v1beta1/gen-doc/api.md#api-v1-beta1-MetricLog).
            For example, to get the first metric value run the following:
            `get_trial_metrics(...)[0].metric.value

        Raises:
            RuntimeError: Unable to get Trial metrics.
        """

        namespace = namespace or self.namespace

        db_manager_address = db_manager_address.split(":")
        channel = grpc.beta.implementations.insecure_channel(
            db_manager_address[0], int(db_manager_address[1])
        )

        with katib_api_pb2.beta_create_DBManager_stub(channel) as client:
            try:
                # When metric name is empty, we select all logs from the Katib DB.
                observation_logs = client.GetObservationLog(
                    katib_api_pb2.GetObservationLogRequest(trial_name=name),
                    timeout=timeout,
                )
            except Exception as e:
                raise RuntimeError(
                    f"Unable to get metrics for Trial {namespace}/{name}. Exception: {e}"
                )

            return observation_logs.observation_log.metric_logs
