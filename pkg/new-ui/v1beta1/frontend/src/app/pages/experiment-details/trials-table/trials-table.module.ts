import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatTableModule } from '@angular/material/table';
import { MatIconModule } from '@angular/material/icon';
import { MatDialogModule } from '@angular/material/dialog';
import { NgxChartsModule } from '@swimlane/ngx-charts';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

import { TrialsTableComponent } from './trials-table.component';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TrialModalComponent } from './trial-modal/trial-modal.component';
import { TrialModalModule } from './trial-modal/trial-modal.module';

@NgModule({
  declarations: [TrialsTableComponent],
  imports: [
    CommonModule,
    MatTableModule,
    MatProgressSpinnerModule,
    MatDialogModule,
    MatIconModule,
    NgxChartsModule,
    MatTooltipModule,
    TrialModalModule,
  ],
  entryComponents: [TrialModalComponent],
  exports: [TrialsTableComponent],
})
export class TrialsTableModule {}
