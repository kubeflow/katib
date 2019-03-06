import pytest

from ..bayesianoptimization.src.algorithm_register import ALGORITHM_REGISTER


@pytest.mark.parametrize("algorithm", ["bayesian_optimization"])
def test_use_previous(algorithm, parameter_config, request_num, X_train, y_train):
    suggestion_config = {}
    boa = ALGORITHM_REGISTER[algorithm](parameter_config, suggestion_config)
    response = boa.get_suggestion(X_train, y_train, request_num)
    assert len(response) == request_num


@pytest.mark.parametrize("algorithm", ["random_search"])
def test_ignore_previous(algorithm, parameter_config, request_num):
    suggestion_config = {}
    boa = ALGORITHM_REGISTER[algorithm](parameter_config, suggestion_config)
    response = boa.get_suggestion(request_num)
    assert len(response) == request_num
