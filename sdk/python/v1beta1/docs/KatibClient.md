# KatibClient

> KatibClient(config_file=None, context=None, client_configuration=None, persist_config=True)

User can load authentication and cluster information from kube-config file and stores them in kubernetes.client.configuration. Parameters are as following:

| Parameter | Description |
| --- | --- |
| config_file | Location of kube-config file. Defaults to `~/.kube/config`. Note that the config_file is needed if user want to operate katib SDK in another remote cluster, user must set `config_file` to load kube-config file explicitly, e.g. `KatibClient(config_file="~/.kube/config")`. |
| context | Set the active context. If is set to None, current_context from config file will be used. |
| client_configuration | The kubernetes.client.Configuration to set configs to. |
| persist_config | If True, config file will be updated when changed (e.g GCP token refresh). |

The APIs for KatibClient are as following:

| Class       | Method                                                      | Description                                                               |
| ----------- | ----------------------------------------------------------- | ------------------------------------------------------------------------- |
| KatibClient | [create_experiment](#create_experiment)                     | Create Katib Experiment                                                   |
| KatibClient | [get_experiment](#get_experiment)                           | Get or watch the specified Experiment or all Experiments in the namespace |
| KatibClient | [delete_experiment](#delete_experiment)                     | Delete specified Experiment                                               |
| KatibClient | [list_experiments](#list_experiments)                       | List all Experiments with status                                          |
| KatibClient | [get_experiment_status](#get_experiment_status)             | Get Experiment status                                                     |
| KatibClient | [is_experiment_succeeded](#is_experiment_succeeded)         | Check if Experiment status is Succeeded                                   |
| KatibClient | [list_trials](#list_trials)                                 | List all Trials of specified Experiment with status                       |
| KatibClient | [get_optimal_hyperparameters](#get_optimal_hyperparameters) | Get currentOptimalTrial with parameterAssignments of an Experiment        |

## create_experiment

> create_experiment(experiment, namespace=None)

### Parameters

| Name | Type | Description | Notes |
| --- | --- | --- | --- |
| experiment | [V1beta1Experiment](V1beta1Experiment.md) | experiment definition | Required |
| namespace | str | Namespace for experiment deploying to. If the `namespace` is not defined, will align with experiment definition, or use current or default namespace if namespace is not specified in experiment definition. | Optional |


### Return type

object

## get_experiment

> get_experiment(name=None, namespace=None)

Get Experiment in the specified namespace

### Parameters

| Name | Type | Description | Notes |
| --- | --- | --- | --- |
| name | str | Experiment name. If the `name` is not specified, will get all experiments in the namespace. | Optional |
| namespace | str | Experiment's namespace. Defaults to current or default namespace. | Optional |

### Return type

object

## delete_experiment

> delete_experiment(name, namespace=None)

### Parameters

| Name      | Type | Description                                                     | Notes    |
| --------- | ---- | --------------------------------------------------------------- | -------- |
| name      | str  | Experiment name.                                                | Required |
| namespace | str  | Experiment namespace. Defaults to current or default namespace. | Optional |

### Return type

object

## list_experiments

> list_experiments(namespace=None)

List all Experiments with status

### Parameters

| Name      | Type | Description                                                     | Notes    |
| --------- | ---- | --------------------------------------------------------------- | -------- |
| namespace | str  | Experiment namespace. Defaults to current or default namespace. | Optional |

### Return type

List

## get_experiment_status

> get_experiment_status(name, namespace=None)

Returns Experiment status, such as Created, Running, Failed, Succeeded or Restarting.

### Parameters

| Name      | Type | Description                                                     | Notes    |
| --------- | ---- | --------------------------------------------------------------- | -------- |
| name      | str  | Experiment name.                                                | Required |
| namespace | str  | Experiment namespace. Defaults to current or default namespace. | Optional |

### Return type

Str

## is_experiment_succeeded

> is_experiment_succeeded(name, namespace=None)

Returns True if Experiment succeeded; false otherwise.

### Parameters

| Name      | Type | Description                                                     | Notes    |
| --------- | ---- | --------------------------------------------------------------- | -------- |
| name      | str  | Experiment name.                                                | Required |
| namespace | str  | Experiment namespace. Defaults to current or default namespace. | Optional |

### Return type

Bool

## list_trials

> list_trials(name, namespace=None)

List all Trials of an Experiment with status

### Parameters

| Name      | Type | Description                                                     | Notes    |
| --------- | ---- | --------------------------------------------------------------- | -------- |
| name      | str  | Experiment name.                                                | Required |
| namespace | str  | Experiment namespace. Defaults to current or default namespace. | Optional |

### Return type

List

## get_optimal_hyperparameters

> get_optimal_hyperparameters(name, namespace=None)

Get currentOptimalTrial with parameterAssignments of an Experiment

### Parameters

| Name      | Type | Description                                                     | Notes    |
| --------- | ---- | --------------------------------------------------------------- | -------- |
| name      | str  | Experiment name.                                                | Required |
| namespace | str  | Experiment namespace. Defaults to current or default namespace. | Optional |

### Return type

Dict
