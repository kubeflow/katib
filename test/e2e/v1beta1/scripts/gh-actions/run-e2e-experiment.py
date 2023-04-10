import argparse
import yaml
import time
import logging
from kubernetes import client
from kubeflow.katib import ApiClient, KatibClient, models
from kubeflow.katib.utils.utils import FakeResponse
from kubeflow.katib.constants import constants

# Experiment timeout is 40 min.
EXPERIMENT_TIMEOUT = 60 * 40

# The default logging config.
logging.basicConfig(level=logging.INFO)


def verify_experiment_results(
    katib_client: KatibClient,
    experiment: models.V1beta1Experiment,
    exp_name: str,
    exp_namespace: str,
):

    # Get the best objective metric.
    best_objective_metric = None
    for metric in experiment.status.current_optimal_trial.observation.metrics:
        if metric.name == experiment.spec.objective.objective_metric_name:
            best_objective_metric = metric
            break

    if best_objective_metric is None:
        raise Exception(
            "Unable to get the best metrics for objective: {}. Current Optimal Trial: {}".format(
                experiment.spec.objective.objective_metric_name,
                experiment.status.current_optimal_trial,
            )
        )

    # Get Experiment Succeeded reason.
    for c in experiment.status.conditions:
        if (
            c.type == constants.EXPERIMENT_CONDITION_SUCCEEDED
            and c.status == constants.CONDITION_STATUS_TRUE
        ):
            succeeded_reason = c.reason
            break

    trials_completed = experiment.status.trials_succeeded or 0
    trials_completed += experiment.status.trials_early_stopped or 0
    max_trial_count = experiment.spec.max_trial_count

    # If Experiment is Succeeded because of Max Trial Reached, all Trials must be completed.
    if (
        succeeded_reason == "ExperimentMaxTrialsReached"
        and trials_completed != max_trial_count
    ):
        raise Exception(
            "All Trials must be Completed. Max Trial count: {}, Experiment status: {}".format(
                max_trial_count, experiment.status
            )
        )

    # If Experiment is Succeeded because of Goal reached, the metrics must be correct.
    if succeeded_reason == "ExperimentGoalReached" and (
        (
            experiment.spec.objective.type == "minimize"
            and float(best_objective_metric.min) > float(experiment.spec.objective.goal)
        )
        or (
            experiment.spec.objective.type == "maximize"
            and float(best_objective_metric.max) < float(experiment.spec.objective.goal)
        )
    ):
        raise Exception(
            "Experiment goal is reached, but metrics are incorrect. "
            f"Experiment objective: {experiment.spec.objective}. "
            f"Experiment best objective metric: {best_objective_metric}"
        )

    # Verify Suggestion's resources. Suggestion name = Experiment name.
    suggestion = katib_client.get_suggestion(exp_name, exp_namespace)

    # For the Never or FromVolume resume policies Suggestion must be Succeeded.
    # For the LongRunning resume policy Suggestion must be always Running.
    for c in suggestion.status.conditions:
        if (
            c.type == constants.EXPERIMENT_CONDITION_SUCCEEDED
            and c.status == constants.CONDITION_STATUS_TRUE
            and experiment.spec.resume_policy == "LongRunning"
        ):
            raise Exception(
                f"Suggestion is Succeeded while Resume Policy is {experiment.spec.resume_policy}."
                f"Suggestion conditions: {suggestion.status.conditions}"
            )
        elif (
            c.type == constants.EXPERIMENT_CONDITION_RUNNING
            and c.status == constants.CONDITION_STATUS_TRUE
            and experiment.spec.resume_policy != "LongRunning"
        ):
            raise Exception(
                f"Suggestion is Running while Resume Policy is {experiment.spec.resume_policy}."
                f"Suggestion conditions: {suggestion.status.conditions}"
            )

    # For Never and FromVolume resume policies verify Suggestion's resources.
    if (
        experiment.spec.resume_policy == "Never"
        or experiment.spec.resume_policy == "FromVolume"
    ):
        resource_name = exp_name + "-" + experiment.spec.algorithm.algorithm_name

        # Suggestion's Service and Deployment should be deleted.
        for i in range(10):
            try:
                client.AppsV1Api().read_namespaced_deployment(
                    resource_name, exp_namespace
                )
            except client.ApiException as e:
                if e.status == 404:
                    break
                else:
                    raise e
            # Deployment deletion might take some time.
            time.sleep(1)
        if i == 10:
            raise Exception(
                "Suggestion Deployment is still alive for Resume Policy: {}".format(
                    experiment.spec.resume_policy
                )
            )

        try:
            client.CoreV1Api().read_namespaced_service(resource_name, exp_namespace)
        except client.ApiException as e:
            if e.status != 404:
                raise e
        else:
            raise Exception(
                "Suggestion Service is still alive for Resume Policy: {}".format(
                    experiment.spec.resume_policy
                )
            )

        # For FromVolume resume policy PVC should not be deleted.
        if experiment.spec.resume_policy == "FromVolume":
            try:
                client.CoreV1Api().read_namespaced_persistent_volume_claim(
                    resource_name, exp_namespace
                )
            except client.ApiException:
                raise Exception("PVC is deleted for FromVolume Resume Policy")


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
