import * as actions from '../actions/nasMonitorActions';

const initialState = {
    filter: '',
    filterType: {
        "Running": true,
        "Failed": true,
        "Succeeded": true,
    },
    jobsList: [
        {
            name: "Job 1",
            status: "Running",
            id: "1", 
        },
        {
            name: "Job 2",
            status: "Failed",
            id: "2", 
        },
        {
            name: "Job 3",
            status: "Succeeded",
            id: "3", 
        }
    ],
    filteredJobsList: [
        {
            name: "Job 1",
            status: "Running",
            id: "1", 
        },
        {
            name: "Job 2",
            status: "Failed",
            id: "2", 
        },
        {
            name: "Job 3",
            status: "Succeeded",
            id: "3", 
        }
    ],
    steps: [
        {
            id: "1",
            name: "Generation 1",
            architecture: "ASD",
            metricsName: "Accuracy",
            metricsValue: "0.99",
            link: "link",
        },
        {
            id: "2",
            name: "Generation 2",
            architecture: "DAS",
            metricsName: "Accuracy",
            metricsValue: "0.99999",
            link: "link",
        },
        {
            id: "3",
            name: "Generation 3",
            architecture: "DAS",
            metricsName: "Accuracy",
            metricsValue: "0.999999999",
            link: "link",
        },
    ]
};

const nasMonitorReducer = (state = initialState, action) => {
    switch (action.type) {
        case actions.FILTER_JOBS:
            const jobs = state.jobsList.slice();
            const newList = jobs.filter(job => job.name.includes(action.filter));

            const avTypes = Object.assign({}, state.filterType);
            var typeKeys = Object.keys(avTypes);

            var avFilters = typeKeys.filter((key) => {
                return avTypes[key]
            });

            const filteredJobs = newList.filter(job => avFilters.includes(job.status));

            return {
                ...state,
                filteredJobsList: filteredJobs,
                filter: action.filter,
            }
        case actions.CHANGE_TYPE:
            const types = Object.assign({}, state.filterType)
            types[action.filter] = action.checked;
            var keys = Object.keys(types);

            var filters = keys.filter((key) => {
                return types[key]
            });
            const jobsList = state.jobsList.slice();
            const filtered = jobsList.filter(job => filters.includes(job.status));
            
            return {
                ...state,
                filterType: types,
                filteredJobsList: filtered,
            }
        default:
            return state;
    }
};


export default nasMonitorReducer;