import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormMetricsCollectorComponent } from './metrics-collector.component';
import { FormModule, EditorModule } from 'kubeflow';
import { ListKeyValueModule } from 'src/app/shared/list-key-value/list-key-value.module';

@NgModule({
  declarations: [FormMetricsCollectorComponent],
  imports: [CommonModule, FormModule, ListKeyValueModule, EditorModule],
  exports: [FormMetricsCollectorComponent],
})
export class FormMetricsCollectorModule {}
