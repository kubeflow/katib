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

import optuna
from collections import defaultdict

from pkg.suggestion.v1beta1.internal.constant import INTEGER, DOUBLE, CATEGORICAL, DISCRETE, MAX_GOAL
from pkg.suggestion.v1beta1.internal.trial import Assignment
from pkg.suggestion.v1beta1.internal.search_space import HyperParameterSearchSpace


class BaseOptunaService(object):
    def __init__(self,
                 algorithm_name="",
                 algorithm_config=None,
                 search_space=None):
        self.algorithm_name = algorithm_name
        self.algorithm_config = algorithm_config
        self.search_space = search_space
        self.assignments_to_optuna_number = defaultdict(list)
        self.recorded_trial_names = set()
        self.study = None
        self._create_study()

    def _create_study(self):
        sampler = self._create_sampler()
        direction = "maximize" if self.search_space.goal == MAX_GOAL else "minimize"

        self.study = optuna.create_study(sampler=sampler, direction=direction)

    def _create_sampler(self):
        if self.algorithm_name == "tpe" or self.algorithm_name == "multivariate-tpe":
            return optuna.samplers.TPESampler(**self.algorithm_config)

        elif self.algorithm_name == "cmaes":
            return optuna.samplers.CmaEsSampler(**self.algorithm_config)

        elif self.algorithm_name == "random":
            return optuna.samplers.RandomSampler(**self.algorithm_config)

        elif self.algorithm_name == "grid":
            combinations = HyperParameterSearchSpace.convert_to_combinations(self.search_space)
            return optuna.samplers.GridSampler(combinations, **self.algorithm_config)

    def get_suggestions(self, trials, current_request_number):
        if len(trials) != 0:
            self._tell(trials)
        return self._ask(current_request_number)

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

    @staticmethod
    def _get_assignments_key(assignments):
        assignments = sorted(assignments, key=lambda a: a.name)
        assignments_str = [f"{a.name}:{a.value}" for a in assignments]
        return ",".join(assignments_str)

    def _get_optuna_search_space(self):
        search_space = {}
        for param in self.search_space.params:
            if param.type == INTEGER:
                search_space[param.name] = optuna.distributions.IntDistribution(int(param.min), int(param.max))
            elif param.type == DOUBLE:
                search_space[param.name] = optuna.distributions.FloatDistribution(float(param.min), float(param.max))
            elif param.type == CATEGORICAL or param.type == DISCRETE:
                search_space[param.name] = optuna.distributions.CategoricalDistribution(param.list)
        return search_space
