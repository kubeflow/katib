export const FETCH_NAS_JOB_INFO_REQUEST = 'FETCH_NAS_JOB_INFO_REQUEST';
export const FETCH_NAS_JOB_INFO_SUCCESS = 'FETCH_NAS_JOB_INFO_SUCCESS';
export const FETCH_NAS_JOB_INFO_FAILURE = 'FETCH_NAS_JOB_INFO_FAILURE';

export const fetchNASJobInfo = (experimentName, namespace) => ({
  type: FETCH_NAS_JOB_INFO_REQUEST,
  experimentName,
  namespace,
});
