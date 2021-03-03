import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DetailsListModule, HeadingSubheadingRowModule } from 'kubeflow';

import { ExperimentDetailsTabComponent } from './experiment-details-tab.component';

@NgModule({
  declarations: [ExperimentDetailsTabComponent],
  imports: [CommonModule, DetailsListModule, HeadingSubheadingRowModule],
  exports: [ExperimentDetailsTabComponent],
})
export class ExperimentDetailsTabModule {}
