import argparse
import logging

from kubeflow.katib import KatibClient
from kubeflow.katib import search
from kubernetes import client
from kubernetes.client import V1Container
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
    parameters = {
        "a": search.int(min=10, max=20),
        "b": search.double(min=0.1, max=0.2)
    }

    # [3] Create a dummy metric collector (DOES NOT HAVE A IMAGE)
    metric_collector = V1Container(
        name="dummy-collector",
        image="dummy-collector:latest",
        command=["python", "/app/dummy-collector.py"],
        args=["--metric-name=result"],
        env=[
            client.V1EnvVar(name="EXPERIMENT_NAME", value=exp_name),
            client.V1EnvVar(name="EXPERIMENT_NAMESPACE", value=exp_namespace)
        ]
    )

    # [4] Create Katib Experiment with 4 Trials and 2 CPUs per Trial.
    # And Wait until Experiment reaches Succeeded condition.
    katib_client.tune(
        name=exp_name,
        namespace=exp_namespace,
        objective=objective,
        parameters=parameters,
        objective_metric_name="result",
        max_trial_count=4,
        resources_per_trial={"cpu": "2"},
        metrics_collector_config={
            "kind": "Custom",
            "customCollector": metric_collector,
        },
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
        "--namespace", type=str, required=True, help="Namespace for the Katib E2E test",
    )
    parser.add_argument(
        "--verbose", action="store_true", help="Verbose output for the Katib E2E test",
    )
    args = parser.parse_args()

    if args.verbose:
        logging.getLogger().setLevel(logging.DEBUG)

    katib_client = KatibClient()

    namespace_labels = client.CoreV1Api().read_namespace(args.namespace).metadata.labels
    if 'katib.kubeflow.org/metrics-collector-injection' not in namespace_labels:
        namespace_labels['katib.kubeflow.org/metrics-collector-injection'] = 'enabled'
        client.CoreV1Api().patch_namespace(args.namespace, {'metadata': {'labels': namespace_labels}})

    # Test with run_e2e_experiment_create_by_tune
    exp_name = "tune-example"
    exp_namespace = args.namespace
    try:
        run_e2e_experiment_create_by_tune(katib_client, exp_name, exp_namespace)
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