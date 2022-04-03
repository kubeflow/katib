import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { FormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatCheckboxModule } from '@angular/material/checkbox';

import { PbtTabComponent } from './pbt-tab.component';

@NgModule({
  declarations: [PbtTabComponent],
  imports: [CommonModule, FormsModule, MatFormFieldModule, MatSelectModule, MatCheckboxModule],
  exports: [PbtTabComponent],
})
export class PbtTabModule {}
