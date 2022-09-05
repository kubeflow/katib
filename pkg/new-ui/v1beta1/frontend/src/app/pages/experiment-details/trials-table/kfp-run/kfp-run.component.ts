import { Component } from '@angular/core';
import { TableColumnComponent } from 'kubeflow/lib/resource-table/component-value/component-value.component';

@Component({
  selector: 'app-kfp-run',
  templateUrl: './kfp-run.component.html',
  styleUrls: ['./kfp-run.component.scss'],
})
export class KfpRunComponent implements TableColumnComponent {
  constructor() {}

  trialKfpRunUrl: string = '';

  set element(experiment: any) {
    if (experiment?.['kfp run']) {
      this.trialKfpRunUrl = `/pipeline/#/runs/details/${experiment['kfp run']}`;
    } else {
      this.trialKfpRunUrl = '';
    }
  }
}
