import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { ExperimentYamlComponent } from './experiment-yaml.component';
import { EditorModule } from 'kubeflow';

@NgModule({
  declarations: [ExperimentYamlComponent],
  imports: [CommonModule, EditorModule],
  exports: [ExperimentYamlComponent],
})
export class ExperimentYamlModule {}
