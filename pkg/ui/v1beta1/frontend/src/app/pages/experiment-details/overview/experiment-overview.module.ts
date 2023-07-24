import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ConditionsTableModule, DetailsListModule } from 'kubeflow';

import { ExperimentOverviewComponent } from './experiment-overview.component';

@NgModule({
  declarations: [ExperimentOverviewComponent],
  imports: [CommonModule, ConditionsTableModule, DetailsListModule],
  exports: [ExperimentOverviewComponent],
})
export class ExperimentOverviewModule {}
