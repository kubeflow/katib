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

export const deleteJob = (experimentName) => ({
    type: DELETE_JOB_REQUEST,
    experimentName,
})

export const OPEN_DELETE_JOB_DIALOG = "OPEN_DELETE_JOB_DIALOG";

export const openDeleteJobDialog = (experimentName) => ({
    type: OPEN_DELETE_JOB_DIALOG, 
    experimentName,
})

export const CLOSE_DELETE_JOB_DIALOG = "CLOSE_DELETE_JOB_DIALOG";

export const closeDeleteDialog = () => ({
    type: CLOSE_DELETE_JOB_DIALOG,
})
