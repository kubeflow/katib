import * as actions from '../actions/generalActions';

const initialState = {
    menuOpen: false,
    snackOpen: false,
    snackText: "",
    deleteDialog: false,
    deleteId: '',
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
        case actions.DELETE_JOB_FAILURE:
            return {
                ...state,
                deleteDialog: false,
                snackOpen: true,
                snackText: "Whoops, something went wrong",
            }
        case actions.DELETE_JOB_SUCCESS:
            return {
                ...state,
                deleteDialog: false,
                snackOpen: true,
                snackText: "Successfully deleted. Press Update button",
            }
        case actions.OPEN_DELETE_DIALOG:
            return {
                ...state,
                deleteDialog: true,
                deleteId: action.id,
            }
        case actions.CLOSE_DELETE_DIALOG:
            return {
                ...state,
                deleteDialog: false,
            }
        default:
            return state;
    }
};

export default generalReducer;