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

import hyperopt
import numpy as np
import logging

from pkg.suggestion.v1beta1.internal.constant import INTEGER, DOUBLE, CATEGORICAL, DISCRETE, MAX_GOAL
from pkg.suggestion.v1beta1.internal.trial import Assignment

logger = logging.getLogger(__name__)

TPE_ALGORITHM_NAME = "tpe"
RANDOM_ALGORITHM_NAME = "random"


class BaseHyperoptService(object):
    def __init__(self,
                 algorithm_name=TPE_ALGORITHM_NAME,
                 algorithm_conf=None,
                 search_space=None):
        self.algorithm_name = algorithm_name
        self.algorithm_conf = algorithm_conf or {}
        # pop common configurations
        random_state = self.algorithm_conf.pop('random_state', None)

        if self.algorithm_name == TPE_ALGORITHM_NAME:
            self.hyperopt_algorithm = hyperopt.tpe.suggest
        elif self.algorithm_name == RANDOM_ALGORITHM_NAME:
            self.hyperopt_algorithm = hyperopt.rand.suggest
        # elif algorithm_name == 'hyperopt-anneal':
        #     self.hyperopt_algorithm = hyperopt.anneal.suggest_batch
        # elif algorithm_name == 'hyperopt-mix':
        #     self.hyperopt_algorithm = hyperopt.mix.suggest

        self.search_space = search_space
        # New hyperopt variables
        self.hyperopt_rstate = np.random.RandomState(random_state)
        self.create_hyperopt_domain()
        self.create_fmin()
        self.is_first_run = True

    def create_hyperopt_domain(self):
        # Construct search space, example: {"x": hyperopt.hp.uniform('x', -10, 10), "x2": hyperopt.hp.uniform('x2', -10, 10)}
        hyperopt_search_space = {}
        for param in self.search_space.params:
            if param.type == INTEGER:
                hyperopt_search_space[param.name] = hyperopt.hp.quniform(
                    param.name,
                    float(param.min),
                    float(param.max),
                    float(param.step))
            elif param.type == DOUBLE:
                hyperopt_search_space[param.name] = hyperopt.hp.uniform(
                    param.name,
                    float(param.min),
                    float(param.max))
            elif param.type == CATEGORICAL or param.type == DISCRETE:
                hyperopt_search_space[param.name] = hyperopt.hp.choice(
                    param.name, param.list)

        self.hyperopt_domain = hyperopt.Domain(
            None, hyperopt_search_space, pass_expr_memo_ctrl=None)

    def create_fmin(self):
        self.fmin = hyperopt.FMinIter(
            self.hyperopt_algorithm,
            self.hyperopt_domain,
            trials=hyperopt.Trials(),
            max_evals=-1,
            rstate=self.hyperopt_rstate,
            verbose=False)

        self.fmin.catch_eval_exceptions = False

    def getSuggestions(self, trials, current_request_number):
        """
        Get the new suggested trials with the given algorithm.
        """

        recorded_trials_names = self.fmin.trials.specs

        hyperopt_trial_new_ids = []
        hyperopt_trial_specs = []
        hyperopt_trial_results = []
        hyperopt_trial_miscs = []
        # Update hyperopt FMin with new completed Trials
        for trial in trials:
            if {"trial-name": trial.name} not in recorded_trials_names:
                # Produce new id for the new Trial
                new_id = self.fmin.trials.new_trial_ids(1)
                hyperopt_trial_new_ids.append(new_id[0])
                hyperopt_trial_miscs_idxs = {}
                # Example: {'l1_normalization': [0.1], 'learning_rate': [0.1], 'hidden2': [1], 'optimizer': [1]}
                hyperopt_trial_miscs_vals = {}

                # Insert Trial assignment to the misc
                hyperopt_trial_misc = dict(
                    tid=new_id[0], cmd=self.hyperopt_domain.cmd, workdir=self.hyperopt_domain.workdir)
                for param in self.search_space.params:
                    parameter_value = None
                    for assignment in trial.assignments:
                        if assignment.name == param.name:
                            parameter_value = assignment.value
                            break
                    if param.type == INTEGER:
                        hyperopt_trial_miscs_idxs[param.name] = new_id
                        hyperopt_trial_miscs_vals[param.name] = [int(parameter_value)]
                    elif param.type == DOUBLE:
                        hyperopt_trial_miscs_idxs[param.name] = new_id
                        hyperopt_trial_miscs_vals[param.name] = [float(parameter_value)]
                    elif param.type == DISCRETE or param.type == CATEGORICAL:
                        index_of_value_in_list = param.list.index(parameter_value)
                        hyperopt_trial_miscs_idxs[param.name] = new_id
                        hyperopt_trial_miscs_vals[param.name] = [index_of_value_in_list]

                hyperopt_trial_misc["idxs"] = hyperopt_trial_miscs_idxs
                hyperopt_trial_misc["vals"] = hyperopt_trial_miscs_vals
                hyperopt_trial_miscs.append(hyperopt_trial_misc)

                # Insert Trial name to the spec
                hyperopt_trial_spec = {
                    "trial-name": trial.name
                }
                hyperopt_trial_specs.append(hyperopt_trial_spec)

                # Insert Trial result to the result
                # TODO: Use negative objective value for loss or not
                # TODO: Do we need to analyse additional_metrics?
                objective_for_hyperopt = float(trial.target_metric.value)
                if self.search_space.goal == MAX_GOAL:
                    # Now hyperopt only supports fmin and we need to reverse objective value for maximization
                    objective_for_hyperopt = -1 * objective_for_hyperopt
                hyperopt_trial_result = {
                    "loss": objective_for_hyperopt,
                    "status": hyperopt.STATUS_OK
                }
                hyperopt_trial_results.append(hyperopt_trial_result)

        if len(trials) > 0:

            # Create new Trial doc
            hyperopt_trials = hyperopt.Trials().new_trial_docs(
                tids=hyperopt_trial_new_ids,
                specs=hyperopt_trial_specs,
                results=hyperopt_trial_results,
                miscs=hyperopt_trial_miscs)

            for i, _ in enumerate(hyperopt_trials):
                hyperopt_trials[i]["state"] = hyperopt.JOB_STATE_DONE

            # Insert new set of Trial to FMin object
            # Example: of inserting doc with tunning lr
            #   [{
            #       'state':2,
            #       'tid':5,
            #       'spec':{
            #          'trial-name':'tpe-48xl8whg'
            #       },
            #       'result':{
            #          'loss':-0.1135,
            #          'status':'ok'
            #       },
            #       'misc':{
            #          'tid':5,
            #          'cmd':('domain_attachment','FMinIter_Domain'),
            #          'workdir':None,
            #          'idxs':{
            #             '--lr':[5]
            #          },
            #          'vals':{
            #             '--lr':[0.025351232898626827]
            #          }
            #       },
            #       'exp_key':None,
            #       'owner':None,
            #       'version':0,
            #       'book_time':None,
            #       'refresh_time':None
            #   }]
            self.fmin.trials.insert_trial_docs(hyperopt_trials)
            self.fmin.trials.refresh()

        # Produce new current_request_number ids to make new Suggestion
        hyperopt_trial_new_ids = self.fmin.trials.new_trial_ids(current_request_number)
        random_state = self.fmin.rstate.randint(2**31 - 1)

        # Trial list that must be deployed
        new_trials = []
        if self.algorithm_name == RANDOM_ALGORITHM_NAME:
            new_trials = self.hyperopt_algorithm(
                new_ids=hyperopt_trial_new_ids,
                domain=self.fmin.domain,
                trials=self.fmin.trials,
                seed=random_state)
        elif self.algorithm_name == TPE_ALGORITHM_NAME:
            # n_startup_jobs indicates for how many Trials we run random suggestion
            # This must be current_request_number value
            # After this tpe suggestion starts analyse Trial info.
            # On the first run we can run suggest just once with n_startup_jobs
            # Next suggest runs must be for each new Trial generation
            if self.is_first_run:
                new_trials = self.hyperopt_algorithm(
                    new_ids=hyperopt_trial_new_ids,
                    domain=self.fmin.domain,
                    trials=self.fmin.trials,
                    seed=random_state,
                    n_startup_jobs=current_request_number,
                    **self.algorithm_conf)
                self.is_first_run = False
            else:
                for i in range(current_request_number):
                    # hyperopt_algorithm always returns one new Trial
                    new_trials.append(self.hyperopt_algorithm(
                        new_ids=[hyperopt_trial_new_ids[i]],
                        domain=self.fmin.domain,
                        trials=self.fmin.trials,
                        seed=random_state,
                        n_startup_jobs=current_request_number,
                        **self.algorithm_conf)[0])

        # Construct return advisor Trials from new hyperopt Trials
        list_of_assignments = []
        for trial in new_trials:
            vals = trial['misc']['vals']
            list_of_assignments.append(BaseHyperoptService.convert(self.search_space, vals))

        if len(list_of_assignments) > 0:
            logger.info("GetSuggestions returns {} new Trial\n".format(len(new_trials)))

        return list_of_assignments

    @staticmethod
    def convert(search_space, vals):
        assignments = []
        for param in search_space.params:
            if param.type == INTEGER:
                assignments.append(Assignment(param.name, int(vals[param.name][0])))
            elif param.type == DOUBLE:
                assignments.append(Assignment(param.name, vals[param.name][0]))
            elif param.type == CATEGORICAL or param.type == DISCRETE:
                assignments.append(
                    Assignment(param.name, param.list[vals[param.name][0]]))
        return assignments
