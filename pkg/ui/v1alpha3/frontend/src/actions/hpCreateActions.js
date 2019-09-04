export const CHANGE_YAML_HP = "CHANGE_YAML_HP";

export const changeYaml = (yaml) => ({
    type: CHANGE_YAML_HP,
    payload: yaml
})

export const CHANGE_META_HP = "CHANGE_META_HP";

export const changeMeta = (name, value) => ({
    type: CHANGE_META_HP,
    name, value,
})

export const CHANGE_SPEC_HP = "CHANGE_SPEC_HP";

export const changeSpec = (name, value) => ({
    type: CHANGE_SPEC_HP,
    name, value
})

export const CHANGE_OBJECTIVE_HP = "CHANGE_OBJECTIVE_HP";

export const changeObjective = (name, value) => ({
    type: CHANGE_OBJECTIVE_HP,
    name, value
})
export const ADD_METRICS_HP = "ADD_METRICS_HP";

export const addMetrics = () => ({
    type: ADD_METRICS_HP,
})

export const DELETE_METRICS_HP = "DELETE_METRICS_HP";

export const deleteMetrics = (index) => ({
    type: DELETE_METRICS_HP,
    index,
})

export const EDIT_METRICS_HP = "EDIT_METRICS_HP";

export const editMetrics = (index, value) => ({
    type: EDIT_METRICS_HP,
    index, value,
})

export const CHANGE_ALGORITHM_NAME_HP = "CHANGE_ALGORITHM_NAME_HP";

export const changeAlgorithmName = (algorithmName) => ({
    type: CHANGE_ALGORITHM_NAME_HP,
    algorithmName,
})

export const ADD_ALGORITHM_SETTING_HP = "ADD_ALGORITHM_SETTING_HP";

export const addAlgorithmSetting = () => ({
    type: ADD_ALGORITHM_SETTING_HP,
})

export const CHANGE_ALGORITHM_SETTING_HP = "CHANGE_ALGORITHM_SETTING_HP";

export const changeAlgorithmSetting = (index, field, value) => ({
    type: CHANGE_ALGORITHM_SETTING_HP,
    field, value, index
})

export const DELETE_ALGORITHM_SETTING_HP = "DELETE_ALGORITHM_SETTING_HP";

export const deleteAlgorithmSetting = (index) => ({
    type: DELETE_ALGORITHM_SETTING_HP,
    index
})

export const ADD_PARAMETER_HP = "CHANGE_PARAMETER_HP";

export const addParameter = () => ({
    type: ADD_PARAMETER_HP,
})

export const EDIT_PARAMETER_HP = "EDIT_PARAMTER_HP";

export const editParameter = (index, field, value) => ({
    type: EDIT_PARAMETER_HP,
    index, field, value,
})

export const DELETE_PARAMETER_HP = "DELETE_PARAMETER_HP";

export const deleteParameter = (index) => ({
    type: DELETE_PARAMETER_HP,
    index,
})

export const ADD_LIST_PARAMETER_HP = "ADD_LIST_PARAMETER_HP";


export const addListParameter = (paramIndex) => ({
    type: ADD_LIST_PARAMETER_HP,
    paramIndex,
})

export const EDIT_LIST_PARAMETER_HP = "EDIT_LIST_PARAMETER_HP";

export const editListParameter = (paramIndex, index, value) => ({
    type: EDIT_LIST_PARAMETER_HP,
    paramIndex, index, value
})

export const DELETE_LIST_PARAMETER_HP = "DELETE_LIST_PARAMETER_HP";

export const deleteListParameter = (paramIndex, index) => ({
    type: DELETE_LIST_PARAMETER_HP,
    paramIndex, index
})

export const CHANGE_TRIAL_HP = "CHANGE_TRIAL_HP";

export const changeTrial = (trial) => ({
    type: CHANGE_TRIAL_HP,
    trial,
})

export const SUBMIT_HP_JOB_REQUEST = "SUBMIT_HP_JOB_REQUEST";
export const SUBMIT_HP_JOB_SUCCESS = "SUBMIT_HP_JOB_SUCCESS";
export const SUBMIT_HP_JOB_FAILURE = "SUBMIT_HP_JOB_FAILURE";

export const submitHPJob = (data) => ({
    type: SUBMIT_HP_JOB_REQUEST,
    data,
})
