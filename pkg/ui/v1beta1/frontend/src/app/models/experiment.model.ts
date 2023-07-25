import { Params } from '@angular/router';

/*
 * UI relative types
 */
interface ExperimentCondition {
  type: string;
  status: string;
  reason: string;
  message: string;
  lastUpdateTime: string;
  lastTransitionTime: string;
}

export interface TrialParameterAssignments {
  name: string;
  value: string;
}

export interface TrialObservationMetrics {
  name: string;
  min: number;
  max: number;
  latest: number;
}

interface ExperimentCurrentOptimalTrial {
  bestTrialName: string;
  parameterAssignments: TrialParameterAssignments[];
  observation: {
    metrics: TrialObservationMetrics[];
  };
}

export interface Experiment {
  name: string;
  namespace: string;
  status: string;
  startTime: string;
  conditions: ExperimentCondition[];
  currentOptimalTrial: ExperimentCurrentOptimalTrial;
  runningTrialList: string[];
  failedTrialList: string[];
  succeededTrialList: string[];
  trials: number;
  trialsSucceeded: number;
  trialsFailed: number;
  trialsRunning: number;
}

export type Experiments = Experiment[];

export interface ExperimentProcessed extends Experiment {
  link: {
    text: string;
    url: string;
    queryParams?: Params | null;
  };
}

export type ExperimentsProcessed = ExperimentProcessed[];
