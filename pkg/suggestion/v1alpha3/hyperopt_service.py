import logging

from pkg.apis.manager.v1alpha3.python import api_pb2
from pkg.apis.manager.v1alpha3.python import api_pb2_grpc
from pkg.suggestion.v1alpha3.internal.search_space import HyperParameter, HyperParameterSearchSpace
from pkg.suggestion.v1alpha3.internal.trial import Trial
from pkg.suggestion.v1alpha3.hyperopt.base_hyperopt_service import BaseHyperoptService

logger = logging.getLogger("HyperoptRandomService")


class HyperoptService(
        api_pb2_grpc.SuggestionServicer):
    def GetSuggestions(self, request, context):
        """
        Main function to provide suggestion.
        """
        base_serice = BaseHyperoptService(
            algorithm_name=request.experiment.spec.algorithm.algorithm_name)
        search_space = HyperParameterSearchSpace.convert(request.experiment)
        trials = Trial.convert(request.trials)
        new_trials = base_serice.getSuggestions(
            search_space, trials, request.request_number)
        return api_pb2.GetSuggestionsReply(
            trials=Trial.generate(new_trials)
        )
