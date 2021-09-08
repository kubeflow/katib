import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ExperimentOptimalTrialComponent } from './experiment-optimal-trial.component';
import { PopoverModule, DetailsListModule } from 'kubeflow';
import { MatDividerModule } from '@angular/material/divider';

@NgModule({
  declarations: [ExperimentOptimalTrialComponent],
  imports: [CommonModule, PopoverModule, DetailsListModule, MatDividerModule],
})
export class ExperimentOptimalTrialModule {}
