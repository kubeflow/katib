# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api.proto](#api.proto)
    - [GetSuggestionsReply](#api.v1.alpha2.GetSuggestionsReply)
    - [GetSuggestionsRequest](#api.v1.alpha2.GetSuggestionsRequest)
    - [ValidateAlgorithmSettingsReply](#api.v1.alpha2.ValidateAlgorithmSettingsReply)
    - [ValidateAlgorithmSettingsRequest](#api.v1.alpha2.ValidateAlgorithmSettingsRequest)
  
  
  
    - [EarlyStopping](#api.v1.alpha2.EarlyStopping)
    - [Manager](#api.v1.alpha2.Manager)
    - [Suggestion](#api.v1.alpha2.Suggestion)
  

- [Scalar Value Types](#scalar-value-types)



<a name="api.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api.proto
Katib API


<a name="api.v1.alpha2.GetSuggestionsReply"></a>

### GetSuggestionsReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trials | [dbif.v1.alpha2.Trial](#dbif.v1.alpha2.Trial) | repeated |  |






<a name="api.v1.alpha2.GetSuggestionsRequest"></a>

### GetSuggestionsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| experiment_name | [string](#string) |  |  |
| algorithm_name | [string](#string) |  |  |
| request_number | [int32](#int32) |  | The number of Suggestion you request at one time. When you set 3 to request_number, you can get three Suggestions at one time. |






<a name="api.v1.alpha2.ValidateAlgorithmSettingsReply"></a>

### ValidateAlgorithmSettingsReply
Return INVALID_ARGUMENT Error if Algorithm Settings are not Valid






<a name="api.v1.alpha2.ValidateAlgorithmSettingsRequest"></a>

### ValidateAlgorithmSettingsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| experiment_spec | [dbif.v1.alpha2.ExperimentSpec](#dbif.v1.alpha2.ExperimentSpec) |  |  |
| algorithm_name | [string](#string) |  |  |





 

 

 


<a name="api.v1.alpha2.EarlyStopping"></a>

### EarlyStopping
TODO

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|


<a name="api.v1.alpha2.Manager"></a>

### Manager
Service for Main API for Katib
For each RPC service, we define mapping to HTTP REST API method.
The mapping includes the URL path, query parameters and request body.
https://cloud.google.com/service-infrastructure/docs/service-management/reference/rpc/google.api#http

*
Register a Experiment to DB.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetSuggestions | [GetSuggestionsRequest](#api.v1.alpha2.GetSuggestionsRequest) | [GetSuggestionsReply](#api.v1.alpha2.GetSuggestionsReply) | Get Suggestions from a Suggestion service. |
| ValidateAlgorithmSettings | [ValidateAlgorithmSettingsRequest](#api.v1.alpha2.ValidateAlgorithmSettingsRequest) | [ValidateAlgorithmSettingsReply](#api.v1.alpha2.ValidateAlgorithmSettingsReply) | Validate AlgorithmSettings in an Experiment. Suggestion service should return INVALID_ARGUMENT Error when the parameter is invalid |


<a name="api.v1.alpha2.Suggestion"></a>

### Suggestion


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetSuggestions | [GetSuggestionsRequest](#api.v1.alpha2.GetSuggestionsRequest) | [GetSuggestionsReply](#api.v1.alpha2.GetSuggestionsReply) |  |
| ValidateAlgorithmSettings | [ValidateAlgorithmSettingsRequest](#api.v1.alpha2.ValidateAlgorithmSettingsRequest) | [ValidateAlgorithmSettingsReply](#api.v1.alpha2.ValidateAlgorithmSettingsReply) |  |

 



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

