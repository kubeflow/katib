import json
import numpy as np
import chocolate as choco
import logging

from pkg.api.v1alpha2.python import api_pb2
from .internal.search_space import *
from .internal.trial import *

logger = logging.getLogger("BaseChocolateService")


class BaseChocolateService(object):
    """
    Refer to https://chocolate.readthedocs.io/
    """

    def __init__(self, algorithm_name=""):
        self.algorithm_name = algorithm_name

    def getSuggestions(self, search_space, trials, request_number):
        """
        Get the new suggested trials with chocolate algorithm.
        """

        # Example: {"x" : choco.uniform(-6, 6), "y" : choco.uniform(-6, 6)}
        chocolate_search_space = {}

        for param in search_space.params:
            if param.type == INTEGER:
                chocolate_search_space[param.name] = choco.quantized_uniform(
                    int(param.min), int(param.max), 1)
            elif param.type == DOUBLE:
                if param.step != None:
                    chocolate_search_space[param.name] = choco.quantized_uniform(
                        float(param.min), float(param.max), float(param.step))
                else:
                    chocolate_search_space[param.name] = choco.uniform(
                        float(param.min), float(param.max))
            elif param.type == CATEGORICAL or param.type == DISCRETE:
                chocolate_search_space[param.name] = choco.choice(param.list)

        conn = choco.SQLiteConnection("sqlite:///my_db.db")
        # Refer to https://chocolate.readthedocs.io/tutorials/algo.html
        if self.algorithm_name == "chocolate-grid":
            sampler = choco.Grid(conn, chocolate_search_space, clear_db=True)
        elif self.algorithm_name == "chocolate-random":
            sampler = choco.Random(conn, chocolate_search_space, clear_db=True)
        elif self.algorithm_name == "chocolate-quasirandom":
            sampler = choco.QuasiRandom(
                conn, chocolate_search_space, clear_db=True)
        elif self.algorithm_name == "chocolate-bayesian-optimization":
            sampler = choco.Bayes(conn, chocolate_search_space, clear_db=True)
        elif self.algorithm_name == "chocolate-CMAES":
            sampler = choco.CMAES(conn, chocolate_search_space, clear_db=True)
        elif self.algorithm_name == "chocolate-MOCMAES":
            mu = 1
            sampler = choco.MOCMAES(conn, chocolate_search_space, mu=mu, clear_db=True)

        for index, trial in enumerate(trials):
            loss_for_choco = float(trial.target_metric.value)
            if search_space.goal == MAX_GOAL:
                loss_for_choco = -1 * loss_for_choco

            entry = {"_chocolate_id": index, "_loss": loss_for_choco}
            for param in search_space.params:
                param_assignment = None
                for assignment in trial.assignments:
                    if param.name == assignment.name:
                        param_assignment = assignment.value
                        break
                if param.type == INTEGER:
                    param_assignment = int(param_assignment)
                elif param.type == DOUBLE:
                    param_assignment = float(param_assignment)
                entry.update({param.name: param_assignment})
            logger.info(entry)
            # Should not use sampler.update(token, loss), because we will create 
            # a new BaseChocolateService instance for every request. Thus we need
            # to insert all previous trials every time.
            conn.insert_result(entry)

        return_trial_list = []

        for i in range(request_number):
            token, chocolate_params = sampler.next()
            return_trial_list.append(
                BaseChocolateService.convert(search_space, chocolate_params))
        return return_trial_list

    @staticmethod
    def convert(search_space, chocolate_params):
        assignments = []
        for i in range(len(search_space.params)):
            param = search_space.params[i]
            if param.type == INTEGER:
                assignments.append(Assignment(param.name, chocolate_params[param.name]))
            elif param.type == DOUBLE:
                assignments.append(Assignment(param.name, chocolate_params[param.name]))
            elif param.type == CATEGORICAL or param.type == DISCRETE:
                assignments.append(Assignment(param.name, chocolate_params[param.name]))
        return Trial(None, assignments, None, None, None)
