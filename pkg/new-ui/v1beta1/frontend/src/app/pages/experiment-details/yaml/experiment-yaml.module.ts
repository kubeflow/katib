import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { ExperimentYamlComponent } from './experiment-yaml.component';
import { AceEditorModule } from 'ng2-ace-editor';

@NgModule({
  declarations: [ExperimentYamlComponent],
  imports: [CommonModule, AceEditorModule],
  exports: [ExperimentYamlComponent],
})
export class ExperimentYamlModule {}
