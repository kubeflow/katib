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

export const ADD_PARAMETER = "ADD_PARAMETER";

export const addParameter = (opIndex) => ({
    type: ADD_PARAMETER,
    opIndex,
})

export const CHANGE_PARAMETER = "CHANGE_PARAMETER";

export const changeParameter = (opIndex, paramIndex, field, value) => ({
    type: CHANGE_PARAMETER,
    opIndex, paramIndex, field, value,
})

export const DELETE_PARAMETER = "DELETE_PARAMETER";

export const deleteParameter = (opIndex, paramIndex) => ({
    type: DELETE_PARAMETER,
    opIndex, paramIndex,
})

export const ADD_LIST_PARAMETER = "ADD_LIST_PARAMETER";

export const addListParameter = (opIndex, paramIndex) => ({
    type: ADD_LIST_PARAMETER,
    opIndex, paramIndex,
})

export const EDIT_LIST_PARAMETER = "EDIT_LIST_PARAMETER";

export const editListParameter = (opIndex, paramIndex, listIndex, value) => ({
    type: EDIT_LIST_PARAMETER,
    opIndex, paramIndex, listIndex, value,
})

export const DELETE_LIST_PARAMETER = "DELETE_LIST_PARAMETER";

export const deleteListParameter = (opIndex, paramIndex, listIndex) => ({
    type: DELETE_LIST_PARAMETER,
    opIndex, paramIndex, listIndex,
})

export const SUBMIT_NAS_JOB_REQUEST = "SUBMIT_NAS_JOB_REQUEST";
export const SUBMIT_NAS_JOB_SUCCESS = "SUBMIT_NAS_JOB_SUCCESS";
export const SUBMIT_NAS_JOB_FAILURE = "SUBMIT_NAS_JOB_FAILURE";

export const submitNASJob = (data) => ({
    type: SUBMIT_NAS_JOB_REQUEST,
    data,
})