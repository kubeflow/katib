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

export const FETCH_JOB_INFO_REQUEST = "FETCH_JOB_INFO_REQUEST"
export const FETCH_JOB_INFO_SUCCESS = "FETCH_JOB_INFO_SUCCESS"
export const FETCH_JOB_INFO_FAILURE = "FETCH_JOB_INFO_FAILURE"

export const fetchJobInfo = (id) => ({
    type: FETCH_JOB_INFO_REQUEST,
    id
})