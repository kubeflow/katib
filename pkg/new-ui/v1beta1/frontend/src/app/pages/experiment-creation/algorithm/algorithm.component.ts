import { Component, OnInit, Input, OnDestroy } from '@angular/core';
import {
  AlgorithmNames,
  NasAlgorithmNames,
} from 'src/app/constants/algorithms-types.const';
import { FormArray, FormGroup, FormControl, Validators } from '@angular/forms';
import { Subscription } from 'rxjs';
import { AlgorithmSettingsMap } from 'src/app/constants/algorithms-settings.const';
import {
  AlgorithmsEnum,
  AlgorithmSettingType,
} from 'src/app/enumerations/algorithms.enum';

@Component({
  selector: 'app-form-algorithm',
  templateUrl: './algorithm.component.html',
  styleUrls: ['./algorithm.component.scss'],
})
export class FormAlgorithmComponent implements OnInit, OnDestroy {
  algorithmSettings: FormArray;
  algorithms: { [key: string]: string } = AlgorithmNames;
  algorithmHasSettings = false;

  private subscriptions: Subscription = new Subscription();

  @Input()
  algorithmForm: FormGroup;

  ngOnInit(): void {
    this.algorithmSettings = this.algorithmForm.get(
      'algorithmSettings',
    ) as FormArray;

    // set the list of algorithm settings once the form loads
    this.setAlgorithmSettings(this.algorithmForm.value.algorithm);

    this.subscriptions.add(
      this.algorithmForm.get('algorithm').valueChanges.subscribe(algo => {
        this.setAlgorithmSettings(algo);
      }),
    );

    this.subscriptions.add(
      this.algorithmForm.get('type').valueChanges.subscribe(type => {
        if (type === 'nas') {
          this.algorithms = NasAlgorithmNames;
          this.algorithmForm.get('algorithm').setValue(AlgorithmsEnum.ENAS);
          return;
        }

        this.algorithms = AlgorithmNames;
        this.algorithmForm.get('algorithm').setValue(AlgorithmsEnum.RANDOM);
      }),
    );
  }

  // form helpers
  setAlgorithmSettings(algo: string) {
    this.algorithmSettings.clear();
    this.algorithmHasSettings = AlgorithmSettingsMap[algo].length !== 0;

    // create the settings
    for (const setting of AlgorithmSettingsMap[algo]) {
      this.addSetting(
        setting.name,
        setting.value,
        setting.type,
        setting.values,
      );
    }
  }

  addSetting(
    name: string,
    value: any,
    type: AlgorithmSettingType,
    values: any[],
  ) {
    this.algorithmSettings.push(
      new FormGroup({
        name: new FormControl(name, Validators.required),
        value: new FormControl(value, []),
        type: new FormControl(type, Validators.required),
        values: new FormControl(values, Validators.required),
      }),
    );
  }

  ngOnDestroy(): void {
    this.subscriptions.unsubscribe();
  }
}
