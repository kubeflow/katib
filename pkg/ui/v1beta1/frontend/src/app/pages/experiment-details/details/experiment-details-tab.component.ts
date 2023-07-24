import {
  ChangeDetectionStrategy,
  Component,
  Input,
  OnChanges,
} from '@angular/core';

import { ChipDescriptor, ListEntry } from 'kubeflow';
import {
  ExperimentK8s,
  ExperimentSpec,
  FeasibleSpaceList,
  FeasibleSpaceMinMax,
} from 'src/app/models/experiment.k8s.model';

@Component({
  selector: 'app-experiment-details-tab',
  templateUrl: './experiment-details-tab.component.html',
  styleUrls: ['./experiment-details-tab.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ExperimentDetailsTabComponent implements OnChanges {
  @Input()
  experiment: ExperimentK8s;
  experimentAlgorithmList: ListEntry[] = [];
  metricsCollectorSpecList: ListEntry[] = [];
  objectiveAdditionalMetrics: string;
  parameters: { chips: ChipDescriptor[]; key: string }[] = [];

  ngOnChanges(): void {
    if (!this.experiment) {
      return;
    }

    const objective = this.experiment.spec.objective;
    let metrics = 'No additional metrics';
    if (
      objective.additionalMetricNames &&
      objective.additionalMetricNames.length
    ) {
      metrics = objective.additionalMetricNames.join(', ');
    }

    this.objectiveAdditionalMetrics = metrics;
    this.parameters = this.generateParametersList(this.experiment.spec);
  }

  generateParametersList(
    spec: ExperimentSpec,
  ): { chips: ChipDescriptor[]; key: string }[] {
    const { parameters } = spec;

    return parameters.map(parameter => {
      const feasibleSpaceList = parameter.feasibleSpace as FeasibleSpaceList;
      const feasibleSpaceMinMax =
        parameter.feasibleSpace as FeasibleSpaceMinMax;

      const chips: ChipDescriptor[] = [
        {
          value: `Parameter type: ${parameter.parameterType}`,
          color: 'primary',
        },
      ];

      if (!!feasibleSpaceList.list) {
        chips.push({
          value: `${feasibleSpaceList.list.join(', ')}`,
          color: 'primary',
        });
      } else {
        chips.push({
          value: `Min: ${feasibleSpaceMinMax.min}`,
          color: 'primary',
        });
        chips.push({
          value: `Max: ${feasibleSpaceMinMax.max}`,
          color: 'primary',
        });
      }

      return {
        key: parameter.name,
        chips,
      };
    });
  }
}
