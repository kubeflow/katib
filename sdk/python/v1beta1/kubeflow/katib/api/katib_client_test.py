import multiprocessing
from typing import List, Optional
from unittest.mock import Mock, patch

import kubeflow.katib as katib
import kubeflow.katib.katib_api_pb2 as katib_api_pb2
import pytest
import transformers
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
from kubeflow.katib.types import types
from kubeflow.storage_initializer.hugging_face import (
    HuggingFaceDatasetParams,
    HuggingFaceModelParams,
    HuggingFaceTrainerParams,
)
from kubeflow.training.models import KubeflowOrgV1PyTorchJob
from kubernetes.client import V1Job, V1ObjectMeta

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


test_tune_data = [
    (
        "missing name",
        {
            "name": None,
            "objective": lambda x: print(f"a={x}"),
            "parameters": {"a": katib.search.int(min=10, max=100)},
        },
        ValueError,
    ),
    (
        "invalid name format",
        {
            "name": "Llama3.1-fine-tune",
        },
        ValueError,
    ),
    (
        "invalid hybrid parameters - objective and model_provider_parameters",
        {
            "name": "tune_test",
            "objective": lambda x: print(f"a={x}"),
            "model_provider_parameters": HuggingFaceModelParams(
                model_uri="hf://google-bert/bert-base-cased",
                transformer_type=transformers.AutoModelForSequenceClassification,
                num_labels=5,
            ),
        },
        ValueError,
    ),
    (
        "missing parameters - no custom objective or external model tuning",
        {
            "name": "tune_test",
        },
        ValueError,
    ),
    (
        "missing parameters in custom objective tuning - lack parameters",
        {
            "name": "tune_test",
            "objective": lambda x: print(f"a={x}"),
        },
        ValueError,
    ),
    (
        "missing parameters in custom objective tuning - lack objective",
        {
            "name": "tune_test",
            "parameters": {"a": katib.search.int(min=10, max=100)},
        },
        ValueError,
    ),
    (
        "missing parameters in external model tuning - lack dataset_provider_parameters "
        "and trainer_parameters",
        {
            "name": "tune_test",
            "model_provider_parameters": HuggingFaceModelParams(
                model_uri="hf://google-bert/bert-base-cased",
                transformer_type=transformers.AutoModelForSequenceClassification,
                num_labels=5,
            ),
        },
        ValueError,
    ),
    (
        "missing parameters in external model tuning - lack model_provider_parameters "
        "and trainer_parameters",
        {
            "name": "tune_test",
            "dataset_provider_parameters": HuggingFaceDatasetParams(
                repo_id="yelp_review_full",
                split="train[:3000]",
            ),
        },
        ValueError,
    ),
    (
        "missing parameters in external model tuning - lack model_provider_parameters "
        "and dataset_provider_parameters",
        {
            "name": "tune_test",
            "trainer_parameters": HuggingFaceTrainerParams(
                training_parameters=transformers.TrainingArguments(
                    output_dir="test_tune_api",
                    learning_rate=katib.search.double(min=1e-05, max=5e-05),
                ),
            ),
        },
        ValueError,
    ),
    (
        "invalid env_per_trial",
        {
            "name": "tune_test",
            "objective": lambda x: print(f"a={x}"),
            "parameters": {"a": katib.search.int(min=10, max=100)},
            "env_per_trial": "invalid",
        },
        ValueError,
    ),
    (
        "invalid model_provider_parameters",
        {
            "name": "tune_test",
            "model_provider_parameters": "invalid",
            "dataset_provider_parameters": HuggingFaceDatasetParams(
                repo_id="yelp_review_full",
                split="train[:3000]",
            ),
            "trainer_parameters": HuggingFaceTrainerParams(
                training_parameters=transformers.TrainingArguments(
                    output_dir="test_tune_api",
                    learning_rate=katib.search.double(min=1e-05, max=5e-05),
                ),
            ),
        },
        ValueError,
    ),
    (
        "invalid dataset_provider_parameters",
        {
            "name": "tune_test",
            "model_provider_parameters": HuggingFaceModelParams(
                model_uri="hf://google-bert/bert-base-cased",
                transformer_type=transformers.AutoModelForSequenceClassification,
                num_labels=5,
            ),
            "dataset_provider_parameters": "invalid",
            "trainer_parameters": HuggingFaceTrainerParams(
                training_parameters=transformers.TrainingArguments(
                    output_dir="test_tune_api",
                    learning_rate=katib.search.double(min=1e-05, max=5e-05),
                ),
            ),
        },
        ValueError,
    ),
    (
        "invalid trainer_parameters",
        {
            "name": "tune_test",
            "model_provider_parameters": HuggingFaceModelParams(
                model_uri="hf://google-bert/bert-base-cased",
                transformer_type=transformers.AutoModelForSequenceClassification,
                num_labels=5,
            ),
            "dataset_provider_parameters": HuggingFaceDatasetParams(
                repo_id="yelp_review_full",
                split="train[:3000]",
            ),
            "trainer_parameters": "invalid",
        },
        ValueError,
    ),
    (
        "pvc creation failed",
        {
            "name": "tune_test",
            "namespace": PVC_FAILED,
            "model_provider_parameters": HuggingFaceModelParams(
                model_uri="hf://google-bert/bert-base-cased",
                transformer_type=transformers.AutoModelForSequenceClassification,
                num_labels=5,
            ),
            "dataset_provider_parameters": HuggingFaceDatasetParams(
                repo_id="yelp_review_full",
                split="train[:3000]",
            ),
            "trainer_parameters": HuggingFaceTrainerParams(
                training_parameters=transformers.TrainingArguments(
                    output_dir="test_tune_api",
                    learning_rate=katib.search.double(min=1e-05, max=5e-05),
                ),
            ),
            "resources_per_trial": types.TrainerResources(
                num_workers=2,
                num_procs_per_worker=2,
                resources_per_worker={"gpu": "2"},
            ),
        },
        RuntimeError,
    ),
    (
        "valid flow with custom objective function and Job as Trial",
        {
            "name": "tune_test",
            "objective": lambda x: print(f"a={x}"),
            "parameters": {"a": katib.search.int(min=10, max=100)},
            "objective_metric_name": "a",
            "resources_per_trial": {"gpu": "2"},
        },
        TEST_RESULT_SUCCESS,
    ),
    (
        "valid flow with custom objective function and PyTorchJob as Trial",
        {
            "name": "tune_test",
            "objective": lambda x: print(f"a={x}"),
            "parameters": {"a": katib.search.int(min=10, max=100)},
            "objective_metric_name": "a",
            "resources_per_trial": types.TrainerResources(
                num_workers=2,
                num_procs_per_worker=2,
                resources_per_worker={"gpu": "2"},
            ),
        },
        TEST_RESULT_SUCCESS,
    ),
    (
        "valid flow with external model tuning",
        {
            "name": "tune_test",
            "model_provider_parameters": HuggingFaceModelParams(
                model_uri="hf://google-bert/bert-base-cased",
                transformer_type=transformers.AutoModelForSequenceClassification,
                num_labels=5,
            ),
            "dataset_provider_parameters": HuggingFaceDatasetParams(
                repo_id="yelp_review_full",
                split="train[:3000]",
            ),
            "trainer_parameters": HuggingFaceTrainerParams(
                training_parameters=transformers.TrainingArguments(
                    output_dir="test_tune_api",
                    learning_rate=katib.search.double(min=1e-05, max=5e-05),
                ),
            ),
            "resources_per_trial": types.TrainerResources(
                num_workers=2,
                num_procs_per_worker=2,
                resources_per_worker={"gpu": "2"},
            ),
            "objective_metric_name": "train_loss",
            "objective_type": "minimize",
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

    with patch.object(
        katib_client, "create_experiment", return_value=Mock()
    ) as mock_create_experiment:
        try:
            katib_client.tune(**kwargs)
            mock_create_experiment.assert_called_once()

            if expected_output == TEST_RESULT_SUCCESS:
                assert expected_output == TEST_RESULT_SUCCESS
                call_args = mock_create_experiment.call_args
                experiment = call_args[0][0]

                if (
                    test_name
                    == "valid flow with custom objective function and Job as Trial"
                ):
                    # Verify input_params
                    args_content = "".join(
                        experiment.spec.trial_template.trial_spec.spec.template.spec.containers[
                            0
                        ].args
                    )
                    assert "'a': '${trialParameters.a}'" in args_content
                    # Verify trial_params
                    assert experiment.spec.trial_template.trial_parameters == [
                        V1beta1TrialParameterSpec(name="a", reference="a"),
                    ]
                    # Verify experiment_params
                    assert experiment.spec.parameters == [
                        V1beta1ParameterSpec(
                            name="a",
                            parameter_type="int",
                            feasible_space=V1beta1FeasibleSpace(min="10", max="100"),
                        ),
                    ]
                    # Verify objective_spec
                    assert experiment.spec.objective == V1beta1ObjectiveSpec(
                        type="maximize",
                        objective_metric_name="a",
                        additional_metric_names=[],
                    )
                    # Verity Trial spec
                    assert isinstance(experiment.spec.trial_template.trial_spec, V1Job)

                elif (
                    test_name
                    == "valid flow with custom objective function and PyTorchJob as Trial"
                ):
                    # Verity Trial spec
                    assert isinstance(
                        experiment.spec.trial_template.trial_spec,
                        KubeflowOrgV1PyTorchJob,
                    )

                elif test_name == "valid flow with external model tuning":
                    # Verify input_params
                    args_content = "".join(
                        experiment.spec.trial_template.trial_spec.spec.pytorch_replica_specs[
                            "Master"
                        ]
                        .template.spec.containers[0]
                        .args
                    )
                    assert (
                        '"learning_rate": "${trialParameters.learning_rate}"'
                        in args_content
                    )
                    # Verify trial_params
                    assert experiment.spec.trial_template.trial_parameters == [
                        V1beta1TrialParameterSpec(
                            name="learning_rate", reference="learning_rate"
                        ),
                    ]
                    # Verify experiment_params
                    assert experiment.spec.parameters == [
                        V1beta1ParameterSpec(
                            name="learning_rate",
                            parameter_type="double",
                            feasible_space=V1beta1FeasibleSpace(
                                min="1e-05", max="5e-05"
                            ),
                        ),
                    ]
                    # Verify objective_spec
                    assert experiment.spec.objective == V1beta1ObjectiveSpec(
                        type="minimize",
                        objective_metric_name="train_loss",
                        additional_metric_names=[],
                    )

        except Exception as e:
            assert type(e) is expected_output
        print("test execution complete")
