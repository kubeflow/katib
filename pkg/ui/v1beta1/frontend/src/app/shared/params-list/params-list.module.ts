import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ParamsListComponent } from './params-list.component';
import { ParameterComponent } from './parameter/parameter.component';
import { AddParamModalComponent } from './add-modal/add-modal.component';

import { MatIconModule } from '@angular/material/icon';
import { MatDividerModule } from '@angular/material/divider';
import { MatRadioModule } from '@angular/material/radio';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';

import { FormModule, PopoverModule, DetailsListModule } from 'kubeflow';
import { ListInputModule } from '../list-input/list-input.module';

@NgModule({
  declarations: [
    ParamsListComponent,
    ParameterComponent,
    AddParamModalComponent,
  ],
  imports: [
    CommonModule,
    FormModule,
    MatIconModule,
    MatDividerModule,
    MatRadioModule,
    PopoverModule,
    DetailsListModule,
    MatSlideToggleModule,
    ListInputModule,
  ],
  exports: [ParamsListComponent],
})
export class ParamsListModule {}
