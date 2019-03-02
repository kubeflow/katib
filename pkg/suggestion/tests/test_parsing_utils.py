import os

import yaml
import pytest
import numpy as np
from box import Box

from pkg.api.python import api_pb2
from ..bayesianoptimization.src import parsing_utils


TEST_DIR = os.path.dirname(os.path.realpath(__file__))


@pytest.fixture
def study_config():
    with open(os.path.join(TEST_DIR, "study_config.yaml"), "r") as f:
        contents = yaml.safe_load(f)
    return Box(contents)


@pytest.fixture
def observations():
    with open(os.path.join(TEST_DIR, "parameter_values.yaml"), "r") as f:
        contents = yaml.safe_load(f)
    return Box(contents)


def test_parsing_utils(study_config, observations):
    name_ids, dim, lower_bounds, upper_bounds, parameter_types, names, discrete_info, categorical_info = \
        parsing_utils.parse_parameter_configs(
            study_config.parameter_configs.configs)
    assert parameter_types == [api_pb2.DOUBLE, api_pb2.INT, api_pb2.DISCRETE, api_pb2.CATEGORICAL]
    assert names == ["x", "y", "fake_discrete", "fake_categorical"]
    assert dim == 5
    assert lower_bounds == [-5.0, -5, 2, 0, 0]
    assert upper_bounds == [5.0, 5, 5, 1, 1]
    assert discrete_info == [{"name": "fake_discrete", "values": [2, 3, 5]}]
    assert categorical_info == \
           [{"name": "fake_categorical", "values": ["true", "false"], "number": 2}]
    x_next = [1.0, 1, 5, "true"]
    parsed_x_next = parsing_utils.parse_x_next(x_next, parameter_types, names, discrete_info, categorical_info)
    assert parsed_x_next == \
           [{"name": "x", "value": 1.0, "type": api_pb2.DOUBLE},
            {"name": "y", "value": 1, "type": api_pb2.INT},
            {"name": "fake_discrete", "value": 5, "type": api_pb2.DISCRETE},
            {"name": "fake_categorical", "value": "true",
             "type": api_pb2.CATEGORICAL}]
    X_train = parsing_utils.parse_previous_observations(observations.parameters, dim, name_ids, parameter_types, categorical_info)
    y_train = parsing_utils.parse_metric(observations.metrics, study_config.optimization_type)
    assert np.allclose(X_train, np.array([[1.0, 1, 2, 1, 0], [1.0, 1, 3, 0, 1]]))
    assert np.allclose(y_train, np.array([1.0, 1.0]))
