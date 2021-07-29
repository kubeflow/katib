# V1beta1EarlyStoppingRule

EarlyStoppingRule represents each rule for early stopping.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**comparison** | **str** | Comparison defines correlation between name and value. | [optional] 
**name** | **str** | Name contains metric name for the rule. | [optional] 
**start_step** | **int** | StartStep defines quantity of intermediate results that should be received before applying the rule. If start step is empty, rule is applied from the first recorded metric. | [optional] 
**value** | **str** | Value contains metric value for the rule. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


