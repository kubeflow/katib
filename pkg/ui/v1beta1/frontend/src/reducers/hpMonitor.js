import * as actions from '../actions/hpMonitorActions';

const initialState = {
  experimentName: '',
  experimentNamespace: 'All namespaces',
  filterType: {
    Created: true,
    Running: true,
    Restarting: true,
    Succeeded: true,
    Failed: true,
  },
  jobsList: [],
  filteredJobsList: [],
  jobData: [],
  trialData: [],
  dialogTrialOpen: false,
  loading: false,
  trialName: '',
};

const hpMonitorReducer = (state = initialState, action) => {
  switch (action.type) {
    case actions.FILTER_JOBS:
      var jobs = state.jobsList.slice();
      var newList = jobs.filter(
        job =>
          job.name.includes(action.experimentName) &&
          (job.namespace === action.experimentNamespace ||
            action.experimentNamespace === 'All namespaces'),
      );
      var types = Object.assign({}, state.filterType);
      var typeKeys = Object.keys(types);

      var filters = typeKeys.filter(key => {
        return types[key];
      });

      var filteredJobs = newList.filter(job => filters.includes(job.status));

      return {
        ...state,
        filteredJobsList: filteredJobs,
        experimentName: action.experimentName,
        experimentNamespace: action.experimentNamespace,
      };
    case actions.CHANGE_TYPE:
      jobs = state.jobsList.slice();
      newList = jobs.filter(
        job =>
          job.name.includes(state.experimentName) &&
          (job.namespace === state.experimentNamespace ||
            state.experimentNamespace === 'All namespaces'),
      );
      types = Object.assign({}, state.filterType);
      types[action.filter] = action.checked;
      typeKeys = Object.keys(types);

      filters = typeKeys.filter(key => {
        return types[key];
      });
      filteredJobs = newList.filter(job => filters.includes(job.status));

      return {
        ...state,
        filterType: types,
        filteredJobsList: filteredJobs,
      };
    case actions.FETCH_HP_JOBS_SUCCESS:
      jobs = action.jobs;
      types = Object.assign({}, state.filterType);
      typeKeys = Object.keys(types);

      filters = typeKeys.filter(key => {
        return types[key];
      });

      filteredJobs = jobs.filter(
        job =>
          filters.includes(job.status) &&
          job.name.includes(state.experimentName) &&
          (job.namespace === state.experimentNamespace ||
            state.experimentNamespace === 'All namespaces'),
      );
      return {
        ...state,
        jobsList: action.jobs,
        filteredJobsList: filteredJobs,
      };
    case actions.FETCH_HP_JOB_INFO_REQUEST:
      return {
        ...state,
        loading: true,
        dialogTrialOpen: false,
      };
    case actions.FETCH_HP_JOB_INFO_SUCCESS:
      return {
        ...state,
        jobData: action.jobData,
        loading: false,
      };
    case actions.FETCH_HP_JOB_INFO_FAILURE:
      return {
        ...state,
        loading: false,
      };
    case actions.FETCH_HP_JOB_TRIAL_INFO_SUCCESS:
      return {
        ...state,
        trialData: action.trialData,
        dialogTrialOpen: true,
        trialName: action.trialName,
      };
    case actions.CLOSE_DIALOG_TRIAL:
      return {
        ...state,
        dialogTrialOpen: false,
      };
    default:
      return state;
  }
};

export default hpMonitorReducer;
