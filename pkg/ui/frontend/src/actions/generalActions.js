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