export const CHANGE_YAML = "CHANGE_YAML";

export const changeYaml = (yaml) => ({
    type: CHANGE_YAML,
    payload: yaml
})

export const CHANGE_META = "CHANGE_META";

export const changeMeta = (name, value) => ({
    type: CHANGE_META,
    name, value,
})

export const CHANGE_SPEC = "CHANGE_SPEC";

export const changeSpec = (name, value) => ({
    type: CHANGE_SPEC,
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

export const CHANGE_WORKER = "CHANGE_WORKER";

export const changeWorker = (worker) => ({
    type: CHANGE_WORKER,
    worker,
})

export const CHANGE_ALGORITHM = "CHANGE_ALGORITHM";

export const changeAlgorithm = (algorithm) => ({
    type: CHANGE_ALGORITHM,
    algorithm,
})

export const CHANGE_REQUEST_NUMBER = "CHANGE_REQUEST_NUMBER";

export const changeRequestNumber = (number) => ({
    type: CHANGE_REQUEST_NUMBER,
    number,
})

export const ADD_SUGGESTION_PARAMETER = "ADD_SUGGESTION_PARAMETER";

export const addSuggestionParameter = () => ({
    type: ADD_SUGGESTION_PARAMETER,
})

export const CHANGE_SUGGESTION_PARAMETER = "CHANGE_SUGGESTION_PARAMETER";

export const changeSuggestionParameter = (index, field, value) => ({
    type: CHANGE_SUGGESTION_PARAMETER,
    field, value, index
})

export const DELETE_SUGGESTION_PARAMETER = "DELETE_SUGGESTION_PARAMETER";

export const deleteSuggestionParameter = (index) => ({
    type: DELETE_SUGGESTION_PARAMETER,
    index
})

export const SUBMIT_HP_JOB_REQUEST = "SUBMIT_HP_JOB_REQUEST";
export const SUBMIT_HP_JOB_SUCCESS = "SUBMIT_HP_JOB_SUCCESS";
export const SUBMIT_HP_JOB_FAILURE = "SUBMIT_HP_JOB_FAILURE";

export const submitHPJob = (data) => ({
    type: SUBMIT_HP_JOB_REQUEST,
    data,
})
