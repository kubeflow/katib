import {
  ChangeDetectionStrategy,
  Component,
  Input,
  OnChanges,
  OnInit,
} from '@angular/core';
import { ChipDescriptor, getCondition } from 'kubeflow';
import { StatusEnum } from 'src/app/enumerations/status.enum';
import { TrialK8s } from 'src/app/models/trial.k8s.model';
import { numberToExponential } from 'src/app/shared/utils';

@Component({
  selector: 'app-trial-overview',
  templateUrl: './trial-overview.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TrialOverviewComponent implements OnInit, OnChanges {
  status: string;
  statusIcon: string;
  completionTime: string;
  performance: ChipDescriptor[];

  @Input()
  trialName: string;

  @Input()
  trial: TrialK8s;

  @Input()
  experimentName: string;

  constructor() {}

  ngOnInit() {
    if (this.trial) {
      const { status, statusIcon } = this.generateTrialStatus(this.trial);
      this.status = status;
      this.statusIcon = statusIcon;
    }
  }

  ngOnChanges(): void {
    if (this.trial) {
      this.generateTrialPropsList(this.trial);
    }
  }

  private generateTrialPropsList(trial: TrialK8s): void {
    this.performance = this.generateTrialMetrics(this.trial);

    const { status, statusIcon } = this.generateTrialStatus(trial);
    this.status = status;
    this.statusIcon = statusIcon;
    this.statusIcon = statusIcon;
    this.completionTime = trial.status?.completionTime;
  }

  private generateTrialStatus(trial: TrialK8s): {
    status: string;
    statusIcon: string;
  } {
    const succeededCondition = getCondition(trial, StatusEnum.SUCCEEDED);

    if (succeededCondition && succeededCondition.status === 'True') {
      return { status: succeededCondition.message, statusIcon: 'check_circle' };
    }

    const failedCondition = getCondition(trial, StatusEnum.FAILED);

    if (failedCondition && failedCondition.status === 'True') {
      return { status: failedCondition.message, statusIcon: 'warning' };
    }

    const runningCondition = getCondition(trial, StatusEnum.RUNNING);

    if (runningCondition && runningCondition.status === 'True') {
      return { status: runningCondition.message, statusIcon: 'schedule' };
    }

    const restartingCondition = getCondition(trial, StatusEnum.RESTARTING);

    if (restartingCondition && restartingCondition.status === 'True') {
      return { status: restartingCondition.message, statusIcon: 'loop' };
    }

    const createdCondition = getCondition(trial, StatusEnum.CREATED);

    if (createdCondition && createdCondition.status === 'True') {
      return {
        status: createdCondition.message,
        statusIcon: 'add_circle_outline',
      };
    }
  }

  private generateTrialMetrics(trial: TrialK8s): ChipDescriptor[] {
    if (!trial.status.observation || !trial.status.observation.metrics) {
      return [];
    }

    const metrics = trial.status.observation.metrics.map(
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
