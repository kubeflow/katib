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
        study_conf = self.getStudyConfig(request.study_id)
        print("Search Space")
        print(study_conf)

        suggestion_parameters = self.getSuggestionParameters(request.param_id)
        print("Suggestion Parameters")
        print(suggestion_parameters)

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
