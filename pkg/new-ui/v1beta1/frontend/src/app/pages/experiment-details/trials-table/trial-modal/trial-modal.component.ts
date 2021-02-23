import { Component, OnInit } from '@angular/core';
import { curveLinear } from 'd3-shape';
import { KWABackendService } from 'src/app/services/backend.service';
import { transformStringResponses } from 'src/app/shared/utils';

interface ChartPoint {
  name: string;
  series: {
    name: any;
    value: number;
  }[];
}

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

          // figure out the min-max values in y-axis
          if (value > this.yScaleMax) {
            this.yScaleMax = value;
          } else {
            this.yScaleMin = value;
          }

          if (this.chartData.find(chart => chart.name === name)) {
            // chart has already some points, append current one
            const index = this.chartData.findIndex(
              chart => chart.name === name,
            );

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
      });
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
}
