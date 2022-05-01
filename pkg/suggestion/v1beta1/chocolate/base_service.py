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

import chocolate as choco
import logging
import base64
import warnings

from pkg.suggestion.v1beta1.internal.constant import MAX_GOAL, INTEGER, DOUBLE, CATEGORICAL, DISCRETE
from pkg.suggestion.v1beta1.internal.trial import Assignment

logger = logging.getLogger(__name__)

DB_ADDRESS = "sqlite:///my_db.db?check_same_thread=False"
DB_FIELD_LOSS = "_loss"
DB_FIELD_CHOCOLATE_ID = "_chocolate_id"
DB_FIELD_TRIAL_NAME = "_trial_name"

DEPRECATED_ALGORITHM_NAME = {
    "chocolate-random": "random",
    "chocolate-quasirandom": "quasirandom",
    "chocolate-bayesian-optimization": "bayesianoptimization",
    "chocolate-mocmaes": "mocmaes",
}


class BaseChocolateService(object):
    """
    Refer to https://chocolate.readthedocs.io/
    """

    def __init__(self, algorithm_name, search_space):
        self.conn = choco.SQLiteConnection(DB_ADDRESS)
        self.search_space = search_space
        self.chocolate_optimizer = None
        self.create_optimizer(algorithm_name)
        # created_trials is the list of dicts with all created trials assignments, loss and trial name
        # _chocolate_id is the ID of the trial, Assignment names are encoded,
        # _loss is the target metric, _trial_name is the Trial name
        # One row example:
        # {'_chocolate_id': 0, 'LS1scg==': 0.001, 'LS1udW0tZXBvY2hz': 1, 'LS1udW0tbGF5ZXJz': 2,
        #   "_loss": "0.97", "_trial_name": "grid-hsdvfdwl"}
        self.created_trials = []
        self.recorded_trials_names = []

    def create_optimizer(self, algorithm_name):

        # Search Space example: {"x" : choco.uniform(-6, 6), "y" : choco.uniform(-6, 6)}
        chocolate_search_space = {}

        for param in self.search_space.params:
            key = BaseChocolateService.encode(param.name)
            # Chocolate quantized_uniform distribution uses half-open interval: [low, high).
            if param.type == INTEGER:
                chocolate_search_space[key] = choco.quantized_uniform(
                    int(param.min), int(param.max) + int(param.step), int(param.step))
            elif param.type == DOUBLE:
                chocolate_search_space[key] = choco.quantized_uniform(
                    float(param.min), float(param.max) + float(param.step), float(param.step))
            # For Categorical and Discrete insert indexes to DB from list of values
            elif param.type == CATEGORICAL or param.type == DISCRETE:
                chocolate_search_space[key] = choco.choice(
                    [idx for idx, _ in enumerate(param.list)])

        if algorithm_name in DEPRECATED_ALGORITHM_NAME:
            warnings.warn(
                "Algorithm name '{}' is deprecated. Please use '{}'.".format(
                    algorithm_name, DEPRECATED_ALGORITHM_NAME[algorithm_name],
                ),
                DeprecationWarning,
            )
            algorithm_name = DEPRECATED_ALGORITHM_NAME[algorithm_name]

        # Refer to https://chocolate.readthedocs.io/tutorials/algo.html
        if algorithm_name == "grid":
            self.chocolate_optimizer = choco.Grid(
                self.conn, chocolate_search_space, clear_db=True)
        # hyperopt-random is the default option in katib.
        elif algorithm_name == "random":
            self.chocolate_optimizer = choco.Random(
                self.conn, chocolate_search_space, clear_db=True)
        elif algorithm_name == "quasirandom":
            self.chocolate_optimizer = choco.QuasiRandom(
                self.conn, chocolate_search_space, clear_db=True)
        elif algorithm_name == "bayesianoptimization":
            self.chocolate_optimizer = choco.Bayes(
                self.conn, chocolate_search_space, clear_db=True)
        # elif self.algorithm_name == "chocolate-CMAES":
        #     self.chocolate_optimizer = choco.CMAES(self.conn, chocolate_search_space, clear_db=True)
        elif algorithm_name == "mocmaes":
            mu = 1
            self.chocolate_optimizer = choco.MOCMAES(
                self.conn, chocolate_search_space, mu=mu, clear_db=True)
        else:
            raise Exception(
                '"Failed to create Chocolate optimizer for the algorithm: {}'.format(algorithm_name))

    def getSuggestions(self, trials, current_request_number, total_request_number):
        """
        Get the new suggested trials with chocolate algorithm.
        """
        logger.info("-" * 100 + "\n")
        logger.info("New GetSuggestions call with total requested {} and currently requesting {} \n".format(
            total_request_number, current_request_number))
        for _, trial in enumerate(trials):
            if trial.name not in self.recorded_trials_names:
                loss_for_choco = float(trial.target_metric.value)
                if self.search_space.goal == MAX_GOAL:
                    loss_for_choco = -1 * loss_for_choco

                trial_assignments_dict = {}
                for param in self.search_space.params:
                    param_assignment = None
                    for assignment in trial.assignments:
                        if param.name == assignment.name:
                            param_assignment = assignment.value
                            break
                    if param.type == INTEGER:
                        param_assignment = int(param_assignment)
                    elif param.type == DOUBLE:
                        param_assignment = float(param_assignment)
                    elif param.type == CATEGORICAL or param.type == DISCRETE:
                        param_assignment = param.list.index(param_assignment)
                    trial_assignments_dict.update({BaseChocolateService.encode(
                        param.name): param_assignment})

                # Finding index for the current Trial Assignments in created_trial list without loss
                new_trial_loss_idx = -1
                i = 0
                while new_trial_loss_idx == -1 and i < len(self.created_trials):
                    # Created Trial must not include loss and must have the same param assignment
                    if ((DB_FIELD_LOSS not in self.created_trials[i] or
                         self.created_trials[i][DB_FIELD_LOSS] is None) and
                            len(trial_assignments_dict.items() & self.created_trials[i].items()) ==
                            len(self.search_space.params)):
                        new_trial_loss_idx = i
                    i += 1

                if new_trial_loss_idx != -1:
                    self.created_trials[new_trial_loss_idx][DB_FIELD_LOSS] = loss_for_choco
                    self.created_trials[new_trial_loss_idx][DB_FIELD_TRIAL_NAME] = trial.name

                    # Update sqlite database with new loss and trial assignments
                    id_filter = {
                        DB_FIELD_CHOCOLATE_ID: self.created_trials[new_trial_loss_idx][DB_FIELD_CHOCOLATE_ID]}
                    self.conn.update_result(
                        id_filter,
                        self.created_trials[new_trial_loss_idx])

                    self.recorded_trials_names.append(trial.name)

                    logger.info("New record in sqlite DB is updated")
                    logger.info("{}\n".format(
                        self.created_trials[new_trial_loss_idx]))
        # Assuming that created_trials are already populated
        # TODO: Handle Restart of algorithm pod
        logger.info("{} Trials created in DB".format(len(self.created_trials)))
        if total_request_number != len(self.created_trials) + current_request_number:
            logger.info("Mismatch in generated trials with k8s suggestions trials")
        new_actual_requested_no = total_request_number - len(self.created_trials)
        prev_generated_no = current_request_number - new_actual_requested_no
        logger.info(
            "In this call, New {} Trials will be generated, {} Trials will be reused from previously generated".format(
                new_actual_requested_no, prev_generated_no))

        list_of_assignments = []
        if prev_generated_no > 0:
            for params in self.created_trials[-prev_generated_no:]:
                if DB_FIELD_TRIAL_NAME in params:
                    logger.error("Trial already updated in selected assignment {}".format(params))
                new_assignment = BaseChocolateService.convert(
                    self.search_space, params)
                list_of_assignments.append(new_assignment)

        for i in range(new_actual_requested_no):
            try:
                token, chocolate_params = self.chocolate_optimizer.next()
                new_assignment = BaseChocolateService.convert(
                    self.search_space, chocolate_params)
                list_of_assignments.append(new_assignment)
                logger.info("New suggested parameters for Trial with chocolate_id: {}".format(
                    token[DB_FIELD_CHOCOLATE_ID]))
                for assignment in new_assignment:
                    logger.info("Name = {}, Value = {}".format(
                        assignment.name, assignment.value))
                logger.info("-" * 50 + "\n")
                # Add new trial assignment with chocolate_id to created trials
                token.update(chocolate_params)
                new_trial_dict = token
                self.created_trials.append(new_trial_dict)

            except StopIteration:
                logger.info(
                    "Chocolate db is exhausted, increase Search Space or decrease maxTrialCount!")

        if len(list_of_assignments) > 0:
            logger.info("GetSuggestions returns {} Trials from requested {} Trials\n\n".format(
                len(list_of_assignments), current_request_number))

        return list_of_assignments

    @staticmethod
    def convert(search_space, chocolate_params):
        assignments = []
        for param in search_space.params:
            key = BaseChocolateService.encode(param.name)
            if param.type == INTEGER:
                assignments.append(Assignment(
                    param.name, chocolate_params[key]))
            elif param.type == DOUBLE:
                assignments.append(Assignment(
                    param.name, chocolate_params[key]))
            elif param.type == CATEGORICAL or param.type == DISCRETE:
                assignments.append(Assignment(
                    param.name, param.list[chocolate_params[key]]))
        return assignments

    @staticmethod
    def encode(name):
        """Encode the name. Chocolate will check if the name contains hyphens.
        Thus we need to encode it.
        """
        return base64.b64encode(name.encode('utf-8')).decode('utf-8')
