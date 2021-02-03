import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { TrialsGraphComponent } from './trials-graph.component';

@NgModule({
  declarations: [TrialsGraphComponent],
  imports: [CommonModule, MatProgressSpinnerModule],
  exports: [TrialsGraphComponent],
})
export class TrialsGraphModule {}
