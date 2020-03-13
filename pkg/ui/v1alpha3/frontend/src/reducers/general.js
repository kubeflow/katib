import * as actions from '../actions/generalActions';
import * as nasCreateActions from '../actions/nasCreateActions';
import * as hpCreateActions from '../actions/hpCreateActions';
import * as hpMonitorActions from '../actions/hpMonitorActions';
import * as nasMonitorActions from '../actions/nasMonitorActions';
import * as templateActions from '../actions/templateActions';

const initialState = {
  menuOpen: false,
  snackOpen: false,
  snackText: '',
  deleteDialog: false,
  deleteId: '',
  namespaces: [],
  globalNamespace: '',
  experiment: {},
  dialogExperimentOpen: false,

  templateNamespace: '',
  templateConfigMapName: '',
  templateName: '',
  trialTemplatesList: [],
  currentTemplateConfigMapsList: [],
  currentTemplateNamesList: [],
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
        snackText: 'Successfully submitted',
      };
    case actions.SUBMIT_YAML_FAILURE:
      return {
        ...state,
        snackOpen: true,
        snackText: action.message,
      };
    case actions.DELETE_EXPERIMENT_FAILURE:
      return {
        ...state,
        deleteDialog: false,
        snackOpen: true,
        snackText: 'Whoops, something went wrong',
      };
    case actions.DELETE_EXPERIMENT_SUCCESS:
      return {
        ...state,
        deleteDialog: false,
        snackOpen: true,
        snackText: 'Successfully deleted. Press Update button',
      };
    case actions.OPEN_DELETE_EXPERIMENT_DIALOG:
      return {
        ...state,
        deleteDialog: true,
        deleteExperimentName: action.name,
        deleteExperimentNamespace: action.namespace,
      };
    case actions.CLOSE_DELETE_EXPERIMENT_DIALOG:
      return {
        ...state,
        deleteDialog: false,
      };
    case nasCreateActions.SUBMIT_NAS_JOB_REQUEST:
      return {
        ...state,
        loading: true,
      };
    case nasCreateActions.SUBMIT_NAS_JOB_SUCCESS:
      return {
        ...state,
        loading: false,
        snackOpen: true,
        snackText: 'Successfully submitted',
      };
    case nasCreateActions.SUBMIT_NAS_JOB_FAILURE:
      return {
        ...state,
        loading: false,
        snackOpen: true,
        snackText: action.message,
      };
    case hpCreateActions.SUBMIT_HP_JOB_REQUEST:
      return {
        ...state,
        loading: true,
      };
    case hpCreateActions.SUBMIT_HP_JOB_SUCCESS:
      return {
        ...state,
        loading: false,
        snackOpen: true,
        snackText: 'Successfully submitted',
      };
    case hpCreateActions.SUBMIT_HP_JOB_FAILURE:
      return {
        ...state,
        loading: false,
        snackOpen: true,
        snackText: action.message,
      };
    case actions.FETCH_NAMESPACES_SUCCESS:
      return {
        ...state,
        namespaces: action.namespaces,
      };
    case actions.CHANGE_GLOBAL_NAMESPACE:
      state.globalNamespace = action.globalNamespace;
      return {
        ...state,
        globalNamespace: action.globalNamespace,
      };
    case actions.FETCH_EXPERIMENT_SUCCESS:
      return {
        ...state,
        experiment: action.experiment,
        dialogExperimentOpen: true,
      };

    case actions.CLOSE_DIALOG_EXPERIMENT:
      return {
        ...state,
        dialogExperimentOpen: false,
      };
    case hpMonitorActions.FETCH_HP_JOB_INFO_REQUEST:
      return {
        ...state,
        dialogExperimentOpen: false,
      };
    case nasMonitorActions.FETCH_NAS_JOB_INFO_REQUEST:
      return {
        ...state,
        dialogExperimentOpen: false,
      };
    case templateActions.FETCH_TRIAL_TEMPLATES_SUCCESS:
      let templates = action.trialTemplatesList;

      let configMapNames = templates[0].ConfigMapsList.map(configMap => configMap.ConfigMapName);

      let templateNames = templates[0].ConfigMapsList[0].TemplatesList.map(
        template => template.Name,
      );

      return {
        ...state,
        trialTemplatesList: templates,
        templateNamespace: templates[0].Namespace,
        templateConfigMapName: templates[0].ConfigMapsList[0].ConfigMapName,
        templateName: templates[0].ConfigMapsList[0].TemplatesList[0].Name,
        currentTemplateConfigMapsList: configMapNames,
        currentTemplateNamesList: templateNames,
      };

    case actions.FILTER_TEMPLATES_EXPERIMENT:
      switch (action.trialConfigMapName) {
        // Case when we change namespace
        case '':
          // Get Namespace index
          let nsIndex = state.trialTemplatesList.findIndex(function(trialTemplate, i) {
            return trialTemplate.Namespace === action.trialNamespace;
          });

          // Get new ConifgMapNames List
          configMapNames = state.trialTemplatesList[nsIndex].ConfigMapsList.map(
            configMap => configMap.ConfigMapName,
          );

          // Get new Template Names List
          console.log('nsIndex: ', nsIndex);
          console.log(
            'templates[nsIndex].ConfigMapsList[0]: ',
            state.trialTemplatesList[nsIndex].ConfigMapsList[0],
          );
          templateNames = state.trialTemplatesList[nsIndex].ConfigMapsList[0].TemplatesList.map(
            template => template.Name,
          );

          return {
            ...state,
            templateNamespace: action.trialNamespace,
            templateConfigMapName: configMapNames[0],
            templateName: templateNames[0],
            currentTemplateConfigMapsList: configMapNames,
            currentTemplateNamesList: templateNames,
          };
        // Case when we change configMap
        default:
          // Get Namespace index
          nsIndex = state.trialTemplatesList.findIndex(function(trialTemplate, i) {
            return trialTemplate.Namespace === action.trialNamespace;
          });

          // Get ConfigMap index
          let cmIndex = state.trialTemplatesList[nsIndex].ConfigMapsList.findIndex(function(
            configMap,
            i,
          ) {
            return configMap.ConfigMapName === action.trialConfigMapName;
          });

          // Get new Template Names List
          templateNames = state.trialTemplatesList[nsIndex].ConfigMapsList[
            cmIndex
          ].TemplatesList.map(template => template.Name);

          return {
            ...state,
            templateNamespace: action.trialNamespace,
            templateConfigMapName: action.trialConfigMapName,
            templateName: templateNames[0],
            currentTemplateNamesList: templateNames,
          };
      }
    case actions.CHANGE_TEMPLATE_NAME:
      return {
        ...state,
        templateName: action.templateName,
      };
    default:
      return state;
  }
};

export default generalReducer;
