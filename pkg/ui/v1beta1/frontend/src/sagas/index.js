import { take, put, call, fork, all } from 'redux-saga/effects';
import axios from 'axios';
import * as templateActions from '../actions/templateActions';
import * as hpMonitorActions from '../actions/hpMonitorActions';
import * as hpCreateActions from '../actions/hpCreateActions';
import * as nasMonitorActions from '../actions/nasMonitorActions';
import * as nasCreateActions from '../actions/nasCreateActions';
import * as generalActions from '../actions/generalActions';

export const submitYaml = function* () {
  while (true) {
    const action = yield take(generalActions.SUBMIT_YAML_REQUEST);
    try {
      let isRightNamespace = false;
      for (const [, value] of Object.entries(action.yaml.split('\n'))) {
        const noSpaceLine = value.replace(/\s/g, '');
        if (noSpaceLine === 'trialTemplate:') {
          break;
        }
        if (
          action.globalNamespace === '' ||
          noSpaceLine === 'namespace:' + action.globalNamespace
        ) {
          isRightNamespace = true;
          break;
        }
      }
      if (isRightNamespace) {
        const result = yield call(goSubmitYaml, action.yaml);
        if (result.status === 200) {
          yield put({
            type: generalActions.SUBMIT_YAML_SUCCESS,
          });
        } else {
          yield put({
            type: generalActions.SUBMIT_YAML_FAILURE,
            message: result.message,
          });
        }
      } else {
        yield put({
          type: generalActions.SUBMIT_YAML_FAILURE,
          message: 'You can submit experiments only in ' + action.globalNamespace + ' namespace!',
        });
      }
    } catch (err) {
      yield put({
        type: generalActions.SUBMIT_YAML_FAILURE,
      });
    }
  }
};

const goSubmitYaml = function* (yaml) {
  try {
    const data = {
      yaml,
    };
    const result = yield call(axios.post, '/katib/submit_yaml/', data);
    return result;
  } catch (err) {
    return {
      status: 500,
      message: err.response.data,
    };
  }
};

export const deleteExperiment = function* () {
  while (true) {
    const action = yield take(generalActions.DELETE_EXPERIMENT_REQUEST);
    try {
      const result = yield call(goDeleteExperiment, action.name, action.namespace);
      if (result.status === 200) {
        // Lower case json keys for all experiments
        let experiments = Object.assign(result.data, {});
        experiments.map((template, i) => {
          return Object.keys(template).forEach(key => {
            const value = template[key];
            delete template[key];
            template[key.toLowerCase()] = value;
          });
        });
        yield put({
          type: generalActions.DELETE_EXPERIMENT_SUCCESS,
        });
        yield put({
          type: generalActions.FETCH_EXPERIMENTS_SUCCESS,
          experiments: experiments,
        });
      } else {
        yield put({
          type: generalActions.DELETE_EXPERIMENT_FAILURE,
        });
      }
    } catch (err) {
      yield put({
        type: generalActions.DELETE_EXPERIMENT_FAILURE,
      });
    }
  }
};

const goDeleteExperiment = function* (name, namespace) {
  try {
    const result = yield call(
      axios.get,
      `/katib/delete_experiment/?experimentName=${name}&namespace=${namespace}`,
    );
    return result;
  } catch (err) {
    yield put({
      type: generalActions.DELETE_EXPERIMENT_FAILURE,
    });
  }
};

export const submitHPJob = function* () {
  while (true) {
    const action = yield take(hpCreateActions.SUBMIT_HP_JOB_REQUEST);
    try {
      const result = yield call(goSubmitHPJob, action.data);
      if (result.status === 200) {
        yield put({
          type: hpCreateActions.SUBMIT_HP_JOB_SUCCESS,
        });
      } else {
        yield put({
          type: hpCreateActions.SUBMIT_HP_JOB_FAILURE,
          message: result.message,
        });
      }
    } catch (err) {
      yield put({
        type: hpCreateActions.SUBMIT_HP_JOB_FAILURE,
      });
    }
  }
};

const goSubmitHPJob = function* (postData) {
  try {
    const data = {
      postData,
    };
    const result = yield call(axios.post, '/katib/submit_hp_job/', data);
    return result;
  } catch (err) {
    return {
      status: 500,
      message: err.response.data,
    };
  }
};

export const fetchExperiments = function* () {
  while (true) {
    yield take(generalActions.FETCH_EXPERIMENTS_REQUEST);
    try {
      const result = yield call(goFetchExperiments);
      if (result.status === 200) {
        let data = Object.assign(result.data, {});
        data.map((template, i) => {
          return Object.keys(template).forEach(key => {
            const value = template[key];
            delete template[key];
            template[key.toLowerCase()] = value;
          });
        });
        yield put({
          type: generalActions.FETCH_EXPERIMENTS_SUCCESS,
          experiments: data,
        });
      } else {
        yield put({
          type: generalActions.FETCH_EXPERIMENTS_FAILURE,
        });
      }
    } catch (err) {
      yield put({
        type: generalActions.FETCH_EXPERIMENTS_FAILURE,
      });
    }
  }
};

const goFetchExperiments = function* () {
  try {
    const result = yield call(axios.get, '/katib/fetch_experiments/');
    return result;
  } catch (err) {
    yield put({
      type: generalActions.FETCH_EXPERIMENTS_FAILURE,
    });
  }
};

export const fetchExperiment = function* () {
  while (true) {
    const action = yield take(generalActions.FETCH_EXPERIMENT_REQUEST);
    try {
      const result = yield call(goFetchExperiment, action.name, action.namespace);
      if (result.status === 200) {
        yield put({
          type: generalActions.FETCH_EXPERIMENT_SUCCESS,
          experiment: result.data,
        });
      } else {
        yield put({
          type: generalActions.FETCH_EXPERIMENT_FAILURE,
        });
      }
    } catch (err) {
      yield put({
        type: generalActions.FETCH_EXPERIMENT_FAILURE,
      });
    }
  }
};

const goFetchExperiment = function* (name, namespace) {
  try {
    const result = yield call(
      axios.get,
      `/katib/fetch_experiment/?experimentName=${name}&namespace=${namespace}`,
    );
    return result;
  } catch (err) {
    yield put({
      type: generalActions.FETCH_EXPERIMENT_FAILURE,
    });
  }
};

export const fetchSuggestion = function* () {
  while (true) {
    const action = yield take(generalActions.FETCH_SUGGESTION_REQUEST);
    try {
      const result = yield call(goFetchSuggestion, action.name, action.namespace);
      if (result.status === 200) {
        yield put({
          type: generalActions.FETCH_SUGGESTION_SUCCESS,
          suggestion: result.data,
        });
      } else {
        yield put({
          type: generalActions.FETCH_SUGGESTION_FAILURE,
        });
      }
    } catch (err) {
      yield put({
        type: generalActions.FETCH_SUGGESTION_FAILURE,
      });
    }
  }
};

const goFetchSuggestion = function* (name, namespace) {
  try {
    const result = yield call(
      axios.get,
      `/katib/fetch_suggestion/?suggestionName=${name}&namespace=${namespace}`,
    );
    return result;
  } catch (err) {
    yield put({
      type: generalActions.FETCH_SUGGESTION_FAILURE,
    });
  }
};

export const fetchHPJobInfo = function* () {
  while (true) {
    const action = yield take(hpMonitorActions.FETCH_HP_JOB_INFO_REQUEST);
    try {
      const result = yield call(goFetchHPJobInfo, action.name, action.namespace);
      if (result.status === 200) {
        let data = result.data.split('\n').map((line, i) => line.split(','));
        yield put({
          type: hpMonitorActions.FETCH_HP_JOB_INFO_SUCCESS,
          jobData: data,
        });
      } else {
        yield put({
          type: hpMonitorActions.FETCH_HP_JOB_INFO_FAILURE,
        });
      }
    } catch (err) {
      yield put({
        type: hpMonitorActions.FETCH_HP_JOB_INFO_FAILURE,
      });
    }
  }
};

const goFetchHPJobInfo = function* (name, namespace) {
  try {
    const result = yield call(
      axios.get,
      `/katib/fetch_hp_job_info/?experimentName=${name}&namespace=${namespace}`,
    );
    return result;
  } catch (err) {
    yield put({
      type: hpMonitorActions.FETCH_HP_JOB_INFO_FAILURE,
    });
  }
};

export const fetchHPJobTrialInfo = function* () {
  while (true) {
    const action = yield take(hpMonitorActions.FETCH_HP_JOB_TRIAL_INFO_REQUEST);
    try {
      const result = yield call(gofetchHPJobTrialInfo, action.trialName, action.namespace);
      if (result.status === 200) {
        let data = result.data.split('\n').map((line, i) => line.split(','));
        yield put({
          type: hpMonitorActions.FETCH_HP_JOB_TRIAL_INFO_SUCCESS,
          trialData: data,
          trialName: action.trialName,
        });
      } else {
        yield put({
          type: hpMonitorActions.FETCH_HP_JOB_TRIAL_INFO_FAILURE,
        });
      }
    } catch (err) {
      yield put({
        type: hpMonitorActions.FETCH_HP_JOB_TRIAL_INFO_FAILURE,
      });
    }
  }
};

const gofetchHPJobTrialInfo = function* (trialName, namespace) {
  try {
    const result = yield call(
      axios.get,
      `/katib/fetch_hp_job_trial_info/?trialName=${trialName}&namespace=${namespace}`,
    );
    return result;
  } catch (err) {
    yield put({
      type: hpMonitorActions.FETCH_HP_JOB_TRIAL_INFO_FAILURE,
    });
  }
};

export const submitNASJob = function* () {
  while (true) {
    const action = yield take(nasCreateActions.SUBMIT_NAS_JOB_REQUEST);
    try {
      const result = yield call(goSubmitNASJob, action.data);
      if (result.status === 200) {
        yield put({
          type: nasCreateActions.SUBMIT_NAS_JOB_SUCCESS,
        });
      } else {
        yield put({
          type: nasCreateActions.SUBMIT_NAS_JOB_FAILURE,
          message: result.message,
        });
      }
    } catch (err) {
      yield put({
        type: nasCreateActions.SUBMIT_NAS_JOB_FAILURE,
      });
    }
  }
};

const goSubmitNASJob = function* (postData) {
  try {
    const data = {
      postData,
    };
    const result = yield call(axios.post, '/katib/submit_nas_job/', data);
    return result;
  } catch (err) {
    return {
      status: 500,
      message: err.response.data,
    };
  }
};

export const fetchNASJobInfo = function* () {
  while (true) {
    const action = yield take(nasMonitorActions.FETCH_NAS_JOB_INFO_REQUEST);
    try {
      const result = yield call(goFetchNASJobInfo, action.experimentName, action.namespace);
      if (result.status === 200) {
        let data = Object.assign(result.data, {});
        data.map((template, i) => {
          return Object.keys(template).forEach(key => {
            const value = template[key];
            delete template[key];
            template[key.toLowerCase()] = value;
          });
        });
        yield put({
          type: nasMonitorActions.FETCH_NAS_JOB_INFO_SUCCESS,
          steps: data,
        });
      } else {
        yield put({
          type: nasMonitorActions.FETCH_NAS_JOB_INFO_FAILURE,
        });
      }
    } catch (err) {
      yield put({
        type: nasMonitorActions.FETCH_NAS_JOB_INFO_FAILURE,
      });
    }
  }
};

const goFetchNASJobInfo = function* (experimentName, namespace) {
  try {
    const result = yield call(
      axios.get,
      `/katib/fetch_nas_job_info/?experimentName=${experimentName}&namespace=${namespace}`,
    );
    return result;
  } catch (err) {
    yield put({
      type: nasMonitorActions.FETCH_NAS_JOB_INFO_FAILURE,
    });
  }
};

export const fetchTrialTemplates = function* () {
  while (true) {
    yield take(templateActions.FETCH_TRIAL_TEMPLATES_REQUEST);
    try {
      const result = yield call(goFetchTrialTemplates);
      if (result.status === 200) {
        yield put({
          type: templateActions.FETCH_TRIAL_TEMPLATES_SUCCESS,
          trialTemplatesData: result.data.Data,
        });
      } else {
        yield put({
          type: templateActions.FETCH_TRIAL_TEMPLATES_FAILURE,
        });
      }
    } catch (err) {
      yield put({
        type: templateActions.FETCH_TRIAL_TEMPLATES_FAILURE,
      });
    }
  }
};

const goFetchTrialTemplates = function* (namespace) {
  try {
    const result = yield call(axios.get, `/katib/fetch_trial_templates`);
    return result;
  } catch (err) {
    yield put({
      type: templateActions.FETCH_TRIAL_TEMPLATES_FAILURE,
    });
  }
};

export const addTemplate = function* () {
  while (true) {
    const action = yield take(templateActions.ADD_TEMPLATE_REQUEST);
    try {
      const result = yield call(
        goAddTemplate,
        action.updatedConfigMapNamespace,
        action.updatedConfigMapName,
        action.updatedConfigMapPath,
        action.updatedTemplateYaml,
      );
      if (result.status === 200) {
        yield put({
          type: templateActions.ADD_TEMPLATE_SUCCESS,
          trialTemplatesData: result.data.Data,
        });
      } else {
        yield put({
          type: templateActions.ADD_TEMPLATE_FAILURE,
          error: result.error,
        });
      }
    } catch (err) {
      yield put({
        type: templateActions.ADD_TEMPLATE_FAILURE,
      });
    }
  }
};

const goAddTemplate = function* (
  updatedConfigMapNamespace,
  updatedConfigMapName,
  updatedConfigMapPath,
  updatedTemplateYaml,
) {
  try {
    const data = {
      updatedConfigMapNamespace,
      updatedConfigMapName,
      updatedConfigMapPath,
      updatedTemplateYaml,
    };
    const result = yield call(axios.post, '/katib/add_template/', data);
    return result;
  } catch (err) {
    return {
      status: 500,
      error: err.response.data,
    };
  }
};

export const editTemplate = function* () {
  while (true) {
    const action = yield take(templateActions.EDIT_TEMPLATE_REQUEST);
    try {
      const result = yield call(
        goEditTemplate,
        action.updatedConfigMapNamespace,
        action.updatedConfigMapName,
        action.configMapPath,
        action.updatedConfigMapPath,
        action.updatedTemplateYaml,
      );
      if (result.status === 200) {
        yield put({
          type: templateActions.EDIT_TEMPLATE_SUCCESS,
          trialTemplatesData: result.data.Data,
        });
      } else {
        yield put({
          type: templateActions.EDIT_TEMPLATE_FAILURE,
          error: result.error,
        });
      }
    } catch (err) {
      yield put({
        type: templateActions.EDIT_TEMPLATE_FAILURE,
      });
    }
  }
};

const goEditTemplate = function* (
  updatedConfigMapNamespace,
  updatedConfigMapName,
  configMapPath,
  updatedConfigMapPath,
  updatedTemplateYaml,
) {
  try {
    const data = {
      updatedConfigMapNamespace,
      updatedConfigMapName,
      configMapPath,
      updatedConfigMapPath,
      updatedTemplateYaml,
    };
    const result = yield call(axios.post, '/katib/edit_template/', data);
    return result;
  } catch (err) {
    return {
      status: 500,
      error: err.response.data,
    };
  }
};

export const deleteTemplate = function* () {
  while (true) {
    const action = yield take(templateActions.DELETE_TEMPLATE_REQUEST);
    try {
      const result = yield call(
        goDeleteTemplate,
        action.updatedConfigMapNamespace,
        action.updatedConfigMapName,
        action.updatedConfigMapPath,
      );
      if (result.status === 200) {
        yield put({
          type: templateActions.DELETE_TEMPLATE_SUCCESS,
          trialTemplatesData: result.data.Data,
        });
      } else {
        yield put({
          type: templateActions.DELETE_TEMPLATE_FAILURE,
          error: result.error,
        });
      }
    } catch (err) {
      yield put({
        type: templateActions.DELETE_TEMPLATE_FAILURE,
      });
    }
  }
};

const goDeleteTemplate = function* (
  updatedConfigMapNamespace,
  updatedConfigMapName,
  updatedConfigMapPath,
) {
  try {
    const data = {
      updatedConfigMapNamespace,
      updatedConfigMapName,
      updatedConfigMapPath,
    };
    const result = yield call(axios.post, '/katib/delete_template/', data);
    return result;
  } catch (err) {
    return {
      status: 500,
      error: err.response.data,
    };
  }
};

export const fetchNamespaces = function* () {
  while (true) {
    yield take(generalActions.FETCH_NAMESPACES_REQUEST);
    try {
      const result = yield call(goFetchNamespaces);
      if (result.status === 200) {
        let data = result.data;
        data.unshift('All namespaces');
        yield put({
          type: generalActions.FETCH_NAMESPACES_SUCCESS,
          namespaces: data,
        });
      } else {
        yield put({
          type: generalActions.FETCH_NAMESPACES_FAILURE,
        });
      }
    } catch (err) {
      yield put({
        type: generalActions.FETCH_NAMESPACES_FAILURE,
      });
    }
  }
};

const goFetchNamespaces = function* () {
  try {
    const result = yield call(axios.get, '/katib/fetch_namespaces');
    return result;
  } catch (err) {
    yield put({
      type: generalActions.FETCH_NAMESPACES_FAILURE,
    });
  }
};

export default function* rootSaga() {
  yield all([
    fork(fetchTrialTemplates),
    fork(fetchExperiments),
    fork(addTemplate),
    fork(editTemplate),
    fork(deleteTemplate),
    fork(submitYaml),
    fork(deleteExperiment),
    fork(submitHPJob),
    fork(submitNASJob),
    fork(fetchHPJobInfo),
    fork(fetchExperiment),
    fork(fetchSuggestion),
    fork(fetchHPJobTrialInfo),
    fork(fetchNASJobInfo),
    fork(fetchNamespaces),
  ]);
}
