import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormObjectiveComponent } from './objective.component';
import { FormModule } from 'kubeflow';
import { MatIconModule } from '@angular/material/icon';
import { ListInputModule } from 'src/app/shared/list-input/list-input.module';
import { MatDividerModule } from '@angular/material/divider';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { MatCheckboxModule } from '@angular/material/checkbox';

@NgModule({
  declarations: [FormObjectiveComponent],
  imports: [
    CommonModule,
    FormModule,
    MatIconModule,
    ListInputModule,
    MatDividerModule,
    MatCheckboxModule,
  ],
  exports: [FormObjectiveComponent],
})
export class FormObjectiveModule {}
