import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { YamlModalComponent } from './yaml-modal.component';
import { MatDialogModule } from '@angular/material/dialog';
import { FormModule, EditorModule } from 'kubeflow';

@NgModule({
  declarations: [YamlModalComponent],
  imports: [CommonModule, MatDialogModule, FormModule, EditorModule],
})
export class YamlModalModule {}
