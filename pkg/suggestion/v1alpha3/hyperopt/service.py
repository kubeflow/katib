import logging
import grpc

from pkg.apis.manager.v1alpha3.python import api_pb2
from pkg.apis.manager.v1alpha3.python import api_pb2_grpc

from pkg.suggestion.v1alpha3.internal.search_space import HyperParameterSearchSpace
from pkg.suggestion.v1alpha3.internal.trial import Trial, Assignment
from pkg.suggestion.v1alpha3.hyperopt.base_service import BaseHyperoptService
from pkg.suggestion.v1alpha3.internal.base_health_service import HealthServicer

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
    __conversion_dict = {
        'tpe': {
            'gamma': lambda x: float(x),
            'prior_weight': lambda x: float(x),
            'n_EI_candidates': lambda x: int(x),
            "random_state": lambda x: int(x),
        },
        "random": {
            "random_state": lambda x: int(x),
        }
    }

    @classmethod
    def convert_algorithm_spec(cls, algorithm_spec):
        ret = {}
        setting_schema = cls.__conversion_dict[algorithm_spec.algorithm_name]
        for s in algorithm_spec.algorithm_setting:
            if s.name in setting_schema:
                ret[s.name] = setting_schema[s.name](s.value)

        return algorithm_spec.algorithm_name, ret

    @classmethod
    def validate_algorithm_spec(cls, algorithm_spec):
        algo_name = algorithm_spec.algorithm_name
        if algo_name == 'tpe':
            return cls._validate_tpe_setting(algorithm_spec.algorithm_setting)
        elif algo_name == 'random':
            return cls._validate_random_setting(algorithm_spec.algorithm_setting)
        else:
            return False, "unknown algorithm name {}".format(algo_name)

    @classmethod
    def _validate_tpe_setting(cls, algorithm_settings):
        for s in algorithm_settings:
            try:
                if s.name == 'gamma':
                    if not 1 > float(s.value) > 0:
                        return False, "gamma should be in the range of (0, 1)"
                elif s.name == 'prior_weight':
                    if not float(s.value) > 0:
                        return False, "prior_weight should be great than zero"
                elif s.name == 'n_EI_candidates':
                    if not int(s.value) > 0:
                        return False, "n_EI_candidates should be great than zero"
                elif s.name == 'random_state':
                    if not int(s.value) >= 0:
                        return False, "random_state should be great or equal than zero"
                else:
                    return False, "unknown setting {} for algorithm tpe".format(s.name)
            except Exception as e:
                return False, "failed to validate {name}({value}): {exception}".format(
                    name=s.name, value=s.value, exception=e)

        return True, ""

    @classmethod
    def _validate_random_setting(cls, algorithm_settings):
        for s in algorithm_settings:
            try:
                if s.name == 'random_state':
                    if not (int(s.value) >= 0):
                        return False, "random_state should be great or equal than zero"
                else:
                    return False, "unknown setting {} for algorithm random".format(s.name)
            except Exception as e:
                return False, "failed to validate {name}({value}): {exception}".format(
                    name=s.name, value=s.value, exception=e)

        return True, ""
