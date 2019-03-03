import pytest

from ..bayesianoptimization.src.algorithm_register import ALGORITHM_REGISTER


@pytest.mark.parametrize("algorithm", ["random_search", "bayesian_optimization"])
def test_boa(algorithm, parameter_config, request_num, correct_X_train, correct_y_train):
    suggestion_config = {}
    boa = ALGORITHM_REGISTER[algorithm](parameter_config, suggestion_config, correct_X_train, correct_y_train)
    response = boa.get_suggestion(request_num)
    assert len(response) == request_num
