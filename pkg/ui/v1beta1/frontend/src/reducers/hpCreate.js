import * as actions from '../actions/hpCreateActions';
import * as constants from '../constants/constants';

const initialState = {
  loading: false,
  commonParametersMetadata: [
    {
      name: 'Name',
      value: 'random-experiment',
      description: 'A name of an Experiment',
    },
    {
      name: 'Namespace',
      value: 'kubeflow',
      description: 'Namespace to deploy an Experiment',
    },
  ],
  commonParametersSpec: [
    {
      name: 'ParallelTrialCount',
      value: '3',
      description: 'How many Trials can be processed in parallel',
    },
    {
      name: 'MaxTrialCount',
      value: '12',
      description: 'Max completed Trials to mark Experiment as succeeded',
    },
    {
      name: 'MaxFailedTrialCount',
      value: '3',
      description: 'Max failed trials to mark Experiment as failed',
    },
    {
      name: 'ResumePolicy',
      value: 'LongRunning',
      description: 'Resume policy describes how the Experiment should be restarted',
    },
  ],
  allResumePolicyTypes: ['Never', 'LongRunning', 'FromVolume'],
  allObjectiveTypes: ['minimize', 'maximize'],
  objective: [
    {
      name: 'Type',
      value: 'maximize',
      description: 'Type of optimization',
    },
    {
      name: 'Goal',
      value: '0.99',
      description: 'Goal of optimization',
    },
    {
      name: 'ObjectiveMetricName',
      value: 'Validation-accuracy',
      description: 'Name for the objective metric',
    },
  ],
  additionalMetricNames: [
    {
      value: 'Train-accuracy',
    },
  ],
  algorithmName: 'random',
  allAlgorithms: ['grid', 'random', 'hyperband', 'bayesianoptimization', 'tpe', 'cmaes'],
  algorithmSettings: [],
  parameters: [
    {
      name: 'lr',
      parameterType: 'double',
      feasibleSpace: 'feasibleSpace',
      min: '0.01',
      max: '0.03',
      list: [],
    },
    {
      name: 'num-layers',
      parameterType: 'int',
      feasibleSpace: 'feasibleSpace',
      min: '2',
      max: '5',
      list: [],
    },
    {
      name: 'optimizer',
      parameterType: 'categorical',
      feasibleSpace: 'list',
      min: '',
      max: '',
      list: [
        {
          value: 'sgd',
        },
        {
          value: 'adam',
        },
        {
          value: 'ftrl',
        },
      ],
    },
  ],
  allParameterTypes: ['int', 'double', 'categorical'],
  currentYaml: '',
  mcSpec: {
    collector: {
      kind: 'StdOut',
    },
    source: {
      filter: {
        metricsFormat: [],
      },
    },
  },
  mcCustomContainerYaml: '',
};

const filterValue = (obj, key) => {
  return obj.findIndex(p => p.name === key);
};

const hpCreateReducer = (state = initialState, action) => {
  switch (action.type) {
    case actions.CHANGE_YAML_HP:
      return {
        ...state,
        currentYaml: action.payload,
      };
    case actions.CHANGE_META_HP:
      let meta = state.commonParametersMetadata.slice();
      var index = filterValue(meta, action.name);
      meta[index].value = action.value;
      return {
        ...state,
        commonParametersMetadata: meta,
      };
    case actions.CHANGE_SPEC_HP:
      let spec = state.commonParametersSpec.slice();
      index = filterValue(spec, action.name);
      spec[index].value = action.value;
      return {
        ...state,
        commonParametersSpec: spec,
      };
    case actions.CHANGE_OBJECTIVE_HP:
      let obj = state.objective.slice();
      index = filterValue(obj, action.name);
      obj[index].value = action.value;
      return {
        ...state,
        objective: obj,
      };
    case actions.ADD_METRICS_HP:
      var additionalMetricNames = state.additionalMetricNames.slice();
      additionalMetricNames.push({
        value: '',
      });
      return {
        ...state,
        additionalMetricNames: additionalMetricNames,
      };
    case actions.DELETE_METRICS_HP:
      additionalMetricNames = state.additionalMetricNames.slice();
      additionalMetricNames.splice(action.index, 1);
      return {
        ...state,
        additionalMetricNames: additionalMetricNames,
      };
    case actions.EDIT_METRICS_HP:
      additionalMetricNames = state.additionalMetricNames.slice();
      additionalMetricNames[action.index].value = action.value;
      return {
        ...state,
        additionalMetricNames: additionalMetricNames,
      };
    case actions.CHANGE_ALGORITHM_NAME_HP:
      return {
        ...state,
        algorithmName: action.algorithmName,
      };
    case actions.ADD_ALGORITHM_SETTING_HP:
      var algorithmSettings = state.algorithmSettings.slice();
      let setting = { name: '', value: '' };
      algorithmSettings.push(setting);
      return {
        ...state,
        algorithmSettings: algorithmSettings,
      };
    case actions.CHANGE_ALGORITHM_SETTING_HP:
      algorithmSettings = state.algorithmSettings.slice();
      algorithmSettings[action.index][action.field] = action.value;
      return {
        ...state,
        algorithmSettings: algorithmSettings,
      };
    case actions.DELETE_ALGORITHM_SETTING_HP:
      algorithmSettings = state.algorithmSettings.slice();
      algorithmSettings.splice(action.index, 1);
      return {
        ...state,
        algorithmSettings: algorithmSettings,
      };
    case actions.ADD_PARAMETER_HP:
      var parameters = state.parameters.slice();
      parameters.push({
        name: '',
        parameterType: '',
        feasibleSpace: 'feasibleSpace',
        min: '',
        max: '',
        list: [],
      });
      return {
        ...state,
        parameters: parameters,
      };
    case actions.EDIT_PARAMETER_HP:
      parameters = state.parameters.slice();
      parameters[action.index][action.field] = action.value;
      return {
        ...state,
        parameters: parameters,
      };
    case actions.DELETE_PARAMETER_HP:
      parameters = state.parameters.slice();
      parameters.splice(action.index, 1);
      return {
        ...state,
        parameters: parameters,
      };
    case actions.ADD_LIST_PARAMETER_HP:
      parameters = state.parameters.slice();
      parameters[action.paramIndex].list.push({
        value: '',
      });
      return {
        ...state,
        parameters: parameters,
      };
    case actions.EDIT_LIST_PARAMETER_HP:
      parameters = state.parameters.slice();
      parameters[action.paramIndex].list[action.index].value = action.value;
      return {
        ...state,
        parameters: parameters,
      };
    case actions.DELETE_LIST_PARAMETER_HP:
      parameters = state.parameters.slice();
      parameters[action.paramIndex].list.splice(action.index, 1);
      return {
        ...state,
        parameters: parameters,
      };
    // Metrics Collector Kind change
    case actions.CHANGE_MC_KIND_HP:
      var newMCSpec = JSON.parse(JSON.stringify(state.mcSpec));
      newMCSpec.collector.kind = action.kind;

      if (
        action.kind === constants.MC_KIND_FILE ||
        action.kind === constants.MC_KIND_TENSORFLOW_EVENT ||
        action.kind === constants.MC_KIND_CUSTOM
      ) {
        let newKind;
        switch (action.kind) {
          case constants.MC_KIND_FILE:
            newKind = constants.MC_FILE_SYSTEM_KIND_FILE;
            break;

          case constants.MC_KIND_TENSORFLOW_EVENT:
            newKind = constants.MC_FILE_SYSTEM_KIND_DIRECTORY;
            break;

          default:
            newKind = constants.MC_FILE_SYSTEM_NO_KIND;
        }
        // File or TF Event Kind
        newMCSpec.source.fileSystemPath = {
          kind: newKind,
          path: '',
        };
      } else if (action.kind === constants.MC_KIND_PROMETHEUS) {
        // Prometheus Kind
        newMCSpec.source.httpGet = {
          port: constants.MC_PROMETHEUS_DEFAULT_PORT,
          path: constants.MC_PROMETHEUS_DEFAULT_PATH,
          scheme: constants.MC_HTTP_GET_HTTP_SCHEME,
          host: '',
          httpHeaders: [],
        };
      }

      return {
        ...state,
        mcSpec: newMCSpec,
        mcCustomContainerYaml: '',
      };
    // File System Path change
    case actions.CHANGE_MC_FILE_SYSTEM_HP:
      newMCSpec = JSON.parse(JSON.stringify(state.mcSpec));
      newMCSpec.source.fileSystemPath.kind = action.kind;
      newMCSpec.source.fileSystemPath.path = action.path;
      return {
        ...state,
        mcSpec: newMCSpec,
      };
    // HTTPGet settings
    case actions.CHANGE_MC_HTTP_GET_HP:
      newMCSpec = JSON.parse(JSON.stringify(state.mcSpec));

      newMCSpec.source.httpGet.port = action.port;
      newMCSpec.source.httpGet.path = action.path;
      newMCSpec.source.httpGet.scheme = action.scheme;
      newMCSpec.source.httpGet.host = action.host;

      return {
        ...state,
        mcSpec: newMCSpec,
      };
    // Collector HTTPGet Headers
    case actions.ADD_MC_HTTP_GET_HEADER_HP:
      newMCSpec = JSON.parse(JSON.stringify(state.mcSpec));
      var currentHeaders = newMCSpec.source.httpGet.httpHeaders.slice();
      let newHeader = { name: '', value: '' };
      currentHeaders.push(newHeader);
      newMCSpec.source.httpGet.httpHeaders = currentHeaders;
      return {
        ...state,
        mcSpec: newMCSpec,
      };
    case actions.CHANGE_MC_HTTP_GET_HEADER_HP:
      newMCSpec = JSON.parse(JSON.stringify(state.mcSpec));
      currentHeaders = newMCSpec.source.httpGet.httpHeaders.slice();
      currentHeaders[action.index][action.fieldName] = action.value;
      newMCSpec.source.httpGet.httpHeaders = currentHeaders;
      return {
        ...state,
        mcSpec: newMCSpec,
      };
    case actions.DELETE_MC_HTTP_GET_HEADER_HP:
      newMCSpec = JSON.parse(JSON.stringify(state.mcSpec));
      currentHeaders = newMCSpec.source.httpGet.httpHeaders.slice();
      currentHeaders.splice(action.index, 1);
      newMCSpec.source.httpGet.httpHeaders = currentHeaders;
      return {
        ...state,
        mcSpec: newMCSpec,
      };
    // Collector Custom container
    case actions.CHANGE_MC_CUSTOM_CONTAINER_HP:
      return {
        ...state,
        mcCustomContainerYaml: action.yamlContainer,
      };
    // Collector Metrics Format
    case actions.ADD_MC_METRICS_FORMAT_HP:
      newMCSpec = JSON.parse(JSON.stringify(state.mcSpec));
      var currentFormats = newMCSpec.source.filter.metricsFormat.slice();
      currentFormats.push('');
      newMCSpec.source.filter.metricsFormat = currentFormats;
      return {
        ...state,
        mcSpec: newMCSpec,
      };
    case actions.CHANGE_MC_METRIC_FORMAT_HP:
      newMCSpec = JSON.parse(JSON.stringify(state.mcSpec));
      currentFormats = newMCSpec.source.filter.metricsFormat.slice();
      currentFormats[action.index] = action.format;
      newMCSpec.source.filter.metricsFormat = currentFormats;
      return {
        ...state,
        mcSpec: newMCSpec,
      };
    case actions.DELETE_MC_METRIC_FORMAT_HP:
      newMCSpec = JSON.parse(JSON.stringify(state.mcSpec));
      currentFormats = newMCSpec.source.filter.metricsFormat.slice();
      currentFormats.splice(action.index, 1);
      newMCSpec.source.filter.metricsFormat = currentFormats;
      return {
        ...state,
        mcSpec: newMCSpec,
      };
    default:
      return state;
  }
};

export default hpCreateReducer;
