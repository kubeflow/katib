# pylint: disable=missing-docstring
import pytest
import numpy as np

from pkg.api.python import api_pb2
from ..bayesianoptimization.src.parameters import ParameterConfig



@pytest.fixture
def request_num():
    return 5

@pytest.fixture()
def correct_lower_bounds():
    return [-5.0, -5, 2, 0, 0]


@pytest.fixture()
def correct_upper_bounds():
    return [5.0, 5, 5, 1, 1]


@pytest.fixture
def correct_names():
    return ["x", "y", "fake_discrete", "fake_categorical"]


@pytest.fixture
def correct_name_ids():
    return {"x": 0, "y": 1, "fake_discrete": 2, "fake_categorical": 3}


@pytest.fixture
def correct_dim():
    return 5


@pytest.fixture
def correct_param_types():
    return [api_pb2.DOUBLE, api_pb2.INT, api_pb2.DISCRETE, api_pb2.CATEGORICAL]


@pytest.fixture
def correct_discrete_info():
    return [{"name": "fake_discrete", "values": [2, 3, 5]}]


@pytest.fixture
def correct_categorical_info():
    return [{"name": "fake_categorical", "values": ["true", "false"], "number": 2}]


@pytest.fixture
def parameter_config(correct_name_ids, correct_dim, correct_lower_bounds,
                     correct_upper_bounds, correct_param_types, correct_names,
                     correct_discrete_info, correct_categorical_info):
    parameter_config = ParameterConfig(name_ids=correct_name_ids,
                                       dim=correct_dim,
                                       lower_bounds=correct_lower_bounds,
                                       upper_bounds=correct_upper_bounds,
                                       parameter_types=correct_param_types,
                                       names=correct_names,
                                       discrete_info=correct_discrete_info,
                                       categorical_info=correct_categorical_info)
    return parameter_config


@pytest.fixture
def correct_X_train():
    return np.array([[1.0, 1, 2, 1, 0], [1.0, 1, 3, 0, 1]])


@pytest.fixture
def correct_y_train():
    return np.array([1.0, 1.0])


def booth_function(X):
    f = (X[:, 0] + 2 * X[:, 1] - 7) ** 2 + (2 * X[:, 0] + X[:, 1] - 5) ** 2
    return f


@pytest.fixture
def lower_bounds():
    return [-5, -5]


@pytest.fixture
def upper_bounds():
    return [5, 5]


@pytest.fixture
def dim():
    return 2


@pytest.fixture
def X_train(lower_bounds, upper_bounds):
    x_range = np.arange(lower_bounds[0], upper_bounds[0] + 1)
    y_range = np.arange(lower_bounds[1], upper_bounds[1] + 1)
    X_train = np.array([(x, y) for x in x_range for y in y_range])
    return X_train


@pytest.fixture
def X_test():
    x_range, y_range = np.arange(-1.5, 2.5), np.arange(-1.5, 2.5)
    X_test = np.array([(x, y) for x in x_range for y in y_range])
    return X_test


@pytest.fixture
def y_train(X_train):
    y_train = -booth_function(X_train)
    return y_train
