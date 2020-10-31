export const TOGGLE_MENU = 'TOGGLE_MENU';

export const toggleMenu = state => {
  return {
    type: TOGGLE_MENU,
    state,
  };
};

export const CLOSE_SNACKBAR = 'CLOSE_SNACKBAR';

export const closeSnackbar = () => {
  return {
    type: CLOSE_SNACKBAR,
  };
};

export const SUBMIT_YAML_REQUEST = 'SUBMIT_YAML_REQUEST';
export const SUBMIT_YAML_FAILURE = 'SUBMIT_YAML_FAILURE';
export const SUBMIT_YAML_SUCCESS = 'SUBMIT_YAML_SUCCESS';

export const submitYaml = (yaml, globalNamespace) => ({
  type: SUBMIT_YAML_REQUEST,
  yaml,
  globalNamespace,
});

export const DELETE_EXPERIMENT_REQUEST = 'DELETE_EXPERIMENT_REQUEST';
export const DELETE_EXPERIMENT_FAILURE = 'DELETE_EXPERIMENT_FAILURE';
export const DELETE_EXPERIMENT_SUCCESS = 'DELETE_EXPERIMENT_SUCCESS';

export const deleteExperiment = (name, namespace) => ({
  type: DELETE_EXPERIMENT_REQUEST,
  name,
  namespace,
});

export const OPEN_DELETE_EXPERIMENT_DIALOG = 'OPEN_DELETE_EXPERIMENT_DIALOG';

export const openDeleteExperimentDialog = (name, namespace) => ({
  type: OPEN_DELETE_EXPERIMENT_DIALOG,
  name,
  namespace,
});

export const CLOSE_DELETE_EXPERIMENT_DIALOG = 'CLOSE_DELETE_EXPERIMENT_DIALOG';

export const closeDeleteExperimentDialog = () => ({
  type: CLOSE_DELETE_EXPERIMENT_DIALOG,
});

export const FETCH_NAMESPACES_REQUEST = 'FETCH_NAMESPACES_REQUEST';
export const FETCH_NAMESPACES_SUCCESS = 'FETCH_NAMESPACES_SUCCESS';
export const FETCH_NAMESPACES_FAILURE = 'FETCH_NAMESPACES_FAILURE';

export const fetchNamespaces = () => ({
  type: FETCH_NAMESPACES_REQUEST,
});

export const CHANGE_GLOBAL_NAMESPACE = 'CHANGE_GLOBAL_NAMESPACE';

export const FETCH_EXPERIMENT_REQUEST = 'FETCH_EXPERIMENT_REQUEST';
export const FETCH_EXPERIMENT_SUCCESS = 'FETCH_EXPERIMENT_SUCCESS';
export const FETCH_EXPERIMENT_FAILURE = 'FETCH_EXPERIMENT_FAILURE';

export const fetchExperiment = (name, namespace) => ({
  type: FETCH_EXPERIMENT_REQUEST,
  name,
  namespace,
});

export const CLOSE_DIALOG_EXPERIMENT = 'CLOSE_DIALOG_EXPERIMENT';

export const closeDialogExperiment = () => ({
  type: CLOSE_DIALOG_EXPERIMENT,
});

export const FETCH_SUGGESTION_REQUEST = 'FETCH_SUGGESTION_REQUEST';
export const FETCH_SUGGESTION_SUCCESS = 'FETCH_SUGGESTION_SUCCESS';
export const FETCH_SUGGESTION_FAILURE = 'FETCH_SUGGESTION_FAILURE';

export const fetchSuggestion = (name, namespace) => ({
  type: FETCH_SUGGESTION_REQUEST,
  name,
  namespace,
});

export const CLOSE_DIALOG_SUGGESTION = 'CLOSE_DIALOG_SUGGESTION';

export const closeDialogSuggestion = () => ({
  type: CLOSE_DIALOG_SUGGESTION,
});

export const FETCH_EXPERIMENTS_REQUEST = 'FETCH_EXPERIMENTS_REQUEST';
export const FETCH_EXPERIMENTS_SUCCESS = 'FETCH_EXPERIMENTS_SUCCESS';
export const FETCH_EXPERIMENTS_FAILURE = 'FETCH_EXPERIMENTS_FAILURE';

export const fetchExperiments = () => ({
  type: FETCH_EXPERIMENTS_REQUEST,
});

export const FILTER_EXPERIMENTS = 'FILTER_EXPERIMENTS';

export const filterExperiments = (experimentName, experimentNamespace) => ({
  type: FILTER_EXPERIMENTS,
  experimentName,
  experimentNamespace,
});

export const CHANGE_STATUS = 'CHANGE_STATUS';

export const changeStatus = (filter, checked) => ({
  type: CHANGE_STATUS,
  filter,
  checked,
});

export const CHANGE_EARLY_STOPPING_ALGORITHM = 'CHANGE_EARLY_STOPPING_ALGORITHM';

export const changeEarlyStoppingAlgorithm = algorithmName => ({
  type: CHANGE_EARLY_STOPPING_ALGORITHM,
  algorithmName,
});

export const ADD_EARLY_STOPPING_SETTING = 'ADD_EARLY_STOPPING_SETTING';

export const addEarlyStoppingSetting = () => ({
  type: ADD_EARLY_STOPPING_SETTING,
});

export const CHANGE_EARLY_STOPPING_SETTING = 'CHANGE_EARLY_STOPPING_SETTING';

export const changeEarlyStoppingSetting = (index, field, value) => ({
  type: CHANGE_EARLY_STOPPING_SETTING,
  index,
  field,
  value,
});

export const DELETE_EARLY_STOPPING_SETTING = 'DELETE_EARLY_STOPPING_SETTING';

export const deleteEarlyStoppingSetting = index => ({
  type: DELETE_EARLY_STOPPING_SETTING,
  index,
});

export const CHANGE_TRIAL_TEMPLATE_SOURCE = 'CHANGE_TRIAL_TEMPLATE_SOURCE';

export const changeTrialTemplateSource = source => ({
  type: CHANGE_TRIAL_TEMPLATE_SOURCE,
  source,
});

export const ADD_PRIMARY_POD_LABEL = 'ADD_PRIMARY_POD_LABEL';

export const addPrimaryPodLabel = () => ({
  type: ADD_PRIMARY_POD_LABEL,
});

export const CHANGE_PRIMARY_POD_LABEL = 'CHANGE_PRIMARY_POD_LABEL';

export const changePrimaryPodLabel = (fieldName, index, value) => ({
  type: CHANGE_PRIMARY_POD_LABEL,
  fieldName,
  index,
  value,
});

export const DELETE_PRIMARY_POD_LABEL = 'DELETE_PRIMARY_POD_LABEL';

export const deletePrimaryPodLabel = index => ({
  type: DELETE_PRIMARY_POD_LABEL,
  index,
});

export const CHANGE_TRIAL_TEMPLATE_SPEC = 'CHANGE_TRIAL_TEMPLATE_SPEC';

export const changeTrialTemplateSpec = (name, value) => ({
  type: CHANGE_TRIAL_TEMPLATE_SPEC,
  name,
  value,
});

export const FILTER_TEMPLATES_EXPERIMENT = 'FILTER_TEMPLATES_EXPERIMENT';

export const filterTemplatesExperiment = (
  configMapNamespaceIndex,
  configMapNameIndex,
  configMapPathIndex,
) => ({
  type: FILTER_TEMPLATES_EXPERIMENT,
  configMapNamespaceIndex,
  configMapNameIndex,
  configMapPathIndex,
});

export const CHANGE_TRIAL_TEMPLATE_YAML = 'CHANGE_TRIAL_TEMPLATE_YAML';

export const changeTrialTemplateYAML = templateYAML => ({
  type: CHANGE_TRIAL_TEMPLATE_YAML,
  templateYAML,
});

export const CHANGE_TRIAL_PARAMETERS = 'CHANGE_TRIAL_PARAMETERS';

export const changeTrialParameters = (index, name, reference, description) => ({
  type: CHANGE_TRIAL_PARAMETERS,
  index,
  name,
  reference,
  description,
});

export const VALIDATION_ERROR = 'VALIDATION_ERROR';

export const validationError = message => ({
  type: VALIDATION_ERROR,
  message,
});
