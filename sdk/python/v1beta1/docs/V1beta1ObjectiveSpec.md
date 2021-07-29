# V1beta1ObjectiveSpec

ObjectiveSpec represents Experiment's objective specification.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**additional_metric_names** | **list[str]** | AdditionalMetricNames represents metrics that should be collected from Trials. This can be empty if we only care about the objective metric. Note: If we adopt a push instead of pull mechanism, this can be omitted completely. | [optional] 
**goal** | **float** | Goal is the Experiment&#39;s objective goal that should be reached. In case of empty goal, Experiment is running until MaxTrialCount &#x3D; TrialsSucceeded. | [optional] 
**metric_strategies** | [**list[V1beta1MetricStrategy]**](V1beta1MetricStrategy.md) | MetricStrategies defines various rules (min, max or latest) to extract metrics values. This field is allowed to missing, experiment defaulter (webhook) will fill it. | [optional] 
**objective_metric_name** | **str** | ObjectiveMetricName represents primary Experiment&#39;s metric to optimize. | [optional] 
**type** | **str** | Type for Experiment optimization. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


