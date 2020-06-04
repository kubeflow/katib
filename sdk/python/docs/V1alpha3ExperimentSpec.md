# V1alpha3ExperimentSpec

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**algorithm** | [**V1alpha3AlgorithmSpec**](V1alpha3AlgorithmSpec.md) | Describes the suggestion algorithm. | [optional] 
**max_failed_trial_count** | **int** | Max failed trials to mark experiment as failed. | [optional] 
**max_trial_count** | **int** | Max completed trials to mark experiment as succeeded | [optional] 
**metrics_collector_spec** | [**V1alpha3MetricsCollectorSpec**](V1alpha3MetricsCollectorSpec.md) | For v1alpha3 we will keep the metrics collector implementation same as v1alpha1. | [optional] 
**nas_config** | [**V1alpha3NasConfig**](V1alpha3NasConfig.md) |  | [optional] 
**objective** | [**V1alpha3ObjectiveSpec**](V1alpha3ObjectiveSpec.md) | Describes the objective of the experiment. | [optional] 
**parallel_trial_count** | **int** | How many trials can be processed in parallel. Defaults to 3 | [optional] 
**parameters** | [**list[V1alpha3ParameterSpec]**](V1alpha3ParameterSpec.md) | List of hyperparameter configurations. | [optional] 
**trial_template** | [**V1alpha3TrialTemplate**](V1alpha3TrialTemplate.md) | Template for each run of the trial. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


