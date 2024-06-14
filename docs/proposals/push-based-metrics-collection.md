# Push-based Metrics Collection Proposal

## Links

- [katib/issues#577([Enhancement Request] Metrics Collector Push-based Implementation)](https://github.com/kubeflow/katib/issues/577)

## Motivation

[Katib](https://github.com/kubeflow/katib) is a Kubernetes-native project for automated machine learning (AutoML). It can not only tune hyperparameters of applications written in any language and natively supports many ML frameworks, but also supports features like early stopping and neural architecture search.

In the procedure of tuning hyperparameters, Metrics Collector, which is implemented as a sidecar container attached to each training container in the [current design](https://github.com/kubeflow/katib/blob/master/docs/proposals/metrics-collector.md), will collect training logs from Trials once the training is complete. Then, the Metrics Collector will parse training logs to get appropriate metrics like accuracy or loss and pass the evaluation results to the HyperParameter tuning algorithm.

However, current implementation of Metrics Collector is pull-based, raising some [design problems](https://github.com/kubeflow/training-operator/issues/722#issuecomment-405669269) such as determining the frequency we scrape the metrics, performance issues like the overhead caused by too many sidecar containers, and restrictions on developing environments which must support sidecar containers. Thus, we should implement a new API for Katib Python SDK to offer users a push-based way to store metrics directly into the Katib DB and resolve those issues raised by pull-based metrics collection.

![](../images/push-based-metrics-collection.png)

Fig.1 Architecture of the new design

### Goals
1. **A new parameter in Python SDK function `tune`**: allow users to specify the method of collecting metrics(push-based/pull-based).

2. **A new interface `report_metrics` in Python SDK**: push the metrics to Katib DB directly.

3. The final metrics of worker pods should be **pushed to Katib DB directly** in the push mode of metrics collection.

### Non-Goals
1. Implement authentication model for Katib DB to push metrics.

2. Support pushing data to different types of storage system(prometheus, self-defined interface etc.)


## API

### New Parameter in Python SDK Function `tune`

We decided to add `metrics_collection_mechanism` to `tune` function in Python SDK.

```Python
def tune(
    self,
    name: str,
    objective: Callable,
    parameters: Dict[str, Any],
    base_image: str = constants.BASE_IMAGE_TENSORFLOW,
    namespace: Optional[str] = None,
    env_per_trial: Optional[Union[Dict[str, str], List[Union[client.V1EnvVar, client.V1EnvFromSource]]]] = None,
    algorithm_name: str = "random",
    algorithm_settings: Union[dict, List[models.V1beta1AlgorithmSetting], None] = None,
    objective_metric_name: str = None,
    additional_metric_names: List[str] = [],
    objective_type: str = "maximize",
    objective_goal: float = None,
    max_trial_count: int = None,
    parallel_trial_count: int = None,
    max_failed_trial_count: int = None,
    resources_per_trial: Union[dict, client.V1ResourceRequirements, None] = None,
    retain_trials: bool = False,
    packages_to_install: List[str] = None,
    pip_index_url: str = "https://pypi.org/simple",
    # The newly added parameter metrics_collector_config.
    # It specifies the config of metrics collector, for example, 
    # metrics_collector_config={"kind": "Push"},
    metrics_collector_config: Dict[str, Any] = None, 
)
```

### New Interface `report_metrics` in Python SDK

```Python
"""Push Metrics Directly to Katib DB
    Katib DB Manager service should be accessible while calling this API.

    If you run this API in-cluster (e.g. from the Kubeflow Notebook) you can
    use the default Katib DB Manager address: `katib-db-manager.kubeflow:6789`.

    If you run this API outside the cluster, you have to port-forward the
    Katib DB Manager before getting the Trial metrics: `kubectl port-forward svc/katib-db-manager -n kubeflow 6789`.
    In that case, you can use this Katib DB Manager address: `localhost:6789`.

    You can use `curl` to verify that Katib DB Manager is reachable: `curl <db-manager-address>`.

    [!!!] Trial name should always be passed into Katib Trials as env variable `KATIB_TRIAL_NAME`.

    Args:
        metrics: Dict of metrics pushed to Katib DB.
            For examle, `metrics = {"loss": 0.01, "accuracy": 0.99}`.
        db-manager-address: Address for the Katib DB Manager in this format: `ip-address:port`.
    
    Raises:
        RuntimeError: Unable to push Trial metrics to Katib DB.
"""
def report_metrics(
    metrics: Dict[str, Any],
    db_manager_address: str = constants.DEFAULT_DB_MANAGER_ADDRESS,
)
```

### A Simple Example:

```Python
import kubeflow.katib as katib

# Step 1. Create an objective function with push-based metrics collection.
def objective(parameters):
    # Import required packages.
    import kubeflow.katib as katib
    # Calculate objective function.
    result = 4 * int(parameters["a"]) - float(parameters["b"]) ** 2
    # Push metrics to Katib DB.
    katib.report_metrics({"result": result})

# Step 2. Create HyperParameter search space.
parameters = {
    "a": katib.search.int(min=10, max=20),
    "b": katib.search.double(min=0.1, max=0.2)
}

# Step 3. Create Katib Experiment with 12 Trials and 2 GPUs per Trial.
katib_client = katib.KatibClient(namespace="kubeflow")
name = "tune-experiment"
katib_client.tune(
    name=name,
    objective=objective,
    parameters=parameters,
    objective_metric_name="result",
    max_trial_count=12,
    resources_per_trial={"gpu": "2"},
    metrics_collector_config={"kind": "Push"},
)

# Step 4. Get the best HyperParameters.
print(katib_client.get_optimal_hyperparameters(name))
```

## Implementation

### Add New Parameter in `tune`

As is mentioned above, we decided to add `metrics_collector_config` to the tune function in Python SDK. Also, we have some changes to be made:

1. Disable injection: set `katib.kubeflow.org/metrics-collector-injection` to `disabled` when the push-based way of metrics collection is adopted so as to disable the injection of the metrics collection sidecar container.

2. Configure the way of metrics collection: set the configuration `spec.metricsCollectionSpec.collector.kind`(specify the way of metrics collection) to `Push`.

3. Rename metrics collector from `None` to `Push`: It's not correct to call push-based metrics collection `None`. We should modify related code to rename it.

4. Write env variables into trial spec: set `KATIB_TRIAL_NAME` for `report_metrics` function to dial db manager.

### New Interface `report_metrics` in Python SDK

We decide to implement this funcion to push metrics directly to Katib DB with the help of grpc. Trial name should always be passed into Katib Trials (and then into this function) as env variable `KATIB_TRIAL_NAME`. 

Also, the function is supposed to be implemented as **global function** because it is called in the user container.

Steps:

1. Wrap metrics into `katib_api_pb2.ReportObservationLogRequest`:

Firstly, convert metrics (in dict format) into `katib_api_pb2.ReportObservationLogRequest` type for the following grpc call, referring to https://github.com/kubeflow/katib/blob/master/pkg/apis/manager/v1beta1/gen-doc/api.md#reportobservationlogrequest

2. Dial Katib DBManager Service

We'll create a DBManager Stub and make a grpc call to report metrics to Katib DB.

### Compatibility Changes in Trial Controller

We need to make appropriate changes in the Trial controller to make sure we insert unavailable value into Katib DB, if user doesn't report metric accidentally. The current implementation handles unavailable metrics in:

```Golang
// If observation is empty metrics collector doesn't finish.
// For early stopping metrics collector are reported logs before Trial status is changed to EarlyStopped.
if jobStatus.Condition == trialutil.JobSucceeded && instance.Status.Observation == nil {
	logger.Info("Trial job is succeeded but metrics are not reported, reconcile requeued")
	return errMetricsNotReported
}
```
1. Distinguish pull-based and push-based metrics collection

We decide to add a if-else statement in the code above to distinguish pull-based and push-based metrics collection. In the push-based collection, the trial does not need to be requeued. Instead, we'll insert a unavailable value to Katib DB.

2. Update the status of trial to `MetricsUnavailable`

In the current implementation of pull-based metrics collection, trials will be re-queued when the metrics collector finds the `.Status.Observation` is empty. However, it's not compatible with push-based metrics collection because the forgotten metrics won't be reported in the new round of reconcile. So, we need to update its status in the function `UpdateTrialStatusCondition` in accomodation with the pull-based metrics collection. The following code will be insert into lines before [trial_controller_util.go#L69](https://github.com/kubeflow/katib/blob/7959ffd54851216dbffba791e1da13c8485d1085/pkg/controller.v1beta1/trial/trial_controller_util.go#L69)


```Golang
else if instance.Spec.MetricCollector.Collector.Kind == "Push" && instance.Status.Obeservation == nil {
    ... // Update the status of this trial to `MetricsUnavailable` and output the reason.
}
```

### Collection of Final Metrics

The final metrics of worker pods should be pushed to Katib DB directly in the push mode of metrics collection.

\#WIP