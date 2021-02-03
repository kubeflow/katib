import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { MatIconModule } from '@angular/material/icon';
import { MatDividerModule } from '@angular/material/divider';
import { MatRadioModule } from '@angular/material/radio';

import { FormModule, PopoverModule, DetailsListModule } from 'kubeflow';

import { FormHyperParametersComponent } from './hyper-parameters.component';
import { ParamsListModule } from 'src/app/shared/params-list/params-list.module';

@NgModule({
  declarations: [FormHyperParametersComponent],
  imports: [
    CommonModule,
    FormModule,
    MatIconModule,
    MatDividerModule,
    MatRadioModule,
    PopoverModule,
    DetailsListModule,
    ParamsListModule,
  ],
  exports: [FormHyperParametersComponent],
})
export class FormHyperParametersModule {}
