import { Component, OnInit, Input, OnDestroy } from '@angular/core';
import { FormGroup, FormArray, FormControl, Validators } from '@angular/forms';
import { EarlyStoppingAlgorithmNames } from 'src/app/constants/algorithms-types.const';
import { Subscription } from 'rxjs';
import { EarlyStoppingSettingsMap } from 'src/app/constants/algorithms-settings.const';
import { AlgorithmSettingType } from 'src/app/enumerations/algorithms.enum';

@Component({
  selector: 'app-form-early-stopping',
  templateUrl: './early-stopping.component.html',
  styleUrls: ['./early-stopping.component.scss'],
})
export class EarlyStoppingComponent implements OnInit, OnDestroy {
  algorithmSettings: FormArray;
  algorithms: { [key: string]: string } = EarlyStoppingAlgorithmNames;
  algorithmHasSettings = false;

  private subscriptions: Subscription = new Subscription();

  @Input() formGroup: FormGroup;

  constructor() {}

  ngOnInit(): void {
    this.algorithmSettings = this.formGroup.get(
      'algorithmSettings',
    ) as FormArray;

    this.subscriptions.add(
      this.formGroup.get('algorithmName').valueChanges.subscribe(algo => {
        this.algorithmSettings.clear();
        this.algorithmHasSettings = EarlyStoppingSettingsMap[algo].length !== 0;

        // create the settings
        for (const setting of EarlyStoppingSettingsMap[algo]) {
          this.addSetting(
            setting.name,
            setting.value,
            setting.type,
            setting.values,
          );
        }
      }),
    );
  }

  // form helpers
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
