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
        default:
            return state;
    }
};

export default generalReducer;