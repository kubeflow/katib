import * as actions from '../actions/nasCreateActions';
import * as constants from '../constants/constants';

const initialState = {
  commonParametersMetadata: [
    {
      name: 'Name',
      value: 'enas-example',
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
      description: 'Max failed Trials to mark Experiment as failed',
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
      value: 'Validation-Accuracy',
      description: 'Name for the objective metric',
    },
  ],
  additionalMetricNames: [],
  metricStrategiesList: ['min', 'max', 'latest'],
  metricStrategies: [
    {
      name: 'Validation-Accuracy',
      strategy: 'max',
    },
  ],
  algorithmName: 'enas',
  allAlgorithms: ['enas', 'darts'],
  algorithmSettings: [
    {
      name: 'controller_hidden_size',
      value: '64',
    },
    {
      name: 'controller_temperature',
      value: '5',
    },
    {
      name: 'controller_tanh_const',
      value: '2.25',
    },
    {
      name: 'controller_entropy_weight',
      value: '1e-5',
    },
    {
      name: 'controller_baseline_decay',
      value: '0.999',
    },
    {
      name: 'controller_learning_rate',
      value: '5e-5',
    },
    {
      name: 'controller_skip_target',
      value: '0.4',
    },
    {
      name: 'controller_skip_weight',
      value: '0.8',
    },
    {
      name: 'controller_train_steps',
      value: '50',
    },
    {
      name: 'controller_log_every_steps',
      value: '10',
    },
  ],
  //Graph Config
  numLayers: '8',
  inputSize: ['32', '32', '3'],
  outputSize: ['10'],
  operations: [
    {
      operationType: 'convolution',
      parameters: [
        {
          name: 'filter_size',
          parameterType: 'categorical',
          feasibleSpace: 'list',
          min: '',
          max: '',
          step: '',
          list: [
            {
              value: '3',
            },
            {
              value: '5',
            },
            {
              value: '7',
            },
          ],
        },
        {
          name: 'num_filter',
          parameterType: 'categorical',
          feasibleSpace: 'list',
          min: '',
          max: '',
          step: '',
          list: [
            {
              value: '32',
            },
            {
              value: '48',
            },
            {
              value: '64',
            },
            {
              value: '96',
            },
            {
              value: '128',
            },
          ],
        },
        {
          name: 'stride',
          parameterType: 'categorical',
          feasibleSpace: 'list',
          min: '',
          max: '',
          step: '',
          list: [
            {
              value: '1',
            },
            {
              value: '2',
            },
          ],
        },
      ],
    },
    {
      operationType: 'separable_convolution',
      parameters: [
        {
          name: 'filter_size',
          parameterType: 'categorical',
          feasibleSpace: 'list',
          min: '',
          max: '',
          step: '',
          list: [
            {
              value: '3',
            },
            {
              value: '5',
            },
            {
              value: '7',
            },
          ],
        },
        {
          name: 'num_filter',
          parameterType: 'categorical',
          feasibleSpace: 'list',
          min: '',
          max: '',
          step: '',
          list: [
            {
              value: '32',
            },
            {
              value: '48',
            },
            {
              value: '64',
            },
            {
              value: '96',
            },
            {
              value: '128',
            },
          ],
        },
        {
          name: 'stride',
          parameterType: 'categorical',
          feasibleSpace: 'list',
          min: '',
          max: '',
          step: '',
          list: [
            {
              value: '1',
            },
            {
              value: '2',
            },
          ],
        },
        {
          name: 'depth_multiplier',
          parameterType: 'categorical',
          feasibleSpace: 'list',
          min: '',
          max: '',
          step: '',
          list: [
            {
              value: '1',
            },
            {
              value: '2',
            },
          ],
        },
      ],
    },
    {
      operationType: 'depthwise_convolution',
      parameters: [
        {
          name: 'filter_size',
          parameterType: 'categorical',
          feasibleSpace: 'list',
          min: '',
          max: '',
          step: '',
          list: [
            {
              value: '3',
            },
            {
              value: '5',
            },
            {
              value: '7',
            },
          ],
        },
        {
          name: 'stride',
          parameterType: 'categorical',
          feasibleSpace: 'list',
          min: '',
          max: '',
          step: '',
          list: [
            {
              value: '1',
            },
            {
              value: '2',
            },
          ],
        },
        {
          name: 'depth_multiplier',
          parameterType: 'categorical',
          feasibleSpace: 'list',
          min: '',
          max: '',
          step: '',
          list: [
            {
              value: '1',
            },
            {
              value: '2',
            },
          ],
        },
      ],
    },
    {
      operationType: 'reduction',
      parameters: [
        {
          name: 'reduction_type',
          parameterType: 'categorical',
          feasibleSpace: 'list',
          min: '',
          max: '',
          step: '',
          list: [
            {
              value: 'max_pooling',
            },
            {
              value: 'avg_pooling',
            },
          ],
        },
        {
          name: 'pool_size',
          parameterType: 'int',
          feasibleSpace: 'feasibleSpace',
          min: '2',
          max: '3',
          step: '1',
          list: [],
        },
      ],
    },
  ],
  allParameterTypes: ['int', 'double', 'categorical'],
  currentYaml: '',
  snackText: '',
  snackOpen: false,
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

// filterValue finds index from array where name == key.
const filterValue = (obj, key) => {
  return obj.findIndex(p => p.name === key);
};

// setMetricStrategies sets metric strategies from objective and additional metrics
const setMetricStrategies = (objective, additionalMetricNames) => {
  let metricStrategies = [];
  // Objective metric - 2 index.
  // Objective type can't be empty.
  if (objective[2].value.trim() !== '') {
    // Set strategy from objective type.
    // Strategy == Objective Type by default.
    let strategy;
    if (objective[0].value === 'minimize') {
      strategy = 'min';
    } else {
      strategy = 'max';
    }
    metricStrategies.push({
      name: objective[2].value,
      strategy: strategy,
    });

    // Add not empty additional metrics.
    additionalMetricNames.forEach(metric => {
      if (metric.trim() !== '') {
        metricStrategies.push({
          name: metric,
          strategy: strategy,
        });
      }
    });
  }
  return metricStrategies;
};

const nasCreateReducer = (state = initialState, action) => {
  switch (action.type) {
    case actions.CHANGE_YAML_NAS:
      return {
        ...state,
        currentYaml: action.payload,
      };
    case actions.CHANGE_META_NAS:
      let meta = state.commonParametersMetadata.slice();
      var index = filterValue(meta, action.name);
      meta[index].value = action.value;
      return {
        ...state,
        commonParametersMetadata: meta,
      };
    case actions.CHANGE_SPEC_NAS:
      let spec = state.commonParametersSpec.slice();
      index = filterValue(spec, action.name);
      spec[index].value = action.value;
      return {
        ...state,
        commonParametersSpec: spec,
      };
    case actions.CHANGE_OBJECTIVE_NAS:
      let newObjective = state.objective.slice();
      index = filterValue(newObjective, action.name);
      newObjective[index].value = action.value;
      // Set new metric strategies.
      var newMetricStrategies = setMetricStrategies(newObjective, state.additionalMetricNames);
      return {
        ...state,
        objective: newObjective,
        metricStrategies: newMetricStrategies,
      };
    case actions.ADD_METRICS_NAS:
      var additionalMetricNames = state.additionalMetricNames.slice();
      additionalMetricNames.push('');
      return {
        ...state,
        additionalMetricNames: additionalMetricNames,
      };
    case actions.DELETE_METRICS_NAS:
      additionalMetricNames = state.additionalMetricNames.slice();
      additionalMetricNames.splice(action.index, 1);
      // Set new metric strategies.
      newMetricStrategies = setMetricStrategies(state.objective, additionalMetricNames);
      return {
        ...state,
        additionalMetricNames: additionalMetricNames,
        metricStrategies: newMetricStrategies,
      };
    case actions.EDIT_METRICS_NAS:
      additionalMetricNames = state.additionalMetricNames.slice();
      additionalMetricNames[action.index] = action.value;
      // Set new metric strategies.
      newMetricStrategies = setMetricStrategies(state.objective, additionalMetricNames);
      return {
        ...state,
        additionalMetricNames: additionalMetricNames,
        metricStrategies: newMetricStrategies,
      };
    case actions.CHANGE_METRIC_STRATEGY_NAS:
      newMetricStrategies = state.metricStrategies.slice();
      newMetricStrategies[action.index].strategy = action.strategy;
      return {
        ...state,
        metricStrategies: newMetricStrategies,
      };
    case actions.CHANGE_ALGORITHM_NAME_NAS:
      return {
        ...state,
        algorithmName: action.algorithmName,
      };
    case actions.ADD_ALGORITHM_SETTING_NAS:
      var algorithmSettings = state.algorithmSettings.slice();
      let setting = { name: '', value: '' };
      algorithmSettings.push(setting);
      return {
        ...state,
        algorithmSettings: algorithmSettings,
      };
    case actions.CHANGE_ALGORITHM_SETTING_NAS:
      algorithmSettings = state.algorithmSettings.slice();
      algorithmSettings[action.index][action.field] = action.value;
      return {
        ...state,
        algorithmSettings: algorithmSettings,
      };
    case actions.DELETE_ALGORITHM_SETTING_NAS:
      algorithmSettings = state.algorithmSettings.slice();
      algorithmSettings.splice(action.index, 1);
      return {
        ...state,
        algorithmSettings: algorithmSettings,
      };
    case actions.EDIT_NUM_LAYERS:
      let numLayers = action.value;
      return {
        ...state,
        numLayers: numLayers,
      };
    case actions.ADD_SIZE:
      var size = state[action.sizeType].slice();
      size.push('0');
      return {
        ...state,
        [action.sizeType]: size,
      };
    case actions.EDIT_SIZE:
      size = state[action.sizeType].slice();
      size[action.index] = action.value;
      return {
        ...state,
        [action.sizeType]: size,
      };
    case actions.DELETE_SIZE:
      size = state[action.sizeType].slice();
      size.splice(action.index, 1);
      return {
        ...state,
        [action.sizeType]: size,
      };
    case actions.ADD_OPERATION:
      var operations = state.operations.slice();
      operations.push({
        operationType: '',
        parameters: [],
      });
      return {
        ...state,
        operations,
      };
    case actions.DELETE_OPERATION:
      operations = state.operations.slice();
      operations.splice(action.index, 1);
      return {
        ...state,
        operations,
      };
    case actions.CHANGE_OPERATION:
      operations = state.operations.slice();
      operations[action.index].operationType = action.value;
      return {
        ...state,
        operations,
      };
    case actions.ADD_PARAMETER_NAS:
      operations = state.operations.slice();
      operations[action.opIndex].parameters.push({
        name: '',
        parameterType: 'categorical',
        feasibleSpace: 'list',
        min: '',
        max: '',
        step: '',
        list: [],
      });
      return {
        ...state,
        operations,
      };
    case actions.CHANGE_PARAMETER_NAS:
      operations = state.operations.slice();
      operations[action.opIndex].parameters[action.paramIndex][action.field] = action.value;
      return {
        ...state,
        operations,
      };
    case actions.DELETE_PARAMETER_NAS:
      operations = state.operations.slice();
      operations[action.opIndex].parameters.splice(action.paramIndex, 1);
      return {
        ...state,
        operations,
      };
    case actions.ADD_LIST_PARAMETER_NAS:
      operations = state.operations.slice();
      operations[action.opIndex].parameters[action.paramIndex].list.push({
        //TODO: Remove it?
        // name: "",
        value: '',
      });
      return {
        ...state,
        operations,
      };
    case actions.DELETE_LIST_PARAMETER_NAS:
      operations = state.operations.slice();
      operations[action.opIndex].parameters[action.paramIndex].list.splice(action.listIndex, 1);
      return {
        ...state,
        operations,
      };
    case actions.EDIT_LIST_PARAMETER_NAS:
      operations = state.operations.slice();
      operations[action.opIndex].parameters[action.paramIndex].list[action.listIndex].value =
        action.value;
      return {
        ...state,
        operations,
      };
    case actions.CLOSE_SNACKBAR:
      return {
        ...state,
        snackOpen: false,
      };
    // Metrics Collector Kind change
    case actions.CHANGE_MC_KIND_NAS:
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
    case actions.CHANGE_MC_FILE_SYSTEM_NAS:
      newMCSpec = JSON.parse(JSON.stringify(state.mcSpec));
      newMCSpec.source.fileSystemPath.kind = action.kind;
      newMCSpec.source.fileSystemPath.path = action.path;
      return {
        ...state,
        mcSpec: newMCSpec,
      };
    // HTTPGet settings
    case actions.CHANGE_MC_HTTP_GET_NAS:
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
    case actions.ADD_MC_HTTP_GET_HEADER_NAS:
      newMCSpec = JSON.parse(JSON.stringify(state.mcSpec));
      var currentHeaders = newMCSpec.source.httpGet.httpHeaders.slice();
      let newHeader = { name: '', value: '' };
      currentHeaders.push(newHeader);
      newMCSpec.source.httpGet.httpHeaders = currentHeaders;
      return {
        ...state,
        mcSpec: newMCSpec,
      };
    case actions.CHANGE_MC_HTTP_GET_HEADER_NAS:
      newMCSpec = JSON.parse(JSON.stringify(state.mcSpec));
      currentHeaders = newMCSpec.source.httpGet.httpHeaders.slice();
      currentHeaders[action.index][action.fieldName] = action.value;
      newMCSpec.source.httpGet.httpHeaders = currentHeaders;
      return {
        ...state,
        mcSpec: newMCSpec,
      };
    case actions.DELETE_MC_HTTP_GET_HEADER_NAS:
      newMCSpec = JSON.parse(JSON.stringify(state.mcSpec));
      currentHeaders = newMCSpec.source.httpGet.httpHeaders.slice();
      currentHeaders.splice(action.index, 1);
      newMCSpec.source.httpGet.httpHeaders = currentHeaders;
      return {
        ...state,
        mcSpec: newMCSpec,
      };
    // Collector Custom container
    case actions.CHANGE_MC_CUSTOM_CONTAINER_NAS:
      return {
        ...state,
        mcCustomContainerYaml: action.yamlContainer,
      };
    // Collector Metrics Format
    case actions.ADD_MC_METRICS_FORMAT_NAS:
      newMCSpec = JSON.parse(JSON.stringify(state.mcSpec));
      var currentFormats = newMCSpec.source.filter.metricsFormat.slice();
      currentFormats.push('');
      newMCSpec.source.filter.metricsFormat = currentFormats;
      return {
        ...state,
        mcSpec: newMCSpec,
      };
    case actions.CHANGE_MC_METRIC_FORMAT_NAS:
      newMCSpec = JSON.parse(JSON.stringify(state.mcSpec));
      currentFormats = newMCSpec.source.filter.metricsFormat.slice();
      currentFormats[action.index] = action.format;
      newMCSpec.source.filter.metricsFormat = currentFormats;
      return {
        ...state,
        mcSpec: newMCSpec,
      };
    case actions.DELETE_MC_METRIC_FORMAT_NAS:
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

export default nasCreateReducer;
