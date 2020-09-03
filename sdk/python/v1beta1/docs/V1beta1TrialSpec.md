# V1beta1TrialSpec

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**metrics_collector** | [**V1beta1MetricsCollectorSpec**](V1beta1MetricsCollectorSpec.md) | Describes how metrics will be collected | [optional] 
**objective** | [**V1beta1ObjectiveSpec**](V1beta1ObjectiveSpec.md) | Describes the objective of the experiment. | [optional] 
**parameter_assignments** | [**list[V1beta1ParameterAssignment]**](V1beta1ParameterAssignment.md) | Key-value pairs for hyperparameters and assignment values. | 
**primary_container_name** | **str** | Name of training container where actual model training is running | [optional] 
**primary_pod_labels** | **dict(str, str)** | Label that determines if pod needs to be injected by Katib sidecar container | [optional] 
**retain_run** | **bool** | Whether to retain the trial run object after completed. | [optional] 
**run_spec** | [**V1UnstructuredUnstructured**](V1UnstructuredUnstructured.md) | Raw text for the trial run spec. This can be any generic Kubernetes runtime object. The trial operator should create the resource as written, and let the corresponding resource controller (e.g. tf-operator) handle the rest. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


