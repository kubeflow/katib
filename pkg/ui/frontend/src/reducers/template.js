import * as actions from '../actions/templateActions';

const initialState = {
    menuOpen: false,
    addOpen: false,
    editOpen: false,
    deleteOpen: false,
    workerTemplates: [
    ],
    collectorTemplates: [
    ],
    newTemplateName: '',
    newTemplateYaml: '',
    currentTemplateIndex: '',
    edittedTemplate: {
        name: '',
        yaml: '',
    },
    currentTemplateName: '',
};

const rootReducer = (state = initialState, action) => {
    switch (action.type) {
        case actions.CLOSE_DIALOG:
            return {
                ...state,
                editOpen: false,
                addOpen: false,
                deleteOpen: false,
            }
        case actions.OPEN_DIALOG:
            switch(action.dialogType) {
                case "delete":
                    switch(action.templateType) {
                        case "worker": 
                            return {
                                ...state,
                                deleteOpen: true,
                                currentTemplateIndex: action.index,
                                currentTemplateName: state.workerTemplates[action.index].name,
                            }
                        case "collector": 
                            return {
                                ...state,
                                deleteOpen: true,
                                currentTemplateIndex: action.index,
                                currentTemplateName: state.collectorTemplates[action.index].name,
                            }
                        default: 
                            return {
                                ...state,
                            }
                    }
                case "edit":
                    switch(action.templateType) {
                        case "worker": 
                            return {
                                ...state,
                                editOpen: true,
                                currentTemplateIndex: action.index,
                                edittedTemplate: state.workerTemplates[action.index],
                            }
                        case "collector": 
                            return {
                                ...state,
                                editOpen: true,
                                currentTemplateIndex: action.index,
                                edittedTemplate: state.collectorTemplates[action.index],
                            }
                        default: 
                            return {
                                ...state,
                            }
                    }
                case "add":
                    return {
                        ...state,
                        addOpen: true,
                    };
                default:
                    return state;
            }
        case actions.CHANGE_TEMPLATE:
            let edittedTemplate = state.edittedTemplate;
            edittedTemplate[action.field] = action.value;
            return {
                ...state,
                edittedTemplate: edittedTemplate,
            }
        case actions.FETCH_WORKER_TEMPLATES_SUCCESS:
            return {
                ...state,
                workerTemplates: action.templates,
            }
        // case actions.FETCH_WORKER_TEMPLATES_FAILURE:
        //     return {
        //         ...state,
        //         snac
        //     }
        case actions.FETCH_COLLECTOR_TEMPLATES_SUCCESS:
            return {
                ...state,
                collectorTemplates: action.templates,
            }
        case actions.ADD_TEMPLATE_SUCCESS:
        case actions.DELETE_TEMPLATE_SUCCESS:
        case actions.EDIT_TEMPLATE_SUCCESS:
            switch (action.templateType) {
                case "worker": 
                    return {
                        ...state,
                        addOpen: false,
                        deleteOpen: false,
                        editOpen: false,
                        workerTemplates: action.templates,
                    } 
                case "collector":
                    return {
                        ...state,
                        addOpen: false,
                        deleteOpen: false,
                        editOpen: false,
                        collectorTemplates: action.templates,
                    }
            }
        case actions.ADD_TEMPLATE_FAILURE:
        case actions.EDIT_TEMPLATE_FAILURE:
        case actions.DELETE_TEMPLATE_FAILURE:
            return {
                ...state,
                addOpen: false,
                deleteOpen: false,
                editOpen: false,
            }
        default:
            return state;
    }
};

export default rootReducer;