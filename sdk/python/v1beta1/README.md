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

## Getting Started

Please follow the [samples](examples/bayesianoptimization-katib-sdk.ipynb) to create, update, delete and get hyperparamaters of  Katib Experiment.

## Documentation for API Endpoints

Class | Method | Description
------------ | -------------  | -------------
[KatibClient](docs/KatibClient.md) | [create_experiment](docs/KatibClient.md#create_experiment) | Create Katib Experiment|
[KatibClient](docs/KatibClient.md) | [get_experiment](docs/KatibClient.md#get_experiment)    | Get or watch the specified Experiment or all Experiment in the namespace |
[KatibClient](docs/KatibClient.md) | [delete_experiment](docs/KatibClient.md#delete_experiment) | Delete specified Experiment |
[KatibClient](docs/KatibClient.md) | [list_experiments](docs/KatibClient.md#list_experiments) | List all Experiments with status |
[KatibClient](docs/KatibClient.md) | [get_experiment_status](docs/KatibClient.md#get_experiment_status) | Get Experiment status|
[KatibClient](docs/KatibClient.md) | [is_experiment_succeeded](docs/KatibClient.md#is_experiment_succeeded) | Check if Experiment status is Succeeded |
[KatibClient](docs/KatibClient.md) | [list_trials](docs/KatibClient.md#list_trials) | List all trials of specified Experiment |
[KatibClient](docs/KatibClient.md) | [get_optimal_hyperparmeters](docs/KatibClient.md#get_optimal_hyperparmeters) | Get currentOptimalTrial with paramaterAssignments of an Experiment|


## Documentation For Models

 - [V1beta1AlgorithmSetting](docs/V1beta1AlgorithmSetting.md)
 - [V1beta1AlgorithmSpec](docs/V1beta1AlgorithmSpec.md)
 - [V1beta1CollectorSpec](docs/V1beta1CollectorSpec.md)
 - [V1beta1ConfigMapSource](docs/V1beta1ConfigMapSource.md)
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
