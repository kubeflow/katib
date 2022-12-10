import { Component, Input, Inject } from '@angular/core';
import { load, dump } from 'js-yaml';
import { SnackBarService, SnackType } from 'kubeflow';
import {
  MatDialog,
  MatDialogRef,
  MAT_DIALOG_DATA,
} from '@angular/material/dialog';

@Component({
  selector: 'app-yaml-modal',
  templateUrl: './yaml-modal.component.html',
  styleUrls: ['./yaml-modal.component.scss'],
})
export class YamlModalComponent {
  public yaml = '';

  constructor(
    public dialogRef: MatDialogRef<YamlModalComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
    private snack: SnackBarService,
  ) {
    this.yaml = dump(data);
  }

  save() {
    try {
      this.dialogRef.close(load(this.yaml));
    } catch (e) {
      this.snack.open(`${e.reason}`, SnackType.Error, 4000);
    }
  }

  close() {
    this.dialogRef.close();
  }
}
