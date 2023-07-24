import { Component, Input } from '@angular/core';
import { FormArray, FormControl, Validators, FormGroup } from '@angular/forms';

@Component({
  selector: 'app-list-key-value',
  templateUrl: './list-key-value.component.html',
  styleUrls: ['./list-key-value.component.scss'],
})
export class ListKeyValueComponent {
  @Input() header: string;
  @Input() addButtonText = 'Add key-value';
  @Input() keyLabel = 'Key';
  @Input() valueLabel = 'Value';
  @Input() formArray: FormArray;

  addCtrl() {
    this.formArray.push(
      new FormGroup({
        key: new FormControl('k', Validators.required),
        value: new FormControl('v', Validators.required),
      }),
    );
  }

  removeCtrl(i: number) {
    this.formArray.removeAt(i);
  }
}
