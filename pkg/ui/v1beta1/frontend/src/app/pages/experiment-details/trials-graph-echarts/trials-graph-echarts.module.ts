import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { TrialsGraphEchartsComponent } from './trials-graph-echarts.component';
import { NgxEchartsModule } from 'ngx-echarts';

@NgModule({
  declarations: [TrialsGraphEchartsComponent],
  imports: [
    CommonModule,
    MatProgressSpinnerModule,
    NgxEchartsModule.forRoot({
      echarts: () => import('echarts'),
    }),
  ],
  exports: [TrialsGraphEchartsComponent],
})
export class TrialsGraphEchartsModule {}
