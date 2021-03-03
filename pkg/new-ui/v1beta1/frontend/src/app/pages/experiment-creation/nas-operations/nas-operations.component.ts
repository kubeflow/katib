import { Component, OnInit, Input } from '@angular/core';
import { FormArray } from '@angular/forms';
import { createNasOperationGroup } from 'src/app/shared/utils';

@Component({
  selector: 'app-form-nas-operations',
  templateUrl: './nas-operations.component.html',
  styleUrls: ['./nas-operations.component.scss'],
})
export class FormNasOperationsComponent {
  @Input() formArray: FormArray;

  constructor() {}

  addCtrl() {
    this.formArray.push(
      createNasOperationGroup({
        operationType: '',
        parameters: [],
      }),
    );
  }

  removeCtrl(i) {
    this.formArray.removeAt(i);
  }
}
