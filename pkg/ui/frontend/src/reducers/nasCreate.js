import * as actions from '../actions/nasCreateActions';

const initialState = {
    commonParametersMetadata: [
        {
            name: "Namespace",
            value: "kubeflow",
            description: "Namespace to deploy a study into"
        },
        {
            name: "Name",
            value: "nasrl-example",
            description: "A name of a study"
        }
    ],
    commonParametersSpec: [
        {
            name: "Name",
            value: "nasrl-example",
            description: "A name of a study"
        },
        // owner is always crd
        {
            name: "OptimizationType",
            value: "maximize",
            description: "Optimization type"
        },
        {
            name: "OptimizationValueName",
            value: "Validation-accuracy",
            description: "A name of metrics to optimize"
        },
        // check for float
        {
            name: "OptimizationGoals",
            value: "0.99",
            description: "A threshold to optimize up to"
        },
        // check for int
        {
            name: "RequestCount",
            value: "4",
            description: "Number of requests"
        },
    ],
    metricsName: [
        {
            value: "accuracy",
        }
    ],
    worker: 'cpuWorkerTemplate.yaml',
    numLayers: '1',
    inputSize: ['32', '32', '3'],
    outputSize: ['10'],
    paramTypes: ["int", "double", "categorical"],
    operations: [
        {
            operationType: "convolution",
            parameterconfigs: [
                {
                    parameterType: "categorical",
                    name: "filter_size",
                    feasible: "list",
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
                    parameterType: "categorical",
                    name: "num_filter",
                    feasible: "list",
                    min: "",
                    max: "",
                    step: "",
                    list: [
                        {
                            value: "32",
                        },
                        {
                            value: "48",
                        }
                    ],
                },
            ]
        }
    ],
    suggestionAlgorithms: ["nasrl", "enas"], // fetch these
    suggestionAlgorithm: "nasrl",
    requestNumber: "3",
    suggestionParameters: [
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
    currentYaml: '',
};

const filterValue = (obj, key) => {
    return obj.findIndex(p => p.name === key)
};

const nasCreateReducer = (state = initialState, action) => {
    console.log(state)
    switch (action.type) {
        case actions.CHANGE_YAML:
            return {
                ...state,
                currentYaml: action.payload,
            }
        case actions.CHANGE_META:
            let meta = state.commonParametersMetadata.slice();
            let index = filterValue(meta, action.name);
            meta[index].value = action.value;
            return {
                ...state,
                commonParametersMetadata: meta,
            }
        case actions.CHANGE_SPEC:
            let spec = state.commonParametersSpec.slice();
            index = filterValue(spec, action.name);
            spec[index].value = action.value;
            return {
                ...state,
                commonParametersSpec: spec,
            }
        case actions.CHANGE_WORKER:
            return {
                ...state,
                worker: action.worker,
            }
        case actions.CHANGE_ALGORITHM:
            return {
                ...state, 
                suggestionAlgorithm: action.algorithm,
            }
        case actions.ADD_SUGGESTION_PARAMETER:
            let suggestionParameters = state.suggestionParameters.slice();
            let parameter = {name: "", value: ""};
            suggestionParameters.push(parameter);
            return {
                ...state,
                suggestionParameters: suggestionParameters,
            }
        case actions.CHANGE_SUGGESTION_PARAMETER:
            suggestionParameters = state.suggestionParameters.slice();
            suggestionParameters[action.index][action.field] = action.value;
            return {
                ...state,
                suggestionParameters: suggestionParameters,
            }
        case actions.DELETE_SUGGESTION_PARAMETER:
            suggestionParameters = state.suggestionParameters.slice();
            suggestionParameters.splice(action.index, 1);
            return {
                ...state,
                suggestionParameters: suggestionParameters,
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
                parameterconfigs: [],
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
        case actions.ADD_PARAMETER:
            operations = state.operations.slice();
            operations[action.opIndex].parameterconfigs.push(
               {name: "", parameterType: ""}
            )
            return {
                ...state,
                operations,
            }
        case actions.CHANGE_PARAMETER:
            operations = state.operations.slice();
            operations[action.opIndex].parameterconfigs[action.paramIndex][action.field] = action.value;
            return {
                ...state,
                operations,
            }
        case actions.DELETE_PARAMETER:
            operations = state.operations.slice();
            operations[action.opIndex].parameterconfigs.splice(action.paramIndex, 1);
            return {
                ...state,
                operations,
            }
        case actions.ADD_LIST_PARAMETER_NAS:
            operations = state.operations.slice();
            operations[action.opIndex].parameterconfigs[action.paramIndex].list.push(
                {
                    name: "",
                    value: "",
                }
            )
            return {
                ...state,
                operations,
            }
        case actions.DELETE_LIST_PARAMETER_NAS:
            operations = state.operations.slice();
            operations[action.opIndex].parameterconfigs[action.paramIndex].list.splice(action.listIndex, 1);
            return {
                ...state,
                operations,
            }
        case actions.EDIT_LIST_PARAMETER_NAS:
            operations = state.operations.slice();
            operations[action.opIndex].parameterconfigs[action.paramIndex].list[action.listIndex] = action.value;
            return {
                ...state,
                operations,
            }
        case actions.CHANGE_REQUEST_NUMBER:
            return {
                ...state,
                requestNumber: action.number,
            }
        case actions.ADD_METRICS_NAS:
            let metricsName = state.metricsName.slice()
            metricsName.push({
                value: "",
            })
            return {
                ...state,
                metricsName: metricsName,
            }
        case actions.DELETE_METRICS_NAS:
            metricsName = state.metricsName.slice()
            metricsName.splice(action.index, 1)
            return {
                ...state,
                metricsName: metricsName,
            }
        case actions.EDIT_METRICS_NAS:
            metricsName = state.metricsName.slice()
            metricsName[action.index].value = action.value
            return {
                ...state,
                metricsName: metricsName,
            }
        default:
            return state;
    }
};

export default nasCreateReducer;