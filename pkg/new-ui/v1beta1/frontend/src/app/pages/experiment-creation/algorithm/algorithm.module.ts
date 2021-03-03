import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormAlgorithmComponent } from './algorithm.component';
import { FormModule } from 'kubeflow';

import { MatIconModule } from '@angular/material/icon';
import { MatRadioModule } from '@angular/material/radio';

@NgModule({
  declarations: [FormAlgorithmComponent],
  imports: [CommonModule, FormModule, MatIconModule, MatRadioModule],
  exports: [FormAlgorithmComponent],
})
export class FormAlgorithmModule {}
