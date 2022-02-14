import * as actions from '../actions/generalActions';
import * as nasCreateActions from '../actions/nasCreateActions';
import * as hpCreateActions from '../actions/hpCreateActions';
import * as hpMonitorActions from '../actions/hpMonitorActions';
import * as nasMonitorActions from '../actions/nasMonitorActions';
import * as templateActions from '../actions/templateActions';

import { TEMPLATE_SOURCE_CONFIG_MAP, TEMPLATE_SOURCE_YAML } from '../constants/constants';

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

  earlyStoppingAlgorithm: 'medianstop',
  allEarlyStoppingAlgorithms: ['medianstop'],
  earlyStoppingSettings: [],

  trialTemplateSourceList: [TEMPLATE_SOURCE_CONFIG_MAP, TEMPLATE_SOURCE_YAML],
  trialTemplateSource: 'ConfigMap',
  primaryPodLabels: [],
  trialTemplateSpec: [
    {
      name: 'PrimaryContainerName',
      value: 'training-container',
      description: 'Name of training container where actual model training is running',
    },
    {
      name: 'SuccessCondition',
      value: 'status.conditions.#(type=="Complete")#|#(status=="True")#',
      description: `Condition when Trial custom resource is succeeded.
      Default value for k8s BatchJob: status.conditions.#(type=="Complete")#|#(status=="True")#.
      Default value for Kubeflow Job (TFJob, PyTorchJob, XGBoostJob, MXJob, MPIJob): status.conditions.#(type=="Succeeded")#|#(status=="True")#.`,
    },
    {
      name: 'FailureCondition',
      value: 'status.conditions.#(type=="Failed")#|#(status=="True")#',
      description: `Condition when Trial custom resource is failed.
      Default value for k8s BatchJob: status.conditions.#(type=="Failed")#|#(status=="True")#.
      Default value for Kubeflow Job (TFJob, PyTorchJob, XGBoostJob, MXJob, MPIJob): status.conditions.#(type=="Failed")#|#(status=="True")#.`,
    },
    {
      name: 'Retain',
      value: 'false',
      description:
        'Retain indicates that Trial resources must be not cleanup. Default value: false',
    },
  ],
  trialTemplateYAML: '',
  trialTemplatesData: [],

  configMapNamespaceIndex: -1,
  configMapNameIndex: -1,
  configMapPathIndex: -1,

  trialParameters: [],

  mcKindsList: ['StdOut', 'File', 'TensorFlowEvent', 'PrometheusMetric', 'Custom', 'None'],
  mcFileSystemKindsList: ['No File System', 'File', 'Directory'],
  mcURISchemesList: ['HTTP', 'HTTPS'],
};

// filterValue finds index from array where name == key.
const filterValue = (obj, key) => {
  return obj.findIndex(p => p.name === key);
};

// getTrialParameters returns Trial parameters for the given YAML.
// It parses only Trial parameters names from the YAML.
const getTrialParameters = YAML => {
  const templateParameterRegex = '\\{trialParameters\\..+?\\}';

  let trialParameters = [];
  let trialParameterNames = [];
  let matchStr = [...YAML.matchAll(templateParameterRegex)];

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
  return trialParameters;
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
    // Experiment early stopping actions.
    case actions.CHANGE_EARLY_STOPPING_ALGORITHM:
      return {
        ...state,
        earlyStoppingAlgorithm: action.algorithmName,
      };
    case actions.ADD_EARLY_STOPPING_SETTING:
      var earlyStoppingSettings = state.earlyStoppingSettings.slice();
      let setting = { name: '', value: '' };
      earlyStoppingSettings.push(setting);
      return {
        ...state,
        earlyStoppingSettings: earlyStoppingSettings,
      };
    case actions.CHANGE_EARLY_STOPPING_SETTING:
      earlyStoppingSettings = state.earlyStoppingSettings.slice();
      earlyStoppingSettings[action.index][action.field] = action.value;
      return {
        ...state,
        earlyStoppingSettings: earlyStoppingSettings,
      };
    case actions.DELETE_EARLY_STOPPING_SETTING:
      earlyStoppingSettings = state.earlyStoppingSettings.slice();
      earlyStoppingSettings.splice(action.index, 1);
      return {
        ...state,
        earlyStoppingSettings: earlyStoppingSettings,
      };
    // Experiment Trial Template actions.
    case templateActions.FETCH_TRIAL_TEMPLATES_SUCCESS:
      var trialTemplatesData = action.trialTemplatesData;

      let configMapNamespaceIndex = -1;
      let configMapNameIndex = -1;
      let configMapPathIndex = -1;
      var trialParameters = [];

      if (
        trialTemplatesData.length > 0 &&
        trialTemplatesData[0].ConfigMaps[0].Templates.length > 0
      ) {
        configMapNamespaceIndex = 0;
        configMapNameIndex = 0;
        configMapPathIndex = 0;

        // Get Trial parameters names from the ConfigMap template YAML
        trialParameters = getTrialParameters(trialTemplatesData[0].ConfigMaps[0].Templates[0].Yaml);
      }
      return {
        ...state,
        trialTemplatesData: trialTemplatesData,
        configMapNamespaceIndex: configMapNamespaceIndex,
        configMapNameIndex: configMapNameIndex,
        configMapPathIndex: configMapPathIndex,
        trialParameters: trialParameters,
      };
    case actions.CHANGE_TRIAL_TEMPLATE_SOURCE:
      return {
        ...state,
        trialTemplateSource: action.source,
        trialTemplateYAML: '',
        trialParameters: [],
      };
    case actions.ADD_PRIMARY_POD_LABEL:
      var newLabels = state.primaryPodLabels.slice();
      newLabels.push({ key: '', value: '' });
      return {
        ...state,
        primaryPodLabels: newLabels,
      };
    case actions.CHANGE_PRIMARY_POD_LABEL:
      newLabels = state.primaryPodLabels.slice();
      newLabels[action.index][action.fieldName] = action.value;
      return {
        ...state,
        primaryPodLabels: newLabels,
      };
    case actions.DELETE_PRIMARY_POD_LABEL:
      newLabels = state.primaryPodLabels.slice();
      newLabels.splice(action.index, 1);
      return {
        ...state,
        primaryPodLabels: newLabels,
      };
    case actions.CHANGE_TRIAL_TEMPLATE_SPEC:
      let newTrialTemplateSpec = state.trialTemplateSpec.slice();
      let index = filterValue(newTrialTemplateSpec, action.name);
      newTrialTemplateSpec[index].value = action.value;
      return {
        ...state,
        trialTemplateSpec: newTrialTemplateSpec,
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

      // Get Parameter names from ConfigMap for Trial parameters.
      // Change only if any ConfigMap information has been changed.
      trialParameters = state.trialParameters.slice();
      if (
        newNamespaceIndex !== state.configMapNamespaceIndex ||
        newNameIndex !== state.configMapNameIndex ||
        newPathIndex !== state.configMapPathIndex
      ) {
        // Get Trial parameters from the YAML.
        let configMap = state.trialTemplatesData[newNamespaceIndex].ConfigMaps[newNameIndex];
        trialParameters = getTrialParameters(configMap.Templates[newPathIndex].Yaml);
      }
      return {
        ...state,
        configMapNamespaceIndex: newNamespaceIndex,
        configMapNameIndex: newNameIndex,
        configMapPathIndex: newPathIndex,
        trialParameters: trialParameters,
      };
    case actions.CHANGE_TRIAL_TEMPLATE_YAML:
      // Get Trial parameters from the YAML.
      trialParameters = getTrialParameters(action.templateYAML);
      return {
        ...state,
        trialTemplateYAML: action.templateYAML,
        trialParameters: trialParameters,
      };
    case actions.CHANGE_TRIAL_PARAMETERS:
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
    case actions.VALIDATION_ERROR:
      return {
        ...state,
        snackOpen: true,
        snackText: action.message,
      };
    default:
      return state;
  }
};

export default generalReducer;
