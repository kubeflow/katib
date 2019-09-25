import logging
import grpc

from pkg.apis.manager.v1alpha3.python import api_pb2
from pkg.apis.manager.v1alpha3.python import api_pb2_grpc
from pkg.apis.manager.health.python import health_pb2

from pkg.suggestion.v1alpha3.internal.search_space import HyperParameter, HyperParameterSearchSpace, DOUBLE
from pkg.suggestion.v1alpha3.internal.trial import Trial, Assignment
from pkg.suggestion.v1alpha3.chocolate.base_chocolate_service import BaseChocolateService
from pkg.suggestion.v1alpha3.base_health_service import HealthServicer

logger = logging.getLogger("ChocolateService")


class ChocolateService(
        api_pb2_grpc.SuggestionServicer, HealthServicer):
    def ValidateAlgorithmSettings(self, request, context):
        algorithm_name = request.experiment.spec.algorithm.algorithm_name
        if algorithm_name == "grid":
            search_space = HyperParameterSearchSpace.convert(
                request.experiment)
            for param in search_space.params:
                if param.type == DOUBLE:
                    if param.step == "" or param.step == None:
                        return self._set_validate_context_error(
                            context, "param {} step is nil".format(param.name))

    def GetSuggestions(self, request, context):
        """
        Main function to provide suggestion.
        """
        base_serice = BaseChocolateService(
            algorithm_name=request.experiment.spec.algorithm.algorithm_name)
        search_space = HyperParameterSearchSpace.convert(request.experiment)
        trials = Trial.convert(request.trials)
        new_assignments = base_serice.getSuggestions(
            search_space, trials, request.request_number)
        return api_pb2.GetSuggestionsReply(
            parameter_assignments=Assignment.generate(new_assignments)
        )

    def _set_validate_context_error(self, context, error_message):
        context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
        context.set_details(error_message)
        logger.info(error_message)
        return api_pb2.ValidateAlgorithmSettingsReply()
