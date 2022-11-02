# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api.proto](#api-proto)
    - [AlgorithmSetting](#api-v1-beta1-AlgorithmSetting)
    - [AlgorithmSpec](#api-v1-beta1-AlgorithmSpec)
    - [DeleteObservationLogReply](#api-v1-beta1-DeleteObservationLogReply)
    - [DeleteObservationLogRequest](#api-v1-beta1-DeleteObservationLogRequest)
    - [EarlyStoppingRule](#api-v1-beta1-EarlyStoppingRule)
    - [EarlyStoppingSetting](#api-v1-beta1-EarlyStoppingSetting)
    - [EarlyStoppingSpec](#api-v1-beta1-EarlyStoppingSpec)
    - [Experiment](#api-v1-beta1-Experiment)
    - [ExperimentSpec](#api-v1-beta1-ExperimentSpec)
    - [ExperimentSpec.ParameterSpecs](#api-v1-beta1-ExperimentSpec-ParameterSpecs)
    - [FeasibleSpace](#api-v1-beta1-FeasibleSpace)
    - [GetEarlyStoppingRulesReply](#api-v1-beta1-GetEarlyStoppingRulesReply)
    - [GetEarlyStoppingRulesRequest](#api-v1-beta1-GetEarlyStoppingRulesRequest)
    - [GetObservationLogReply](#api-v1-beta1-GetObservationLogReply)
    - [GetObservationLogRequest](#api-v1-beta1-GetObservationLogRequest)
    - [GetSuggestionsReply](#api-v1-beta1-GetSuggestionsReply)
    - [GetSuggestionsReply.ParameterAssignments](#api-v1-beta1-GetSuggestionsReply-ParameterAssignments)
    - [GetSuggestionsReply.ParameterAssignments.LabelsEntry](#api-v1-beta1-GetSuggestionsReply-ParameterAssignments-LabelsEntry)
    - [GetSuggestionsRequest](#api-v1-beta1-GetSuggestionsRequest)
    - [GraphConfig](#api-v1-beta1-GraphConfig)
    - [Metric](#api-v1-beta1-Metric)
    - [MetricLog](#api-v1-beta1-MetricLog)
    - [NasConfig](#api-v1-beta1-NasConfig)
    - [NasConfig.Operations](#api-v1-beta1-NasConfig-Operations)
    - [ObjectiveSpec](#api-v1-beta1-ObjectiveSpec)
    - [Observation](#api-v1-beta1-Observation)
    - [ObservationLog](#api-v1-beta1-ObservationLog)
    - [Operation](#api-v1-beta1-Operation)
    - [Operation.ParameterSpecs](#api-v1-beta1-Operation-ParameterSpecs)
    - [ParameterAssignment](#api-v1-beta1-ParameterAssignment)
    - [ParameterSpec](#api-v1-beta1-ParameterSpec)
    - [ReportObservationLogReply](#api-v1-beta1-ReportObservationLogReply)
    - [ReportObservationLogRequest](#api-v1-beta1-ReportObservationLogRequest)
    - [SetTrialStatusReply](#api-v1-beta1-SetTrialStatusReply)
    - [SetTrialStatusRequest](#api-v1-beta1-SetTrialStatusRequest)
    - [Trial](#api-v1-beta1-Trial)
    - [TrialSpec](#api-v1-beta1-TrialSpec)
    - [TrialSpec.LabelsEntry](#api-v1-beta1-TrialSpec-LabelsEntry)
    - [TrialSpec.ParameterAssignments](#api-v1-beta1-TrialSpec-ParameterAssignments)
    - [TrialStatus](#api-v1-beta1-TrialStatus)
    - [ValidateAlgorithmSettingsReply](#api-v1-beta1-ValidateAlgorithmSettingsReply)
    - [ValidateAlgorithmSettingsRequest](#api-v1-beta1-ValidateAlgorithmSettingsRequest)
    - [ValidateEarlyStoppingSettingsReply](#api-v1-beta1-ValidateEarlyStoppingSettingsReply)
    - [ValidateEarlyStoppingSettingsRequest](#api-v1-beta1-ValidateEarlyStoppingSettingsRequest)
  
    - [ComparisonType](#api-v1-beta1-ComparisonType)
    - [ObjectiveType](#api-v1-beta1-ObjectiveType)
    - [ParameterType](#api-v1-beta1-ParameterType)
    - [TrialStatus.TrialConditionType](#api-v1-beta1-TrialStatus-TrialConditionType)
  
    - [DBManager](#api-v1-beta1-DBManager)
    - [EarlyStopping](#api-v1-beta1-EarlyStopping)
    - [Suggestion](#api-v1-beta1-Suggestion)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api.proto
Katib GRPC API v1beta1


<a name="api-v1-beta1-AlgorithmSetting"></a>

### AlgorithmSetting
HP or NAS algorithm settings.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="api-v1-beta1-AlgorithmSpec"></a>

### AlgorithmSpec
HP or NAS algorithm specification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| algorithm_name | [string](#string) |  |  |
| algorithm_settings | [AlgorithmSetting](#api-v1-beta1-AlgorithmSetting) | repeated |  |






<a name="api-v1-beta1-DeleteObservationLogReply"></a>

### DeleteObservationLogReply







<a name="api-v1-beta1-DeleteObservationLogRequest"></a>

### DeleteObservationLogRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trial_name | [string](#string) |  |  |






<a name="api-v1-beta1-EarlyStoppingRule"></a>

### EarlyStoppingRule
EarlyStoppingRule represents single early stopping rule.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the rule. Usually, metric name. |
| value | [string](#string) |  | Value of the metric. |
| comparison | [ComparisonType](#api-v1-beta1-ComparisonType) |  | Correlation between name and value, one of equal, less or greater |
| start_step | [int32](#int32) |  | Defines quantity of intermediate results that should be received before applying the rule. If start step is empty, rule is applied from the first recorded metric. |






<a name="api-v1-beta1-EarlyStoppingSetting"></a>

### EarlyStoppingSetting
Early stopping algorithm settings.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="api-v1-beta1-EarlyStoppingSpec"></a>

### EarlyStoppingSpec
Early stopping algorithm specification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| algorithm_name | [string](#string) |  |  |
| algorithm_settings | [EarlyStoppingSetting](#api-v1-beta1-EarlyStoppingSetting) | repeated |  |






<a name="api-v1-beta1-Experiment"></a>

### Experiment
Structure for a single Experiment.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name for the Experiment. |
| spec | [ExperimentSpec](#api-v1-beta1-ExperimentSpec) |  | Experiment specification. |






<a name="api-v1-beta1-ExperimentSpec"></a>

### ExperimentSpec
Specification of an Experiment. Experiment represents a single optimization run over a feasible space. 
Each Experiment contains a configuration describing the feasible space, as well as a set of Trials.
It is assumed that objective function f(x) does not change in the course of an Experiment.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| parameter_specs | [ExperimentSpec.ParameterSpecs](#api-v1-beta1-ExperimentSpec-ParameterSpecs) |  |  |
| objective | [ObjectiveSpec](#api-v1-beta1-ObjectiveSpec) |  | Objective specification for the Experiment. |
| algorithm | [AlgorithmSpec](#api-v1-beta1-AlgorithmSpec) |  | HP or NAS algorithm specification for the Experiment. |
| early_stopping | [EarlyStoppingSpec](#api-v1-beta1-EarlyStoppingSpec) |  | Early stopping specification for the Experiment. |
| parallel_trial_count | [int32](#int32) |  | How many Trials can be processed in parallel. |
| max_trial_count | [int32](#int32) |  | Max completed Trials to mark Experiment as succeeded. |
| nas_config | [NasConfig](#api-v1-beta1-NasConfig) |  | NAS configuration for the Experiment. |






<a name="api-v1-beta1-ExperimentSpec-ParameterSpecs"></a>

### ExperimentSpec.ParameterSpecs
List of ParameterSpec.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| parameters | [ParameterSpec](#api-v1-beta1-ParameterSpec) | repeated |  |






<a name="api-v1-beta1-FeasibleSpace"></a>

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






<a name="api-v1-beta1-GetEarlyStoppingRulesReply"></a>

### GetEarlyStoppingRulesReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| early_stopping_rules | [EarlyStoppingRule](#api-v1-beta1-EarlyStoppingRule) | repeated |  |






<a name="api-v1-beta1-GetEarlyStoppingRulesRequest"></a>

### GetEarlyStoppingRulesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| experiment | [Experiment](#api-v1-beta1-Experiment) |  |  |
| trials | [Trial](#api-v1-beta1-Trial) | repeated |  |
| db_manager_address | [string](#string) |  |  |






<a name="api-v1-beta1-GetObservationLogReply"></a>

### GetObservationLogReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| observation_log | [ObservationLog](#api-v1-beta1-ObservationLog) |  |  |






<a name="api-v1-beta1-GetObservationLogRequest"></a>

### GetObservationLogRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trial_name | [string](#string) |  |  |
| metric_name | [string](#string) |  |  |
| start_time | [string](#string) |  | The start of the time range. RFC3339 format |
| end_time | [string](#string) |  | The end of the time range. RFC3339 format |






<a name="api-v1-beta1-GetSuggestionsReply"></a>

### GetSuggestionsReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| parameter_assignments | [GetSuggestionsReply.ParameterAssignments](#api-v1-beta1-GetSuggestionsReply-ParameterAssignments) | repeated |  |
| algorithm | [AlgorithmSpec](#api-v1-beta1-AlgorithmSpec) |  |  |
| early_stopping_rules | [EarlyStoppingRule](#api-v1-beta1-EarlyStoppingRule) | repeated |  |






<a name="api-v1-beta1-GetSuggestionsReply-ParameterAssignments"></a>

### GetSuggestionsReply.ParameterAssignments



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| assignments | [ParameterAssignment](#api-v1-beta1-ParameterAssignment) | repeated |  |
| trial_name | [string](#string) |  | Optional field to override the trial name |
| labels | [GetSuggestionsReply.ParameterAssignments.LabelsEntry](#api-v1-beta1-GetSuggestionsReply-ParameterAssignments-LabelsEntry) | repeated | Optional field to add labels to the generated Trials |






<a name="api-v1-beta1-GetSuggestionsReply-ParameterAssignments-LabelsEntry"></a>

### GetSuggestionsReply.ParameterAssignments.LabelsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="api-v1-beta1-GetSuggestionsRequest"></a>

### GetSuggestionsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| experiment | [Experiment](#api-v1-beta1-Experiment) |  |  |
| trials | [Trial](#api-v1-beta1-Trial) | repeated | All completed trials owned by the experiment. |
| current_request_number | [int32](#int32) |  | The number of Suggestions requested at one time. When you set 3 to current_request_number, you get three Suggestions at one time. |
| total_request_number | [int32](#int32) |  | The number of Suggestions requested till now. |






<a name="api-v1-beta1-GraphConfig"></a>

### GraphConfig
GraphConfig contains a config of DAG


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| num_layers | [int32](#int32) |  | Number of layers |
| input_sizes | [int32](#int32) | repeated | Dimensions of input size |
| output_sizes | [int32](#int32) | repeated | Dimensions of output size |






<a name="api-v1-beta1-Metric"></a>

### Metric



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="api-v1-beta1-MetricLog"></a>

### MetricLog



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| time_stamp | [string](#string) |  | RFC3339 format |
| metric | [Metric](#api-v1-beta1-Metric) |  |  |






<a name="api-v1-beta1-NasConfig"></a>

### NasConfig
NasConfig contains a config of NAS job


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| graph_config | [GraphConfig](#api-v1-beta1-GraphConfig) |  | Config of DAG |
| operations | [NasConfig.Operations](#api-v1-beta1-NasConfig-Operations) |  | List of Operation |






<a name="api-v1-beta1-NasConfig-Operations"></a>

### NasConfig.Operations



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operation | [Operation](#api-v1-beta1-Operation) | repeated |  |






<a name="api-v1-beta1-ObjectiveSpec"></a>

### ObjectiveSpec
Objective specification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [ObjectiveType](#api-v1-beta1-ObjectiveType) |  | Type of optimization. |
| goal | [double](#double) |  | Goal of optimization, can be empty. |
| objective_metric_name | [string](#string) |  | Primary metric name for the optimization. |
| additional_metric_names | [string](#string) | repeated | List of additional metrics to record from Trial. This can be empty if we only care about the objective metric. |






<a name="api-v1-beta1-Observation"></a>

### Observation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| metrics | [Metric](#api-v1-beta1-Metric) | repeated |  |






<a name="api-v1-beta1-ObservationLog"></a>

### ObservationLog



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| metric_logs | [MetricLog](#api-v1-beta1-MetricLog) | repeated |  |






<a name="api-v1-beta1-Operation"></a>

### Operation
Config for operations in DAG


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operation_type | [string](#string) |  | Type of operation in DAG |
| parameter_specs | [Operation.ParameterSpecs](#api-v1-beta1-Operation-ParameterSpecs) |  |  |






<a name="api-v1-beta1-Operation-ParameterSpecs"></a>

### Operation.ParameterSpecs
List of ParameterSpec


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| parameters | [ParameterSpec](#api-v1-beta1-ParameterSpec) | repeated |  |






<a name="api-v1-beta1-ParameterAssignment"></a>

### ParameterAssignment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="api-v1-beta1-ParameterSpec"></a>

### ParameterSpec
Config for a hyperparameter.
Katib will create each Hyper parameter from this config.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the parameter. |
| parameter_type | [ParameterType](#api-v1-beta1-ParameterType) |  | Type of the parameter. |
| feasible_space | [FeasibleSpace](#api-v1-beta1-FeasibleSpace) |  | FeasibleSpace for the parameter. |






<a name="api-v1-beta1-ReportObservationLogReply"></a>

### ReportObservationLogReply







<a name="api-v1-beta1-ReportObservationLogRequest"></a>

### ReportObservationLogRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trial_name | [string](#string) |  |  |
| observation_log | [ObservationLog](#api-v1-beta1-ObservationLog) |  |  |






<a name="api-v1-beta1-SetTrialStatusReply"></a>

### SetTrialStatusReply







<a name="api-v1-beta1-SetTrialStatusRequest"></a>

### SetTrialStatusRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trial_name | [string](#string) |  |  |






<a name="api-v1-beta1-Trial"></a>

### Trial
Structure for a single Trial.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name for the Trial. |
| spec | [TrialSpec](#api-v1-beta1-TrialSpec) |  | Trial specification. |
| status | [TrialStatus](#api-v1-beta1-TrialStatus) |  | Trial status. |






<a name="api-v1-beta1-TrialSpec"></a>

### TrialSpec
Specification of a Trial. It represents Trial&#39;s parameter assignments and objective.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| objective | [ObjectiveSpec](#api-v1-beta1-ObjectiveSpec) |  | Objective specification for the Trial. |
| parameter_assignments | [TrialSpec.ParameterAssignments](#api-v1-beta1-TrialSpec-ParameterAssignments) |  | List of assignments generated for the Trial. |
| labels | [TrialSpec.LabelsEntry](#api-v1-beta1-TrialSpec-LabelsEntry) | repeated | Map of labels assigned to the Trial |






<a name="api-v1-beta1-TrialSpec-LabelsEntry"></a>

### TrialSpec.LabelsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="api-v1-beta1-TrialSpec-ParameterAssignments"></a>

### TrialSpec.ParameterAssignments
List of ParameterAssignment


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| assignments | [ParameterAssignment](#api-v1-beta1-ParameterAssignment) | repeated |  |






<a name="api-v1-beta1-TrialStatus"></a>

### TrialStatus
Current Trial status. It contains Trial&#39;s latest condition, start time, completion time, observation.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| start_time | [string](#string) |  | Trial start time in RFC3339 format |
| completion_time | [string](#string) |  | Trial completion time in RFC3339 format |
| condition | [TrialStatus.TrialConditionType](#api-v1-beta1-TrialStatus-TrialConditionType) |  | Trial current condition. It is equal to the latest Trial CR condition. |
| observation | [Observation](#api-v1-beta1-Observation) |  | The best Trial observation in logs. |






<a name="api-v1-beta1-ValidateAlgorithmSettingsReply"></a>

### ValidateAlgorithmSettingsReply
Return INVALID_ARGUMENT Error if Algorithm Settings are not Valid






<a name="api-v1-beta1-ValidateAlgorithmSettingsRequest"></a>

### ValidateAlgorithmSettingsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| experiment | [Experiment](#api-v1-beta1-Experiment) |  |  |






<a name="api-v1-beta1-ValidateEarlyStoppingSettingsReply"></a>

### ValidateEarlyStoppingSettingsReply
Return INVALID_ARGUMENT Error if Early Stopping Settings are not Valid






<a name="api-v1-beta1-ValidateEarlyStoppingSettingsRequest"></a>

### ValidateEarlyStoppingSettingsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| early_stopping | [EarlyStoppingSpec](#api-v1-beta1-EarlyStoppingSpec) |  |  |





 


<a name="api-v1-beta1-ComparisonType"></a>

### ComparisonType


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN_COMPARISON | 0 | Unknown comparison, not used |
| EQUAL | 1 | Equal comparison, e.g. accuracy = 0.7 |
| LESS | 2 | Less comparison, e.g. accuracy &lt; 0.7 |
| GREATER | 3 | Greater comparison, e.g. accuracy &gt; 0.7 |



<a name="api-v1-beta1-ObjectiveType"></a>

### ObjectiveType
Direction of optimization. Minimize or Maximize.

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 | Undefined type and not used. |
| MINIMIZE | 1 | Minimize |
| MAXIMIZE | 2 | Maximize |



<a name="api-v1-beta1-ParameterType"></a>

### ParameterType
Types of value for HyperParameter.

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN_TYPE | 0 | Undefined type and not used. |
| DOUBLE | 1 | Double float type. Use &#34;Max/Min&#34;. |
| INT | 2 | Int type. Use &#34;Max/Min&#34;. |
| DISCRETE | 3 | Discrete number type. Use &#34;List&#34; as float. |
| CATEGORICAL | 4 | Categorical type. Use &#34;List&#34; as string. |



<a name="api-v1-beta1-TrialStatus-TrialConditionType"></a>

### TrialStatus.TrialConditionType
Trial can be in one of 8 conditions.
TODO (andreyvelich): Remove unused conditions.

| Name | Number | Description |
| ---- | ------ | ----------- |
| CREATED | 0 |  |
| RUNNING | 1 |  |
| SUCCEEDED | 2 |  |
| KILLED | 3 |  |
| FAILED | 4 |  |
| METRICSUNAVAILABLE | 5 |  |
| EARLYSTOPPED | 6 |  |
| UNKNOWN | 7 |  |


 

 


<a name="api-v1-beta1-DBManager"></a>

### DBManager
DBManager service defines APIs to manage Katib database.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ReportObservationLog | [ReportObservationLogRequest](#api-v1-beta1-ReportObservationLogRequest) | [ReportObservationLogReply](#api-v1-beta1-ReportObservationLogReply) | Report a log of Observations for a Trial. The log consists of timestamp and value of metric. Katib store every log of metrics. You can see accuracy curve or other metric logs on UI. |
| GetObservationLog | [GetObservationLogRequest](#api-v1-beta1-GetObservationLogRequest) | [GetObservationLogReply](#api-v1-beta1-GetObservationLogReply) | Get all log of Observations for a Trial. |
| DeleteObservationLog | [DeleteObservationLogRequest](#api-v1-beta1-DeleteObservationLogRequest) | [DeleteObservationLogReply](#api-v1-beta1-DeleteObservationLogReply) | Delete all log of Observations for a Trial. |


<a name="api-v1-beta1-EarlyStopping"></a>

### EarlyStopping
EarlyStopping service defines APIs to manage Katib Early Stopping algorithms

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetEarlyStoppingRules | [GetEarlyStoppingRulesRequest](#api-v1-beta1-GetEarlyStoppingRulesRequest) | [GetEarlyStoppingRulesReply](#api-v1-beta1-GetEarlyStoppingRulesReply) |  |
| SetTrialStatus | [SetTrialStatusRequest](#api-v1-beta1-SetTrialStatusRequest) | [SetTrialStatusReply](#api-v1-beta1-SetTrialStatusReply) |  |
| ValidateEarlyStoppingSettings | [ValidateEarlyStoppingSettingsRequest](#api-v1-beta1-ValidateEarlyStoppingSettingsRequest) | [ValidateEarlyStoppingSettingsReply](#api-v1-beta1-ValidateEarlyStoppingSettingsReply) |  |


<a name="api-v1-beta1-Suggestion"></a>

### Suggestion
Suggestion service defines APIs to manage Katib Suggestion from HP or NAS algorithms

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetSuggestions | [GetSuggestionsRequest](#api-v1-beta1-GetSuggestionsRequest) | [GetSuggestionsReply](#api-v1-beta1-GetSuggestionsReply) |  |
| ValidateAlgorithmSettings | [ValidateAlgorithmSettingsRequest](#api-v1-beta1-ValidateAlgorithmSettingsRequest) | [ValidateAlgorithmSettingsReply](#api-v1-beta1-ValidateAlgorithmSettingsReply) |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

