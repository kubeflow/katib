import { Component, Input } from '@angular/core';
import { FormArray } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { AddParamModalComponent } from './add-modal/add-modal.component';
import { createParameterGroup } from '../utils';

@Component({
  selector: 'app-shared-params-list',
  templateUrl: './params-list.component.html',
  styleUrls: ['./params-list.component.scss'],
})
export class ParamsListComponent {
  @Input() paramsArray: FormArray;

  constructor(private dialog: MatDialog) {}

  removeParam(i: number) {
    this.paramsArray.removeAt(i);
  }

  addParam() {
    const newParamGroup = createParameterGroup({
      name: '',
      parameterType: 'int',
      feasibleSpace: {
        min: '1',
        max: '64',
        step: '1',
      },
    });

    const sub = this.dialog
      .open(AddParamModalComponent, { data: newParamGroup })
      .afterClosed()
      .subscribe(group => {
        sub.unsubscribe();

        if (group) {
          this.paramsArray.push(group);
        }
      });
  }
}
