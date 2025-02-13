import argparse
import logging

from kubeflow.katib import KatibClient, search
from kubeflow.katib.types.types import TrainerResources
from kubernetes import client
from verify import verify_experiment_results

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
    logging.debug(katib_client.get_experiment(exp_name, exp_namespace))
    logging.debug(katib_client.get_suggestion(exp_name, exp_namespace))


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
    logging.debug(katib_client.get_experiment(exp_name, exp_namespace))
    logging.debug(katib_client.get_suggestion(exp_name, exp_namespace))


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
