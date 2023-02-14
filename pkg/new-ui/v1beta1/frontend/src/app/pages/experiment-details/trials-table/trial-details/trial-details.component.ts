import { Component, OnDestroy, OnInit } from '@angular/core';
import { KWABackendService } from 'src/app/services/backend.service';
import { transformStringResponses } from 'src/app/shared/utils';
import { ActivatedRoute, Router } from '@angular/router';
import { TrialK8s } from 'src/app/models/trial.k8s.model';
import { Subscription } from 'rxjs';
import { StatusEnum } from 'src/app/enumerations/status.enum';
import { ExponentialBackoff, getCondition, NamespaceService } from 'kubeflow';

interface ChartPoint {
  name: string;
  type: string;
  data: [string, number][];
}

@Component({
  selector: 'app-trial-details',
  templateUrl: './trial-details.component.html',
  styleUrls: ['./trial-details.component.scss'],
})
export class TrialDetailsComponent implements OnInit, OnDestroy {
  initOpts = {
    renderer: 'svg',
  };

  trialName: string;
  namespace: string;
  pageLoading = true;
  trialDetails: TrialK8s;
  experimentName: string;
  showTrialGraph: boolean = false;
  options: {};
  trialLogs: string;
  logsRequestError: string;
  chartData: ChartPoint[] = [];
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
      const time = detail[timeIndex];
      const type = 'line';

      // figure out the min-max values in y-axis
      if (value > this.yScaleMax) {
        this.yScaleMax = value;
      }

      if (value < this.yScaleMin) {
        this.yScaleMin = value;
      }

      if (this.chartData.find(chart => chart.name === name)) {
        // chart has already some points, append current one
        const index = this.chartData.findIndex(chart => chart.name === name);

        this.chartData[index].data.push([time, value]);
      } else {
        // first point of the chart
        this.chartData.push({
          name,
          type,
          data: [[time, value]],
        });
      }
    });

    this.yScaleMin = Math.floor(this.yScaleMin * 10) / 10;
    this.yScaleMax = Math.ceil(this.yScaleMax * 10) / 10;

    this.options = this.createGraphOptions(
      this.chartData,
      this.yScaleMin,
      this.yScaleMax,
    );
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
        this.pageLoading = false;
      });

    this.backendService.getTrialLogs(this.trialName, this.namespace).subscribe(
      logs => {
        this.trialLogs = logs;
        this.logsRequestError = null;
      },
      error => {
        this.trialLogs = null;
        this.logsRequestError = error;
      },
    );
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

  createGraphOptions(
    chartData: ChartPoint[],
    yScaleMin: number,
    yScaleMax: number,
  ) {
    // Set the options value that echarts need to create the graph
    let graphOptions = {
      legend: {
        data: ['Train-accuracy', 'Validation-accuracy'],
      },
      tooltip: {
        trigger: 'axis',
      },
      toolbox: {
        show: true,
        feature: {
          dataZoom: {
            yAxisIndex: 'none',
          },
          dataView: {
            readOnly: true,
            buttonColor: '#1e88e5',
            optionToContent: function (opt) {
              var series = opt.series;
              var table = '';
              var tmp;
              for (const sr of series) {
                tmp = table;
                table =
                  '<table style="width:100%; table-layout: fixed;"><tbody><tr>' +
                  '<td style="font-weight: bold; width:10%">Timestamp</td>' +
                  '<td style="font-weight: bold; width:30%">' +
                  sr.name +
                  '</td>' +
                  '</tr>';

                for (var i = 0; i < sr.data.length; i++) {
                  table +=
                    '<tr>' +
                    '<td>' +
                    sr.data[i][0].replace('T', ' ').substring(0, 19) +
                    '</td>' +
                    '<td>' +
                    sr.data[i][1] +
                    '</td>' +
                    '</tr>';
                }
                table += '</tbody></table><br>';
                table = tmp + table;
              }

              return table;
            },
          },
          saveAsImage: {},
        },
      },
      xAxis: [{ type: 'time' }],
      yAxis: [{ type: 'value', min: yScaleMin, max: yScaleMax }],
      series: chartData,
    };
    return graphOptions;
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
