import { K8sObject } from 'kubeflow';
import { V1Container } from '@kubernetes/client-node';

/*
 * K8s object definitions
 */
export const TRIAL_KIND = 'Trial';
export const TRIAL_APIVERSION = 'kubeflow.org/v1beta1';

export interface TrialK8s extends K8sObject {
  spec?: TrialSpec;
  status?: TrialStatus;
}

export interface TrialSpec {
  metricsCollector: MetricsCollector;
  objective: Objective;
  parameterAssignments: { name: string; value: number }[];
  primaryContainerName: string;
  successCondition: string;
  failureCondition: string;
  runSpec: K8sObject;
}

export interface MetricsCollector {
  collector?: CollectorSpec;
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

export interface Objective {
  type: ObjectiveType;
  goal: number;
  objectiveMetricName: string;
  additionalMetricNames: string[];
  metricStrategies: MetricStrategy[];
}

export type ObjectiveType = 'maximize' | 'minimize';

export interface MetricStrategy {
  name: string;
  value: string;
}

export interface RunSpec {}

/*
 * status
 */

interface TrialStatus {
  startTime: string;
  completionTime: string;
  conditions: TrialStatusCondition[];
  observation: {
    metrics: {
      name: string;
      latest: string;
      min: string;
      max: string;
    }[];
  };
}

interface TrialStatusCondition {
  type: string;
  status: string;
  reason: string;
  message: string;
  lastUpdateTime: string;
  lastTransitionTime: string;
}
