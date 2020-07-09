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
  namespaces: [],
  globalNamespace: '',

  experimentName: '',
  experimentNamespace: 'All namespaces',
  filterStatus: {
    Created: true,
    Running: true,
    Restarting: true,
    Succeeded: true,
    Failed: true,
  },
  experiments: [],
  filteredExperiments: [],

  experiment: {},
  dialogExperimentOpen: false,
  suggestion: {},
  dialogSuggestionOpen: false,

  trialTemplatesData: [],

  configMapNamespaceIndex: -1,
  configMapNameIndex: -1,
  configMapPathIndex: -1,

  trialParameters: [],

  mcKindsList: ['StdOut', 'File', 'TensorFlowEvent', 'PrometheusMetric', 'Custom', 'None'],
  mcFileSystemKindsList: ['No File System', 'File', 'Directory'],
  mcURISchemesList: ['HTTP', 'HTTPS'],
};

const templateParameterRegex = '\\{trialParameters\\..+?\\}';

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
        snackText: 'Successfully deleted',
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
    case actions.FETCH_SUGGESTION_SUCCESS:
      return {
        ...state,
        suggestion: action.suggestion,
        dialogSuggestionOpen: true,
      };
    case actions.CLOSE_DIALOG_SUGGESTION:
      return {
        ...state,
        dialogSuggestionOpen: false,
      };
    case hpMonitorActions.FETCH_HP_JOB_INFO_REQUEST:
      return {
        ...state,
        dialogExperimentOpen: false,
        dialogSuggestionOpen: false,
      };
    case nasMonitorActions.FETCH_NAS_JOB_INFO_REQUEST:
      return {
        ...state,
        dialogExperimentOpen: false,
        dialogSuggestionOpen: false,
      };
    case templateActions.FETCH_TRIAL_TEMPLATES_SUCCESS:
      var trialTemplatesData = action.trialTemplatesData;

      let configMapNamespaceIndex = -1;
      let configMapNameIndex = -1;
      let configMapPathIndex = -1;

      if (
        trialTemplatesData.length > 0 &&
        trialTemplatesData[0].ConfigMaps[0].Templates.length > 0
      ) {
        configMapNamespaceIndex = 0;
        configMapNameIndex = 0;
        configMapPathIndex = 0;
      }

      // Get Parameter names from ConfigMap for Trial parameters
      var yaml = trialTemplatesData[0].ConfigMaps[0].Templates[0].Yaml;
      var trialParameters = [];
      var trialParameterNames = [];

      var matchStr = [...yaml.matchAll(templateParameterRegex)];
      matchStr.forEach(param => {
        let newParameter = param[0].slice(param[0].indexOf('.') + 1, param[0].indexOf('}'));
        if (!trialParameterNames.includes(newParameter)) {
          trialParameterNames.push(newParameter);
          trialParameters.push({
            name: newParameter,
            reference: '',
            description: '',
          });
        }
      });

      return {
        ...state,
        trialTemplatesData: trialTemplatesData,
        configMapNamespaceIndex: configMapNamespaceIndex,
        configMapNameIndex: configMapNameIndex,
        configMapPathIndex: configMapPathIndex,
        trialParameters: trialParameters,
      };
    case actions.FILTER_TEMPLATES_EXPERIMENT:
      let newNamespaceIndex = 0;
      let newNameIndex = 0;
      let newPathIndex = 0;

      if (action.configMapNamespaceIndex !== state.configMapNamespaceIndex) {
        newNamespaceIndex = action.configMapNamespaceIndex;
      } else if (action.configMapNameIndex !== state.configMapNameIndex) {
        newNamespaceIndex = action.configMapNamespaceIndex;
        newNameIndex = action.configMapNameIndex;
      } else {
        newNamespaceIndex = action.configMapNamespaceIndex;
        newNameIndex = action.configMapNameIndex;
        newPathIndex = action.configMapPathIndex;
      }

      // Get Parameter names from ConfigMap for Trial parameters
      // Change only if any ConfigMap information has been changed
      trialParameters = state.trialParameters.slice();
      if (
        newNamespaceIndex !== state.configMapNamespaceIndex ||
        newNameIndex !== state.configMapNameIndex ||
        newPathIndex !== state.configMapPathIndex
      ) {
        trialTemplatesData = state.trialTemplatesData;
        yaml =
          trialTemplatesData[newNamespaceIndex].ConfigMaps[newNameIndex].Templates[newPathIndex]
            .Yaml;
        trialParameterNames = [];
        trialParameters = [];
        matchStr = [...yaml.matchAll(templateParameterRegex)];
        matchStr.forEach(param => {
          let newParameter = param[0].slice(param[0].indexOf('.') + 1, param[0].indexOf('}'));
          if (!trialParameterNames.includes(newParameter)) {
            trialParameterNames.push(newParameter);
            trialParameters.push({
              name: newParameter,
              reference: '',
              description: '',
            });
          }
        });
      }
      return {
        ...state,
        configMapNamespaceIndex: newNamespaceIndex,
        configMapNameIndex: newNameIndex,
        configMapPathIndex: newPathIndex,
        trialParameters: trialParameters,
      };
    case actions.VALIDATION_ERROR:
      return {
        ...state,
        snackOpen: true,
        snackText: action.message,
      };
    case actions.EDIT_TRIAL_PARAMETERS:
      let newParams = state.trialParameters.slice();
      newParams[action.index].name = action.name;
      newParams[action.index].reference = action.reference;
      newParams[action.index].description = action.description;
      return {
        ...state,
        trialParameters: newParams,
      };
    case actions.FETCH_EXPERIMENTS_SUCCESS:
      var experiments = action.experiments;

      var statuses = Object.assign({}, state.filterStatus);
      var statusKeys = Object.keys(statuses);
      var newFilterStatus = statusKeys.filter(key => {
        return statuses[key];
      });

      var filteredExperiments = experiments.filter(
        experiment =>
          newFilterStatus.includes(experiment.status) &&
          experiment.name.includes(state.experimentName) &&
          (experiment.namespace === state.experimentNamespace ||
            state.experimentNamespace === 'All namespaces'),
      );
      return {
        ...state,
        experiments: action.experiments,
        filteredExperiments: filteredExperiments,
      };
    case actions.FILTER_EXPERIMENTS:
      experiments = state.experiments.slice();

      statuses = Object.assign({}, state.filterStatus);
      statusKeys = Object.keys(statuses);
      newFilterStatus = statusKeys.filter(key => {
        return statuses[key];
      });

      filteredExperiments = experiments.filter(
        experiment =>
          newFilterStatus.includes(experiment.status) &&
          experiment.name.includes(action.experimentName) &&
          (experiment.namespace === action.experimentNamespace ||
            action.experimentNamespace === 'All namespaces'),
      );
      return {
        ...state,
        filteredExperiments: filteredExperiments,
        experimentName: action.experimentName,
        experimentNamespace: action.experimentNamespace,
      };
    case actions.CHANGE_STATUS:
      experiments = state.experiments.slice();

      statuses = Object.assign({}, state.filterStatus);
      statuses[action.filter] = action.checked;
      statusKeys = Object.keys(statuses);
      newFilterStatus = statusKeys.filter(key => {
        return statuses[key];
      });

      filteredExperiments = experiments.filter(
        experiment =>
          newFilterStatus.includes(experiment.status) &&
          experiment.name.includes(state.experimentName) &&
          (experiment.namespace === state.experimentNamespace ||
            state.experimentNamespace === 'All namespaces'),
      );
      return {
        ...state,
        filterStatus: statuses,
        filteredExperiments: filteredExperiments,
      };
    case templateActions.ADD_TEMPLATE_SUCCESS:
      return {
        ...state,
        snackOpen: true,
        snackText: 'Successfully added new Template',
      };
    case templateActions.DELETE_TEMPLATE_SUCCESS:
      return {
        ...state,
        snackOpen: true,
        snackText: 'Successfully deleted Template',
      };
    case templateActions.EDIT_TEMPLATE_SUCCESS:
      return {
        ...state,
        snackOpen: true,
        snackText: 'Successfully edited Template',
      };
    case templateActions.ADD_TEMPLATE_FAILURE:
      return {
        ...state,
        snackOpen: true,
        snackText: 'Add Template failed: ' + action.error,
      };
    case templateActions.EDIT_TEMPLATE_FAILURE:
      return {
        ...state,
        snackOpen: true,
        snackText: 'Edit Template failed: ' + action.error,
      };
    case templateActions.DELETE_TEMPLATE_FAILURE:
      return {
        ...state,
        snackOpen: true,
        snackText: 'Delete Template failed: ' + action.error,
      };
    default:
      return state;
  }
};

export default generalReducer;
