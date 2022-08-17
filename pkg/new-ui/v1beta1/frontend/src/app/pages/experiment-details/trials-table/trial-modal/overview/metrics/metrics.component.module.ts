import { NgModule } from '@angular/core';
import { ConditionsTableModule, DetailsListModule } from 'kubeflow';
import { TrialModalMetricsComponent } from './metrics.component';

@NgModule({
  declarations: [TrialModalMetricsComponent],
  imports: [ConditionsTableModule, DetailsListModule],
  exports: [TrialModalMetricsComponent],
})
export class TrialModalMetricsModule {}
