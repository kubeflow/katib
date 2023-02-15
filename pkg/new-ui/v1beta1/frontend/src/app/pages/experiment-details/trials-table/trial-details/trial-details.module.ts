import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatTableModule } from '@angular/material/table';
import { MatIconModule } from '@angular/material/icon';
import { MatDialogModule } from '@angular/material/dialog';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatTabsModule } from '@angular/material/tabs';
import { TrialOverviewModule } from './trial-overview/trial-overview.module';
import { TrialDetailsComponent } from './trial-details.component';
import { TrialYamlModule } from './trial-yaml/trial-yaml.module';

import {
  TitleActionsToolbarModule,
  LoadingSpinnerModule,
  PanelModule,
} from 'kubeflow';
import { NgxEchartsModule } from 'ngx-echarts';
import { TrialLogsModule } from './trial-logs/trial-logs.module';

@NgModule({
  declarations: [TrialDetailsComponent],
  imports: [
    TrialOverviewModule,
    CommonModule,
    MatTableModule,
    MatProgressSpinnerModule,
    MatDialogModule,
    MatIconModule,
    MatTooltipModule,
    MatTabsModule,
    TitleActionsToolbarModule,
    LoadingSpinnerModule,
    PanelModule,
    TitleActionsToolbarModule,
    TrialYamlModule,
    NgxEchartsModule.forRoot({
      echarts: () => import('echarts'),
    }),
    TrialLogsModule,
  ],
  exports: [TrialDetailsComponent],
})
export class TrialDetailsModule {}
