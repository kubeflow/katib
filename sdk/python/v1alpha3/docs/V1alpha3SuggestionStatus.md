# V1alpha3SuggestionStatus

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**algorithm_settings** | [**list[V1alpha3AlgorithmSetting]**](V1alpha3AlgorithmSetting.md) | Algorithmsettings set by the algorithm services. | [optional] 
**completion_time** | [**V1Time**](V1Time.md) | Represents time when the Suggestion was completed. It is not guaranteed to be set in happens-before order across separate operations. It is represented in RFC3339 form and is in UTC. | [optional] 
**conditions** | [**list[V1alpha3SuggestionCondition]**](V1alpha3SuggestionCondition.md) | List of observed runtime conditions for this Suggestion. | [optional] 
**last_reconcile_time** | [**V1Time**](V1Time.md) | Represents last time when the Suggestion was reconciled. It is not guaranteed to be set in happens-before order across separate operations. It is represented in RFC3339 form and is in UTC. | [optional] 
**start_time** | [**V1Time**](V1Time.md) | Represents time when the Suggestion was acknowledged by the Suggestion controller. It is not guaranteed to be set in happens-before order across separate operations. It is represented in RFC3339 form and is in UTC. | [optional] 
**suggestion_count** | **int** | Number of suggestion results | [optional] 
**suggestions** | [**list[V1alpha3TrialAssignment]**](V1alpha3TrialAssignment.md) | Suggestion results | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


