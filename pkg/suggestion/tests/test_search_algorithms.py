import itertools

import numpy as np

from ..bayesianoptimization.src.algorithm_register import ALGORITHM_REGISTER


def test_bayesian_optimization(parameter_config, request_num, X_train, y_train):
    suggestion_config = {}
    boa = ALGORITHM_REGISTER["bayesian_optimization"](parameter_config, suggestion_config)
    response = boa.get_suggestion(X_train, y_train, request_num)
    assert len(response) == request_num


def test_grid_search(parameter_config):
    suggestion_config = {"grid_size": 2}
    request_num = 24
    X_train, y_train = np.zeros(shape=(0, 2)), np.zeros(shape=(0,))
    boa = ALGORITHM_REGISTER["grid_search"](parameter_config, suggestion_config)
    response = boa.get_suggestion(X_train, y_train, request_num)
    correct_response = itertools.product((-5.0, 5.0), (-5, 5), (2, 3, 5), ("true", "false"))
    assert response == list(correct_response)


def test_random_search(parameter_config, request_num):
    suggestion_config = {}
    boa = ALGORITHM_REGISTER["random_search"](parameter_config, suggestion_config)
    response = boa.get_suggestion(request_num)
    assert len(response) == request_num
