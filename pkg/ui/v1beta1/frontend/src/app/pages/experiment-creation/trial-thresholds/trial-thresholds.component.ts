import { Component, Input } from '@angular/core';
import { FormGroup } from '@angular/forms';

@Component({
  selector: 'app-form-trial-thresholds',
  templateUrl: './trial-thresholds.component.html',
  styleUrls: ['./trial-thresholds.component.scss'],
})
export class FormTrialThresholdsComponent {
  @Input()
  formGroup: FormGroup;
}
