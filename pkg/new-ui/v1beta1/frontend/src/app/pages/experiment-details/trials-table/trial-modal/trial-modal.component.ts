import { Component, OnInit } from '@angular/core';
import { curveBasis } from 'd3-shape';
import { KWABackendService } from 'src/app/services/backend.service';
import { transformStringResponses } from 'src/app/shared/utils';

@Component({
  selector: 'app-trial-modal',
  templateUrl: './trial-modal.component.html',
  styleUrls: ['./trial-modal.component.scss'],
})
export class TrialModalComponent implements OnInit {
  trialName: string;
  namespace: string;
  dataLoaded: boolean;

  // chart's options
  view = [700, 600];
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
  chartData: { name: string; series: { name: string; value: number }[] }[] = [];
  curve = curveBasis;
  yScaleMax = 0;
  yScaleMin = 1;

  constructor(private backendService: KWABackendService) {}

  ngOnInit() {
    this.backendService
      .getTrial(this.trialName, this.namespace)
      .subscribe(response => {
        const { types, details } = transformStringResponses(response);
        const nameIndex = types.findIndex(type => type === 'Metric name');
        const timeIndex = types.findIndex(type => type === 'Time');
        const valueIndex = types.findIndex(type => type === 'Value');

        details.forEach(detail => {
          const name = detail[nameIndex];
          const value = +detail[valueIndex];
          const time = new Date(detail[timeIndex]);
          const formattedDate = `${time.getHours()}:${time.getMinutes()}`;

          if (value > this.yScaleMax) {
            this.yScaleMax = value;
          }

          if (value < this.yScaleMin) {
            this.yScaleMin = value;
          }

          if (this.chartData.find(chart => chart.name === name)) {
            const index = this.chartData.findIndex(
              chart => chart.name === name,
            );
            this.chartData[index].series.push({
              name: formattedDate,
              value,
            });
          } else {
            this.chartData.push({
              name,
              series: [{ name: formattedDate, value }],
            });
          }
        });

        this.yScaleMin = Math.floor(this.yScaleMin * 10) / 10;
        this.yScaleMax = Math.ceil(this.yScaleMax * 10) / 10;
        this.dataLoaded = true;
      });
  }
}
