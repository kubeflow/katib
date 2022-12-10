import { NgModule } from '@angular/core';
import { ConditionsTableModule, DetailsListModule } from 'kubeflow';
import { TrialMetricsComponent } from './metrics.component';

@NgModule({
  declarations: [TrialMetricsComponent],
  imports: [ConditionsTableModule, DetailsListModule],
  exports: [TrialMetricsComponent],
})
export class TrialMetricsModule {}
