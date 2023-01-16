import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import lowerCase from 'lodash-es/lowerCase';
import capitalize from 'lodash-es/capitalize';
import { ExperimentK8s } from 'src/app/models/experiment.k8s.model';

@Component({
  selector: 'app-trials-graph-echarts',
  templateUrl: './trials-graph-echarts.component.html',
  styleUrls: ['./trials-graph-echarts.component.scss'],
})
export class TrialsGraphEchartsComponent implements OnChanges {
  initOpts = {
    renderer: 'svg',
  };

  options: any;
  dataArray = [];
  dataToDisplay = [];
  dataAllInfo = [];
  tooltipHeaders = [];
  tooltipDataToDisplay = [];
  parallelAxis = [];
  maxAxisValue = [];
  color = [];
  numberOfmetricStrategies: number;

  @Input()
  experimentTrialsCsv: string;

  @Input()
  experiment: ExperimentK8s;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    // Re-render the graph only when we detect changes to the Trials data, received from the backend
    if (!changes.experimentTrialsCsv || !this.experimentTrialsCsv) {
      return;
    }

    this.initializeData();

    this.numberOfmetricStrategies =
      this.experiment.spec.objective.metricStrategies.length; // the number of the output metrics
    let lines = this.experimentTrialsCsv.split('\n');
    let axes = lines[0].split(',');
    let excludeFromGraph = ['trialName', 'Status', 'KFP Run'];
    let excludeFromTooltipHeaders = ['KFP Run'];
    let axesToDisplay = axes.filter(axis => !excludeFromGraph.includes(axis));

    // In case of having additional metrics, move these at the end
    if (this.numberOfmetricStrategies > 1) {
      for (let i = 1; i < this.numberOfmetricStrategies; i++) {
        axesToDisplay.push(axesToDisplay.splice(1, 1)[0]);
      }
    }
    // Move the target metric at the end and duplicate it
    axesToDisplay.push(axesToDisplay.shift());
    axesToDisplay.push(axesToDisplay[axesToDisplay.length - 1]);

    // Set tooltip headers that includes both trial name and status values
    this.tooltipHeaders = axes.filter(
      axis => !excludeFromTooltipHeaders.includes(axis),
    );

    this.dataArray = this.convertCsvToArray(lines, axes);
    this.color = this.createColorHeatmap(this.tooltipHeaders);
    this.parallelAxis = this.createParallelAxis(axesToDisplay, this.dataArray);
    this.prepareGraphData();

    this.options = this.createGraphOptions(
      this.dataToDisplay,
      this.parallelAxis,
      this.color,
      this.dataAllInfo,
      this.tooltipDataToDisplay,
      this.tooltipHeaders,
    );
  }

  // Reset the lists
  initializeData() {
    this.dataArray = [];
    this.dataToDisplay = [];
    this.dataAllInfo = [];
    this.tooltipHeaders = [];
    this.tooltipDataToDisplay = [];
    this.parallelAxis = [];
    this.maxAxisValue = [];
    this.color = [];
  }

  convertCsvToArray(lines, axes) {
    let array = [];
    for (let i = 1; i < lines.length; i++) {
      let obj = {};
      let currentline = lines[i].split(',');
      for (let j = 0; j < axes.length; j++) {
        obj[axes[j]] = currentline[j];
      }
      array.push(obj);
    }
    return array;
  }

  // Set the heatmap color based on the output metric
  createColorHeatmap(tooltipHeaders) {
    let heatmapColor = ['#1a2a6c', '#b21f1f', '#fdbb2d'];
    if (tooltipHeaders[2].includes('loss')) {
      heatmapColor.reverse();
    }
    return heatmapColor;
  }

  createParallelAxis(axesToDisplay, data) {
    // Set the maximum value of each axis
    for (let i = 0; i < axesToDisplay.length; i++) {
      const max =
        Math.max(...data.map(item => item[axesToDisplay[i]])) +
        0.1 * Math.max(...data.map(item => item[axesToDisplay[i]]));
      this.maxAxisValue.push(max);
    }

    // Set the parallel axes of the graph
    let parallelAxisArray = [];
    let parallelAxisObj = {};
    for (let i = 0; i < this.experiment.spec.parameters.length; i++) {
      // In case of having a metric of type categorical, we have to be explicit and set the type of
      // this parallel axis to category since it appears in its own unique way
      if (this.experiment.spec.parameters[i].parameterType === 'categorical') {
        parallelAxisObj = {
          dim: i,
          name: lowerCase(axesToDisplay[i]),
          type: 'category',
        };
      } else {
        parallelAxisObj = {
          dim: i,
          name: lowerCase(axesToDisplay[i]),
          max: this.maxAxisValue[i],
          axisLabel: {
            showMaxLabel: false,
          },
        };
      }
      parallelAxisArray.push(parallelAxisObj);
    }

    for (
      let j = this.experiment.spec.parameters.length;
      j < axesToDisplay.length;
      j++
    ) {
      parallelAxisObj = {
        dim: j,
        name: lowerCase(axesToDisplay[j]),
        max: this.maxAxisValue[j],
        axisLabel: {
          showMaxLabel: false,
        },
      };
      if (j === axesToDisplay.length - 1) {
        parallelAxisObj = {
          dim: j,
          max: this.maxAxisValue[j],
          axisTick: {
            show: false,
          },
          axisLine: {
            show: false,
          },
          axisLabel: {
            // show: false // doesn't work
            align: 'right',
            margin: 1000000,
          },
        };
      }
      parallelAxisArray.push(parallelAxisObj);
    }
    return parallelAxisArray;
  }

  // Set data needed for creating the graph
  prepareGraphData() {
    let trialToDisplay = [];
    let dataAll = [];
    let trialNameStatus = [];
    let trialMetrics = [];
    let tooltipData = [];
    this.dataArray.forEach(trial => {
      delete trial['KFP Run'];

      Object.keys(trial).forEach(axis => {
        if (axis === 'trialName' || axis === 'Status') {
          trialNameStatus.push(trial[axis]);
        } else {
          trialToDisplay.push(trial[axis]);
          trialMetrics.push(trial[axis]);
        }
      });

      // In case of having additional metrics, move these at the end
      if (this.numberOfmetricStrategies > 1) {
        for (let i = 1; i < this.numberOfmetricStrategies; i++) {
          trialToDisplay.push(trialToDisplay.splice(1, 1)[0]);
        }
      }
      // Move the target metric at the end and duplicate it
      trialToDisplay.push(trialToDisplay.shift());
      trialToDisplay.push(trialToDisplay[trialToDisplay.length - 1]);
      if (trialToDisplay[this.parallelAxis.length - 1] !== '') {
        dataAll = trialNameStatus.concat(trialToDisplay);
        tooltipData = trialNameStatus.concat(trialMetrics);
        this.dataAllInfo.push(dataAll);
        this.tooltipDataToDisplay.push(tooltipData);
        this.dataToDisplay.push(trialToDisplay);
      }
      trialNameStatus = [];
      trialToDisplay = [];
      trialMetrics = [];
    });
  }

  createGraphOptions(
    dataToDisplay,
    parallelAxis,
    color,
    dataAllInfo,
    tooltipDataToDisplay,
    tooltipHeaders,
  ) {
    // Set the options value that echarts need to create the graph
    let graphOptions = {
      tooltip: {
        // O(n^2)
        formatter: function (params) {
          return createTooltipText(
            params,
            dataAllInfo,
            tooltipDataToDisplay,
            tooltipHeaders,
          );
        },
        padding: 10,
        borderWidth: 1,
      },
      toolbox: {
        show: true,
        feature: {
          restore: {},
          saveAsImage: {},
        },
      },
      parallelAxis: parallelAxis,
      visualMap: {
        min: 0,
        max:
          this.maxAxisValue[parallelAxis.length - 1] === 0
            ? 1
            : this.maxAxisValue[parallelAxis.length - 1],
        precision: 2,
        dimension: parallelAxis.length - 1,
        inRange: {
          color: color,
        },
        itemHeight: 482,
        itemWidth: 40,
        right: 40,
        bottom: 50,
        align: 'left',
      },
      series: {
        type: 'parallel',
        lineStyle: {
          width: 2,
          opacity: 0.5,
        },
        smooth: true,
        emphasis: {
          focus: 'self',
          lineStyle: {
            width: 3,
            opacity: 1,
          },
        },
        data: dataToDisplay,
      },
    };
    return graphOptions;
  }
}

function createTooltipText(
  params,
  dataAllInfo,
  tooltipDataToDisplay,
  tooltipHeaders,
): string {
  for (let i = 0; i < dataAllInfo.length; i++) {
    const included = dataAllInfo[i].filter(value =>
      params.data.includes(value),
    );
    if (included.length === params.data.length) {
      params.data = tooltipDataToDisplay[i];
      let tooltip = '';
      for (let i = 0; i < tooltipHeaders.length; i++) {
        tooltip +=
          '<b>' +
          capitalize(lowerCase(tooltipHeaders[i])) +
          ': ' +
          '</b>' +
          params.data[i] +
          '<br/>';
      }
      return tooltip;
    }
  }
}
