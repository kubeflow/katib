export const CLOSE_DIALOG = "CLOSE_DIALOG";

export const closeDialog = (dialogType) => ({
    type: CLOSE_DIALOG,
    dialogType,
})

export const OPEN_DIALOG = "OPEN_DIALOG";

export const openDialog = (dialogType, index = -1, templateType = -1) => ({
    type: OPEN_DIALOG,
    dialogType, index, templateType
})

export const CHANGE_TEMPLATE = "CHANGE_TEMPLATE";

export const changeTemplate = (field, value) => ({
    type: CHANGE_TEMPLATE,
    field, value
})

export const FETCH_WORKER_TEMPLATES_REQUEST = "FETCH_WORKER_TEMPLATES_REQUEST"
export const FETCH_WORKER_TEMPLATES_SUCCESS = "FETCH_WORKER_TEMPLATES_SUCCESS"
export const FETCH_WORKER_TEMPLATES_FAILURE = "FETCH_WORKER_TEMPLATES_FAILURE"

export const fetchWorkerTemplates = () => ({
    type: FETCH_WORKER_TEMPLATES_REQUEST,
})

export const FETCH_COLLECTOR_TEMPLATES_REQUEST = "FETCH_COLLECTOR_TEMPLATES_REQUEST"
export const FETCH_COLLECTOR_TEMPLATES_SUCCESS = "FETCH_COLLECTOR_TEMPLATES_SUCCESS"
export const FETCH_COLLECTOR_TEMPLATES_FAILURE = "FETCH_COLLECTOR_TEMPLATES_FAILURE"

export const fetchCollectorTemplates = () => ({
    type: FETCH_COLLECTOR_TEMPLATES_REQUEST,
})

export const ADD_TEMPLATE_REQUEST = "ADD_TEMPLATE_REQUEST"
export const ADD_TEMPLATE_SUCCESS = "ADD_TEMPLATE_SUCCESS"
export const ADD_TEMPLATE_FAILURE = "ADD_TEMPLATE_FAILURE"

export const addTemplate = (name, yaml, kind) => ({
    type: ADD_TEMPLATE_REQUEST,
    name, yaml, kind
})

export const EDIT_TEMPLATE_REQUEST = "EDIT_TEMPLATE_REQUEST"
export const EDIT_TEMPLATE_SUCCESS = "EDIT_TEMPLATE_SUCCESS"
export const EDIT_TEMPLATE_FAILURE = "EDIT_TEMPLATE_FAILURE"

export const editTemplate = (name, yaml, kind) => ({
    type: EDIT_TEMPLATE_REQUEST,
    name, yaml, kind
})

export const DELETE_TEMPLATE_REQUEST = "DELETE_TEMPLATE_REQUEST"
export const DELETE_TEMPLATE_SUCCESS = "DELETE_TEMPLATE_SUCCESS"
export const DELETE_TEMPLATE_FAILURE = "DELETE_TEMPLATE_FAILURE"

export const deleteTemplate = (name, templateType) => ({
    type: DELETE_TEMPLATE_REQUEST,
    name, templateType
})

