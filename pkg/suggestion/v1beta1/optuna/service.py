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

from collections import defaultdict
import threading

import optuna

from pkg.apis.manager.v1beta1.python import api_pb2
from pkg.apis.manager.v1beta1.python import api_pb2_grpc
from pkg.suggestion.v1beta1.internal.constant import INTEGER, DOUBLE, CATEGORICAL, DISCRETE, MAX_GOAL
from pkg.suggestion.v1beta1.internal.search_space import HyperParameterSearchSpace
from pkg.suggestion.v1beta1.internal.trial import Trial, Assignment
from pkg.suggestion.v1beta1.internal.base_health_service import HealthServicer


class OptunaService(api_pb2_grpc.SuggestionServicer, HealthServicer):

    def __init__(self):
        super(OptunaService, self).__init__()
        self.study = None
        self.search_space = None
        self.recorded_trial_names = set()
        self.assignments_to_optuna_number = defaultdict(list)
        self.lock = threading.Lock()

    def GetSuggestions(self, request, context):
        """
        Main function to provide suggestion.
        """
        with self.lock:
            if self.study is None:
                self.search_space = HyperParameterSearchSpace.convert(request.experiment)
                self.study = self._create_study(request.experiment.spec.algorithm, self.search_space)

            trials = Trial.convert(request.trials)

            if len(trials) != 0:
                self._tell(trials)
            list_of_assignments = self._ask(request.current_request_number)

            return api_pb2.GetSuggestionsReply(
                parameter_assignments=Assignment.generate(list_of_assignments)
            )

    def _create_study(self, algorithm_spec, search_space):
        sampler = self._create_sampler(algorithm_spec)
        direction = "maximize" if search_space.goal == MAX_GOAL else "minimize"

        study = optuna.create_study(sampler=sampler, direction=direction)

        return study

    def _create_sampler(self, algorithm_spec):
        name = algorithm_spec.algorithm_name
        settings = {s.name: s.value for s in algorithm_spec.algorithm_settings}

        if name == "tpe" or name == "multivariate-tpe":
            kwargs = {}
            for k, v in settings.items():
                if k == "n_startup_trials":
                    kwargs["n_startup_trials"] = int(v)
                elif k == "n_ei_candidates":
                    kwargs["n_ei_candidates"] = int(v)
                elif k == "random_state":
                    kwargs["seed"] = int(v)
                else:
                    raise ValueError("Unknown name for {}: {}".format(name, k))

            kwargs["multivariate"] = name == "multivariate-tpe"
            kwargs["constant_liar"] = True

            sampler = optuna.samplers.TPESampler(**kwargs)

        elif name == "cmaes":
            kwargs = {}
            for k, v in settings.items():
                if k == "restart_strategy":
                    kwargs["restart_strategy"] = v
                elif k == "sigma":
                    kwargs["sigma0"] = float(v)
                elif k == "random_state":
                    kwargs["seed"] = int(v)
                else:
                    raise ValueError("Unknown name for {}: {}".format(name, k))

            sampler = optuna.samplers.CmaEsSampler(**kwargs)

        elif name == "random":
            kwargs = {}
            for k, v in settings.items():
                if k == "random_state":
                    kwargs["seed"] = int(v)
                else:
                    raise ValueError("Unknown name for {}: {}".format(name, k))

            sampler = optuna.samplers.RandomSampler(**kwargs)

        else:
            raise ValueError("Unknown algorithm name: {}".format(name))

        return sampler

    def _ask(self, current_request_number):
        list_of_assignments = []
        for _ in range(current_request_number):
            optuna_trial = self.study.ask(fixed_distributions=self._get_optuna_search_space())

            assignments = [Assignment(k, v) for k, v in optuna_trial.params.items()]
            list_of_assignments.append(assignments)

            assignments_key = self._get_assignments_key(assignments)
            self.assignments_to_optuna_number[assignments_key].append(optuna_trial.number)

        return list_of_assignments

    def _tell(self, trials):
        for trial in trials:
            if trial.name not in self.recorded_trial_names:
                self.recorded_trial_names.add(trial.name)

                value = float(trial.target_metric.value)
                assignments_key = self._get_assignments_key(trial.assignments)
                optuna_trial_numbers = self.assignments_to_optuna_number[assignments_key]

                if len(optuna_trial_numbers) != 0:
                    trial_number = optuna_trial_numbers.pop(0)
                    self.study.tell(trial_number, value)
                else:
                    raise ValueError("An unknown trial has been passed in the GetSuggestion request.")

    def _get_assignments_key(self, assignments):
        assignments = sorted(assignments, key=lambda a: a.name)
        assignments_str = [f"{a.name}:{a.value}" for a in assignments]
        return ",".join(assignments_str)

    def _get_optuna_search_space(self):
        search_space = {}
        for param in self.search_space.params:
            if param.type == INTEGER:
                search_space[param.name] = optuna.distributions.IntUniformDistribution(int(param.min), int(param.max))
            elif param.type == DOUBLE:
                search_space[param.name] = optuna.distributions.UniformDistribution(float(param.min), float(param.max))
            elif param.type == CATEGORICAL or param.type == DISCRETE:
                search_space[param.name] = optuna.distributions.CategoricalDistribution(param.list)
        return search_space

    def _get_casted_assignment_value(self, assignment):
        for param in self.search_space.params:
            if param.name == assignment.name:
                if param.type == INTEGER:
                    return int(assignment.value)
                elif param.type == DOUBLE:
                    return float(assignment.value)
                elif param.type == CATEGORICAL or param.type == DISCRETE:
                    return assignment.value
                else:
                    raise ValueError("Unknown parameter type: {}".format(param.type))
        raise ValueError("Parameter not found in the search space: {}".format(param.name))
