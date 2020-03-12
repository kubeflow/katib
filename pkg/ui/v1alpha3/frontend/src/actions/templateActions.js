export const CLOSE_DIALOG = 'CLOSE_DIALOG';

export const closeDialog = () => ({
  type: CLOSE_DIALOG,
});

export const OPEN_DIALOG = 'OPEN_DIALOG';

export const openDialog = (
  dialogType,
  namespace = '',
  configMapName = '',
  templateName = '',
  templateYaml = '',
) => ({
  type: OPEN_DIALOG,
  dialogType,
  namespace,
  configMapName,
  templateName,
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

export const addTemplate = (edittedNamespace, edittedConfigMapName, edittedName, edittedYaml) => ({
  type: ADD_TEMPLATE_REQUEST,
  edittedNamespace,
  edittedConfigMapName,
  edittedName,
  edittedYaml,
});

export const EDIT_TEMPLATE_REQUEST = 'EDIT_TEMPLATE_REQUEST';
export const EDIT_TEMPLATE_SUCCESS = 'EDIT_TEMPLATE_SUCCESS';
export const EDIT_TEMPLATE_FAILURE = 'EDIT_TEMPLATE_FAILURE';

export const editTemplate = (
  edittedNamespace,
  edittedConfigMapName,
  currentName,
  edittedName,
  edittedYaml,
) => ({
  type: EDIT_TEMPLATE_REQUEST,
  edittedNamespace,
  edittedConfigMapName,
  currentName,
  edittedName,
  edittedYaml,
});

export const DELETE_TEMPLATE_REQUEST = 'DELETE_TEMPLATE_REQUEST';
export const DELETE_TEMPLATE_SUCCESS = 'DELETE_TEMPLATE_SUCCESS';
export const DELETE_TEMPLATE_FAILURE = 'DELETE_TEMPLATE_FAILURE';

export const deleteTemplate = (edittedNamespace, edittedConfigMapName, edittedName) => ({
  type: DELETE_TEMPLATE_REQUEST,
  edittedNamespace,
  edittedConfigMapName,
  edittedName,
});

export const CHANGE_TEMPLATE = 'CHANGE_TEMPLATE';

export const changeTemplate = (
  edittedTemplateNamespace,
  edittedTemplateConfigMapName,
  edittedTemplateName,
  edittedTemplateYaml,
  edittedTemplateConfigMapSelectList,
) => ({
  type: CHANGE_TEMPLATE,
  edittedTemplateNamespace,
  edittedTemplateConfigMapName,
  edittedTemplateName,
  edittedTemplateYaml,
  edittedTemplateConfigMapSelectList,
});
