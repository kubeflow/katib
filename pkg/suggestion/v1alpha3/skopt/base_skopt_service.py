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

    def __init__(self, algorithm_name="skopt-bayesian-optimization",
                 base_estimator="GP",
                 n_initial_points=10,
                 acq_func="gp_hedge",
                 acq_optimizer="auto",
                 random_state=None):
        self.base_estimator = base_estimator
        self.n_initial_points = n_initial_points
        self.acq_func = acq_func
        self.acq_optimizer = acq_optimizer
        self.random_state = random_state
        self.algorithm_name = algorithm_name

    def getSuggestions(self, search_space, trials, request_number):
        """
        Get the new suggested trials with skopt algorithm.
        """

        skopt_search_space = []

        for param in search_space.params:
            if param.type == INTEGER:
                skopt_search_space.append(skopt.space.Integer(
                    int(param.min), int(param.max), name=param.name))
            elif param.type == DOUBLE:
                skopt_search_space.append(skopt.space.Real(
                    float(param.min), float(param.max), "log-uniform", name=param.name))
            elif param.type == CATEGORICAL or param.type == DISCRETE:
                skopt_search_space.append(
                    skopt.space.Categorical(param.list, name=param.name))

        if self.algorithm_name != "bayesianoptimization":
            raise Exception(
                '"Failed to create the algortihm: {}'.format(self.algorithm_name))
        skopt_optimizer = skopt.Optimizer(skopt_search_space,
                                          base_estimator=self.base_estimator,
                                          n_initial_points=self.n_initial_points,
                                          acq_func=self.acq_func,
                                          acq_optimizer=self.acq_optimizer,
                                          random_state=self.random_state)

        for trial in trials:
            skopt_suggested = []
            for param in search_space.params:
                parameter_value = None
                for assignment in trial.assignments:
                    if assignment.name == param.name:
                        parameter_value = assignment.value
                        break
                if param.type == INTEGER:
                    skopt_suggested.append(int(parameter_value))
                elif param.type == DOUBLE:
                    skopt_suggested.append(float(parameter_value))
                else:
                    skopt_suggested.append(parameter_value)

            loss_for_skopt = float(trial.target_metric.value)
            if search_space.goal == MAX_GOAL:
                loss_for_skopt = -1 * loss_for_skopt

            skopt_optimizer.tell(skopt_suggested, loss_for_skopt)

        return_trial_list = []

        for i in range(request_number):
            skopt_suggested = skopt_optimizer.ask()
            return_trial_list.append(
                BaseSkoptService.convert(search_space, skopt_suggested))
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
