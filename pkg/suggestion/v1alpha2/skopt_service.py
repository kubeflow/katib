import logging
import grpc
from pkg.api.v1alpha2.python import api_pb2
from pkg.api.v1alpha2.python import api_pb2_grpc
from . import parsing_util, base_skopt_service
from .internal.search_space import HyperParameter, HyperParameterSearchSpace
from .internal.trial import Trial

logger = logging.getLogger("HyperoptRandomService")


class SkoptService(
        api_pb2_grpc.SuggestionServicer):
    def GetSuggestions(self, request, context):
        """
        Main function to provide suggestion.
        """
        base_serice = base_skopt_service.BaseSkoptService(
            algorithm_name=request.experiment.spec.algorithm.algorithm_name)
        search_space = HyperParameterSearchSpace.convert(request.experiment)
        trials = Trial.convert(request.trials)
        new_trials = base_serice.getSuggestions(search_space, trials, request.request_number)
        return api_pb2.GetSuggestionsReply(
            trials=Trial.generate(new_trials)
        )
