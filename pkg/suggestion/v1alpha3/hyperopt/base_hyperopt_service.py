import hyperopt
import numpy as np
import logging

from pkg.suggestion.v1alpha3.internal.search_space import *
from pkg.suggestion.v1alpha3.internal.trial import *

logger = logging.getLogger("BaseHyperoptService")

class BaseHyperoptService(object):
    def __init__(self, algorithm_name="tpe", random_state=None):
        self.random_state = random_state
        if algorithm_name == 'tpe':
            self.hyperopt_algorithm = hyperopt.tpe.suggest
        elif algorithm_name == 'random':
            self.hyperopt_algorithm = hyperopt.rand.suggest
        # elif algorithm_name == 'hyperopt-anneal':
        #     self.hyperopt_algorithm = hyperopt.anneal.suggest_batch
        # elif algorithm_name == 'hyperopt-mix':
        #     self.hyperopt_algorithm = hyperopt.mix.suggest
        else:
            raise Exception('"Failed to create the algortihm: {}'.format(algorithm_name))

    def getSuggestions(self, search_space, trials, request_number):
        """
        Get the new suggested trials with the given algorithm.
        """
        # Construct search space, example: {"x": hyperopt.hp.uniform('x', -10, 10), "x2": hyperopt.hp.uniform('x2', -10, 10)}
        hyperopt_search_space = {}
        for param in search_space.params:
            if param.type == INTEGER:
                hyperopt_search_space[param.name] = hyperopt.hp.quniform(
                    param.name,
                    float(param.min),
                    float(param.max), 1)
            elif param.type == DOUBLE:
                hyperopt_search_space[param.name] = hyperopt.hp.uniform(
                    param.name,
                    float(param.min),
                    float(param.max))
            elif param.type == CATEGORICAL \
                    or param.type == DISCRETE:
                hyperopt_search_space[param.name] = hyperopt.hp.choice(
                    param.name, param.list)
         # New hyperopt variables
        hyperopt_rstate = np.random.RandomState(self.random_state)
        hyperopt_domain = hyperopt.Domain(
            None, hyperopt_search_space, pass_expr_memo_ctrl=None)

        hyperopt_trial_specs = []
        hyperopt_trial_results = []
        # Example: # Example: [{'tid': 0, 'idxs': {'l1_normalization': [0], 'learning_rate': [0], 'hidden2': [0], 'optimizer': [0]}, 'cmd': ('domain_attachment', 'FMinIter_Domain'), 'vals': {'l1_normalization': [0.1], 'learning_rate': [0.1], 'hidden2': [1], 'optimizer': [1]}, 'workdir': None}]
        hyperopt_trial_miscs = []
        hyperopt_trial_new_ids = []

        # Update hyperopt for trained trials with completed advisor trials
        completed_hyperopt_trials = hyperopt.Trials()
        for trial in trials:
             # Example: {'l1_normalization': [0], 'learning_rate': [0], 'hidden2': [0], 'optimizer': [0]}
            hyperopt_trial_miscs_idxs = {}
            # Example: {'l1_normalization': [0.1], 'learning_rate': [0.1], 'hidden2': [1], 'optimizer': [1]}
            hyperopt_trial_miscs_vals = {}
            new_id = trial.name
            hyperopt_trial_new_ids.append(new_id)
            hyperopt_trial_misc = dict(
                tid=new_id, cmd=hyperopt_domain.cmd, workdir=hyperopt_domain.workdir)
            for param in search_space.params:
                parameter_value = None
                for assignment in trial.assignments:
                    if assignment.name == param.name:
                        parameter_value = assignment.value
                        break
                if param.type == INTEGER:
                    hyperopt_trial_miscs_idxs[param.name] = [new_id]
                    hyperopt_trial_miscs_vals[param.name] = [
                        parameter_value]
                elif param.type == DOUBLE:
                    hyperopt_trial_miscs_idxs[param.name] = [new_id]
                    hyperopt_trial_miscs_vals[param.name] = [
                        parameter_value]
                elif param.type == DISCRETE or param.type == CATEGORICAL:
                    index_of_value_in_list = param.list.index(parameter_value)
                    hyperopt_trial_miscs_idxs[param.name] = [trial.name]
                    hyperopt_trial_miscs_vals[param.name] = [
                        index_of_value_in_list
                    ]

            hyperopt_trial_specs.append(None)

            hyperopt_trial_misc["idxs"] = hyperopt_trial_miscs_idxs
            hyperopt_trial_misc["vals"] = hyperopt_trial_miscs_vals
            hyperopt_trial_miscs.append(hyperopt_trial_misc)

            # TODO: Use negative objective value for loss or not
            objective_for_hyperopt = float(trial.target_metric.value)
            if search_space.goal == MAX_GOAL:
                # Now hyperopt only supports fmin and we need to reverse objective value for maximization
                objective_for_hyperopt = -1 * objective_for_hyperopt
            hyperopt_trial_result = {
                "loss": objective_for_hyperopt,
                "status": hyperopt.STATUS_OK
            }
            hyperopt_trial_results.append(hyperopt_trial_result)
        if len(trials) > 0:
            # Example: {'refresh_time': datetime.datetime(2018, 9, 18, 12, 6, 41, 922000), 'book_time': datetime.datetime(2018, 9, 18, 12, 6, 41, 922000), 'misc': {'tid': 0, 'idxs': {'x2': [0], 'x': [0]}, 'cmd': ('domain_attachment', 'FMinIter_Domain'), 'vals': {'x2': [-8.137088361136204], 'x': [-4.849028446711832]}, 'workdir': None}, 'state': 2, 'tid': 0, 'exp_key': None, 'version': 0, 'result': {'status': 'ok', 'loss': 14.849028446711833}, 'owner': None, 'spec': None}
            hyperopt_trials = completed_hyperopt_trials.new_trial_docs(
                hyperopt_trial_new_ids, hyperopt_trial_specs, hyperopt_trial_results, hyperopt_trial_miscs)
            for current_hyperopt_trials in hyperopt_trials:
                current_hyperopt_trials["state"] = hyperopt.JOB_STATE_DONE

            completed_hyperopt_trials.insert_trial_docs(hyperopt_trials)
            completed_hyperopt_trials.refresh()

        rval = hyperopt.FMinIter(
            self.hyperopt_algorithm,
            hyperopt_domain,
            completed_hyperopt_trials,
            max_evals=-1,
            rstate=hyperopt_rstate,
            verbose=0)
        rval.catch_eval_exceptions = False

        new_ids = rval.trials.new_trial_ids(request_number)

        rval.trials.refresh()

        random_state = rval.rstate.randint(2**31 - 1)
        new_trials = self.hyperopt_algorithm(
            new_ids, rval.domain, completed_hyperopt_trials, random_state)
        rval.trials.refresh()

        # Construct return advisor trials from new hyperopt trials
        list_of_assignments = []
        for i in range(request_number):
            vals = new_trials[i]['misc']['vals']
            list_of_assignments.append(BaseHyperoptService.convert(search_space, vals))
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
