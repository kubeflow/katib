# V1beta1ExperimentSpec

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**algorithm** | [**V1beta1AlgorithmSpec**](V1beta1AlgorithmSpec.md) | Describes the suggestion algorithm. | [optional] 
**early_stopping** | [**V1beta1EarlyStoppingSpec**](V1beta1EarlyStoppingSpec.md) | Describes the early stopping algorithm. | [optional] 
**max_failed_trial_count** | **int** | Max failed trials to mark experiment as failed. | [optional] 
**max_trial_count** | **int** | Max completed trials to mark experiment as succeeded | [optional] 
**metrics_collector_spec** | [**V1beta1MetricsCollectorSpec**](V1beta1MetricsCollectorSpec.md) | Describes the specification of the metrics collector | [optional] 
**nas_config** | [**V1beta1NasConfig**](V1beta1NasConfig.md) |  | [optional] 
**objective** | [**V1beta1ObjectiveSpec**](V1beta1ObjectiveSpec.md) | Describes the objective of the experiment. | [optional] 
**parallel_trial_count** | **int** | How many trials can be processed in parallel. Defaults to 3 | [optional] 
**parameters** | [**list[V1beta1ParameterSpec]**](V1beta1ParameterSpec.md) | List of hyperparameter configurations. | [optional] 
**resume_policy** | **str** | Describes resuming policy which usually take effect after experiment terminated. | [optional] 
**trial_template** | [**V1beta1TrialTemplate**](V1beta1TrialTemplate.md) | Template for each run of the trial. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


