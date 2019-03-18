import numpy as np

import api_pb2
from bayesianoptimization.src import parsing_utils


def test_parsing_study_config(study_config, dim, names, name_ids,
                              param_types, lower_bounds,
                              upper_bounds, discrete_info,
                              categorical_info):
    parsed_config = parsing_utils.parse_parameter_configs(study_config.parameter_configs.configs)
    assert parsed_config.name_ids == name_ids
    assert parsed_config.parameter_types == param_types
    assert parsed_config.names == names
    assert parsed_config.dim == dim
    assert np.allclose(parsed_config.lower_bounds.ravel(), lower_bounds)
    assert np.allclose(parsed_config.upper_bounds.ravel(), upper_bounds)
    assert parsed_config.discrete_info == discrete_info
    assert parsed_config.categorical_info == categorical_info


def test_parsing_x_next(param_types, names,
                        discrete_info, categorical_info):
    x_next = [1.0, 1, 5, "true"]
    parsed_x_next = parsing_utils.parse_x_next_vector(x_next,
                                                      param_types,
                                                      names,
                                                      discrete_info,
                                                      categorical_info)
    assert parsed_x_next == \
           [{"name": "x", "value": 1.0, "type": api_pb2.DOUBLE},
            {"name": "y", "value": 1, "type": api_pb2.INT},
            {"name": "fake_discrete", "value": 5, "type": api_pb2.DISCRETE},
            {"name": "fake_categorical", "value": "true",
             "type": api_pb2.CATEGORICAL}]


def test_parsing_past_observations(observations, dim,
                                   name_ids, param_types,
                                   categorical_info, X_train):
    X_train = parsing_utils.parse_previous_observations(
        observations.parameters,
        dim,
        name_ids,
        param_types,
        categorical_info)
    assert np.allclose(X_train, X_train)


def test_parsing_past_metrics(study_config, observations, y_train):
    y_train = parsing_utils.parse_metric(observations.metrics, study_config.optimization_type)
    assert np.allclose(y_train, y_train)
