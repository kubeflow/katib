# V1beta1TrialTemplate

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**config_map** | [**V1beta1ConfigMapSource**](V1beta1ConfigMapSource.md) | ConfigMap spec represents a reference to ConfigMap | [optional] 
**retain** | **bool** | Retain indicates that Trial resources must be not cleanup | [optional] 
**trial_parameters** | [**list[V1beta1TrialParameterSpec]**](V1beta1TrialParameterSpec.md) | List of parameres that are used in Trial template | [optional] 
**trial_spec** | [**V1UnstructuredUnstructured**](V1UnstructuredUnstructured.md) | TrialSpec represents Trial template in unstructured format | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


