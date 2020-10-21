// Here we store global constants for React frontend

export const EXPERIMENT_TYPE_HP = 'hp';
export const EXPERIMENT_TYPE_NAS = 'nas';

export const MC_KIND_STDOUT = 'StdOut';
export const MC_KIND_FILE = 'File';
export const MC_KIND_TENSORFLOW_EVENT = 'TensorFlowEvent';
export const MC_KIND_PROMETHEUS = 'PrometheusMetric';
export const MC_KIND_CUSTOM = 'Custom';
export const MC_KIND_NONE = 'None';

export const MC_FILE_SYSTEM_KIND_FILE = 'File';
export const MC_FILE_SYSTEM_KIND_DIRECTORY = 'Directory';
export const MC_FILE_SYSTEM_NO_KIND = 'No File System';

export const MC_HTTP_GET_HTTP_SCHEME = 'HTTP';

export const MC_PROMETHEUS_DEFAULT_PORT = 8080;
export const MC_PROMETHEUS_DEFAULT_PATH = '/metrics';

export const LINK_HP_CREATE = '/katib/hp';
export const LINK_HP_MONITOR = '/katib/hp_monitor';
export const LINK_NAS_CREATE = '/katib/nas';
export const LINK_NAS_MONITOR = '/katib/nas_monitor';
export const LINK_TRIAL_TEMPLATE = '/katib/trial';

export const GENERAL_MODULE = 'general';
export const HP_CREATE_MODULE = 'hpCreate';
export const HP_MONITOR_MODULE = 'hpMonitor';
export const NAS_CREATE_MODULE = 'nasCreate';
export const NAS_MONITOR_MODULE = 'nasMonitor';
export const TEMPLATE_MODULE = 'template';

export const TEMPLATE_SOURCE_CONFIG_MAP = 'ConfigMap';
export const TEMPLATE_SOURCE_YAML = 'YAML';
