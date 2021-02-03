import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { AbstractControl } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { AddParamModalComponent } from '../add-modal/add-modal.component';

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

    return Array.isArray(this.paramCtrl.get('value').value);
  }

  get name() {
    return this.paramCtrl.get('name').value;
  }

  get type() {
    return this.paramCtrl.get('type').value;
  }

  get min() {
    return this.paramCtrl.get('value').value.min;
  }

  get max() {
    return this.paramCtrl.get('value').value.max;
  }

  get step() {
    return this.paramCtrl.get('value').value.step;
  }

  get stepSign() {
    if (this.step >= 0) {
      return ', +';
    }

    return '';
  }

  get listValue(): any[] {
    return this.paramCtrl.get('value').value;
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
          this.paramCtrl.get('type').setValue(group.get('type').value);
          this.paramCtrl.get('name').setValue(group.get('name').value);
          this.paramCtrl.get('value').setValue(group.get('value').value);
        }
      });
  }
}
