# Copyright 2022 The Kubeflow Authors.
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

import logging
import multiprocessing
from datetime import datetime
from typing import Iterable, Optional

import grpc
from kubernetes import client, config

from pkg.apis.manager.v1beta1.python import api_pb2, api_pb2_grpc

logger = logging.getLogger()
logging.basicConfig(level=logging.INFO)

STATUS_EARLY_STOPPED = "EarlyStopped"
KUBEFLOW_GROUP = "kubeflow.org"
KATIB_VERSION = "v1beta1"
TRIAL_PLURAL = "trials"
APISERVER_TIMEOUT = 120

DEFAULT_NAMESPACE = "default"

SUCCEEDED_TRIAL = api_pb2.TrialStatus.TrialConditionType.SUCCEEDED


class MedianStopService(api_pb2_grpc.EarlyStoppingServicer):
    def __init__(self):
        super(MedianStopService, self).__init__()
        self.is_first_run = True
        # Default settings
        self.min_trials_required = 3
        self.start_step = 4
        # trials_avg_history is the dict with succeeded Trials history where
        # key = Trial name, value = average value for "start_step" reported metrics.
        self.trials_avg_history = {}

        # Assume that Trial namespace = Suggestion namespace.
        try:
            with open(
                "/var/run/secrets/kubernetes.io/serviceaccount/namespace", "r"
            ) as f:
                self.namespace = f.readline()
                # Set config and api instance for k8s client.
                config.load_incluster_config()
        # This is used when service is not running in k8s, e.g. for unit tests.
        except Exception as e:
            logger.info(
                '{}. Service is not running in Kubernetes Pod, "{}" namespace is used'.format(
                    e, DEFAULT_NAMESPACE
                )
            )
            self.namespace = DEFAULT_NAMESPACE
            # Set config and api instance for k8s client.
            config.load_kube_config()

        self.api_instance = client.CustomObjectsApi()

    def ValidateEarlyStoppingSettings(
        self,
        request: api_pb2.ValidateEarlyStoppingSettingsRequest,
        context: grpc.ServicerContext,
    ) -> api_pb2.ValidateEarlyStoppingSettingsReply:
        is_valid, message = self.validate_early_stopping_spec(request.early_stopping)
        if not is_valid:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(message)
            logger.error(message)
        return api_pb2.ValidateEarlyStoppingSettingsReply()

    def validate_early_stopping_spec(self, early_stopping_spec):
        algorithm_name = early_stopping_spec.algorithm_name
        if algorithm_name == "medianstop":
            return self.validate_medianstop_setting(
                early_stopping_spec.algorithm_settings
            )
        else:
            return False, "unknown algorithm name {}".format(algorithm_name)

    @staticmethod
    def validate_medianstop_setting(early_stopping_settings):
        for setting in early_stopping_settings:
            try:
                if setting.name == "min_trials_required":
                    if not (int(setting.value) > 0):
                        return (
                            False,
                            "min_trials_required must be greater than zero (>0)",
                        )
                elif setting.name == "start_step":
                    if not (int(setting.value) >= 1):
                        return (
                            False,
                            "start_step must be greater or equal than one (>=1)",
                        )
                else:
                    return False, "unknown setting {} for algorithm medianstop".format(
                        setting.name
                    )

            except Exception as e:
                return False, "failed to validate {}({}): {}".format(
                    setting.name, setting.value, e
                )

        return True, ""

    def GetEarlyStoppingRules(
        self,
        request: api_pb2.GetEarlyStoppingRulesRequest,
        context: grpc.ServicerContext,
    ) -> api_pb2.GetSuggestionsReply:
        logger.info("Get new early stopping rules")

        # Get required values for the first call.
        if self.is_first_run:
            self.is_first_run = False
            # Get early stopping settings.
            self.get_early_stopping_settings(
                request.experiment.spec.early_stopping.algorithm_settings
            )
            logger.info(
                "Median stopping settings are: min_trials_required: {}, start_step: {}".format(
                    self.min_trials_required, self.start_step
                )
            )

            # Get comparison type and objective metric
            if request.experiment.spec.objective.type == api_pb2.MAXIMIZE:
                self.comparison = api_pb2.LESS
            else:
                self.comparison = api_pb2.GREATER
            self.objective_metric = (
                request.experiment.spec.objective.objective_metric_name
            )

            # Get DB manager address. It should have host and port.
            # For example: katib-db-manager.kubeflow:6789 - default one.
            self.db_manager_address = request.db_manager_address.split(":")
            if len(self.db_manager_address) != 2:
                raise Exception(
                    "Invalid Katib DB manager service address: {}".format(
                        self.db_manager_address
                    )
                )

        early_stopping_rules = []

        median = self.get_median_value(request.trials)
        if median is not None:
            early_stopping_rules.append(
                api_pb2.EarlyStoppingRule(
                    name=self.objective_metric,
                    value=str(median),
                    comparison=self.comparison,
                    start_step=self.start_step,
                )
            )

        logger.info(
            "New early stopping rules are:\n {}\n\n".format(early_stopping_rules)
        )
        return api_pb2.GetEarlyStoppingRulesReply(
            early_stopping_rules=early_stopping_rules
        )

    def get_early_stopping_settings(
        self, early_stopping_settings: Iterable[api_pb2.EarlyStoppingSetting]
    ):
        for setting in early_stopping_settings:
            if setting.name == "min_trials_required":
                self.min_trials_required = int(setting.value)
            elif setting.name == "start_step":
                self.start_step = int(setting.value)

    def get_median_value(self, trials: Iterable[api_pb2.Trial]) -> Optional[float]:
        for trial in trials:
            # Get metrics only for the new succeeded Trials.
            if (
                trial.name not in self.trials_avg_history
                and trial.status.condition == SUCCEEDED_TRIAL
            ):
                with grpc.insecure_channel(
                    f"{self.db_manager_address[0]}:{self.db_manager_address[1]}"
                ) as channel:
                    stub = api_pb2_grpc.DBManagerStub(channel)
                    get_log_response: api_pb2.GetObservationLogReply = (
                        stub.GetObservationLog(
                            api_pb2.GetObservationLogRequest(
                                trial_name=trial.name, metric_name=self.objective_metric
                            ),
                            timeout=APISERVER_TIMEOUT,
                        )
                    )

                # Get only first start_step metrics.
                # Since metrics are collected consistently and ordered by time,
                # we slice top start_step metrics.
                first_x_logs = get_log_response.observation_log.metric_logs[
                    : self.start_step
                ]
                metric_sum = 0
                for log in first_x_logs:
                    metric_sum += float(log.metric.value)

                # Get average metric value for the Trial.
                new_average = metric_sum / len(first_x_logs)
                self.trials_avg_history[trial.name] = new_average
                logger.info(
                    "Adding new succeeded Trial: {} with average metrics value: {}".format(
                        trial.name, new_average
                    )
                )
                logger.info(
                    "Trials average log history: {}".format(self.trials_avg_history)
                )

        # If count of succeeded Trials is greater than min_trials_required, calculate median.
        if len(self.trials_avg_history) >= self.min_trials_required:
            median = sum(list(self.trials_avg_history.values())) / len(
                self.trials_avg_history
            )
            logger.info("Generate new Median value: {}".format(median))
            return median
        # Else, return None.
        logger.info(
            "Count of succeeded Trials: {} is less than min_trials_required: {}".format(
                len(self.trials_avg_history), self.min_trials_required
            )
        )
        return None

    def SetTrialStatus(
        self, request: api_pb2.SetTrialStatusRequest, context: grpc.ServicerContext
    ) -> api_pb2.SetTrialStatusReply:
        trial_name = request.trial_name

        logger.info("Update status for Trial: {}".format(trial_name))

        # TODO (andreyvelich): Move this part to Katib SDK ?
        # Get Trial object
        thread = self.api_instance.get_namespaced_custom_object(
            KUBEFLOW_GROUP,
            KATIB_VERSION,
            self.namespace,
            TRIAL_PLURAL,
            trial_name,
            async_req=True,
        )

        trial = None
        try:
            trial = thread.get(APISERVER_TIMEOUT)
        except multiprocessing.TimeoutError:
            raise Exception("Timeout trying to get Katib Trial")
        except Exception as e:
            raise Exception(
                "Get Trial: {} in namespace: {} failed. Exception: {}".format(
                    trial_name, self.namespace, e
                )
            )

        time_now = datetime.now().strftime("%Y-%m-%dT%H:%M:%SZ")

        early_stopped_condition = {
            "type": STATUS_EARLY_STOPPED,
            "status": "True",
            "reason": "TrialEarlyStopped",
            "message": "Trial is early stopped",
            "lastUpdateTime": time_now,
            "lastTransitionTime": time_now,
        }
        trial["status"]["conditions"].append(early_stopped_condition)

        # Update Trial object with early stopped status
        try:
            self.api_instance.patch_namespaced_custom_object_status(
                KUBEFLOW_GROUP,
                KATIB_VERSION,
                self.namespace,
                TRIAL_PLURAL,
                trial_name,
                trial,
                async_req=True,
            )
        except Exception as e:
            raise Exception(
                "Update status for Trial: {} in namespace: {} failed. Exception: {}".format(
                    trial_name, self.namespace, e
                )
            )

        logger.info(
            "Changed status to: {} for Trial: {} in namespace: {}\n\n".format(
                STATUS_EARLY_STOPPED, trial_name, self.namespace
            )
        )

        return api_pb2.SetTrialStatusReply()
