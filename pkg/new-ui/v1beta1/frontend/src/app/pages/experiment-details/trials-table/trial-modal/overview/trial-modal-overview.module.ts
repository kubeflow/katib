import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  ConditionsTableModule,
  DetailsListModule,
  HeadingSubheadingRowModule,
} from 'kubeflow';
import { TrialModalMetricsModule } from './metrics/metrics.component.module';
import { TrialModalOverviewComponent } from './trial-modal-overview.component';

@NgModule({
  declarations: [TrialModalOverviewComponent],
  imports: [
    CommonModule,
    ConditionsTableModule,
    DetailsListModule,
    HeadingSubheadingRowModule,
    TrialModalMetricsModule,
  ],
  exports: [TrialModalOverviewComponent],
})
export class TrialModalOverviewModule {}
