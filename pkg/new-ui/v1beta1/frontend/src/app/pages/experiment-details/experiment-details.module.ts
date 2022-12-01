import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTabsModule } from '@angular/material/tabs';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import {
  TitleActionsToolbarModule,
  LoadingSpinnerModule,
  PanelModule,
} from 'kubeflow';

import { ExperimentDetailsComponent } from './experiment-details.component';
import { TrialsTableModule } from './trials-table/trials-table.module';
import { ExperimentOverviewModule } from './overview/experiment-overview.module';
import { ExperimentDetailsTabModule } from './details/experiment-details-tab.module';
import { ExperimentYamlModule } from './yaml/experiment-yaml.module';
import { TrialsGraphEchartsModule } from './trials-graph-echarts/trials-graph-echarts.module';
import { KubeflowModule } from 'kubeflow';

@NgModule({
  declarations: [ExperimentDetailsComponent],
  imports: [
    CommonModule,
    TrialsTableModule,
    MatButtonModule,
    MatTabsModule,
    MatIconModule,
    LoadingSpinnerModule,
    PanelModule,
    ExperimentOverviewModule,
    ExperimentDetailsTabModule,
    MatProgressSpinnerModule,
    ExperimentYamlModule,
    TitleActionsToolbarModule,
    TrialsGraphEchartsModule,
    KubeflowModule,
  ],
  exports: [ExperimentDetailsComponent],
})
export class ExperimentDetailsModule {}
