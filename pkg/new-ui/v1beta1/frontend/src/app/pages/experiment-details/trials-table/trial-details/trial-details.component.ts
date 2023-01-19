import { Component, OnDestroy, OnInit } from '@angular/core';
import { curveLinear } from 'd3-shape';
import { KWABackendService } from 'src/app/services/backend.service';
import { transformStringResponses } from 'src/app/shared/utils';
import { ActivatedRoute, Router } from '@angular/router';
import { TrialK8s } from 'src/app/models/trial.k8s.model';
import { Subscription } from 'rxjs';
import { StatusEnum } from 'src/app/enumerations/status.enum';
import { ExponentialBackoff, getCondition, NamespaceService } from 'kubeflow';

interface ChartPoint {
  name: string;
  series: {
    name: any;
    value: number;
  }[];
}

@Component({
  selector: 'app-trial-details',
  templateUrl: './trial-details.component.html',
  styleUrls: ['./trial-details.component.scss'],
})
export class TrialDetailsComponent implements OnInit, OnDestroy {
  trialName: string;
  namespace: string;
  dataLoaded: boolean;
  trialDetails: TrialK8s;
  experimentName: string;
  showTrialGraph: boolean = false;

  // chart's options
  view = [700, 500];
  legend = true;
  legendTitle = '';
  animations = true;
  xAxis = true;
  yAxis = true;
  showYAxisLabel = true;
  showXAxisLabel = true;
  xAxisLabel = 'Datetime';
  yAxisLabel = 'Value';
  timeline = true;
  chartData: ChartPoint[] = [];
  curve = curveLinear;
  yScaleMax = 0;
  yScaleMin = 1;

  constructor(
    private activatedRoute: ActivatedRoute,
    private router: Router,
    private backendService: KWABackendService,
    private namespaceService: NamespaceService,
  ) {}

  private poller: ExponentialBackoff;

  private subs = new Subscription();

  ngOnInit() {
    this.trialName = this.activatedRoute.snapshot.params.trialName;
    this.experimentName = this.activatedRoute.snapshot.params.experimentName;

    this.subs.add(
      this.namespaceService.getSelectedNamespace().subscribe(namespace => {
        this.namespace = namespace;
        this.updateTrialInfo();
      }),
    );
  }

  ngOnDestroy(): void {
    this.subs.unsubscribe();
  }

  processTrialInfo(response: any): void {
    const { types, details } = transformStringResponses(response);
    const nameIndex = types.findIndex(type => type === 'Metric name');
    const timeIndex = types.findIndex(type => type === 'Time');
    const valueIndex = types.findIndex(type => type === 'Value');

    details.forEach(detail => {
      const name = detail[nameIndex];
      const value = +detail[valueIndex];
      const time = new Date(detail[timeIndex]);

      // figure out the min-max values in y-axis
      if (value > this.yScaleMax) {
        this.yScaleMax = value;
      } else {
        this.yScaleMin = value;
      }

      if (this.chartData.find(chart => chart.name === name)) {
        // chart has already some points, append current one
        const index = this.chartData.findIndex(chart => chart.name === name);

        this.chartData[index].series.push({
          //name: formattedDate,
          name: time,
          value,
        });
      } else {
        // first point of the chart
        this.chartData.push({
          name,
          series: [{ name: time, value }],
        });
      }
    });

    this.yScaleMin = Math.floor(this.yScaleMin * 10) / 10;
    this.yScaleMax = Math.ceil(this.yScaleMax * 10) / 10;
    this.dataLoaded = true;
  }

  private updateTrialInfo() {
    this.backendService
      .getTrial(this.trialName, this.namespace)
      .subscribe(response => this.processTrialInfo(response));

    this.backendService
      .getTrialInfo(this.trialName, this.namespace)
      .subscribe((response: TrialK8s) => {
        this.trialDetails = response;
        const status = this.trialStatus(response);

        if (status && status === StatusEnum.SUCCEEDED) {
          this.showTrialGraph = true;
        }

        if (
          status &&
          !(status === StatusEnum.FAILED || status === StatusEnum.SUCCEEDED)
        ) {
          // if the status of the trial is not succeeded either failed
          // then start polling the trial
          this.startTrialInfoPolling();
          this.startTrialPolling();
        }
      });
  }

  private trialStatus(trial: TrialK8s): StatusEnum {
    const succeededCondition = getCondition(trial, StatusEnum.SUCCEEDED);

    if (succeededCondition && succeededCondition.status === 'True') {
      return StatusEnum.SUCCEEDED;
    }

    const failedCondition = getCondition(trial, StatusEnum.FAILED);

    if (failedCondition && failedCondition.status === 'True') {
      return StatusEnum.FAILED;
    }

    const runningCondition = getCondition(trial, StatusEnum.RUNNING);

    if (runningCondition && runningCondition.status === 'True') {
      return StatusEnum.RUNNING;
    }

    const restartingCondition = getCondition(trial, StatusEnum.RESTARTING);

    if (restartingCondition && restartingCondition.status === 'True') {
      return StatusEnum.RESTARTING;
    }

    const createdCondition = getCondition(trial, StatusEnum.CREATED);

    if (createdCondition && createdCondition.status === 'True') {
      return StatusEnum.CREATED;
    }
  }

  private startTrialInfoPolling() {
    this.poller = new ExponentialBackoff({
      interval: 5000,
      retries: 1,
      maxInterval: 5001,
    });

    // Poll for new data and reset the poller if different data is found
    this.subs.add(
      this.poller.start().subscribe(() => {
        this.backendService
          .getTrialInfo(this.trialName, this.namespace)
          .subscribe(response => {
            this.trialDetails = response;
            const status = this.trialStatus(response);

            if (status && status === StatusEnum.SUCCEEDED) {
              this.showTrialGraph = true;
            }
          });
      }),
    );
  }

  private startTrialPolling() {
    this.poller = new ExponentialBackoff({
      interval: 5000,
      retries: 1,
      maxInterval: 5001,
    });

    // Poll for new data and reset the poller if different data is found
    this.subs.add(
      this.poller.start().subscribe(() => {
        this.backendService
          .getTrial(this.trialName, this.namespace)
          .subscribe(response => this.processTrialInfo(response));
      }),
    );
  }

  public xAxisFormat(time: Date) {
    function zeroPad(n: number): string {
      if (n < 10) {
        return `0${n}`;
      }

      return n.toString();
    }

    const hours = zeroPad(time.getHours());
    const minutes = zeroPad(time.getMinutes());
    const seconds = zeroPad(time.getSeconds());
    return `${hours}:${minutes}:${seconds}`;
  }

  returnToExperimentDetails() {
    this.router.navigate(
      [`/experiment/${this.namespace}/${this.experimentName}`],
      {
        queryParams: { tab: 'trials' },
      },
    );
  }
}
