import * as actions from '../actions/hpMonitorActions';

const initialState = {
    experimentName: '',
    experimentNamespace: '',
    filterType: {
        "Created": true,
        "Running": true,
        "Restarting": true,
        "Succeeded": true,
        "Failed": true,
    },
    jobsList: [
    ],
    filteredJobsList: [
    ],
    jobData: [
    ],
    trialData: [
    ],
    dialogOpen: false,
    loading: false,
};

const hpMonitorReducer = (state = initialState, action) => {
    switch (action.type) {
        case actions.FILTER_JOBS:
            let jobs = state.jobsList.slice();
            let newList = jobs.filter(job =>
                (
                    job.name.includes(action.experimentName) &&
                    (
                        job.namespace == action.experimentNamespace ||
                        action.experimentNamespace == "All namespaces" ||
                        action.experimentNamespace.length == 0
                    )
                )
            )
            let types = Object.assign({}, state.filterType);
            var typeKeys = Object.keys(types);

            var filters = typeKeys.filter((key) => {
                return types[key]
            });

            let filteredJobs = newList.filter(job => filters.includes(job.status));

            return {
                ...state,
                filteredJobsList: filteredJobs,
                experimentName: action.experimentName,
                experimentNamespace: action.experimentNamespace
            }
        case actions.CHANGE_TYPE:
            jobs = state.jobsList.slice();
            newList = jobs.filter(job =>
                (
                    job.name.includes(state.experimentName) &&
                    (
                        job.namespace == state.experimentNamespace ||
                        state.experimentNamespace == "All namespaces" ||
                        state.experimentNamespace.length == 0
                    )
                )
            )
            types = Object.assign({}, state.filterType)
            types[action.filter] = action.checked;
            typeKeys = Object.keys(types);

            filters = typeKeys.filter((key) => {
                return types[key]
            });
            filteredJobs = newList.filter(job => filters.includes(job.status));
            
            return {
                ...state,
                filterType: types,
                filteredJobsList: filteredJobs,
            }
        case actions.FETCH_HP_JOBS_SUCCESS:
            jobs = action.jobs
            types = Object.assign({}, state.filterType);
            typeKeys = Object.keys(types);

            filters = typeKeys.filter((key) => {
                return types[key]
            });

            filteredJobs = jobs.filter(job =>
                (
                    filters.includes(job.status) &&
                    job.name.includes(state.experimentName) &&
                    (
                        job.namespace == state.experimentNamespace ||
                        state.experimentNamespace == "All namespaces" ||
                        state.experimentNamespace.length == 0
                    )
                )
            )
            return {
                ...state,
                jobsList: action.jobs,
                filteredJobsList: filteredJobs,
            }
        case actions.FETCH_HP_JOB_INFO_REQUEST:
            return {
                ...state,
                loading: true,
            }
        case actions.FETCH_HP_JOB_INFO_SUCCESS:
            return {
                ...state,
                jobData: action.jobData,
                loading: false,
            }
        case actions.FETCH_HP_JOB_INFO_FAILURE:
            return {
                ...state,
                loading: false,
            }
        case actions.FETCH_HP_JOB_REQUEST:
            return {
                ...state,
                loading: true,
            }
        case actions.FETCH_HP_JOB_SUCCESS:
            return {
                ...state,
                experiment: action.experiment,
                loading: false,
            }
        case actions.FETCH_HP_JOB_FAILURE:
            return {
                ...state,
                loading: false,
            }
        case actions.FETCH_HP_JOB_TRIAL_INFO_SUCCESS:
            return {
                ...state,
                trialData: action.trialData,
                dialogOpen: true,
            }
        case actions.CLOSE_DIALOG:
            return {
                ...state,
                dialogOpen: false,
            }
        default:
            return state;
    }
};

export default hpMonitorReducer;
