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