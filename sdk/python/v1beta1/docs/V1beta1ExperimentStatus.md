# V1beta1ExperimentStatus

ExperimentStatus is the current status of an Experiment.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**completion_time** | **datetime** |  | [optional] 
**conditions** | [**list[V1beta1ExperimentCondition]**](V1beta1ExperimentCondition.md) | List of observed runtime conditions for this Experiment. | [optional] 
**current_optimal_trial** | [**V1beta1OptimalTrial**](V1beta1OptimalTrial.md) |  | [optional] 
**early_stopped_trial_list** | **list[str]** | List of trial names which have been early stopped. | [optional] 
**failed_trial_list** | **list[str]** | List of trial names which have already failed. | [optional] 
**killed_trial_list** | **list[str]** | List of trial names which have been killed. | [optional] 
**last_reconcile_time** | **datetime** |  | [optional] 
**metrics_unavailable_trial_list** | **list[str]** | List of trial names which have been metrics unavailable | [optional] 
**pending_trial_list** | **list[str]** | List of trial names which are pending. | [optional] 
**running_trial_list** | **list[str]** | List of trial names which are running. | [optional] 
**start_time** | **datetime** |  | [optional] 
**succeeded_trial_list** | **list[str]** | List of trial names which have already succeeded. | [optional] 
**trial_metrics_unavailable** | **int** | How many trials are currently metrics unavailable. | [optional] 
**trials** | **int** | Trials is the total number of trials owned by the experiment. | [optional] 
**trials_early_stopped** | **int** | How many trials are currently early stopped. | [optional] 
**trials_failed** | **int** | How many trials have failed. | [optional] 
**trials_killed** | **int** | How many trials have been killed. | [optional] 
**trials_pending** | **int** | How many trials are currently pending. | [optional] 
**trials_running** | **int** | How many trials are currently running. | [optional] 
**trials_succeeded** | **int** | How many trials have succeeded. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


