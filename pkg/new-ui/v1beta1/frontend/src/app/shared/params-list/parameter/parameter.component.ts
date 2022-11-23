import { Component, Input, Output, EventEmitter } from '@angular/core';
import { FormGroup } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { AddParamModalComponent } from '../add-modal/add-modal.component';
import { ParameterType } from 'src/app/models/experiment.k8s.model';

@Component({
  selector: 'app-shared-parameter',
  templateUrl: './parameter.component.html',
  styleUrls: ['./parameter.component.scss'],
})
export class ParameterComponent {
  @Input() paramFormGroup: FormGroup;
  @Output() delete = new EventEmitter<boolean>();

  constructor(private dialog: MatDialog) {}

  get isListValue() {
    if (!this.paramFormGroup) {
      return false;
    }

    return this.type === 'discrete' || this.type === 'categorical';
  }

  get name(): string {
    return this.paramFormGroup.get('name').value;
  }

  get type(): ParameterType {
    return this.paramFormGroup.get('parameterType').value;
  }

  get min() {
    return this.paramFormGroup.get('feasibleSpace').value.min;
  }

  get max() {
    return this.paramFormGroup.get('feasibleSpace').value.max;
  }

  get step() {
    return this.paramFormGroup.get('feasibleSpace').value.step;
  }

  get stepSign() {
    if (this.step > 0) {
      return ', +';
    }

    if (this.step < 0) {
      return ', ';
    }

    return '';
  }

  get listValue(): any[] {
    return this.paramFormGroup.get('feasibleSpace').value.list;
  }

  get listStr() {
    return this.listValue.join(', ');
  }

  editParam() {
    const sub = this.dialog
      .open(AddParamModalComponent, { data: this.paramFormGroup })
      .afterClosed()
      .subscribe(group => {
        sub.unsubscribe();

        if (group) {
          this.paramFormGroup.get('name').setValue(group.get('name').value);
          this.paramFormGroup
            .get('parameterType')
            .setValue(group.get('parameterType').value);
          this.paramFormGroup.setControl(
            'feasibleSpace',
            group.get('feasibleSpace'),
          );
        }
      });
  }
}
