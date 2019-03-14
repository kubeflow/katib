import { take, put, call, fork, select, all, takeEvery } from 'redux-saga/effects';
import axios from 'axios';
import * as templateActions from '../actions/templateActions';
import * as hpMonitorActions from '../actions/hpMonitorActions';
import * as nasMonitorActions from '../actions/nasMonitorActions';
import * as generalActions from '../actions/generalActions';


export const fetchWorkerTemplates = function *() {
    while (true) {
        const action = yield take(templateActions.FETCH_WORKER_TEMPLATES_REQUEST);
        try {
            const result = yield call(
                goFetchWorkerTemplates
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
                    type: templateActions.FETCH_WORKER_TEMPLATES_SUCCESS,
                    templates: data
                })
            } else {
                yield put({
                    type: templateActions.FETCH_WORKER_TEMPLATES_FAILURE,
                }) 
            }
        } catch (err) {
            yield put({
                type: templateActions.FETCH_WORKER_TEMPLATES_FAILURE,
            })
        }
    }
}

const goFetchWorkerTemplates = function *() {
    try {
        const result = yield call(
            axios.get,
            'http://127.0.0.1:9303/katib/fetch_worker_templates/',
        )
        return result
    } catch (err) {
        yield put({
            type: templateActions.FETCH_WORKER_TEMPLATES_FAILURE,
        })
    }
}

export const fetchCollectorTemplates = function *() {
    while (true) {
        const action = yield take(templateActions.FETCH_COLLECTOR_TEMPLATES_REQUEST);
        try {
            const result = yield call(
                goFetchCollectorTemplates
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
                    type: templateActions.FETCH_COLLECTOR_TEMPLATES_SUCCESS,
                    templates: data
                })
            } else {
                yield put({
                    type: templateActions.FETCH_COLLECTOR_TEMPLATES_FAILURE,
                }) 
            }
        } catch (err) {
            yield put({
                type: templateActions.FETCH_COLLECTOR_TEMPLATES_FAILURE,
            })
        }
    }
}

const goFetchCollectorTemplates = function *() {
    try {
        const result = yield call(
            axios.get,
            'http://127.0.0.1:9303/katib/fetch_collector_templates/',
        )
        return result
    } catch (err) {
        yield put({
            type: templateActions.FETCH_WORKER_TEMPLATES_FAILURE,
        })
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
            'http://127.0.0.1:9303/katib/fetch_hp_jobs/',
        )
        return result
    } catch (err) {
        yield put({
            type: hpMonitorActions.FETCH_HP_JOBS_FAILURE,
        })
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
            'http://127.0.0.1:9303/katib/fetch_nas_jobs/',
        )
        return result
    } catch (err) {
        yield put({
            type: nasMonitorActions.FETCH_NAS_JOBS_FAILURE,
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

const goAddTemplate = function *(name, yaml, kind) {
    try {
        const data = {
            name, yaml, kind
        }
        const result = yield call(
            axios.post,
            'http://127.0.0.1:9303/katib/update_template/',
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

const goEditTemplate = function *(name, yaml, kind) {
    try {
        const data = {
            name, yaml, kind
        }
        const result = yield call(
            axios.post,
            'http://127.0.0.1:9303/katib/update_template/',
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
                action.templateType,
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

const goDeleteTemplate = function *(name, kind) {
    try {
        const data = {
            name, kind
        }
        const result = yield call(
            axios.post,
            'http://127.0.0.1:9303/katib/delete_template/',
            data,
        )
        return result
    } catch (err) {
        yield put({
            type: templateActions.DELETE_TEMPLATE_FAILURE,
        })
    }
}


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
            'http://127.0.0.1:9303/katib/submit_yaml/',
            data,
        )
        return result
    } catch (err) {
        yield put({
            type: generalActions.SUBMIT_YAML_FAILURE,
        })
    }
}

export const fetchHPJobInfo = function *() {
    while (true) {
        const action = yield take(hpMonitorActions.FETCH_JOB_INFO_REQUEST);
        try {
            const result = yield call(
                goFetchHPJobInfo,
                action.id
            )
            if (result.status === 200) {
                let data = result.data.split("\n").map((line, i) => line.split(','))
                yield put({
                    type: hpMonitorActions.FETCH_JOB_INFO_SUCCESS,
                    jobData: data
                })
            } else {
                yield put({
                    type: hpMonitorActions.FETCH_JOB_INFO_FAILURE,
                }) 
            }
        } catch (err) {
            yield put({
                type: hpMonitorActions.FETCH_JOB_INFO_FAILURE,
            })
        }
    }
}

const goFetchHPJobInfo = function *(id) {
    try {
        const result = yield call(
            axios.get,
            `http://127.0.0.1:9303/katib/fetch_hp_job_info/?id=${id}`,
        )
        return result
    } catch (err) {
        yield put({
            type: hpMonitorActions.FETCH_JOB_INFO_FAILURE,
        })
    }
}

export const fetchNASJobInfo = function *() {
    while (true) {
        const action = yield take(nasMonitorActions.FETCH_JOB_INFO_REQUEST);
        try {
            const result = yield call(
                goFetchNASJobInfo,
                action.id
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
                    type: nasMonitorActions.FETCH_JOB_INFO_SUCCESS,
                    steps: data,
                })
            } else {
                yield put({
                    type: nasMonitorActions.FETCH_JOB_INFO_FAILURE,
                }) 
            }
        } catch (err) {
            yield put({
                type: nasMonitorActions.FETCH_JOB_INFO_FAILURE,
            })
        }
    }
}

const goFetchNASJobInfo = function *(id) {
    try {
        const result = yield call(
            axios.get,
            `http://127.0.0.1:9303/katib/fetch_nas_job_info/?id=${id}`,
        )
        return result
    } catch (err) {
        yield put({
            type: nasMonitorActions.FETCH_JOB_INFO_FAILURE,
        })
    }
}

export const fetchWorkerInfo = function *() {
    while (true) {
        const action = yield take(hpMonitorActions.FETCH_WORKER_INFO_REQUEST);
        try {
            const result = yield call(
                goFetchWorkerInfo,
                action.studyID,
                action.workerID
            )
            if (result.status === 200) {
                let data = result.data.split("\n").map((line, i) => line.split(','))
                yield put({
                    type: hpMonitorActions.FETCH_WORKER_INFO_SUCCESS,
                    workerData: data
                })
            } else {
                yield put({
                    type: hpMonitorActions.FETCH_WORKER_INFO_FAILURE,
                }) 
            }
        } catch (err) {
            yield put({
                type: hpMonitorActions.FETCH_WORKER_INFO_FAILURE,
            })
        }
    }
}

const goFetchWorkerInfo = function *(studyID, workerID) {
    try {
        const result = yield call(
            axios.get,
            `http://127.0.0.1:9303/katib/fetch_worker_info/?studyID=${studyID}&workerID=${workerID}`,
        )
        return result
    } catch (err) {
        yield put({
            type: hpMonitorActions.FETCH_WORKER_INFO_FAILURE,
        })
    }
}

export default function* rootSaga() {
    yield all([
        fork(fetchWorkerTemplates),
        fork(fetchCollectorTemplates),
        fork(fetchHPJobs),
        fork(fetchNASJobs),
        fork(addTemplate), 
        fork(editTemplate),
        fork(deleteTemplate),
        fork(submitYaml),
        fork(fetchHPJobInfo),
        fork(fetchWorkerInfo),
        fork(fetchNASJobInfo)
    ]);
};
