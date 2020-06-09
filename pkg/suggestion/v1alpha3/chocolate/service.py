import logging
import grpc

from pkg.apis.manager.v1alpha3.python import api_pb2
from pkg.apis.manager.v1alpha3.python import api_pb2_grpc

from pkg.suggestion.v1alpha3.internal.constant import DOUBLE
from pkg.suggestion.v1alpha3.internal.search_space import HyperParameterSearchSpace
from pkg.suggestion.v1alpha3.internal.trial import Trial, Assignment
from pkg.suggestion.v1alpha3.chocolate.base_service import BaseChocolateService
from pkg.suggestion.v1alpha3.internal.base_health_service import HealthServicer

logger = logging.getLogger(__name__)


class ChocolateService(api_pb2_grpc.SuggestionServicer, HealthServicer):
    def __init__(self):
        super(ChocolateService, self).__init__()
        self.base_service = None
        self.is_first_run = True

    def ValidateAlgorithmSettings(self, request, context):
        algorithm_name = request.experiment.spec.algorithm.algorithm_name
        if algorithm_name == "grid":
            search_space = HyperParameterSearchSpace.convert(
                request.experiment)
            for param in search_space.params:
                if param.type == DOUBLE:
                    if param.step == "" or param.step is None:
                        return self._set_validate_context_error(
                            context, "param {} step is nil".format(param.name))
        return api_pb2.ValidateAlgorithmSettingsReply()

    def GetSuggestions(self, request, context):
        """
        Main function to provide suggestion.
        """

        if self.is_first_run:
            search_space = HyperParameterSearchSpace.convert(
                request.experiment)
            self.base_service = BaseChocolateService(
                algorithm_name=request.experiment.spec.algorithm.algorithm_name,
                search_space=search_space)
            self.is_first_run = False

        trials = Trial.convert(request.trials)
        try:
            new_assignments = self.base_service.getSuggestions(
                trials, request.request_number)
        except StopIteration as e:
            return self._set_get_suggestions_context_error(
                context, grpc.StatusCode.NOT_FOUND, str(e))
        return api_pb2.GetSuggestionsReply(
            parameter_assignments=Assignment.generate(new_assignments)
        )

    def _set_validate_context_error(self, context, error_message):
        context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
        context.set_details(error_message)
        logger.info(error_message)
        return api_pb2.ValidateAlgorithmSettingsReply()

    def _set_get_suggestions_context_error(self, context, code, details):
        context.set_code(code)
        context.set_details(details)
        logger.info("GetSuggestions raised an error: (%s, %s) %s",
                    code.value[0], code.value[1], details)
        return api_pb2.GetSuggestionsReply()
