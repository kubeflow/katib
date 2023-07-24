export interface ConfigMapTemplateResponse {
  Path: 'string';
  Yaml: 'string';
}

export interface ConfigMapBody {
  ConfigMapName: string;
  Templates: ConfigMapTemplateResponse[];
}

export interface ConfigMapResponse {
  ConfigMapNamespace: string;
  ConfigMaps: ConfigMapBody[];
}

export interface TrialTemplateResponse {
  Data: ConfigMapResponse[];
}
