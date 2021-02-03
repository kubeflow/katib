import { Component, OnInit, OnDestroy, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { FormGroup, FormControl, Validators, FormArray } from '@angular/forms';
import { Subscription } from 'rxjs';
import { createParameterGroup } from '../../utils';

@Component({
  selector: 'app-add-param-modal',
  templateUrl: './add-modal.component.html',
  styleUrls: ['./add-modal.component.scss'],
})
export class AddParamModalComponent implements OnInit, OnDestroy {
  originalFormGroup: FormGroup;
  formGroup: FormGroup;

  subs = new Subscription();

  get isList(): boolean {
    const tp = this.formGroup.get('type').value;
    return tp === 'discrete' || tp === 'categorical';
  }

  get inputStep(): number {
    const tp = this.formGroup.get('type').value;
    if (tp === 'int') {
      return 1;
    }

    return 0.01;
  }

  constructor(
    private dialog: MatDialogRef<AddParamModalComponent>,
    @Inject(MAT_DIALOG_DATA) public data: FormGroup,
  ) {
    this.originalFormGroup = data;
    this.formGroup = createParameterGroup(data.value);
  }

  ngOnInit() {
    this.subs.add(
      this.formGroup.get('type').valueChanges.subscribe(type => {
        if (this.isList) {
          this.formGroup.removeControl('value');
          this.formGroup.addControl(
            'value',
            new FormArray([], Validators.required),
          );
          return;
        }

        this.formGroup.removeControl('value');
        this.formGroup.addControl(
          'value',
          new FormGroup({
            min: new FormControl('1', Validators.required),
            max: new FormControl('64', Validators.required),
            step: new FormControl('', []),
          }),
        );
      }),
    );
  }

  ngOnDestroy() {
    this.subs.unsubscribe();
  }

  save() {
    this.dialog.close(this.formGroup);
  }

  cancel() {
    this.dialog.close();
  }
}
