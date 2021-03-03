import { Component, Input, Inject } from '@angular/core';
import { load, dump } from 'js-yaml';
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
  ) {
    this.yaml = dump(data);
  }

  save() {
    this.dialogRef.close(load(this.yaml));
  }

  close() {
    this.dialogRef.close();
  }
}
