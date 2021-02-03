import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import {
  NamespaceSelectModule,
  ResourceTableModule,
  ConfirmDialogModule,
} from 'kubeflow';

import { ExperimentsComponent } from './experiments.component';
import { ExperimentOptimalTrialModule } from './optimal-trial/experiment-optimal-trial.module';

@NgModule({
  declarations: [ExperimentsComponent],
  imports: [
    CommonModule,
    NamespaceSelectModule,
    ResourceTableModule,

    ConfirmDialogModule,
    ExperimentOptimalTrialModule,
    MatCardModule,
    ConfirmDialogModule,
    MatCardModule,
  ],
  exports: [ExperimentsComponent],
})
export class ExperimentsModule {}
