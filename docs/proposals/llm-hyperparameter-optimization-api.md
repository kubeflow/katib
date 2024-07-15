# HyperParameter Optimization API for LLM Fine-Tuning

- [HyperParameter Optimization API for LLM Fine-Tuning](#hyperparameter-optimization-api-for-llm-fine-tuning)
  * [Links](#links)
  * [Motivation](#motivation)
  * [Goals](#goals)
  * [Non-Goals](#non-goals)
  * [Design for API](#design-for-api)
    + [Example](#example)
  * [Implementation](#implementation)

## Links

- [katib/issues#2291 (Tuning API in Katib for LLMs)](https://github.com/kubeflow/katib/issues/2291)

## Motivation

The rapid advancements and growing popularity of Large Language Models (LLMs) have driven an increased need for effective LLMOps in Kubernetes environments. To address this, we developed a [train API](https://www.kubeflow.org/docs/components/training/user-guides/fine-tuning/) within the Training Python SDK, simplifying the process of fine-tuning LLMs using distributed PyTorchJob workers. However, hyperparameter optimization remains a crucial yet labor-intensive task for enhancing model performance. Automating this tuning process through a dedicated API will facilitate efficient and scalable exploration of hyperparameters, ultimately improving model performance and reducing manual effort.

## Goals

Our goal is to develop a high-level API for tuning hyperparameters of LLMs that simplifies the process of hyperparameter optimization in Kubernetes. This API will seamlessly integrate with external platforms like HuggingFace and S3 for importing pretrained models and datasets. By specifying parameters for the training objective, trial configurations, and PyTorch worker configurations, the API will automate the creation of experiments and execution of trials. This abstraction of Kubernetes infrastructure complexities will enable data scientists to optimize hyperparameters efficiently and effectively.

## Non-Goals

1. Incorporate early stopping strategy into the API to optimize training efficiency.
2. Expand support for distributed training in frameworks beyond PyTorch by leveraging their distributed training capabilities.
3. Support adding custom providers through configmap or CRD approach to enhance flexibility.
4. Enable users to deploy tuned models for inference within their applications or seamlessly integrate them into existing NLP pipelines for specialized tasks.

## Design for API

![Design for API](images/design_api.jpg)

```python
import kubeflow.katib as katib
from kubeflow.katib import KatibClient

class KatibClient(object):
	
	def tune(
		self,
		name: str, 
		namespace: Optional[str] = None,
		model_provider_parameters: Optional[HuggingFaceModelParams] = None,
		dataset_provider_parameters: Optional[Union[HuggingFaceDatasetParams, S3DatasetParams]] = None,
		trainer_parameters: Union[HuggingFaceTrainerParams, Dict[str, Any]] = None,
		storage_config: Dict[str, Optional[Union[str, List[str]]]] = {
            "size": constants.PVC_DEFAULT_SIZE,
            "storage_class": None,
            "access_modes": constants.PVC_DEFAULT_ACCESS_MODES,
        },
		objective: Optional[Callable] = None,
		base_image: Optional[str] = None,
		algorithm_name: str = "random",
		algorithm_settings: Union[dict, List[models.V1beta1AlgorithmSetting], None] = None,
		objective_metric_name: str = "eval_accuracy",
		additional_metric_names: List[str] = [],
		objective_type: str = "maximize",
		objective_goal: float = None,
		max_trial_count: int = None,
		parallel_trial_count: int = None,
		max_failed_trial_count: int = None,
		resources_per_trial = Union[dict, client.V1ResourceRequirements, types.TrainerResources, None] = None,
		retain_trials: bool = False,
		env_per_trial: Optional[Union[Dict[str, str], List[Union[client.V1EnvVar, client.V1EnvFromSource]]]] = None,
		packages_to_install: List[str] = None,
		pip_index_url: str = "https://pypi.org/simple",
	):
		"""
        Initiates a hyperparameter tuning experiment in Katib.
		Model, dataset and parameters can be configured using one of the following options:
		- Using the Storage Initializer: Specify `model_provider_parameters`, `dataset_provider_parameters`, and `trainer_parameters`. This option downloads models and datasets from external platforms like HuggingFace and S3, and utilizes `Trainer.train()` in HuggingFace to train the model.
		- Defining a custom objective function: Specify the `objective` parameter to define your own objective function, and use the `base_image` parameter to execute the objective function.

        Parameters:
		- name: Name for the experiment.
		- namespace: Namespace for the experiment. Defaults to the namespace of the 'KatibClient' object.
		- model_provider_parameters: Parameters for providing the model. Compatible with model providers like HuggingFace.
		- dataset_provider_parameters: Parameters for providing the dataset. Compatible with dataset providers like HuggingFace or S3.
		- trainer_parameters: Parameters for configuring the training process, including settings for hyperparameters search space.
		- storage_config: Configuration for Storage Initializer PVC to download pre-trained model and dataset.
		- objective: Objective function that Katib uses to train the model.
		- base_image: Image to use when executing the objective function.
		- algorithm_name: Tuning algorithm name (e.g., 'random', 'bayesian').
		- algorithm_settings: Settings for the tuning algorithm.
		- objective_metric_name: Primary metric to optimize.
		- additional_metric_names: List of additional metrics to collect.
		- objective_type: Optimization direction for the objective metric, "minimize" or "maximize".
		- objective_goal: Desired value of the objective metric.
		- max_trial_count: Maximum number of trials to run.
		- parallel_trial_count: Number of trials to run in parallel.
		- max_failed_trial_count: Maximum number of allowed failed trials.
		- resources_per_trial: Resources assigned to per trial, which can be specified using one of the following options:
			- Non-distributed Training: Specify a kubernetes.client.V1ResourceRequirements object or a dicitionary that includes one or more of the following keys: `cpu`, `memory`, or `gpu` (other keys will be ignored).
			- Distributed Training in Pytorch: Specify a types.TrainerResources, which includes the following parameters:
				- num_workers: Number of PyTorchJob workers.
				- num_procs_per_worker: Number of processes per PyTorchJob worker.
				- resources_per_worker: Resources assigned to per PyTorchJob worker container, specified as either a kubernetes.client.V1ResourceRequirements object or a dicitionary that includes one or more of the following keys: `cpu`, `memory`, or `gpu` (other keys will be ignored).
		- retain_trials: Whether to retain trial resources after completion.
		- env_per_trial: Environment variables for worker containers.
		- packages_to_install: Additional Python packages to install.
		- pip_index_url: URL of the PyPI index for installing packages.
        """
        pass  # Implementation logic for initiating the experiment
```

### Example

```python
import transformers
from peft import LoraConfig

from kubeflow.storage_initializer.hugging_face import (
	HuggingFaceModelParams,
	HuggingFaceDatasetParams,
	HuggingFaceTrainerParams,
)

import kubeflow.katib as katib
from kubeflow.katib import KatibClient

# Create a Katib client
katib_client = KatibClient(namespace='kubeflow')

# Run the tuning Experiment
exp_name = "llm-hp-optimization-test"
katib_client.tune(
	name = exp_name,
	# BERT model URI and type of Transformer to train it.
	model_provider_parameters = HuggingFaceModelParams(
		model_uri = "hf://google-bert/bert-base-cased",
		transformer_type = transformers.AutoModelForSequenceClassification,
	),
	# Use 3000 samples from Yelp dataset.
	dataset_provider_parameters = HuggingFaceDatasetParams(
		repo_id = "yelp_review_full",
		split = "train[:3000]",
	),
	# Specify HuggingFace Trainer parameters.
	trainer_parameters = HuggingFaceTrainerParams(
		training_parameters = transformers.TrainingArguments(
			output_dir = "test_trainer",
			save_strategy = "no",
			eval_strategy = "epoch",
			disable_tqdm = True,
			log_level = "info",
			learning_rate = katib.search.double(min=1e-05, max=5e-05),
			per_device_train_batch_size = katib.search.categorical([8, 16, 32]),
			num_train_epochs = katib.search.int(min=1, max=10),
			weight_decay = katib.search.double(min=0.0, max=1.0),
		),
		# Set LoRA config to reduce number of trainable model parameters.
		lora_config = LoraConfig(
			r = katib.search.int(min=8, max=32),
			lora_alpha = 8,
			lora_dropout = 0.1,
			bias = "none",
		),
	),	
	objective_metric_name = "eval_accuracy", 
	objective_type = "maximize", 
	algorithm_name = "random",
	max_trial_count = 50,
	parallel_trial_count = 2,
	resources_per_trial = types.TrainerResources(
		num_workers = 4,
		num_procs_per_worker = 2,
		resources_per_worker = {
			"gpu": 2,
			"cpu": 1,
			"memory": "10G",
		},
	), 
	# For non-distributed training in PyTorch, specify resources like this: resources_per_trial = {"gpu": 2, "cpu": 10, "memory": "20G"}
)

# Get the best hyperparameters
best_hps = katib_client.get_optimal_hyperparameters(exp_name)
```

## Implementation

By passing the specified parameters, this API will automate hyperparameter optimization for LLMs. The implementation will focus on the following aspects:

**Model and Dataset Management**: We will leverage the [storage_initializer](https://github.com/kubeflow/training-operator/tree/master/sdk/python/kubeflow/storage_initializer) from the Training Python SDK for seamless integration of pretrained models and datasets from platforms like HuggingFace and S3. This component manages downloading and storing pretrained models and datasets via a PersistentVolumeClaim (PVC), which is shared across containers, ensuring efficient access to the pretrained model and dataset without redundant downloads.

**Hyperparameter Configuration**: Users specify training parameters and the hyperparameters to be optimized within `trainer_parameters`. The API will first traverse `trainer_parameters.training_parameters` and `trainer_parameters.lora_config` to identify tunable hyperparameters and set up their values for the Experiment and Trials. These parameters are then passed as `args` to the container spec of workers.

```python
# Traverse and set up hyperparameters
input_params = {}
experiment_params = []
trial_params = []

training_args = trainer_parameters.training_parameters
for p_name, p_value in training_args.to_dict().items():
	if not hasattr(training_args, p_name):
		logger.warning(f"Training parameter {p_name} is not supported by the current transformer.")
		continue
	if isinstance(p_value, models.V1beta1ParameterSpec):
		value = f"${{trialParameters.{p_name}}}"
		setattr(training_args, p_name, value)
		p_value.name = p_name
		experiment_params.append(p_value)
		trial_params.append(models.V1beta1TrialParameterSpec(name=p_name, reference=p_name))
	elif p_value is not None:
		value = type(old_attr)(p_value)
		setattr(training_args, p_name, value)
input_params['training_args'] = training_args

# Note: Repeat similar logic for `lora_config`

# create container spec of worker
container_spec = client.V1Container(
	...
	args=[
		"--model_uri",
		model_provider_parameters.model_uri,
		"--transformer_type",
		model_provider_parameters.transformer_type.__name__,
		"--model_dir",
		"REPLACE_WITH_ACTUAL_MODEL_PATH", 
		"--dataset_dir",
		"REPLACE_WITH_ACTUAL_DATASET_PATH",
		"--lora_config",
		json.dumps(input_params['lora_config'].__dict__, cls=utils.SetEncoder),
		"--training_parameters",
		json.dumps(input_params['training_args'].to_dict()),
	],
	...
)
```

**Hyperparameter Optimization**: This API will create an Experiment that defines the search space for identified tunable hyperparameters, the objective metric, optimization algorithm, etc. The Experiment will orchestrate the hyperparameter tuning process, generating Trials for each configuratin. Each Trial will be implemented as a Kubernete PyTorchJob, with the `trialTemplate` specifying the exact values for hyperparameters. The `trialTemplate` will also define master and worker containers, facilitating effective resource distribution and parallel execution of Trials. Trial results will then be fed back to the Experiment, which will evaluate the outcomes to identify the optimal set of hyperparameters.

 **Dependencies Update**: To reuse existing assets from the Training Python SDK and integrate packages from HuggingFace, dependencies will be added to the `setup.py` of the Katib Python SDK as follows:

```python
setuptools.setup(
	...// Configurations of the package
	extras_require={
		"huggingface": ["kubeflow-training[huggingface]==1.8.0rc1"],
	},
)
```
