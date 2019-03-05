import pytest

from ..bayesianoptimization.src.algorithm_register import ALGORITHM_REGISTER


@pytest.mark.parametrize("algorithm", ["random_search", "bayesian_optimization"])
def test_boa(algorithm, parameter_config, request_num, X_train, y_train):
    suggestion_config = {}
    boa = ALGORITHM_REGISTER[algorithm](parameter_config, suggestion_config, X_train, y_train)
    response = boa.get_suggestion(request_num)
    assert len(response) == request_num
