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

 - [V1alpha3AlgorithmSetting](docs/V1alpha3AlgorithmSetting.md)
 - [V1alpha3AlgorithmSpec](docs/V1alpha3AlgorithmSpec.md)
 - [V1alpha3CollectorSpec](docs/V1alpha3CollectorSpec.md)
 - [V1alpha3EarlyStoppingSetting](docs/V1alpha3EarlyStoppingSetting.md)
 - [V1alpha3EarlyStoppingSpec](docs/V1alpha3EarlyStoppingSpec.md)
 - [V1alpha3Experiment](docs/V1alpha3Experiment.md)
 - [V1alpha3ExperimentCondition](docs/V1alpha3ExperimentCondition.md)
 - [V1alpha3ExperimentList](docs/V1alpha3ExperimentList.md)
 - [V1alpha3ExperimentSpec](docs/V1alpha3ExperimentSpec.md)
 - [V1alpha3ExperimentStatus](docs/V1alpha3ExperimentStatus.md)
 - [V1alpha3FeasibleSpace](docs/V1alpha3FeasibleSpace.md)
 - [V1alpha3FileSystemPath](docs/V1alpha3FileSystemPath.md)
 - [V1alpha3FilterSpec](docs/V1alpha3FilterSpec.md)
 - [V1alpha3GoTemplate](docs/V1alpha3GoTemplate.md)
 - [V1alpha3GraphConfig](docs/V1alpha3GraphConfig.md)
 - [V1alpha3Metric](docs/V1alpha3Metric.md)
 - [V1alpha3MetricsCollectorSpec](docs/V1alpha3MetricsCollectorSpec.md)
 - [V1alpha3NasConfig](docs/V1alpha3NasConfig.md)
 - [V1alpha3ObjectiveSpec](docs/V1alpha3ObjectiveSpec.md)
 - [V1alpha3Observation](docs/V1alpha3Observation.md)
 - [V1alpha3Operation](docs/V1alpha3Operation.md)
 - [V1alpha3OptimalTrial](docs/V1alpha3OptimalTrial.md)
 - [V1alpha3ParameterAssignment](docs/V1alpha3ParameterAssignment.md)
 - [V1alpha3ParameterSpec](docs/V1alpha3ParameterSpec.md)
 - [V1alpha3SourceSpec](docs/V1alpha3SourceSpec.md)
 - [V1alpha3Suggestion](docs/V1alpha3Suggestion.md)
 - [V1alpha3SuggestionCondition](docs/V1alpha3SuggestionCondition.md)
 - [V1alpha3SuggestionList](docs/V1alpha3SuggestionList.md)
 - [V1alpha3SuggestionSpec](docs/V1alpha3SuggestionSpec.md)
 - [V1alpha3SuggestionStatus](docs/V1alpha3SuggestionStatus.md)
 - [V1alpha3TemplateSpec](docs/V1alpha3TemplateSpec.md)
 - [V1alpha3Trial](docs/V1alpha3Trial.md)
 - [V1alpha3TrialAssignment](docs/V1alpha3TrialAssignment.md)
 - [V1alpha3TrialCondition](docs/V1alpha3TrialCondition.md)
 - [V1alpha3TrialList](docs/V1alpha3TrialList.md)
 - [V1alpha3TrialSpec](docs/V1alpha3TrialSpec.md)
 - [V1alpha3TrialStatus](docs/V1alpha3TrialStatus.md)
 - [V1alpha3TrialTemplate](docs/V1alpha3TrialTemplate.md)


## Documentation For Authorization

 All endpoints do not require authorization.


## Author

prem0912
