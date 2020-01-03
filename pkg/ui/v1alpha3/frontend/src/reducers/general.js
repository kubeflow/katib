import * as actions from '../actions/generalActions';
import * as nasCreateActions from '../actions/nasCreateActions';
import * as hpCreateActions from '../actions/hpCreateActions';

const initialState = {
    menuOpen: false,
    snackOpen: false,
    snackText: "",
    deleteDialog: false,
    deleteId: '',
    namespaces: [
    ],
    globalNamespace: ""
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
                snackText: action.message,
            }
        case actions.DELETE_EXPERIMENT_FAILURE:
            return {
                ...state,
                deleteDialog: false,
                snackOpen: true,
                snackText: "Whoops, something went wrong",
            }
        case actions.DELETE_EXPERIMENT_SUCCESS:
            return {
                ...state,
                deleteDialog: false,
                snackOpen: true,
                snackText: "Successfully deleted. Press Update button",
            }
        case actions.OPEN_DELETE_EXPERIMENT_DIALOG:
            return {
                ...state,
                deleteDialog: true,
                deleteExperimentName: action.name,
                deleteExperimentNamespace: action.namespace,
            }
        case actions.CLOSE_DELETE_EXPERIMENT_DIALOG:
            return {
                ...state,
                deleteDialog: false,
            }
        case nasCreateActions.SUBMIT_NAS_JOB_REQUEST:
            return {
                ...state,
                loading: true,
            }
        case nasCreateActions.SUBMIT_NAS_JOB_SUCCESS:
            return {
                ...state,
                loading: false,
                snackOpen: true,
                snackText: "Successfully submitted",
            }
        case nasCreateActions.SUBMIT_NAS_JOB_FAILURE:
            return {
                ...state,
                loading: false,
                snackOpen: true,
                snackText: action.message,
            }
        case hpCreateActions.SUBMIT_HP_JOB_REQUEST:
            return {
                ...state,
                loading: true,
            }
        case hpCreateActions.SUBMIT_HP_JOB_SUCCESS:
            return {
                ...state,
                loading: false,
                snackOpen: true,
                snackText: "Successfully submitted",
            }
        case hpCreateActions.SUBMIT_HP_JOB_FAILURE:
            return {
                ...state,
                loading: false,
                snackOpen: true,
                snackText: action.message,
            }
        case actions.FETCH_NAMESPACES_SUCCESS:
            return {
                ...state,
                namespaces: action.namespaces
            }
        case actions.CHANGE_GLOBAL_NAMESPACE:
            state.globalNamespace = action.globalNamespace
            return {
                ...state,
                globalNamespace: action.globalNamespace
            }
        default:
            return state;
    }
};

export default generalReducer;
