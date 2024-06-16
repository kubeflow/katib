# HyperParameter Optimization API for LLM Fine-Tuning

- [HyperParameter Optimization API for LLM Fine-Tuning](#hyperparameter-optimization-api-for-llm-fine-tuning)
  * [Links](#links)
  * [Motivation](#motivation)
  * [Goal](#goal)
  * [Design for API](#design-for-api)
    + [Example](#example)
  * [Implementation](#implementation)
  * [Advanced Functionalities](#advanced-functionalities)

## Links

- [katib/issues#2291 (Tuning API in Katib for LLMs)](https://github.com/kubeflow/katib/issues/2291)

## Motivation

The rapid advancements and growing popularity of Large Language Models (LLMs) have driven an increased need for effective LLMOps in Kubernetes environments. To address this, we developed a [train API](https://www.kubeflow.org/docs/components/training/user-guides/fine-tuning/) within the Training Python SDK, simplifying the process of fine-tuning LLMs using distributed PyTorchJob workers. However, hyperparameter optimization remains a crucial yet labor-intensive task for enhancing model performance. Automating this tuning process through a dedicated API will facilitate efficient and scalable exploration of hyperparameters, ultimately improving model performance and reducing manual effort.

## Goal

We aim to develop a high-level API for tuning hyperparameters of LLMs. This API will abstract infrastructure complexities within Kubernetes, enabling data scientists to easily optimize hyperparameters. It will leverage the train API for seamless integration of pretrained models and datasets from platforms like HuggingFace and S3.

## Design for API

![Design for API](../images/design_api.jpg)

```python
import kubeflow.katib as katib
from kubeflow.katib import KatibClient

class KatibClient(object):
	
	def tune(
		name, 
		hp_space,
		objective,
		objective_metric_name,
		objective_type,
		algorithm_name,
		namespace: Optional,
		additional_metric_names: Optional,
		objective_goal: Optional,
		algorithm_setting: Optional,
		max_trail_count: Optional,
		max_failed_trail_count: Optional,
		parallel_trail_count: Optional,
		resources_per_trail: Optional,
		num_workers: Optional,
		num_procs_per_worker: Optional,
		resources_per_worker: Optional,
		retain_trials: Optional,
		env_per_trial, # TBD
		storage_config, # TBD
		base_image, # TBD
		packages_to_install, # TBD
		pip_index_url, # TBD
	):
		"""
        Initiates a hyperparameter tuning experiment in Katib.

        Parameters:
        - name: Name for the experiment.
        - hp_space: Dictionary defining the hyperparameter search space.
        - objective: Objective function to be optimized by Katib.
        - objective_metric_name: Primary metric to optimize.
        - objective_type: Optimization direction for the objective metric, "minimize" or "maximize".
        - algorithm_name: Tuning algorithm name (e.g., 'random', 'bayesian').
        - namespace: Namespace for the experiment. Defaults to the namespace of the 'KatibClient' object.
		- additional_metric_names: List of additional metrics to collect.
		- objective_goal: Desired value of the objective metric.
		- algorithm_settings: Settings for the tuning algorithm.
        - max_trial_count: Maximum number of trials to run.
        - max_failed_trial_count: Maximum number of allowed failed trials.
        - parallel_trial_count: Number of trials to run in parallel.
        - resources_per_trial: Resources required per trial.
        - num_workers: Number of PyTorchJob workers for distributed jobs.
        - num_procs_per_worker: Number of processes per worker for distributed jobs.
        - resources_per_worker: Resources allocated per worker.
        - env_per_trial: Environment variables for worker containers.
        - storage_config: PVC configuration for pre-trained model and dataset storage.
        - retain_trials: Whether to retain trial resources after completion.
        - base_image: Base Docker image for running the objective function.
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
from kubeflow.training import TrainingClient
import kubeflow.katib as katib
from kubeflow.katib import KatibClient

# Create a Katib client
katib_client = KatibClient(namespace='kubeflow')

# Define an objective function leveraging the train API
def objective(exp_name):
	result = TrainingClient().train(
		name=exp_name,
		# BERT model URI and type of Transformer to train it.
		model_provider_parameters=HuggingFaceModelParams(
			model_uri="hf://google-bert/bert-base-cased",
			transformer_type=transformers.AutoModelForSequenceClassification,
		),
		# Use 3000 samples from Yelp dataset.
		dataset_provider_parameters=HuggingFaceDatasetParams(
			repo_id="yelp_review_full",
			split="train[:3000]",
		),
		# Specify HuggingFace Trainer parameters.
		trainer_parameters=HuggingFaceTrainerParams(
			training_parameters=transformers.TrainingArguments(
				output_dir="test_trainer",
				save_strategy="no",
				eval_strategy="epoch",
				disable_tqdm=True,
				log_level="info",
			),
			# Set LoRA config to reduce number of trainable model parameters.
			lora_config=LoraConfig(
				r=8,
				lora_alpha=8,
				lora_dropout=0.1,
				bias="none",
			),
    	),	
	)
	return result.metrics

# Define the hyperparameters search space
hp_space = {
	"learning_rate": katib.search.double(min=1e-05, max=5e-05),
	"per_device_train_batch_size": katib.search.int(min=8, max=64),
	"num_train_epochs": katib.search.int(min=1, max=10),
	"weight_decay": katib.search.double(min=0.0, max=1.0),
	"lora_config.r": katib.search.int(min=8, max=32)
}

# Run the tuning Experiment
exp_name = "llm-hp-optimization-test"
katib_client.tune(
	name=exp_name,
	hp_space=hp_space,
	objective=objective(exp_name), 
	objective_metric_name="eval_accuracy", 
	objective_type="maximize", 
	algorithm_name="random",
	max_trial_count=50,
	parallel_trial_count=2,
	resources_per_trial={
		"gpu": 8,
		"cpu": 20,
		"memory": "40G",
	},
	num_workers=4,
	num_procs_per_worker=2,
	resources_per_worker={
		"gpu": 2,
		"cpu": 5,
		"memory": "10G",
	}
)

# Get the best hyperparameters
best_hps = katib_client.get_optimal_hyperparameters(exp_name)
```

## Implementation

We will utilize the [train API](https://github.com/kubeflow/training-operator/blob/6ce4d57d699a76c3d043917bd0902c931f14080f/sdk/python/kubeflow/training/api/training_client.py#L96) from the Training Operator to create the objective function for hyperparameter optimization. This involves importing pretrained models and datasets from external platforms like HuggingFace and S3, and enabling seamless model training. The `tune` function will then automate the exploration of hyperparameter configurations to optimize LLM performance.

**Model and Dataset Downloading**: The [storage_initializer](https://github.com/kubeflow/training-operator/tree/6ce4d57d699a76c3d043917bd0902c931f14080f/sdk/python/kubeflow/storage_initializer) in the Training Operator will manage downloading and storing pretrained models and datasets via a PersistentVolumeClaim (PVC). This volume will be shared across containers, ensuring efficient access to the pretrained model and dataset without redundant downloads.

**Experiment and Trial Configuration**: Similar to the [existing tune API](https://github.com/kubeflow/katib/blob/0d190b94373c2f8f6150bf17d6dfa3698f4b2961/sdk/python/v1beta1/kubeflow/katib/api/katib_client.py#L152), we will create an experiment that defines the objective metric, search space, optimization algorithm, etc. The Experiment orchestrates the hyperparameter optimization process, including generating Trials, tracking their results, and identifying the optimal hyperparameters. Each Trial, implemented as a Kubernetes Job, runs model training with specific hyperparameters. In the trial specification for PyTorchJob, we will define the init container, master, and worker containers, allowing effective distribution across workers. Trial results will be fed back into the Experiment for evaluation, which then generates new Trials to further explore the hyperparameter space or concludes the tuning process.

## Advanced Functionalities

1. Incorporate early stopping strategy into the API to optimize training efficiency.
2. Expand support for distributed training in frameworks beyond PyTorch by leveraging their distributed training capabilities.
3. Support adding custom providers through configmap or CRD approach to enhance flexibility.
4. Enable users to deploy tuned models for inference within their applications or seamlessly integrate them into existing NLP pipelines for specialized tasks.

_#WIP_