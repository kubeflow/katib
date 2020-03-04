export const FILTER_JOBS = 'FILTER_JOBS';

export const filterJobs = (experimentName, experimentNamespace) => ({
  type: FILTER_JOBS,
  experimentName,
  experimentNamespace,
});

export const CHANGE_TYPE = 'CHANGE_TYPE';

export const changeType = (filter, checked) => ({
  type: CHANGE_TYPE,
  filter,
  checked,
});

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

export const FETCH_HP_JOB_REQUEST = 'FETCH_HP_JOB_REQUEST';
export const FETCH_HP_JOB_SUCCESS = 'FETCH_HP_JOB_SUCCESS';
export const FETCH_HP_JOB_FAILURE = 'FETCH_HP_JOB_FAILURE';

export const fetchHPJob = (name, namespace) => ({
  type: FETCH_HP_JOB_REQUEST,
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

export const CLOSE_DIALOG_EXPERIMENT = 'CLOSE_DIALOG_EXPERIMENT';

export const closeDialogExperiment = () => ({
  type: CLOSE_DIALOG_EXPERIMENT,
});
