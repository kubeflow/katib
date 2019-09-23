import logging

from pkg.apis.manager.v1alpha3.python import api_pb2
from pkg.apis.manager.v1alpha3.python import api_pb2_grpc
from pkg.suggestion.v1alpha3.internal.search_space import HyperParameter, HyperParameterSearchSpace
from pkg.suggestion.v1alpha3.internal.trial import Trial, Assignment
from pkg.suggestion.v1alpha3.skopt.base_skopt_service import BaseSkoptService
from pkg.suggestion.v1alpha3.base_health_service import HealthServicer


logger = logging.getLogger("SkoptService")


class SkoptService(
        api_pb2_grpc.SuggestionServicer, HealthServicer):
    def GetSuggestions(self, request, context):
        """
        Main function to provide suggestion.
        """
        name, config = OptimizerConfiguration.convertAlgorithmSpec(
            request.experiment.spec.algorithm)
        base_service = BaseSkoptService(
            algorithm_name=name,
            base_estimator=config.base_estimator,
            n_initial_points=config.n_initial_points,
            acq_func=config.acq_func,
            acq_optimizer=config.acq_optimizer,
            random_state=config.random_state)
        search_space = HyperParameterSearchSpace.convert(request.experiment)
        trials = Trial.convert(request.trials)
        new_trials = base_service.getSuggestions(
            search_space, trials, request.request_number)
        return api_pb2.GetSuggestionsReply(
            parameter_assignments=Assignment.generate(new_trials)
        )


class OptimizerConfiguration(object):
    def __init__(self, base_estimator="GP",
                 n_initial_points=10,
                 acq_func="gp_hedge",
                 acq_optimizer="auto",
                 random_state=None):
        self.base_estimator = base_estimator
        self.n_initial_points = n_initial_points
        self.acq_func = acq_func
        self.acq_optimizer = acq_optimizer
        self.random_state = random_state

    @staticmethod
    def convertAlgorithmSpec(algorithm_spec):
        optmizer = OptimizerConfiguration()
        for s in algorithm_spec.algorithm_setting:
            if s.name == "base_estimator":
                optmizer.base_estimator = s.value
            elif s.name == "n_initial_points":
                optmizer.n_initial_points = int(s.value)
            elif s.name == "acq_func":
                optmizer.acq_func = s.value
            elif s.name == "acq_optimizer":
                optmizer.acq_optimizer = s.value
            elif s.name == "random_state":
                optmizer.random_state = int(s.value)
        return algorithm_spec.algorithm_name, optmizer
