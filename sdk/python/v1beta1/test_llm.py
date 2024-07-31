import kubeflow.katib as katib
from kubeflow.katib import KatibClient

import transformers
from peft import LoraConfig

from kubeflow.storage_initializer.hugging_face import (
	HuggingFaceModelParams,
	HuggingFaceDatasetParams,
	HuggingFaceTrainerParams,
)

cl = KatibClient(namespace="kubeflow")


# [3] Create Katib Experiment with 12 Trials and 2 CPUs per Trial.
name = "llm-experiment"
cl.tune(
	name = name,
	# BERT model URI and type of Transformer to train it.
	model_provider_parameters = HuggingFaceModelParams(
		model_uri = "hf://google-bert/bert-base-cased",
		transformer_type = transformers.AutoModelForSequenceClassification,
	),
	# Use 3000 samples from Yelp dataset.
	dataset_provider_parameters = HuggingFaceDatasetParams(
		repo_id = "yelp_review_full",
		split = "train[:8]",
	),
	# Specify HuggingFace Trainer parameters.
	trainer_parameters = HuggingFaceTrainerParams(
		training_parameters = transformers.TrainingArguments(
			output_dir = "test_tune_api",
			save_strategy = "no",
			learning_rate = katib.search.double(min=1e-05, max=5e-05),
            #no_cuda=True, #if you use cpu instead of gpu
            #use_cpu=True, #if you use cpu instead of gpu
            num_train_epochs=1,
		),
		# Set LoRA config to reduce number of trainable model parameters.
		lora_config = LoraConfig(
			r = katib.search.int(min=8, max=32),
			lora_alpha = 8,
			lora_dropout = 0.1,
			bias = "none",
		),
	),	
	objective_metric_name = "train_loss", 
	objective_type = "minimize", 
	algorithm_name = "random",
	max_trial_count = 1,
	parallel_trial_count = 1,
    resources_per_trial={
        "cpu": "4",
        "memory": "10G",
    },
)

# [4] Wait until Katib Experiment is complete
cl.wait_for_experiment_condition(name=name)

# [5] Get the best hyperparameters.
#print(cl.get_optimal_hyperparameters(name))