# Kubeflow Katib SDK

Python SDK for Kubeflow Katib

## Requirements.

Python 2.7 and 3.4+

## Installation & Usage

### pip install

```sh
pip install kubeflow-katib
```

Then import package:

```python
from kubeflow import katib
```

### Setuptools

Install via [Setuptools](http://pypi.python.org/pypi/setuptools).

```sh
python setup.py install --user
```

(or `sudo python setup.py install` to install the package for all users)

### Publish new SDK version to PyPi

Our SDK is located in [`kubeflow-katib` PyPi package](https://pypi.org/project/kubeflow-katib/).
Katib Python SDK is published as part of the Katib patch releases.
You can check the release process [here](../../../scripts/v1beta1/release.sh).
For each Katib patch release, we upload a new SDK version to the PyPi.
The SDK version is equal to the Katib version.

## Getting Started

Please follow the [examples](../../../examples/v1beta1/sdk) to learn more about Katib SDK.

## Documentation for API Endpoints

| Class                 | Method                                 | Description                                                 |
| --------------------- | -------------------------------------- | ----------------------------------------------------------- |
| [KatibClient][client] | [create_experiment][create]            | Create the Katib Experiment                                 |
| [KatibClient][client] | [get_experiment][get_e]                | Get the Katib Experiment                                    |
| [KatibClient][client] | [get_suggestion][get_s]                | Get the Katib Suggestion                                    |
| [KatibClient][client] | [delete_experiment][delete]            | Delete the Katib Experiment                                 |
| [KatibClient][client] | [list_experiments][list_e]             | List all Katib Experiments                                  |
| [KatibClient][client] | [get_experiment_status][get_status]    | Get the Experiment current status                           |
| [KatibClient][client] | [is_experiment_succeeded][is_suc]      | Check if Experiment has succeeded                           |
| [KatibClient][client] | [list_trials][list_t]                  | List all Experiment's Trials                                |
| [KatibClient][client] | [get_success_trial_details][get_suc_t] | Get the Trial details that have succeeded for an Experiment |
| [KatibClient][client] | [get_optimal_hyperparameters][opt_hp]  | Get the current optimal Trial from the Experiment           |

[client]: docs/KatibClient.md
[create]: docs/KatibClient.md#create_experiment
[get_e]: docs/KatibClient.md#get_experiment
[get_s]: docs/KatibClient.md#get_suggestion
[delete]: docs/KatibClient.md#delete_experiment
[list_e]: docs/KatibClient.md#list_experiments
[get_status]: docs/KatibClient.md#get_experiment_status
[is_suc]: docs/KatibClient.md#is_experiment_succeeded
[list_t]: docs/KatibClient.md#list_trials
[get_suc_t]: docs/KatibClient.md#get_success_trial_details
[opt_hp]: docs/KatibClient.md#get_optimal_hyperparameters

## Documentation For Models

- [V1beta1AlgorithmSetting](docs/V1beta1AlgorithmSetting.md)
- [V1beta1AlgorithmSpec](docs/V1beta1AlgorithmSpec.md)
- [V1beta1CollectorSpec](docs/V1beta1CollectorSpec.md)
- [V1beta1ConfigMapSource](docs/V1beta1ConfigMapSource.md)
- [V1beta1EarlyStoppingRule](docs/V1beta1EarlyStoppingRule.md)
- [V1beta1EarlyStoppingSetting](docs/V1beta1EarlyStoppingSetting.md)
- [V1beta1EarlyStoppingSpec](docs/V1beta1EarlyStoppingSpec.md)
- [V1beta1Experiment](docs/V1beta1Experiment.md)
- [V1beta1ExperimentCondition](docs/V1beta1ExperimentCondition.md)
- [V1beta1ExperimentList](docs/V1beta1ExperimentList.md)
- [V1beta1ExperimentSpec](docs/V1beta1ExperimentSpec.md)
- [V1beta1ExperimentStatus](docs/V1beta1ExperimentStatus.md)
- [V1beta1FeasibleSpace](docs/V1beta1FeasibleSpace.md)
- [V1beta1FileSystemPath](docs/V1beta1FileSystemPath.md)
- [V1beta1FilterSpec](docs/V1beta1FilterSpec.md)
- [V1beta1GraphConfig](docs/V1beta1GraphConfig.md)
- [V1beta1Metric](docs/V1beta1Metric.md)
- [V1beta1MetricStrategy](docs/V1beta1MetricStrategy.md)
- [V1beta1MetricsCollectorSpec](docs/V1beta1MetricsCollectorSpec.md)
- [V1beta1NasConfig](docs/V1beta1NasConfig.md)
- [V1beta1ObjectiveSpec](docs/V1beta1ObjectiveSpec.md)
- [V1beta1Observation](docs/V1beta1Observation.md)
- [V1beta1Operation](docs/V1beta1Operation.md)
- [V1beta1OptimalTrial](docs/V1beta1OptimalTrial.md)
- [V1beta1ParameterAssignment](docs/V1beta1ParameterAssignment.md)
- [V1beta1ParameterSpec](docs/V1beta1ParameterSpec.md)
- [V1beta1SourceSpec](docs/V1beta1SourceSpec.md)
- [V1beta1Suggestion](docs/V1beta1Suggestion.md)
- [V1beta1SuggestionCondition](docs/V1beta1SuggestionCondition.md)
- [V1beta1SuggestionList](docs/V1beta1SuggestionList.md)
- [V1beta1SuggestionSpec](docs/V1beta1SuggestionSpec.md)
- [V1beta1SuggestionStatus](docs/V1beta1SuggestionStatus.md)
- [V1beta1Trial](docs/V1beta1Trial.md)
- [V1beta1TrialAssignment](docs/V1beta1TrialAssignment.md)
- [V1beta1TrialCondition](docs/V1beta1TrialCondition.md)
- [V1beta1TrialList](docs/V1beta1TrialList.md)
- [V1beta1TrialParameterSpec](docs/V1beta1TrialParameterSpec.md)
- [V1beta1TrialSource](docs/V1beta1TrialSource.md)
- [V1beta1TrialSpec](docs/V1beta1TrialSpec.md)
- [V1beta1TrialStatus](docs/V1beta1TrialStatus.md)
- [V1beta1TrialTemplate](docs/V1beta1TrialTemplate.md)

## Documentation For Authorization

All endpoints do not require authorization.

## Author

prem0912
