import numpy as np

from ..bayesianoptimization.src.bayesian_optimization_algorithm import BOAlgorithm


def test_boa(parameter_config, request_num, correct_X_train, correct_y_train):
    suggestion_config = {}
    boa = BOAlgorithm(parameter_config, suggestion_config, correct_X_train, correct_y_train)
    response = boa.get_suggestion(request_num)
    assert len(response) == request_num
