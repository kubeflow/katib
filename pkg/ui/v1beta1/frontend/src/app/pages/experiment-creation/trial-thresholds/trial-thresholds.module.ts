import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormTrialThresholdsComponent } from './trial-thresholds.component';
import { FormModule } from 'kubeflow';

@NgModule({
  declarations: [FormTrialThresholdsComponent],
  imports: [CommonModule, FormModule],
  exports: [FormTrialThresholdsComponent],
})
export class FormTrialThresholdsModule {}
