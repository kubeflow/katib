import * as actions from '../actions/hpMonitorActions';

const initialState = {
  jobData: [],
  trialData: [],
  dialogTrialOpen: false,
  loading: false,
  trialName: '',
};

const hpMonitorReducer = (state = initialState, action) => {
  switch (action.type) {
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
