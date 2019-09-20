import logging

from pkg.apis.manager.v1alpha3.python import api_pb2
from pkg.apis.manager.v1alpha3.python import api_pb2_grpc
from pkg.apis.manager.health.python import health_pb2

from pkg.suggestion.v1alpha3.internal.search_space import HyperParameter, HyperParameterSearchSpace
from pkg.suggestion.v1alpha3.internal.trial import Trial, Assignment
from pkg.suggestion.v1alpha3.hyperopt.base_hyperopt_service import BaseHyperoptService
from pkg.suggestion.v1alpha3.base_health_service import HealthServicer

logger = logging.getLogger("HyperoptRandomService")


class HyperoptService(
        api_pb2_grpc.SuggestionServicer, HealthServicer):
    def GetSuggestions(self, request, context):
        """
        Main function to provide suggestion.
        """
        base_serice = BaseHyperoptService(
            algorithm_name=request.experiment.spec.algorithm.algorithm_name)
        search_space = HyperParameterSearchSpace.convert(request.experiment)
        trials = Trial.convert(request.trials)
        new_assignments = base_serice.getSuggestions(
            search_space, trials, request.request_number)
        return api_pb2.GetSuggestionsReply(
            parameter_assignments=Assignment.generate(new_assignments)
        )
