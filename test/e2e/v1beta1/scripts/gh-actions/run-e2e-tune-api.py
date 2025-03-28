import argparse
import logging
from pprint import pformat

import kubeflow.katib as katib
from kubeflow.katib import KatibClient, search
from kubeflow.katib.types.types import TrainerResources
from kubernetes import client
from verify import verify_experiment_results
import os

# Experiment timeout is 40 min.
EXPERIMENT_TIMEOUT = 60 * 40

# The default logging config.
logging.basicConfig(level=logging.INFO)

def run_e2e_experiment_create_by_tune(
    katib_client: KatibClient,
    exp_name: str,
    exp_namespace: str,
):
    # Create Katib Experiment and wait until it is finished.
    logging.debug("Creating Experiment: {}/{}".format(exp_namespace, exp_name))

    # Use the test case from get-started tutorial.
    # https://www.kubeflow.org/docs/components/katib/getting-started/#getting-started-with-katib-python-sdk
    # [1] Create an objective function.
    def objective(parameters):
        import time

        time.sleep(5)
        result = 4 * int(parameters["a"]) - float(parameters["b"]) ** 2
        print(f"result={result}")

    # [2] Create hyperparameter search space.
    parameters = {"a": search.int(min=10, max=20), "b": search.double(min=0.1, max=0.2)}

    # [3] Create Katib Experiment with 4 Trials and 2 CPUs per Trial.
    # And Wait until Experiment reaches Succeeded condition.
    katib_client.tune(
        name=exp_name,
        namespace=exp_namespace,
        objective=objective,
        parameters=parameters,
        objective_metric_name="result",
        max_trial_count=4,
        resources_per_trial={"cpu": "100m"},
    )
    experiment = katib_client.wait_for_experiment_condition(
        exp_name, exp_namespace, timeout=EXPERIMENT_TIMEOUT
    )

    # Verify the Experiment results.
    verify_experiment_results(katib_client, experiment, exp_name, exp_namespace)

    # Print the Experiment and Suggestion.
    logging.debug("Experiment:\n%s", pformat(katib_client.get_experiment(exp_name, exp_namespace)))
    logging.debug("Suggestion:\n%s", pformat(katib_client.get_suggestion(exp_name, exp_namespace)))

def run_e2e_experiment_create_by_tune_pytorchjob(
    katib_client: KatibClient,
    exp_name: str,
    exp_namespace: str,
):
    # Create Katib Experiment and wait until it is finished.
    logging.debug("Creating Experiment: {}/{}".format(exp_namespace, exp_name))

    # Verify the PyTorchJob distributed.
    def objective(parameters):
        import os
        import time

        import torch.distributed as dist

        # Setup PyTorch distributed.
        dist.init_process_group(backend="gloo")

        print(
            "PyTorch Dist. WORLD_SIZE: {}, RANK: {}, LOCAL_RANK: {}".format(
                dist.get_world_size(), dist.get_rank(), os.getenv("LOCAL_RANK")
            )
        )

        time.sleep(5)
        # Only get results from the process with RANK=0.
        if dist.get_rank() == 0:
            result = 4 * int(parameters["a"]) - float(parameters["b"]) ** 2
            print(f"result={result}")
        dist.destroy_process_group()

    # Create Katib Experiment with 3 Trials. Every Trial runs PyTorchJob with 2 workers.
    katib_client.tune(
        name=exp_name,
        namespace=exp_namespace,
        objective=objective,
        parameters={
            "a": search.int(min=10, max=20),
            "b": search.double(min=0.1, max=0.2),
        },
        objective_metric_name="result",
        max_trial_count=3,
        parallel_trial_count=2,
        resources_per_trial=TrainerResources(
            num_workers=2,
            num_procs_per_worker=2,
            resources_per_worker={"cpu": "100m"},
        ),
    )

    experiment = katib_client.wait_for_experiment_condition(
        exp_name, exp_namespace, timeout=EXPERIMENT_TIMEOUT
    )

    # Verify the Experiment results.
    verify_experiment_results(katib_client, experiment, exp_name, exp_namespace)

    # Print the Experiment and Suggestion.
    logging.debug("Experiment:\n%s", pformat(katib_client.get_experiment(exp_name, exp_namespace)))
    logging.debug("Suggestion:\n%s", pformat(katib_client.get_suggestion(exp_name, exp_namespace)))

def run_e2e_experiment_create_by_tune_with_llm_optimization(
    katib_client: KatibClient,
    exp_name: str,
    exp_namespace: str,
):
    import transformers
    from kubeflow.storage_initializer.hugging_face import (
        HuggingFaceDatasetParams,
        HuggingFaceModelParams,
        HuggingFaceTrainerParams,
    )
    from peft import LoraConfig

    # Create Katib Experiment and wait until it is finished.
    logging.debug("Creating Experiment: {}/{}".format(exp_namespace, exp_name))
    
    # Use the test case from fine-tuning API tutorial.
    # https://www.kubeflow.org/docs/components/training/user-guides/fine-tuning/
    # Create Katib Experiment.
    # And Wait until Experiment reaches Succeeded condition.
    katib_client.tune(
        name=exp_name,
        namespace=exp_namespace,
        # BERT model URI and type of Transformer to train it.
        model_provider_parameters=HuggingFaceModelParams(
            model_uri="hf://google-bert/bert-base-cased",
            transformer_type=transformers.AutoModelForSequenceClassification,
            num_labels=5,
        ),
        # In order to save test time, use 8 samples from Yelp dataset.
        dataset_provider_parameters=HuggingFaceDatasetParams(
            repo_id="yelp_review_full",
            split="train[:8]",
        ),
        # Specify HuggingFace Trainer parameters.
        trainer_parameters=HuggingFaceTrainerParams(
            training_parameters=transformers.TrainingArguments(
                output_dir="test_tune_api",
                save_strategy="no",
                learning_rate = search.double(min=1e-05, max=5e-05),
                num_train_epochs=1,
            ),
            # Set LoRA config to reduce number of trainable model parameters.
            lora_config=LoraConfig(
                r = search.int(min=8, max=32),
                lora_alpha=8,
                lora_dropout=0.1,
                bias="none",
            ),
        ),
        objective_metric_name = "train_loss", 
        objective_type = "minimize", 
        algorithm_name = "random",
        max_trial_count = 1,
        parallel_trial_count = 1,
        resources_per_trial=katib.TrainerResources(
            num_workers=1,
            num_procs_per_worker=1,
            resources_per_worker={"cpu": "2", "memory": "10G",},
        ),
        storage_config={
            "size": "10Gi",
            "access_modes": ["ReadWriteOnce"],
        },
        retain_trials=True,
    )
    experiment = katib_client.wait_for_experiment_condition(
        exp_name, exp_namespace, timeout=EXPERIMENT_TIMEOUT
    )

    # Verify the Experiment results.
    verify_experiment_results(katib_client, experiment, exp_name, exp_namespace)

    # Print the Experiment and Suggestion.
    logging.debug("Experiment:\n%s", pformat(katib_client.get_experiment(exp_name, exp_namespace)))
    logging.debug("Suggestion:\n%s", pformat(katib_client.get_suggestion(exp_name, exp_namespace)))

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--namespace",
        type=str,
        required=True,
        help="Namespace for the Katib E2E test",
    )
    parser.add_argument(
        "--verbose",
        action="store_true",
        help="Verbose output for the Katib E2E test",
    )
    args = parser.parse_args()

    if args.verbose:
        logging.getLogger().setLevel(logging.DEBUG)

    katib_client = KatibClient()

    namespace_labels = client.CoreV1Api().read_namespace(args.namespace).metadata.labels
    if "katib.kubeflow.org/metrics-collector-injection" not in namespace_labels:
        namespace_labels["katib.kubeflow.org/metrics-collector-injection"] = "enabled"
        client.CoreV1Api().patch_namespace(
            args.namespace, {"metadata": {"labels": namespace_labels}}
        )

    # Test with run_e2e_experiment_create_by_tune
    exp_name = "tune-example"
    exp_namespace = args.namespace
    try:
        run_e2e_experiment_create_by_tune(katib_client, exp_name, exp_namespace)
        logging.info("---------------------------------------------------------------")
        logging.info(
            f"E2E is succeeded for Experiment created by tune: {exp_namespace}/{exp_name}"
        )
    except Exception as e:
        logging.info("---------------------------------------------------------------")
        logging.info(
            f"E2E is failed for Experiment created by tune: {exp_namespace}/{exp_name}"
        )
        raise e
    finally:
        # Delete the Experiment.
        logging.info("---------------------------------------------------------------")
        logging.info("---------------------------------------------------------------")
        katib_client.delete_experiment(exp_name, exp_namespace)

    # Test with run_e2e_experiment_create_by_tune_pytorchjob
    exp_name = "tune-example-pytorchjob"
    exp_namespace = args.namespace
    try:
        run_e2e_experiment_create_by_tune_pytorchjob(
            katib_client, exp_name, exp_namespace
        )
        logging.info("---------------------------------------------------------------")
        logging.info(
            f"E2E is succeeded for Experiment created by tune with PyTorchJob: {exp_namespace}/{exp_name}"
        )
    except Exception as e:
        logging.info("---------------------------------------------------------------")
        logging.info(
            f"E2E is failed for Experiment created by tune with PyTorchJob: {exp_namespace}/{exp_name}"
        )
        raise e
    finally:
        # Delete the Experiment.
        logging.info("---------------------------------------------------------------")
        logging.info("---------------------------------------------------------------")
        katib_client.delete_experiment(exp_name, exp_namespace)

    exp_name = "tune-example-llm-optimization"
    exp_namespace = args.namespace
    try:
        run_e2e_experiment_create_by_tune_with_llm_optimization(katib_client, exp_name, exp_namespace)
        logging.info("---------------------------------------------------------------")
        logging.info(f"E2E is succeeded for Experiment created by tune: {exp_namespace}/{exp_name}")
    except Exception as e:
        logging.info("---------------------------------------------------------------")
        logging.info(f"E2E is failed for Experiment created by tune: {exp_namespace}/{exp_name}")
        raise e
    finally:
        # Delete the Experiment.
        logging.info("---------------------------------------------------------------")
        logging.info("---------------------------------------------------------------")
        katib_client.delete_experiment(exp_name, exp_namespace)
