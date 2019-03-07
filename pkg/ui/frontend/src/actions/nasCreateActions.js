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