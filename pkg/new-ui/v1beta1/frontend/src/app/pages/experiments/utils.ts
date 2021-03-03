import { Experiments, Experiment } from '../../models/experiment.model';
import { Status, STATUS_TYPE, ChipsListValue, ChipDescriptor } from 'kubeflow';
import { StatusEnum } from 'src/app/enumerations/status.enum';

export function parseStatus(exp: Experiment): Status {
  const statusCol = {
    phase: STATUS_TYPE.ERROR,
    state: '',
    message: exp.status,
  };

  if (exp.status === StatusEnum.SUCCEEDED) {
    statusCol.phase = STATUS_TYPE.READY;
  }

  if (exp.status === StatusEnum.FAILED) {
    statusCol.phase = STATUS_TYPE.WARNING;
  }

  if (exp.status === StatusEnum.CREATED) {
    statusCol.phase = STATUS_TYPE.UNAVAILABLE;
  }

  if (exp.status === StatusEnum.RUNNING) {
    statusCol.phase = STATUS_TYPE.UNAVAILABLE;
  }

  return statusCol;
}

export function parseTotalTrials(exp: Experiment): number {
  return exp.trials ? exp.trials : 0;
}

export function parseSucceededTrials(exp: Experiment): number {
  return exp.trialsSucceeded ? exp.trialsSucceeded : 0;
}

export function parseRunningTrials(exp: Experiment): number {
  return exp.trialsRunning ? exp.trialsRunning : 0;
}

export function parseFailedTrials(exp: Experiment): number {
  return exp.trialsFailed ? exp.trialsFailed : 0;
}
