import {
  ChangeDetectionStrategy,
  Component,
  Input,
  OnChanges,
} from '@angular/core';
import { ChipDescriptor, getCondition } from 'kubeflow';
import { ExperimentK8s } from 'src/app/models/experiment.k8s.model';
import { ObjectiveTypeEnum } from 'src/app/enumerations/objective-type.enum';
import { StatusEnum } from 'src/app/enumerations/status.enum';
import { numberToExponential } from 'src/app/shared/utils';

@Component({
  selector: 'app-experiment-overview',
  templateUrl: './experiment-overview.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ExperimentOverviewComponent implements OnChanges {
  status: string;
  statusIcon: string;
  bestTrialName: string;
  bestTrialPerformance: ChipDescriptor[];
  userGoal: string;
  runningTrials: number;
  failedTrials: number;
  succeededTrials: number;
  parameters: ChipDescriptor[] = [];

  @Input()
  experimentName: string;

  @Input()
  experiment: ExperimentK8s;

  constructor() {}

  ngOnChanges(): void {
    if (this.experiment) {
      this.generateExperimentPropsList(this.experiment);
    }
  }

  private generateExperimentPropsList(experiment: ExperimentK8s): void {
    const optimalTrialExists: boolean =
      experiment.status.currentOptimalTrial &&
      !!experiment.status.currentOptimalTrial.bestTrialName;

    const { status, statusIcon } = this.generateExperimentStatus(experiment);
    this.status = status;
    this.statusIcon = statusIcon;

    this.bestTrialName = optimalTrialExists
      ? experiment.status.currentOptimalTrial.bestTrialName
      : 'No optimal trial yet';

    this.generateExperimentBestParameters(experiment);

    this.bestTrialPerformance = this.generateExperimentBestMetrics(
      experiment,
      optimalTrialExists,
    );

    this.userGoal = `${experiment.spec.objective.objectiveMetricName} ${
      experiment.spec.objective.type === ObjectiveTypeEnum.maximize ? '>' : '<'
    } ${experiment.spec.objective.goal}`;

    this.runningTrials = experiment.status.runningTrialList
      ? experiment.status.runningTrialList.length
      : 0;

    this.failedTrials = experiment.status.failedTrialList
      ? experiment.status.failedTrialList.length
      : 0;

    this.succeededTrials = experiment.status.succeededTrialList
      ? experiment.status.succeededTrialList.length
      : 0;
  }

  private generateExperimentStatus(experiment: ExperimentK8s): {
    status: string;
    statusIcon: string;
  } {
    const succeededCondition = getCondition(experiment, StatusEnum.SUCCEEDED);

    if (succeededCondition && succeededCondition.status === 'True') {
      return { status: succeededCondition.message, statusIcon: 'check_circle' };
    }

    const failedCondition = getCondition(experiment, StatusEnum.FAILED);

    if (failedCondition && failedCondition.status === 'True') {
      return { status: failedCondition.message, statusIcon: 'warning' };
    }

    const runningCondition = getCondition(experiment, StatusEnum.RUNNING);

    if (runningCondition && runningCondition.status === 'True') {
      return { status: runningCondition.message, statusIcon: 'schedule' };
    }

    const restartingCondition = getCondition(experiment, StatusEnum.RESTARTING);

    if (restartingCondition && restartingCondition.status === 'True') {
      return { status: restartingCondition.message, statusIcon: 'loop' };
    }

    const createdCondition = getCondition(experiment, StatusEnum.CREATED);

    if (createdCondition && createdCondition.status === 'True') {
      return {
        status: createdCondition.message,
        statusIcon: 'add_circle_outline',
      };
    }
  }

  private generateExperimentBestParameters(
    experiment: ExperimentK8s,
  ): ChipDescriptor[] {
    this.parameters = [];

    if (!experiment.status.currentOptimalTrial.parameterAssignments) {
      return;
    }

    const parameters =
      experiment.status.currentOptimalTrial.parameterAssignments.map(
        param =>
          `${param.name}: ${
            !isNaN(+param.value)
              ? numberToExponential(+param.value, 6)
              : param.value
          }`,
      );

    for (const c of parameters) {
      const chip: ChipDescriptor = { value: c, color: 'primary' };
      this.parameters.push(chip);
    }
  }

  private generateExperimentBestMetrics(
    experiment: ExperimentK8s,
    optimalTrialExists,
  ): ChipDescriptor[] {
    if (!optimalTrialExists) {
      return [];
    }

    const metrics =
      experiment.status.currentOptimalTrial.observation.metrics.map(
        metric =>
          `${metric.name}:  ${
            !isNaN(+metric.latest)
              ? numberToExponential(+metric.latest, 6)
              : metric.latest
          }`,
      );

    return metrics.map(m => {
      return { value: m, color: 'primary', tooltip: 'Latest value' };
    });
  }
}
