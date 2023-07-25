import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatTableModule } from '@angular/material/table';
import { MatIconModule } from '@angular/material/icon';
import { MatDialogModule } from '@angular/material/dialog';

import { TrialsTableComponent } from './trials-table.component';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatButtonModule } from '@angular/material/button';
import { KubeflowModule } from 'kubeflow';
import { KfpRunComponent } from './kfp-run/kfp-run.component';
import { TrialDetailsModule } from './trial-details/trial-details.module';

@NgModule({
  declarations: [TrialsTableComponent, KfpRunComponent],
  imports: [
    CommonModule,
    MatTableModule,
    MatDialogModule,
    MatIconModule,
    MatTooltipModule,
    MatButtonModule,
    KubeflowModule,
    TrialDetailsModule,
  ],
  exports: [TrialsTableComponent],
})
export class TrialsTableModule {}
