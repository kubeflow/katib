import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormAlgorithmComponent } from './algorithm.component';
import { FormModule } from 'kubeflow';

import { MatIconModule } from '@angular/material/icon';
import { MatRadioModule } from '@angular/material/radio';
import { FormAlgorithmSettingComponent } from './setting/setting.component';

@NgModule({
  declarations: [FormAlgorithmComponent, FormAlgorithmSettingComponent],
  imports: [CommonModule, FormModule, MatIconModule, MatRadioModule],
  exports: [FormAlgorithmComponent, FormAlgorithmSettingComponent],
})
export class FormAlgorithmModule {}
