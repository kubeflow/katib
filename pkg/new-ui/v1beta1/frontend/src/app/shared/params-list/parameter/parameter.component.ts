import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { AbstractControl } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { AddParamModalComponent } from '../add-modal/add-modal.component';
import {
  FeasibleSpaceMinMax,
  ParameterType,
} from 'src/app/models/experiment.k8s.model';

@Component({
  selector: 'app-shared-parameter',
  templateUrl: './parameter.component.html',
  styleUrls: ['./parameter.component.scss'],
})
export class ParameterComponent implements OnInit {
  @Input() paramCtrl: AbstractControl;
  @Output() delete = new EventEmitter<boolean>();

  constructor(private dialog: MatDialog) {}

  ngOnInit() {}

  get isListValue() {
    if (!this.paramCtrl) {
      return false;
    }

    return this.type === 'discrete' || this.type === 'categorical';
  }

  get name(): string {
    return this.paramCtrl.get('name').value;
  }

  get type(): ParameterType {
    return this.paramCtrl.get('parameterType').value;
  }

  get min() {
    return this.paramCtrl.get('feasibleSpace').value.min;
  }

  get max() {
    return this.paramCtrl.get('feasibleSpace').value.max;
  }

  get step() {
    return this.paramCtrl.get('feasibleSpace').value.step;
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
    return this.paramCtrl.get('feasibleSpace').value.list;
  }

  get listStr() {
    return this.listValue.join(', ');
  }

  editParam() {
    const sub = this.dialog
      .open(AddParamModalComponent, { data: this.paramCtrl })
      .afterClosed()
      .subscribe(group => {
        sub.unsubscribe();

        if (group) {
          this.paramCtrl.get('name').setValue(group.get('name').value);
          this.paramCtrl
            .get('parameterType')
            .setValue(group.get('parameterType').value);
          this.paramCtrl
            .get('feasibleSpace')
            .setValue(group.get('feasibleSpace').value);
        }
      });
  }
}
