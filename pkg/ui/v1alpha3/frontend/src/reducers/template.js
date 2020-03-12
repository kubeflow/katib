import * as actions from '../actions/templateActions';

const initialState = {
  addOpen: false,
  editOpen: false,
  deleteOpen: false,
  trialTemplatesList: [],
  filteredTrialTemplatesList: [],
  currentTemplateName: '',
  edittedTemplateNamespace: '',
  edittedTemplateConfigMapName: '',
  edittedTemplateName: '',
  edittedTemplateYaml: '',
  loading: false,
  edittedTemplateConfigMapSelectList: [],
  filteredNamespace: 'All namespaces',
  filteredConfigMapName: '',
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
      let templates = state.trialTemplatesList;

      //Filter ConfigMap
      let filteredConfigMaps = [];
      for (let i = 0; i < templates.length; i++) {
        let configMapsList = [];
        for (let j = 0; j < templates[i].ConfigMapsList.length; j++) {
          if (templates[i].ConfigMapsList[j].ConfigMapName.includes(state.filteredConfigMapName)) {
            configMapsList.push(templates[i].ConfigMapsList[j]);
          }
        }
        if (configMapsList.length != 0) {
          let newNamespaceBlock = {};
          newNamespaceBlock.Namespace = templates[i].Namespace;
          newNamespaceBlock.ConfigMapsList = configMapsList;
          filteredConfigMaps.push(newNamespaceBlock);
        }
      }

      // Filter Namespace
      let filteredTemplates = filteredConfigMaps.filter(
        template =>
          template.Namespace == state.filteredNamespace ||
          state.filteredNamespace == 'All namespaces',
      );
      return {
        ...state,
        trialTemplatesList: action.trialTemplatesList,
        filteredTrialTemplatesList: filteredTemplates,
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
        filteredTrialTemplatesList: action.trialTemplatesList,
        filteredNamespace: 'All namespaces',
        filteredConfigMapName: '',
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
    case actions.FILTER_TEMPLATES:
      templates = state.trialTemplatesList;

      //Filter ConfigMap
      filteredConfigMaps = [];
      for (let i = 0; i < templates.length; i++) {
        let configMapsList = [];
        for (let j = 0; j < templates[i].ConfigMapsList.length; j++) {
          if (templates[i].ConfigMapsList[j].ConfigMapName.includes(action.filteredConfigMapName)) {
            configMapsList.push(templates[i].ConfigMapsList[j]);
          }
        }
        if (configMapsList.length != 0) {
          let newNamespaceBlock = {};
          newNamespaceBlock.Namespace = templates[i].Namespace;
          newNamespaceBlock.ConfigMapsList = configMapsList;
          filteredConfigMaps.push(newNamespaceBlock);
        }
      }

      //Filter Namespace
      filteredTemplates = filteredConfigMaps.filter(
        template =>
          template.Namespace == action.filteredNamespace ||
          action.filteredNamespace == 'All namespaces',
      );

      return {
        ...state,
        filteredNamespace: action.filteredNamespace,
        filteredConfigMapName: action.filteredConfigMapName,
        filteredTrialTemplatesList: filteredTemplates,
      };
    default:
      return state;
  }
};

export default rootReducer;
