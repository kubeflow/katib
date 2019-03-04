import * as actions from '../actions/templateActions';

const initialState = {
    menuOpen: false,
    addOpen: false,
    editOpen: false,
    deleteOpen: false,
    workerTemplates: [
        {
            name: "Worker Test 1",
            yaml: "ASDASDASDASDasdas djasnd asjdnj akjsdnajsknd jas dsajk nashknd kasnd askjnd aks dnask dnask dnsak j2quw jqoi qwi jna sjljnas jklaskln daklsjls aljkd asj a",
        },
        {
            name: "Worker Test 2",
            yaml: "ASDASDASDASDasdas djasnd asjdnj akjsdnajsknd jas dsajk nashknd kasnd askjnd aks dnask dnask dnsak j2quw jqoi qwi jna sjljnas jklaskln daklsjls aljkd asj a",
        },
    ],
    collectorTemplates: [
        {
            name: "Collector Test 1",
            yaml: "ASDASDASDASDasdas djasnd asjdnj akjsdnajsknd jas dsajk nashknd kasnd askjnd aks dnask dnask dnsak j2quw jqoi qwi jna sjljnas jklaskln daklsjls aljkd asj a",
        },
        {
            name: "Collector Test 2",
            yaml: "ASDASDASDASDasdas djasnd asjdnj akjsdnajsknd jas dsajk nashknd kasnd askjnd aks dnask dnask dnsak j2quw jqoi qwi jna sjljnas jklaskln daklsjls aljkd asj a",
        },
    ],
    newTemplateName: '',
    newTemplateYaml: '',
    currentTemplateIndex: '',
    edittedTemplate: {},
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
                    return {
                        ...state,
                        deleteOpen: true,
                        currentTemplateIndex: action.index,
                    };
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
            console.log(edittedTemplate)
            return {
                ...state,
                edittedTemplate: edittedTemplate,
            }
        case actions.DELETE_TEMPLATE:
            switch(action.templateType) {
                case "worker":
                    let workers = state.workerTemplates.slice();
                    workers.splice(action.index, 1);
                    return {
                        ...state,
                        workerTemplates: workers,
                    }
                case "collector":
                    let collectors = state.collectorTemplates.slice();
                    collectors.splice(action.index, 1);
                    return {
                        ...state,
                        collectorTemplates: collectors,
                        deleteOpen: false,
                    }
                default:
                    return {
                        ...state,
                    }
            }
        default:
            return state;
    }
};

export default rootReducer;