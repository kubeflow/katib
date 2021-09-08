import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { YamlModalComponent } from './yaml-modal.component';
import { MatDialogModule } from '@angular/material/dialog';
import { FormModule } from 'kubeflow';
import { AceEditorModule } from 'ng2-ace-editor';

@NgModule({
  declarations: [YamlModalComponent],
  imports: [CommonModule, MatDialogModule, FormModule, AceEditorModule],
})
export class YamlModalModule {}
