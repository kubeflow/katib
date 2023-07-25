import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormMetadataComponent } from './metadata.component';
import { FormModule } from 'kubeflow';

@NgModule({
  declarations: [FormMetadataComponent],
  imports: [CommonModule, FormModule],
  exports: [FormMetadataComponent],
})
export class FormMetadataModule {}
