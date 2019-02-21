import yaml
import pytest
import numpy as np
from box import Box

from pkg.api.python import api_pb2
from ..bayesianoptimization.src.algorithm_manager import AlgorithmManager



@pytest.fixture
def study_config():
    with open("tests/study_config.yaml", "r", encoding="utf-8") as f:
        contents = yaml.safe_load(f)
    return Box(contents)


@pytest.fixture
def observations():
    with open("tests/parameter_values.yaml", "r", encoding="utf-8") as f:
        contents = yaml.safe_load(f)
    return Box(contents)


def test_algorithm_manager(study_config, observations):
    study_id = "test_id"
    x_next = [1.0, 1.0]
    manager = AlgorithmManager(study_id, study_config,
                               observations.parameters, observations.metrics)
    assert manager.study_id == study_id
    assert manager.study_config == study_config
    assert manager.goal == api_pb2.MAXIMIZE
    assert all(t == api_pb2.DOUBLE for t in manager.types)
    assert manager.names == ["x", "y"]
    assert manager.dim == 2
    assert manager.lower_bound == [-5.0, -5.0]
    assert manager.upper_bound == [5.0, 5.0]
    assert np.allclose(manager.X_train, np.array([[1.0, 1.0], [1.0, 1.0]]))
    assert np.allclose(manager.y_train, np.array([1.0, 1.0]))
    parsed_x_next = manager.parse_x_next(x_next)
    x_next_dict = manager.convert_to_dict(parsed_x_next)
    assert x_next_dict == \
           [{"name": "x", "value": 1.0, "type": api_pb2.DOUBLE},
            {"name": "y", "value": 1.0, "type": api_pb2.DOUBLE}]
