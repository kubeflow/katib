import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  ConditionsTableModule,
  DetailsListModule,
  HeadingSubheadingRowModule,
} from 'kubeflow';
import { TrialMetricsModule } from './metrics/metrics.component.module';
import { TrialOverviewComponent } from './trial-overview.component';

@NgModule({
  declarations: [TrialOverviewComponent],
  imports: [
    CommonModule,
    ConditionsTableModule,
    DetailsListModule,
    HeadingSubheadingRowModule,
    TrialMetricsModule,
  ],
  exports: [TrialOverviewComponent],
})
export class TrialOverviewModule {}
