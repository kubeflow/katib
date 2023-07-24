import { Component, OnInit, OnDestroy, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { FormGroup, FormControl, Validators, FormArray } from '@angular/forms';
import { Subscription } from 'rxjs';
import { createParameterGroup, createFeasibleSpaceGroup } from '../../utils';
import {
  FeasibleSpace,
  FeasibleSpaceMinMax,
  FeasibleSpaceList,
} from 'src/app/models/experiment.k8s.model';

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
    const tp = this.formGroup.get('parameterType').value;
    return tp === 'discrete' || tp === 'categorical';
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
      this.formGroup.get('parameterType').valueChanges.subscribe(type => {
        let fsMinMax: FeasibleSpaceMinMax = { min: '1', max: '64', step: '' };
        let fsList: FeasibleSpaceList = { list: [] };
        let fs: FeasibleSpace;

        fs = fsMinMax;
        if (this.isList) {
          fs = fsList;
        }

        this.formGroup.removeControl('feasibleSpace');
        this.formGroup.addControl(
          'feasibleSpace',
          createFeasibleSpaceGroup(type, fs),
        );
        return;
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
