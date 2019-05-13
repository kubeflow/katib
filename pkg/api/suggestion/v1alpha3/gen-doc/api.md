# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [suggestion.proto](#suggestion.proto)
    - [AlgorithmSetting](#api.v1.alpha3.AlgorithmSetting)
    - [AlgorithmSpec](#api.v1.alpha3.AlgorithmSpec)
    - [Experiment](#api.v1.alpha3.Experiment)
    - [ExperimentSpec](#api.v1.alpha3.ExperimentSpec)
    - [FeasibleSpace](#api.v1.alpha3.FeasibleSpace)
    - [GetSuggestionsReply](#api.v1.alpha3.GetSuggestionsReply)
    - [GetSuggestionsRequest](#api.v1.alpha3.GetSuggestionsRequest)
    - [Metric](#api.v1.alpha3.Metric)
    - [ObjectiveSpec](#api.v1.alpha3.ObjectiveSpec)
    - [Observation](#api.v1.alpha3.Observation)
    - [ParameterAssignment](#api.v1.alpha3.ParameterAssignment)
    - [ParameterAssignments](#api.v1.alpha3.ParameterAssignments)
    - [ParameterSpec](#api.v1.alpha3.ParameterSpec)
    - [ParameterSpecs](#api.v1.alpha3.ParameterSpecs)
    - [Trial](#api.v1.alpha3.Trial)
    - [TrialSpec](#api.v1.alpha3.TrialSpec)
    - [TrialStatus](#api.v1.alpha3.TrialStatus)
  
    - [ObjectiveType](#api.v1.alpha3.ObjectiveType)
    - [ParameterType](#api.v1.alpha3.ParameterType)
  
  
    - [Suggestion](#api.v1.alpha3.Suggestion)
  

- [Scalar Value Types](#scalar-value-types)



<a name="suggestion.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## suggestion.proto



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






<a name="api.v1.alpha3.Experiment"></a>

### Experiment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| experiment_spec | [ExperimentSpec](#api.v1.alpha3.ExperimentSpec) |  |  |






<a name="api.v1.alpha3.ExperimentSpec"></a>

### ExperimentSpec



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| algorithm | [AlgorithmSpec](#api.v1.alpha3.AlgorithmSpec) |  |  |
| parameter_specs | [ParameterSpecs](#api.v1.alpha3.ParameterSpecs) |  |  |
| objective | [ObjectiveSpec](#api.v1.alpha3.ObjectiveSpec) |  |  |






<a name="api.v1.alpha3.FeasibleSpace"></a>

### FeasibleSpace



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| max | [string](#string) |  | Max Value |
| min | [string](#string) |  | Minimum Value |
| list | [string](#string) | repeated | List of Values. |
| step | [string](#string) |  | Step for double or int parameter |






<a name="api.v1.alpha3.GetSuggestionsReply"></a>

### GetSuggestionsReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trials | [Trial](#api.v1.alpha3.Trial) | repeated | trials should be created in the next run. |






<a name="api.v1.alpha3.GetSuggestionsRequest"></a>

### GetSuggestionsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| experiment | [Experiment](#api.v1.alpha3.Experiment) |  |  |
| trials | [Trial](#api.v1.alpha3.Trial) | repeated | all completed trials owned by the experiment. |
| request_number | [int32](#int32) |  |  |






<a name="api.v1.alpha3.Metric"></a>

### Metric



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="api.v1.alpha3.ObjectiveSpec"></a>

### ObjectiveSpec



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [ObjectiveType](#api.v1.alpha3.ObjectiveType) |  |  |
| goal | [double](#double) |  |  |
| objective_metric_name | [string](#string) |  |  |






<a name="api.v1.alpha3.Observation"></a>

### Observation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| metrics | [Metric](#api.v1.alpha3.Metric) | repeated |  |






<a name="api.v1.alpha3.ParameterAssignment"></a>

### ParameterAssignment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="api.v1.alpha3.ParameterAssignments"></a>

### ParameterAssignments



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| assignments | [ParameterAssignment](#api.v1.alpha3.ParameterAssignment) | repeated |  |






<a name="api.v1.alpha3.ParameterSpec"></a>

### ParameterSpec



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the parameter. |
| parameter_type | [ParameterType](#api.v1.alpha3.ParameterType) |  | Type of the parameter. |
| feasible_space | [FeasibleSpace](#api.v1.alpha3.FeasibleSpace) |  | FeasibleSpace for the parameter. |






<a name="api.v1.alpha3.ParameterSpecs"></a>

### ParameterSpecs



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| parameters | [ParameterSpec](#api.v1.alpha3.ParameterSpec) | repeated |  |






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
| parameter_assignments | [ParameterAssignments](#api.v1.alpha3.ParameterAssignments) |  |  |
| run_spec | [string](#string) |  |  |






<a name="api.v1.alpha3.TrialStatus"></a>

### TrialStatus



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| observation | [Observation](#api.v1.alpha3.Observation) |  | The best observation in logs. |





 


<a name="api.v1.alpha3.ObjectiveType"></a>

### ObjectiveType


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 | Undefined type and not used. |
| MINIMIZE | 1 | Minimize |
| MAXIMIZE | 2 | Maximize |



<a name="api.v1.alpha3.ParameterType"></a>

### ParameterType


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN_TYPE | 0 | Undefined type and not used. |
| DOUBLE | 1 | Double float type. Use &#34;Max/Min&#34;. |
| INT | 2 | Int type. Use &#34;Max/Min&#34;. |
| DISCRETE | 3 | Discrete number type. Use &#34;List&#34; as float. |
| CATEGORICAL | 4 | Categorical type. Use &#34;List&#34; as string. |


 

 


<a name="api.v1.alpha3.Suggestion"></a>

### Suggestion


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetSuggestions | [GetSuggestionsRequest](#api.v1.alpha3.GetSuggestionsRequest) | [GetSuggestionsReply](#api.v1.alpha3.GetSuggestionsReply) |  |

 



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

