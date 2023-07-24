import { Experiment } from 'src/app/models/experiment.model';
import lowerCase from 'lodash-es/lowerCase';
import { numberToExponential } from 'src/app/shared/utils';

export interface KeyValuePair {
  name: string;
  value: string;
}

export function parseOptimalMetric(exp: Experiment): KeyValuePair[] {
  if (
    !exp ||
    !exp.currentOptimalTrial ||
    !exp.currentOptimalTrial.observation.metrics
  ) {
    return [];
  }

  return exp.currentOptimalTrial.observation.metrics.map(metric => {
    const value = isNaN(+metric.latest)
      ? metric.latest.toString()
      : numberToExponential(Number(metric.latest), 6);

    const name =
      lowerCase(metric.name).charAt(0).toUpperCase() +
      lowerCase(metric.name).slice(1);

    return { name, value };
  });
}

export function parseOptimalParameters(exp: Experiment): KeyValuePair[] {
  if (
    !exp ||
    !exp.currentOptimalTrial ||
    !exp.currentOptimalTrial.parameterAssignments
  ) {
    return [];
  }

  return exp.currentOptimalTrial.parameterAssignments.map(param => {
    const value = isNaN(+param.value)
      ? param.value.toString()
      : numberToExponential(+param.value, 6);

    const name =
      lowerCase(param.name).charAt(0).toUpperCase() +
      lowerCase(param.name).slice(1);

    return { name, value };
  });
}

export function trackByParamFn(index: number, row: KeyValuePair) {
  return `${row.name}: ${row.value}`;
}
