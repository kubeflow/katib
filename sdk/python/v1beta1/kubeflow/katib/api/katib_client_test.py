import multiprocessing
import pytest
from unittest.mock import patch, Mock
from kubeflow.katib import models
from kubeflow.katib.api.katib_client import KatibClient

@pytest.fixture
def katib_client():
    return KatibClient()

@pytest.fixture
def valid_experiment():
    return models.V1beta1Experiment(
        metadata=models.V1ObjectMeta(name="test-experiment"),
        spec=models.V1beta1ExperimentSpec()
    )

class TestCreateExperiment:
    @patch('kubeflow.katib.api.katib_client.client.CustomObjectsApi.create_namespaced_custom_object')
    def test_create_experiment_success(self, mock_create_namespaced_custom_object, katib_client, valid_experiment):
        mock_create_namespaced_custom_object.return_value = {"metadata": {"name": "test-experiment"}}
        experiment = katib_client.create_experiment(valid_experiment)
        assert experiment == "Experiment default/test-experiment has been created"

    @patch('kubeflow.katib.api.katib_client.client.CustomObjectsApi.create_namespaced_custom_object')
    def test_create_experiment_timeout_error(self, mock_create_namespaced_custom_object, katib_client, valid_experiment):
        mock_create_namespaced_custom_object.side_effect = multiprocessing.TimeoutError()
        with pytest.raises(TimeoutError):
            katib_client.create_experiment(valid_experiment, namespace="timeout")

    @patch('kubeflow.katib.api.katib_client.client.CustomObjectsApi.create_namespaced_custom_object')
    def test_create_experiment_runtime_error(self, mock_create_namespaced_custom_object, katib_client, valid_experiment):
        mock_create_namespaced_custom_object.side_effect = RuntimeError()
        with pytest.raises(RuntimeError):
            katib_client.create_experiment(valid_experiment, namespace="runtime")

class TestGetExperiment:
    @patch('kubeflow.katib.api.katib_client.client.CustomObjectsApi.get_namespaced_custom_object')
    def test_get_experiment_success(self, mock_get_namespaced_custom_object, katib_client, valid_experiment):
        mock_response = {"metadata": {"name": "test-experiment"}, "spec": {}, "status": {}}
        mock_get_namespaced_custom_object.return_value = mock_response
        experiment = katib_client.get_experiment("test-experiment")
        assert experiment == valid_experiment

    @patch('kubeflow.katib.api.katib_client.client.CustomObjectsApi.get_namespaced_custom_object')
    def test_get_experiment_timeout_error(self, mock_get_namespaced_custom_object, katib_client):
        mock_get_namespaced_custom_object.side_effect = multiprocessing.TimeoutError()
        with pytest.raises(TimeoutError):
            katib_client.get_experiment("test-experiment", namespace="timeout")

    @patch('kubeflow.katib.api.katib_client.client.CustomObjectsApi.get_namespaced_custom_object')
    def test_get_experiment_runtime_error(self, mock_get_namespaced_custom_object, katib_client):
        mock_get_namespaced_custom_object.side_effect = RuntimeError()
        with pytest.raises(RuntimeError):
            katib_client.get_experiment("test-experiment", namespace="runtime")
