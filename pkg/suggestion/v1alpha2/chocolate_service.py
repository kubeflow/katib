import logging
import grpc
from pkg.api.v1alpha2.python import api_pb2
from pkg.api.v1alpha2.python import api_pb2_grpc
from . import parsing_util, base_chocolate_service
from .internal.search_space import HyperParameter, HyperParameterSearchSpace
from .internal.trial import Trial

logger = logging.getLogger("ChocolateService")


class ChocolateService(
        api_pb2_grpc.SuggestionServicer):
    def GetSuggestions(self, request, context):
        """
        Main function to provide suggestion.
        """
        name = request.experiment.spec.algorithm.algorithm_name
        base_service = base_chocolate_service.BaseChocolateService(
            algorithm_name=name)
        search_space = HyperParameterSearchSpace.convert(request.experiment)
        trials = Trial.convert(request.trials)
        new_trials = base_service.getSuggestions(
            search_space, trials, request.request_number)
        return api_pb2.GetSuggestionsReply(
            trials=Trial.generate(new_trials)
        )
