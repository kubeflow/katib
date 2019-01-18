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
        print("BEFORE GET STUDY CONFIG")
        print(request)
        study_conf = self.getStudyConfig(request.study_id)
        print("Study Config")
        print(study_conf)
        print(type(study_conf))
        return api_pb2.GetSuggestionsReply(trials=trials)
    
    def getStudyConfig(self, studyID):
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            gsrep = client.GetStudy(api_pb2.GetStudyRequest(study_id=studyID, job_type="NAS"), 10)
            return gsrep.study_config
