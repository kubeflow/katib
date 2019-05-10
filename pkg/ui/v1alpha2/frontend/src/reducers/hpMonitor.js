import * as actions from '../actions/hpMonitorActions';

const initialState = {
    filter: '',
    filterType: {
        "Running": true,
        "Failed": true,
        "Completed": true,
    },
    jobsList: [
    ],
    filteredJobsList: [
    ],
    jobData: [
            [
              "TrialName",
              "accuracy",
              "Validation-accuracy",
              "--lr",
              "--num-layers",
              "--optimizer"
            ],
            [
              "md8cb40c3e01a911",
              "0.105313",
              "0.113854",
              "0.0257",
              "3",
              "ftrl"
            ],
            [
              "l027d9a4162276bf",
              "0.105313",
              "0.113854",
              "0.0146",
              "5",
              "ftrl"
            ],
            [
              "ne8e8f581eee5f4b",
              "0.105313",
              "0.113854",
              "0.0139",
              "3",
              "ftrl"
            ],
            [
              "ze0db6b764c36e83",
              "0.972344",
              "0.960888",
              "0.0287",
              "5",
              "adam"
            ],
            [
              "y07c6f7dbeb469bb",
              "0.105313",
              "0.113854",
              "0.0270",
              "5",
              "ftrl"
            ],
            [
              "ead3ade22bbbcf19",
              "0.981094",
              "0.966660",
              "0.0239",
              "4",
              "adam"
            ],
            [
              "q1f2c679a9543c28",
              "0.105313",
              "0.113854",
              "0.0113",
              "2",
              "ftrl"
            ],
            [
              "s9b4188ba1bc5407",
              "0.999687",
              "0.981290",
              "0.0190",
              "2",
              "sgd"
            ],
            [
              "v7aad7071b2d98d6",
              "0.105313",
              "0.113854",
              "0.0111",
              "3",
              "ftrl"
            ],
            [
              "g47fb08fec5bcd15",
              "0.984219",
              "0.969248",
              "0.0204",
              "4",
              "adam"
            ],
            [
              "ocd06e6f2089a52b",
              "0.105313",
              "0.113854",
              "0.0255",
              "5",
              "ftrl"
            ],
            [
              "a38994c9a725922d",
              "0.995313",
              "0.971338",
              "0.0103",
              "5",
              "adam"
            ],
            [
              ""
            ]
    ],
    trialData: [
      [
        "symbol",
        "time",
        "value"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:15",
        "0.114171"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:16",
        "0.111406"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:17",
        "0.114687"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:18",
        "0.112344"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:19",
        "0.113906"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:20",
        "0.112344"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:21",
        "0.113906"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:22",
        "0.115937"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:23",
        "0.105313"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:24",
        "0.115312"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:25",
        "0.116646"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:26",
        "0.115312"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:27",
        "0.116646"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:28",
        "0.115312"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:29",
        "0.116646"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:30",
        "0.102188"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:31",
        "0.116646"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:32",
        "0.114687"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:33",
        "0.112344"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:34",
        "0.113906"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:35",
        "0.115937"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:36",
        "0.105313"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:37",
        "0.115937"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:38",
        "0.116646"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:39",
        "0.115312"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:40",
        "0.116646"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:41",
        "0.102188"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:42",
        "0.116646"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:43",
        "0.102188"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:44",
        "0.116646"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:45",
        "0.114687"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:46",
        "0.111406"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:47",
        "0.113906"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:48",
        "0.112344"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:49",
        "0.113906"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:50",
        "0.112344"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:51",
        "0.113906"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:52",
        "0.112344"
      ],
      [
        "accuracy",
        "2019-04-23T19:14:53",
        "0.105313"
      ],
      [
        "Validation-accuracy",
        "2019-04-23T19:14:17",
        "0.113854"
      ],
      [
        "Validation-accuracy",
        "2019-04-23T19:14:19",
        "0.113854"
      ],
      [
        "Validation-accuracy",
        "2019-04-23T19:14:21",
        "0.113854"
      ],
      [
        "Validation-accuracy",
        "2019-04-23T19:14:23",
        "0.113854"
      ],
      [
        "Validation-accuracy",
        "2019-04-23T19:14:25",
        "0.113854"
      ],
      [
        "Validation-accuracy",
        "2019-04-23T19:14:27",
        "0.113854"
      ],
      [
        "Validation-accuracy",
        "2019-04-23T19:14:29",
        "0.113854"
      ],
      [
        "Validation-accuracy",
        "2019-04-23T19:14:30",
        "0.113854"
      ],
      [
        "Validation-accuracy",
        "2019-04-23T19:14:32",
        "0.113854"
      ],
      [
        "Validation-accuracy",
        "2019-04-23T19:14:34",
        "0.113854"
      ],
      [
        "Validation-accuracy",
        "2019-04-23T19:14:36",
        "0.113854"
      ],
      [
        "Validation-accuracy",
        "2019-04-23T19:14:38",
        "0.113854"
      ],
      [
        "Validation-accuracy",
        "2019-04-23T19:14:39",
        "0.113854"
      ],
      [
        "Validation-accuracy",
        "2019-04-23T19:14:41",
        "0.113854"
      ],
      [
        "Validation-accuracy",
        "2019-04-23T19:14:43",
        "0.113854"
      ],
      [
        "Validation-accuracy",
        "2019-04-23T19:14:45",
        "0.113854"
      ],
      [
        "Validation-accuracy",
        "2019-04-23T19:14:47",
        "0.113854"
      ],
      [
        "Validation-accuracy",
        "2019-04-23T19:14:49",
        "0.113854"
      ],
      [
        "Validation-accuracy",
        "2019-04-23T19:14:51",
        "0.113854"
      ],
      [
        "Validation-accuracy",
        "2019-04-23T19:14:53",
        "0.113854"
      ],
      [
        ""
      ]
    ],
    dialogOpen: false,
    loading: false,
};

const hpMonitorReducer = (state = initialState, action) => {
    switch (action.type) {
        case actions.FILTER_JOBS:
            let jobs = state.jobsList.slice();
            let newList = jobs.filter(job => job.name.includes(action.filter));

            let avTypes = Object.assign({}, state.filterType);
            var typeKeys = Object.keys(avTypes);

            var avFilters = typeKeys.filter((key) => {
                return avTypes[key]
            });

            let filteredJobs = newList.filter(job => avFilters.includes(job.status));

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
        case actions.FETCH_HP_JOBS_SUCCESS:
            jobs = action.jobs
            avTypes = Object.assign({}, state.filterType);
            typeKeys = Object.keys(avTypes);

            avFilters = typeKeys.filter((key) => {
                return avTypes[key]
            });

            filteredJobs = jobs.filter(job => avFilters.includes(job.status));
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