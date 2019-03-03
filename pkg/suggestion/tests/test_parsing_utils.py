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


def test_parsing_study_config(study_config, correct_dim, correct_names, correct_name_ids,
                              correct_param_types, correct_lower_bounds,
                              correct_upper_bounds, correct_discrete_info,
                              correct_categorical_info):
    parsed_config = parsing_utils.parse_parameter_configs(study_config.parameter_configs.configs)
    assert parsed_config.name_ids == correct_name_ids
    assert parsed_config.parameter_types == correct_param_types
    assert parsed_config.names == correct_names
    assert parsed_config.dim == correct_dim
    assert parsed_config.lower_bounds == correct_lower_bounds
    assert parsed_config.upper_bounds == correct_upper_bounds
    assert parsed_config.discrete_info == correct_discrete_info
    assert parsed_config.categorical_info == correct_categorical_info


def test_parsing_x_next(correct_param_types, correct_names,
                        correct_discrete_info, correct_categorical_info):
    x_next = [1.0, 1, 5, "true"]
    parsed_x_next = parsing_utils.parse_x_next(x_next,
                                               correct_param_types,
                                               correct_names,
                                               correct_discrete_info,
                                               correct_categorical_info)
    assert parsed_x_next == \
           [{"name": "x", "value": 1.0, "type": api_pb2.DOUBLE},
            {"name": "y", "value": 1, "type": api_pb2.INT},
            {"name": "fake_discrete", "value": 5, "type": api_pb2.DISCRETE},
            {"name": "fake_categorical", "value": "true",
             "type": api_pb2.CATEGORICAL}]


def test_parsing_past_observations(observations, correct_dim,
                                   correct_name_ids, correct_param_types,
                                   correct_categorical_info, correct_X_train):
    X_train = parsing_utils.parse_previous_observations(
        observations.parameters,
        correct_dim,
        correct_name_ids,
        correct_param_types,
        correct_categorical_info)
    assert np.allclose(X_train, correct_X_train)


def test_parsing_past_metrics(study_config, observations, correct_y_train):
    y_train = parsing_utils.parse_metric(observations.metrics, study_config.optimization_type)
    assert np.allclose(y_train, correct_y_train)
