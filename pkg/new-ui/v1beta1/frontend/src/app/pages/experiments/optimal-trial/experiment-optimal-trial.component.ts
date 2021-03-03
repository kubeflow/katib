import { Component } from '@angular/core';
import { Experiment } from 'src/app/models/experiment.model';
import { TableColumnComponent } from 'kubeflow/lib/resource-table/component-value/component-value.component';
import {
  KeyValuePair,
  trackByParamFn,
  parseOptimalMetric,
  parseOptimalParameters,
} from './utils';

@Component({
  selector: 'app-experiment-optimal-trial',
  templateUrl: './experiment-optimal-trial.component.html',
  styleUrls: ['./experiment-optimal-trial.component.scss'],
})
export class ExperimentOptimalTrialComponent implements TableColumnComponent {
  public params: KeyValuePair[] = [];
  public metrics: KeyValuePair[] = [];
  public trackByParamFn = trackByParamFn;

  set element(experiment: Experiment) {
    this.metrics = parseOptimalMetric(experiment);
    this.params = parseOptimalParameters(experiment);
  }
}
