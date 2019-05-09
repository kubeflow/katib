export const FILTER_JOBS = "FILTER_JOBS";

export const filterJobs = (filter) => ({
    type: FILTER_JOBS,
    filter,
})

export const CHANGE_TYPE = "CHANGE_TYPE";

export const changeType = (filter, checked) => ({
    type: CHANGE_TYPE, 
    filter, checked
})

export const FETCH_HP_JOBS_REQUEST = "FETCH_HP_JOBS_REQUEST";
export const FETCH_HP_JOBS_SUCCESS = "FETCH_HP_JOBS_SUCCESS";
export const FETCH_HP_JOBS_FAILURE = "FETCH_HP_JOBS_FAILURE";

export const fetchHPJobs = () => ({
    type: FETCH_HP_JOBS_REQUEST,
})

export const FETCH_HP_JOB_INFO_REQUEST = "FETCH_HP_JOB_INFO_REQUEST";
export const FETCH_HP_JOB_INFO_SUCCESS = "FETCH_HP_JOB_INFO_SUCCESS";
export const FETCH_HP_JOB_INFO_FAILURE = "FETCH_HP_JOB_INFO_FAILURE";

export const fetchHPJobInfo = (name) => ({
    type: FETCH_HP_JOB_INFO_REQUEST,
    name
})

export const FETCH_TRIAL_INFO_REQUEST = "FETCH_TRIAL_INFO_REQUEST";
export const FETCH_TRIAL_INFO_SUCCESS = "FETCH_TRIAL_INFO_SUCCESS";
export const FETCH_TRIAL_INFO_FAILURE = "FETCH_TRIAL_INFO_FAILURE";

export const fetchTrialInfo = (studyID, workerID) => ({
    type: FETCH_TRIAL_INFO_REQUEST,
    studyID, workerID
})

export const CLOSE_DIALOG = "CLOSE_DIALOG";

export const closeDialog = () => ({
    type: CLOSE_DIALOG,
})
