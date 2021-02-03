import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormNasOperationsComponent } from './nas-operations.component';
import { OperationComponent } from './operation/operation.component';
import { FormModule } from 'kubeflow';
import { ParamsListModule } from 'src/app/shared/params-list/params-list.module';
import { MatDividerModule } from '@angular/material/divider';

@NgModule({
  declarations: [FormNasOperationsComponent, OperationComponent],
  imports: [CommonModule, FormModule, ParamsListModule, MatDividerModule],
  exports: [FormNasOperationsComponent],
})
export class FormNasOperationsModule {}
