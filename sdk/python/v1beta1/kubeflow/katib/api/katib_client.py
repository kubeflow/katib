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

import multiprocessing
from typing import Callable, List, Dict, Any
import inspect
import textwrap

from kubeflow.katib import models
from kubeflow.katib.api_client import ApiClient
from kubeflow.katib.constants import constants
from kubeflow.katib.utils import utils
from kubernetes import client, config


class KatibClient(object):
    def __init__(
        self,
        config_file=None,
        context=None,
        client_configuration=None,
        persist_config=True,
    ):
        """KatibClient constructor.

        :param config_file: Name of the kube-config file. Defaults to ~/.kube/config.
        :param context: Set the active context. Defaults to current_context from the kube-config.
        :param client_configuration: The kubernetes.client.Configuration to set configs to.
        :param persist_config: If True, config file will be updated when changed.
        """

        self.in_cluster = None
        if config_file or not utils.is_running_in_k8s():
            config.load_kube_config(
                config_file=config_file,
                context=context,
                client_configuration=client_configuration,
                persist_config=persist_config,
            )
            self.in_cluster = False
        else:
            config.load_incluster_config()
            self.in_cluster = True

        self.api_instance = client.CustomObjectsApi()
        self.api_client = ApiClient()

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
        self, exp_object, namespace=utils.get_default_target_namespace()
    ):
        """Create the Katib Experiment.

        :param exp_object: Experiment object.
        :param namespace: Experiment namespace.
        If the namespace is None, it takes namespace from the Experiment or "default".
        """

        try:
            self.api_instance.create_namespaced_custom_object(
                constants.KUBEFLOW_GROUP,
                constants.KATIB_VERSION,
                namespace,
                constants.EXPERIMENT_PLURAL,
                exp_object,
            )
        except client.rest.ApiException as e:
            raise RuntimeError(
                "Exception when calling CustomObjectsApi->create_namespaced_custom_object:\
         %s\n"
                % e
            )

        # TODO (andreyvelich): Use proper logger.
        print("Experiment {} has been created".format(exp_object.metadata.name))

        if self._is_ipython():
            if self.in_cluster:
                import IPython

                IPython.display.display(
                    IPython.display.HTML(
                        "Katib Experiment {} "
                        'link <a href="/_/katib/#/katib/hp_monitor/{}/{}" target="_blank">here</a>'.format(
                            exp_object.metadata.name,
                            namespace,
                            exp_object.metadata.name,
                        )
                    )
                )

    def tune(
        self,
        # TODO (andreyvelich): How to be consistent with other APIs (name) ?
        name: str,
        objective: Callable,
        parameters: Dict[str, Any],
        base_image: str = constants.BASE_IMAGE_TENSORFLOW,
        namespace: str = utils.get_default_target_namespace(),
        algorithm_name: str = "random",
        objective_metric_name: str = None,
        additional_metric_names: List[str] = [],
        objective_type: str = "maximize",
        objective_goal: float = None,
        max_trial_count: int = None,
        parallel_trial_count: int = None,
        max_failed_trial_count: int = None,
        retain_trials: bool = False,
        packages_to_install: List[str] = None,
        pip_index_url: str = "https://pypi.org/simple",
    ):
        """Create HyperParameter Tuning Katib Experiment from the objective function.

        Args:
            name: Name for the Experiment.
            objective: Objective function that Katib uses to train the model.
                This function must be Callable and it must have only one dict argument.
                Katib uses this argument to send HyperParameters to the function.
                The function should not use any code declared outside of the function
                definition. Import statements must be added inside the function.
            parameters: Dict of HyperParameters to tune your Experiment. You
                should use Katib SDK to define the search space for these parameters.

                For example: `parameters = {"lr": katib.search.double(min=0.1, max=0.2)}`

                Also, you can use these parameters to define input for your
                objective function.
            base_image: Image to use when executing the objective function.
            namespace: Namespace for the Experiment.
            algorithm_name: Search algorithm for the HyperParameter tuning.
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
            retain_trials: Whether Trials' resources (e.g. pods) are deleted after Succeeded state.
            packages_to_install: List of Python packages to install in addition
                to the base image packages. These packages are installed before
                executing the objective function.
            pip_index_url: The PyPI url from which to install Python packages.
        """

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
        experiment.spec.algorithm = models.V1beta1AlgorithmSpec(
            algorithm_name=algorithm_name
        )

        # Add Trial budget to the Katib Experiment.
        if max_trial_count is not None:
            experiment.spec.max_trial_count = max_trial_count
        if parallel_trial_count is not None:
            experiment.spec.parallel_trial_count = parallel_trial_count
        if max_failed_trial_count is not None:
            experiment.spec.max_failed_trial_count = max_failed_trial_count

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
        for p_name, p_value in parameters.items():
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

        # Create Trial specification.
        trial_spec = client.V1Job(
            api_version="batch/v1",
            kind="Job",
            spec=client.V1JobSpec(
                template=client.V1PodTemplateSpec(
                    metadata=models.V1ObjectMeta(
                        annotations={"sidecar.istio.io/inject": "false"}
                    ),
                    spec=client.V1PodSpec(
                        restart_policy="Never",
                        containers=[
                            client.V1Container(
                                name=constants.DEFAULT_PRIMARY_CONTAINER_NAME,
                                image=base_image,
                                command=["bash", "-c"],
                                args=[exec_script],
                            )
                        ],
                    ),
                )
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
        self.create_experiment(exp_object=experiment, namespace=namespace)

    # TODO (andreyvelich): Get Experiment should always return one Experiment.
    # Use list_experiments to return Experiment list.
    # That function should return Experiment object.
    def get_experiment(self, name=None, namespace=None):
        """Get the Katib Experiment.

        :param name: Experiment name.
        If the name is None returns all Experiments in the namespace.
        :param namespace: Experiment namespace.
        If the namespace is `None`, it takes namespace from the Experiment object or "default".


        :return: Experiment object.
        :rtype: dict
        """
        if namespace is None:
            namespace = utils.get_default_target_namespace()

        if name:
            thread = self.api_instance.get_namespaced_custom_object(
                constants.KUBEFLOW_GROUP,
                constants.KATIB_VERSION,
                namespace,
                constants.EXPERIMENT_PLURAL,
                name,
                async_req=True,
            )

            katibexp = None
            try:
                katibexp = thread.get(constants.APISERVER_TIMEOUT)
            except multiprocessing.TimeoutError:
                raise RuntimeError("Timeout trying to get katib experiment.")
            except client.rest.ApiException as e:
                raise RuntimeError(
                    "Exception when calling CustomObjectsApi->get_namespaced_custom_object:\
          %s\n"
                    % e
                )
            except Exception as e:
                raise RuntimeError(
                    "There was a problem to get experiment {0} in namespace {1}. Exception: \
          {2} ".format(
                        name, namespace, e
                    )
                )

        else:
            thread = self.api_instance.list_namespaced_custom_object(
                constants.KUBEFLOW_GROUP,
                constants.KATIB_VERSION,
                namespace,
                constants.EXPERIMENT_PLURAL,
                async_req=True,
            )

            katibexp = None
            try:
                katibexp = thread.get(constants.APISERVER_TIMEOUT)
            except multiprocessing.TimeoutError:
                raise RuntimeError("Timeout trying to get Experiment.")
            except client.rest.ApiException as e:
                raise RuntimeError(
                    "Exception when calling CustomObjectsApi->list_namespaced_custom_object:\
          %s\n"
                    % e
                )
            except Exception as e:
                raise RuntimeError(
                    "There was a problem to get experiment in namespace {0}. \
          Exception: {1} ".format(
                        namespace, e
                    )
                )

        return katibexp

    def get_suggestion(self, name=None, namespace=None):
        """Get the Katib Suggestion.

        :param name: Suggestion name.
        If the name is None returns all Suggestion in the namespace.
        :param namespace: Suggestion namespace.
        If the namespace is None, it takes namespace from the Suggestion object or "default".

        :return: Suggestion object.
        :rtype: dict
        """

        if namespace is None:
            namespace = utils.get_default_target_namespace()

        if name:
            thread = self.api_instance.get_namespaced_custom_object(
                constants.KUBEFLOW_GROUP,
                constants.KATIB_VERSION,
                namespace,
                constants.SUGGESTION_PLURAL,
                name,
                async_req=True,
            )

            katib_suggestion = None
            try:
                katib_suggestion = thread.get(constants.APISERVER_TIMEOUT)
            except multiprocessing.TimeoutError:
                raise RuntimeError("Timeout trying to get Katib suggestion")
            except client.rest.ApiException as e:
                raise RuntimeError(
                    "Exception when calling CustomObjectsApi->get_namespaced_custom_object:\
          %s\n"
                    % e
                )
            except Exception as e:
                raise RuntimeError(
                    "There was a problem to get suggestion {0} in namespace {1}. Exception: \
          {2} ".format(
                        name, namespace, e
                    )
                )

        else:
            thread = self.api_instance.list_namespaced_custom_object(
                constants.KUBEFLOW_GROUP,
                constants.KATIB_VERSION,
                namespace,
                constants.SUGGESTION_PLURAL,
                async_req=True,
            )

            katib_suggestion = None
            try:
                katib_suggestion = thread.get(constants.APISERVER_TIMEOUT)
            except multiprocessing.TimeoutError:
                raise RuntimeError("Timeout trying to get Katib suggestion")
            except client.rest.ApiException as e:
                raise RuntimeError(
                    "Exception when calling CustomObjectsApi->list_namespaced_custom_object:\
          %s\n"
                    % e
                )
            except Exception as e:
                raise RuntimeError(
                    "There was a problem to get suggestions in namespace {0}. \
          Exception: {1} ".format(
                        namespace, e
                    )
                )

        return katib_suggestion

    def delete_experiment(self, name, namespace=utils.get_default_target_namespace()):
        """Delete the Katib Experiment.

        :param name: Experiment name.
        :param namespace: Experiment namespace.
        If the namespace is None, it takes namespace from the Experiment object or "default".
        """

        try:
            self.api_instance.delete_namespaced_custom_object(
                constants.KUBEFLOW_GROUP,
                constants.KATIB_VERSION,
                namespace,
                constants.EXPERIMENT_PLURAL,
                name,
                body=client.V1DeleteOptions(),
            )
        except client.rest.ApiException as e:
            raise RuntimeError(
                "Exception when calling CustomObjectsApi->delete_namespaced_custom_object:\
         %s\n"
                % e
            )

        # TODO (andreyvelich): Use proper logger.
        print("Experiment {} has been deleted".format(name))

    def list_experiments(self, namespace=None):
        """List all Katib Experiments.

        :param namespace: Experiments namespace.
        If the namespace is None, it takes "default" namespace.

        :return: List of Experiment objects.
        :rtype: list[V1beta1Experiment]
        """
        if namespace is None:
            namespace = utils.get_default_target_namespace()

        thread = self.api_instance.list_namespaced_custom_object(
            constants.KUBEFLOW_GROUP,
            constants.KATIB_VERSION,
            namespace=namespace,
            plural=constants.EXPERIMENT_PLURAL,
            async_req=True,
        )

        katibexp = None
        result = []
        try:
            katibexp = thread.get(constants.APISERVER_TIMEOUT)
            result = [
                self.api_client.deserialize(
                    utils.FakeResponse(item), models.V1beta1Experiment
                )
                for item in katibexp.get("items")
            ]

        except multiprocessing.TimeoutError:
            raise RuntimeError("Timeout trying to get katib experiment.")
        except client.rest.ApiException as e:
            raise RuntimeError(
                "Exception when calling CustomObjectsApi->get_namespaced_custom_object:\
          %s\n"
                % e
            )
        except Exception as e:
            raise RuntimeError(
                "There was a problem to get experiments in namespace {0}. Exception: \
          {1} ".format(
                    namespace, e
                )
            )
        return result

    def get_experiment_status(self, name, namespace=None):
        """Get the Experiment current status.

        :param name: Experiment name.
        :param namespace: Experiment namespace.
        If the namespace is None, it takes "default" namespace.

        :return: Current Experiment status.
        :rtype: str
        """
        if namespace is None:
            namespace = utils.get_default_target_namespace()

        katibexp = self.get_experiment(name, namespace=namespace)
        last_condition = katibexp.get("status", {}).get("conditions", [])[-1]
        return last_condition.get("type", "")

    def is_experiment_succeeded(self, name, namespace=None):
        """Check if Experiment has succeeded.

        :param name: Experiment name.
        :param namespace: Experiment namespace.
        If the namespace is None, it takes "default" namespace.

        :return: Whether Experiment has succeeded or not.
        :rtype: bool
        """
        experiment_status = self.get_experiment_status(name, namespace=namespace)
        return experiment_status.lower() == "succeeded"

    def list_trials(self, name=None, namespace=None):
        """List all Experiment's Trials.

        :param name: Experiment name.
        :param namespace: Experiments namespace.
        If the namespace is None, it takes "default" namespace.

        :return: List of Trial objects
        :rtype: list[V1beta1Trial]
        """
        if namespace is None:
            namespace = utils.get_default_target_namespace()

        thread = self.api_instance.list_namespaced_custom_object(
            constants.KUBEFLOW_GROUP,
            constants.KATIB_VERSION,
            namespace=namespace,
            plural=constants.TRIAL_PLURAL,
            async_req=True,
        )

        katibtrial = None
        result = []
        try:
            katibtrial = thread.get(constants.APISERVER_TIMEOUT)

            for item in katibtrial.get("items"):
                if (
                    name is not None
                    and item.get("metadata", {}).get("ownerReferences")[0].get("name")
                    != name
                ):
                    continue

                result.append(
                    self.api_client.deserialize(
                        utils.FakeResponse(item), models.V1beta1Trial
                    )
                )
        except multiprocessing.TimeoutError:
            raise RuntimeError("Timeout trying to get katib experiment.")
        except client.rest.ApiException as e:
            raise RuntimeError(
                "Exception when calling CustomObjectsApi->list_namespaced_custom_object:\
          %s\n"
                % e
            )
        except Exception as e:
            raise RuntimeError(
                "There was a problem to get experiment {0} in namespace {1}. Exception: \
          {2} ".format(
                    name, namespace, e
                )
            )
        return result

    def get_success_trial_details(self, name=None, namespace=None):
        """Get the Trial details that have succeeded for an Experiment.

        :param name: Experiment name.
        :param namespace: Experiment namespace.
        If the namespace is None, it takes namespace from the Experiment or "default".

        :return: Trial names with the hyperparameters and metrics.
        :type: list[dict]
        """
        if namespace is None:
            namespace = utils.get_default_target_namespace()

        thread = self.api_instance.list_namespaced_custom_object(
            constants.KUBEFLOW_GROUP,
            constants.KATIB_VERSION,
            namespace=namespace,
            plural=constants.TRIAL_PLURAL,
            async_req=True,
        )

        katibtrial = None
        result = []
        try:
            katibtrial = thread.get(constants.APISERVER_TIMEOUT)

            for item in katibtrial.get("items"):
                status = item.get("status", {}).get("conditions", [])[-1].get("type")
                if status != "Succeeded":
                    continue

                if (
                    name is not None
                    and item.get("metadata", {}).get("ownerReferences")[0].get("name")
                    != name
                ):
                    continue

                output = {}
                output["name"] = item.get("metadata", {}).get("name")
                output["hyperparameters"] = item.get("spec", {}).get(
                    "parameterAssignments", []
                )
                output["metrics"] = (
                    item.get("status", {}).get("observation", {}).get("metrics", [])
                )
                result.append(output)
        except multiprocessing.TimeoutError:
            raise RuntimeError(
                "Timeout trying to get succeeded trials of the katib experiment."
            )
        except client.rest.ApiException as e:
            raise RuntimeError(
                "Exception when calling CustomObjectsApi->list_namespaced_custom_object:\
          %s\n"
                % e
            )
        except Exception as e:
            raise RuntimeError(
                "There was a problem to get experiment {0} in namespace {1}. Exception: \
          {2} ".format(
                    name, namespace, e
                )
            )

        return result

    def get_optimal_hyperparameters(self, name=None, namespace=None):
        """Get the current optimal Trial from the Experiment.

        :param name: Experiment name.
        :param namespace: Experiment namespace.

        :return: Current optimal Trial for the Experiment.
        :rtype: dict
        """
        if namespace is None:
            namespace = utils.get_default_target_namespace()

        katibexp = self.get_experiment(name, namespace=namespace)
        result = {}
        result["currentOptimalTrial"] = katibexp.get("status", {}).get(
            "currentOptimalTrial"
        )

        return result
