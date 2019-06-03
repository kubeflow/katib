import * as actions from '../actions/nasCreateActions';

const initialState = {
    commonParametersMetadata: [
        {
            name: "Name",
            value: "nasrl-example",
            description: "A name of an experiment"
        },
        {
            name: "Namespace",
            value: "kubeflow",
            description: "Namespace to deploy an experiment"
        }
    ],
    commonParametersSpec: [
        {
            name: "ParallelTrialCount",
            value: "3",
            description: "How many trials can be processed in parallel"
        },
        {
            name: "MaxTrialCount",
            value: "12",
            description: "Max completed trials to mark experiment as succeeded"
        },
        {
            name: "MaxFailedTrialCount",
            value: "3",
            description: "Max failed trials to mark experiment as failed"
        }
    ],
    objective: [
        {
            name: "Type",
            value: "maximize",
            description: "Type of optimization"
        },
        {
            name: "Goal",
            value: "0.99",
            description: "Goal of optimization"
        },
        {
            name: "ObjectiveMetricName",
            value: "Validation-Accuracy",
            description: "Name for the objective metric"
        }
    ],
    additionalMetricNames: [],
    algorithmName: "nasrl",
    allAlgorithms: ["nasrl", "envelopenet"],
    algorithmSettings: [
        {
            name: "lstm_num_cells",
            value: "64",
        },
        {
            name: "lstm_num_layers",
            value: "1",
        },
        {
            name: "lstm_keep_prob",
            value: "1.0",
        },
        {
            name: "optimizer",
            value: "adam",
        },
        {
            name: "init_learning_rate",
            value: "1e-3",
        },
        {
            name: "lr_decay_start",
            value: "0",
        },
        {
            name: "lr_decay_every",
            value: "1000",
        },
        {
            name: "lr_decay_rate",
            value: "0.9",
        },
        {
            name: "skip-target",
            value: "0.4",
        },
        {
            name: "skip-weight",
            value: "0.8",
        },
        {
            name: "l2_reg",
            value: "0",
        },
        {
            name: "entropy_weight",
            value: "1e-4",
        },
        {
            name: "baseline_decay",
            value: "0.9999",
        }
    ],
    //Graph Config
    numLayers: '1',
    inputSize: ['32', '32', '3'],
    outputSize: ['10'],
    operations: [
        {
            operationType: "convolution",
            parameters: [
                {
                    name: "filter_size",
                    parameterType: "categorical",                  
                    feasibleSpace: "list",
                    min: "",
                    max: "",
                    step: "",
                    list: [
                        {
                            value: "3",
                        },
                        {
                            value: "5",
                        },
                        {
                            value: "7",
                        }
                    ],
                },
                {
                    name: "num_filter",
                    parameterType: "categorical",
                    feasibleSpace: "list",
                    min: "",
                    max: "",
                    step: "",
                    list: [
                        {
                            value: "32",
                        },
                        {
                            value: "48",
                        },
                        {
                            value: "64",
                        },
                        {
                            value: "96",
                        },
                        {
                            value: "128",
                        }
                    ],
                },
                {
                    name: "stride",
                    parameterType: "categorical",
                    feasibleSpace: "list",
                    min: "",
                    max: "",
                    step: "",
                    list: [
                        {
                            value: "1",
                        },
                        {
                            value: "2",
                        }
                    ],
                },
            ]
        },
        {
            operationType: "separable_convolution",
            parameters: [
                {
                    name: "filter_size",
                    parameterType: "categorical",                  
                    feasibleSpace: "list",
                    min: "",
                    max: "",
                    step: "",
                    list: [
                        {
                            value: "3",
                        },
                        {
                            value: "5",
                        },
                        {
                            value: "7",
                        },
                    ],
                },
                {
                    name: "num_filter",
                    parameterType: "categorical",
                    feasibleSpace: "list",
                    min: "",
                    max: "",
                    step: "",
                    list: [
                        {
                            value: "32",
                        },
                        {
                            value: "48",
                        },
                        {
                            value: "64",
                        },
                        {
                            value: "96",
                        },
                        {
                            value: "128",
                        },
                    ],
                },
                {
                    name: "stride",
                    parameterType: "categorical",
                    feasibleSpace: "list",
                    min: "",
                    max: "",
                    step: "",
                    list: [
                        {
                            value: "1",
                        },
                        {
                            value: "2",
                        },
                    ],
                },
                {
                    name: "depth_multiplier",
                    parameterType: "categorical",
                    feasibleSpace: "list",
                    min: "",
                    max: "",
                    step: "",
                    list: [
                        {
                            value: "1",
                        },
                        {
                            value: "2",
                        },
                    ],
                },
            ],
        },
        {
            operationType: "depthwise_convolution",
            parameters: [
                {
                    name: "filter_size",
                    parameterType: "categorical",                  
                    feasibleSpace: "list",
                    min: "",
                    max: "",
                    step: "",
                    list: [
                        {
                            value: "3",
                        },
                        {
                            value: "5",
                        },
                        {
                            value: "7",
                        },
                    ],
                },
                {
                    name: "stride",
                    parameterType: "categorical",
                    feasibleSpace: "list",
                    min: "",
                    max: "",
                    step: "",
                    list: [
                        {
                            value: "1",
                        },
                        {
                            value: "2",
                        },
                    ],
                },
                {
                    name: "depth_multiplier",
                    parameterType: "categorical",
                    feasibleSpace: "list",
                    min: "",
                    max: "",
                    step: "",
                    list: [
                        {
                            value: "1",
                        },
                        {
                            value: "2",
                        },
                    ],
                },
            ],
        },
        {
            operationType: "reduction",
            parameters: [
                {
                    name: "reduction_type",
                    parameterType: "categorical",                  
                    feasibleSpace: "list",
                    min: "",
                    max: "",
                    step: "",
                    list: [
                        {
                            value: "max_pooling",
                        },
                        {
                            value: "avg_pooling",
                        },
                    ],
                },
                {
                    name: "pool_size",
                    parameterType: "int",
                    feasibleSpace: "feasibleSpace",
                    min: "2",
                    max: "3",
                    step: "1",
                    list: [],
                },
            ],
        },
    ],
    allParameterTypes: ["int", "double", "categorical"],
    trial: 'nasRLTrialTemplate.yaml',
    currentYaml: '',
    snackText: '',
    snackOpen: false,
};

const filterValue = (obj, key) => {
    return obj.findIndex(p => p.name === key)
};

const nasCreateReducer = (state = initialState, action) => {
    switch (action.type) {
        case actions.CHANGE_YAML_NAS:
            return {
                ...state,
                currentYaml: action.payload,
            }
        case actions.CHANGE_META_NAS:
            let meta = state.commonParametersMetadata.slice();
            let index = filterValue(meta, action.name);
            meta[index].value = action.value;
            return {
                ...state,
                commonParametersMetadata: meta,
            }
        case actions.CHANGE_SPEC_NAS:
            let spec = state.commonParametersSpec.slice();
            index = filterValue(spec, action.name);
            spec[index].value = action.value;
            return {
                ...state,
                commonParametersSpec: spec,
            }
        case actions.CHANGE_OBJECTIVE_NAS:
            let obj = state.objective.slice();
            index = filterValue(obj, action.name);
            obj[index].value = action.value;
            return {
                ...state,
                objective: obj,
            }
        case actions.ADD_METRICS_NAS:
            let additionalMetricNames = state.additionalMetricNames.slice()
            additionalMetricNames.push({
                value: "",
            })
            return {
                ...state,
                additionalMetricNames: additionalMetricNames,
            }
        case actions.DELETE_METRICS_NAS:
            additionalMetricNames = state.additionalMetricNames.slice()
            additionalMetricNames.splice(action.index, 1)
            return {
                ...state,
                additionalMetricNames: additionalMetricNames,
            }
        case actions.EDIT_METRICS_NAS:
            additionalMetricNames = state.additionalMetricNames.slice()
            additionalMetricNames[action.index].value = action.value
            return {
                ...state,
                additionalMetricNames: additionalMetricNames,
            }
        case actions.CHANGE_ALGORITHM_NAME_NAS:
            return {
                ...state, 
                algorithmName: action.algorithmName,
            }
        case actions.ADD_ALGORITHM_SETTING_NAS:
            let algorithmSettings = state.algorithmSettings.slice();
            let setting = {name: "", value: ""};
            algorithmSettings.push(setting);
            return {
                ...state,
                algorithmSettings: algorithmSettings,
            }
        case actions.CHANGE_ALGORITHM_SETTING_NAS:
            algorithmSettings = state.algorithmSettings.slice();
            algorithmSettings[action.index][action.field] = action.value;
            return {
                ...state,
                algorithmSettings: algorithmSettings,
            }
        case actions.DELETE_ALGORITHM_SETTING_NAS:
            algorithmSettings = state.algorithmSettings.slice();
            algorithmSettings.splice(action.index, 1);
            return {
                ...state,
                algorithmSettings: algorithmSettings,
            }
        case actions.EDIT_NUM_LAYERS:
            let numLayers = action.value;
            return {
                ...state,
                numLayers: numLayers
            }
        case actions.ADD_SIZE:
            let size = state[action.sizeType].slice();
            size.push("0")
            return {
                ...state,
                [action.sizeType]: size,
            }
        case actions.EDIT_SIZE:
            size = state[action.sizeType].slice();
            size[action.index] = action.value;
            return {
                ...state,
                [action.sizeType]: size,
            }
        case actions.DELETE_SIZE:
            size = state[action.sizeType].slice();
            size.splice(action.index, 1);
            return {
                ...state,
                [action.sizeType]: size,
            }
        case actions.ADD_OPERATION:
            let operations = state.operations.slice();
            operations.push({
                operationType: "",
                parameters: [],
            });
            return {
                ...state,
                operations,
            }
        case actions.DELETE_OPERATION:
            operations = state.operations.slice();
            operations.splice(action.index, 1);
            return {
                ...state,
                operations,
            }
        case actions.CHANGE_OPERATION:
            operations = state.operations.slice();
            operations[action.index].operationType = action.value;
            return {
                ...state,
                operations,
            }
        case actions.ADD_PARAMETER_NAS:
            operations = state.operations.slice();
            operations[action.opIndex].parameters.push(
               {
                name: "",
                parameterType: "categorical",
                feasibleSpace: "list",
                min: "",
                max: "",
                step: "",
                list: [
                ],
               }
            )
            return {
                ...state,
                operations,
            }
        case actions.CHANGE_PARAMETER_NAS:
            operations = state.operations.slice();
            operations[action.opIndex].parameters[action.paramIndex][action.field] = action.value;
            return {
                ...state,
                operations,
            }
        case actions.DELETE_PARAMETER_NAS:
            operations = state.operations.slice();
            operations[action.opIndex].parameters.splice(action.paramIndex, 1);
            return {
                ...state,
                operations,
            }
        case actions.ADD_LIST_PARAMETER_NAS:
            operations = state.operations.slice();
            operations[action.opIndex].parameters[action.paramIndex].list.push(
                {
                    //TODO: Remove it?
                    // name: "",
                    value: "",
                }
            )
            return {
                ...state,
                operations,
            }
        case actions.DELETE_LIST_PARAMETER_NAS:
            operations = state.operations.slice();
            operations[action.opIndex].parameters[action.paramIndex].list.splice(action.listIndex, 1);
            return {
                ...state,
                operations,
            }
        case actions.EDIT_LIST_PARAMETER_NAS:
            operations = state.operations.slice();
            operations[action.opIndex].parameters[action.paramIndex].list[action.listIndex].value = action.value;
            return {
                ...state,
                operations,
            }
        case actions.CHANGE_TRIAL_NAS:
            return {
                ...state,
                trial: action.trial,
            }
        case actions.CLOSE_SNACKBAR:
            return {
                ...state,
                snackOpen: false,
            }
        default:
            return state;
    }
};

export default nasCreateReducer;
