import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ListInputComponent } from './list-input.component';

import { MatIconModule } from '@angular/material/icon';

import { FormModule } from 'kubeflow';

@NgModule({
  declarations: [ListInputComponent],
  imports: [CommonModule, FormModule, MatIconModule],
  exports: [ListInputComponent],
})
export class ListInputModule {}
