# V1beta1TrialTemplate

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**config_map** | [**V1beta1ConfigMapSource**](V1beta1ConfigMapSource.md) | ConfigMap spec represents a reference to ConfigMap | [optional] 
**failure_condition** | **str** | Condition when trial custom resource is failed. Condition must be in GJSON format, ref https://github.com/tidwall/gjson. For example for BatchJob: status.conditions.#(type&#x3D;&#x3D;\&quot;Failed\&quot;)#|#(status&#x3D;&#x3D;\&quot;True\&quot;)# | [optional] 
**primary_container_name** | **str** | Name of training container where actual model training is running | [optional] 
**primary_pod_labels** | **dict(str, str)** | Labels that determines if pod needs to be injected by Katib sidecar container | [optional] 
**retain** | **bool** | Retain indicates that trial resources must be not cleanup | [optional] 
**success_condition** | **str** | Condition when trial custom resource is succeeded. Condition must be in GJSON format, ref https://github.com/tidwall/gjson. For example for BatchJob: status.conditions.#(type&#x3D;&#x3D;\&quot;Complete\&quot;)#|#(status&#x3D;&#x3D;\&quot;True\&quot;)# | [optional] 
**trial_parameters** | [**list[V1beta1TrialParameterSpec]**](V1beta1TrialParameterSpec.md) | List of parameters that are used in trial template | [optional] 
**trial_spec** | [**V1UnstructuredUnstructured**](V1UnstructuredUnstructured.md) | TrialSpec represents trial template in unstructured format | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


