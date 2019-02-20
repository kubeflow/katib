# Protocol Documentation
<a name="top"/>

## Table of Contents

- [api.proto](#api.proto)
    - [CreateStudyReply](#api.CreateStudyReply)
    - [CreateStudyRequest](#api.CreateStudyRequest)
    - [CreateTrialReply](#api.CreateTrialReply)
    - [CreateTrialRequest](#api.CreateTrialRequest)
    - [DataSetInfo](#api.DataSetInfo)
    - [DeleteStudyReply](#api.DeleteStudyReply)
    - [DeleteStudyRequest](#api.DeleteStudyRequest)
    - [EarlyStoppingParameter](#api.EarlyStoppingParameter)
    - [EarlyStoppingParameterSet](#api.EarlyStoppingParameterSet)
    - [FeasibleSpace](#api.FeasibleSpace)
    - [GetEarlyStoppingParameterListReply](#api.GetEarlyStoppingParameterListReply)
    - [GetEarlyStoppingParameterListRequest](#api.GetEarlyStoppingParameterListRequest)
    - [GetEarlyStoppingParametersReply](#api.GetEarlyStoppingParametersReply)
    - [GetEarlyStoppingParametersRequest](#api.GetEarlyStoppingParametersRequest)
    - [GetMetricsReply](#api.GetMetricsReply)
    - [GetMetricsRequest](#api.GetMetricsRequest)
    - [GetSavedModelReply](#api.GetSavedModelReply)
    - [GetSavedModelRequest](#api.GetSavedModelRequest)
    - [GetSavedModelsReply](#api.GetSavedModelsReply)
    - [GetSavedModelsRequest](#api.GetSavedModelsRequest)
    - [GetSavedStudiesReply](#api.GetSavedStudiesReply)
    - [GetSavedStudiesRequest](#api.GetSavedStudiesRequest)
    - [GetShouldStopWorkersReply](#api.GetShouldStopWorkersReply)
    - [GetShouldStopWorkersRequest](#api.GetShouldStopWorkersRequest)
    - [GetStudyListReply](#api.GetStudyListReply)
    - [GetStudyListRequest](#api.GetStudyListRequest)
    - [GetStudyReply](#api.GetStudyReply)
    - [GetStudyRequest](#api.GetStudyRequest)
    - [GetSuggestionParameterListReply](#api.GetSuggestionParameterListReply)
    - [GetSuggestionParameterListRequest](#api.GetSuggestionParameterListRequest)
    - [GetSuggestionParametersReply](#api.GetSuggestionParametersReply)
    - [GetSuggestionParametersRequest](#api.GetSuggestionParametersRequest)
    - [GetSuggestionsReply](#api.GetSuggestionsReply)
    - [GetSuggestionsRequest](#api.GetSuggestionsRequest)
    - [GetTrialReply](#api.GetTrialReply)
    - [GetTrialRequest](#api.GetTrialRequest)
    - [GetTrialsReply](#api.GetTrialsReply)
    - [GetTrialsRequest](#api.GetTrialsRequest)
    - [GetWorkerFullInfoReply](#api.GetWorkerFullInfoReply)
    - [GetWorkerFullInfoRequest](#api.GetWorkerFullInfoRequest)
    - [GetWorkersReply](#api.GetWorkersReply)
    - [GetWorkersRequest](#api.GetWorkersRequest)
    - [Metrics](#api.Metrics)
    - [MetricsLog](#api.MetricsLog)
    - [MetricsLogSet](#api.MetricsLogSet)
    - [MetricsValueTime](#api.MetricsValueTime)
    - [ModelInfo](#api.ModelInfo)
    - [Parameter](#api.Parameter)
    - [ParameterConfig](#api.ParameterConfig)
    - [RegisterWorkerReply](#api.RegisterWorkerReply)
    - [RegisterWorkerRequest](#api.RegisterWorkerRequest)
    - [ReportMetricsLogsReply](#api.ReportMetricsLogsReply)
    - [ReportMetricsLogsRequest](#api.ReportMetricsLogsRequest)
    - [SaveModelReply](#api.SaveModelReply)
    - [SaveModelRequest](#api.SaveModelRequest)
    - [SaveStudyReply](#api.SaveStudyReply)
    - [SaveStudyRequest](#api.SaveStudyRequest)
    - [SetEarlyStoppingParametersReply](#api.SetEarlyStoppingParametersReply)
    - [SetEarlyStoppingParametersRequest](#api.SetEarlyStoppingParametersRequest)
    - [SetSuggestionParametersReply](#api.SetSuggestionParametersReply)
    - [SetSuggestionParametersRequest](#api.SetSuggestionParametersRequest)
    - [StopSuggestionReply](#api.StopSuggestionReply)
    - [StopSuggestionRequest](#api.StopSuggestionRequest)
    - [StopWorkersReply](#api.StopWorkersReply)
    - [StopWorkersRequest](#api.StopWorkersRequest)
    - [StudyConfig](#api.StudyConfig)
    - [StudyConfig.ParameterConfigs](#api.StudyConfig.ParameterConfigs)
    - [StudyOverview](#api.StudyOverview)
    - [SuggestionParameter](#api.SuggestionParameter)
    - [SuggestionParameterSet](#api.SuggestionParameterSet)
    - [Tag](#api.Tag)
    - [Trial](#api.Trial)
    - [UpdateWorkerStateReply](#api.UpdateWorkerStateReply)
    - [UpdateWorkerStateRequest](#api.UpdateWorkerStateRequest)
    - [Worker](#api.Worker)
    - [WorkerFullInfo](#api.WorkerFullInfo)
  
    - [OptimizationType](#api.OptimizationType)
    - [ParameterType](#api.ParameterType)
    - [State](#api.State)
  
  
    - [EarlyStopping](#api.EarlyStopping)
    - [Manager](#api.Manager)
    - [Suggestion](#api.Suggestion)
  

- [Scalar Value Types](#scalar-value-types)



<a name="api.proto"/>
<p align="right"><a href="#top">Top</a></p>

## api.proto
Katib API


<a name="api.CreateStudyReply"/>

### CreateStudyReply
Return generated StudyID.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_id | [string](#string) |  |  |






<a name="api.CreateStudyRequest"/>

### CreateStudyRequest
Create a Study from Study Config.
Generate an unique ID and store the Study to DB.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_config | [StudyConfig](#api.StudyConfig) |  |  |






<a name="api.CreateTrialReply"/>

### CreateTrialReply
Return generated TrialID.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trial_id | [string](#string) |  |  |






<a name="api.CreateTrialRequest"/>

### CreateTrialRequest
Create a Trial from Trial Config.
Generate an unique ID and store the Trial to DB.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trial | [Trial](#api.Trial) |  |  |






<a name="api.DataSetInfo"/>

### DataSetInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| path | [string](#string) |  |  |






<a name="api.DeleteStudyReply"/>

### DeleteStudyReply
Return deleted Study ID.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_id | [string](#string) |  |  |






<a name="api.DeleteStudyRequest"/>

### DeleteStudyRequest
Delete a Study from DB by Study ID.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_id | [string](#string) |  |  |






<a name="api.EarlyStoppingParameter"/>

### EarlyStoppingParameter
Parameter for EarlyStopping service. Key-value format.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of Parameter. |
| value | [string](#string) |  | Value of Parameter. |






<a name="api.EarlyStoppingParameterSet"/>

### EarlyStoppingParameterSet



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| param_id | [string](#string) |  |  |
| early_stopping_algorithm | [string](#string) |  |  |
| early_stopping_parameters | [EarlyStoppingParameter](#api.EarlyStoppingParameter) | repeated |  |






<a name="api.FeasibleSpace"/>

### FeasibleSpace
Feasible space for optimization.
Int and Double type use Max/Min.
Discrete and Categorical type use List.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| max | [string](#string) |  | Max Value |
| min | [string](#string) |  | Minimum Value |
| list | [string](#string) | repeated | List of Values. |






<a name="api.GetEarlyStoppingParameterListReply"/>

### GetEarlyStoppingParameterListReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| early_stopping_parameter_sets | [EarlyStoppingParameterSet](#api.EarlyStoppingParameterSet) | repeated |  |






<a name="api.GetEarlyStoppingParameterListRequest"/>

### GetEarlyStoppingParameterListRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_id | [string](#string) |  |  |






<a name="api.GetEarlyStoppingParametersReply"/>

### GetEarlyStoppingParametersReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| early_stopping_parameters | [EarlyStoppingParameter](#api.EarlyStoppingParameter) | repeated |  |






<a name="api.GetEarlyStoppingParametersRequest"/>

### GetEarlyStoppingParametersRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| param_id | [string](#string) |  |  |






<a name="api.GetMetricsReply"/>

### GetMetricsReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| metrics_log_sets | [MetricsLogSet](#api.MetricsLogSet) | repeated |  |






<a name="api.GetMetricsRequest"/>

### GetMetricsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_id | [string](#string) |  |  |
| worker_ids | [string](#string) | repeated |  |
| metrics_names | [string](#string) | repeated |  |






<a name="api.GetSavedModelReply"/>

### GetSavedModelReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| model | [ModelInfo](#api.ModelInfo) |  |  |






<a name="api.GetSavedModelRequest"/>

### GetSavedModelRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_name | [string](#string) |  |  |
| worker_id | [string](#string) |  |  |






<a name="api.GetSavedModelsReply"/>

### GetSavedModelsReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| models | [ModelInfo](#api.ModelInfo) | repeated |  |






<a name="api.GetSavedModelsRequest"/>

### GetSavedModelsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_name | [string](#string) |  |  |






<a name="api.GetSavedStudiesReply"/>

### GetSavedStudiesReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| studies | [StudyOverview](#api.StudyOverview) | repeated |  |






<a name="api.GetSavedStudiesRequest"/>

### GetSavedStudiesRequest







<a name="api.GetShouldStopWorkersReply"/>

### GetShouldStopWorkersReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| should_stop_worker_ids | [string](#string) | repeated |  |






<a name="api.GetShouldStopWorkersRequest"/>

### GetShouldStopWorkersRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_id | [string](#string) |  |  |
| early_stopping_algorithm | [string](#string) |  |  |
| param_id | [string](#string) |  |  |






<a name="api.GetStudyListReply"/>

### GetStudyListReply
Return a overview list of Studies.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_overviews | [StudyOverview](#api.StudyOverview) | repeated |  |






<a name="api.GetStudyListRequest"/>

### GetStudyListRequest
Get all Study Configs from DB.






<a name="api.GetStudyReply"/>

### GetStudyReply
Return a config of specified Study.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_config | [StudyConfig](#api.StudyConfig) |  |  |






<a name="api.GetStudyRequest"/>

### GetStudyRequest
Get a Study Config from DB by ID of Study.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_id | [string](#string) |  |  |






<a name="api.GetSuggestionParameterListReply"/>

### GetSuggestionParameterListReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| suggestion_parameter_sets | [SuggestionParameterSet](#api.SuggestionParameterSet) | repeated |  |






<a name="api.GetSuggestionParameterListRequest"/>

### GetSuggestionParameterListRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_id | [string](#string) |  |  |






<a name="api.GetSuggestionParametersReply"/>

### GetSuggestionParametersReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| suggestion_parameters | [SuggestionParameter](#api.SuggestionParameter) | repeated |  |






<a name="api.GetSuggestionParametersRequest"/>

### GetSuggestionParametersRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| param_id | [string](#string) |  |  |






<a name="api.GetSuggestionsReply"/>

### GetSuggestionsReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trials | [Trial](#api.Trial) | repeated |  |






<a name="api.GetSuggestionsRequest"/>

### GetSuggestionsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_id | [string](#string) |  |  |
| suggestion_algorithm | [string](#string) |  |  |
| request_number | [int32](#int32) |  |  |
| log_worker_ids | [string](#string) | repeated |  |
| param_id | [string](#string) |  |  |






<a name="api.GetTrialReply"/>

### GetTrialReply
Return a trial configuration by specified trial ID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trial | [Trial](#api.Trial) |  |  |






<a name="api.GetTrialRequest"/>

### GetTrialRequest
Get a trial configuration from DB by trial ID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trial_id | [string](#string) |  |  |






<a name="api.GetTrialsReply"/>

### GetTrialsReply
Return a trial list in specified Study.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trials | [Trial](#api.Trial) | repeated |  |






<a name="api.GetTrialsRequest"/>

### GetTrialsRequest
Get a Trial Configs from DB by ID of Study.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_id | [string](#string) |  |  |






<a name="api.GetWorkerFullInfoReply"/>

### GetWorkerFullInfoReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| worker_full_infos | [WorkerFullInfo](#api.WorkerFullInfo) | repeated |  |






<a name="api.GetWorkerFullInfoRequest"/>

### GetWorkerFullInfoRequest
Get a full information related to specified Workers.
It includes Worker Config, HyperParameters and Metrics Logs.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_id | [string](#string) |  |  |
| trial_id | [string](#string) |  |  |
| worker_id | [string](#string) |  |  |
| only_latest_log | [bool](#bool) |  |  |






<a name="api.GetWorkersReply"/>

### GetWorkersReply
Return a Worker list by specified condition.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| workers | [Worker](#api.Worker) | repeated |  |






<a name="api.GetWorkersRequest"/>

### GetWorkersRequest
Get a configs and status of a Worker from DB by ID of Study, Trial or Worker.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_id | [string](#string) |  |  |
| trial_id | [string](#string) |  |  |
| worker_id | [string](#string) |  |  |






<a name="api.Metrics"/>

### Metrics
Metrics of a worker


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of metrics. |
| value | [string](#string) |  | Value of metrics. Double float. |






<a name="api.MetricsLog"/>

### MetricsLog
Metrics logs of a worker


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of metrics. |
| values | [MetricsValueTime](#api.MetricsValueTime) | repeated | Log of metrics. Ordered by time series. |






<a name="api.MetricsLogSet"/>

### MetricsLogSet
Logs of metrics for a worker.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| worker_id | [string](#string) |  | ID of the corresponding worker. |
| metrics_logs | [MetricsLog](#api.MetricsLog) | repeated | Logs of metrics. |
| worker_status | [State](#api.State) |  | Status of the corresponding worker. |






<a name="api.MetricsValueTime"/>

### MetricsValueTime
Metrics of a worker with timestamp


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| time | [string](#string) |  | Timestamp RFC3339 format. |
| value | [string](#string) |  | Value of metrics. Double float. |






<a name="api.ModelInfo"/>

### ModelInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_name | [string](#string) |  |  |
| worker_id | [string](#string) |  |  |
| parameters | [Parameter](#api.Parameter) | repeated |  |
| metrics | [Metrics](#api.Metrics) | repeated |  |
| model_path | [string](#string) |  |  |






<a name="api.Parameter"/>

### Parameter
Value of a Hyper parameter.
This will be created from a correcponding Config.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the parameter. |
| parameter_type | [ParameterType](#api.ParameterType) |  | Type of the parameter. |
| value | [string](#string) |  | Value of the parameter. |






<a name="api.ParameterConfig"/>

### ParameterConfig
Config for a Hyper parameter.
Katib will create each Hyper parameter from this config.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the parameter. |
| parameter_type | [ParameterType](#api.ParameterType) |  | Type of the parameter. |
| feasible | [FeasibleSpace](#api.FeasibleSpace) |  | FeasibleSpace for the parameter. |






<a name="api.RegisterWorkerReply"/>

### RegisterWorkerReply
Return generated WorkerID.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| worker_id | [string](#string) |  |  |






<a name="api.RegisterWorkerRequest"/>

### RegisterWorkerRequest
Create a Worker from Worker Config.
Generate an unique ID and store the Worker to DB.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| worker | [Worker](#api.Worker) |  |  |






<a name="api.ReportMetricsLogsReply"/>

### ReportMetricsLogsReply







<a name="api.ReportMetricsLogsRequest"/>

### ReportMetricsLogsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_id | [string](#string) |  |  |
| metrics_log_sets | [MetricsLogSet](#api.MetricsLogSet) | repeated |  |






<a name="api.SaveModelReply"/>

### SaveModelReply







<a name="api.SaveModelRequest"/>

### SaveModelRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| model | [ModelInfo](#api.ModelInfo) |  |  |
| data_set | [DataSetInfo](#api.DataSetInfo) |  |  |
| tensor_board | [bool](#bool) |  |  |






<a name="api.SaveStudyReply"/>

### SaveStudyReply







<a name="api.SaveStudyRequest"/>

### SaveStudyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_name | [string](#string) |  |  |
| owner | [string](#string) |  |  |
| description | [string](#string) |  |  |






<a name="api.SetEarlyStoppingParametersReply"/>

### SetEarlyStoppingParametersReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| param_id | [string](#string) |  |  |






<a name="api.SetEarlyStoppingParametersRequest"/>

### SetEarlyStoppingParametersRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_id | [string](#string) |  |  |
| early_stopping_algorithm | [string](#string) |  |  |
| param_id | [string](#string) |  |  |
| early_stopping_parameters | [EarlyStoppingParameter](#api.EarlyStoppingParameter) | repeated |  |






<a name="api.SetSuggestionParametersReply"/>

### SetSuggestionParametersReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| param_id | [string](#string) |  |  |






<a name="api.SetSuggestionParametersRequest"/>

### SetSuggestionParametersRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_id | [string](#string) |  |  |
| suggestion_algorithm | [string](#string) |  |  |
| param_id | [string](#string) |  |  |
| suggestion_parameters | [SuggestionParameter](#api.SuggestionParameter) | repeated |  |






<a name="api.StopSuggestionReply"/>

### StopSuggestionReply







<a name="api.StopSuggestionRequest"/>

### StopSuggestionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_id | [string](#string) |  |  |






<a name="api.StopWorkersReply"/>

### StopWorkersReply







<a name="api.StopWorkersRequest"/>

### StopWorkersRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| study_id | [string](#string) |  |  |
| worker_ids | [string](#string) | repeated |  |
| is_complete | [bool](#bool) |  |  |






<a name="api.StudyConfig"/>

### StudyConfig
Config of a Study. Study represents a single optimization run over a feasible space. 
Each Study contains a configuration describing the feasible space, as well as a set of Trials.
It is assumed that objective function f(x) does not change in the course of a Study.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of Study. |
| owner | [string](#string) |  | Owner of Study. |
| optimization_type | [OptimizationType](#api.OptimizationType) |  | Optimization type. |
| optimization_goal | [double](#double) |  | Goal of optimization value. |
| parameter_configs | [StudyConfig.ParameterConfigs](#api.StudyConfig.ParameterConfigs) |  | List of ParameterConfig |
| access_permissions | [string](#string) | repeated | Access Permission |
| tags | [Tag](#api.Tag) | repeated | Tag for Study |
| objective_value_name | [string](#string) |  | Name of objective value. |
| metrics | [string](#string) | repeated | List of metrics name. |
| jobId | [string](#string) |  | ID of studyjob that is created from this config. |






<a name="api.StudyConfig.ParameterConfigs"/>

### StudyConfig.ParameterConfigs
List of ParameterConfig


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| configs | [ParameterConfig](#api.ParameterConfig) | repeated |  |






<a name="api.StudyOverview"/>

### StudyOverview
Overview of a study. For UI.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of Study. |
| owner | [string](#string) |  | Owner of Study. |
| id | [string](#string) |  | Study ID. |
| description | [string](#string) |  | Discretption of Study. |






<a name="api.SuggestionParameter"/>

### SuggestionParameter
Parameter for Suggestion service. Key-value format.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of Parameter. |
| value | [string](#string) |  | Value of Parameter. |






<a name="api.SuggestionParameterSet"/>

### SuggestionParameterSet



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| param_id | [string](#string) |  |  |
| suggestion_algorithm | [string](#string) |  |  |
| suggestion_parameters | [SuggestionParameter](#api.SuggestionParameter) | repeated |  |






<a name="api.Tag"/>

### Tag
Tag for each resource.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of tag. |
| value | [string](#string) |  | Value of tag. |






<a name="api.Trial"/>

### Trial
A set of Hyperparameter.
In a study, multiple trials are evaluated by workers.
Suggestion service will generate next trials.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trial_id | [string](#string) |  | Trial ID. |
| study_id | [string](#string) |  | Study ID. |
| parameter_set | [Parameter](#api.Parameter) | repeated | Hyperparameter set |
| objective_value | [string](#string) |  | Objective Value |
| tags | [Tag](#api.Tag) | repeated | Tags of Trial. |






<a name="api.UpdateWorkerStateReply"/>

### UpdateWorkerStateReply







<a name="api.UpdateWorkerStateRequest"/>

### UpdateWorkerStateRequest
Update a Status of Worker.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| worker_id | [string](#string) |  |  |
| status | [State](#api.State) |  |  |






<a name="api.Worker"/>

### Worker
A process of evaluation for a trial.
Types of worker supported by Katib are k8s Job, TF-Job, and Pytorch-Job.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| worker_id | [string](#string) |  | Worker ID. |
| study_id | [string](#string) |  | Study ID. |
| trial_id | [string](#string) |  | Trial ID. |
| Type | [string](#string) |  | Type of Worker |
| status | [State](#api.State) |  | Status of Worker. |
| TemplatePath | [string](#string) |  | Path for the manufest template of Worker. |
| tags | [Tag](#api.Tag) | repeated | Tags of Worker. |






<a name="api.WorkerFullInfo"/>

### WorkerFullInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Worker | [Worker](#api.Worker) |  |  |
| parameter_set | [Parameter](#api.Parameter) | repeated |  |
| metrics_logs | [MetricsLog](#api.MetricsLog) | repeated |  |





 


<a name="api.OptimizationType"/>

### OptimizationType
Direction of optimization. Minimize or Maximize.

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN_OPTIMIZATION | 0 | Undefined type and not used. |
| MINIMIZE | 1 | Minimize |
| MAXIMIZE | 2 | Maximize |



<a name="api.ParameterType"/>

### ParameterType
Types of value for HyperParameter.

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN_TYPE | 0 | Undefined type and not used. |
| DOUBLE | 1 | Double float type. Use &#34;Max/Min&#34;. |
| INT | 2 | Int type. Use &#34;Max/Min&#34;. |
| DISCRETE | 3 | Discrete number type. Use &#34;List&#34; as float. |
| CATEGORICAL | 4 | Categorical type. Use &#34;List&#34; as string. |



<a name="api.State"/>

### State
Status code for worker.
This value is stored as TINYINT in MySQL.

| Name | Number | Description |
| ---- | ------ | ----------- |
| PENDING | 0 | Pending. Created but not running. |
| RUNNING | 1 | Running. |
| COMPLETED | 2 | Completed. |
| KILLED | 3 | Killed. Not failed. |
| ERROR | 120 | Error. |


 

 


<a name="api.EarlyStopping"/>

### EarlyStopping


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetShouldStopWorkers | [GetShouldStopWorkersRequest](#api.GetShouldStopWorkersRequest) | [GetShouldStopWorkersReply](#api.GetShouldStopWorkersRequest) |  |


<a name="api.Manager"/>

### Manager
Service for Main API for Katib
For each RPC service, we define mapping to HTTP REST API method.
The mapping includes the URL path, query parameters and request body.
https://cloud.google.com/service-infrastructure/docs/service-management/reference/rpc/google.api#http

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateStudy | [CreateStudyRequest](#api.CreateStudyRequest) | [CreateStudyReply](#api.CreateStudyRequest) | Create a Study from Study Config. Generate a unique ID and store the Study to DB. |
| GetStudy | [GetStudyRequest](#api.GetStudyRequest) | [GetStudyReply](#api.GetStudyRequest) | Get a Study Config from DB by ID of Study. |
| DeleteStudy | [DeleteStudyRequest](#api.DeleteStudyRequest) | [DeleteStudyReply](#api.DeleteStudyRequest) | Delete a Study from DB by Study ID. |
| GetStudyList | [GetStudyListRequest](#api.GetStudyListRequest) | [GetStudyListReply](#api.GetStudyListRequest) | Get all Study Configs from DB. |
| CreateTrial | [CreateTrialRequest](#api.CreateTrialRequest) | [CreateTrialReply](#api.CreateTrialRequest) | Create a Trial from Trial Config. Generate a unique ID and store the Trial to DB. |
| GetTrials | [GetTrialsRequest](#api.GetTrialsRequest) | [GetTrialsReply](#api.GetTrialsRequest) | Get a Trial Configs from DB by ID of Study. |
| GetTrial | [GetTrialRequest](#api.GetTrialRequest) | [GetTrialReply](#api.GetTrialRequest) | Get a Trial Configuration from DB by ID of Trial. |
| RegisterWorker | [RegisterWorkerRequest](#api.RegisterWorkerRequest) | [RegisterWorkerReply](#api.RegisterWorkerRequest) | Create a Worker from Worker Config. Generate a unique ID and store the Worker to DB. |
| GetWorkers | [GetWorkersRequest](#api.GetWorkersRequest) | [GetWorkersReply](#api.GetWorkersRequest) | Get a Worker Configs and Status from DB by ID of Study, Trial or Worker. |
| UpdateWorkerState | [UpdateWorkerStateRequest](#api.UpdateWorkerStateRequest) | [UpdateWorkerStateReply](#api.UpdateWorkerStateRequest) | Update a Status of Worker. |
| GetWorkerFullInfo | [GetWorkerFullInfoRequest](#api.GetWorkerFullInfoRequest) | [GetWorkerFullInfoReply](#api.GetWorkerFullInfoRequest) | Get full information related to specified Workers. It includes Worker Config, HyperParameters and Metrics Logs. |
| GetSuggestions | [GetSuggestionsRequest](#api.GetSuggestionsRequest) | [GetSuggestionsReply](#api.GetSuggestionsRequest) | Get Suggestions from a Suggestion service. |
| GetShouldStopWorkers | [GetShouldStopWorkersRequest](#api.GetShouldStopWorkersRequest) | [GetShouldStopWorkersReply](#api.GetShouldStopWorkersRequest) |  |
| GetMetrics | [GetMetricsRequest](#api.GetMetricsRequest) | [GetMetricsReply](#api.GetMetricsRequest) | Get metrics of workers. You can get all logs of metrics since start of the worker. |
| SetSuggestionParameters | [SetSuggestionParametersRequest](#api.SetSuggestionParametersRequest) | [SetSuggestionParametersReply](#api.SetSuggestionParametersRequest) | Create or Update parameter set for a suggestion service. If you specify an ID of parameter set, it will update the parameter set by your request. If you don&#39;t specify an ID, it will create a new parameter set for corresponding study and suggestion service. The parameters are key-value format. |
| GetSuggestionParameters | [GetSuggestionParametersRequest](#api.GetSuggestionParametersRequest) | [GetSuggestionParametersReply](#api.GetSuggestionParametersRequest) | Get suggestion parameter set from DB specified. |
| GetSuggestionParameterList | [GetSuggestionParameterListRequest](#api.GetSuggestionParameterListRequest) | [GetSuggestionParameterListReply](#api.GetSuggestionParameterListRequest) | Get all suggestion parameter sets from DB. |
| SetEarlyStoppingParameters | [SetEarlyStoppingParametersRequest](#api.SetEarlyStoppingParametersRequest) | [SetEarlyStoppingParametersReply](#api.SetEarlyStoppingParametersRequest) |  |
| GetEarlyStoppingParameters | [GetEarlyStoppingParametersRequest](#api.GetEarlyStoppingParametersRequest) | [GetEarlyStoppingParametersReply](#api.GetEarlyStoppingParametersRequest) |  |
| GetEarlyStoppingParameterList | [GetEarlyStoppingParameterListRequest](#api.GetEarlyStoppingParameterListRequest) | [GetEarlyStoppingParameterListReply](#api.GetEarlyStoppingParameterListRequest) |  |
| SaveStudy | [SaveStudyRequest](#api.SaveStudyRequest) | [SaveStudyReply](#api.SaveStudyRequest) |  |
| SaveModel | [SaveModelRequest](#api.SaveModelRequest) | [SaveModelReply](#api.SaveModelRequest) |  |
| ReportMetricsLogs | [ReportMetricsLogsRequest](#api.ReportMetricsLogsRequest) | [ReportMetricsLogsReply](#api.ReportMetricsLogsRequest) | Report a logs of metrics for workers. The logs for each worker must have timestamp and must be ordered in time series. When the log you reported are already reported before, it will be dismissed and get no error. |
| GetSavedStudies | [GetSavedStudiesRequest](#api.GetSavedStudiesRequest) | [GetSavedStudiesReply](#api.GetSavedStudiesRequest) |  |
| GetSavedModels | [GetSavedModelsRequest](#api.GetSavedModelsRequest) | [GetSavedModelsReply](#api.GetSavedModelsRequest) |  |


<a name="api.Suggestion"/>

### Suggestion


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetSuggestions | [GetSuggestionsRequest](#api.GetSuggestionsRequest) | [GetSuggestionsReply](#api.GetSuggestionsRequest) |  |

 



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

