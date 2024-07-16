from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ParameterType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    UNKNOWN_TYPE: _ClassVar[ParameterType]
    DOUBLE: _ClassVar[ParameterType]
    INT: _ClassVar[ParameterType]
    DISCRETE: _ClassVar[ParameterType]
    CATEGORICAL: _ClassVar[ParameterType]

class ObjectiveType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    UNKNOWN: _ClassVar[ObjectiveType]
    MINIMIZE: _ClassVar[ObjectiveType]
    MAXIMIZE: _ClassVar[ObjectiveType]

class ComparisonType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    UNKNOWN_COMPARISON: _ClassVar[ComparisonType]
    EQUAL: _ClassVar[ComparisonType]
    LESS: _ClassVar[ComparisonType]
    GREATER: _ClassVar[ComparisonType]
UNKNOWN_TYPE: ParameterType
DOUBLE: ParameterType
INT: ParameterType
DISCRETE: ParameterType
CATEGORICAL: ParameterType
UNKNOWN: ObjectiveType
MINIMIZE: ObjectiveType
MAXIMIZE: ObjectiveType
UNKNOWN_COMPARISON: ComparisonType
EQUAL: ComparisonType
LESS: ComparisonType
GREATER: ComparisonType

class Experiment(_message.Message):
    __slots__ = ("name", "spec")
    NAME_FIELD_NUMBER: _ClassVar[int]
    SPEC_FIELD_NUMBER: _ClassVar[int]
    name: str
    spec: ExperimentSpec
    def __init__(self, name: _Optional[str] = ..., spec: _Optional[_Union[ExperimentSpec, _Mapping]] = ...) -> None: ...

class ExperimentSpec(_message.Message):
    __slots__ = ("parameter_specs", "objective", "algorithm", "early_stopping", "parallel_trial_count", "max_trial_count", "nas_config")
    class ParameterSpecs(_message.Message):
        __slots__ = ("parameters",)
        PARAMETERS_FIELD_NUMBER: _ClassVar[int]
        parameters: _containers.RepeatedCompositeFieldContainer[ParameterSpec]
        def __init__(self, parameters: _Optional[_Iterable[_Union[ParameterSpec, _Mapping]]] = ...) -> None: ...
    PARAMETER_SPECS_FIELD_NUMBER: _ClassVar[int]
    OBJECTIVE_FIELD_NUMBER: _ClassVar[int]
    ALGORITHM_FIELD_NUMBER: _ClassVar[int]
    EARLY_STOPPING_FIELD_NUMBER: _ClassVar[int]
    PARALLEL_TRIAL_COUNT_FIELD_NUMBER: _ClassVar[int]
    MAX_TRIAL_COUNT_FIELD_NUMBER: _ClassVar[int]
    NAS_CONFIG_FIELD_NUMBER: _ClassVar[int]
    parameter_specs: ExperimentSpec.ParameterSpecs
    objective: ObjectiveSpec
    algorithm: AlgorithmSpec
    early_stopping: EarlyStoppingSpec
    parallel_trial_count: int
    max_trial_count: int
    nas_config: NasConfig
    def __init__(self, parameter_specs: _Optional[_Union[ExperimentSpec.ParameterSpecs, _Mapping]] = ..., objective: _Optional[_Union[ObjectiveSpec, _Mapping]] = ..., algorithm: _Optional[_Union[AlgorithmSpec, _Mapping]] = ..., early_stopping: _Optional[_Union[EarlyStoppingSpec, _Mapping]] = ..., parallel_trial_count: _Optional[int] = ..., max_trial_count: _Optional[int] = ..., nas_config: _Optional[_Union[NasConfig, _Mapping]] = ...) -> None: ...

class ParameterSpec(_message.Message):
    __slots__ = ("name", "parameter_type", "feasible_space")
    NAME_FIELD_NUMBER: _ClassVar[int]
    PARAMETER_TYPE_FIELD_NUMBER: _ClassVar[int]
    FEASIBLE_SPACE_FIELD_NUMBER: _ClassVar[int]
    name: str
    parameter_type: ParameterType
    feasible_space: FeasibleSpace
    def __init__(self, name: _Optional[str] = ..., parameter_type: _Optional[_Union[ParameterType, str]] = ..., feasible_space: _Optional[_Union[FeasibleSpace, _Mapping]] = ...) -> None: ...

class FeasibleSpace(_message.Message):
    __slots__ = ("max", "min", "list", "step")
    MAX_FIELD_NUMBER: _ClassVar[int]
    MIN_FIELD_NUMBER: _ClassVar[int]
    LIST_FIELD_NUMBER: _ClassVar[int]
    STEP_FIELD_NUMBER: _ClassVar[int]
    max: str
    min: str
    list: _containers.RepeatedScalarFieldContainer[str]
    step: str
    def __init__(self, max: _Optional[str] = ..., min: _Optional[str] = ..., list: _Optional[_Iterable[str]] = ..., step: _Optional[str] = ...) -> None: ...

class ObjectiveSpec(_message.Message):
    __slots__ = ("type", "goal", "objective_metric_name", "additional_metric_names")
    TYPE_FIELD_NUMBER: _ClassVar[int]
    GOAL_FIELD_NUMBER: _ClassVar[int]
    OBJECTIVE_METRIC_NAME_FIELD_NUMBER: _ClassVar[int]
    ADDITIONAL_METRIC_NAMES_FIELD_NUMBER: _ClassVar[int]
    type: ObjectiveType
    goal: float
    objective_metric_name: str
    additional_metric_names: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, type: _Optional[_Union[ObjectiveType, str]] = ..., goal: _Optional[float] = ..., objective_metric_name: _Optional[str] = ..., additional_metric_names: _Optional[_Iterable[str]] = ...) -> None: ...

class AlgorithmSpec(_message.Message):
    __slots__ = ("algorithm_name", "algorithm_settings")
    ALGORITHM_NAME_FIELD_NUMBER: _ClassVar[int]
    ALGORITHM_SETTINGS_FIELD_NUMBER: _ClassVar[int]
    algorithm_name: str
    algorithm_settings: _containers.RepeatedCompositeFieldContainer[AlgorithmSetting]
    def __init__(self, algorithm_name: _Optional[str] = ..., algorithm_settings: _Optional[_Iterable[_Union[AlgorithmSetting, _Mapping]]] = ...) -> None: ...

class AlgorithmSetting(_message.Message):
    __slots__ = ("name", "value")
    NAME_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    name: str
    value: str
    def __init__(self, name: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...

class EarlyStoppingSpec(_message.Message):
    __slots__ = ("algorithm_name", "algorithm_settings")
    ALGORITHM_NAME_FIELD_NUMBER: _ClassVar[int]
    ALGORITHM_SETTINGS_FIELD_NUMBER: _ClassVar[int]
    algorithm_name: str
    algorithm_settings: _containers.RepeatedCompositeFieldContainer[EarlyStoppingSetting]
    def __init__(self, algorithm_name: _Optional[str] = ..., algorithm_settings: _Optional[_Iterable[_Union[EarlyStoppingSetting, _Mapping]]] = ...) -> None: ...

class EarlyStoppingSetting(_message.Message):
    __slots__ = ("name", "value")
    NAME_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    name: str
    value: str
    def __init__(self, name: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...

class NasConfig(_message.Message):
    __slots__ = ("graph_config", "operations")
    class Operations(_message.Message):
        __slots__ = ("operation",)
        OPERATION_FIELD_NUMBER: _ClassVar[int]
        operation: _containers.RepeatedCompositeFieldContainer[Operation]
        def __init__(self, operation: _Optional[_Iterable[_Union[Operation, _Mapping]]] = ...) -> None: ...
    GRAPH_CONFIG_FIELD_NUMBER: _ClassVar[int]
    OPERATIONS_FIELD_NUMBER: _ClassVar[int]
    graph_config: GraphConfig
    operations: NasConfig.Operations
    def __init__(self, graph_config: _Optional[_Union[GraphConfig, _Mapping]] = ..., operations: _Optional[_Union[NasConfig.Operations, _Mapping]] = ...) -> None: ...

class GraphConfig(_message.Message):
    __slots__ = ("num_layers", "input_sizes", "output_sizes")
    NUM_LAYERS_FIELD_NUMBER: _ClassVar[int]
    INPUT_SIZES_FIELD_NUMBER: _ClassVar[int]
    OUTPUT_SIZES_FIELD_NUMBER: _ClassVar[int]
    num_layers: int
    input_sizes: _containers.RepeatedScalarFieldContainer[int]
    output_sizes: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, num_layers: _Optional[int] = ..., input_sizes: _Optional[_Iterable[int]] = ..., output_sizes: _Optional[_Iterable[int]] = ...) -> None: ...

class Operation(_message.Message):
    __slots__ = ("operation_type", "parameter_specs")
    class ParameterSpecs(_message.Message):
        __slots__ = ("parameters",)
        PARAMETERS_FIELD_NUMBER: _ClassVar[int]
        parameters: _containers.RepeatedCompositeFieldContainer[ParameterSpec]
        def __init__(self, parameters: _Optional[_Iterable[_Union[ParameterSpec, _Mapping]]] = ...) -> None: ...
    OPERATION_TYPE_FIELD_NUMBER: _ClassVar[int]
    PARAMETER_SPECS_FIELD_NUMBER: _ClassVar[int]
    operation_type: str
    parameter_specs: Operation.ParameterSpecs
    def __init__(self, operation_type: _Optional[str] = ..., parameter_specs: _Optional[_Union[Operation.ParameterSpecs, _Mapping]] = ...) -> None: ...

class Trial(_message.Message):
    __slots__ = ("name", "spec", "status")
    NAME_FIELD_NUMBER: _ClassVar[int]
    SPEC_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    name: str
    spec: TrialSpec
    status: TrialStatus
    def __init__(self, name: _Optional[str] = ..., spec: _Optional[_Union[TrialSpec, _Mapping]] = ..., status: _Optional[_Union[TrialStatus, _Mapping]] = ...) -> None: ...

class TrialSpec(_message.Message):
    __slots__ = ("objective", "parameter_assignments", "labels")
    class ParameterAssignments(_message.Message):
        __slots__ = ("assignments",)
        ASSIGNMENTS_FIELD_NUMBER: _ClassVar[int]
        assignments: _containers.RepeatedCompositeFieldContainer[ParameterAssignment]
        def __init__(self, assignments: _Optional[_Iterable[_Union[ParameterAssignment, _Mapping]]] = ...) -> None: ...
    class LabelsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    OBJECTIVE_FIELD_NUMBER: _ClassVar[int]
    PARAMETER_ASSIGNMENTS_FIELD_NUMBER: _ClassVar[int]
    LABELS_FIELD_NUMBER: _ClassVar[int]
    objective: ObjectiveSpec
    parameter_assignments: TrialSpec.ParameterAssignments
    labels: _containers.ScalarMap[str, str]
    def __init__(self, objective: _Optional[_Union[ObjectiveSpec, _Mapping]] = ..., parameter_assignments: _Optional[_Union[TrialSpec.ParameterAssignments, _Mapping]] = ..., labels: _Optional[_Mapping[str, str]] = ...) -> None: ...

class ParameterAssignment(_message.Message):
    __slots__ = ("name", "value")
    NAME_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    name: str
    value: str
    def __init__(self, name: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...

class TrialStatus(_message.Message):
    __slots__ = ("start_time", "completion_time", "condition", "observation")
    class TrialConditionType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        CREATED: _ClassVar[TrialStatus.TrialConditionType]
        RUNNING: _ClassVar[TrialStatus.TrialConditionType]
        SUCCEEDED: _ClassVar[TrialStatus.TrialConditionType]
        KILLED: _ClassVar[TrialStatus.TrialConditionType]
        FAILED: _ClassVar[TrialStatus.TrialConditionType]
        METRICSUNAVAILABLE: _ClassVar[TrialStatus.TrialConditionType]
        EARLYSTOPPED: _ClassVar[TrialStatus.TrialConditionType]
        UNKNOWN: _ClassVar[TrialStatus.TrialConditionType]
    CREATED: TrialStatus.TrialConditionType
    RUNNING: TrialStatus.TrialConditionType
    SUCCEEDED: TrialStatus.TrialConditionType
    KILLED: TrialStatus.TrialConditionType
    FAILED: TrialStatus.TrialConditionType
    METRICSUNAVAILABLE: TrialStatus.TrialConditionType
    EARLYSTOPPED: TrialStatus.TrialConditionType
    UNKNOWN: TrialStatus.TrialConditionType
    START_TIME_FIELD_NUMBER: _ClassVar[int]
    COMPLETION_TIME_FIELD_NUMBER: _ClassVar[int]
    CONDITION_FIELD_NUMBER: _ClassVar[int]
    OBSERVATION_FIELD_NUMBER: _ClassVar[int]
    start_time: str
    completion_time: str
    condition: TrialStatus.TrialConditionType
    observation: Observation
    def __init__(self, start_time: _Optional[str] = ..., completion_time: _Optional[str] = ..., condition: _Optional[_Union[TrialStatus.TrialConditionType, str]] = ..., observation: _Optional[_Union[Observation, _Mapping]] = ...) -> None: ...

class Observation(_message.Message):
    __slots__ = ("metrics",)
    METRICS_FIELD_NUMBER: _ClassVar[int]
    metrics: _containers.RepeatedCompositeFieldContainer[Metric]
    def __init__(self, metrics: _Optional[_Iterable[_Union[Metric, _Mapping]]] = ...) -> None: ...

class Metric(_message.Message):
    __slots__ = ("name", "value")
    NAME_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    name: str
    value: str
    def __init__(self, name: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...

class ReportObservationLogRequest(_message.Message):
    __slots__ = ("trial_name", "observation_log")
    TRIAL_NAME_FIELD_NUMBER: _ClassVar[int]
    OBSERVATION_LOG_FIELD_NUMBER: _ClassVar[int]
    trial_name: str
    observation_log: ObservationLog
    def __init__(self, trial_name: _Optional[str] = ..., observation_log: _Optional[_Union[ObservationLog, _Mapping]] = ...) -> None: ...

class ReportObservationLogReply(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class ObservationLog(_message.Message):
    __slots__ = ("metric_logs",)
    METRIC_LOGS_FIELD_NUMBER: _ClassVar[int]
    metric_logs: _containers.RepeatedCompositeFieldContainer[MetricLog]
    def __init__(self, metric_logs: _Optional[_Iterable[_Union[MetricLog, _Mapping]]] = ...) -> None: ...

class MetricLog(_message.Message):
    __slots__ = ("time_stamp", "metric")
    TIME_STAMP_FIELD_NUMBER: _ClassVar[int]
    METRIC_FIELD_NUMBER: _ClassVar[int]
    time_stamp: str
    metric: Metric
    def __init__(self, time_stamp: _Optional[str] = ..., metric: _Optional[_Union[Metric, _Mapping]] = ...) -> None: ...

class GetObservationLogRequest(_message.Message):
    __slots__ = ("trial_name", "metric_name", "start_time", "end_time")
    TRIAL_NAME_FIELD_NUMBER: _ClassVar[int]
    METRIC_NAME_FIELD_NUMBER: _ClassVar[int]
    START_TIME_FIELD_NUMBER: _ClassVar[int]
    END_TIME_FIELD_NUMBER: _ClassVar[int]
    trial_name: str
    metric_name: str
    start_time: str
    end_time: str
    def __init__(self, trial_name: _Optional[str] = ..., metric_name: _Optional[str] = ..., start_time: _Optional[str] = ..., end_time: _Optional[str] = ...) -> None: ...

class GetObservationLogReply(_message.Message):
    __slots__ = ("observation_log",)
    OBSERVATION_LOG_FIELD_NUMBER: _ClassVar[int]
    observation_log: ObservationLog
    def __init__(self, observation_log: _Optional[_Union[ObservationLog, _Mapping]] = ...) -> None: ...

class DeleteObservationLogRequest(_message.Message):
    __slots__ = ("trial_name",)
    TRIAL_NAME_FIELD_NUMBER: _ClassVar[int]
    trial_name: str
    def __init__(self, trial_name: _Optional[str] = ...) -> None: ...

class DeleteObservationLogReply(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class GetSuggestionsRequest(_message.Message):
    __slots__ = ("experiment", "trials", "current_request_number", "total_request_number")
    EXPERIMENT_FIELD_NUMBER: _ClassVar[int]
    TRIALS_FIELD_NUMBER: _ClassVar[int]
    CURRENT_REQUEST_NUMBER_FIELD_NUMBER: _ClassVar[int]
    TOTAL_REQUEST_NUMBER_FIELD_NUMBER: _ClassVar[int]
    experiment: Experiment
    trials: _containers.RepeatedCompositeFieldContainer[Trial]
    current_request_number: int
    total_request_number: int
    def __init__(self, experiment: _Optional[_Union[Experiment, _Mapping]] = ..., trials: _Optional[_Iterable[_Union[Trial, _Mapping]]] = ..., current_request_number: _Optional[int] = ..., total_request_number: _Optional[int] = ...) -> None: ...

class GetSuggestionsReply(_message.Message):
    __slots__ = ("parameter_assignments", "algorithm", "early_stopping_rules")
    class ParameterAssignments(_message.Message):
        __slots__ = ("assignments", "trial_name", "labels")
        class LabelsEntry(_message.Message):
            __slots__ = ("key", "value")
            KEY_FIELD_NUMBER: _ClassVar[int]
            VALUE_FIELD_NUMBER: _ClassVar[int]
            key: str
            value: str
            def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
        ASSIGNMENTS_FIELD_NUMBER: _ClassVar[int]
        TRIAL_NAME_FIELD_NUMBER: _ClassVar[int]
        LABELS_FIELD_NUMBER: _ClassVar[int]
        assignments: _containers.RepeatedCompositeFieldContainer[ParameterAssignment]
        trial_name: str
        labels: _containers.ScalarMap[str, str]
        def __init__(self, assignments: _Optional[_Iterable[_Union[ParameterAssignment, _Mapping]]] = ..., trial_name: _Optional[str] = ..., labels: _Optional[_Mapping[str, str]] = ...) -> None: ...
    PARAMETER_ASSIGNMENTS_FIELD_NUMBER: _ClassVar[int]
    ALGORITHM_FIELD_NUMBER: _ClassVar[int]
    EARLY_STOPPING_RULES_FIELD_NUMBER: _ClassVar[int]
    parameter_assignments: _containers.RepeatedCompositeFieldContainer[GetSuggestionsReply.ParameterAssignments]
    algorithm: AlgorithmSpec
    early_stopping_rules: _containers.RepeatedCompositeFieldContainer[EarlyStoppingRule]
    def __init__(self, parameter_assignments: _Optional[_Iterable[_Union[GetSuggestionsReply.ParameterAssignments, _Mapping]]] = ..., algorithm: _Optional[_Union[AlgorithmSpec, _Mapping]] = ..., early_stopping_rules: _Optional[_Iterable[_Union[EarlyStoppingRule, _Mapping]]] = ...) -> None: ...

class ValidateAlgorithmSettingsRequest(_message.Message):
    __slots__ = ("experiment",)
    EXPERIMENT_FIELD_NUMBER: _ClassVar[int]
    experiment: Experiment
    def __init__(self, experiment: _Optional[_Union[Experiment, _Mapping]] = ...) -> None: ...

class ValidateAlgorithmSettingsReply(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class GetEarlyStoppingRulesRequest(_message.Message):
    __slots__ = ("experiment", "trials", "db_manager_address")
    EXPERIMENT_FIELD_NUMBER: _ClassVar[int]
    TRIALS_FIELD_NUMBER: _ClassVar[int]
    DB_MANAGER_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    experiment: Experiment
    trials: _containers.RepeatedCompositeFieldContainer[Trial]
    db_manager_address: str
    def __init__(self, experiment: _Optional[_Union[Experiment, _Mapping]] = ..., trials: _Optional[_Iterable[_Union[Trial, _Mapping]]] = ..., db_manager_address: _Optional[str] = ...) -> None: ...

class GetEarlyStoppingRulesReply(_message.Message):
    __slots__ = ("early_stopping_rules",)
    EARLY_STOPPING_RULES_FIELD_NUMBER: _ClassVar[int]
    early_stopping_rules: _containers.RepeatedCompositeFieldContainer[EarlyStoppingRule]
    def __init__(self, early_stopping_rules: _Optional[_Iterable[_Union[EarlyStoppingRule, _Mapping]]] = ...) -> None: ...

class EarlyStoppingRule(_message.Message):
    __slots__ = ("name", "value", "comparison", "start_step")
    NAME_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    COMPARISON_FIELD_NUMBER: _ClassVar[int]
    START_STEP_FIELD_NUMBER: _ClassVar[int]
    name: str
    value: str
    comparison: ComparisonType
    start_step: int
    def __init__(self, name: _Optional[str] = ..., value: _Optional[str] = ..., comparison: _Optional[_Union[ComparisonType, str]] = ..., start_step: _Optional[int] = ...) -> None: ...

class ValidateEarlyStoppingSettingsRequest(_message.Message):
    __slots__ = ("early_stopping",)
    EARLY_STOPPING_FIELD_NUMBER: _ClassVar[int]
    early_stopping: EarlyStoppingSpec
    def __init__(self, early_stopping: _Optional[_Union[EarlyStoppingSpec, _Mapping]] = ...) -> None: ...

class ValidateEarlyStoppingSettingsReply(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class SetTrialStatusRequest(_message.Message):
    __slots__ = ("trial_name",)
    TRIAL_NAME_FIELD_NUMBER: _ClassVar[int]
    trial_name: str
    def __init__(self, trial_name: _Optional[str] = ...) -> None: ...

class SetTrialStatusReply(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
