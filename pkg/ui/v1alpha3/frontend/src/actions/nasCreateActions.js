export const CHANGE_YAML_NAS = "CHANGE_YAML_NAS";

export const changeYaml = (yaml) => ({
    type: CHANGE_YAML_NAS,
    payload: yaml
})

export const CHANGE_META_NAS = "CHANGE_META_NAS";

export const changeMeta = (name, value) => ({
    type: CHANGE_META_NAS,
    name, value,
})

export const CHANGE_SPEC_NAS = "CHANGE_SPEC_NAS";

export const changeSpec = (name, value) => ({
    type: CHANGE_SPEC_NAS,
    name, value
})

export const CHANGE_OBJECTIVE_NAS = "CHANGE_OBJECTIVE_NAS";

export const changeObjective = (name, value) => ({
    type: CHANGE_OBJECTIVE_NAS,
    name, value
})

export const ADD_METRICS_NAS = "ADD_METRICS_NAS";

export const addMetrics = () => ({
    type: ADD_METRICS_NAS,
})

export const DELETE_METRICS_NAS = "DELETE_METRICS_NAS";

export const deleteMetrics = (index) => ({
    type: DELETE_METRICS_NAS,
    index,
})

export const EDIT_METRICS_NAS = "EDIT_METRICS_NAS";

export const editMetrics = (index, value) => ({
    type: EDIT_METRICS_NAS,
    index, value,
})

export const CHANGE_ALGORITHM_NAME_NAS = "CHANGE_ALGORITHM_NAME_NAS";

export const changeAlgorithmName = (algorithmName) => ({
    type: CHANGE_ALGORITHM_NAME_NAS,
    algorithmName,
})

export const ADD_ALGORITHM_SETTING_NAS = "ADD_ALGORITHM_SETTING_NAS";

export const addAlgorithmSetting = () => ({
    type: ADD_ALGORITHM_SETTING_NAS,
})

export const CHANGE_ALGORITHM_SETTING_NAS = "CHANGE_ALGORITHM_SETTING_NAS";

export const changeAlgorithmSetting = (index, field, value) => ({
    type: CHANGE_ALGORITHM_SETTING_NAS,
    field, value, index
})

export const DELETE_ALGORITHM_SETTING_NAS = "DELETE_ALGORITHM_SETTING_NAS";

export const deleteAlgorithmSetting = (index) => ({
    type: DELETE_ALGORITHM_SETTING_NAS,
    index
})

export const EDIT_NUM_LAYERS = "EDIT_NUM_LAYERS"

export const editNumLayers = (value) => ({
    type: EDIT_NUM_LAYERS,
    value
})

export const ADD_SIZE = "ADD_SIZE";

export const addSize = (sizeType) => ({
    type: ADD_SIZE,
    sizeType,
})

export const EDIT_SIZE = "EDIT_SIZE";

export const editSize = (sizeType, index, value) => ({
    type: EDIT_SIZE,
    sizeType, index, value,
})

export const DELETE_SIZE = "DELETE_SIZE";

export const deleteSize = (sizeType, index) => ({
    type: DELETE_SIZE,
    sizeType, index,
})

export const ADD_OPERATION = "ADD_OPERATION";

export const addOperation = () => ({
    type: ADD_OPERATION,
})

export const DELETE_OPERATION = "DELETE_OPERATION";

export const deleteOperation = (index) => ({
    type: DELETE_OPERATION,
    index,
})

export const CHANGE_OPERATION = "CHANGE_OPERATION";

export const changeOperation = (index, value) => ({
    type: CHANGE_OPERATION,
    index, value,
})

export const ADD_PARAMETER_NAS = "ADD_PARAMETER_NAS";

export const addParameter = (opIndex) => ({
    type: ADD_PARAMETER_NAS,
    opIndex,
})

export const CHANGE_PARAMETER_NAS = "CHANGE_PARAMETER_NAS";

export const changeParameter = (opIndex, paramIndex, field, value) => ({
    type: CHANGE_PARAMETER_NAS,
    opIndex, paramIndex, field, value,
})

export const DELETE_PARAMETER_NAS = "DELETE_PARAMETER_NAS";

export const deleteParameter = (opIndex, paramIndex) => ({
    type: DELETE_PARAMETER_NAS,
    opIndex, paramIndex,
})


export const ADD_LIST_PARAMETER_NAS = "ADD_LIST_PARAMETER_NAS";

export const addListParameter = (opIndex, paramIndex) => ({
    type: ADD_LIST_PARAMETER_NAS,
    opIndex, paramIndex,
})

export const EDIT_LIST_PARAMETER_NAS = "EDIT_LIST_PARAMETER_NAS";

export const editListParameter = (opIndex, paramIndex, listIndex, value) => ({
    type: EDIT_LIST_PARAMETER_NAS,
    opIndex, paramIndex, listIndex, value,
})

export const DELETE_LIST_PARAMETER_NAS = "DELETE_LIST_PARAMETER_NAS";

export const deleteListParameter = (opIndex, paramIndex, listIndex) => ({
    type: DELETE_LIST_PARAMETER_NAS,
    opIndex, paramIndex, listIndex,
})

export const CHANGE_TRIAL_NAS = "CHANGE_TRIAL_NAS";

export const changeTrial = (trial) => ({
    type: CHANGE_TRIAL_NAS,
    trial,
})

export const CHANGE_TRIAL_NAMESPACE_NAS = "CHANGE_TRIAL_NAMESPACE_HP";

export const changeTrialNamespace = (namespace) => ({
    type: CHANGE_TRIAL_NAMESPACE_NAS,
    namespace,
})

export const SUBMIT_NAS_JOB_REQUEST = "SUBMIT_NAS_JOB_REQUEST";
export const SUBMIT_NAS_JOB_SUCCESS = "SUBMIT_NAS_JOB_SUCCESS";
export const SUBMIT_NAS_JOB_FAILURE = "SUBMIT_NAS_JOB_FAILURE";

export const submitNASJob = (data) => ({
    type: SUBMIT_NAS_JOB_REQUEST,
    data,
})

export const CLOSE_SNACKBAR = "CLOSE_SNACKBAR";

export const closeSnackbar = () => ({
    type: CLOSE_SNACKBAR,
})
