from concurrent import futures
import random
import string
import time
import grpc
import numpy as np
import sys

from api import api_pb2
from api import api_pb2_grpc
from suggestion.BO.bayesian_optimization_algorithm import BOAlgorithm
from suggestion.algorithm_manager import AlgorithmManager


class BayesianService(api_pb2_grpc.SuggestionServicer):
    def __init__(self):
        # {
        #     study_id:[
        #         {
        #             trial_id:
        #             metric:
        #             parameters: []
        #         }
        #     ]
        # }
        self.trial_hist = {}
        # {
        #     study_id:{
        #         N:
        #     }
        # }
        self.service_params = {}

    def GenerateTrials(self, request, context):
        if request.study_id not in self.trial_hist.keys():
            self.trial_hist[request.study_id] = []
        X_train = []
        y_train = []

        for x in request.completed_trials:
            for trial in self.trial_hist[x.study_id]:
                if trial["trial_id"] == x.trial_id:
                    trial["metric"] = x.objective_value

        for x in self.trial_hist[request.study_id]:
            if x["metric"] is not None:
                X_train.append(x["parameters"])
                y_train.append(x["metric"])

        algo_manager = AlgorithmManager(
            study_id=request.study_id,
            study_config=request.configs,
            X_train=X_train,
            y_train=y_train,
        )

        lowerbound = np.array(algo_manager.lower_bound)
        upperbound = np.array(algo_manager.upper_bound)
        # print("lowerbound", lowerbound)
        # print("upperbound", upperbound)
        alg = BOAlgorithm(
            dim=algo_manager.dim,
            N=self.service_params[request.study_id]["N"],
            lowerbound=lowerbound,
            upperbound=upperbound,
            X_train=algo_manager.X_train,
            y_train=algo_manager.y_train,
        )
        x_next = alg.get_suggestion().squeeze()

        # todo: maybe there is a better way to generate a trial_id
        trial_id = ''.join(random.sample(string.ascii_letters + string.digits, 12))
        self.trial_hist[request.study_id].append(dict({
            "trial_id": trial_id,
            "parameters": x_next,
            "metric": None,
        }))
        # print(x_next)

        x_next = algo_manager.parse_x_next(x_next)
        x_next = algo_manager.convert_to_dict(x_next)
        trial = api_pb2.Trial(
            trial_id=trial_id,
            study_id=request.study_id,
            parameter_set=[
                api_pb2.Parameter(
                    name=x["name"],
                    value=str(x["value"]),
                    parameter_type=x["type"],
                ) for x in x_next
            ],
            status=api_pb2.PENDING,
            eval_logs=[],
        )
        # print(self.trial_hist)

        return api_pb2.GenerateTrialsReply(
            trials=[trial],
            completed=False,
        )

    def SetSuggestionParameters(self, request, context):
        if request.study_id not in self.service_params.keys():
            self.service_params[request.study_id] = {}
        for param in request.suggestion_parameters:
            if param.name != "N":
                print("unknown parameter name")
                sys.exit(1)
            self.service_params[request.study_id][param.name] = int(param.value)
        return api_pb2.SetSuggestionParametersReply()

    def StopSuggestion(self, request, context):
        if request.study_id in self.service_params.keys():
            del self.service_params[request.study_id]
            del self.trial_hist[request.study_id]
        return api_pb2.StopStudyReply()
