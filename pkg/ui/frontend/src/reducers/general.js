import * as actions from '../actions/generalActions';

const initialState = {
    menuOpen: false,
    snackOpen: false,
    snackText: "",
};

const generalReducer = (state = initialState, action) => {
    switch (action.type) {
        case actions.TOGGLE_MENU:
            return {
                ...state,
                menuOpen: action.state,
            };
        case actions.CLOSE_SNACKBAR:
            return {
                ...state,
                snackOpen: false,
            };
        case actions.SUBMIT_YAML_SUCCESS:
            return {
                ...state,
                snackOpen: true,
                snackText: "Successfully submitted",
            }
        case actions.SUBMIT_YAML_FAILURE:
            return {
                ...state,
                snackOpen: true,
                snackText: "Whoops, something went wrong",
            }
        default:
            return state;
    }
};

export default generalReducer;