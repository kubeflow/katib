# V1beta1TrialStatus

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**completion_time** | **datetime** | Represents time when the Trial was completed. It is not guaranteed to be set in happens-before order across separate operations. It is represented in RFC3339 form and is in UTC | [optional] 
**conditions** | [**list[V1beta1TrialCondition]**](V1beta1TrialCondition.md) | List of observed runtime conditions for this Trial. | [optional] 
**last_reconcile_time** | **datetime** | Represents last time when the Trial was reconciled. It is not guaranteed to be set in happens-before order across separate operations. It is represented in RFC3339 form and is in UTC. | [optional] 
**observation** | [**V1beta1Observation**](V1beta1Observation.md) | Results of the Trial - objectives and other metrics values. | [optional] 
**start_time** | **datetime** | Represents time when the Trial was acknowledged by the Trial controller. It is not guaranteed to be set in happens-before order across separate operations. It is represented in RFC3339 form and is in UTC | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


