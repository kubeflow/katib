# V1alpha3ExperimentStatus

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**completion_time** | [**V1Time**](V1Time.md) | Represents time when the Experiment was completed. It is not guaranteed to be set in happens-before order across separate operations. It is represented in RFC3339 form and is in UTC. | [optional] 
**conditions** | [**list[V1alpha3ExperimentCondition]**](V1alpha3ExperimentCondition.md) | List of observed runtime conditions for this Experiment. | [optional] 
**current_optimal_trial** | [**V1alpha3OptimalTrial**](V1alpha3OptimalTrial.md) | Current optimal trial parameters and observations. | [optional] 
**failed_trial_list** | **list[str]** | List of trial names which have already failed. | [optional] 
**killed_trial_list** | **list[str]** | List of trial names which have been killed. | [optional] 
**last_reconcile_time** | [**V1Time**](V1Time.md) | Represents last time when the Experiment was reconciled. It is not guaranteed to be set in happens-before order across separate operations. It is represented in RFC3339 form and is in UTC. | [optional] 
**pending_trial_list** | **list[str]** | List of trial names which are pending. | [optional] 
**running_trial_list** | **list[str]** | List of trial names which are running. | [optional] 
**start_time** | [**V1Time**](V1Time.md) | Represents time when the Experiment was acknowledged by the Experiment controller. It is not guaranteed to be set in happens-before order across separate operations. It is represented in RFC3339 form and is in UTC. | [optional] 
**succeeded_trial_list** | **list[str]** | List of trial names which have already succeeded. | [optional] 
**trials** | **int** | Trials is the total number of trials owned by the experiment. | [optional] 
**trials_failed** | **int** | How many trials have failed. | [optional] 
**trials_killed** | **int** | How many trials have been killed. | [optional] 
**trials_pending** | **int** | How many trials are currently pending. | [optional] 
**trials_running** | **int** | How many trials are currently running. | [optional] 
**trials_succeeded** | **int** | How many trials have succeeded. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


