import * as actions from '../actions/nasMonitorActions';

const initialState = {
  experimentName: '',
  experimentNamespace: '',
  filter: '',
  filterType: {
    Created: true,
    Running: true,
    Restarting: true,
    Succeeded: true,
    Failed: true,
  },
  jobsList: [],
  filteredJobsList: [],
  loading: false,
  steps: [],
};

const nasMonitorReducer = (state = initialState, action) => {
  switch (action.type) {
    case actions.FILTER_JOBS:
      let jobs = state.jobsList.slice();
      let newList = jobs.filter(
        job =>
          job.name.includes(action.experimentName) &&
          (job.namespace == action.experimentNamespace ||
            action.experimentNamespace == 'All namespaces' ||
            action.experimentNamespace.length == 0),
      );
      let types = Object.assign({}, state.filterType);
      var typeKeys = Object.keys(types);

      var filters = typeKeys.filter(key => {
        return types[key];
      });

      let filteredJobs = newList.filter(job => filters.includes(job.status));

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
          (job.namespace == state.experimentNamespace ||
            state.experimentNamespace == 'All namespaces' ||
            state.experimentNamespace.length == 0),
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
    case actions.FETCH_NAS_JOBS_SUCCESS:
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
          (job.namespace == state.experimentNamespace ||
            state.experimentNamespace == 'All namespaces' ||
            state.experimentNamespace.length == 0),
      );
      return {
        ...state,
        jobsList: action.jobs,
        filteredJobsList: filteredJobs,
      };
    case actions.FETCH_NAS_JOB_INFO_REQUEST:
      return {
        ...state,
        loading: true,
      };
    case actions.FETCH_NAS_JOB_INFO_FAILURE:
      return {
        ...state,
        loading: false,
      };
    case actions.FETCH_NAS_JOB_INFO_SUCCESS:
      return {
        ...state,
        loading: false,
        steps: action.steps,
      };
    default:
      return state;
  }
};

export default nasMonitorReducer;
