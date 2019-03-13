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
    workerSpec: ["Test 1", "Test 2"], // fetch names from backend 
    worker: '',
    numLayers: '1',
    inputSize: ['32', '32', '3'],
    outputSize: ['10'],
    suggestionAlgorithms: ["rl", "enas"], // fetch these
    suggestionAlgorithm: "",
    suggestionParameters: [
        {
            name: "DefaultGrid",
            value: "1",
        },
        {
            name: "--lr",
            value: "4",
        },
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
        default:
            return state;
    }
};

export default nasCreateReducer;