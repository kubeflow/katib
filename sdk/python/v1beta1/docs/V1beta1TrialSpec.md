# V1beta1TrialSpec

TrialSpec is the specification of a Trial.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**early_stopping_rules** | [**list[V1beta1EarlyStoppingRule]**](V1beta1EarlyStoppingRule.md) | Rules for early stopping techniques. Each rule should be met to early stop Trial. | [optional] 
**failure_condition** | **str** | Condition when trial custom resource is failed. Condition must be in GJSON format, ref https://github.com/tidwall/gjson. For example for BatchJob: status.conditions.#(type&#x3D;&#x3D;\&quot;Failed\&quot;)#|#(status&#x3D;&#x3D;\&quot;True\&quot;)# | [optional] 
**labels** | **dict(str, str)** | Labels that provide additional metadata for services (e.g. Suggestions tracking) | [optional] 
**metrics_collector** | [**V1beta1MetricsCollectorSpec**](V1beta1MetricsCollectorSpec.md) |  | [optional] 
**objective** | [**V1beta1ObjectiveSpec**](V1beta1ObjectiveSpec.md) |  | [optional] 
**parameter_assignments** | [**list[V1beta1ParameterAssignment]**](V1beta1ParameterAssignment.md) | Key-value pairs for hyperparameters and assignment values. | [optional] 
**primary_container_name** | **str** | Name of training container where actual model training is running | [optional] 
**primary_pod_labels** | **dict(str, str)** | Label that determines if pod needs to be injected by Katib sidecar container | [optional] 
**retain_run** | **bool** | Whether to retain the trial run object after completed. | [optional] 
**run_spec** | **object** |  | [optional] 
**success_condition** | **str** | Condition when trial custom resource is succeeded. Condition must be in GJSON format, ref https://github.com/tidwall/gjson. For example for BatchJob: status.conditions.#(type&#x3D;&#x3D;\&quot;Complete\&quot;)#|#(status&#x3D;&#x3D;\&quot;True\&quot;)# | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


