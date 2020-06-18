export const CLOSE_DIALOG = 'CLOSE_DIALOG';

export const closeDialog = () => ({
  type: CLOSE_DIALOG,
});

export const OPEN_DIALOG = 'OPEN_DIALOG';

export const openDialog = (
  dialogType,
  configMapNamespace = '',
  configMapName = '',
  configMapPath = '',
  templateYaml = '',
) => ({
  type: OPEN_DIALOG,
  dialogType,
  configMapNamespace,
  configMapName,
  configMapPath,
  templateYaml,
});

export const FETCH_TRIAL_TEMPLATES_REQUEST = 'FETCH_TRIAL_TEMPLATES_REQUEST';
export const FETCH_TRIAL_TEMPLATES_SUCCESS = 'FETCH_TRIAL_TEMPLATES_SUCCESS';
export const FETCH_TRIAL_TEMPLATES_FAILURE = 'FETCH_TRIAL_TEMPLATES_FAILURE';

export const fetchTrialTemplates = () => ({
  type: FETCH_TRIAL_TEMPLATES_REQUEST,
});

export const ADD_TEMPLATE_REQUEST = 'ADD_TEMPLATE_REQUEST';
export const ADD_TEMPLATE_SUCCESS = 'ADD_TEMPLATE_SUCCESS';
export const ADD_TEMPLATE_FAILURE = 'ADD_TEMPLATE_FAILURE';

export const addTemplate = (
  updatedConfigMapNamespace,
  updatedConfigMapName,
  updatedConfigMapPath,
  updatedTemplateYaml,
) => ({
  type: ADD_TEMPLATE_REQUEST,
  updatedConfigMapNamespace,
  updatedConfigMapName,
  updatedConfigMapPath,
  updatedTemplateYaml,
});

export const EDIT_TEMPLATE_REQUEST = 'EDIT_TEMPLATE_REQUEST';
export const EDIT_TEMPLATE_SUCCESS = 'EDIT_TEMPLATE_SUCCESS';
export const EDIT_TEMPLATE_FAILURE = 'EDIT_TEMPLATE_FAILURE';

export const editTemplate = (
  updatedConfigMapNamespace,
  updatedConfigMapName,
  configMapPath,
  updatedConfigMapPath,
  updatedTemplateYaml,
) => ({
  type: EDIT_TEMPLATE_REQUEST,
  updatedConfigMapNamespace,
  updatedConfigMapName,
  configMapPath,
  updatedConfigMapPath,
  updatedTemplateYaml,
});

export const DELETE_TEMPLATE_REQUEST = 'DELETE_TEMPLATE_REQUEST';
export const DELETE_TEMPLATE_SUCCESS = 'DELETE_TEMPLATE_SUCCESS';
export const DELETE_TEMPLATE_FAILURE = 'DELETE_TEMPLATE_FAILURE';

export const deleteTemplate = (
  updatedConfigMapNamespace,
  updatedConfigMapName,
  updatedConfigMapPath,
) => ({
  type: DELETE_TEMPLATE_REQUEST,
  updatedConfigMapNamespace,
  updatedConfigMapName,
  updatedConfigMapPath,
});

export const CHANGE_TEMPLATE = 'CHANGE_TEMPLATE';

export const changeTemplate = (
  updatedConfigMapNamespace,
  updatedConfigMapName,
  updatedConfigMapPath,
  updatedTemplateYaml,
) => ({
  type: CHANGE_TEMPLATE,
  updatedConfigMapNamespace,
  updatedConfigMapName,
  updatedConfigMapPath,
  updatedTemplateYaml,
});

export const FILTER_TEMPLATES = 'FILTER_TEMPLATES';

export const filterTemplates = (filteredConfigMapNamespace, filteredConfigMapName) => ({
  type: FILTER_TEMPLATES,
  filteredConfigMapNamespace,
  filteredConfigMapName,
});
