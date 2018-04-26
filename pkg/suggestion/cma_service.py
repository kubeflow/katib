import random
import string

import grpc
import numpy as np

from pkg.api.python import api_pb2
from pkg.api.python import api_pb2_grpc
from pkg.suggestion.cma.src.algorithm_manager import AlgorithmManager
from pkg.suggestion.cma.src.cma_algorithm import CMAES

_ONE_DAY_IN_SECONDS = 60 * 60 * 24


class CMAService(api_pb2_grpc.SuggestionServicer):
    def __init__(self):
        """
        format for self.population:
        {
            study_id:{
                internal_params:{
                    path_sigma:
                    path_c
                    C
                    sigma
                    mean
                }
                population:[
                {
                    trial_id:
                    metric:
                    parameters: []
                    penalty:
                }
                ]
            }
        }
        """
        self.population = {}

    def GenerateTrials(self, request, context):
        trials = []
        X_train = []
        y_train = []
        algo_manager = AlgorithmManager(
            study_id=request.study_id,
            study_config=request.configs,
            X_train=X_train,
            y_train=y_train,
        )
        lowerbound = np.array(algo_manager.lower_bound)
        upperbound = np.array(algo_manager.upper_bound)

        cma = CMAES(
            dim=algo_manager.dim,
            upperbound=upperbound,
            lowerbound=lowerbound,
        )
        # get suggestion for this study for the first time
        if request.study_id not in self.population.keys():
            self.population[request.study_id] = {}
            self.population[request.study_id]["population"] = []
            self.population[request.study_id]["internal_params"] = {}

            path_sigma, path_c, C, sigma, mean = cma.init_params()
            self.population[request.study_id]["internal_params"]["path_sigma"] = path_sigma
            self.population[request.study_id]["internal_params"]["path_c"] = path_c
            self.population[request.study_id]["internal_params"]["C"] = C
            self.population[request.study_id]["internal_params"]["sigma"] = sigma
            self.population[request.study_id]["internal_params"]["mean"] = mean

        # this study already have a population to try
        else:
            for trial in request.completed_trials:
                for p in self.population[request.study_id]["population"]:
                    if trial.trial_id == p["trial_id"]:
                        # the algorithm is originally for minimization
                        if request.configs.optimization_type == api_pb2.MAXIMIZE:
                            p["metric"] = -float(trial.objective_value)
                        else:
                            p["metric"] = float(trial.objective_value)
            """
            metrics
            [
                {
                    "x": [],
                    "y": ,
                    "penalty":
                }
            ]
            """
            # the algorithm cannot continue without all trials in the population are evaluated
            metrics = []
            for p in self.population[request.study_id]["population"]:
                if p["metric"] is None:
                    context.set_code(grpc.StatusCode.UNKNOWN)
                    context.set_details("all trials in the population should be evaluated")
                    return api_pb2.GenerateTrialsReply(
                        trials=[],
                        completed=False,
                    )
                metrics.append(dict(
                    x=p["parameters"],
                    y=p["metric"],
                    penalty=p["penalty"],
                ))
            next_path_sigma, next_path_c, next_C, next_sigma, next_mean = cma.report_metric(
                objective_dict=metrics,
                mean=self.population[request.study_id]["internal_params"]["mean"],
                sigma=self.population[request.study_id]["internal_params"]["sigma"],
                C=self.population[request.study_id]["internal_params"]["C"],
                path_sigma=self.population[request.study_id]["internal_params"]["path_sigma"],
                path_c=self.population[request.study_id]["internal_params"]["path_c"],
            )
            self.population[request.study_id]["internal_params"]["path_sigma"] = next_path_sigma
            self.population[request.study_id]["internal_params"]["path_c"] = next_path_c
            self.population[request.study_id]["internal_params"]["C"] = next_C
            self.population[request.study_id]["internal_params"]["sigma"] = next_sigma
            self.population[request.study_id]["internal_params"]["mean"] = next_mean

        """
        raw_suggestions:
        [
            {
                "suggestion":[]
                "penalty":
            }
        ]
        """
        raw_suggestions = cma.get_suggestion(
            mean=self.population[request.study_id]["internal_params"]["mean"],
            sigma=self.population[request.study_id]["internal_params"]["sigma"],
            C=self.population[request.study_id]["internal_params"]["C"],
        )

        for raw_suggestion in raw_suggestions:
            # record the intermediate step
            trial_id = ''.join(random.sample(string.ascii_letters + string.digits, 12))
            self.population[request.study_id]["population"].append(dict(
                trial_id=trial_id,
                metric=None,
                parameters=raw_suggestion["suggestion"],
                penalty=raw_suggestion["penalty"],
            ))

            # parse the raw suggestions to desired format
            trial = algo_manager.parse_x_next(raw_suggestion["suggestion"])
            trial = algo_manager.convert_to_dict(trial)
            trials.append(api_pb2.Trial(
                trial_id=trial_id,
                study_id=request.study_id,
                parameter_set=[
                    api_pb2.Parameter(
                        name=x["name"],
                        value=str(x["value"]),
                        parameter_type=x["type"],
                    ) for x in trial
                ],
                status=api_pb2.PENDING,
                eval_logs=[],
            ))

        return api_pb2.GenerateTrialsReply(
            trials=trials,
            completed=False,
        )

    def StopSuggestion(self, request, context):
        if request.study_id in self.population.keys():
            del self.population[request.study_id]
        return api_pb2.StopStudyReply()

    def SetSuggestionParameters(self, request, context):
        return api_pb2.SetSuggestionParametersReply()
