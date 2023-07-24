import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatDividerModule } from '@angular/material/divider';

import { FormTrialTemplateComponent } from './trial-template.component';
import { FormModule, PopoverModule, EditorModule } from 'kubeflow';
import { ListKeyValueModule } from 'src/app/shared/list-key-value/list-key-value.module';
import { TrialParameterComponent } from './trial-parameter/trial-parameter.component';

@NgModule({
  declarations: [FormTrialTemplateComponent, TrialParameterComponent],
  imports: [
    CommonModule,
    FormModule,
    ListKeyValueModule,
    MatDividerModule,
    PopoverModule,
    EditorModule,
  ],
  exports: [FormTrialTemplateComponent],
})
export class FormTrialTemplateModule {}
