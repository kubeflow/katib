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

export const ADD_PARAMETER = "CHANGE_PARAMETER";

export const addParameter = () => ({
    type: ADD_PARAMETER,
})

export const EDIT_PARAMETER = "EDIT_PARAMTER";

export const editParameter = (index, field, value) => ({
    type: EDIT_PARAMETER,
    index, field, value,
})

export const DELETE_PARAMETER = "DELETE_PARAMETER";

export const deleteParameter = (index) => ({
    type: DELETE_PARAMETER,
    index,
})

export const ADD_LIST_PARAMETER = "ADD_LIST_PARAMETER";

export const addListParameter = (paramIndex) => ({
    type: ADD_LIST_PARAMETER,
    paramIndex,
})

export const EDIT_LIST_PARAMETER = "EDIT_LIST_PARAMETER";

export const editListParameter = (paramIndex, index, value) => ({
    type: EDIT_LIST_PARAMETER,
    paramIndex, index, value
})

export const DELETE_LIST_PARAMETER = "DELETE_LIST_PARAMETER";

export const deleteListParameter = (paramIndex, index) => ({
    type: DELETE_LIST_PARAMETER,
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