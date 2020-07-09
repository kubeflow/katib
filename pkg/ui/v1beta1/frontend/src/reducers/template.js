import * as actions from '../actions/templateActions';

const initialState = {
  addOpen: false,
  editOpen: false,
  deleteOpen: false,
  trialTemplatesData: [],
  filteredTrialTemplatesData: [],

  updatedConfigMapNamespace: '',
  updatedConfigMapName: '',
  configMapPath: '',
  updatedConfigMapPath: '',
  updatedTemplateYaml: '',

  loading: false,

  filteredConfigMapNamespace: 'All namespaces',
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
            updatedConfigMapNamespace: action.configMapNamespace,
            updatedConfigMapName: action.configMapName,
            updatedConfigMapPath: '',
            updatedTemplateYaml: '',
          };
        case 'edit':
          return {
            ...state,
            editOpen: true,
            updatedConfigMapNamespace: action.configMapNamespace,
            updatedConfigMapName: action.configMapName,
            configMapPath: action.configMapPath,
            updatedConfigMapPath: action.configMapPath,
            updatedTemplateYaml: action.templateYaml,
          };
        case 'delete':
          return {
            ...state,
            deleteOpen: true,
            updatedConfigMapNamespace: action.configMapNamespace,
            updatedConfigMapName: action.configMapName,
            updatedConfigMapPath: action.configMapPath,
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
        trialTemplatesData: action.trialTemplatesData,
        filteredTrialTemplatesData: action.trialTemplatesData,
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
        trialTemplatesData: action.trialTemplatesData,
        filteredTrialTemplatesData: action.trialTemplatesData,
        filteredConfigMapNamespace: 'All namespaces',
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
        filteredConfigMapNamespace: 'All namespaces',
        filteredConfigMapName: '',
      };
    case actions.CHANGE_TEMPLATE:
      return {
        ...state,
        updatedConfigMapNamespace: action.updatedConfigMapNamespace,
        updatedConfigMapName: action.updatedConfigMapName,
        updatedConfigMapPath: action.updatedConfigMapPath,
        updatedTemplateYaml: action.updatedTemplateYaml,
      };
    case actions.FILTER_TEMPLATES:
      //Filter ConfigMap Name
      let templatesData = state.trialTemplatesData;
      let filteredConfigMaps = [];
      for (let i = 0; i < templatesData.length; i++) {
        let configMaps = [];
        for (let j = 0; j < templatesData[i].ConfigMaps.length; j++) {
          if (templatesData[i].ConfigMaps[j].ConfigMapName.includes(action.filteredConfigMapName)) {
            configMaps.push(templatesData[i].ConfigMaps[j]);
          }
        }
        if (configMaps.length !== 0) {
          let newNamespaceBlock = {};
          newNamespaceBlock.ConfigMapNamespace = templatesData[i].ConfigMapNamespace;
          newNamespaceBlock.ConfigMaps = configMaps;
          filteredConfigMaps.push(newNamespaceBlock);
        }
      }

      //Filter Namespace
      let filteredTrialTemplatesData = filteredConfigMaps.filter(
        template =>
          template.ConfigMapNamespace === action.filteredConfigMapNamespace ||
          action.filteredConfigMapNamespace === 'All namespaces',
      );

      return {
        ...state,
        filteredConfigMapNamespace: action.filteredConfigMapNamespace,
        filteredConfigMapName: action.filteredConfigMapName,
        filteredTrialTemplatesData: filteredTrialTemplatesData,
      };
    default:
      return state;
  }
};

export default rootReducer;
