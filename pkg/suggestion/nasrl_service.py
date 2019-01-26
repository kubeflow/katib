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

        return api_pb2.GetSuggestionsReply(trials=trials)
