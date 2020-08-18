# Copyright 2019 The Kubeflow Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
import multiprocessing
import time

from kubernetes import client, config

from kubeflow.katib.constants import constants
from kubeflow.katib.utils import utils


class KatibClient(object):

    def __init__(self, config_file=None, context=None,  # pylint: disable=too-many-arguments
                 client_configuration=None, persist_config=True):
        """
        katibclient constructor
        :param config_file: kubeconfig file, defaults to ~/.kube/config
        :param context: kubernetes context
        :param client_configuration: kubernetes configuration object
        :param persist_config:
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
        """
        Create the katib experiment
        :param exp_object: experiment object
        :param namespace: defaults to current or default namespace
        :return: created experiment dict
        """

        if namespace is None:
            namespace = utils.set_katib_namespace(exp_object)
        try:
            outputs = self.api_instance.create_namespaced_custom_object(
                constants.EXPERIMENT_GROUP,
                constants.EXPERIMENT_VERSION,
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
                html = \
                    ('Katib Experiment link <a href="/_/katib/#/katib/hp_monitor/%s/%s" target="_blank">here</a>'
                     % (namespace, exp_object.metadata.name))
                IPython.display.display(IPython.display.HTML(html))
        self.in_cluster = None
        return outputs

    def get_experiment(self, name=None, namespace=None):
        """
        Get single experiment or all experiment
        :param name: existing experiment name optional
        :param namespace: defaults to current or default namespace
        :return: experiment dict
        """
        if namespace is None:
            namespace = utils.get_default_target_namespace()

        if name:
            thread = self.api_instance.get_namespaced_custom_object(
                constants.EXPERIMENT_GROUP,
                constants.EXPERIMENT_VERSION,
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
                constants.EXPERIMENT_GROUP,
                constants.EXPERIMENT_VERSION,
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

    def delete_experiment(self, name, namespace=None):
        """
        Delete experiment
        :param name: experiment name required
        :param namespace: defaults to current or default namespace
        :return: status dict
        """
        if namespace is None:
            namespace = utils.get_default_target_namespace()

        try:
            return self.api_instance.delete_namespaced_custom_object(
                constants.EXPERIMENT_GROUP,
                constants.EXPERIMENT_VERSION,
                namespace,
                constants.EXPERIMENT_PLURAL,
                name,
                body=client.V1DeleteOptions())
        except client.rest.ApiException as e:
            raise RuntimeError(
                "Exception when calling CustomObjectsApi->delete_namespaced_custom_object:\
         %s\n" % e)

    def list_experiments(self, namespace=None):
        """
        List all experiments
        :param namespace: defaults to current or default namespace
        :return: list of experiment name with status as list
        """
        if namespace is None:
            namespace = utils.get_default_target_namespace()

        thread = self.api_instance.list_namespaced_custom_object(
            constants.EXPERIMENT_GROUP,
            constants.EXPERIMENT_VERSION,
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
                "There was a problem to get experiment {0} in namespace {1}. Exception: \
          {2} ".format(name, namespace, e))
        return result

    def get_experiment_status(self, name, namespace=None):
        """Returns experiment status, such as Running, Failed or Succeeded.
        Args:
          :param name: An experiment name. required
          :param namespace: defaults to current or default namespace.
          :return: status str
        """
        if namespace is None:
            namespace = utils.get_default_target_namespace()

        katibexp = self.get_experiment(name, namespace=namespace)
        last_condition = katibexp.get("status", {}).get("conditions", [])[-1]
        return last_condition.get("type", "")

    def is_experiment_succeeded(self, name, namespace=None):
        """Returns true if experiment succeeded; false otherwise.
        Args:
          :param name: An experiment name. required
          :param namespace: defaults to current or default namespace. optional
          :return: status bool
        """
        experiment_status = self.get_experiment_status(
            name, namespace=namespace)
        return experiment_status.lower() == "succeeded"

    def list_trials(self, name=None, namespace=None):
        """
        Get trials of an experiment
        :param name: existing experiment name
        :param namespace: defaults to current or default namespace
        :return: trials name with status as list
        """
        if namespace is None:
            namespace = utils.get_default_target_namespace()

        thread = self.api_instance.list_namespaced_custom_object(
            constants.EXPERIMENT_GROUP,
            constants.EXPERIMENT_VERSION,
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

    def get_optimal_hyperparmeters(self, name=None, namespace=None):
        """
        Get status, currentOptimalTrial with paramaterAssignments
        :param name: existing experiment name
        :param namespace: defaults to current or default namespace
        :return: dict with status, currentOptimalTrial with paramaterAssignments of an experiment
        """
        if namespace is None:
            namespace = utils.get_default_target_namespace()

        katibexp = self.get_experiment(name, namespace=namespace)
        result = {}
        result["currentOptimalTrial"] = katibexp.get(
            "status", {}).get("currentOptimalTrial")

        return result

