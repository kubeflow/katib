import json

import grpc
import numpy as np

from pkg.api.python import api_pb2
from pkg.api.python import api_pb2_grpc
from pkg.suggestion.cma.src.algorithm_manager import AlgorithmManager
from pkg.suggestion.cma.src.cma_algorithm import CMAES

_ONE_DAY_IN_SECONDS = 60 * 60 * 24


class CMAService(api_pb2_grpc.SuggestionServicer):
    def __init__(self):
        channel = grpc.insecure_channel("localhost:50051")
        self.stub = api_pb2_grpc.ManagerStub(channel)

    def GetSuggestions(self, request, context):
        trials = []
        ret = self.stub.GetStudy(api_pb2.GetStudyRequest(
            study_id=request.study_id,
        ))

        algo_manager = AlgorithmManager(
            study_id=request.study_id,
            study_config=ret.study_config,
            X_train=[],
            y_train=[],
        )
        lowerbound = np.array(algo_manager.lower_bound)
        upperbound = np.array(algo_manager.upper_bound)

        cma = CMAES(
            dim=algo_manager.dim,
            upperbound=upperbound,
            lowerbound=lowerbound,
        )

        param_names = ['population', 'path_sigma', 'path_c', 'C', 'sigma', 'mean']
        param_info = {}
        for p in param_names:
            param_info[p] = dict(
                id="",
                value=""
            )

        ret = self.stub.GetSuggestionParameterList(api_pb2.GetSuggestionParameterListRequest(
            study_id=request.study_id,
        ))

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
        metrics = []

        path_sigma, path_c, C, sigma, mean = cma.init_params()
        for param in ret.suggestion_parameter_set:
            if param.param_name == "path_sigma":
                path_sigma = np.array(json.loads(param.suggestion_parameters[0].value))
                param_info["path_sigma"]["id"] = param.param_id

            elif param.param_name == "path_c":
                path_c = np.array(json.loads(param.suggestion_parameters[0].value))
                param_info["path_c"]["id"] = param.param_id

            elif param.param_name == "C":
                C = np.array(json.loads(param.suggestion_parameters[0].value))
                param_info["C"]["id"] = param.param_id

            elif param.param_name == "sigma":
                sigma = np.array(json.loads(param.suggestion_parameters[0].value))
                param_info["sigma"]["id"] = param.param_id

            elif param.param_name == "mean":
                mean = np.array(json.loads(param.suggestion_parameters[0].value))
                param_info["mean"]["id"] = param.param_id

            elif param.param_name == "population":
                param_info["population"]["id"] = param.param_id
                for p in param.suggestion_parameters:
                    value = json.loads(p.value)
                    if value["y"] == "":
                        ret = self.stub.GetTrial(api_pb2.GetTrialRequest(
                            trial_id=value["trial_id"],
                        ))

                        # the algorithm cannot continue without all trials in the population are evaluated
                        if ret.trial.objective_value == "":
                            context.set_code(grpc.StatusCode.UNKNOWN)
                            context.set_details("all trials in the population should be evaluated")
                            return api_pb2.GetSuggestionsReply(
                                trials=[],
                            )

                        # the algorithm is originally for minimization
                        if algo_manager.goal == api_pb2.MAXIMIZE:
                            y = -float(ret.trial.objective_value)
                        else:
                            y = float(ret.trial.objective_value)
                        metrics.append(dict(
                            x=np.array(json.loads(value["x"])),
                            y=y,
                            penalty=value["penalty"],
                        ))

        param_info["path_sigma"]["value"], param_info["path_c"]["value"], param_info["C"]["value"], \
        param_info["sigma"]["value"], param_info["mean"]["value"] = cma.report_metric(
            objective_dict=metrics,
            mean=mean,
            sigma=sigma,
            C=C,
            path_sigma=path_sigma,
            path_c=path_c,
        )

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
            mean=param_info["mean"]["value"],
            sigma=param_info["sigma"]["value"],
            C=param_info["C"]["value"],
        )

        suggestion_params = []
        for raw_suggestion in raw_suggestions:
            # parse the raw suggestions to desired format
            trial = algo_manager.parse_x_next(raw_suggestion["suggestion"])
            trial = algo_manager.convert_to_dict(trial)
            new_trial = api_pb2.Trial(
                study_id=request.study_id,
                parameter_set=[
                    api_pb2.Parameter(
                        name=x["name"],
                        value=str(x["value"]),
                        parameter_type=x["type"],
                    ) for x in trial
                ],
                status=api_pb2.PENDING,
                objective_value="",
            )
            ret = self.stub.CreateTrial(api_pb2.CreateTrialRequest(
                trial=new_trial
            ))
            new_trial.trial_id = ret.trial_id
            trials.append(new_trial)

            value = dict(
                trial_id=ret.trial_id,
                x=str(raw_suggestion["suggestion"].tolist()),
                y="",
                penalty=raw_suggestion["penalty"],
            )
            suggestion_params.append(api_pb2.SuggestionParameter(
                name="population",
                value=json.dumps(value)
            ))

        ret = self.stub.SetSuggestionParameters(api_pb2.SetSuggestionParametersRequest(
            study_id=request.study_id,
            param_id=param_info["population"]["id"],
            suggestion_algorithm=request.suggestion_algorithm,
            suggestion_parameters=suggestion_params,
        ))

        for param_name, info in param_info.items():
            if param_name != "population":
                ret = self.stub.SetSuggestionParameters(api_pb2.SetSuggestionParametersRequest(
                    study_id=request.study_id,
                    param_id=info["id"],
                    suggestion_algorithm=request.suggestion_algorithm,
                    suggestion_parameters=[api_pb2.SuggestionParameter(
                        name=param_name,
                        value=str(info["value"].tolist())
                    )]
                ))

        return api_pb2.GetSuggestionsReply(
            trials=trials,
        )
