import time

from kubeflow.katib import KatibClient, models
from kubeflow.katib.constants import constants
from kubernetes import client


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
