export const FETCH_HP_JOBS_REQUEST = 'FETCH_HP_JOBS_REQUEST';
export const FETCH_HP_JOBS_SUCCESS = 'FETCH_HP_JOBS_SUCCESS';
export const FETCH_HP_JOBS_FAILURE = 'FETCH_HP_JOBS_FAILURE';

export const fetchHPJobs = () => ({
  type: FETCH_HP_JOBS_REQUEST,
});

export const FETCH_HP_JOB_INFO_REQUEST = 'FETCH_HP_JOB_INFO_REQUEST';
export const FETCH_HP_JOB_INFO_SUCCESS = 'FETCH_HP_JOB_INFO_SUCCESS';
export const FETCH_HP_JOB_INFO_FAILURE = 'FETCH_HP_JOB_INFO_FAILURE';

export const fetchHPJobInfo = (name, namespace) => ({
  type: FETCH_HP_JOB_INFO_REQUEST,
  name,
  namespace,
});

export const FETCH_HP_JOB_TRIAL_INFO_REQUEST = 'FETCH_HP_JOB_TRIAL_INFO_REQUEST';
export const FETCH_HP_JOB_TRIAL_INFO_SUCCESS = 'FETCH_HP_JOB_TRIAL_INFO_SUCCESS';
export const FETCH_HP_JOB_TRIAL_INFO_FAILURE = 'FETCH_HP_JOB_TRIAL_INFO_FAILURE';

export const fetchHPJobTrialInfo = (trialName, namespace) => ({
  type: FETCH_HP_JOB_TRIAL_INFO_REQUEST,
  trialName,
  namespace,
});

export const CLOSE_DIALOG_TRIAL = 'CLOSE_DIALOG_TRIAL';

export const closeDialogTrial = () => ({
  type: CLOSE_DIALOG_TRIAL,
});
