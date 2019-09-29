import { take, put, call, fork, select, all, takeEvery } from 'redux-saga/effects';
import axios from 'axios';
import * as templateActions from '../actions/templateActions';
import * as hpMonitorActions from '../actions/hpMonitorActions';
import * as hpCreateActions from '../actions/hpCreateActions';
import * as nasMonitorActions from '../actions/nasMonitorActions';
import * as nasCreateActions from '../actions/nasCreateActions';
import * as generalActions from '../actions/generalActions';


export const submitYaml = function *() {
    while (true) {
        const action = yield take(generalActions.SUBMIT_YAML_REQUEST);
        try {
            const result = yield call(
                goSubmitYaml,
                action.yaml
            )
            if (result.status === 200) {
                yield put({
                    type: generalActions.SUBMIT_YAML_SUCCESS,
                })
            } else {
                yield put({
                    type: generalActions.SUBMIT_YAML_FAILURE,
                    message: result.message,
                }) 
            }
        } catch (err) {
            yield put({
                type: generalActions.SUBMIT_YAML_FAILURE,
            })
        }
    }
}

const goSubmitYaml = function *(yaml) {
    try {
        const data = {
            yaml
        }
        const result = yield call(
            axios.post,
            '/katib/submit_yaml/',
            data,
        )
        return result
    } catch (err) {
        return {
            status: 500,
            message: err.response.data,
        }
    }
}

export const deleteExperiment = function *() {
    while (true) {
        const action = yield take(generalActions.DELETE_EXPERIMENT_REQUEST);
        try {
            const result = yield call(
                goDeleteExperiment,
                action.experimentName
            )
            if (result.status === 200) {
                yield put({
                    type: generalActions.DELETE_EXPERIMENT_SUCCESS,
                })
            } else {
                yield put({
                    type: generalActions.DELETE_EXPERIMENT_FAILURE,
                }) 
            }
        } catch (err) {
            yield put({
                type: generalActions.DELETE_EXPERIMENT_FAILURE,
            })
        }
    }
}

const goDeleteExperiment = function *(experimentName) {
    try {
        const result = yield call(
            axios.get,
            `/katib/delete_experiment/?experimentName=${experimentName}`,
        )
        return result
    } catch (err) {
        yield put({
            type: generalActions.DELETE_EXPERIMENT_FAILURE,
        })
    }
}

export const submitHPJob = function *() {
    while (true) {
        const action = yield take(hpCreateActions.SUBMIT_HP_JOB_REQUEST);
        try {
            const result = yield call(
                goSubmitHPJob,
                action.data
            )
            if (result.status === 200) {
                yield put({
                    type: hpCreateActions.SUBMIT_HP_JOB_SUCCESS,
                })
            } else {
                yield put({
                    type: hpCreateActions.SUBMIT_HP_JOB_FAILURE,
                    message: result.message,
                }) 
            }
        } catch (err) {
            yield put({
                type: hpCreateActions.SUBMIT_HP_JOB_FAILURE,
            })
        }
    }
}

const goSubmitHPJob = function *(postData) {
    try {
        const data = {
            postData
        }
        const result = yield call(
            axios.post,
            '/katib/submit_hp_job/',
            data,
        )
        return result
    } catch (err) {
        return {
            status: 500,
            message: err.response.data,
        }
    }
}

export const fetchHPJobs = function *() {
    while (true) {
        const action = yield take(hpMonitorActions.FETCH_HP_JOBS_REQUEST);
        try {
            const result = yield call(
                goFetchHPJobs
            )
            if (result.status === 200) {
                let data = Object.assign(result.data, {})
                data.map((template, i) => {
                    Object.keys(template).forEach(key => {
                        const value = template[key];
                        delete template[key];
                        template[key.toLowerCase()] = value;
                    });
                })
                yield put({
                    type: hpMonitorActions.FETCH_HP_JOBS_SUCCESS,
                    jobs: data
                })
            } else {
                yield put({
                    type: hpMonitorActions.FETCH_HP_JOBS_FAILURE,
                }) 
            }
        } catch (err) {
            yield put({
                type: hpMonitorActions.FETCH_HP_JOBS_FAILURE,
            })
        }
    }
}

const goFetchHPJobs = function *() {
    try {
        const result = yield call(
            axios.get,
            '/katib/fetch_hp_jobs/',
        )
        return result
    } catch (err) {
        yield put({
            type: hpMonitorActions.FETCH_HP_JOBS_FAILURE,
        })
    }
}

export const fetchHPJobInfo = function *() {
    while (true) {
        const action = yield take(hpMonitorActions.FETCH_HP_JOB_INFO_REQUEST);
        try {
            const result = yield call(
                goFetchHPJobInfo,
                action.experimentName
            )
            if (result.status === 200) {
                let data = result.data.split("\n").map((line, i) => line.split(','))
                yield put({
                    type: hpMonitorActions.FETCH_HP_JOB_INFO_SUCCESS,
                    jobData: data
                })
            } else {
                yield put({
                    type: hpMonitorActions.FETCH_HP_JOB_INFO_FAILURE,
                }) 
            }
        } catch (err) {
            yield put({
                type: hpMonitorActions.FETCH_HP_JOB_INFO_FAILURE,
            })
        }
    }
}

const goFetchHPJobInfo = function *(experimentName) {
    try {
        const result = yield call(
            axios.get,
            `/katib/fetch_hp_job_info/?experimentName=${experimentName}`,
        )
        return result
    } catch (err) {
        yield put({
            type: hpMonitorActions.FETCH_HP_JOB_INFO_FAILURE,
        })
    }
}

export const fetchHPJobTrialInfo = function *() {
    while (true) {
        const action = yield take(hpMonitorActions.FETCH_HP_JOB_TRIAL_INFO_REQUEST);
        try {
            const result = yield call(
                gofetchHPJobTrialInfo,
                action.trialName
            )
            if (result.status === 200) {
                let data = result.data.split("\n").map((line, i) => line.split(','))
                yield put({
                    type: hpMonitorActions.FETCH_HP_JOB_TRIAL_INFO_SUCCESS,
                    trialData: data
                })
            } else {
                yield put({
                    type: hpMonitorActions.FETCH_HP_JOB_TRIAL_INFO_FAILURE,
                }) 
            }
        } catch (err) {
            yield put({
                type: hpMonitorActions.FETCH_HP_JOB_TRIAL_INFO_FAILURE,
            })
        }
    }
}

const gofetchHPJobTrialInfo = function *(trialName) {
    try {
        const result = yield call(
            axios.get,
            `/katib/fetch_hp_job_trial_info/?trialName=${trialName}`,
        )
        return result
    } catch (err) {
        yield put({
            type: hpMonitorActions.FETCH_HP_JOB_TRIAL_INFO_FAILURE,
        })
    }
}

export const submitNASJob = function *() {
    while (true) {
        const action = yield take(nasCreateActions.SUBMIT_NAS_JOB_REQUEST);
        try {
            const result = yield call(
                goSubmitNASJob,
                action.data
            )
            if (result.status === 200) {
                yield put({
                    type: nasCreateActions.SUBMIT_NAS_JOB_SUCCESS,
                })
            } else {
                yield put({
                    type: nasCreateActions.SUBMIT_NAS_JOB_FAILURE,
                    message: result.message,
                }) 
            }
        } catch (err) {
            yield put({
                type: nasCreateActions.SUBMIT_NAS_JOB_FAILURE,
            })
        }
    }
}

const goSubmitNASJob = function *(postData) {
    try {
        const data = {
            postData
        }
        const result = yield call(
            axios.post,
            '/katib/submit_nas_job/',
            data,
        )
        return result
    } catch (err) {
        return {
            status: 500,
            message: err.response.data,
        }
    }
}


export const fetchNASJobs = function *() {
    while (true) {
        const action = yield take(nasMonitorActions.FETCH_NAS_JOBS_REQUEST);
        try {
            const result = yield call(
                goFetchNASJobs
            )
            if (result.status === 200) {
                let data = Object.assign(result.data, {})
                data.map((template, i) => {
                    Object.keys(template).forEach(key => {
                        const value = template[key];
                        delete template[key];
                        template[key.toLowerCase()] = value;
                    });
                })
                yield put({
                    type: nasMonitorActions.FETCH_NAS_JOBS_SUCCESS,
                    jobs: data
                })
            } else {
                yield put({
                    type: nasMonitorActions.FETCH_NAS_JOBS_FAILURE,
                }) 
            }
        } catch (err) {
            yield put({
                type: nasMonitorActions.FETCH_NAS_JOBS_FAILURE,
            })
        }
    }
}

const goFetchNASJobs = function *() {
    try {
        const result = yield call(
            axios.get,
            '/katib/fetch_nas_jobs/',
        )
        return result
    } catch (err) {
        yield put({
            type: nasMonitorActions.FETCH_NAS_JOBS_FAILURE,
        })
    }
}

export const fetchNASJobInfo = function *() {
    while (true) {
        const action = yield take(nasMonitorActions.FETCH_NAS_JOB_INFO_REQUEST);
        try {
            const result = yield call(
                goFetchNASJobInfo,
                action.experimentName
            )
            if (result.status === 200) {
                let data = Object.assign(result.data, {})
                data.map((template, i) => {
                    Object.keys(template).forEach(key => {
                        const value = template[key];
                        delete template[key];
                        template[key.toLowerCase()] = value;
                    });
                })
                yield put({
                    type: nasMonitorActions.FETCH_NAS_JOB_INFO_SUCCESS,
                    steps: data,
                })
            } else {
                yield put({
                    type: nasMonitorActions.FETCH_NAS_JOB_INFO_FAILURE,
                }) 
            }
        } catch (err) {
            yield put({
                type: nasMonitorActions.FETCH_NAS_JOB_INFO_FAILURE,
            })
        }
    }
}

const goFetchNASJobInfo = function *(experimentName) {
    try {
        const result = yield call(
            axios.get,
            `/katib/fetch_nas_job_info/?experimentName=${experimentName}`,
        )
        return result
    } catch (err) {
        yield put({
            type: nasMonitorActions.FETCH_NAS_JOB_INFO_FAILURE,
        })
    }
}


export const fetchTrialTemplates = function *() {
    while (true) {
        const action = yield take(templateActions.FETCH_TRIAL_TEMPLATES_REQUEST);
        try {
            const result = yield call(
                goFetchTrialTemplates,
                action.namespace
            )
            if (result.status === 200) {
                let data = Object.assign(result.data, {})
                data.map((template, i) => {
                    Object.keys(template).forEach(key => {
                        const value = template[key];
                        delete template[key];
                        template[key.toLowerCase()] = value;
                    });
                })
                yield put({
                    type: templateActions.FETCH_TRIAL_TEMPLATES_SUCCESS,
                    templates: data
                })
            } else {
                yield put({
                    type: templateActions.FETCH_TRIAL_TEMPLATES_FAILURE,
                }) 
            }
        } catch (err) {
            yield put({
                type: templateActions.FETCH_TRIAL_TEMPLATES_FAILURE,
            })
        }
    }
}

const goFetchTrialTemplates = function *(namespace) {
    try {
        const result = yield call(
            axios.get,
            `/katib/fetch_trial_templates/?namespace=${namespace}`,
        )
        return result
    } catch (err) {
        yield put({
            type: templateActions.FETCH_TRIAL_TEMPLATES_FAILURE,
        })
    }
}

export const addTemplate = function *() {
    while (true) {
        const action = yield take(templateActions.ADD_TEMPLATE_REQUEST);
        try {
            const result = yield call(
                goAddTemplate,
                action.name,
                action.yaml,
                action.kind,
                action.action
            )
            if (result.status === 200) {
                let data = Object.assign(result.data.Data, {})
                data.map((template, i) => {
                    Object.keys(template).forEach(key => {
                        const value = template[key];
                        delete template[key];
                        template[key.toLowerCase()] = value;
                    });
                })
                yield put({
                    type: templateActions.ADD_TEMPLATE_SUCCESS,
                    templates: data,
                    templateType: result.data.TemplateType
                })
            } else {
                yield put({
                    type: templateActions.ADD_TEMPLATE_FAILURE,
                }) 
            }
        } catch (err) {
            yield put({
                type: templateActions.ADD_TEMPLATE_FAILURE,
            })
        }
    }
}

const goAddTemplate = function *(name, yaml, kind, action) {
    try {
        const data = {
            name, yaml, kind, action
        }
        const result = yield call(
            axios.post,
            '/katib/update_template/',
            data,
        )
        return result
    } catch (err) {
        yield put({
            type: templateActions.ADD_TEMPLATE_FAILURE,
        })
    }
}

export const editTemplate = function *() {
    while (true) {
        const action = yield take(templateActions.EDIT_TEMPLATE_REQUEST);
        try {
            const result = yield call(
                goEditTemplate,
                action.name,
                action.yaml,
                action.kind,
                action.action
            )
            if (result.status === 200) {
                let data = Object.assign(result.data.Data, {})
                data.map((template, i) => {
                    Object.keys(template).forEach(key => {
                        const value = template[key];
                        delete template[key];
                        template[key.toLowerCase()] = value;
                    });
                })
                yield put({
                    type: templateActions.EDIT_TEMPLATE_SUCCESS,
                    templates: data,
                    templateType: result.data.TemplateType
                })
            } else {
                yield put({
                    type: templateActions.EDIT_TEMPLATE_FAILURE,
                }) 
            }
        } catch (err) {
            yield put({
                type: templateActions.EDIT_TEMPLATE_FAILURE,
            })
        }
    }
}

const goEditTemplate = function *(name, yaml, kind, action) {
    try {
        const data = {
            name, yaml, kind, action
        }
        const result = yield call(
            axios.post,
            '/katib/update_template/',
            data,
        )
        return result
    } catch (err) {
        yield put({
            type: templateActions.EDIT_TEMPLATE_FAILURE,
        })
    }
}

export const deleteTemplate = function *() {
    while (true) {
        const action = yield take(templateActions.DELETE_TEMPLATE_REQUEST);
        try {
            const result = yield call(
                goDeleteTemplate,
                action.name,
                action.kind,
                action.action
            )
            if (result.status === 200) {
                let data = Object.assign(result.data.Data, {})
                data.map((template, i) => {
                    Object.keys(template).forEach(key => {
                        const value = template[key];
                        delete template[key];
                        template[key.toLowerCase()] = value;
                    });
                })
                yield put({
                    type: templateActions.DELETE_TEMPLATE_SUCCESS,
                    templates: data,
                    templateType: result.data.TemplateType
                })
            } else {
                yield put({
                    type: templateActions.DELETE_TEMPLATE_FAILURE,
                }) 
            }
        } catch (err) {
            yield put({
                type: templateActions.DELETE_TEMPLATE_FAILURE,
            })
        }
    }
}

const goDeleteTemplate = function *(name, kind, action) {
    try {
        const data = {
            name, kind, action
        }
        const result = yield call(
            axios.post,
            '/katib/update_template/',
            data,
        )
        return result
    } catch (err) {
        yield put({
            type: templateActions.DELETE_TEMPLATE_FAILURE,
        })
    }
}

export default function* rootSaga() {
    yield all([
        fork(fetchTrialTemplates),
        fork(fetchHPJobs),
        fork(fetchNASJobs),
        fork(addTemplate), 
        fork(editTemplate),
        fork(deleteTemplate),
        fork(submitYaml),
        fork(deleteExperiment),
        fork(submitHPJob),
        fork(submitNASJob),
        fork(fetchHPJobInfo),
        fork(fetchHPJobTrialInfo),
        fork(fetchNASJobInfo)
    ]);
};
