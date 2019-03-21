import * as actions from '../actions/nasCreateActions';
import { stat } from 'fs';

const initialState = {
    commonParametersMetadata: [
        {
            name: "Namespace",
            value: "kubeflow",
            description: "Namespace to deploy a study into"
        },
        {
            name: "Name",
            value: "random-example",
            description: "A name of a study"
        }
    ],
    commonParametersSpec: [
        {
            name: "Name",
            value: "random-example",
            description: "A name of a study"
        },
        // owner is always crd
        {
            name: "OptimizationType",
            value: "Maximize",
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
        // list here
        {
            name: "MetricsName",
            value: "list here",
            description: "A name of a study"
        }
    ],
    //  specify NASCONFIG?
    // select!
    workerSpec: ["cpuWorkerTemplate.yaml", "Test 2"], // fetch names from backend 
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
                    parameterType: "double",
                    name: "Shit2",
                    feasible: "feasible",
                    min: "",
                    max: "",
                    step: "",
                    list: [
                        {
                            value: "",
                        }
                    ],
                },
                {
                    parameterType: "int",
                    name: "Shit",
                    feasible: "list",
                    min: "",
                    max: "",
                    step: "",
                    list: [],
                },
            ]
        }
    ],
    suggestionAlgorithms: ["nasrl", "enas"], // fetch these
    suggestionAlgorithm: "nasrl",
    suggestionParameters: [
    ],
    currentYaml: '',
};

const filterValue = (obj, key) => {
    return obj.findIndex(p => p.name === key)
};

const nasCreateReducer = (state = initialState, action) => {
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
        case actions.ADD_LIST_PARAMETER:
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
        case actions.DELETE_LIST_PARAMETER:
            operations = state.operations.slice();
            operations[action.opIndex].parameterconfigs[action.paramIndex].list.splice(action.listIndex, 1);
            return {
                ...state,
                operations,
            }
        case actions.EDIT_LIST_PARAMETER:
            operations = state.operations.slice();
            operations[action.opIndex].parameterconfigs[action.paramIndex].list[action.listIndex] = action.value;
            return {
                ...state,
                operations,
            }
        default:
            return state;
    }
};

export default nasCreateReducer;