import unittest
from unittest import TestCase
from unittest.mock import Mock
from unittest.mock import patch

from kubeflow.katib import KatibClient
from kubeflow.katib import models
import kubeflow.katib as katib
from kubeflow.storage_initializer.hugging_face import HuggingFaceDatasetParams
from kubeflow.storage_initializer.hugging_face import HuggingFaceModelParams
from kubeflow.storage_initializer.hugging_face import HuggingFaceTrainerParams
from kubeflow.training import models as training_models
from kubernetes import client
from kubernetes.client.exceptions import ApiException
from peft import LoraConfig
import transformers


class TestTuneAPI(TestCase):
    # Create an instance of the KatibClient
    def setUp(self):
        self.katib_client = KatibClient(namespace="default")

    # Test input
    # Test for missing required parameters
    def test_tune_missing_name(self):
        with self.assertRaises(ValueError) as context:
            self.katib_client.tune(
                name=None,
                objective=lambda x: x,
                parameters={
                    "a": katib.search.int(min=10, max=100),
                    "b": katib.search.double(min=0.1, max=0.2),
                },
            )

        self.assertIn("Please specify name for the Experiment.", str(context.exception))

    # Test for invalid hyperparameter optimization configuration
    # Case 1: Set two options: 1) external models and datasets; 2) custom objective at the same time
    def test_tune_invalid_with_model_provider_and_objective(self):
        with self.assertRaises(ValueError) as context:
            self.katib_client.tune(
                name="experiment",
                model_provider_parameters=Mock(),
                objective=lambda x: x,
            )

        self.assertIn("Invalid configuration", str(context.exception))

    def test_tune_invalid_with_dataset_provider_and_objective(self):
        with self.assertRaises(ValueError) as context:
            self.katib_client.tune(
                name="experiment",
                dataset_provider_parameters=Mock(),
                objective=lambda x: x,
            )

        self.assertIn("Invalid configuration", str(context.exception))

    def test_tune_invalid_with_trainer_parameters_and_objective(self):
        with self.assertRaises(ValueError) as context:
            self.katib_client.tune(
                name="experiment",
                trainer_parameters=Mock(),
                objective=lambda x: x,
            )

        self.assertIn("Invalid configuration", str(context.exception))

    def test_tune_invalid_with_model_provider_and_parameters(self):
        with self.assertRaises(ValueError) as context:
            self.katib_client.tune(
                name="experiment",
                model_provider_parameters=Mock(),
                parameters={"lr": Mock()},
            )

        self.assertIn("Invalid configuration", str(context.exception))

    def test_tune_invalid_with_dataset_provider_and_parameters(self):
        with self.assertRaises(ValueError) as context:
            self.katib_client.tune(
                name="experiment",
                dataset_provider_parameters=Mock(),
                parameters={"lr": Mock()},
            )

        self.assertIn("Invalid configuration", str(context.exception))

    def test_tune_invalid_with_trainer_parameters_and_parameters(self):
        with self.assertRaises(ValueError) as context:
            self.katib_client.tune(
                name="experiment",
                trainer_parameters=Mock(),
                parameters={"lr": Mock()},
            )

        self.assertIn("Invalid configuration", str(context.exception))

    # Case 2: Missing parameters when choosing one option
    def test_tune_invalid_with_only_model_provider(self):
        with self.assertRaises(ValueError) as context:
            self.katib_client.tune(
                name="experiment",
                model_provider_parameters=Mock(),
            )

        self.assertIn("One of the required parameters is None", str(context.exception))

    def test_tune_invalid_with_only_dataset_provider(self):
        with self.assertRaises(ValueError) as context:
            self.katib_client.tune(
                name="experiment",
                dataset_provider_parameters=Mock(),
            )

        self.assertIn("One of the required parameters is None", str(context.exception))

    def test_tune_invalid_with_only_trainer_parameters(self):
        with self.assertRaises(ValueError) as context:
            self.katib_client.tune(
                name="experiment",
                trainer_parameters=Mock(),
            )

        self.assertIn("One of the required parameters is None", str(context.exception))

    def test_tune_invalid_with_only_objective(self):
        with self.assertRaises(ValueError) as context:
            self.katib_client.tune(
                name="experiment",
                objective=lambda x: x,
            )

        self.assertIn("One of the required parameters is None", str(context.exception))

    def test_tune_invalid_with_only_parameters(self):
        with self.assertRaises(ValueError) as context:
            self.katib_client.tune(
                name="experiment",
                parameters={"lr": Mock()},
            )

        self.assertIn("One of the required parameters is None", str(context.exception))

    # Case 3: No parameters provided
    def test_tune_no_parameters(self):
        with self.assertRaises(ValueError) as context:
            self.katib_client.tune(name="experiment")

        self.assertIn("Invalid configuration", str(context.exception))

    # Test for invalid parameters
    # Case 1: Invalid env_per_trial
    def test_tune_invalid_env_per_trial(self):
        with self.assertRaises(ValueError) as context:
            self.katib_client.tune(
                name="experiment",
                objective=lambda x: x,
                parameters={
                    "a": katib.search.int(min=10, max=100),
                    "b": katib.search.double(min=0.1, max=0.2),
                },
                env_per_trial=[123],  # Invalid type
            )

        self.assertIn("Incorrect value for env_per_trial", str(context.exception))

    # Case 2: Invalid resources_per_trial.num_workers (for distributed training)
    def test_tune_invalid_resources_per_trial_value(self):
        with self.assertRaises(ValueError) as context:
            self.katib_client.tune(
                name="experiment",
                objective=lambda x: x,
                parameters={
                    "a": katib.search.int(min=10, max=100),
                    "b": katib.search.double(min=0.1, max=0.2),
                },
                resources_per_trial=katib.TrainerResources(
                    num_workers=0,  # Invalid value, should be at least 1
                    num_procs_per_worker=1,
                    resources_per_worker={"cpu": "1", "memory": "1Gi"},
                ),
            )

        self.assertIn(
            "At least one Worker for PyTorchJob must be set", str(context.exception)
        )

    # Case 3: Invalid model_provider_parameters
    def test_tune_invalid_model_provider_parameters(self):
        with self.assertRaises(ValueError) as context:
            self.katib_client.tune(
                name="experiment",
                model_provider_parameters=123,  # Invalid type, should be an instance of HuggingFaceModelParams
                dataset_provider_parameters=HuggingFaceDatasetParams(
                    repo_id="yelp_review_full",
                    split="train[:8]",
                ),
                trainer_parameters=HuggingFaceTrainerParams(
                    training_parameters=transformers.TrainingArguments(
                        output_dir="test_tune_api",
                        save_strategy="no",
                        learning_rate=katib.search.double(min=1e-05, max=5e-05),
                        num_train_epochs=1,
                        logging_dir="test_tune_api/logs",
                    ),
                    # Set LoRA config to reduce number of trainable model parameters.
                    lora_config=LoraConfig(
                        r=katib.search.int(min=8, max=32),
                        lora_alpha=8,
                        lora_dropout=0.1,
                        bias="none",
                    ),
                ),
            )

        self.assertIn(
            "Model provider parameters must be an instance of HuggingFaceModelParams",
            str(context.exception),
        )

    # Case 4: Invalid dataset_provider_parameters
    def test_tune_invalid_dataset_provider_parameters(self):
        with self.assertRaises(ValueError) as context:
            self.katib_client.tune(
                name="experiment",
                model_provider_parameters=HuggingFaceModelParams(
                    model_uri="hf://google-bert/bert-base-cased",
                    transformer_type=transformers.AutoModelForSequenceClassification,
                    num_labels=5,
                ),
                dataset_provider_parameters=123,  # Invalid type, should be an instance of HuggingFaceDatasetParameters or S3DatasetParams
                trainer_parameters=HuggingFaceTrainerParams(
                    training_parameters=transformers.TrainingArguments(
                        output_dir="test_tune_api",
                        save_strategy="no",
                        learning_rate=katib.search.double(min=1e-05, max=5e-05),
                        num_train_epochs=1,
                        logging_dir="test_tune_api/logs",
                    ),
                    # Set LoRA config to reduce number of trainable model parameters.
                    lora_config=LoraConfig(
                        r=katib.search.int(min=8, max=32),
                        lora_alpha=8,
                        lora_dropout=0.1,
                        bias="none",
                    ),
                ),
            )

        self.assertIn(
            "Dataset provider parameters must be an instance of S3DatasetParams or HuggingFaceDatasetParams",
            str(context.exception),
        )

    # Case 5: Invalid trainer_parameters.training_parameters
    def test_tune_invalid_trainer_parameters_training_parameters(self):
        with self.assertRaises(TypeError) as context:
            self.katib_client.tune(
                name="experiment",
                model_provider_parameters=HuggingFaceModelParams(
                    model_uri="hf://google-bert/bert-base-cased",
                    transformer_type=transformers.AutoModelForSequenceClassification,
                    num_labels=5,
                ),
                dataset_provider_parameters=HuggingFaceDatasetParams(
                    repo_id="yelp_review_full",
                    split="train[:8]",
                ),
                trainer_parameters=HuggingFaceTrainerParams(
                    training_parameters=transformers.TrainingArguments(
                        output_dir="test_tune_api",
                        not_a_valid_parameter="no",
                    ),
                    lora_config=LoraConfig(),
                ),
            )

        self.assertIn(
            "TrainingArguments.__init__() got an unexpected keyword argument",
            str(context.exception),
        )

    # Case 6: Invalid trainer_parameters.lora_config
    def test_tune_invalid_trainer_parameters_lora_config(self):
        with self.assertRaises(TypeError) as context:
            self.katib_client.tune(
                name="experiment",
                model_provider_parameters=HuggingFaceModelParams(
                    model_uri="hf://google-bert/bert-base-cased",
                    transformer_type=transformers.AutoModelForSequenceClassification,
                ),
                dataset_provider_parameters=HuggingFaceDatasetParams(
                    repo_id="yelp_review_full",
                    split="train[:8]",
                ),
                trainer_parameters=HuggingFaceTrainerParams(
                    training_parameters=transformers.TrainingArguments(
                        output_dir="test_tune_api",
                    ),
                    lora_config=LoraConfig(
                        not_a_valid_parameter="no",
                    ),
                ),
            )

        self.assertIn(
            "LoraConfig.__init__() got an unexpected keyword argument",
            str(context.exception),
        )

    # Test functionality
    # Test PVC creation
    # Case 1: PVC successfully created
    @patch("kubernetes.client.CoreV1Api.create_namespaced_persistent_volume_claim")
    @patch("kubernetes.client.CoreV1Api.list_namespaced_persistent_volume_claim")
    @patch("kubeflow.katib.KatibClient.create_experiment")
    def test_pvc_creation(self, mock_create_experiment, mock_list_pvc, mock_create_pvc):
        mock_create_pvc.return_value = Mock()
        mock_list_pvc.return_value = Mock(items=[])
        mock_create_experiment.return_value = Mock()

        exp_name = "experiment"
        storage_config = {
            "size": "10Gi",
            "access_modes": ["ReadWriteOnce"],
        }
        self.katib_client.tune(
            name=exp_name,
            # BERT model URI and type of Transformer to train it.
            model_provider_parameters=HuggingFaceModelParams(
                model_uri="hf://google-bert/bert-base-cased",
                transformer_type=transformers.AutoModelForSequenceClassification,
            ),
            # Use 3000 samples from Yelp dataset.
            dataset_provider_parameters=HuggingFaceDatasetParams(
                repo_id="yelp_review_full",
                split="train[:8]",
            ),
            # Specify HuggingFace Trainer parameters.
            trainer_parameters=HuggingFaceTrainerParams(
                training_parameters=transformers.TrainingArguments(
                    output_dir="test_tune_api",
                    save_strategy="no",
                    learning_rate=katib.search.double(min=1e-05, max=5e-05),
                    num_train_epochs=1,
                    logging_dir="test_tune_api/logs",
                ),
                # Set LoRA config to reduce number of trainable model parameters.
                lora_config=LoraConfig(
                    r=katib.search.int(min=8, max=32),
                    lora_alpha=8,
                    lora_dropout=0.1,
                    bias="none",
                ),
            ),
            objective_metric_name="accuracy",
            storage_config=storage_config,
        )

        expected_pvc_spec = models.V1PersistentVolumeClaim(
            api_version="v1",
            kind="PersistentVolumeClaim",
            metadata={"name": exp_name, "namespace": "default"},
            spec=models.V1PersistentVolumeClaimSpec(
                access_modes=storage_config["access_modes"],
                resources=models.V1ResourceRequirements(
                    requests={"storage": storage_config["size"]}
                ),
            ),
        )

        mock_create_pvc.assert_called_once_with(
            namespace="default", body=expected_pvc_spec
        )

    # Case 2: PVC already exists
    @patch("kubernetes.client.CoreV1Api.create_namespaced_persistent_volume_claim")
    @patch("kubernetes.client.CoreV1Api.list_namespaced_persistent_volume_claim")
    @patch("kubeflow.katib.KatibClient.create_experiment")
    def test_pvc_creation_with_existing_pvc(
        self, mock_create_experiment, mock_list_pvc, mock_create_pvc
    ):
        # Simulate an ApiException being raised when trying to create a PVC
        mock_create_pvc.side_effect = ApiException(status=409, reason="Already exists")

        # Simulate existing PVC in the list
        mock_existing_pvc = Mock()
        mock_existing_pvc.metadata.name = "test-pvc"
        mock_list_pvc.return_value = Mock(items=[mock_existing_pvc])

        mock_create_experiment.return_value = Mock()

        exp_name = "test-pvc"
        storage_config = {
            "size": "10Gi",
            "access_modes": ["ReadWriteOnce"],
        }
        self.katib_client.tune(
            name=exp_name,
            # BERT model URI and type of Transformer to train it.
            model_provider_parameters=HuggingFaceModelParams(
                model_uri="hf://google-bert/bert-base-cased",
                transformer_type=transformers.AutoModelForSequenceClassification,
            ),
            # Use 3000 samples from Yelp dataset.
            dataset_provider_parameters=HuggingFaceDatasetParams(
                repo_id="yelp_review_full",
                split="train[:8]",
            ),
            # Specify HuggingFace Trainer parameters.
            trainer_parameters=HuggingFaceTrainerParams(
                training_parameters=transformers.TrainingArguments(
                    output_dir="test_tune_api",
                    save_strategy="no",
                    learning_rate=katib.search.double(min=1e-05, max=5e-05),
                    num_train_epochs=1,
                    logging_dir="test_tune_api/logs",
                ),
                # Set LoRA config to reduce number of trainable model parameters.
                lora_config=LoraConfig(
                    r=katib.search.int(min=8, max=32),
                    lora_alpha=8,
                    lora_dropout=0.1,
                    bias="none",
                ),
            ),
            objective_metric_name="accuracy",
            storage_config=storage_config,
        )

        # Assert that create_namespaced_persistent_volume_claim was called once
        mock_create_pvc.assert_called_once()

        # Assert that list_namespaced_persistent_volume_claim was called to check existing PVCs
        mock_list_pvc.assert_called_once_with("default")

        # Ensure no exception is raised since the PVC already exists
        self.assertTrue(mock_list_pvc.return_value.items[0].metadata.name == exp_name)

    # Case 3: PVC creation fails
    @patch("kubernetes.client.CoreV1Api.create_namespaced_persistent_volume_claim")
    @patch("kubernetes.client.CoreV1Api.list_namespaced_persistent_volume_claim")
    @patch("kubeflow.katib.KatibClient.create_experiment")
    def test_pvc_creation_fails(
        self, mock_create_experiment, mock_list_pvc, mock_create_pvc
    ):
        # Simulate an ApiException being raised when trying to create a PVC
        mock_create_pvc.side_effect = ApiException(
            status=500, reason="Internal Server Error"
        )

        # Simulate no existing PVC in the list
        mock_list_pvc.return_value = Mock(items=[])

        mock_create_experiment.return_value = Mock()

        exp_name = "test-pvc"
        storage_config = {
            "size": "10Gi",
            "access_modes": ["ReadWriteOnce"],
        }
        with self.assertRaises(RuntimeError) as context:
            self.katib_client.tune(
                name=exp_name,
                # BERT model URI and type of Transformer to train it.
                model_provider_parameters=HuggingFaceModelParams(
                    model_uri="hf://google-bert/bert-base-cased",
                    transformer_type=transformers.AutoModelForSequenceClassification,
                ),
                # Use 3000 samples from Yelp dataset.
                dataset_provider_parameters=HuggingFaceDatasetParams(
                    repo_id="yelp_review_full",
                    split="train[:8]",
                ),
                # Specify HuggingFace Trainer parameters.
                trainer_parameters=HuggingFaceTrainerParams(
                    training_parameters=transformers.TrainingArguments(
                        output_dir="test_tune_api",
                        save_strategy="no",
                        learning_rate=katib.search.double(min=1e-05, max=5e-05),
                        num_train_epochs=1,
                        logging_dir="test_tune_api/logs",
                    ),
                    # Set LoRA config to reduce number of trainable model parameters.
                    lora_config=LoraConfig(
                        r=katib.search.int(min=8, max=32),
                        lora_alpha=8,
                        lora_dropout=0.1,
                        bias="none",
                    ),
                ),
                objective_metric_name="accuracy",
                storage_config=storage_config,
            )

        # Assert that the appropriate error message is raised
        self.assertIn("failed to create PVC", str(context.exception))

        # Assert that create_namespaced_persistent_volume_claim was called once
        mock_create_pvc.assert_called_once()

        # Assert that list_namespaced_persistent_volume_claim was called once
        mock_list_pvc.assert_called_once_with("default")

    # Test container, pod, job/pytorchjob, trial template, and experiment creation
    # Case 1: Custom objective
    @patch("kubeflow.katib.KatibClient.create_experiment")
    def test_experiment_creation_with_custom_objective(
        self, mock_create_experiment
    ):
        self.katib_client.tune(
            name="experiment",
            objective=lambda x: x,
            parameters={
                "a": katib.search.int(min=10, max=100),
                "b": katib.search.double(min=0.1, max=0.2),
            },
            objective_metric_name="accuracy",
            objective_goal=0.9,
            max_trial_count=10,
            parallel_trial_count=2,
            max_failed_trial_count=1,
            resources_per_trial={"cpu": "1", "memory": "1Gi"},
        )

        mock_create_experiment.assert_called_once()
        args, kwargs = mock_create_experiment.call_args
        experiment = args[0]

        expected_container = [
            models.V1Container(
                name="training-container",
                image="docker.io/tensorflow/tensorflow:2.13.0",
                command=["bash", "-c"],
                args=[
                    "\n"
                    "program_path=$(mktemp -d)\n"
                    "read -r -d '' SCRIPT << EOM\n"
                    "\n"
                    "objective=lambda x: x,\n"
                    "\n"
                    "<lambda>({'a': '${trialParameters.a}', 'b': '${trialParameters.b}'})\n"
                    "\n"
                    "EOM\n"
                    'printf "%s" "$SCRIPT" > "$program_path/ephemeral_script.py"\n'
                    'python3 -u "$program_path/ephemeral_script.py"'
                ],
                resources=models.V1ResourceRequirements(
                    requests={"cpu": "1", "memory": "1Gi"},
                    limits={"cpu": "1", "memory": "1Gi"},
                ),
            )
        ]

        expected_pod = models.V1PodTemplateSpec(
            metadata=models.V1ObjectMeta(
                annotations={"sidecar.istio.io/inject": "false"}
            ),
            spec=models.V1PodSpec(
                containers=expected_container,
                restart_policy="Never",
            ),
        )

        expected_job = client.V1Job(
            api_version="batch/v1",
            kind="Job",
            spec=client.V1JobSpec(
                template=expected_pod,
            ),
        )

        expected_trial_template = models.V1beta1TrialTemplate(
            primary_container_name="training-container",
            trial_parameters=[
                models.V1beta1TrialParameterSpec(name="a", reference="a"),
                models.V1beta1TrialParameterSpec(name="b", reference="b"),
            ],
            retain=False,
            trial_spec=expected_job,
        )

        expected_parameters = [
            models.V1beta1ParameterSpec(
                name="a",
                parameter_type="int",
                feasible_space=models.V1beta1FeasibleSpace(min="10", max="100"),
            ),
            models.V1beta1ParameterSpec(
                name="b",
                parameter_type="double",
                feasible_space=models.V1beta1FeasibleSpace(min="0.1", max="0.2"),
            ),
        ]

        self.assertEqual(experiment.spec.objective.type, "maximize")
        self.assertEqual(experiment.spec.objective.objective_metric_name, "accuracy")
        self.assertEqual(experiment.spec.objective.goal, 0.9)
        self.assertEqual(experiment.spec.algorithm.algorithm_name, "random")
        self.assertEqual(experiment.spec.max_trial_count, 10)
        self.assertEqual(experiment.spec.parallel_trial_count, 2)
        self.assertEqual(experiment.spec.max_failed_trial_count, 1)
        self.assertEqual(experiment.spec.parameters, expected_parameters)
        self.assertEqual(
            experiment.spec.trial_template.trial_spec.spec.template.spec.containers,
            expected_container,
        )
        self.assertEqual(
            experiment.spec.trial_template.trial_spec.spec.template, expected_pod
        )
        self.assertEqual(experiment.spec.trial_template.trial_spec, expected_job)
        self.assertEqual(experiment.spec.trial_template, expected_trial_template)

    # Case 2: External models and datasets
    @patch("kubeflow.katib.KatibClient.create_experiment")
    def test_experiment_creation_with_external_model(
        self, mock_create_experiment
    ):
        exp_name = "experiment"
        self.katib_client.tune(
            name=exp_name,
            # BERT model URI and type of Transformer to train it.
            model_provider_parameters=HuggingFaceModelParams(
                model_uri="hf://google-bert/bert-base-cased",
                transformer_type=transformers.AutoModelForSequenceClassification,
                num_labels=5,
            ),
            # Use 3000 samples from Yelp dataset.
            dataset_provider_parameters=HuggingFaceDatasetParams(
                repo_id="yelp_review_full",
                split="train[:8]",
            ),
            # Specify HuggingFace Trainer parameters.
            trainer_parameters=HuggingFaceTrainerParams(
                training_parameters=transformers.TrainingArguments(
                    output_dir="test_tune_api",
                    save_strategy="no",
                    learning_rate=katib.search.double(min=1e-05, max=5e-05),
                    num_train_epochs=1,
                    logging_dir="test_tune_api/logs",
                ),
                # Set LoRA config to reduce number of trainable model parameters.
                lora_config=LoraConfig(
                    r=katib.search.int(min=8, max=32),
                    lora_alpha=8,
                    lora_dropout=0.1,
                    bias="none",
                ),
            ),
            objective_metric_name="accuracy",
            objective_goal=0.9,
            max_trial_count=10,
            parallel_trial_count=2,
            max_failed_trial_count=1,
            resources_per_trial=katib.TrainerResources(
                num_workers=3,
                num_procs_per_worker=1,
                resources_per_worker={"cpu": "1", "memory": "1Gi"},
            ),
        )

        mock_create_experiment.assert_called_once()
        args, kwargs = mock_create_experiment.call_args
        experiment = args[0]

        expected_init_container = [
            models.V1Container(
                name="storage-initializer",
                image="docker.io/kubeflow/storage-initializer",
                args=[
                    "--model_provider",
                    "hf",
                    "--model_provider_parameters",
                    '{"model_uri": "hf://google-bert/bert-base-cased", "transformer_type": "AutoModelForSequenceClassification", '
                    '"access_token": null, "num_labels": 5}',
                    "--dataset_provider",
                    "hf",
                    "--dataset_provider_parameters",
                    '{"repo_id": "yelp_review_full", "access_token": null, "split": "train[:8]"}',
                ],
                volume_mounts=[
                    training_models.V1VolumeMount(
                        name="storage-initializer",
                        mount_path="/workspace",
                    )
                ],
            )
        ]

        expected_container = [
            models.V1Container(
                name="pytorch",
                image="docker.io/kubeflow/trainer-huggingface",
                args=[
                    "--model_uri",
                    "hf://google-bert/bert-base-cased",
                    "--transformer_type",
                    "AutoModelForSequenceClassification",
                    "--num_labels",
                    "5",
                    "--model_dir",
                    "/workspace/model",
                    "--dataset_dir",
                    "/workspace/dataset",
                    "--lora_config",
                    '\'{"peft_type": "LORA", "base_model_name_or_path": null, "task_type": null, '
                    '"inference_mode": false, "r": "${trialParameters.r}", "target_modules": null, '
                    '"lora_alpha": 8, "lora_dropout": 0.1, "fan_in_fan_out": false, "bias": "none", '
                    '"modules_to_save": null, "init_lora_weights": true}\'',
                    "--training_parameters",
                    '\'{"output_dir": "test_tune_api", "overwrite_output_dir": false, "do_train": '
                    'false, "do_eval": false, "do_predict": false, "evaluation_strategy": "no", '
                    '"prediction_loss_only": false, "per_device_train_batch_size": 8, '
                    '"per_device_eval_batch_size": 8, "per_gpu_train_batch_size": null, '
                    '"per_gpu_eval_batch_size": null, "gradient_accumulation_steps": 1, '
                    '"eval_accumulation_steps": null, "eval_delay": 0, "learning_rate": '
                    '"${trialParameters.learning_rate}", "weight_decay": 0.0, "adam_beta1": 0.9, '
                    '"adam_beta2": 0.999, "adam_epsilon": 1e-08, "max_grad_norm": 1.0, '
                    '"num_train_epochs": 1, "max_steps": -1, "lr_scheduler_type": "linear", '
                    '"lr_scheduler_kwargs": {}, "warmup_ratio": 0.0, "warmup_steps": 0, '
                    '"log_level": "passive", "log_level_replica": "warning", "log_on_each_node": '
                    'true, "logging_dir": "test_tune_api/logs", "logging_strategy": "steps", '
                    '"logging_first_step": false, "logging_steps": 500, "logging_nan_inf_filter": '
                    'true, "save_strategy": "no", "save_steps": 500, "save_total_limit": null, '
                    '"save_safetensors": true, "save_on_each_node": false, "save_only_model": '
                    'false, "no_cuda": false, "use_cpu": false, "use_mps_device": false, "seed": '
                    '42, "data_seed": null, "jit_mode_eval": false, "use_ipex": false, "bf16": '
                    'false, "fp16": false, "fp16_opt_level": "O1", "half_precision_backend": '
                    '"auto", "bf16_full_eval": false, "fp16_full_eval": false, "tf32": null, '
                    '"local_rank": 0, "ddp_backend": null, "tpu_num_cores": null, '
                    '"tpu_metrics_debug": false, "debug": [], "dataloader_drop_last": false, '
                    '"eval_steps": null, "dataloader_num_workers": 0, "dataloader_prefetch_factor": '
                    'null, "past_index": -1, "run_name": "test_tune_api", "disable_tqdm": false, '
                    '"remove_unused_columns": true, "label_names": null, "load_best_model_at_end": '
                    'false, "metric_for_best_model": null, "greater_is_better": null, '
                    '"ignore_data_skip": false, "fsdp": [], "fsdp_min_num_params": 0, '
                    '"fsdp_config": {"min_num_params": 0, "xla": false, "xla_fsdp_v2": false, '
                    '"xla_fsdp_grad_ckpt": false}, "fsdp_transformer_layer_cls_to_wrap": null, '
                    '"accelerator_config": {"split_batches": false, "dispatch_batches": null, '
                    '"even_batches": true, "use_seedable_sampler": true}, "deepspeed": null, '
                    '"label_smoothing_factor": 0.0, "optim": "adamw_torch", "optim_args": null, '
                    '"adafactor": false, "group_by_length": false, "length_column_name": "length", '
                    '"report_to": ["tensorboard"], "ddp_find_unused_parameters": null, '
                    '"ddp_bucket_cap_mb": null, "ddp_broadcast_buffers": null, '
                    '"dataloader_pin_memory": true, "dataloader_persistent_workers": false, '
                    '"skip_memory_metrics": true, "use_legacy_prediction_loop": false, '
                    '"push_to_hub": false, "resume_from_checkpoint": null, "hub_model_id": null, '
                    '"hub_strategy": "every_save", "hub_token": "<HUB_TOKEN>", "hub_private_repo": '
                    'false, "hub_always_push": false, "gradient_checkpointing": false, '
                    '"gradient_checkpointing_kwargs": null, "include_inputs_for_metrics": false, '
                    '"fp16_backend": "auto", "push_to_hub_model_id": null, '
                    '"push_to_hub_organization": null, "push_to_hub_token": "<PUSH_TO_HUB_TOKEN>", '
                    '"mp_parameters": "", "auto_find_batch_size": false, "full_determinism": '
                    'false, "torchdynamo": null, "ray_scope": "last", "ddp_timeout": 1800, '
                    '"torch_compile": false, "torch_compile_backend": null, "torch_compile_mode": '
                    'null, "dispatch_batches": null, "split_batches": null, "include_tokens_per_'
                    'second": false, "include_num_input_tokens_seen": false, '
                    '"neftune_noise_alpha": null}\'',
                ],
                resources=models.V1ResourceRequirements(
                    requests={"cpu": "1", "memory": "1Gi"},
                    limits={"cpu": "1", "memory": "1Gi"},
                ),
                volume_mounts=[
                    training_models.V1VolumeMount(
                        name="storage-initializer",
                        mount_path="/workspace",
                    )
                ],
            )
        ]

        expected_master_pod = models.V1PodTemplateSpec(
            metadata=models.V1ObjectMeta(
                annotations={"sidecar.istio.io/inject": "false"}
            ),
            spec=models.V1PodSpec(
                init_containers=expected_init_container,
                containers=expected_container,
                volumes=[
                    models.V1Volume(
                        name="storage-initializer",
                        persistent_volume_claim=models.V1PersistentVolumeClaimVolumeSource(
                            claim_name=exp_name
                        ),
                    )
                ],
            ),
        )

        expected_worker_pod = models.V1PodTemplateSpec(
            metadata=models.V1ObjectMeta(
                annotations={"sidecar.istio.io/inject": "false"}
            ),
            spec=models.V1PodSpec(
                containers=expected_container,
                volumes=[
                    models.V1Volume(
                        name="storage-initializer",
                        persistent_volume_claim=models.V1PersistentVolumeClaimVolumeSource(
                            claim_name=exp_name
                        ),
                    )
                ],
            ),
        )

        expected_job = training_models.KubeflowOrgV1PyTorchJob(
            api_version="kubeflow.org/v1",
            kind="PyTorchJob",
            spec=training_models.KubeflowOrgV1PyTorchJobSpec(
                run_policy=training_models.KubeflowOrgV1RunPolicy(
                    clean_pod_policy=None
                ),
                pytorch_replica_specs={
                    "Master": training_models.KubeflowOrgV1ReplicaSpec(
                        replicas=1,
                        template=expected_master_pod,
                    ),
                    "Worker": training_models.KubeflowOrgV1ReplicaSpec(
                        replicas=2,
                        template=expected_worker_pod,
                    ),
                },
                nproc_per_node="1",
            ),
        )

        expected_trial_template = models.V1beta1TrialTemplate(
            primary_container_name="pytorch",
            trial_parameters=[
                models.V1beta1TrialParameterSpec(
                    name="learning_rate", reference="learning_rate"
                ),
                models.V1beta1TrialParameterSpec(name="r", reference="r"),
            ],
            retain=False,
            trial_spec=expected_job,
        )

        expected_parameters = [
            models.V1beta1ParameterSpec(
                name="learning_rate",
                parameter_type="double",
                feasible_space=models.V1beta1FeasibleSpace(min="1e-05", max="5e-05"),
            ),
            models.V1beta1ParameterSpec(
                name="r",
                parameter_type="int",
                feasible_space=models.V1beta1FeasibleSpace(min="8", max="32"),
            ),
        ]

        self.assertEqual(experiment.spec.objective.type, "maximize")
        self.assertEqual(experiment.spec.objective.objective_metric_name, "accuracy")
        self.assertEqual(experiment.spec.objective.goal, 0.9)
        self.assertEqual(experiment.spec.algorithm.algorithm_name, "random")
        self.assertEqual(experiment.spec.max_trial_count, 10)
        self.assertEqual(experiment.spec.parallel_trial_count, 2)
        self.assertEqual(experiment.spec.max_failed_trial_count, 1)
        self.assertEqual(experiment.spec.parameters, expected_parameters)
        self.assertEqual(
            experiment.spec.trial_template.trial_spec.spec.pytorch_replica_specs[
                "Master"
            ].template.spec.init_containers,
            expected_init_container,
        )
        self.assertEqual(
            experiment.spec.trial_template.trial_spec.spec.pytorch_replica_specs[
                "Master"
            ].template.spec.containers,
            expected_container,
        )
        self.assertEqual(
            experiment.spec.trial_template.trial_spec.spec.pytorch_replica_specs[
                "Master"
            ].replicas,
            1,
        )
        self.assertEqual(
            experiment.spec.trial_template.trial_spec.spec.pytorch_replica_specs[
                "Master"
            ].template,
            expected_master_pod,
        )
        self.assertEqual(
            experiment.spec.trial_template.trial_spec.spec.pytorch_replica_specs[
                "Worker"
            ].template.spec.containers,
            expected_container,
        )
        self.assertEqual(
            experiment.spec.trial_template.trial_spec.spec.pytorch_replica_specs[
                "Worker"
            ].replicas,
            2,
        )
        self.assertEqual(
            experiment.spec.trial_template.trial_spec.spec.pytorch_replica_specs[
                "Worker"
            ].template,
            expected_worker_pod,
        )
        self.assertEqual(experiment.spec.trial_template.trial_spec, expected_job)
        self.assertEqual(experiment.spec.trial_template, expected_trial_template)

if __name__ == "__main__":
    unittest.main()
