import multiprocessing
from typing import List, Optional
from unittest.mock import Mock, patch

import kubeflow.katib.katib_api_pb2 as katib_api_pb2
import pytest
from kubeflow.katib import (
    KatibClient,
    V1beta1AlgorithmSpec,
    V1beta1Experiment,
    V1beta1ExperimentSpec,
    V1beta1FeasibleSpace,
    V1beta1ObjectiveSpec,
    V1beta1ParameterSpec,
    V1beta1TrialParameterSpec,
    V1beta1TrialTemplate,
)
from kubeflow.katib.constants import constants
from kubernetes.client import V1ObjectMeta

TEST_RESULT_SUCCESS = "success"


class ConflictException(Exception):
    def __init__(self):
        self.status = 409


def create_namespaced_custom_object_response(*args, **kwargs):
    if args[2] == "timeout":
        raise multiprocessing.TimeoutError()
    elif args[2] == "conflict":
        raise ConflictException()
    elif args[2] == "runtime":
        raise Exception()
    elif args[2] in ("test", "test-name"):
        return {"metadata": {"name": "experiment-mnist-ci-test"}}
    elif args[2] == "test-generate-name":
        return {"metadata": {"name": "12345-experiment-mnist-ci-test"}}


def get_observation_log_response(*args, **kwargs):
    if kwargs.get("timeout") == 0:
        raise TimeoutError
    elif args[0].trial_name == "invalid":
        raise RuntimeError
    else:
        return katib_api_pb2.GetObservationLogReply(
            observation_log=katib_api_pb2.ObservationLog(
                metric_logs=[
                    katib_api_pb2.MetricLog(
                        time_stamp="2024-07-29T15:09:08Z",
                        metric=katib_api_pb2.Metric(name="result", value="0.99"),
                    )
                ]
            )
        )


def generate_trial_template() -> V1beta1TrialTemplate:
    trial_spec = {
        "apiVersion": "batch/v1",
        "kind": "Job",
        "spec": {
            "template": {
                "metadata": {"annotations": {"sidecar.istio.io/inject": "false"}},
                "spec": {
                    "containers": [
                        {
                            "name": "training-container",
                            "image": "docker.io/kubeflowkatib/pytorch-mnist-cpu:v0.14.0",
                            "command": [
                                "python3",
                                "/opt/pytorch-mnist/mnist.py",
                                "--epochs=1",
                                "--batch-size=64",
                                "--lr=${trialParameters.learningRate}",
                                "--momentum=${trialParameters.momentum}",
                            ],
                        }
                    ],
                    "restartPolicy": "Never",
                },
            }
        },
    }

    return V1beta1TrialTemplate(
        primary_container_name="training-container",
        trial_parameters=[
            V1beta1TrialParameterSpec(
                name="learningRate",
                description="Learning rate for the training model",
                reference="lr",
            ),
            V1beta1TrialParameterSpec(
                name="momentum",
                description="Momentum for the training model",
                reference="momentum",
            ),
        ],
        trial_spec=trial_spec,
    )


def generate_experiment(
    metadata: V1ObjectMeta,
    algorithm_spec: V1beta1AlgorithmSpec,
    objective_spec: V1beta1ObjectiveSpec,
    parameters: List[V1beta1ParameterSpec],
    trial_template: V1beta1TrialTemplate,
) -> V1beta1Experiment:
    return V1beta1Experiment(
        api_version=constants.API_VERSION,
        kind=constants.EXPERIMENT_KIND,
        metadata=metadata,
        spec=V1beta1ExperimentSpec(
            max_trial_count=3,
            parallel_trial_count=2,
            max_failed_trial_count=1,
            algorithm=algorithm_spec,
            objective=objective_spec,
            parameters=parameters,
            trial_template=trial_template,
        ),
    )


def create_experiment(
    name: Optional[str] = None, generate_name: Optional[str] = None
) -> V1beta1Experiment:
    experiment_namespace = "test"

    if name is not None:
        metadata = V1ObjectMeta(name=name, namespace=experiment_namespace)
    elif generate_name is not None:
        metadata = V1ObjectMeta(
            generate_name=generate_name, namespace=experiment_namespace
        )
    else:
        metadata = V1ObjectMeta(namespace=experiment_namespace)

    algorithm_spec = V1beta1AlgorithmSpec(algorithm_name="random")

    objective_spec = V1beta1ObjectiveSpec(
        type="minimize",
        goal=0.001,
        objective_metric_name="loss",
    )

    parameters = [
        V1beta1ParameterSpec(
            name="lr",
            parameter_type="double",
            feasible_space=V1beta1FeasibleSpace(min="0.01", max="0.06"),
        ),
        V1beta1ParameterSpec(
            name="momentum",
            parameter_type="double",
            feasible_space=V1beta1FeasibleSpace(min="0.5", max="0.9"),
        ),
    ]

    trial_template = generate_trial_template()

    experiment = generate_experiment(
        metadata, algorithm_spec, objective_spec, parameters, trial_template
    )
    return experiment


test_create_experiment_data = [
    (
        "experiment name and generate_name missing",
        {"experiment": create_experiment()},
        ValueError,
    ),
    (
        "create_namespaced_custom_object timeout error",
        {
            "experiment": create_experiment(name="experiment-mnist-ci-test"),
            "namespace": "timeout",
        },
        TimeoutError,
    ),
    (
        "create_namespaced_custom_object conflict error",
        {
            "experiment": create_experiment(name="experiment-mnist-ci-test"),
            "namespace": "conflict",
        },
        Exception,
    ),
    (
        "create_namespaced_custom_object runtime error",
        {
            "experiment": create_experiment(name="experiment-mnist-ci-test"),
            "namespace": "runtime",
        },
        RuntimeError,
    ),
    (
        "valid flow with experiment type V1beta1Experiment and name",
        {
            "experiment": create_experiment(name="experiment-mnist-ci-test"),
            "namespace": "test-name",
        },
        TEST_RESULT_SUCCESS,
    ),
    (
        "valid flow with experiment type V1beta1Experiment and generate_name",
        {
            "experiment": create_experiment(generate_name="experiment-mnist-ci-test"),
            "namespace": "test-generate-name",
        },
        TEST_RESULT_SUCCESS,
    ),
    (
        "valid flow with experiment JSON and name",
        {
            "experiment": {
                "metadata": {
                    "name": "experiment-mnist-ci-test",
                }
            },
            "namespace": "test-name",
        },
        TEST_RESULT_SUCCESS,
    ),
    (
        "valid flow with experiment JSON and generate_name",
        {
            "experiment": {
                "metadata": {
                    "generate_name": "experiment-mnist-ci-test",
                }
            },
            "namespace": "test-generate-name",
        },
        TEST_RESULT_SUCCESS,
    ),
]


test_get_trial_metrics_data = [
    (
        "valid trial name",
        {"name": "example", "namespace": "valid", "timeout": constants.DEFAULT_TIMEOUT},
        [
            katib_api_pb2.MetricLog(
                time_stamp="2024-07-29T15:09:08Z",
                metric=katib_api_pb2.Metric(name="result", value="0.99"),
            )
        ],
    ),
    (
        "invalid trial name",
        {
            "name": "invalid",
            "namespace": "invalid",
            "timeout": constants.DEFAULT_TIMEOUT,
        },
        RuntimeError,
    ),
    (
        "GetObservationLog timeout error",
        {"name": "example", "namespace": "valid", "timeout": 0},
        RuntimeError,
    ),
]


@pytest.fixture
def katib_client():
    with patch(
        "kubernetes.client.CustomObjectsApi",
        return_value=Mock(
            create_namespaced_custom_object=Mock(
                side_effect=create_namespaced_custom_object_response
            )
        ),
    ), patch("kubernetes.config.load_kube_config", return_value=Mock()), patch(
        "kubeflow.katib.katib_api_pb2_grpc.DBManagerStub",
        return_value=Mock(
            GetObservationLog=Mock(side_effect=get_observation_log_response)
        ),
    ):
        client = KatibClient()
        yield client


@pytest.mark.parametrize(
    "test_name,kwargs,expected_output", test_create_experiment_data
)
def test_create_experiment(katib_client, test_name, kwargs, expected_output):
    """
    test create_experiment function of katib client
    """
    print("\n\nExecuting test:", test_name)
    try:
        katib_client.create_experiment(**kwargs)
        assert expected_output == TEST_RESULT_SUCCESS
    except Exception as e:
        assert type(e) is expected_output
    print("test execution complete")


@pytest.mark.parametrize(
    "test_name,kwargs,expected_output", test_get_trial_metrics_data
)
def test_get_trial_metrics(katib_client, test_name, kwargs, expected_output):
    """
    test get_trial_metrics function of katib client
    """
    print("\n\nExecuting test:", test_name)
    try:
        metrics = katib_client.get_trial_metrics(**kwargs)
        for i in range(len(metrics)):
            assert metrics[i] == expected_output[i]
    except Exception as e:
        assert type(e) is expected_output
    print("test execution complete")
