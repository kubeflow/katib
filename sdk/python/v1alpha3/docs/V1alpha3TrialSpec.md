# V1alpha3TrialSpec

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**metrics_collector** | [**V1alpha3MetricsCollectorSpec**](V1alpha3MetricsCollectorSpec.md) | Describes how metrics will be collected | [optional] 
**objective** | [**V1alpha3ObjectiveSpec**](V1alpha3ObjectiveSpec.md) | Describes the objective of the experiment. | [optional] 
**parameter_assignments** | [**list[V1alpha3ParameterAssignment]**](V1alpha3ParameterAssignment.md) | Key-value pairs for hyperparameters and assignment values. | 
**retain_run** | **bool** | Whether to retain the trial run object after completed. | [optional] 
**run_spec** | **str** | Raw text for the trial run spec. This can be any generic Kubernetes runtime object. The trial operator should create the resource as written, and let the corresponding resource controller (e.g. tf-operator) handle the rest. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


