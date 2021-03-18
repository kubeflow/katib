# KatibClient

> KatibClient(config_file=None, context=None, client_configuration=None, persist_config=True)

User can load authentication and cluster information from kube-config file and
stores them in kubernetes.client.configuration. Parameters are as following:

| Parameter            | Description                                                               |
| -------------------- | ------------------------------------------------------------------------- |
| config_file          | Name of the kube-config file. Defaults to ~/.kube/config.                 |
| context              | Set the active context. Defaults to current_context from the kube-config. |
| client_configuration | The kubernetes.client.Configuration to set configs to.                    |
| persist_config       | If True, config file will be updated when changed.                        |

The APIs for KatibClient are as following:

| Class       | Method                                                      | Description                                                 |
| ----------- | ----------------------------------------------------------- | ----------------------------------------------------------- |
| KatibClient | [create_experiment](#create_experiment)                     | Create the Katib Experiment                                 |
| KatibClient | [get_experiment](#get_experiment)                           | Get the Katib Experiment                                    |
| KatibClient | [get_suggestion](#get_suggestion)                           | Get the Katib Suggestion                                    |
| KatibClient | [delete_experiment](#delete_experiment)                     | Delete the Katib Experiment                                 |
| KatibClient | [list_experiments](#list_experiments)                       | List all Katib Experiments                                  |
| KatibClient | [get_experiment_status](#get_experiment_status)             | Get the Experiment current status                           |
| KatibClient | [is_experiment_succeeded](#is_experiment_succeeded)         | Check if Experiment has succeeded                           |
| KatibClient | [list_trials](#list_trials)                                 | List all Experiment's Trials                                |
| KatibClient | [get_success_trial_details](#get_success_trial_details)     | Get the Trial details that have succeeded for an Experiment |
| KatibClient | [get_optimal_hyperparameters](#get_optimal_hyperparameters) | Get the current optimal Trial from the Experiment           |

## create_experiment

> create_experiment(exp_object, namespace=None)

Create the Katib Experiment.
If the namespace is `None`, it takes namespace from the Experiment or "default".

Return the created Experiment.

### Parameters

| Name       | Type                                      | Description           | Notes    |
| ---------- | ----------------------------------------- | --------------------- | -------- |
| exp_object | [V1beta1Experiment](V1beta1Experiment.md) | Experiment object.    | Required |
| namespace  | str                                       | Experiment namespace. | Optional |

### Return type

dict

## get_experiment

> get_experiment(name=None, namespace=None)

Get the Katib Experiment.
If the name is `None` returns all Experiments in the namespace.
If the namespace is `None`, it takes namespace from the Experiment object or "default".

Return the Experiment object.

### Parameters

| Name      | Type | Description           | Notes    |
| --------- | ---- | --------------------- | -------- |
| name      | str  | Experiment name.      | Optional |
| namespace | str  | Experiment namespace. | Optional |

### Return type

dict

## get_suggestion

> get_suggestion(name=None, namespace=None)

Get the Katib Suggestion.
If the name is `None` returns all Suggestion in the namespace.
If the namespace is `None`, it takes namespace from the Suggestion object or "default".

Return the Suggestion object.

### Parameters

| Name      | Type | Description           | Notes    |
| --------- | ---- | --------------------- | -------- |
| name      | str  | Suggestion name.      | Optional |
| namespace | str  | Suggestion namespace. | Optional |

### Return type

dict

## delete_experiment

> delete_experiment(name, namespace=None)

Delete the Katib Experiment
If the namespace is `None`, it takes namespace from the Experiment object or "default".

Return the deleted Experiment object.

### Parameters

| Name      | Type | Description           | Notes    |
| --------- | ---- | --------------------- | -------- |
| name      | str  | Experiment name.      | Required |
| namespace | str  | Experiment namespace. | Optional |

### Return type

dict

## list_experiments

> list_experiments(namespace=None)

List all Katib Experiments.
If the namespace is `None`, it takes "default" namespace.

Return list of Experiment names with the statuses.

### Parameters

| Name      | Type | Description           | Notes    |
| --------- | ---- | --------------------- | -------- |
| namespace | str  | Experiment namespace. | Optional |

### Return type

list[dict]

## get_experiment_status

> get_experiment_status(name, namespace=None)

Get the Experiment current status.
If the namespace is `None`, it takes "default" namespace.

Return the current Experiment status.

### Parameters

| Name      | Type | Description           | Notes    |
| --------- | ---- | --------------------- | -------- |
| name      | str  | Experiment name.      | Required |
| namespace | str  | Experiment namespace. | Optional |

### Return type

str

## is_experiment_succeeded

> is_experiment_succeeded(name, namespace=None)

Check if Experiment has succeeded.
If the namespace is `None`, it takes "default" namespace.

Return whether Experiment has succeeded or not.

### Parameters

| Name      | Type | Description           | Notes    |
| --------- | ---- | --------------------- | -------- |
| name      | str  | Experiment name.      | Required |
| namespace | str  | Experiment namespace. | Optional |

### Return type

bool

## list_trials

> list_trials(name, namespace=None)

List all Experiment's Trials.
If the namespace is `None`, it takes "default" namespace.

Return list of Trial names with the statuses.

### Parameters

| Name      | Type | Description           | Notes    |
| --------- | ---- | --------------------- | -------- |
| name      | str  | Experiment name.      | Required |
| namespace | str  | Experiment namespace. | Optional |

### Return type

list[dict]

## get_success_trial_details

> get_success_trial_details(name, namespace=None)

Get the Trial details that have succeeded for an Experiment.
If the namespace is `None`, it takes namespace from the Experiment or "default".

Return Trial names with the hyperparameters and metrics.

### Parameters

| Name      | Type | Description           | Notes    |
| --------- | ---- | --------------------- | -------- |
| name      | str  | Experiment name.      | Required |
| namespace | str  | Experiment namespace. | Optional |

### Return type

list[dict]

## get_optimal_hyperparameters

> get_optimal_hyperparameters(name, namespace=None)

Get the current optimal Trial from the Experiment.
If the namespace is `None`, it takes namespace from the Experiment or "default".

Return current optimal Trial for the Experiment.

### Parameters

| Name      | Type | Description           | Notes    |
| --------- | ---- | --------------------- | -------- |
| name      | str  | Experiment name.      | Required |
| namespace | str  | Experiment namespace. | Optional |

### Return type

dict
