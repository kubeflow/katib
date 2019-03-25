export const TOGGLE_MENU = "TOGGLE_MENU";

export const toggleMenu = (state) => {
    return {
        type: TOGGLE_MENU,
        state,
    };
};

export const CLOSE_SNACKBAR = "CLOSE_SNACKBAR";

export const closeSnackbar = () => {
    return {
        type: CLOSE_SNACKBAR,
    };
};

export const SUBMIT_YAML_REQUEST = "SUBMIT_YAML_REQUEST";
export const SUBMIT_YAML_FAILURE = "SUBMIT_YAML_FAILURE";
export const SUBMIT_YAML_SUCCESS = "SUBMIT_YAML_SUCCESS";

export const submitYaml = (yaml) => ({
    type: SUBMIT_YAML_REQUEST,
    yaml,
})

export const DELETE_JOB_REQUEST = "DELETE_JOB_REQUEST";
export const DELETE_JOB_FAILURE = "DELETE_JOB_FAILURE";
export const DELETE_JOB_SUCCESS = "DELETE_JOB_SUCCESS";

export const deleteJob = (id) => ({
    type: DELETE_JOB_REQUEST,
    id,
})

export const OPEN_DELETE_DIALOG = "OPEN_DELETE_DIALOG";

export const openDeleteDialog = (id) => ({
    type: OPEN_DELETE_DIALOG, 
    id,
})

export const CLOSE_DELETE_DIALOG = "CLOSE_DELETE_DIALOG";

export const closeDeleteDialog = () => ({
    type: CLOSE_DELETE_DIALOG,
})