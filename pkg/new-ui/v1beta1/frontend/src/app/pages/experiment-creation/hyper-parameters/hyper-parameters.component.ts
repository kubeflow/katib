import { Component, Input } from '@angular/core';
import { FormArray } from '@angular/forms';

@Component({
  selector: 'app-form-hyper-parameters',
  templateUrl: './hyper-parameters.component.html',
  styleUrls: ['./hyper-parameters.component.scss'],
})
export class FormHyperParametersComponent {
  @Input() hyperParamsArray: FormArray;

  constructor() {}

  get combinations() {
    const params = this.hyperParamsArray.value;
    if (!params.length) {
      return 0;
    }

    let confs = 1;
    let currentConfs = 0;
    for (const param of params) {
      if (Array.isArray(param.value)) {
        currentConfs = param.value.length;
      } else {
        const min = param.value.min;
        const max = param.value.max;
        const step = param.value.step;

        currentConfs = Math.ceil((max - min) / step) + 1;
      }

      if (currentConfs === 0) {
        continue;
      }

      confs *= currentConfs;
    }

    return confs;
  }
}
