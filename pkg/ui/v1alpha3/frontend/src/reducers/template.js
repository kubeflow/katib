import * as actions from '../actions/templateActions';

const initialState = {
  addOpen: false,
  editOpen: false,
  deleteOpen: false,
  //TODO: Delete it
  // trialTemplates: [],
  trialTemplatesList: [],
  currentTemplateName: '',
  edittedTemplateNamespace: '',
  edittedTemplateConfigMapName: '',
  edittedTemplateName: '',
  edittedTemplateYaml: '',
  loading: false,
  edittedTemplateConfigMapSelectList: [],
};

const rootReducer = (state = initialState, action) => {
  switch (action.type) {
    case actions.CLOSE_DIALOG:
      return {
        ...state,
        editOpen: false,
        addOpen: false,
        deleteOpen: false,
      };
    case actions.OPEN_DIALOG:
      switch (action.dialogType) {
        case 'add':
          return {
            ...state,
            addOpen: true,
            edittedTemplateNamespace: action.namespace,
            edittedTemplateConfigMapName: action.configMapName,
            edittedTemplateName: '',
            edittedTemplateYaml: '',
          };
        case 'edit':
          return {
            ...state,
            editOpen: true,
            edittedTemplateNamespace: action.namespace,
            edittedTemplateConfigMapName: action.configMapName,
            edittedTemplateName: action.templateName,
            edittedTemplateYaml: action.templateYaml,
            currentTemplateName: action.templateName,
          };
        case 'delete':
          return {
            ...state,
            deleteOpen: true,
            edittedTemplateNamespace: action.namespace,
            edittedTemplateConfigMapName: action.configMapName,
            edittedTemplateName: action.templateName,
          };
        default:
          return state;
      }
    case actions.FETCH_TRIAL_TEMPLATES_REQUEST:
      return {
        ...state,
        loading: true,
      };
    case actions.FETCH_TRIAL_TEMPLATES_FAILURE:
      return {
        ...state,
        loading: false,
      };
    case actions.FETCH_TRIAL_TEMPLATES_SUCCESS:
      return {
        ...state,
        trialTemplatesList: action.trialTemplatesList,
        loading: false,
      };
    case actions.ADD_TEMPLATE_SUCCESS:
    case actions.DELETE_TEMPLATE_SUCCESS:
    case actions.EDIT_TEMPLATE_SUCCESS:
      return {
        ...state,
        addOpen: false,
        deleteOpen: false,
        editOpen: false,
        trialTemplatesList: action.trialTemplatesList,
      };

    case actions.ADD_TEMPLATE_FAILURE:
    case actions.EDIT_TEMPLATE_FAILURE:
    case actions.DELETE_TEMPLATE_FAILURE:
      return {
        ...state,
        addOpen: false,
        deleteOpen: false,
        editOpen: false,
      };
    case actions.CHANGE_TEMPLATE:
      return {
        ...state,
        edittedTemplateNamespace: action.edittedTemplateNamespace,
        edittedTemplateConfigMapName: action.edittedTemplateConfigMapName,
        edittedTemplateName: action.edittedTemplateName,
        edittedTemplateYaml: action.edittedTemplateYaml,
        edittedTemplateConfigMapSelectList: action.edittedTemplateConfigMapSelectList,
      };
    default:
      return state;
  }
};

export default rootReducer;
