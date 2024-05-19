# LLM Hyperparameters Tuning API Proposal

- [LLM Hyperparameters Tuning API Proposal](#llm-hyperparameters-tuning-api-proposal)
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

LLMs are experiencing rapid growth in popularity and evolution, leading to an increased demand for LLMOps in Kubernetes environment. Hyperparameter tuning can greatly improve model performance but is often labor-intensive for humans. We want to provide an LLM hyperparameters tuning API which automates the process and is user-friendly for data scientists.

## Goal

We aim to develop a high-level API for tuning LLM hyperparameters, which abstracts away infrastructure complexities for data scientists, and facilitates easy hyperparameter tuning by enabling the seamless import of pretrained models and datasets from platforms like HuggingFace and S3, thereby improving accessibility and ease of use.

## Design for API

![Design for API](https://github.com/helenxie-bit/kubeflow/blob/1cdff9f77d68df824bea128e419dc612e172173b/design_api.jpg)

```python
# Reuse existing assets in Training Operator for importing models and datasets
import kubeflow.storage_initializer.hugging_face
import kubeflow.storage_initializer.s3

from kubeflow import katib
katib_client = katib.KatibClient()

# Arguments related to the model provider including credentials
model_args = modelProviderClass()

# Arguments related to the dataset provider including credentials 
dataset_args = datasetProviderClass()

# Define hyperparameters and search space
parameters = {
	"learning_rate": katib.search.double(min=[], max=[]),
	"batch_size": katib.search.int(min=[], max=[]),
	"num_epoch": katib.search.int(min=[], max=[]),
	"weight_decay": katib.search.double(min=[], max=[])
	...
}

# Create experiment for tuning
exp_name = "llm-tune-test"
katib_client.tune(
	name=exp_name,
	model=model_args,
	dataset=dataset_args,
	parameters=parameters, 
	objective_metric_name: str = None,
	additional_metric_names: List[str] = [],
	objective_type: str = "maximize", 
	objective_goal: float = None,
	algorithm_name: str = "random",
	resources_per_trial: Union[dict, client.V1ResourceRequirements, None] = None, 
	parallel_trial_count: int = None,
	max_trial_count: int = None,
	max_failed_trail_count: int = None
)

# Monitor experiment progress
katib_client.get_experiment_conditions(exp_name)
katib_client.get_success_trail_details(exp_name)

# Get the best hyperparameters
best_hps = katib_client.get_optimal_hyperparameters(exp_name)
```

### Example

```python
import kubeflow.storage_initializer.hugging_face
from kubeflow import katib

# Import model and dataset
model_args = HuggingFaceModelParams(
        model_uri="hf://google-bert/bert-base-cased",
        transformer_type=transformers.AutoModelForSequenceClassification,
)

dataset_args = HfDatasetParams(
        repo_id="yelp_review_full",
        split="train[:3000]",
)

# Define hyperparameters and search space
parameters = {
	"learning_rate": katib.search.double(min=1e-05, max=5e-05),
	"batch_size": katib.search.int(min=8, max=64),
	"num_epoch": katib.search.int(min=1, max=10),
	"weight_decay": katib.search.double(min=0.0, max=1.0)
	...
}

# Create experiment for tuning
exp_name = "llm-tune-test"
katib_client.tune(
	name=exp_name,
	model=model_args, 
	dataset=dataset_args, 
	parameters=parameters, 
	objective_metric_name="accuracy", 
	objective_type="maximize", 
	algorithm_name="random",
	resources_per_trial={"gpu": "2"},
	parallel_trial_count=2,
	max_trial_count=50
)

# Get the best hyperparameters
best_hps = katib_client.get_optimal_hyperparameters(exp_name)
```

## Implementation

We will leverage [existing assets](https://github.com/kubeflow/training-operator/tree/687f0c9d2f5cf5dcc97dec87c869ec7f1309d07c/sdk/python/kubeflow/storage_initializer) from the Training Operator to enable the import of pre-trained models and datasets from external platforms (e.g., HuggingFace, S3, etc.). This involves utilizing a storage initializer (PVC) to download pre-trained models and datasets, then sharing this volume with other containers responsible for computation. Additionally, we will define the init container as well as containers for the master and workers in the trail specification for PyTorchJob, similar to the implementation of the train API in Training Operator. This allows us to distribute trials to different workers effectively.

Regarding hyperparameter tuning, the API will feature a tune function similar to the [previous one](https://github.com/kubeflow/katib/tree/master/sdk/python/v1beta1). This function enables users to specify parameters such as hyperparameters, objective metric, algorithm name, etc., facilitating the exploration of various hyperparameter configurations to optimize LLM performance. However, the objective function in the previous tune API will be replaced by the model provider and dataset provider, aligning with the import process for models and datasets.

## Advanced Functionalities

1. Incorporate early stopping strategy into the API to optimize training efficiency.
2. Expand support for distributed training in frameworks beyond PyTorch by leveraging their distributed training capabilities.
3. Support adding custom providers through configmap or CRD approach to enhance flexibility.
4. Enable users to deploy tuned models for inference within their applications or seamlessly integrate them into existing NLP pipelines for specialized tasks.

_#WIP_
