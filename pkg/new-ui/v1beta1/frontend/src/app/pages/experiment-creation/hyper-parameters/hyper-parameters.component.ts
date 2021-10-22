import { Component, Input } from '@angular/core';
import { FormArray } from '@angular/forms';
import {
  ParameterSpec,
  FeasibleSpaceList,
  FeasibleSpaceMinMax,
} from 'src/app/models/experiment.k8s.model';

@Component({
  selector: 'app-form-hyper-parameters',
  templateUrl: './hyper-parameters.component.html',
  styleUrls: ['./hyper-parameters.component.scss'],
})
export class FormHyperParametersComponent {
  @Input() hyperParamsArray: FormArray;

  constructor() {}

  get combinations() {
    const params = this.hyperParamsArray.value as ParameterSpec[];
    if (!params.length) {
      return 0;
    }

    let confs = 1;
    let currentConfs = 0;
    for (const param of params) {
      if (
        param.parameterType === 'discrete' ||
        param.parameterType === 'categorical'
      ) {
        currentConfs = (param.feasibleSpace as FeasibleSpaceList).list.length;
      } else {
        const fs = param.feasibleSpace as FeasibleSpaceMinMax;

        const min = fs.min;
        const max = fs.max;
        const step = fs.step;

        // don't calculate the combinations is step is omitted
        if (step === '') {
          return null;
        }

        try {
          parseFloat(min);
          parseFloat(max);
          parseFloat(step);
        } catch (e) {
          console.log('Could not convert min/max/step to number');
          return null;
        }

        currentConfs =
          Math.abs(
            Math.ceil((parseFloat(max) - parseFloat(min)) / parseFloat(step)),
          ) + 1;
      }

      if (currentConfs === 0) {
        continue;
      }

      confs *= currentConfs;
    }

    return confs;
  }
}
