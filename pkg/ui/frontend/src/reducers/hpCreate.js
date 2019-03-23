import * as actions from '../actions/hpCreateActions';

const initialState = {
    loading: false,
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
    //  specify NASCONFIG?
    parameterConfig: [
        {
            parameterType: "double",
            name: "--lr",
            feasible: "feasible",
            min: "0.01",
            max: "0.03",
            list: [],
        },
        {
            parameterType: "int",
            name: "--num-layers",
            feasible: "feasible",
            min: "2",
            max: "5",
            list: [
                {
                    value: "",
                }
            ],
        },
        {
            parameterType: "categorical",
            name: "--optimizer",
            feasible: "list",
            min: "",
            max: "",
            list: [
                {
                    value: "sgd",
                },
                {
                    value: "adam",
                },
                {
                    value: "ftrl"
                }
            ],
        },
    ],
    paramTypes: ["int", "double", "categorical"],
    worker: "cpuWorkerTemplate.yaml",
    suggestionAlgorithms: ["grid", "random", "hyperband"], // fetch these
    suggestionAlgorithm: "random",
    requestNumber: "3",
    suggestionParameters: [
    ],
    currentYaml: '',
};

const filterValue = (obj, key) => {
    return obj.findIndex(p => p.name === key)
};

const hpCreateReducer = (state = initialState, action) => {
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
        case actions.ADD_PARAMETER:
            let params = state.parameterConfig.slice();
            params.push({name: "", parameterType: ""})
            return {
                ...state,
                parameterConfig: params,
            }
        case actions.EDIT_PARAMETER:
            params = state.parameterConfig.slice();
            params[action.index][action.field] = action.value;
            return {
                ...state,
                parameterConfig: params,
            }
        case actions.DELETE_PARAMETER:
            params = state.parameterConfig.slice();
            params.splice(action.index, 1);
            return {
                ...state,
                parameterConfig: params,
            }
        case actions.ADD_LIST_PARAMETER_HP:
            params = state.parameterConfig.slice();
            params[action.paramIndex].list.push({
                value: "",
            })
            return {
                ...state,
                parameterConfig: params
            }
        case actions.EDIT_LIST_PARAMETER_HP:
            params = state.parameterConfig.slice();
            params[action.paramIndex].list[action.index] = action.value;
            return {
                ...state,
                parameterConfig: params
            }
        case actions.DELETE_LIST_PARAMETER_HP:
            params = state.parameterConfig.slice();
            params[action.paramIndex].list.splice(action.index, 1);
            return {
                ...state,
                parameterConfig: params
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
            let parameter = {name: "", parameterType: ""};
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
        case actions.SUBMIT_HP_JOB_REQUEST:
            return {
                ...state,
                loading: true,
            }
        case actions.SUBMIT_HP_JOB_SUCCESS:
            return {
                ...state,
                loading: false,
            }
        case actions.SUBMIT_HP_JOB_FAILURE:
            return {
                ...state,
                loading: false,
            }
        case actions.CHANGE_REQUEST_NUMBER:
            return {
                ...state,
                requestNumber: action.number,
            }
        case actions.ADD_METRICS_HP:
            let metricsName = state.metricsName.slice()
            metricsName.push({
                value: "",
            })
            return {
                ...state,
                metricsName: metricsName,
            }
        case actions.DELETE_METRICS_HP:
            metricsName = state.metricsName.slice()
            metricsName.splice(action.index, 1)
            return {
                ...state,
                metricsName: metricsName,
            }
        case actions.EDIT_METRICS_HP:
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

export default hpCreateReducer;