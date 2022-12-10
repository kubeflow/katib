import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { TrialYamlComponent } from './trial-yaml.component';
import { EditorModule } from 'kubeflow';

@NgModule({
  declarations: [TrialYamlComponent],
  imports: [CommonModule, EditorModule],
  exports: [TrialYamlComponent],
})
export class TrialYamlModule {}
