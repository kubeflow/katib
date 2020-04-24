import logging
import grpc

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
        name, config = OptimizerConfiguration.convert_algorithm_spec(
            request.experiment.spec.algorithm)

        if self.is_first_run:
            search_space = HyperParameterSearchSpace.convert(request.experiment)
            self.base_service = BaseHyperoptService(
                algorithm_name=name,
                algorithm_conf=config,
                search_space=search_space)
            self.is_first_run = False

        trials = Trial.convert(request.trials)
        new_assignments = self.base_service.getSuggestions(trials, request.request_number)
        return api_pb2.GetSuggestionsReply(
            parameter_assignments=Assignment.generate(new_assignments)
        )

    def ValidateAlgorithmSettings(self, request, context):
        is_valid, message = OptimizerConfiguration.validate_algorithm_spec(
            request.experiment.spec.algorithm)
        if not is_valid:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(message)
            logger.error(message)
        return api_pb2.ValidateAlgorithmSettingsReply()


class OptimizerConfiguration:
    __schema_dict = {
        'tpe': {
            'gamma': (lambda x: float(x), lambda x: 1 > float(x) > 0),
            'prior_weight': (lambda x: float(x), lambda x: float(x) > 0),
            'n_EI_candidates': (lambda x: int(x), lambda x: int(x) > 0),
            "random_state": (lambda x: int(x), lambda x: int(x) >= 0),
        },
        "random": {
            "random_state": (lambda x: int(x), lambda x: int(x) >= 0)
        }
    }

    @classmethod
    def convert_algorithm_spec(cls, algorithm_spec):
        ret = {}
        setting_schema = cls.__schema_dict[algorithm_spec.algorithm_name]
        for s in algorithm_spec.algorithm_setting:
            if s.name in setting_schema:
                ret[s.name] = setting_schema[s.name][0](s.value)

        return algorithm_spec.algorithm_name, ret

    @classmethod
    def validate_algorithm_spec(cls, algorithm_spec):
        algo_name = algorithm_spec.algorithm_name
        if algo_name not in cls.__schema_dict:
            return False, "unknown algorithm name %s" % algo_name

        setting_schema = cls.__schema_dict[algo_name]
        for s in algorithm_spec.algorithm_setting:
            if s.name not in setting_schema:
                return False, "unknown setting %s for algorithm %s" % (s.name, algo_name)
            try:
                if not setting_schema[s.name][1](s.value):
                    return False, "invalid value %s for setting %s" % (s.value, s.name)
            except Exception as e:
                return False, "invalid value %s for setting %s" % (s.value, s.name)

        return True, ""
