export const CHANGE_YAML_HP = 'CHANGE_YAML_HP';

export const changeYaml = yaml => ({
  type: CHANGE_YAML_HP,
  payload: yaml,
});

export const CHANGE_META_HP = 'CHANGE_META_HP';

export const changeMeta = (name, value) => ({
  type: CHANGE_META_HP,
  name,
  value,
});

export const CHANGE_SPEC_HP = 'CHANGE_SPEC_HP';

export const changeSpec = (name, value) => ({
  type: CHANGE_SPEC_HP,
  name,
  value,
});

export const CHANGE_OBJECTIVE_HP = 'CHANGE_OBJECTIVE_HP';

export const changeObjective = (name, value) => ({
  type: CHANGE_OBJECTIVE_HP,
  name,
  value,
});
export const ADD_METRICS_HP = 'ADD_METRICS_HP';

export const addMetrics = () => ({
  type: ADD_METRICS_HP,
});

export const DELETE_METRICS_HP = 'DELETE_METRICS_HP';

export const deleteMetrics = index => ({
  type: DELETE_METRICS_HP,
  index,
});

export const EDIT_METRICS_HP = 'EDIT_METRICS_HP';

export const editMetrics = (index, value) => ({
  type: EDIT_METRICS_HP,
  index,
  value,
});

export const CHANGE_METRIC_STRATEGY_HP = 'CHANGE_METRIC_STRATEGY_HP';

export const metricStrategyChange = (index, strategy) => ({
  type: CHANGE_METRIC_STRATEGY_HP,
  index,
  strategy,
});

export const CHANGE_ALGORITHM_NAME_HP = 'CHANGE_ALGORITHM_NAME_HP';

export const changeAlgorithmName = algorithmName => ({
  type: CHANGE_ALGORITHM_NAME_HP,
  algorithmName,
});

export const ADD_ALGORITHM_SETTING_HP = 'ADD_ALGORITHM_SETTING_HP';

export const addAlgorithmSetting = () => ({
  type: ADD_ALGORITHM_SETTING_HP,
});

export const CHANGE_ALGORITHM_SETTING_HP = 'CHANGE_ALGORITHM_SETTING_HP';

export const changeAlgorithmSetting = (index, field, value) => ({
  type: CHANGE_ALGORITHM_SETTING_HP,
  field,
  value,
  index,
});

export const DELETE_ALGORITHM_SETTING_HP = 'DELETE_ALGORITHM_SETTING_HP';

export const deleteAlgorithmSetting = index => ({
  type: DELETE_ALGORITHM_SETTING_HP,
  index,
});

export const ADD_PARAMETER_HP = 'CHANGE_PARAMETER_HP';

export const addParameter = () => ({
  type: ADD_PARAMETER_HP,
});

export const EDIT_PARAMETER_HP = 'EDIT_PARAMTER_HP';

export const editParameter = (index, field, value) => ({
  type: EDIT_PARAMETER_HP,
  index,
  field,
  value,
});

export const DELETE_PARAMETER_HP = 'DELETE_PARAMETER_HP';

export const deleteParameter = index => ({
  type: DELETE_PARAMETER_HP,
  index,
});

export const ADD_LIST_PARAMETER_HP = 'ADD_LIST_PARAMETER_HP';

export const addListParameter = paramIndex => ({
  type: ADD_LIST_PARAMETER_HP,
  paramIndex,
});

export const EDIT_LIST_PARAMETER_HP = 'EDIT_LIST_PARAMETER_HP';

export const editListParameter = (paramIndex, index, value) => ({
  type: EDIT_LIST_PARAMETER_HP,
  paramIndex,
  index,
  value,
});

export const DELETE_LIST_PARAMETER_HP = 'DELETE_LIST_PARAMETER_HP';

export const deleteListParameter = (paramIndex, index) => ({
  type: DELETE_LIST_PARAMETER_HP,
  paramIndex,
  index,
});

export const SUBMIT_HP_JOB_REQUEST = 'SUBMIT_HP_JOB_REQUEST';
export const SUBMIT_HP_JOB_SUCCESS = 'SUBMIT_HP_JOB_SUCCESS';
export const SUBMIT_HP_JOB_FAILURE = 'SUBMIT_HP_JOB_FAILURE';

export const submitHPJob = data => ({
  type: SUBMIT_HP_JOB_REQUEST,
  data,
});

export const CHANGE_MC_KIND_HP = 'CHANGE_MC_KIND_HP';

export const changeMCKindHP = kind => ({
  type: CHANGE_MC_KIND_HP,
  kind,
});

export const CHANGE_MC_FILE_SYSTEM_HP = 'CHANGE_MC_FILE_SYSTEM_HP';

export const changeMCFileSystemHP = (kind, path) => ({
  type: CHANGE_MC_FILE_SYSTEM_HP,
  kind,
  path,
});

export const CHANGE_MC_HTTP_GET_HP = 'CHANGE_MC_HTTP_GET_HP';

export const changeMCHttpGetHP = (port, path, scheme, host) => ({
  type: CHANGE_MC_HTTP_GET_HP,
  port,
  path,
  scheme,
  host,
});

export const ADD_MC_HTTP_GET_HEADER_HP = 'ADD_MC_HTTP_GET_HEADER_HP';

export const addMCHttpGetHeaderHP = () => ({
  type: ADD_MC_HTTP_GET_HEADER_HP,
});

export const CHANGE_MC_HTTP_GET_HEADER_HP = 'CHANGE_MC_HTTP_GET_HEADER_HP';

export const changeMCHttpGetHeaderHP = (fieldName, value, index) => ({
  type: CHANGE_MC_HTTP_GET_HEADER_HP,
  fieldName,
  value,
  index,
});

export const DELETE_MC_HTTP_GET_HEADER_HP = 'DELETE_MC_HTTP_GET_HEADER_HP';

export const deleteMCHttpGetHeaderHP = index => ({
  type: DELETE_MC_HTTP_GET_HEADER_HP,
  index,
});

export const ADD_MC_METRICS_FORMAT_HP = 'ADD_MC_METRICS_FORMAT_HP';

export const addMCMetricsFormatHP = () => ({
  type: ADD_MC_METRICS_FORMAT_HP,
});

export const CHANGE_MC_METRIC_FORMAT_HP = 'CHANGE_MC_METRIC_FORMAT_HP';

export const changeMCMetricsFormatHP = (format, index) => ({
  type: CHANGE_MC_METRIC_FORMAT_HP,
  format,
  index,
});

export const DELETE_MC_METRIC_FORMAT_HP = 'DELETE_MC_METRIC_FORMAT_HP';

export const deleteMCMetricsFormatHP = index => ({
  type: DELETE_MC_METRIC_FORMAT_HP,
  index,
});

export const CHANGE_MC_CUSTOM_CONTAINER_HP = 'CHANGE_MC_CUSTOM_CONTAINER_HP';

export const changeMCCustomContainerHP = yamlContainer => ({
  type: CHANGE_MC_CUSTOM_CONTAINER_HP,
  yamlContainer,
});
