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

import datetime
import logging

import skopt

from pkg.suggestion.v1beta1.internal.constant import (
    CATEGORICAL,
    DISCRETE,
    DOUBLE,
    INTEGER,
    MAX_GOAL,
)
from pkg.suggestion.v1beta1.internal.trial import Assignment

logger = logging.getLogger(__name__)


class BaseSkoptService(object):
    """
    Refer to https://github.com/scikit-optimize/scikit-optimize .
    """

    def __init__(
        self,
        base_estimator="GP",
        n_initial_points=10,
        acq_func="gp_hedge",
        acq_optimizer="auto",
        random_state=None,
        search_space=None,
    ):
        self.base_estimator = base_estimator
        self.n_initial_points = n_initial_points
        self.acq_func = acq_func
        self.acq_optimizer = acq_optimizer
        self.random_state = random_state
        self.search_space = search_space
        self.skopt_optimizer = None
        self.create_optimizer()
        self.succeeded_trials = 0
        # List of recorded Trials names
        self.recorded_trials_names = []

    def create_optimizer(self):
        skopt_search_space = []

        for param in self.search_space.params:
            if param.type == INTEGER:
                skopt_search_space.append(
                    skopt.space.Integer(int(param.min), int(param.max), name=param.name)
                )
            elif param.type == DOUBLE:
                skopt_search_space.append(
                    skopt.space.Real(
                        float(param.min),
                        float(param.max),
                        "log-uniform",
                        name=param.name,
                    )
                )
            elif param.type == CATEGORICAL or param.type == DISCRETE:
                skopt_search_space.append(
                    skopt.space.Categorical(param.list, name=param.name)
                )

        self.skopt_optimizer = skopt.Optimizer(
            skopt_search_space,
            base_estimator=self.base_estimator,
            n_initial_points=self.n_initial_points,
            acq_func=self.acq_func,
            acq_optimizer=self.acq_optimizer,
            random_state=self.random_state,
        )

    def getSuggestions(self, trials, current_request_number):
        """
        Get the new suggested trials with skopt algorithm.
        """
        logger.info("-" * 100 + "\n")
        logger.info(
            "New GetSuggestions call with current request number: {}\n".format(
                current_request_number
            )
        )
        skopt_suggested = []
        loss_for_skopt = []
        if len(trials) > self.succeeded_trials or self.succeeded_trials == 0:
            self.succeeded_trials = len(trials)
            if self.succeeded_trials != 0:
                logger.info(
                    "Succeeded Trials changed: {}\n".format(self.succeeded_trials)
                )
            for trial in trials:
                if trial.name not in self.recorded_trials_names:
                    self.recorded_trials_names.append(trial.name)
                    trial_assignment = []
                    for param in self.search_space.params:
                        parameter_value = None
                        for assignment in trial.assignments:
                            if assignment.name == param.name:
                                parameter_value = assignment.value
                                break
                        if param.type == INTEGER:
                            trial_assignment.append(int(parameter_value))
                        elif param.type == DOUBLE:
                            trial_assignment.append(float(parameter_value))
                        else:
                            trial_assignment.append(parameter_value)
                    skopt_suggested.append(trial_assignment)
                    loss_value = float(trial.target_metric.value)
                    if self.search_space.goal == MAX_GOAL:
                        loss_value = -1 * loss_value
                    loss_for_skopt.append(loss_value)

            if loss_for_skopt != [] and skopt_suggested != []:
                logger.info("Running Optimizer tell to record observation")
                logger.info("Evaluated parameters: {}".format(skopt_suggested))
                logger.info("Objective values: {}\n".format(loss_for_skopt))
                t1 = datetime.datetime.now()
                self.skopt_optimizer.tell(skopt_suggested, loss_for_skopt)
                logger.info(
                    "Optimizer tell method takes {} seconds".format(
                        (datetime.datetime.now() - t1).seconds
                    )
                )
                logger.info(
                    "List of recorded Trials names: {}\n".format(
                        self.recorded_trials_names
                    )
                )

        else:
            logger.error(
                "Succeeded Trials didn't change: {}\n".format(self.succeeded_trials)
            )

        logger.info("Running Optimizer ask to query new parameters for Trials\n")

        return_trial_list = []

        skopt_suggested = self.skopt_optimizer.ask(n_points=current_request_number)
        for suggestion in skopt_suggested:
            logger.info("New suggested parameters for Trial: {}".format(suggestion))
            return_trial_list.append(
                BaseSkoptService.convert(self.search_space, suggestion)
            )

        logger.info(
            "GetSuggestions returns {} new Trials\n\n".format(len(return_trial_list))
        )
        return return_trial_list

    @staticmethod
    def convert(search_space, skopt_suggested):
        assignments = []
        for i in range(len(search_space.params)):
            param = search_space.params[i]
            if param.type == INTEGER:
                assignments.append(Assignment(param.name, skopt_suggested[i]))
            elif param.type == DOUBLE:
                assignments.append(Assignment(param.name, skopt_suggested[i]))
            elif param.type == CATEGORICAL or param.type == DISCRETE:
                assignments.append(Assignment(param.name, skopt_suggested[i]))
        return assignments
