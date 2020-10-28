# V1beta1SuggestionStatus

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**algorithm_settings** | [**list[V1beta1AlgorithmSetting]**](V1beta1AlgorithmSetting.md) | AlgorithmSettings defines HP or NAS algorithm settings which suggestion gRPC service returns. These settings overwrites Experiment&#39;s settings before the gRPC request. It can be empty if settings haven&#39;t been changed. | [optional] 
**completion_time** | [**V1Time**](V1Time.md) | Represents time when the Suggestion was completed. It is not guaranteed to be set in happens-before order across separate operations. It is represented in RFC3339 form and is in UTC. | [optional] 
**conditions** | [**list[V1beta1SuggestionCondition]**](V1beta1SuggestionCondition.md) | List of observed runtime conditions for this Suggestion. | [optional] 
**last_reconcile_time** | [**V1Time**](V1Time.md) | Represents last time when the Suggestion was reconciled. It is not guaranteed to be set in happens-before order across separate operations. It is represented in RFC3339 form and is in UTC. | [optional] 
**start_time** | [**V1Time**](V1Time.md) | Represents time when the Suggestion was acknowledged by the Suggestion controller. It is not guaranteed to be set in happens-before order across separate operations. It is represented in RFC3339 form and is in UTC. | [optional] 
**suggestion_count** | **int** | Number of suggestion results | [optional] 
**suggestions** | [**list[V1beta1TrialAssignment]**](V1beta1TrialAssignment.md) | Suggestion results | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


