import json
import numpy as np
import skopt
import logging

from pkg.suggestion.v1alpha3.internal.search_space import *
from pkg.suggestion.v1alpha3.internal.trial import *

logger = logging.getLogger("BaseSkoptService")


class BaseSkoptService(object):
    """
    Refer to https://github.com/scikit-optimize/scikit-optimize .
    """

    def __init__(self,
                 base_estimator="GP",
                 n_initial_points=10,
                 acq_func="gp_hedge",
                 acq_optimizer="auto",
                 random_state=None,
                 search_space=None):
        self.base_estimator = base_estimator
        self.n_initial_points = n_initial_points
        self.acq_func = acq_func
        self.acq_optimizer = acq_optimizer
        self.random_state = random_state
        self.search_space = search_space
        self.skopt_optimizer = None
        self.create_optimizer()
        self.succeeded_trials = 0
        # Dict of recorded trials where key = loss value, value = List of trial assignment list for this value
        self.recorded_trials = {}

    def create_optimizer(self):
        skopt_search_space = []

        for param in self.search_space.params:
            if param.type == INTEGER:
                skopt_search_space.append(skopt.space.Integer(
                    int(param.min), int(param.max), name=param.name))
            elif param.type == DOUBLE:
                skopt_search_space.append(skopt.space.Real(
                    float(param.min), float(param.max), "log-uniform", name=param.name))
            elif param.type == CATEGORICAL or param.type == DISCRETE:
                skopt_search_space.append(
                    skopt.space.Categorical(param.list, name=param.name))

        self.skopt_optimizer = skopt.Optimizer(
            skopt_search_space,
            base_estimator=self.base_estimator,
            n_initial_points=self.n_initial_points,
            acq_func=self.acq_func,
            acq_optimizer=self.acq_optimizer,
            random_state=self.random_state)

    def getSuggestions(self, trials, request_number):
        """
        Get the new suggested trials with skopt algorithm.
        """

        skopt_suggested = []
        loss_for_skopt = []
        return_trial_list = []
        if len(trials) > self.succeeded_trials or self.succeeded_trials == 0:
            self.succeeded_trials = len(trials)
            logger.info("Succeeded Trials: {}".format(self.succeeded_trials))
            for trial in trials:
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

                loss_value = float(trial.target_metric.value)
                if self.search_space.goal == MAX_GOAL:
                    loss_value = -1 * loss_value

                # Update only for new values
                if loss_value not in self.recorded_trials:
                    self.recorded_trials[loss_value] = []
                    self.recorded_trials[loss_value].append(trial_assignment)
                    skopt_suggested.append(trial_assignment)
                    loss_for_skopt.append(loss_value)
                elif trial_assignment not in self.recorded_trials[loss_value]:
                    self.recorded_trials[loss_value].append(trial_assignment)
                    skopt_suggested.append(trial_assignment)
                    loss_for_skopt.append(loss_value)

            if loss_for_skopt != [] and skopt_suggested != []:
                self.skopt_optimizer.tell(skopt_suggested, loss_for_skopt)

            for i in range(request_number):
                skopt_suggested = self.skopt_optimizer.ask()
                return_trial_list.append(
                    BaseSkoptService.convert(self.search_space, skopt_suggested))

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
        for a in assignments:
            logger.info("Generate new Trial with Assignment")
            logger.info(a)
        return assignments
