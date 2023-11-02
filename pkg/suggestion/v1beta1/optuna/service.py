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

import threading
import grpc
import logging
import itertools

from pkg.apis.manager.v1beta1.python import api_pb2
from pkg.apis.manager.v1beta1.python import api_pb2_grpc
from pkg.suggestion.v1beta1.internal.search_space import HyperParameterSearchSpace
from pkg.suggestion.v1beta1.internal.trial import Trial, Assignment
from pkg.suggestion.v1beta1.optuna.base_service import BaseOptunaService
from pkg.suggestion.v1beta1.internal.base_health_service import HealthServicer

logger = logging.getLogger(__name__)


class OptunaService(api_pb2_grpc.SuggestionServicer, HealthServicer):
    def __init__(self):
        super(OptunaService, self).__init__()
        self.lock = threading.Lock()
        self.base_service = None

    def GetSuggestions(self, request, context):
        """
        Main function to provide suggestion.
        """
        with self.lock:
            name, config = OptimizerConfiguration.convert_algorithm_spec(
                request.experiment.spec.algorithm
            )
            if self.base_service is None:
                search_space = HyperParameterSearchSpace.convert(request.experiment)
                self.base_service = BaseOptunaService(
                    algorithm_name=name,
                    algorithm_config=config,
                    search_space=search_space,
                )

            trials = Trial.convert(request.trials)
            list_of_assignments = self.base_service.get_suggestions(
                trials, request.current_request_number
            )
            return api_pb2.GetSuggestionsReply(
                parameter_assignments=Assignment.generate(list_of_assignments)
            )

    def ValidateAlgorithmSettings(self, request, context):
        is_valid, message = OptimizerConfiguration.validate_algorithm_spec(
            request.experiment
        )
        if not is_valid:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(message)
            logger.error(message)
        return api_pb2.ValidateAlgorithmSettingsReply()


class OptimizerConfiguration(object):
    __conversion_dict = {
        "tpe": {
            "n_startup_trials": lambda x: int(x),
            "n_ei_candidates": lambda x: int(x),
            "seed": lambda x: int(x),
            "constant_liar": True,
        },
        "multivariate-tpe": {
            "n_startup_trials": lambda x: int(x),
            "n_ei_candidates": lambda x: int(x),
            "seed": lambda x: int(x),
            "multivariate": "multivariate-tpe",
            "constant_liar": True,
        },
        "cmaes": {
            "restart_strategy": lambda x: None if x == "None" or x == "none" else x,
            "sigma0": lambda x: float(x),
            "seed": lambda x: int(x),
        },
        "random": {
            "seed": lambda x: int(x),
        },
        "grid": {
            "seed": lambda x: int(x),
        },
    }

    @classmethod
    def convert_algorithm_spec(cls, algorithm_spec):
        config = {}

        algorithm_name = algorithm_spec.algorithm_name
        setting_schema = cls.__conversion_dict[algorithm_name]
        for s in algorithm_spec.algorithm_settings:
            if s.name in setting_schema:
                config[s.name] = setting_schema[s.name](s.value)
            elif s.name == "sigma":
                config["sigma0"] = setting_schema["sigma0"](s.value)
            elif s.name == "random_state":
                config["seed"] = setting_schema["seed"](s.value)

        if algorithm_name == "tpe" or algorithm_name == "multivariate-tpe":
            config["constant_liar"] = setting_schema["constant_liar"]
        if algorithm_name == "multivariate-tpe":
            config["multivariate"] = setting_schema["multivariate"]

        return algorithm_spec.algorithm_name, config

    @classmethod
    def validate_algorithm_spec(cls, experiment):
        algorithm_spec = experiment.spec.algorithm
        algorithm_name = algorithm_spec.algorithm_name
        algorithm_settings = algorithm_spec.algorithm_settings
        parameters = experiment.spec.parameter_specs.parameters

        if algorithm_name == "tpe" or algorithm_name == "multivariate-tpe":
            return cls._validate_tpe_setting(algorithm_spec)
        elif algorithm_name == "cmaes":
            return cls._validate_cmaes_setting(algorithm_settings, parameters)
        elif algorithm_name == "random":
            return cls._validate_random_setting(algorithm_settings)
        elif algorithm_name == "grid":
            return cls._validate_grid_setting(experiment)
        else:
            return False, "unknown algorithm name {}".format(algorithm_name)

    @classmethod
    def _validate_tpe_setting(cls, algorithm_spec):
        algorithm_name = algorithm_spec.algorithm_name
        algorithm_settings = algorithm_spec.algorithm_settings

        for s in algorithm_settings:
            try:
                if s.name in ["n_startup_trials", "n_ei_candidates", "random_state"]:
                    if not int(s.value) >= 0:
                        return False, "{} should be greate or equal than zero".format(
                            s.name
                        )
                else:
                    return False, "unknown setting {} for algorithm {}".format(
                        s.name, algorithm_name
                    )
            except Exception as e:
                return False, "failed to validate {name}({value}): {exception}".format(
                    name=s.name, value=s.value, exception=e
                )

        return True, ""

    @classmethod
    def _validate_cmaes_setting(cls, algorithm_settings, parameters):
        for s in algorithm_settings:
            try:
                if s.name == "restart_strategy":
                    if s.value not in ["ipop", "None", "none"]:
                        return (
                            False,
                            "restart_strategy {} is not supported in CMAES optimization".format(
                                s.value
                            ),
                        )
                elif s.name == "sigma":
                    if not float(s.value) >= 0:
                        return False, "sigma should be greate or equal than zero"
                elif s.name == "random_state":
                    if not int(s.value) >= 0:
                        return False, "random_state should be greate or equal than zero"
                else:
                    return False, "unknown setting {} for algorithm cmaes".format(
                        s.name
                    )

            except Exception as e:
                return False, "failed to validate {name}({value}): {exception}".format(
                    name=s.name, value=s.value, exception=e
                )

        cnt = 0
        for p in parameters:
            if p.parameter_type == api_pb2.DOUBLE or p.parameter_type == api_pb2.INT:
                cnt += 1
        if cnt < 2:
            return (
                False,
                "cmaes only supports two or more dimensional continuous search space.",
            )

        return True, ""

    @classmethod
    def _validate_random_setting(cls, algorithm_settings):
        for s in algorithm_settings:
            try:
                if s.name == "random_state":
                    if not int(s.value) >= 0:
                        return False, ""
                else:
                    return False, "unknown setting {} for algorithm random".format(
                        s.name
                    )

            except Exception as e:
                return False, "failed to validate {name}({value}): {exception}".format(
                    name=s.name, value=s.value, exception=e
                )

        return True, ""

    @classmethod
    def _validate_grid_setting(cls, experiment):
        algorithm_settings = experiment.spec.algorithm.algorithm_settings
        search_space = HyperParameterSearchSpace.convert(experiment)

        for s in algorithm_settings:
            try:
                if s.name == "random_state":
                    if not int(s.value) >= 0:
                        return False, ""
                else:
                    return False, "unknown setting {} for algorithm grid".format(s.name)

            except Exception as e:
                return False, "failed to validate {name}({value}): {exception}".format(
                    name=s.name, value=s.value, exception=e
                )

        try:
            combinations = HyperParameterSearchSpace.convert_to_combinations(
                search_space
            )
            num_combinations = len(list(itertools.product(*combinations.values())))
            max_trial_count = experiment.spec.max_trial_count
            if max_trial_count > num_combinations:
                return (
                    False,
                    "Max Trial Count: {max_trial} > all possible search combinations: {combinations}".format(
                        max_trial=max_trial_count, combinations=num_combinations
                    ),
                )

        except Exception as e:
            return (
                False,
                "failed to validate parameters({parameters}): {exception}".format(
                    parameters=search_space.params, exception=e
                ),
            )

        return True, ""
