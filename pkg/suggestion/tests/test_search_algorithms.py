import itertools

import numpy as np

from bayesianoptimization.src.algorithm_register import ALGORITHM_REGISTER
from bayesianoptimization.src.parsing_utils import parse_x_next_tuple


def test_bayesian_optimization(parameter_config, request_num, X_train, y_train):
    suggestion_config = {}
    boa = ALGORITHM_REGISTER["bayesian_optimization"](parameter_config, suggestion_config)
    response = boa.get_suggestion(X_train, y_train, request_num)
    assert len(response) == request_num


def test_grid_search(parameter_config):
    suggestion_config = {"DefaultGrid": 2}
    request_num = 24
    X_train, y_train = np.zeros(shape=(0, 2)), np.zeros(shape=(0,))
    boa = ALGORITHM_REGISTER["grid_search"](parameter_config, suggestion_config)
    response = boa.get_suggestion(X_train, y_train, request_num)
    correct_combinations = \
        itertools.product(("-5.0", "5.0"), ("-5", "5"), ("2", "3", "5"), ("true", "false"))
    correct_response = [parse_x_next_tuple(comb, parameter_config.parameter_types, parameter_config.names)
                        for comb in correct_combinations]
    assert response == correct_response


def test_random_search(parameter_config, request_num):
    suggestion_config = {}
    boa = ALGORITHM_REGISTER["random_search"](parameter_config, suggestion_config)
    response = boa.get_suggestion(request_num)
    assert len(response) == request_num
