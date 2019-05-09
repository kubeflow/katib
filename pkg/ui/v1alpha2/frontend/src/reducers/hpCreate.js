import * as actions from '../actions/hpCreateActions';

const initialState = {
    loading: false,
    commonParametersMetadata: [
        {
            name: "Name",
            value: "random-experiment",
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
            value: "Validation-accuracy",
            description: "Name for the objective metric"
        }
    ],
    additionalMetricNames: [
        {
            value: "accuracy"
        }
    ],
    algorithmName: [ "random" ],
    allAlgorithms: ["grid", "random", "hyperband"],
    algorithmSettings: [

    ],
    parameters: [
        {
            name: "--lr",
            parameterType: "double",
            feasibleSpace: "feasibleSpace",
            min: "0.01",
            max: "0.03",
            list: [],
        },
        {
            name: "--num-layers",
            parameterType: "int",
            feasibleSpace: "feasibleSpace",
            min: "2",
            max: "5",
            list: [],
        },
        {
            name: "--optimizer",
            parameterType: "categorical",
            feasibleSpace: "list",
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
    allParameterTypes: ["int", "double", "categorical"],
    trial: "cpuTrialTemplate.yaml",
    currentYaml: '',
};

const filterValue = (obj, key) => {
    return obj.findIndex(p => p.name === key)
};

const hpCreateReducer = (state = initialState, action) => {
    switch (action.type) {
        case actions.CHANGE_YAML_HP:
            return {
                ...state,
                currentYaml: action.payload,
            }
        case actions.CHANGE_META_HP:
            let meta = state.commonParametersMetadata.slice();
            let index = filterValue(meta, action.name);
            meta[index].value = action.value;
            return {
                ...state,
                commonParametersMetadata: meta,
            }
        case actions.CHANGE_SPEC_HP:
            let spec = state.commonParametersSpec.slice();
            index = filterValue(spec, action.name);
            spec[index].value = action.value;
            return {
                ...state,
                commonParametersSpec: spec,
            }
        case actions.CHANGE_OBJECTIVE_HP:
            let obj = state.objective.slice();
            index = filterValue(obj, action.name);
            obj[index].value = action.value;
            return {
                ...state,
                objective: obj,
            }
        case actions.ADD_METRICS_HP:
            let additionalMetricNames = state.additionalMetricNames.slice()
            additionalMetricNames.push({
                value: "",
            })
            return {
                ...state,
                additionalMetricNames: additionalMetricNames,
            }
        case actions.DELETE_METRICS_HP:
            additionalMetricNames = state.additionalMetricNames.slice()
            additionalMetricNames.splice(action.index, 1)
            return {
                ...state,
                additionalMetricNames: additionalMetricNames,
            }
        case actions.EDIT_METRICS_HP:
            additionalMetricNames = state.additionalMetricNames.slice()
            additionalMetricNames[action.index].value = action.value
            return {
                ...state,
                additionalMetricNames: additionalMetricNames,
            }
        case actions.CHANGE_ALGORITHM_NAME_HP:
            return {
                ...state, 
                algorithmName: action.algorithmName,
            }
        case actions.ADD_ALGORITHM_SETTING_HP:
            let algorithmSettings = state.algorithmSettings.slice();
            let setting = {name: "", value: ""};
            algorithmSettings.push(setting);
            return {
                ...state,
                algorithmSettings: algorithmSettings,
            }
        case actions.CHANGE_ALGORITHM_SETTING_HP:
            algorithmSettings = state.algorithmSettings.slice();
            algorithmSettings[action.index][action.field] = action.value;
            return {
                ...state,
                algorithmSettings: algorithmSettings,
            }
        case actions.DELETE_ALGORITHM_SETTING_HP:
            algorithmSettings = state.algorithmSettings.slice();
            algorithmSettings.splice(action.index, 1);
            return {
                ...state,
                algorithmSettings: algorithmSettings,
            }
        case actions.ADD_PARAMETER_HP:
            let parameters = state.parameters.slice();
            parameters.push({
                name: "",
                parameterType: "",
                feasibleSpace: "feasibleSpace",
                min: "",
                max: "",
                list: [],
            })
            return {
                ...state,
                parameters: parameters,
            }
        case actions.EDIT_PARAMETER_HP:
            parameters = state.parameters.slice();
            parameters[action.index][action.field] = action.value;
            return {
                ...state,
                parameters: parameters,
            }
        case actions.DELETE_PARAMETER_HP:
            parameters = state.parameters.slice();
            parameters.splice(action.index, 1);
            return {
                ...state,
                parameters: parameters,
            }
        case actions.ADD_LIST_PARAMETER_HP:
            parameters = state.parameters.slice();
            parameters[action.paramIndex].list.push({
                value: "",
            })
            return {
                ...state,
                parameters: parameters
            }
        case actions.EDIT_LIST_PARAMETER_HP:
            parameters = state.parameters.slice();
            parameters[action.paramIndex].list[action.index].value = action.value;
            return {
                ...state,
                parameters: parameters
            }
        case actions.DELETE_LIST_PARAMETER_HP:
            parameters = state.parameters.slice();  
            parameters[action.paramIndex].list.splice(action.index, 1);
            return {
                ...state,
                parameters: parameters
            }
        case actions.CHANGE_TRIAL_HP:
            return {
                ...state,
                trial: action.trial,
            }
        default:
            return state;
    }
};

export default hpCreateReducer;