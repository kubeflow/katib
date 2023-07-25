import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TrialLogsComponent } from './trial-logs.component';
import {
  HeadingSubheadingRowModule,
  KubeflowModule,
  LogsViewerModule,
} from 'kubeflow';

@NgModule({
  declarations: [TrialLogsComponent],
  imports: [
    CommonModule,
    KubeflowModule,
    HeadingSubheadingRowModule,
    LogsViewerModule,
  ],
  exports: [TrialLogsComponent],
})
export class TrialLogsModule {}
