import { Component, OnInit, Input } from '@angular/core';
import {
  FormArray,
  FormBuilder,
  FormGroup,
  FormControl,
  Validators,
} from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { AddParamModalComponent } from './add-modal/add-modal.component';

@Component({
  selector: 'app-shared-params-list',
  templateUrl: './params-list.component.html',
  styleUrls: ['./params-list.component.scss'],
})
export class ParamsListComponent implements OnInit {
  @Input() paramsArray: FormArray;

  constructor(private builder: FormBuilder, private dialog: MatDialog) {}

  ngOnInit() {}

  removeParam(i: number) {
    this.paramsArray.removeAt(i);
  }

  addParam() {
    const newParamGroup = new FormGroup({
      name: new FormControl('', Validators.required),
      type: new FormControl('int'),
      value: new FormGroup({
        min: new FormControl('1', Validators.required),
        max: new FormControl('64', Validators.required),
        step: new FormControl('', []),
      }),
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
