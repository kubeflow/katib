# pylint: disable=missing-docstring
import os

import yaml
import pytest
from box import Box
import numpy as np

import api_pb2
from katib_suggestion.parameters import ParameterConfig
from katib_suggestion.model.gp import GaussianProcessModel


TEST_DIR = os.path.dirname(os.path.realpath(__file__))


@pytest.fixture
def request_num():
    return 5

@pytest.fixture()
def lower_bounds():
    return [-5.0, -5, 2, 0, 0]


@pytest.fixture()
def upper_bounds():
    return [5.0, 5, 5, 1, 1]


@pytest.fixture
def names():
    return ["x", "y", "fake_discrete", "fake_categorical"]


@pytest.fixture
def name_ids():
    return {"x": 0, "y": 1, "fake_discrete": 2, "fake_categorical": 3}


@pytest.fixture
def dim():
    return 5


@pytest.fixture
def param_types():
    return [api_pb2.DOUBLE, api_pb2.INT, api_pb2.DISCRETE, api_pb2.CATEGORICAL]


@pytest.fixture
def discrete_info():
    return [{"name": "fake_discrete", "values": [2, 3, 5]}]


@pytest.fixture
def categorical_info():
    return [{"name": "fake_categorical", "values": ["true", "false"], "number": 2}]


@pytest.fixture
def parameter_config(name_ids, dim, lower_bounds,
                     upper_bounds, param_types, names,
                     discrete_info, categorical_info):
    parameter_config = ParameterConfig(name_ids=name_ids,
                                       dim=dim,
                                       lower_bounds=lower_bounds,
                                       upper_bounds=upper_bounds,
                                       parameter_types=param_types,
                                       names=names,
                                       discrete_info=discrete_info,
                                       categorical_info=categorical_info)
    return parameter_config


@pytest.fixture
def X_train():
    return np.array([[1.0, 1, 2, 1, 0], [1.0, 1, 3, 0, 1]])


@pytest.fixture
def y_train():
    return np.array([1.0, 1.0])


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


@pytest.fixture
def model(X_train, y_train):
    model = GaussianProcessModel()
    model.fit(X_train, y_train)
    return model


@pytest.fixture
def scaler(parameter_config):
    scaler = parameter_config.create_scaler()
    return scaler


@pytest.fixture
def rpc_request():
    return Box({"param_id": "test_param",
                "study_id": "test_study",
                "request_number": 2})
