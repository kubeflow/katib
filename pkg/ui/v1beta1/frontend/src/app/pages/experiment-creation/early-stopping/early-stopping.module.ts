import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { EarlyStoppingComponent } from './early-stopping.component';
import { FormModule } from 'kubeflow';
import { FormAlgorithmModule } from '../algorithm/algorithm.module';

@NgModule({
  declarations: [EarlyStoppingComponent],
  imports: [CommonModule, FormModule, FormAlgorithmModule],
  exports: [EarlyStoppingComponent],
})
export class FormEarlyStoppingModule {}
