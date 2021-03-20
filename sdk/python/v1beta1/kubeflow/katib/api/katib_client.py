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

from kubernetes import client, config

from kubeflow.katib.constants import constants
from kubeflow.katib.utils import utils


class KatibClient(object):

    def __init__(self, config_file=None, context=None,
                 client_configuration=None, persist_config=True):
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
                persist_config=persist_config)
            self.in_cluster = False
        else:
            config.load_incluster_config()
            self.in_cluster = True

        self.api_instance = client.CustomObjectsApi()

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

    def create_experiment(self, exp_object, namespace=None):
        """Create the Katib Experiment.

        :param exp_object: Experiment object.
        :param namespace: Experiment namespace.
        If the namespace is None, it takes namespace from the Experiment or "default".

        :return: Created Experiment.
        :rtype: dict
        """

        if namespace is None:
            namespace = utils.set_katib_namespace(exp_object)
        try:
            outputs = self.api_instance.create_namespaced_custom_object(
                constants.KUBEFLOW_GROUP,
                constants.KATIB_VERSION,
                namespace,
                constants.EXPERIMENT_PLURAL,
                exp_object)
        except client.rest.ApiException as e:
            raise RuntimeError(
                "Exception when calling CustomObjectsApi->create_namespaced_custom_object:\
         %s\n" % e)

        if self._is_ipython():
            if self.in_cluster:
                import IPython
                IPython.display.display(IPython.display.HTML(
                    'Katib Experiment {} '
                    'link <a href="/_/katib/#/katib/hp_monitor/{}/{}" target="_blank">here</a>'.format(
                        exp_object.metadata.name, namespace, exp_object.metadata.name)
                ))
        return outputs

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
                async_req=True)

            katibexp = None
            try:
                katibexp = thread.get(constants.APISERVER_TIMEOUT)
            except multiprocessing.TimeoutError:
                raise RuntimeError("Timeout trying to get katib experiment.")
            except client.rest.ApiException as e:
                raise RuntimeError(
                    "Exception when calling CustomObjectsApi->get_namespaced_custom_object:\
          %s\n" % e)
            except Exception as e:
                raise RuntimeError(
                    "There was a problem to get experiment {0} in namespace {1}. Exception: \
          {2} ".format(name, namespace, e))

        else:
            thread = self.api_instance.list_namespaced_custom_object(
                constants.KUBEFLOW_GROUP,
                constants.KATIB_VERSION,
                namespace,
                constants.EXPERIMENT_PLURAL,
                async_req=True)

            katibexp = None
            try:
                katibexp = thread.get(constants.APISERVER_TIMEOUT)
            except multiprocessing.TimeoutError:
                raise RuntimeError("Timeout trying to get Experiment.")
            except client.rest.ApiException as e:
                raise RuntimeError(
                    "Exception when calling CustomObjectsApi->list_namespaced_custom_object:\
          %s\n" % e)
            except Exception as e:
                raise RuntimeError(
                    "There was a problem to get experiment in namespace {0}. \
          Exception: {1} ".format(namespace, e))

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
                async_req=True)

            katib_suggestion = None
            try:
                katib_suggestion = thread.get(constants.APISERVER_TIMEOUT)
            except multiprocessing.TimeoutError:
                raise RuntimeError("Timeout trying to get Katib suggestion")
            except client.rest.ApiException as e:
                raise RuntimeError(
                    "Exception when calling CustomObjectsApi->get_namespaced_custom_object:\
          %s\n" % e)
            except Exception as e:
                raise RuntimeError(
                    "There was a problem to get suggestion {0} in namespace {1}. Exception: \
          {2} ".format(name, namespace, e))

        else:
            thread = self.api_instance.list_namespaced_custom_object(
                constants.KUBEFLOW_GROUP,
                constants.KATIB_VERSION,
                namespace,
                constants.SUGGESTION_PLURAL,
                async_req=True)

            katib_suggestion = None
            try:
                katib_suggestion = thread.get(constants.APISERVER_TIMEOUT)
            except multiprocessing.TimeoutError:
                raise RuntimeError("Timeout trying to get Katib suggestion")
            except client.rest.ApiException as e:
                raise RuntimeError(
                    "Exception when calling CustomObjectsApi->list_namespaced_custom_object:\
          %s\n" % e)
            except Exception as e:
                raise RuntimeError(
                    "There was a problem to get suggestions in namespace {0}. \
          Exception: {1} ".format(namespace, e))

        return katib_suggestion

    def delete_experiment(self, name, namespace=None):
        """Delete the Katib Experiment.

        :param name: Experiment name.
        :param namespace: Experiment namespace.
        If the namespace is None, it takes namespace from the Experiment object or "default".

        :return: Deleted Experiment object.
        :rtype: dict
        """
        if namespace is None:
            namespace = utils.get_default_target_namespace()

        try:
            return self.api_instance.delete_namespaced_custom_object(
                constants.KUBEFLOW_GROUP,
                constants.KATIB_VERSION,
                namespace,
                constants.EXPERIMENT_PLURAL,
                name,
                body=client.V1DeleteOptions())
        except client.rest.ApiException as e:
            raise RuntimeError(
                "Exception when calling CustomObjectsApi->delete_namespaced_custom_object:\
         %s\n" % e)

    def list_experiments(self, namespace=None):
        """List all Katib Experiments.

        :param namespace: Experiments namespace.
        If the namespace is None, it takes "default" namespace.

        :return: List of Experiment names with the statuses.
        :rtype: list[dict]
        """
        if namespace is None:
            namespace = utils.get_default_target_namespace()

        thread = self.api_instance.list_namespaced_custom_object(
            constants.KUBEFLOW_GROUP,
            constants.KATIB_VERSION,
            namespace=namespace,
            plural=constants.EXPERIMENT_PLURAL,
            async_req=True)

        katibexp = None
        try:
            katibexp = thread.get(constants.APISERVER_TIMEOUT)
            result = []
            for i in katibexp.get("items"):
                output = {}
                output["name"] = i.get("metadata", {}).get("name")
                output["status"] = i.get("status", {}).get(
                    "conditions", [])[-1].get("type")
                result.append(output)
        except multiprocessing.TimeoutError:
            raise RuntimeError("Timeout trying to get katib experiment.")
        except client.rest.ApiException as e:
            raise RuntimeError(
                "Exception when calling CustomObjectsApi->get_namespaced_custom_object:\
          %s\n" % e)
        except Exception as e:
            raise RuntimeError(
                "There was a problem to get experiments in namespace {1}. Exception: \
          {2} ".format(namespace, e))
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
        experiment_status = self.get_experiment_status(
            name, namespace=namespace)
        return experiment_status.lower() == "succeeded"

    def list_trials(self, name=None, namespace=None):
        """List all Experiment's Trials.

        :param name: Experiment name.
        :param namespace: Experiments namespace.
        If the namespace is None, it takes "default" namespace.

        :return: List of Trial names with the statuses.
        :rtype: list[dict]
        """
        if namespace is None:
            namespace = utils.get_default_target_namespace()

        thread = self.api_instance.list_namespaced_custom_object(
            constants.KUBEFLOW_GROUP,
            constants.KATIB_VERSION,
            namespace=namespace,
            plural=constants.TRIAL_PLURAL,
            async_req=True)

        katibtrial = None
        try:
            katibtrial = thread.get(constants.APISERVER_TIMEOUT)
            result = []
            for i in katibtrial.get("items"):
                output = {}
                if i.get("metadata", {}).get("ownerReferences")[0].get("name") == name:
                    output["name"] = i.get("metadata", {}).get("name")
                    output["status"] = i.get("status", {}).get(
                        "conditions", [])[-1].get("type")
                    result.append(output)
        except multiprocessing.TimeoutError:
            raise RuntimeError("Timeout trying to getkatib experiment.")
        except client.rest.ApiException as e:
            raise RuntimeError(
                "Exception when calling CustomObjectsApi->get_namespaced_custom_object:\
          %s\n" % e)
        except Exception as e:
            raise RuntimeError(
                "There was a problem to get experiment {0} in namespace {1}. Exception: \
          {2} ".format(name, namespace, e))
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
            async_req=True)

        katibtrial = None
        try:
            katibtrial = thread.get(constants.APISERVER_TIMEOUT)
            result = []
            for i in katibtrial.get("items"):
                output = {}
                if i.get("metadata", {}).get("ownerReferences")[0].get("name") == name:
                    status = i.get("status", {}).get("conditions", [])[-1].get("type")
                    if status == "Succeeded":
                        output["name"] = i.get("metadata", {}).get("name")
                        output["hyperparameters"] = i.get("spec", {}).get("parameterAssignments", [])
                        output["metrics"] = (
                            i.get("status", {})
                            .get("observation", {})
                            .get("metrics", [])
                        )
                        result.append(output)
        except multiprocessing.TimeoutError:
            raise RuntimeError("Timeout trying to getkatib experiment.")
        except client.rest.ApiException as e:
            raise RuntimeError(
                "Exception when calling CustomObjectsApi->get_namespaced_custom_object:\
          %s\n" % e)
        except Exception as e:
            raise RuntimeError(
                "There was a problem to get experiment {0} in namespace {1}. Exception: \
          {2} ".format(name, namespace, e))
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
        result["currentOptimalTrial"] = katibexp.get(
            "status", {}).get("currentOptimalTrial")

        return result
