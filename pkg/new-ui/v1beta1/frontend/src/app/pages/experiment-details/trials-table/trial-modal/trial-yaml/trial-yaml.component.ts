import { Component, Input } from '@angular/core';
import { TrialK8s } from 'src/app/models/trial.k8s.model';
import { dump } from 'js-yaml';

@Component({
  selector: 'app-trial-yaml',
  templateUrl: './trial-yaml.component.html',
  styleUrls: ['./trial-yaml.component.scss'],
})
export class TrialYamlComponent {
  public yaml = '';

  @Input()
  set trialJson(trial: TrialK8s) {
    this.yaml = dump(trial);
  }
}
