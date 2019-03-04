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

export const DELETE_TEMPLATE = "DELETE_TEMPLATE";

export const deleteTemplate = (templateType) => ({
    type: DELETE_TEMPLATE,
    templateType,
})
