import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormNasGraphComponent } from './nas-graph.component';
import { FormModule } from 'kubeflow';
import { ListInputModule } from 'src/app/shared/list-input/list-input.module';

@NgModule({
  declarations: [FormNasGraphComponent],
  imports: [CommonModule, FormModule, ListInputModule],
  exports: [FormNasGraphComponent],
})
export class FormNasGraphModule {}
