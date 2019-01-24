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

    def GetSuggestions(self, request, context):
        trials = []
        print("INSIDE GET SUGGESTION WITH request %s" %request)

        study_conf = self.getStudyConfig(request.study_id)
        suggestion_parameters = self.getSuggestionParameters(request.param_id)
        
        trials.append(api_pb2.Trial(
                study_id=request.study_id,
                parameter_set=[
                    api_pb2.Parameter(
                        name="Test_name",
                        value="Test_value",
                        parameter_type= api_pb2.CATEGORICAL,
                    )
                ]
            )
        )
        print("TRIALS CREATED")
        
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            print("INSIDE CLIENT")
            for i, t in enumerate(trials):
                ctrep = client.CreateTrial(api_pb2.CreateTrialRequest(trial=t), 10)
                trials[i].trial_id = ctrep.trial_id
            print("TRIALS INSERTED")
            print(ctrep.trial_id)
        
        print("END OF SUGGESTION")

        return api_pb2.GetSuggestionsReply(trials=trials)
    
    def getStudyConfig(self, studyID):
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            gsrep = client.GetStudy(api_pb2.GetStudyRequest(study_id=studyID, job_type="NAS"), 10)
            return gsrep.study_config

    def getSuggestionParameters(self, paramID):
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            gsprep = client.GetSuggestionParameters(api_pb2.GetSuggestionParametersRequest(param_id=paramID), 10)
            return gsprep.suggestion_parameters
