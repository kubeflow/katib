import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KubeflowModule } from 'kubeflow';
import { ExperimentsComponent } from './experiments.component';
import { ExperimentOptimalTrialModule } from './optimal-trial/experiment-optimal-trial.module';

@NgModule({
  declarations: [ExperimentsComponent],
  imports: [CommonModule, ExperimentOptimalTrialModule, KubeflowModule],
  exports: [ExperimentsComponent],
})
export class ExperimentsModule {}
