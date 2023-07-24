import { K8sObject } from 'kubeflow';
import { V1Container } from '@kubernetes/client-node';

/*
 * K8s object definitions
 */
export const EXPERIMENT_KIND = 'Experiment';
export const EXPERIMENT_APIVERSION = 'kubeflow.org/v1beta1';

export interface ExperimentK8s extends K8sObject {
  spec?: ExperimentSpec;
  status?: ExperimentStatus;
}

export interface ExperimentSpec {
  parallelTrialCount?: number;
  maxTrialCount?: number;
  maxFailedTrialCount?: number;
  resumePolicy?: ResumePolicyType;
  objective?: ObjectiveSpec;
  algorithm?: AlgorithmSpec;
  earlyStopping?: AlgorithmSpec;
  parameters?: ParameterSpec[];
  metricsCollectorSpec?: MetricsCollectorSpec;
  trialTemplate?: TrialTemplateSpec;
  nasConfig?: NasConfig;
}

export type ResumePolicyType = 'Never' | 'LongRunning' | 'FromVolume';

export interface ObjectiveSpec {
  type: ObjectiveType;
  goal: number;
  objectiveMetricName: string;
  additionalMetricNames: string[];
  metricStrategies: MetricStrategy[];
}

export type ObjectiveType = 'maximize' | 'minimize';

export interface AlgorithmSpec {
  algorithmName: string;
  algorithmSettings: AlgorithmSetting[];
}

export interface AlgorithmSetting {
  name: string;
  value: string;
}

export interface MetricStrategy {
  name: string;
  value: string;
}
export interface FeasibleSpaceMinMax {
  max: string;
  min: string;
  step: string;
}

export interface FeasibleSpaceList {
  list: string[];
}

export interface ParameterSpec {
  name: string;
  parameterType: ParameterType;
  feasibleSpace: FeasibleSpace;
}

export type FeasibleSpace = FeasibleSpaceMinMax | FeasibleSpaceList;

export type ParameterType = 'int' | 'double' | 'discrete' | 'categorical';

export interface NasConfig {
  graphConfig?: GraphConfig;
  operations?: NasOperation[];
}

export interface GraphConfig {
  numLayers?: number;
  inputSizes?: number[];
  outputSizes?: number[];
}

export interface NasOperation {
  operationType: string;
  parameters: ParameterSpec[];
}

export interface MetricsCollectorSpec {
  source?: SourceSpec;
  collector?: CollectorSpec;
}

export interface SourceSpec {
  httpGet?: HttpGet;
  fileSystemPath?: FileSystemPath;
  filter?: FilterSpec;
}

export interface CollectorSpec {
  kind: CollectorKind;
  customCollector: V1Container;
}

export type CollectorKind =
  | 'StdOut'
  | 'File'
  | 'TensorFlowEvent'
  | 'PrometheusMetric'
  | 'Custom'
  | 'None';

export interface HttpGet {
  host?: string;
  port?: number;
  path?: string;
  scheme?: string;
  httpHeaders?: HttpHeader[];
}

export interface HttpHeader {
  name: string;
  value: string;
}

export interface FileSystemPath {
  path: string;
  kind: FileSystemKind;
}

export type FileSystemKind = 'Directory' | 'File';

export interface FilterSpec {
  metricsFormat?: string[];
}

export interface TrialTemplateSpec {
  retain: boolean;
  trialSource: K8sObject;
  trialParameters: TrialParameter[];
  primaryPodLabels: { [key: string]: string };
  primaryContainerName: string;
  successCondition: string;
  failureCondition: string;
}

export interface TrialParameter {
  name: string;
  description: string;
  reference: string;
}

/*
 * status
 */
interface ExperimentStatusCondition {
  type: string;
  status: boolean;
  reason: string;
  message: string;
  lastUpdateTime: string;
  lastTransitionTime: string;
}

interface CurrentOptimalTrial {
  bestTrialName: string;
  parameterAssignments: { name: string; value: number }[];
  observation: {
    metrics: {
      name: string;
      latest: number;
      min: number;
      max: string;
    }[];
  };
}

interface ExperimentStatus {
  startTime: string;
  completionTime: string;
  conditions: ExperimentStatusCondition[];
  currentOptimalTrial: CurrentOptimalTrial;
  succeededTrialList: string[];
  runningTrialList: string[];
  failedTrialList: string[];
  trials: number;
  trialsSucceeded: number;
}
