import random
import string

import grpc
import numpy as np

from pkg.api.python import api_pb2
from pkg.api.python import api_pb2_grpc
import logging
from logging import getLogger, StreamHandler, INFO, DEBUG


class NasrlService(api_pb2_grpc.SuggestionServicer):
    def __init__(self, logger=None):
        self.manager_addr = "vizier-core"
        self.manager_port = 6789
        self.current_trial_id = ""

    def GetSuggestions(self, request, context):
        trials = []
        print("INSIDE GET SUGGESTION WITH request %s" %request)

        study_conf = self.getStudyConfig(request.study_id)
        suggestion_parameters = self.getSuggestionParameters(request.param_id)
        
        if self.current_trial_id != "":
            print("Current trial id is: {}".format(self.current_trial_id))
            print("Getting Trial")
            channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
            with api_pb2.beta_create_Manager_stub(channel) as client:
                tr_rep = client.GetTrial(api_pb2.GetTrialRequest(trial_id=self.current_trial_id), 10)
                print("Reponse is: ")
                print(tr_rep)
        trials.append(api_pb2.Trial(
                study_id=request.study_id,
                parameter_set=[
                    api_pb2.Parameter(
                        name="test_name",
                        value="test_value",
                        parameter_type= api_pb2.CATEGORICAL,
                    ),
                    api_pb2.Parameter(
                        name="katib",
                        value="test",
                        parameter_type= api_pb2.CATEGORICAL,
                    )
                ], 
            )
        )
        print("TRIALS CREATED")
        
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            print("INSIDE CLIENT")
            for i, t in enumerate(trials):
                ctrep = client.CreateTrial(api_pb2.CreateTrialRequest(trial=t), 10)
                trials[i].trial_id = ctrep.trial_id
                self.current_trial_id = ctrep.trial_id
            print("TRIALS INSERTED")
            print(ctrep.trial_id)

        print("METRICS LOG START")

        self.getEvalHistory(request.study_id)

        print("METRICS LOG END")

        print("END OF SUGGESTION")

        return api_pb2.GetSuggestionsReply(trials=trials)
    
    def getStudyConfig(self, studyID):
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            gsrep = client.GetStudy(api_pb2.GetStudyRequest(study_id=studyID), 10)
            return gsrep.study_config

    def getSuggestionParameters(self, paramID):
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            gsprep = client.GetSuggestionParameters(api_pb2.GetSuggestionParametersRequest(param_id=paramID), 10)
            return gsprep.suggestion_parameters


    def getEvalHistory(self, studyID):
        worker_hist = []
   
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            gwfrep = client.GetWorkerFullInfo(api_pb2.GetWorkerFullInfoRequest(study_id=studyID, only_latest_log=True), 10)
            worker_hist = gwfrep.worker_full_infos

        for w in worker_hist:
            if w.Worker.status == api_pb2.COMPLETED:
                for ml in w.metrics_logs:
                    print(ml)
                    print
