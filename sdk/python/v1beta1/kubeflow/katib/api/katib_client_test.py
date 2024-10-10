import multiprocessing
from typing import List, Optional
from unittest.mock import MagicMock, Mock, patch

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

PVC_FAILED = "pvc creation failed"

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


def create_namespaced_persistent_volume_claim_response(*args, **kwargs):
    if kwargs.get("namespace") == PVC_FAILED:
        raise Exception("PVC creation failed")
    else:
        return {"metadata": {"name": "tune_test"}}


def list_namespaced_persistent_volume_claim_response(*args, **kwargs):
    if kwargs.get("namespace") == PVC_FAILED:
        mock_pvc = Mock()
        mock_pvc.metadata.name = "pvc_failed"
        mock_list = Mock()
        mock_list.items = [mock_pvc]
    else:
        mock_pvc = Mock()
        mock_pvc.metadata.name = "tune_test"
        mock_list = Mock()
        mock_list.items = [mock_pvc]
    return mock_list


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


# Mock classes for testing
class MockTransformerType:
    __name__ = "MockTransformerType"


class HuggingFaceModelParams:
    def __init__(
        self,
        model_uri=None,
        transformer_type=MockTransformerType,
        access_token=None,
        num_labels=None,
    ):
        self.model_uri = model_uri
        self.transformer_type = transformer_type
        self.access_token = access_token
        self.num_labels = num_labels


class HuggingFaceDatasetParams:
    def __init__(self, repo_id=None, access_token=None, split=None):
        self.repo_id = repo_id
        self.access_token = access_token
        self.split = split


class HuggingFaceTrainerParams:
    def __init__(self, training_parameters=None, lora_config=None):
        self.training_parameters = training_parameters
        self.lora_config = lora_config


class S3DatasetParams:
    def __init__(
        self,
        endpoint_url=None,
        bucket_name=None,
        file_key=None,
        region_name=None,
        access_key=None,
        secret_key=None,
    ):
        self.endpoint_url = endpoint_url
        self.bucket_name = bucket_name
        self.file_key = file_key
        self.region_name = region_name
        self.access_key = access_key
        self.secret_key = secret_key


class KubeflowOrgV1PyTorchJobSpec:
    def __init__(
        self,
        elastic_policy=None,
        nproc_per_node=None,
        pytorch_replica_specs={},
        run_policy=None,
    ):
        self.elastic_policy = elastic_policy
        self.nproc_per_node = nproc_per_node
        self.pytorch_replica_specs = pytorch_replica_specs
        self.run_policy = run_policy


class KubeflowOrgV1PyTorchJob:
    def __init__(
        self,
        api_version=None,
        kind=None,
        metadata=None,
        spec=KubeflowOrgV1PyTorchJobSpec,
        status=None,
    ):
        self.api_version = api_version
        self.kind = kind
        self.metadata = metadata
        self.spec = spec
        self.status = status


test_tune_data = [
    (
        "missing name",
        {
            "name": None,
            "objective": lambda x: x,
            "parameters": {"param": "value"},
        },
        ValueError,
    ),
    (
        "invalid hybrid parameters - objective and model_provider_parameters",
        {
            "name": "tune_test",
            "objective": lambda x: x,
            "model_provider_parameters": HuggingFaceModelParams(),
        },
        ValueError,
    ),
    (
        "missing parameters",
        {
            "name": "tune_test",
        },
        ValueError,
    ),
    (
        "missing parameters in custom objective tuning - lack parameters",
        {
            "name": "tune_test",
            "objective": lambda x: x,
        },
        ValueError,
    ),
    (
        "missing parameters in custom objective tuning - lack objective",
        {
            "name": "tune_test",
            "parameters": {"param": "value"},
        },
        ValueError,
    ),
    (
        "missing parameters in external model tuning - lack dataset_provider_parameters and trainer_parameters",
        {
            "name": "tune_test",
            "model_provider_parameters": HuggingFaceModelParams(),
        },
        ValueError,
    ),
    (
        "missing parameters in external model tuning - lack model_provider_parameters and trainer_parameters",
        {
            "name": "tune_test",
            "dataset_provider_parameters": HuggingFaceDatasetParams(),
        },
        ValueError,
    ),
    (
        "missing parameters in external model tuning - lack model_provider_parameters and dataset_provider_parameters",
        {
            "name": "tune_test",
            "trainer_parameters": HuggingFaceTrainerParams(),
        },
        ValueError,
    ),
    (
        "invalid env_per_trial",
        {
            "name": "tune_test",
            "objective": lambda x: x,
            "parameters": {"param": "value"},
            "env_per_trial": "invalid",
        },
        ValueError,
    ),
    (
        "invalid model_provider_parameters",
        {
            "name": "tune_test",
            "model_provider_parameters": "invalid",
            "dataset_provider_parameters": HuggingFaceDatasetParams(),
            "trainer_parameters": HuggingFaceTrainerParams(),
        },
        ValueError,
    ),
    (
        "invalid dataset_provider_parameters",
        {
            "name": "tune_test",
            "model_provider_parameters": HuggingFaceModelParams(),
            "dataset_provider_parameters": "invalid",
            "trainer_parameters": HuggingFaceTrainerParams(),
        },
        ValueError,
    ),
    (
        "invalid trainer_parameters",
        {
            "name": "tune_test",
            "model_provider_parameters": HuggingFaceModelParams(),
            "dataset_provider_parameters": HuggingFaceDatasetParams(),
            "trainer_parameters": "invalid",
        },
        ValueError,
    ),
    (
        "pvc creation failed",
        {
            "name": "tune_test",
            "namespace": PVC_FAILED,
            "model_provider_parameters": HuggingFaceModelParams(),
            "dataset_provider_parameters": HuggingFaceDatasetParams(),
            "trainer_parameters": HuggingFaceTrainerParams(),
        },
        RuntimeError,
    ),
    (
        "valid flow with custom objective tuning",
        {
            "name": "tune_test",
            "namespace": "tune",
            "objective": lambda x: x,
            "parameters": {"param": "value"},
        },
        TEST_RESULT_SUCCESS,
    ),
    (
        "valid flow with external model tuning",
        {
            "name": "tune_test",
            "namespace": "tune",
            "model_provider_parameters": HuggingFaceModelParams(),
            "dataset_provider_parameters": HuggingFaceDatasetParams(),
            "trainer_parameters": HuggingFaceTrainerParams(),
        },
        TEST_RESULT_SUCCESS,
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
    ), patch(
        "kubernetes.client.CoreV1Api",
        return_value=Mock(
            create_namespaced_persistent_volume_claim=Mock(
                side_effect=create_namespaced_persistent_volume_claim_response
            ),
            list_namespaced_persistent_volume_claim=Mock(
                side_effect=list_namespaced_persistent_volume_claim_response
            ),
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


@pytest.mark.parametrize("test_name,kwargs,expected_output", test_tune_data)
def test_tune(katib_client, test_name, kwargs, expected_output):
    """
    test tune function of katib client
    """
    print("\n\nExecuting test:", test_name)

    PYTORCHJOB_KIND = "PyTorchJob"
    JOB_PARAMETERS = {
        "PyTorchJob": {
            "model": "KubeflowOrgV1PyTorchJob",
            "plural": "pytorchjobs",
            "container": "pytorch",
            "base_image": "docker.io/pytorch/pytorch:2.1.2-cuda11.8-cudnn8-runtime",
        }
    }

    with patch.dict(
        "sys.modules",
        {
            "kubeflow.storage_initializer": Mock(),
            "kubeflow.storage_initializer.hugging_face": Mock(),
            "kubeflow.storage_initializer.s3": Mock(),
            "kubeflow.storage_initializer.constants": Mock(),
            "kubeflow.training": MagicMock(),
            "kubeflow.training.models": Mock(),
            "kubeflow.training.utils": Mock(),
            "kubeflow.training.constants": Mock(),
            "kubeflow.training.constants.constants": Mock(),
        },
    ), patch(
        "kubeflow.storage_initializer.hugging_face.HuggingFaceModelParams",
        HuggingFaceModelParams,
    ), patch(
        "kubeflow.storage_initializer.hugging_face.HuggingFaceDatasetParams",
        HuggingFaceDatasetParams,
    ), patch(
        "kubeflow.storage_initializer.hugging_face.HuggingFaceTrainerParams",
        HuggingFaceTrainerParams,
    ), patch(
        "kubeflow.storage_initializer.s3.S3DatasetParams", S3DatasetParams
    ), patch(
        "kubeflow.training.models.KubeflowOrgV1PyTorchJob", KubeflowOrgV1PyTorchJob
    ), patch(
        "kubeflow.training.constants.constants.JOB_PARAMETERS", JOB_PARAMETERS
    ), patch(
        "kubeflow.training.constants.constants.PYTORCHJOB_KIND", PYTORCHJOB_KIND
    ), patch(
        "kubeflow.katib.utils.utils.get_trial_substitutions_from_trainer",
        return_value={"param": "value"},
    ), patch.object(
        katib_client, "create_experiment", return_value=Mock()
    ) as mock_create_experiment:
        try:
            katib_client.tune(**kwargs)
            mock_create_experiment.assert_called_once()
            assert expected_output == TEST_RESULT_SUCCESS
        except Exception as e:
            assert type(e) is expected_output
        print("test execution complete")
