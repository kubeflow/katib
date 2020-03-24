import json
import numpy as np
import chocolate as choco
import logging
import base64

from pkg.suggestion.v1alpha3.internal.search_space import *
from pkg.suggestion.v1alpha3.internal.trial import *

logger = logging.getLogger(__name__)


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
            key = BaseChocolateService.encode(param.name)
            if param.type == INTEGER:
                chocolate_search_space[key] = choco.quantized_uniform(
                    int(param.min), int(param.max), 1)
            elif param.type == DOUBLE:
                chocolate_search_space[key] = choco.quantized_uniform(
                    float(param.min), float(param.max), float(param.step))
            elif param.type == CATEGORICAL:
                chocolate_search_space[key] = choco.choice(param.list)
            else:
                chocolate_search_space[key] = choco.choice(
                    [float(e) for e in param.list])

        conn = choco.SQLiteConnection("sqlite:///my_db.db")
        # Refer to https://chocolate.readthedocs.io/tutorials/algo.html
        if self.algorithm_name == "grid":
            sampler = choco.Grid(conn, chocolate_search_space, clear_db=True)
        # hyperopt-random is the default option in katib.
        elif self.algorithm_name == "chocolate-random":
            sampler = choco.Random(conn, chocolate_search_space, clear_db=True)
        elif self.algorithm_name == "chocolate-quasirandom":
            sampler = choco.QuasiRandom(
                conn, chocolate_search_space, clear_db=True)
        elif self.algorithm_name == "chocolate-bayesian-optimization":
            sampler = choco.Bayes(conn, chocolate_search_space, clear_db=True)
        # elif self.algorithm_name == "chocolate-CMAES":
        #     sampler = choco.CMAES(conn, chocolate_search_space, clear_db=True)
        elif self.algorithm_name == "chocolate-MOCMAES":
            mu = 1
            sampler = choco.MOCMAES(
                conn, chocolate_search_space, mu=mu, clear_db=True)
        else:
            raise Exception(
                '"Failed to create the algortihm: {}'.format(self.algorithm_name))

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
                entry.update({BaseChocolateService.encode(
                    param.name): param_assignment})
            logger.info(entry)
            # Should not use sampler.update(token, loss), because we will create
            # a new BaseChocolateService instance for every request. Thus we need
            # to insert all previous trials every time.
            conn.insert_result(entry)

        list_of_assignments = []

        for i in range(request_number):
            try:
                token, chocolate_params = sampler.next()
                list_of_assignments.append(
                    BaseChocolateService.convert(search_space, chocolate_params))
            except StopIteration:
                logger.info("Chocolate db is exhausted, increase Search Space or decrease maxTrialCount!")
        return list_of_assignments

    @staticmethod
    def convert(search_space, chocolate_params):
        assignments = []
        for i in range(len(search_space.params)):
            param = search_space.params[i]
            key = BaseChocolateService.encode(param.name)
            if param.type == INTEGER:
                assignments.append(Assignment(
                    param.name, chocolate_params[key]))
            elif param.type == DOUBLE:
                assignments.append(Assignment(
                    param.name, chocolate_params[key]))
            elif param.type == CATEGORICAL or param.type == DISCRETE:
                assignments.append(Assignment(
                    param.name, chocolate_params[key]))
        return assignments

    @staticmethod
    def encode(name):
        """Encode the name. Chocolate will check if the name contains hyphens.
        Thus we need to encode it.
        """
        return base64.b64encode(name.encode('utf-8')).decode('utf-8')
