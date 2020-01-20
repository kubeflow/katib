
# Latest version of the docs

## For the most up-to-date docs, see the [Katib API reference docs](https://www.kubeflow.org/docs/reference/katib/).

# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api.proto](#api.proto)
    - [AlgorithmSetting](#api.v1.alpha3.AlgorithmSetting)
    - [AlgorithmSpec](#api.v1.alpha3.AlgorithmSpec)
    - [DeleteObservationLogReply](#api.v1.alpha3.DeleteObservationLogReply)
    - [DeleteObservationLogRequest](#api.v1.alpha3.DeleteObservationLogRequest)
    - [EarlyStoppingSpec](#api.v1.alpha3.EarlyStoppingSpec)
    - [Experiment](#api.v1.alpha3.Experiment)
    - [ExperimentSpec](#api.v1.alpha3.ExperimentSpec)
    - [ExperimentSpec.ParameterSpecs](#api.v1.alpha3.ExperimentSpec.ParameterSpecs)
    - [FeasibleSpace](#api.v1.alpha3.FeasibleSpace)
    - [GetObservationLogReply](#api.v1.alpha3.GetObservationLogReply)
    - [GetObservationLogRequest](#api.v1.alpha3.GetObservationLogRequest)
    - [GetSuggestionsReply](#api.v1.alpha3.GetSuggestionsReply)
    - [GetSuggestionsReply.ParameterAssignments](#api.v1.alpha3.GetSuggestionsReply.ParameterAssignments)
    - [GetSuggestionsRequest](#api.v1.alpha3.GetSuggestionsRequest)
    - [GraphConfig](#api.v1.alpha3.GraphConfig)
    - [Metric](#api.v1.alpha3.Metric)
    - [MetricLog](#api.v1.alpha3.MetricLog)
    - [NasConfig](#api.v1.alpha3.NasConfig)
    - [NasConfig.Operations](#api.v1.alpha3.NasConfig.Operations)
    - [ObjectiveSpec](#api.v1.alpha3.ObjectiveSpec)
    - [Observation](#api.v1.alpha3.Observation)
    - [ObservationLog](#api.v1.alpha3.ObservationLog)
    - [Operation](#api.v1.alpha3.Operation)
    - [Operation.ParameterSpecs](#api.v1.alpha3.Operation.ParameterSpecs)
    - [ParameterAssignment](#api.v1.alpha3.ParameterAssignment)
    - [ParameterSpec](#api.v1.alpha3.ParameterSpec)
    - [ReportObservationLogReply](#api.v1.alpha3.ReportObservationLogReply)
    - [ReportObservationLogRequest](#api.v1.alpha3.ReportObservationLogRequest)
    - [Trial](#api.v1.alpha3.Trial)
    - [TrialSpec](#api.v1.alpha3.TrialSpec)
    - [TrialSpec.ParameterAssignments](#api.v1.alpha3.TrialSpec.ParameterAssignments)
    - [TrialStatus](#api.v1.alpha3.TrialStatus)
    - [ValidateAlgorithmSettingsReply](#api.v1.alpha3.ValidateAlgorithmSettingsReply)
    - [ValidateAlgorithmSettingsRequest](#api.v1.alpha3.ValidateAlgorithmSettingsRequest)
  
    - [ObjectiveType](#api.v1.alpha3.ObjectiveType)
    - [ParameterType](#api.v1.alpha3.ParameterType)
    - [TrialStatus.TrialConditionType](#api.v1.alpha3.TrialStatus.TrialConditionType)
  
  
    - [EarlyStopping](#api.v1.alpha3.EarlyStopping)
    - [Manager](#api.v1.alpha3.Manager)
    - [Suggestion](#api.v1.alpha3.Suggestion)
  

- [Scalar Value Types](#scalar-value-types)



<a name="api.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api.proto
Katib API


<a name="api.v1.alpha3.AlgorithmSetting"></a>

### AlgorithmSetting



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="api.v1.alpha3.AlgorithmSpec"></a>

### AlgorithmSpec



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| algorithm_name | [string](#string) |  |  |
| algorithm_setting | [AlgorithmSetting](#api.v1.alpha3.AlgorithmSetting) | repeated |  |
| early_stopping_spec | [EarlyStoppingSpec](#api.v1.alpha3.EarlyStoppingSpec) |  |  |






<a name="api.v1.alpha3.DeleteObservationLogReply"></a>

### DeleteObservationLogReply







<a name="api.v1.alpha3.DeleteObservationLogRequest"></a>

### DeleteObservationLogRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trial_name | [string](#string) |  |  |






<a name="api.v1.alpha3.EarlyStoppingSpec"></a>

### EarlyStoppingSpec
TODO






<a name="api.v1.alpha3.Experiment"></a>

### Experiment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of Experiment. This is unique in DB. |
| spec | [ExperimentSpec](#api.v1.alpha3.ExperimentSpec) |  |  |






<a name="api.v1.alpha3.ExperimentSpec"></a>

### ExperimentSpec
Spec of a Experiment. Experiment represents a single optimization run over a feasible space. 
Each Experiment contains a configuration describing the feasible space, as well as a set of Trials.
It is assumed that objective function f(x) does not change in the course of a Experiment.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| parameter_specs | [ExperimentSpec.ParameterSpecs](#api.v1.alpha3.ExperimentSpec.ParameterSpecs) |  |  |
| objective | [ObjectiveSpec](#api.v1.alpha3.ObjectiveSpec) |  |  |
| algorithm | [AlgorithmSpec](#api.v1.alpha3.AlgorithmSpec) |  |  |
| trial_template | [string](#string) |  |  |
| metrics_collector_spec | [string](#string) |  |  |
| parallel_trial_count | [int32](#int32) |  |  |
| max_trial_count | [int32](#int32) |  |  |
| nas_config | [NasConfig](#api.v1.alpha3.NasConfig) |  |  |






<a name="api.v1.alpha3.ExperimentSpec.ParameterSpecs"></a>

### ExperimentSpec.ParameterSpecs
List of ParameterSpec


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| parameters | [ParameterSpec](#api.v1.alpha3.ParameterSpec) | repeated |  |






<a name="api.v1.alpha3.FeasibleSpace"></a>

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






<a name="api.v1.alpha3.GetObservationLogReply"></a>

### GetObservationLogReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| observation_log | [ObservationLog](#api.v1.alpha3.ObservationLog) |  |  |






<a name="api.v1.alpha3.GetObservationLogRequest"></a>

### GetObservationLogRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trial_name | [string](#string) |  |  |
| metric_name | [string](#string) |  |  |
| start_time | [string](#string) |  | The start of the time range. RFC3339 format |
| end_time | [string](#string) |  | The end of the time range. RFC3339 format |






<a name="api.v1.alpha3.GetSuggestionsReply"></a>

### GetSuggestionsReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| parameter_assignments | [GetSuggestionsReply.ParameterAssignments](#api.v1.alpha3.GetSuggestionsReply.ParameterAssignments) | repeated |  |
| algorithm | [AlgorithmSpec](#api.v1.alpha3.AlgorithmSpec) |  |  |






<a name="api.v1.alpha3.GetSuggestionsReply.ParameterAssignments"></a>

### GetSuggestionsReply.ParameterAssignments



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| assignments | [ParameterAssignment](#api.v1.alpha3.ParameterAssignment) | repeated |  |






<a name="api.v1.alpha3.GetSuggestionsRequest"></a>

### GetSuggestionsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| experiment | [Experiment](#api.v1.alpha3.Experiment) |  |  |
| trials | [Trial](#api.v1.alpha3.Trial) | repeated | all completed trials owned by the experiment. |
| request_number | [int32](#int32) |  | The number of Suggestion you request at one time. When you set 3 to request_number, you can get three Suggestions at one time. |






<a name="api.v1.alpha3.GraphConfig"></a>

### GraphConfig
GraphConfig contains a config of DAG


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| num_layers | [int32](#int32) |  | Number of layers |
| input_sizes | [int32](#int32) | repeated | Dimensions of input size |
| output_sizes | [int32](#int32) | repeated | Dimensions of output size |






<a name="api.v1.alpha3.Metric"></a>

### Metric



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="api.v1.alpha3.MetricLog"></a>

### MetricLog



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| time_stamp | [string](#string) |  | RFC3339 format |
| metric | [Metric](#api.v1.alpha3.Metric) |  |  |






<a name="api.v1.alpha3.NasConfig"></a>

### NasConfig
NasConfig contains a config of NAS job


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| graph_config | [GraphConfig](#api.v1.alpha3.GraphConfig) |  | Config of DAG |
| operations | [NasConfig.Operations](#api.v1.alpha3.NasConfig.Operations) |  | List of Operation |






<a name="api.v1.alpha3.NasConfig.Operations"></a>

### NasConfig.Operations



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operation | [Operation](#api.v1.alpha3.Operation) | repeated |  |






<a name="api.v1.alpha3.ObjectiveSpec"></a>

### ObjectiveSpec



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [ObjectiveType](#api.v1.alpha3.ObjectiveType) |  |  |
| goal | [double](#double) |  |  |
| objective_metric_name | [string](#string) |  |  |
| additional_metric_names | [string](#string) | repeated | This can be empty if we only care about the objective metric. |






<a name="api.v1.alpha3.Observation"></a>

### Observation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| metrics | [Metric](#api.v1.alpha3.Metric) | repeated |  |






<a name="api.v1.alpha3.ObservationLog"></a>

### ObservationLog



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| metric_logs | [MetricLog](#api.v1.alpha3.MetricLog) | repeated |  |






<a name="api.v1.alpha3.Operation"></a>

### Operation
Config for operations in DAG


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operation_type | [string](#string) |  | Type of operation in DAG |
| parameter_specs | [Operation.ParameterSpecs](#api.v1.alpha3.Operation.ParameterSpecs) |  |  |






<a name="api.v1.alpha3.Operation.ParameterSpecs"></a>

### Operation.ParameterSpecs
List of ParameterSpec


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| parameters | [ParameterSpec](#api.v1.alpha3.ParameterSpec) | repeated |  |






<a name="api.v1.alpha3.ParameterAssignment"></a>

### ParameterAssignment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="api.v1.alpha3.ParameterSpec"></a>

### ParameterSpec
Config for a Hyper parameter.
Katib will create each Hyper parameter from this config.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the parameter. |
| parameter_type | [ParameterType](#api.v1.alpha3.ParameterType) |  | Type of the parameter. |
| feasible_space | [FeasibleSpace](#api.v1.alpha3.FeasibleSpace) |  | FeasibleSpace for the parameter. |






<a name="api.v1.alpha3.ReportObservationLogReply"></a>

### ReportObservationLogReply







<a name="api.v1.alpha3.ReportObservationLogRequest"></a>

### ReportObservationLogRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trial_name | [string](#string) |  |  |
| observation_log | [ObservationLog](#api.v1.alpha3.ObservationLog) |  |  |






<a name="api.v1.alpha3.Trial"></a>

### Trial



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| spec | [TrialSpec](#api.v1.alpha3.TrialSpec) |  |  |
| status | [TrialStatus](#api.v1.alpha3.TrialStatus) |  |  |






<a name="api.v1.alpha3.TrialSpec"></a>

### TrialSpec



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| experiment_name | [string](#string) |  |  |
| objective | [ObjectiveSpec](#api.v1.alpha3.ObjectiveSpec) |  |  |
| parameter_assignments | [TrialSpec.ParameterAssignments](#api.v1.alpha3.TrialSpec.ParameterAssignments) |  |  |
| run_spec | [string](#string) |  |  |
| metrics_collector_spec | [string](#string) |  |  |






<a name="api.v1.alpha3.TrialSpec.ParameterAssignments"></a>

### TrialSpec.ParameterAssignments
List of ParameterAssignment


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| assignments | [ParameterAssignment](#api.v1.alpha3.ParameterAssignment) | repeated |  |






<a name="api.v1.alpha3.TrialStatus"></a>

### TrialStatus



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| start_time | [string](#string) |  | RFC3339 format |
| completion_time | [string](#string) |  | RFC3339 format |
| condition | [TrialStatus.TrialConditionType](#api.v1.alpha3.TrialStatus.TrialConditionType) |  |  |
| observation | [Observation](#api.v1.alpha3.Observation) |  | The best observation in logs. |






<a name="api.v1.alpha3.ValidateAlgorithmSettingsReply"></a>

### ValidateAlgorithmSettingsReply
Return INVALID_ARGUMENT Error if Algorithm Settings are not Valid






<a name="api.v1.alpha3.ValidateAlgorithmSettingsRequest"></a>

### ValidateAlgorithmSettingsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| experiment | [Experiment](#api.v1.alpha3.Experiment) |  |  |





 


<a name="api.v1.alpha3.ObjectiveType"></a>

### ObjectiveType
Direction of optimization. Minimize or Maximize.

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 | Undefined type and not used. |
| MINIMIZE | 1 | Minimize |
| MAXIMIZE | 2 | Maximize |



<a name="api.v1.alpha3.ParameterType"></a>

### ParameterType
Types of value for HyperParameter.

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN_TYPE | 0 | Undefined type and not used. |
| DOUBLE | 1 | Double float type. Use &#34;Max/Min&#34;. |
| INT | 2 | Int type. Use &#34;Max/Min&#34;. |
| DISCRETE | 3 | Discrete number type. Use &#34;List&#34; as float. |
| CATEGORICAL | 4 | Categorical type. Use &#34;List&#34; as string. |



<a name="api.v1.alpha3.TrialStatus.TrialConditionType"></a>

### TrialStatus.TrialConditionType


| Name | Number | Description |
| ---- | ------ | ----------- |
| CREATED | 0 |  |
| RUNNING | 1 |  |
| SUCCEEDED | 2 |  |
| KILLED | 3 |  |
| FAILED | 4 |  |
| UNKNOWN | 5 |  |


 

 


<a name="api.v1.alpha3.EarlyStopping"></a>

### EarlyStopping
TODO

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|


<a name="api.v1.alpha3.Manager"></a>

### Manager
Service for Main API for Katib
For each RPC service, we define mapping to HTTP REST API method.
The mapping includes the URL path, query parameters and request body.
https://cloud.google.com/service-infrastructure/docs/service-management/reference/rpc/google.api#http

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ReportObservationLog | [ReportObservationLogRequest](#api.v1.alpha3.ReportObservationLogRequest) | [ReportObservationLogReply](#api.v1.alpha3.ReportObservationLogReply) | Report a log of Observations for a Trial. The log consists of timestamp and value of metric. Katib store every log of metrics. You can see accuracy curve or other metric logs on UI. |
| GetObservationLog | [GetObservationLogRequest](#api.v1.alpha3.GetObservationLogRequest) | [GetObservationLogReply](#api.v1.alpha3.GetObservationLogReply) | Get all log of Observations for a Trial. |
| DeleteObservationLog | [DeleteObservationLogRequest](#api.v1.alpha3.DeleteObservationLogRequest) | [DeleteObservationLogReply](#api.v1.alpha3.DeleteObservationLogReply) | Delete all log of Observations for a Trial. |


<a name="api.v1.alpha3.Suggestion"></a>

### Suggestion


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetSuggestions | [GetSuggestionsRequest](#api.v1.alpha3.GetSuggestionsRequest) | [GetSuggestionsReply](#api.v1.alpha3.GetSuggestionsReply) |  |
| ValidateAlgorithmSettings | [ValidateAlgorithmSettingsRequest](#api.v1.alpha3.ValidateAlgorithmSettingsRequest) | [ValidateAlgorithmSettingsReply](#api.v1.alpha3.ValidateAlgorithmSettingsReply) |  |

 



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

