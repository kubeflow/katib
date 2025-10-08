# V1beta1SuggestionSpec

SuggestionSpec is the specification of a Suggestion.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**algorithm** | [**V1beta1AlgorithmSpec**](V1beta1AlgorithmSpec.md) |  | [optional] 
**early_stopping** | [**V1beta1EarlyStoppingSpec**](V1beta1EarlyStoppingSpec.md) |  | [optional] 
**requests** | **int** | Number of suggestions requested. | [optional] 
**resume_policy** | **str** | ResumePolicy describes resuming policy which usually take effect after experiment terminated. Default value is Never. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


