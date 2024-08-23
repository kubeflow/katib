import argparse
import logging

import yaml
from kubeflow.katib import ApiClient, KatibClient, models
from kubeflow.katib.constants import constants
from kubeflow.katib.utils.utils import FakeResponse
from kubernetes import client
from verify import verify_experiment_results

# Experiment timeout is 40 min.
EXPERIMENT_TIMEOUT = 60 * 40

# The default logging config.
logging.basicConfig(level=logging.INFO)


def run_e2e_experiment(
    katib_client: KatibClient,
    experiment: models.V1beta1Experiment,
    exp_name: str,
    exp_namespace: str,
):

    # Create Katib Experiment and wait until it is finished.
    logging.debug(
        "Creating Experiment: {}/{} with MaxTrialCount: {}, ParallelTrialCount: {}".format(
            exp_namespace,
            exp_name,
            experiment.spec.max_trial_count,
            experiment.spec.parallel_trial_count,
        )
    )

    # Wait until Experiment reaches Succeeded condition.
    katib_client.create_experiment(experiment, exp_namespace)
    experiment = katib_client.wait_for_experiment_condition(
        exp_name, exp_namespace, timeout=EXPERIMENT_TIMEOUT
    )

    # Test resume feature for "FromVolume" and "LongRunning" Experiments.
    if exp_name == "from-volume-resume" or exp_name == "long-running-resume":
        max_trial_count = experiment.spec.max_trial_count + 1
        parallel_trial_count = experiment.spec.parallel_trial_count + 1
        logging.debug(
            f"Restarting Experiment {exp_namespace}/{exp_name} "
            f"with MaxTrialCount: {max_trial_count} and ParallelTrialCount: {parallel_trial_count}"
        )

        # Modify Experiment budget.
        katib_client.edit_experiment_budget(
            exp_name, exp_namespace, max_trial_count, parallel_trial_count
        )
        # Wait until Experiment is Restarted.
        katib_client.wait_for_experiment_condition(
            exp_name,
            exp_namespace,
            constants.EXPERIMENT_CONDITION_RESTARTING,
            EXPERIMENT_TIMEOUT,
        )
        # Wait until Experiment is Succeeded.
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
        "--experiment-path",
        type=str,
        required=True,
        help="Path to the Katib Experiment.",
    )
    parser.add_argument(
        "--namespace", type=str, required=True, help="Namespace for the Katib E2E test",
    )
    parser.add_argument(
        "--trial-pod-annotations", type=str, help="Annotation for the pod created by trial",
    )
    parser.add_argument(
        "--verbose", action="store_true", help="Verbose output for the Katib E2E test",
    )
    args = parser.parse_args()

    if args.verbose:
        logging.getLogger().setLevel(logging.DEBUG)

    logging.info("---------------------------------------------------------------")
    logging.info("---------------------------------------------------------------")
    logging.info(f"Start E2E test for the Katib Experiment: {args.experiment_path}")

    # Read Experiment YAML to Fake Response object.
    with open(args.experiment_path, "r") as file:
        experiment = FakeResponse(yaml.safe_load(file))

    # Replace batch size to number of epochs for faster execution.
    experiment.data = experiment.data.replace("--batch-size=64", "--num-epochs=2")

    # Convert to the Katib Experiment object.
    experiment = ApiClient().deserialize(experiment, "V1beta1Experiment")
    experiment.metadata.namespace = args.namespace
    exp_name = experiment.metadata.name
    exp_namespace = experiment.metadata.namespace

    # Set Trial threshold for Katib Experiments.
    MAX_TRIAL_COUNT = 2
    PARALLEL_TRIAL_COUNT = 1
    MAX_FAILED_TRIAL_COUNT = 0

    # For one random search Experiment we test parallel execution.
    if experiment.metadata.name == "random":
        MAX_TRIAL_COUNT += 1
        PARALLEL_TRIAL_COUNT += 1
        if args.trial_pod_annotations:
            kind = experiment.spec.trial_template.trial_spec['kind']
            if kind != "Job":
                raise NotImplementedError(f'Trail pod annotations not implemented for {kind}!')

            trial_spec_metadata = experiment.spec.trial_template.trial_spec['spec']['template'].get('metadata', {})
            trial_spec_pod_annotations = trial_spec_metadata.get('annotations', {})
            trial_spec_pod_annotations.update(eval(args.trial_pod_annotations))
            trial_spec_metadata['annotations'] = trial_spec_pod_annotations
            experiment.spec.trial_template.trial_spec['spec']['template']['metadata'] = trial_spec_metadata

    # Hyperband will validate the parallel trial count, thus we should not change it.
    # We don't need to test parallel Trials for Darts.
    if (
        experiment.spec.algorithm.algorithm_name != "hyperband"
        and experiment.spec.algorithm.algorithm_name != "darts"
    ):
        experiment.spec.max_trial_count = MAX_TRIAL_COUNT
        experiment.spec.parallel_trial_count = PARALLEL_TRIAL_COUNT
        experiment.spec.max_failed_trial_count = MAX_FAILED_TRIAL_COUNT

    katib_client = KatibClient()

    namespace_labels = client.CoreV1Api().read_namespace(args.namespace).metadata.labels
    if 'katib.kubeflow.org/metrics-collector-injection' not in namespace_labels:
        namespace_labels['katib.kubeflow.org/metrics-collector-injection'] = 'enabled'
        client.CoreV1Api().patch_namespace(args.namespace, {'metadata': {'labels': namespace_labels}})

    try:
        run_e2e_experiment(katib_client, experiment, exp_name, exp_namespace)
        logging.info("---------------------------------------------------------------")
        logging.info(f"E2E is succeeded for Experiment: {exp_namespace}/{exp_name}")
    except Exception as e:
        logging.info("---------------------------------------------------------------")
        logging.info(f"E2E is failed for Experiment: {exp_namespace}/{exp_name}")
        raise e
    finally:
        # Delete the Experiment.
        logging.info("---------------------------------------------------------------")
        logging.info("---------------------------------------------------------------")
        katib_client.delete_experiment(exp_name, exp_namespace)
