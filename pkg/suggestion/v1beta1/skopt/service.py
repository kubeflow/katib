import logging

from pkg.apis.manager.v1beta1.python import api_pb2
from pkg.apis.manager.v1beta1.python import api_pb2_grpc
from pkg.suggestion.v1beta1.internal.search_space import HyperParameterSearchSpace
from pkg.suggestion.v1beta1.internal.trial import Trial, Assignment
from pkg.suggestion.v1beta1.skopt.base_service import BaseSkoptService
from pkg.suggestion.v1beta1.internal.base_health_service import HealthServicer


logger = logging.getLogger(__name__)


class SkoptService(api_pb2_grpc.SuggestionServicer, HealthServicer):

    def __init__(self):
        super(SkoptService, self).__init__()
        self.base_service = None
        self.is_first_run = True

    def GetSuggestions(self, request, context):
        """
        Main function to provide suggestion.
        """
        algorithm_name, config = OptimizerConfiguration.convertAlgorithmSpec(
            request.experiment.spec.algorithm)
        if algorithm_name != "bayesianoptimization":
            raise Exception("Failed to create the algorithm: {}".format(algorithm_name))

        if self.is_first_run:
            search_space = HyperParameterSearchSpace.convert(request.experiment)
            self.base_service = BaseSkoptService(
                base_estimator=config.base_estimator,
                n_initial_points=config.n_initial_points,
                acq_func=config.acq_func,
                acq_optimizer=config.acq_optimizer,
                random_state=config.random_state,
                search_space=search_space)
            self.is_first_run = False

        trials = Trial.convert(request.trials)
        new_trials = self.base_service.getSuggestions(trials, request.request_number)
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
        optimizer = OptimizerConfiguration()
        for s in algorithm_spec.algorithm_setting:
            if s.name == "base_estimator":
                optimizer.base_estimator = s.value
            elif s.name == "n_initial_points":
                optimizer.n_initial_points = int(s.value)
            elif s.name == "acq_func":
                optimizer.acq_func = s.value
            elif s.name == "acq_optimizer":
                optimizer.acq_optimizer = s.value
            elif s.name == "random_state":
                optimizer.random_state = int(s.value)
        return algorithm_spec.algorithm_name, optimizer
