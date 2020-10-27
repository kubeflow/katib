export const CHANGE_YAML_NAS = 'CHANGE_YAML_NAS';

export const changeYaml = yaml => ({
  type: CHANGE_YAML_NAS,
  payload: yaml,
});

export const CHANGE_META_NAS = 'CHANGE_META_NAS';

export const changeMeta = (name, value) => ({
  type: CHANGE_META_NAS,
  name,
  value,
});

export const CHANGE_SPEC_NAS = 'CHANGE_SPEC_NAS';

export const changeSpec = (name, value) => ({
  type: CHANGE_SPEC_NAS,
  name,
  value,
});

export const CHANGE_OBJECTIVE_NAS = 'CHANGE_OBJECTIVE_NAS';

export const changeObjective = (name, value) => ({
  type: CHANGE_OBJECTIVE_NAS,
  name,
  value,
});

export const ADD_METRICS_NAS = 'ADD_METRICS_NAS';

export const addMetrics = () => ({
  type: ADD_METRICS_NAS,
});

export const DELETE_METRICS_NAS = 'DELETE_METRICS_NAS';

export const deleteMetrics = index => ({
  type: DELETE_METRICS_NAS,
  index,
});

export const EDIT_METRICS_NAS = 'EDIT_METRICS_NAS';

export const editMetrics = (index, value) => ({
  type: EDIT_METRICS_NAS,
  index,
  value,
});

export const CHANGE_METRIC_STRATEGY_NAS = 'CHANGE_METRIC_STRATEGY_NAS';

export const metricStrategyChange = (index, strategy) => ({
  type: CHANGE_METRIC_STRATEGY_NAS,
  index,
  strategy,
});

export const CHANGE_ALGORITHM_NAME_NAS = 'CHANGE_ALGORITHM_NAME_NAS';

export const changeAlgorithmName = algorithmName => ({
  type: CHANGE_ALGORITHM_NAME_NAS,
  algorithmName,
});

export const ADD_ALGORITHM_SETTING_NAS = 'ADD_ALGORITHM_SETTING_NAS';

export const addAlgorithmSetting = () => ({
  type: ADD_ALGORITHM_SETTING_NAS,
});

export const CHANGE_ALGORITHM_SETTING_NAS = 'CHANGE_ALGORITHM_SETTING_NAS';

export const changeAlgorithmSetting = (index, field, value) => ({
  type: CHANGE_ALGORITHM_SETTING_NAS,
  field,
  value,
  index,
});

export const DELETE_ALGORITHM_SETTING_NAS = 'DELETE_ALGORITHM_SETTING_NAS';

export const deleteAlgorithmSetting = index => ({
  type: DELETE_ALGORITHM_SETTING_NAS,
  index,
});

export const EDIT_NUM_LAYERS = 'EDIT_NUM_LAYERS';

export const editNumLayers = value => ({
  type: EDIT_NUM_LAYERS,
  value,
});

export const ADD_SIZE = 'ADD_SIZE';

export const addSize = sizeType => ({
  type: ADD_SIZE,
  sizeType,
});

export const EDIT_SIZE = 'EDIT_SIZE';

export const editSize = (sizeType, index, value) => ({
  type: EDIT_SIZE,
  sizeType,
  index,
  value,
});

export const DELETE_SIZE = 'DELETE_SIZE';

export const deleteSize = (sizeType, index) => ({
  type: DELETE_SIZE,
  sizeType,
  index,
});

export const ADD_OPERATION = 'ADD_OPERATION';

export const addOperation = () => ({
  type: ADD_OPERATION,
});

export const DELETE_OPERATION = 'DELETE_OPERATION';

export const deleteOperation = index => ({
  type: DELETE_OPERATION,
  index,
});

export const CHANGE_OPERATION = 'CHANGE_OPERATION';

export const changeOperation = (index, value) => ({
  type: CHANGE_OPERATION,
  index,
  value,
});

export const ADD_PARAMETER_NAS = 'ADD_PARAMETER_NAS';

export const addParameter = opIndex => ({
  type: ADD_PARAMETER_NAS,
  opIndex,
});

export const CHANGE_PARAMETER_NAS = 'CHANGE_PARAMETER_NAS';

export const changeParameter = (opIndex, paramIndex, field, value) => ({
  type: CHANGE_PARAMETER_NAS,
  opIndex,
  paramIndex,
  field,
  value,
});

export const DELETE_PARAMETER_NAS = 'DELETE_PARAMETER_NAS';

export const deleteParameter = (opIndex, paramIndex) => ({
  type: DELETE_PARAMETER_NAS,
  opIndex,
  paramIndex,
});

export const ADD_LIST_PARAMETER_NAS = 'ADD_LIST_PARAMETER_NAS';

export const addListParameter = (opIndex, paramIndex) => ({
  type: ADD_LIST_PARAMETER_NAS,
  opIndex,
  paramIndex,
});

export const EDIT_LIST_PARAMETER_NAS = 'EDIT_LIST_PARAMETER_NAS';

export const editListParameter = (opIndex, paramIndex, listIndex, value) => ({
  type: EDIT_LIST_PARAMETER_NAS,
  opIndex,
  paramIndex,
  listIndex,
  value,
});

export const DELETE_LIST_PARAMETER_NAS = 'DELETE_LIST_PARAMETER_NAS';

export const deleteListParameter = (opIndex, paramIndex, listIndex) => ({
  type: DELETE_LIST_PARAMETER_NAS,
  opIndex,
  paramIndex,
  listIndex,
});

export const SUBMIT_NAS_JOB_REQUEST = 'SUBMIT_NAS_JOB_REQUEST';
export const SUBMIT_NAS_JOB_SUCCESS = 'SUBMIT_NAS_JOB_SUCCESS';
export const SUBMIT_NAS_JOB_FAILURE = 'SUBMIT_NAS_JOB_FAILURE';

export const submitNASJob = data => ({
  type: SUBMIT_NAS_JOB_REQUEST,
  data,
});

export const CLOSE_SNACKBAR = 'CLOSE_SNACKBAR';

export const closeSnackbar = () => ({
  type: CLOSE_SNACKBAR,
});

export const CHANGE_MC_KIND_NAS = 'CHANGE_MC_KIND_NAS';

export const changeMCKindNAS = kind => ({
  type: CHANGE_MC_KIND_NAS,
  kind,
});

export const CHANGE_MC_FILE_SYSTEM_NAS = 'CHANGE_MC_FILE_SYSTEM_NAS';

export const changeMCFileSystemNAS = (kind, path) => ({
  type: CHANGE_MC_FILE_SYSTEM_NAS,
  kind,
  path,
});

export const CHANGE_MC_HTTP_GET_NAS = 'CHANGE_MC_HTTP_GET_NAS';

export const changeMCHttpGetNAS = (port, path, scheme, host) => ({
  type: CHANGE_MC_HTTP_GET_NAS,
  port,
  path,
  scheme,
  host,
});

export const ADD_MC_HTTP_GET_HEADER_NAS = 'ADD_MC_HTTP_GET_HEADER_NAS';

export const addMCHttpGetHeaderNAS = () => ({
  type: ADD_MC_HTTP_GET_HEADER_NAS,
});

export const CHANGE_MC_HTTP_GET_HEADER_NAS = 'CHANGE_MC_HTTP_GET_HEADER_NAS';

export const changeMCHttpGetHeaderNAS = (fieldName, value, index) => ({
  type: CHANGE_MC_HTTP_GET_HEADER_NAS,
  fieldName,
  value,
  index,
});

export const DELETE_MC_HTTP_GET_HEADER_NAS = 'DELETE_MC_HTTP_GET_HEADER_NAS';

export const deleteMCHttpGetHeaderNAS = index => ({
  type: DELETE_MC_HTTP_GET_HEADER_NAS,
  index,
});

export const ADD_MC_METRICS_FORMAT_NAS = 'ADD_MC_METRICS_FORMAT_NAS';

export const addMCMetricsFormatNAS = () => ({
  type: ADD_MC_METRICS_FORMAT_NAS,
});

export const CHANGE_MC_METRIC_FORMAT_NAS = 'CHANGE_MC_METRIC_FORMAT_NAS';

export const changeMCMetricsFormatNAS = (format, index) => ({
  type: CHANGE_MC_METRIC_FORMAT_NAS,
  format,
  index,
});

export const DELETE_MC_METRIC_FORMAT_NAS = 'DELETE_MC_METRIC_FORMAT_NAS';

export const deleteMCMetricsFormatNAS = index => ({
  type: DELETE_MC_METRIC_FORMAT_NAS,
  index,
});

export const CHANGE_MC_CUSTOM_CONTAINER_NAS = 'CHANGE_MC_CUSTOM_CONTAINER_NAS';

export const changeMCCustomContainerNAS = yamlContainer => ({
  type: CHANGE_MC_CUSTOM_CONTAINER_NAS,
  yamlContainer,
});
