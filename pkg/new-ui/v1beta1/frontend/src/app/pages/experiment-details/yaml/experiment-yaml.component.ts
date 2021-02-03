import { Component, Input } from '@angular/core';
import { ExperimentK8s } from 'src/app/models/experiment.k8s.model';
import { dump } from 'js-yaml';

@Component({
  selector: 'app-experiment-yaml',
  templateUrl: './experiment-yaml.component.html',
  styleUrls: ['./experiment-yaml.component.scss'],
})
export class ExperimentYamlComponent {
  public yaml = '';

  @Input()
  set experimentJson(exp: ExperimentK8s) {
    this.yaml = dump(exp);
  }
}
