from unittest.mock import patch

import pytest
from kubeflow.katib import report_metrics
from kubeflow.katib.constants import constants

TEST_RESULT_SUCCESS = "success"
ENV_VARIABLE_EMPTY = True
ENV_VARIABLE_NOT_EMPTY = False


def report_observation_log_response(*args, **kwargs):
    if kwargs.get("timeout") == 0:
        raise TimeoutError


test_report_metrics_data = [
    (
        "valid metrics with float type",
        {"metrics": {"result": 0.99}, "timeout": constants.DEFAULT_TIMEOUT},
        TEST_RESULT_SUCCESS,
        ENV_VARIABLE_NOT_EMPTY,
    ),
    (
        "valid metrics with string type",
        {"metrics": {"result": "0.99"}, "timeout": constants.DEFAULT_TIMEOUT},
        TEST_RESULT_SUCCESS,
        ENV_VARIABLE_NOT_EMPTY,
    ),
    (
        "valid metrics with int type",
        {"metrics": {"result": 1}, "timeout": constants.DEFAULT_TIMEOUT},
        TEST_RESULT_SUCCESS,
        ENV_VARIABLE_NOT_EMPTY,
    ),
    (
        "ReportObservationLog timeout error",
        {"metrics": {"result": 0.99}, "timeout": 0},
        RuntimeError,
        ENV_VARIABLE_NOT_EMPTY,
    ),
    (
        "invalid metrics with type string",
        {"metrics": {"result": "abc"}, "timeout": constants.DEFAULT_TIMEOUT},
        ValueError,
        ENV_VARIABLE_NOT_EMPTY,
    ),
    (
        "Trial name is not passed to env variables",
        {"metrics": {"result": 0.99}, "timeout": constants.DEFAULT_TIMEOUT},
        ValueError,
        ENV_VARIABLE_EMPTY,
    ),
]


@pytest.fixture
def mock_getenv(request):
    with patch("os.getenv") as mock:
        if request.param is ENV_VARIABLE_EMPTY:
            mock.side_effect = ValueError
        else:
            mock.return_value = "example"
        yield mock


@pytest.fixture
def mock_get_current_k8s_namespace():
    with patch("kubeflow.katib.utils.utils.get_current_k8s_namespace") as mock:
        mock.return_value = "test"
        yield mock


@pytest.fixture
def mock_report_observation_log():
    with patch("kubeflow.katib.katib_api_pb2_grpc.DBManagerStub") as mock:
        mock_instance = mock.return_value
        mock_instance.ReportObservationLog.side_effect = report_observation_log_response
        yield mock_instance


@pytest.mark.parametrize(
    "test_name,kwargs,expected_output,mock_getenv",
    test_report_metrics_data,
    indirect=["mock_getenv"],
)
def test_report_metrics(
    test_name,
    kwargs,
    expected_output,
    mock_getenv,
    mock_get_current_k8s_namespace,
    mock_report_observation_log,
):
    """
    test report_metrics function
    """
    print("\n\nExecuting test:", test_name)
    try:
        report_metrics(**kwargs)
        assert expected_output == TEST_RESULT_SUCCESS
    except Exception as e:
        assert type(e) is expected_output
    print("test execution complete")
