import logging
import grpc
from pkg.api.v1alpha2.python import api_pb2
from pkg.api.v1alpha2.python import api_pb2_grpc
from . import parsing_util, base_hyperopt_service
from .internal.search_space import HyperParameter, HyperParameterSearchSpace
from .internal.trial import Trial

logger = logging.getLogger("HyperoptRandomService")

class HyperoptRandomService(
        api_pb2_grpc.SuggestionServicer,
        base_hyperopt_service.BaseHyperoptService):
    def __init__(self):
        super(
            HyperoptRandomService, self).__init__(algorithm_name="random")

    def GetSuggestions(self, request, context):
        """
        Main function to provide suggestion.
        """
        search_space = HyperParameterSearchSpace.convert(request.experiment)
        trials = Trial.convert(request.trials)
        new_trials = super().getSuggestions(search_space, trials, request.request_number)
        return api_pb2.GetSuggestionsReply(
            trials=Trial.generate(new_trials)
        )
