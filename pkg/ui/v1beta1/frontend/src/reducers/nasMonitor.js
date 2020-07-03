import * as actions from '../actions/nasMonitorActions';

const initialState = {
  loading: false,
  steps: [],
};

const nasMonitorReducer = (state = initialState, action) => {
  switch (action.type) {
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
