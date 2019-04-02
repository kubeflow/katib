# Protocol Documentation
<a name="top"/>

## Table of Contents

- [api.proto](#api.proto)
    - [AlgorithmSetting](#api.v1.alpha2.AlgorithmSetting)
    - [AlgorithmSpec](#api.v1.alpha2.AlgorithmSpec)
    - [DeleteExperimentReply](#api.v1.alpha2.DeleteExperimentReply)
    - [DeleteExperimentRequest](#api.v1.alpha2.DeleteExperimentRequest)
    - [DeleteTrialReply](#api.v1.alpha2.DeleteTrialReply)
    - [DeleteTrialRequest](#api.v1.alpha2.DeleteTrialRequest)
    - [EarlyStoppingSpec](#api.v1.alpha2.EarlyStoppingSpec)
    - [Experiment](#api.v1.alpha2.Experiment)
    - [ExperimentSpec](#api.v1.alpha2.ExperimentSpec)
    - [ExperimentStatus](#api.v1.alpha2.ExperimentStatus)
    - [ExperimentSummary](#api.v1.alpha2.ExperimentSummary)
    - [FeasibleSpace](#api.v1.alpha2.FeasibleSpace)
    - [GetAlgorithmExtraSettingsReply](#api.v1.alpha2.GetAlgorithmExtraSettingsReply)
    - [GetAlgorithmExtraSettingsRequest](#api.v1.alpha2.GetAlgorithmExtraSettingsRequest)
    - [GetExperimentListReply](#api.v1.alpha2.GetExperimentListReply)
    - [GetExperimentListRequest](#api.v1.alpha2.GetExperimentListRequest)
    - [GetExperimentReply](#api.v1.alpha2.GetExperimentReply)
    - [GetExperimentRequest](#api.v1.alpha2.GetExperimentRequest)
    - [GetObservationLogReply](#api.v1.alpha2.GetObservationLogReply)
    - [GetObservationLogRequest](#api.v1.alpha2.GetObservationLogRequest)
    - [GetSuggestionsReply](#api.v1.alpha2.GetSuggestionsReply)
    - [GetSuggestionsRequest](#api.v1.alpha2.GetSuggestionsRequest)
    - [GetTrialListReply](#api.v1.alpha2.GetTrialListReply)
    - [GetTrialListRequest](#api.v1.alpha2.GetTrialListRequest)
    - [GetTrialReply](#api.v1.alpha2.GetTrialReply)
    - [GetTrialRequest](#api.v1.alpha2.GetTrialRequest)
    - [GraphConfig](#api.v1.alpha2.GraphConfig)
    - [Metric](#api.v1.alpha2.Metric)
    - [MetricLog](#api.v1.alpha2.MetricLog)
    - [NasConfig](#api.v1.alpha2.NasConfig)
    - [NasConfig.Operations](#api.v1.alpha2.NasConfig.Operations)
    - [ObjectiveSpec](#api.v1.alpha2.ObjectiveSpec)
    - [Observation](#api.v1.alpha2.Observation)
    - [ObservationLog](#api.v1.alpha2.ObservationLog)
    - [Operation](#api.v1.alpha2.Operation)
    - [ParameterAssignment](#api.v1.alpha2.ParameterAssignment)
    - [ParameterSpec](#api.v1.alpha2.ParameterSpec)
    - [RegisterExperimentReply](#api.v1.alpha2.RegisterExperimentReply)
    - [RegisterExperimentRequest](#api.v1.alpha2.RegisterExperimentRequest)
    - [RegisterTrialReply](#api.v1.alpha2.RegisterTrialReply)
    - [RegisterTrialRequest](#api.v1.alpha2.RegisterTrialRequest)
    - [ReportObservationLogReply](#api.v1.alpha2.ReportObservationLogReply)
    - [ReportObservationLogRequest](#api.v1.alpha2.ReportObservationLogRequest)
    - [Trial](#api.v1.alpha2.Trial)
    - [TrialSpec](#api.v1.alpha2.TrialSpec)
    - [TrialStatus](#api.v1.alpha2.TrialStatus)
    - [UpdateAlgorithmExtraSettingsReply](#api.v1.alpha2.UpdateAlgorithmExtraSettingsReply)
    - [UpdateAlgorithmExtraSettingsRequest](#api.v1.alpha2.UpdateAlgorithmExtraSettingsRequest)
    - [UpdateExperimentStatusReply](#api.v1.alpha2.UpdateExperimentStatusReply)
    - [UpdateExperimentStatusRequest](#api.v1.alpha2.UpdateExperimentStatusRequest)
    - [UpdateTrialStatusReply](#api.v1.alpha2.UpdateTrialStatusReply)
    - [UpdateTrialStatusRequest](#api.v1.alpha2.UpdateTrialStatusRequest)
    - [ValidateAlgorithmSettingsReply](#api.v1.alpha2.ValidateAlgorithmSettingsReply)
    - [ValidateAlgorithmSettingsRequest](#api.v1.alpha2.ValidateAlgorithmSettingsRequest)
  
    - [ExperimentStatus.ExperimentConditionType](#api.v1.alpha2.ExperimentStatus.ExperimentConditionType)
    - [ObjectiveType](#api.v1.alpha2.ObjectiveType)
    - [ParameterType](#api.v1.alpha2.ParameterType)
    - [TrialStatus.TrialConditionType](#api.v1.alpha2.TrialStatus.TrialConditionType)
  
  
    - [EarlyStopping](#api.v1.alpha2.EarlyStopping)
    - [Manager](#api.v1.alpha2.Manager)
    - [Suggestion](#api.v1.alpha2.Suggestion)
  

- [Scalar Value Types](#scalar-value-types)



<a name="api.proto"/>
<p align="right"><a href="#top">Top</a></p>

## api.proto
Katib API


<a name="api.v1.alpha2.AlgorithmSetting"/>

### AlgorithmSetting



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="api.v1.alpha2.AlgorithmSpec"/>

### AlgorithmSpec



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| algorithm_name | [string](#string) |  |  |
| algorithm_setting | [AlgorithmSetting](#api.v1.alpha2.AlgorithmSetting) | repeated |  |
| early_stopping_spec | [EarlyStoppingSpec](#api.v1.alpha2.EarlyStoppingSpec) |  |  |






<a name="api.v1.alpha2.DeleteExperimentReply"/>

### DeleteExperimentReply







<a name="api.v1.alpha2.DeleteExperimentRequest"/>

### DeleteExperimentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| experiment_name | [string](#string) |  |  |






<a name="api.v1.alpha2.DeleteTrialReply"/>

### DeleteTrialReply







<a name="api.v1.alpha2.DeleteTrialRequest"/>

### DeleteTrialRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trial_name | [string](#string) |  |  |






<a name="api.v1.alpha2.EarlyStoppingSpec"/>

### EarlyStoppingSpec
TODO






<a name="api.v1.alpha2.Experiment"/>

### Experiment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of Experiment. This is unique in DB. |
| experiment_spec | [ExperimentSpec](#api.v1.alpha2.ExperimentSpec) |  |  |
| experiment_status | [ExperimentStatus](#api.v1.alpha2.ExperimentStatus) |  |  |






<a name="api.v1.alpha2.ExperimentSpec"/>

### ExperimentSpec
Spec of a Experiment. Experiment represents a single optimization run over a feasible space. 
Each Experiment contains a configuration describing the feasible space, as well as a set of Trials.
It is assumed that objective function f(x) does not change in the course of a Experiment.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| parameters | [ParameterSpec](#api.v1.alpha2.ParameterSpec) | repeated |  |
| objective | [ObjectiveSpec](#api.v1.alpha2.ObjectiveSpec) |  |  |
| algorithm | [AlgorithmSpec](#api.v1.alpha2.AlgorithmSpec) |  |  |
| trial_template | [string](#string) |  |  |
| parallel_trial_count | [int32](#int32) |  |  |
| max_trial_count | [int32](#int32) |  |  |
| metrics_collector_type | [string](#string) |  |  |
| nas_config | [NasConfig](#api.v1.alpha2.NasConfig) |  |  |






<a name="api.v1.alpha2.ExperimentStatus"/>

### ExperimentStatus



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| start_time | [string](#string) |  | RFC3339 format |
| completion_time | [string](#string) |  | RFC3339 format |
| condition | [ExperimentStatus.ExperimentConditionType](#api.v1.alpha2.ExperimentStatus.ExperimentConditionType) |  |  |






<a name="api.v1.alpha2.ExperimentSummary"/>

### ExperimentSummary



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| experiment_name | [string](#string) |  |  |
| status | [ExperimentStatus](#api.v1.alpha2.ExperimentStatus) |  |  |






<a name="api.v1.alpha2.FeasibleSpace"/>

### FeasibleSpace
Feasible space for optimization.
Int and Double type use Max/Min.
Discrete and Categorical type use List.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| max | [string](#string) |  | Max Value |
| min | [string](#string) |  | Minimum Value |
| list | [string](#string) | repeated | List of Values. |
| step | [string](#string) |  | Step for double or int parameter |






<a name="api.v1.alpha2.GetAlgorithmExtraSettingsReply"/>

### GetAlgorithmExtraSettingsReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| extra_algorithm_settings | [AlgorithmSetting](#api.v1.alpha2.AlgorithmSetting) | repeated |  |






<a name="api.v1.alpha2.GetAlgorithmExtraSettingsRequest"/>

### GetAlgorithmExtraSettingsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| experiment_name | [string](#string) |  |  |






<a name="api.v1.alpha2.GetExperimentListReply"/>

### GetExperimentListReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| experiment_summaries | [ExperimentSummary](#api.v1.alpha2.ExperimentSummary) | repeated |  |






<a name="api.v1.alpha2.GetExperimentListRequest"/>

### GetExperimentListRequest







<a name="api.v1.alpha2.GetExperimentReply"/>

### GetExperimentReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| experiment | [Experiment](#api.v1.alpha2.Experiment) |  |  |






<a name="api.v1.alpha2.GetExperimentRequest"/>

### GetExperimentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| experiment_name | [string](#string) |  |  |






<a name="api.v1.alpha2.GetObservationLogReply"/>

### GetObservationLogReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| observation_log | [ObservationLog](#api.v1.alpha2.ObservationLog) |  |  |






<a name="api.v1.alpha2.GetObservationLogRequest"/>

### GetObservationLogRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trial_name | [string](#string) |  |  |






<a name="api.v1.alpha2.GetSuggestionsReply"/>

### GetSuggestionsReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trials | [Trial](#api.v1.alpha2.Trial) | repeated |  |






<a name="api.v1.alpha2.GetSuggestionsRequest"/>

### GetSuggestionsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| experiment_name | [string](#string) |  |  |
| algorithm_name | [string](#string) |  |  |
| request_number | [int32](#int32) |  | The number of Suggestion you request at one time. When you set 3 to request_number, you can get three Suggestions at one time. |






<a name="api.v1.alpha2.GetTrialListReply"/>

### GetTrialListReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trials | [Trial](#api.v1.alpha2.Trial) | repeated |  |






<a name="api.v1.alpha2.GetTrialListRequest"/>

### GetTrialListRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| experiment_name | [string](#string) |  |  |
| start_time | [string](#string) |  | The start of the time range. RFC3339 format |
| end_time | [string](#string) |  | The end of the time range. RFC3339 format |
| filter | [string](#string) |  |  |






<a name="api.v1.alpha2.GetTrialReply"/>

### GetTrialReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trial | [Trial](#api.v1.alpha2.Trial) |  |  |






<a name="api.v1.alpha2.GetTrialRequest"/>

### GetTrialRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trial_name | [string](#string) |  |  |






<a name="api.v1.alpha2.GraphConfig"/>

### GraphConfig
GraphConfig contains a config of DAG


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| num_layers | [int32](#int32) |  | Number of layers |
| input_sizes | [int32](#int32) | repeated | Dimensions of input size |
| output_sizes | [int32](#int32) | repeated | Dimensions of output size |






<a name="api.v1.alpha2.Metric"/>

### Metric



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="api.v1.alpha2.MetricLog"/>

### MetricLog



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| time_stamp | [string](#string) |  | RFC3339 format |
| metric | [Metric](#api.v1.alpha2.Metric) |  |  |






<a name="api.v1.alpha2.NasConfig"/>

### NasConfig
NasConfig contains a config of NAS job


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| graph_config | [GraphConfig](#api.v1.alpha2.GraphConfig) |  | Config of DAG |
| operations | [NasConfig.Operations](#api.v1.alpha2.NasConfig.Operations) |  | List of Operation |






<a name="api.v1.alpha2.NasConfig.Operations"/>

### NasConfig.Operations



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operation | [Operation](#api.v1.alpha2.Operation) | repeated |  |






<a name="api.v1.alpha2.ObjectiveSpec"/>

### ObjectiveSpec



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [ObjectiveType](#api.v1.alpha2.ObjectiveType) |  |  |
| goal | [float](#float) |  |  |
| objective_metric_name | [string](#string) |  |  |
| additional_metrics_names | [string](#string) | repeated | This can be empty if we only care about the objective metric. |






<a name="api.v1.alpha2.Observation"/>

### Observation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| metrics | [Metric](#api.v1.alpha2.Metric) | repeated |  |






<a name="api.v1.alpha2.ObservationLog"/>

### ObservationLog



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| metric_logs | [MetricLog](#api.v1.alpha2.MetricLog) | repeated |  |






<a name="api.v1.alpha2.Operation"/>

### Operation
Config for operations in DAG


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operation_type | [string](#string) |  | Type of operation in DAG |
| parameters | [ParameterSpec](#api.v1.alpha2.ParameterSpec) | repeated | List of ParameterSpec

/ List of ParameterSpec |






<a name="api.v1.alpha2.ParameterAssignment"/>

### ParameterAssignment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="api.v1.alpha2.ParameterSpec"/>

### ParameterSpec
Config for a Hyper parameter.
Katib will create each Hyper parameter from this config.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the parameter. |
| parameter_type | [ParameterType](#api.v1.alpha2.ParameterType) |  | Type of the parameter. |
| feasible_space | [FeasibleSpace](#api.v1.alpha2.FeasibleSpace) |  | FeasibleSpace for the parameter. |






<a name="api.v1.alpha2.RegisterExperimentReply"/>

### RegisterExperimentReply







<a name="api.v1.alpha2.RegisterExperimentRequest"/>

### RegisterExperimentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| experiment | [Experiment](#api.v1.alpha2.Experiment) |  |  |






<a name="api.v1.alpha2.RegisterTrialReply"/>

### RegisterTrialReply







<a name="api.v1.alpha2.RegisterTrialRequest"/>

### RegisterTrialRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trial | [Trial](#api.v1.alpha2.Trial) |  |  |






<a name="api.v1.alpha2.ReportObservationLogReply"/>

### ReportObservationLogReply







<a name="api.v1.alpha2.ReportObservationLogRequest"/>

### ReportObservationLogRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trial_name | [string](#string) |  |  |
| observation_log | [ObservationLog](#api.v1.alpha2.ObservationLog) |  |  |






<a name="api.v1.alpha2.Trial"/>

### Trial



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| spec | [TrialSpec](#api.v1.alpha2.TrialSpec) |  |  |
| status | [TrialStatus](#api.v1.alpha2.TrialStatus) |  |  |






<a name="api.v1.alpha2.TrialSpec"/>

### TrialSpec



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| experiment_name | [string](#string) |  |  |
| parameter_assignments | [ParameterAssignment](#api.v1.alpha2.ParameterAssignment) | repeated |  |
| run_spec | [string](#string) |  |  |






<a name="api.v1.alpha2.TrialStatus"/>

### TrialStatus



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| start_time | [string](#string) |  | RFC3339 format |
| completion_time | [string](#string) |  | RFC3339 format |
| condition | [TrialStatus.TrialConditionType](#api.v1.alpha2.TrialStatus.TrialConditionType) |  |  |
| observation | [Observation](#api.v1.alpha2.Observation) |  | The best observation in logs. |






<a name="api.v1.alpha2.UpdateAlgorithmExtraSettingsReply"/>

### UpdateAlgorithmExtraSettingsReply







<a name="api.v1.alpha2.UpdateAlgorithmExtraSettingsRequest"/>

### UpdateAlgorithmExtraSettingsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| experiment_name | [string](#string) |  |  |
| extra_algorithm_settings | [AlgorithmSetting](#api.v1.alpha2.AlgorithmSetting) | repeated |  |






<a name="api.v1.alpha2.UpdateExperimentStatusReply"/>

### UpdateExperimentStatusReply







<a name="api.v1.alpha2.UpdateExperimentStatusRequest"/>

### UpdateExperimentStatusRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| experiment_name | [string](#string) |  |  |
| new_status | [ExperimentStatus](#api.v1.alpha2.ExperimentStatus) |  |  |






<a name="api.v1.alpha2.UpdateTrialStatusReply"/>

### UpdateTrialStatusReply







<a name="api.v1.alpha2.UpdateTrialStatusRequest"/>

### UpdateTrialStatusRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trial_name | [string](#string) |  |  |
| new_status | [TrialStatus](#api.v1.alpha2.TrialStatus) |  |  |






<a name="api.v1.alpha2.ValidateAlgorithmSettingsReply"/>

### ValidateAlgorithmSettingsReply
Return INVALID_ARGUMENT Error if Algorithm Settings are not Valid






<a name="api.v1.alpha2.ValidateAlgorithmSettingsRequest"/>

### ValidateAlgorithmSettingsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| experiment_spec | [ExperimentSpec](#api.v1.alpha2.ExperimentSpec) |  |  |
| algorithm_name | [string](#string) |  |  |





 


<a name="api.v1.alpha2.ExperimentStatus.ExperimentConditionType"/>

### ExperimentStatus.ExperimentConditionType


| Name | Number | Description |
| ---- | ------ | ----------- |
| CREATED | 0 |  |
| RUNNING | 1 |  |
| RESTARTING | 2 |  |
| SUCCEEDED | 3 |  |
| FAILED | 4 |  |



<a name="api.v1.alpha2.ObjectiveType"/>

### ObjectiveType
Direction of optimization. Minimize or Maximize.

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 | Undefined type and not used. |
| MINIMIZE | 1 | Minimize |
| MAXIMIZE | 2 | Maximize |



<a name="api.v1.alpha2.ParameterType"/>

### ParameterType
Types of value for HyperParameter.

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN_TYPE | 0 | Undefined type and not used. |
| DOUBLE | 1 | Double float type. Use &#34;Max/Min&#34;. |
| INT | 2 | Int type. Use &#34;Max/Min&#34;. |
| DISCRETE | 3 | Discrete number type. Use &#34;List&#34; as float. |
| CATEGORICAL | 4 | Categorical type. Use &#34;List&#34; as string. |



<a name="api.v1.alpha2.TrialStatus.TrialConditionType"/>

### TrialStatus.TrialConditionType


| Name | Number | Description |
| ---- | ------ | ----------- |
| PENDING | 0 |  |
| RUNNING | 1 |  |
| COMPLETED | 2 |  |
| KILLED | 3 |  |
| FAILED | 4 |  |


 

 


<a name="api.v1.alpha2.EarlyStopping"/>

### EarlyStopping
TODO

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|


<a name="api.v1.alpha2.Manager"/>

### Manager
Service for Main API for Katib
For each RPC service, we define mapping to HTTP REST API method.
The mapping includes the URL path, query parameters and request body.
https://cloud.google.com/service-infrastructure/docs/service-management/reference/rpc/google.api#http

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| RegisterExperiment | [RegisterExperimentRequest](#api.v1.alpha2.RegisterExperimentRequest) | [RegisterExperimentReply](#api.v1.alpha2.RegisterExperimentRequest) | Register a Experiment to DB. |
| DeleteExperiment | [DeleteExperimentRequest](#api.v1.alpha2.DeleteExperimentRequest) | [DeleteExperimentReply](#api.v1.alpha2.DeleteExperimentRequest) | Delete a Experiment from DB by name. |
| GetExperiment | [GetExperimentRequest](#api.v1.alpha2.GetExperimentRequest) | [GetExperimentReply](#api.v1.alpha2.GetExperimentRequest) | Get a Experiment from DB by name. |
| GetExperimentList | [GetExperimentListRequest](#api.v1.alpha2.GetExperimentListRequest) | [GetExperimentListReply](#api.v1.alpha2.GetExperimentListRequest) | Get a summary list of Experiment from DB. The summary includes name and condition. |
| UpdateExperimentStatus | [UpdateExperimentStatusRequest](#api.v1.alpha2.UpdateExperimentStatusRequest) | [UpdateExperimentStatusReply](#api.v1.alpha2.UpdateExperimentStatusRequest) | Update Status of a experiment. |
| UpdateAlgorithmExtraSettings | [UpdateAlgorithmExtraSettingsRequest](#api.v1.alpha2.UpdateAlgorithmExtraSettingsRequest) | [UpdateAlgorithmExtraSettingsReply](#api.v1.alpha2.UpdateAlgorithmExtraSettingsRequest) | Update AlgorithmExtraSettings. The ExtraSetting is created if it does not exist, otherwise it is overwrited. |
| GetAlgorithmExtraSettings | [GetAlgorithmExtraSettingsRequest](#api.v1.alpha2.GetAlgorithmExtraSettingsRequest) | [GetAlgorithmExtraSettingsReply](#api.v1.alpha2.GetAlgorithmExtraSettingsRequest) | Get all AlgorithmExtraSettings. |
| RegisterTrial | [RegisterTrialRequest](#api.v1.alpha2.RegisterTrialRequest) | [RegisterTrialReply](#api.v1.alpha2.RegisterTrialRequest) | Register a Trial to DB. ID will be filled by manager automatically. |
| DeleteTrial | [DeleteTrialRequest](#api.v1.alpha2.DeleteTrialRequest) | [DeleteTrialReply](#api.v1.alpha2.DeleteTrialRequest) | Delete a Trial from DB by ID. |
| GetTrialList | [GetTrialListRequest](#api.v1.alpha2.GetTrialListRequest) | [GetTrialListReply](#api.v1.alpha2.GetTrialListRequest) | Get a list of Trial from DB by name of a Experiment. |
| GetTrial | [GetTrialRequest](#api.v1.alpha2.GetTrialRequest) | [GetTrialReply](#api.v1.alpha2.GetTrialRequest) | Get a Trial from DB by ID of Trial. |
| UpdateTrialStatus | [UpdateTrialStatusRequest](#api.v1.alpha2.UpdateTrialStatusRequest) | [UpdateTrialStatusReply](#api.v1.alpha2.UpdateTrialStatusRequest) | Update Status of a trial. |
| ReportObservationLog | [ReportObservationLogRequest](#api.v1.alpha2.ReportObservationLogRequest) | [ReportObservationLogReply](#api.v1.alpha2.ReportObservationLogRequest) | Report a log of Observations for a Trial. The log consists of timestamp and value of metric. Katib store every log of metrics. You can see accuracy curve or other metric logs on UI. |
| GetObservationLog | [GetObservationLogRequest](#api.v1.alpha2.GetObservationLogRequest) | [GetObservationLogReply](#api.v1.alpha2.GetObservationLogRequest) | Get all log of Observations for a Trial. |
| GetSuggestions | [GetSuggestionsRequest](#api.v1.alpha2.GetSuggestionsRequest) | [GetSuggestionsReply](#api.v1.alpha2.GetSuggestionsRequest) | Get Suggestions from a Suggestion service. |
| ValidateAlgorithmSettings | [ValidateAlgorithmSettingsRequest](#api.v1.alpha2.ValidateAlgorithmSettingsRequest) | [ValidateAlgorithmSettingsReply](#api.v1.alpha2.ValidateAlgorithmSettingsRequest) | Validate AlgorithmSettings in an Experiment. Suggestion service should return INVALID_ARGUMENT Error when the parameter is invalid |


<a name="api.v1.alpha2.Suggestion"/>

### Suggestion


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetSuggestions | [GetSuggestionsRequest](#api.v1.alpha2.GetSuggestionsRequest) | [GetSuggestionsReply](#api.v1.alpha2.GetSuggestionsRequest) |  |
| ValidateAlgorithmSettings | [ValidateAlgorithmSettingsRequest](#api.v1.alpha2.ValidateAlgorithmSettingsRequest) | [ValidateAlgorithmSettingsReply](#api.v1.alpha2.ValidateAlgorithmSettingsRequest) |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | --------- | ----------- |
| <a name="double" /> double |  | double | double | float |
| <a name="float" /> float |  | float | float | float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long |
| <a name="bool" /> bool |  | bool | boolean | boolean |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str |

