import { Component, Input } from '@angular/core';
import { FormArray, FormControl, Validators } from '@angular/forms';

@Component({
  selector: 'app-shared-list-input',
  templateUrl: './list-input.component.html',
  styleUrls: ['./list-input.component.scss'],
})
export class ListInputComponent {
  @Input() header: string;
  @Input() valueLabel = 'Value';
  @Input() formArray: FormArray;
  @Input() addValueText = 'Add value';
  @Input() requiredValue = true;

  addCtrl() {
    const validators = this.requiredValue ? Validators.required : [];
    this.formArray.push(new FormControl('', validators));
  }

  removeCtrl(i: number) {
    this.formArray.removeAt(i);
  }
}
