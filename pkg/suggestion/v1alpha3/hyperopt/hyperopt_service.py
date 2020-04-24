import logging

from pkg.apis.manager.v1alpha3.python import api_pb2
from pkg.apis.manager.v1alpha3.python import api_pb2_grpc

from pkg.suggestion.v1alpha3.internal.search_space import HyperParameterSearchSpace
from pkg.suggestion.v1alpha3.internal.trial import Trial, Assignment
from pkg.suggestion.v1alpha3.hyperopt.base_hyperopt_service import BaseHyperoptService
from pkg.suggestion.v1alpha3.base_health_service import HealthServicer

logger = logging.getLogger(__name__)


class HyperoptService(api_pb2_grpc.SuggestionServicer, HealthServicer):

    def __init__(self):
        super(HyperoptService, self).__init__()
        self.base_service = None
        self.is_first_run = True

    def GetSuggestions(self, request, context):
        """
        Main function to provide suggestion.
        """
        name, config = OptimizerConfiguration.convertAlgorithmSpec(
            request.experiment.spec.algorithm)

        if self.is_first_run:
            search_space = HyperParameterSearchSpace.convert(request.experiment)
            self.base_service = BaseHyperoptService(
                algorithm_name=name,
                random_state=config.random_state,
                search_space=search_space)
            self.is_first_run = False

        trials = Trial.convert(request.trials)
        new_assignments = self.base_service.getSuggestions(trials, request.request_number)
        return api_pb2.GetSuggestionsReply(
            parameter_assignments=Assignment.generate(new_assignments)
        )

    def ValidateAlgorithmSettings(self, request, context):
        return api_pb2.ValidateAlgorithmSettingsReply()


class OptimizerConfiguration(object):
    def __init__(self, random_state=None):
        self.random_state = random_state

    @staticmethod
    def convertAlgorithmSpec(algorithm_spec):
        optimizer = OptimizerConfiguration()
        for s in algorithm_spec.algorithm_setting:
            if s.name == "random_state":
                optimizer.random_state = int(s.value)
        return algorithm_spec.algorithm_name, optimizer
