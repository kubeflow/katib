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

import grpc

from pkg.apis.manager.v1beta1.python import api_pb2
from pkg.apis.manager.v1beta1.python import api_pb2_grpc
from pkg.suggestion.v1beta1.internal.base_health_service import HealthServicer
from pkg.suggestion.v1beta1.internal.search_space import HyperParameterSearchSpace
from pkg.suggestion.v1beta1.internal.trial import Trial, Assignment
from pkg.suggestion.v1beta1.skopt.base_service import BaseSkoptService

logger = logging.getLogger(__name__)


class SkoptService(api_pb2_grpc.SuggestionServicer, HealthServicer):

    def __init__(self):
        super(SkoptService, self).__init__()
        self.base_service = None
        self.is_first_run = True

    def GetSuggestions(self, request, context):
        """
        Main function to provide suggestion.
        """
        algorithm_name, config = OptimizerConfiguration.convert_algorithm_spec(
            request.experiment.spec.algorithm)

        if self.is_first_run:
            search_space = HyperParameterSearchSpace.convert(request.experiment)
            self.base_service = BaseSkoptService(
                base_estimator=config.base_estimator,
                n_initial_points=config.n_initial_points,
                acq_func=config.acq_func,
                acq_optimizer=config.acq_optimizer,
                random_state=config.random_state,
                search_space=search_space)
            self.is_first_run = False

        trials = Trial.convert(request.trials)
        new_trials = self.base_service.getSuggestions(trials, request.current_request_number)
        return api_pb2.GetSuggestionsReply(
            parameter_assignments=Assignment.generate(new_trials)
        )

    def ValidateAlgorithmSettings(self, request, context):
        is_valid, message = OptimizerConfiguration.validate_algorithm_spec(
            request.experiment.spec.algorithm)
        if not is_valid:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(message)
            logger.error(message)
        return api_pb2.ValidateAlgorithmSettingsReply()


class OptimizerConfiguration(object):
    def __init__(self, base_estimator="GP",
                 n_initial_points=10,
                 acq_func="gp_hedge",
                 acq_optimizer="auto",
                 random_state=None):
        self.base_estimator = base_estimator
        self.n_initial_points = n_initial_points
        self.acq_func = acq_func
        self.acq_optimizer = acq_optimizer
        self.random_state = random_state

    @staticmethod
    def convert_algorithm_spec(algorithm_spec):
        optimizer = OptimizerConfiguration()
        for s in algorithm_spec.algorithm_settings:
            if s.name == "base_estimator":
                optimizer.base_estimator = s.value
            elif s.name == "n_initial_points":
                optimizer.n_initial_points = int(s.value)
            elif s.name == "acq_func":
                optimizer.acq_func = s.value
            elif s.name == "acq_optimizer":
                optimizer.acq_optimizer = s.value
            elif s.name == "random_state":
                optimizer.random_state = int(s.value)
        return algorithm_spec.algorithm_name, optimizer

    @classmethod
    def validate_algorithm_spec(cls, algorithm_spec):
        algo_name = algorithm_spec.algorithm_name

        if algo_name == "bayesianoptimization":
            return cls._validate_bayesianoptimization_setting(algorithm_spec.algorithm_settings)
        else:
            return False, "unknown algorithm name {}".format(algo_name)

    @classmethod
    def _validate_bayesianoptimization_setting(cls, algorithm_settings):
        for s in algorithm_settings:
            try:
                if s.name == "base_estimator":
                    if s.value not in ["GP", "RF", "ET", "GBRT"]:
                        return False, "base_estimator {} is not supported in Bayesian optimization".format(s.value)
                elif s.name == "n_initial_points":
                    if not (int(s.value) >= 0):
                        return False, "n_initial_points should be great or equal than zero"
                elif s.name == "acq_func":
                    if s.value not in ["gp_hedge", "LCB", "EI", "PI", "EIps", "PIps"]:
                        return False, "acq_func {} is not supported in Bayesian optimization".format(s.value)
                elif s.name == "acq_optimizer":
                    if s.value not in ["auto", "sampling", "lbfgs"]:
                        return False, "acq_optimizer {} is not supported in Bayesian optimization".format(s.value)
                elif s.name == "random_state":
                    if not (int(s.value) >= 0):
                        return False, "random_state should be great or equal than zero"
                else:
                    return False, "unknown setting {} for algorithm bayesianoptimization".format(s.name)
            except Exception as e:
                return False, "failed to validate {name}({value}): {exception}".format(name=s.name, value=s.value,
                                                                                       exception=e)

        return True, ""
