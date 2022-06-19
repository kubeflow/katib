import {
  ChangeDetectionStrategy,
  Component,
  Input,
  OnChanges,
  OnInit,
  AfterViewInit,
  SimpleChanges,
  ElementRef,
  ViewChild,
} from '@angular/core';
import lowerCase from 'lodash-es/lowerCase';
import { safeDivision } from 'src/app/shared/utils';
import { ExperimentK8s } from '../../../models/experiment.k8s.model';

import { Inject } from '@angular/core';
import { DOCUMENT } from '@angular/common';

import { StatusEnum } from 'src/app/enumerations/status.enum';

declare let d3: any;

type PbtPoint = {
  trialName: string;
  parentUid: string;
  parameters: Object; // all y-axis possible values (parameters + alternativeMetrics)
  generation: number; // generation
  metricValue: number; // evaluation metric
};

@Component({
  selector: 'app-experiment-pbt-tab',
  templateUrl: './pbt-tab.component.html',
  styleUrls: ['./pbt-tab.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class PbtTabComponent implements OnChanges, OnInit, AfterViewInit {
  @ViewChild('pbtGraph') graphWrapper: ElementRef; // used to manipulate svg dom

  graph: any; // svg dom
  selectableNames: string[]; // list of parameters/metrics for UI dropdown
  selectedName: string; // user selected parameter/metric
  displayTrace: boolean;

  private labelData: { [trialName: string]: PbtPoint } = {};
  private trialData: { [trialName: string]: Object } = {};
  private parameterNames: string[]; // parameter names
  private goalName: string = '';
  private data: PbtPoint[][] = []; // data sorted by generation and segment
  private graphHelper: any; // graph metadata and auxiliary info

  @Input()
  experiment: ExperimentK8s;

  @Input()
  labelCsv: string[] = [];

  @Input()
  experimentTrialsCsv: string[] = [];

  constructor(@Inject(DOCUMENT) private document: Document) {
    this.graphHelper = {};
    this.graphHelper.margin = { top: 10, right: 30, bottom: 30, left: 60 };
    this.graphHelper.width =
      460 - this.graphHelper.margin.left - this.graphHelper.margin.right;
    this.graphHelper.height =
      400 - this.graphHelper.margin.top - this.graphHelper.margin.bottom;
    // Track angular initialization since using @Input from template
    this.graphHelper.ngInit = false;
  }

  ngOnInit(): void {
    if (this.experiment.spec.algorithm.algorithmName != 'pbt') {
      // Prevent initialization if not Pbt
      return;
    }

    // Create full list of parameters
    this.parameterNames = this.experiment.spec.parameters.map(param => {
      return param.name;
    });
    // Identify goal
    this.goalName = this.experiment.spec.objective.objectiveMetricName;
    // Create full list of selectable names
    this.selectableNames = [...this.parameterNames];
    if (
      this.experiment.spec.objective.additionalMetricNames &&
      this.experiment.spec.objective.additionalMetricNames.length > 0
    ) {
      this.selectableNames = [
        ...this.selectableNames,
        ...this.experiment.spec.objective.additionalMetricNames,
      ];
    }
    // Create converters for all possible y-axes
    this.graphHelper.yMeta = {};
    for (const param of this.experiment.spec.parameters) {
      this.graphHelper.yMeta[param.name] = {};
      if (
        param.parameterType == 'discrete' ||
        param.parameterType == 'categorical'
      ) {
        this.graphHelper.yMeta[param.name].transform = x => {
          return x;
        };
        this.graphHelper.yMeta[param.name].isNumber = false;
      } else if (param.parameterType == 'double') {
        this.graphHelper.yMeta[param.name].transform = x => {
          return parseFloat(x);
        };
        this.graphHelper.yMeta[param.name].isNumber = true;
      } else {
        this.graphHelper.yMeta[param.name].transform = x => {
          return parseInt(x);
        };
        this.graphHelper.yMeta[param.name].isNumber = true;
      }
    }
    if (this.experiment.spec.objective.additionalMetricNames) {
      for (const metricName of this.experiment.spec.objective
        .additionalMetricNames) {
        if (this.graphHelper.yMeta.hasOwnProperty(metricName)) {
          console.warn(
            'Additional metric name conflict with parameter name; ignoring metric:',
            metricName,
          );
          continue;
        }
        this.graphHelper.yMeta[metricName] = {};
        this.graphHelper.yMeta[metricName].transform = x => {
          return parseFloat(x);
        };
        this.graphHelper.yMeta[metricName].isNumber = true;
      }
    }

    this.graphHelper.ngInit = true;
  }

  ngAfterViewInit(): void {
    if (this.experiment.spec.algorithm.algorithmName != 'pbt') {
      // Remove pbt tab and tab content
      const tabs = document.querySelectorAll('.mat-tab-labels .mat-tab-label');
      for (let i = 0; i < tabs.length; i++) {
        if (
          tabs[i]
            .querySelector('.mat-tab-label-content')
            .innerHTML.includes('PBT')
        ) {
          const tabId = tabs[i].getAttribute('id');
          const tabBodyId = tabId.replace('label', 'content');
          const tabBody = document.querySelector('#' + tabBodyId);
          tabBody.remove();
          tabs[i].remove();
          break;
        }
      }
      return;
    }
    // Specify default choice for dropdown menu
    this.selectedName = this.selectableNames[0];
    // Specify default trace view
    this.displayTrace = false;
  }

  onDropdownChange() {
    // Trigger graph redraw on dropdown change event
    this.clearGraph();
    this.updateGraph();
  }

  onTraceChange() {
    // Trigger graph redraw on trace change event
    // TODO: could use d3.select(..).remove() instead of recreating
    this.clearGraph();
    this.updateGraph();
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (!this.graphHelper || !this.graphHelper.ngInit) {
      console.warn(
        'graphHelper not initialized yet, attempting manual call to ngOnInit()',
      );
      this.ngOnInit();
    }
    // Recompute formatted plotting points on data input changes
    let updatePoints = false;
    if (changes.experimentTrialsCsv && this.experimentTrialsCsv) {
      let trialArr = d3.csv.parse(this.experimentTrialsCsv);
      for (let trial of trialArr) {
        if (
          trial['Status'] == StatusEnum.SUCCEEDED &&
          !this.trialData.hasOwnProperty(trial['trialName'])
        ) {
          this.trialData[trial['trialName']] = trial;
          updatePoints = true;
        }
      }
    }

    if (changes.labelCsv && this.labelCsv) {
      let labelArr = d3.csv.parse(this.labelCsv);
      for (let label of labelArr) {
        if (this.labelData.hasOwnProperty(label['trialName'])) {
          continue;
        }
        let newEntry: PbtPoint = {
          trialName: label['trialName'],
          parentUid: label['pbt.suggestion.katib.kubeflow.org/parent'],
          generation: parseInt(
            label['pbt.suggestion.katib.kubeflow.org/generation'],
          ),
          parameters: undefined,
          metricValue: undefined,
        };
        this.labelData[newEntry.trialName] = newEntry;
        updatePoints = true;
      }
    }

    if (updatePoints) {
      // Lazy; reprocess all points
      let points: PbtPoint[] = [];
      Object.values(this.labelData).forEach(entry => {
        let point = {} as PbtPoint;
        point.trialName = entry.trialName;
        point.generation = entry.generation;
        point.parentUid = entry.parentUid;

        // Find corresponding trial data
        let trial = this.trialData[point.trialName];
        if (trial !== undefined) {
          point.metricValue = parseFloat(trial[this.goalName]);
          point.parameters = {};
          for (let p of this.selectableNames) {
            point.parameters[p] = this.graphHelper.yMeta[p].transform(trial[p]);
          }
        } else {
          point.metricValue = undefined;
        }
        points.push(point);
      });

      // Generate segments
      // Group seeds
      let remaining_points = {};
      for (let p of points) {
        remaining_points[p.trialName] = p;
      }

      this.data = [];
      while (Object.keys(remaining_points).length > 0) {
        let seeds = this.maxGenerationTrials(remaining_points);
        for (let seed of seeds) {
          let segment = [];
          let v = seed;
          while (v) {
            segment.push(v);
            let delete_entry = v.trialName;
            v = remaining_points[v.parentUid];
            delete remaining_points[delete_entry];
          }
          this.data.push(segment);
        }
      }

      this.clearGraph();
      this.updateGraph();
    }
  }

  private maxGenerationTrials(d) {
    let end_seeds = [];
    let max_generation = 0;
    for (let k of Object.keys(d)) {
      if (d[k].generation > max_generation) {
        max_generation = d[k].generation;
        end_seeds = [];
        end_seeds.push(d[k]);
      } else if (d[k].generation === max_generation) {
        end_seeds.push(d[k]);
      }
    }
    return end_seeds;
  }

  private clearGraph() {
    // Clear any existing views from graphs object
    if (this.graph) {
      // this.graph.remove(); // d3 remove from DOM
      d3.select(this.graphWrapper.nativeElement).select('svg').remove();
    }
    this.graph = undefined;
  }

  private createGraph() {
    this.graph = d3
      .select(this.graphWrapper.nativeElement)
      .append('svg')
      .attr(
        'width',
        this.graphHelper.width +
          this.graphHelper.margin.left +
          this.graphHelper.margin.right,
      )
      .attr(
        'height',
        this.graphHelper.height +
          this.graphHelper.margin.top +
          this.graphHelper.margin.bottom,
      )
      .append('g')
      .attr(
        'transform',
        'translate(' +
          this.graphHelper.margin.left +
          ',' +
          this.graphHelper.margin.top +
          ')',
      );
  }

  private getRangeX() {
    const xValues = Object.values(this.labelData).map(entry => {
      return entry.generation;
    });
    return d3.scale
      .linear()
      .domain(d3.extent(xValues))
      .range([0, this.graphHelper.width]);
  }

  private getRangeY(key) {
    if (this.selectableNames.includes(key)) {
      const paramValues = Object.keys(this.trialData).map(trialName => {
        return this.graphHelper.yMeta[key].transform(
          this.trialData[trialName][key],
        );
      });

      if (this.graphHelper.yMeta[key].isNumber) {
        return d3.scale
          .linear()
          .domain(d3.extent(paramValues))
          .range([this.graphHelper.height, 0]);
      } else {
        paramValues.sort((a, b) => a - b);
        return d3.scale
          .ordinal()
          .domain(paramValues)
          .range([this.graphHelper.height, 0]);
      }
    } else {
      console.error('Key(' + key + ') not found in y-axis list');
    }
  }

  private getColorScaleZ() {
    let values = [];
    for (const segment of this.data) {
      for (const point of segment) {
        values.push(point.metricValue);
      }
    }
    return d3.scale
      .linear()
      .domain(d3.extent(values))
      .interpolate(d3.interpolateHcl)
      .range([d3.rgb('#cfd8dc'), d3.rgb('#263238')]);
  }

  private updateGraph() {
    if (!this.graphHelper || !this.graphHelper.ngInit) {
      // ngOnInit not called yet
      return;
    }
    if (!this.data.length || this.data.length == 0) {
      // Data not initialized
      return;
    }
    if (!this.graph) {
      this.createGraph();
    }

    // Add X axis
    let xAxis = this.getRangeX();
    this.graph
      .append('g')
      .attr('transform', 'translate(0,' + this.graphHelper.height + ')')
      .call(d3.svg.axis().scale(xAxis).orient('bottom'));
    this.graph
      .append('text')
      .attr('text-anchor', 'middle')
      .attr('x', this.graphHelper.width / 2)
      .attr('y', this.graphHelper.height + 30)
      .text('Generation');
    // Add Y axis
    let yAxis = this.getRangeY(this.selectedName);
    this.graph.append('g').call(d3.svg.axis().scale(yAxis).orient('left'));
    this.graph
      .append('text')
      .attr('text-anchor', 'middle')
      .attr('x', -this.graphHelper.height / 2)
      .attr('y', -50)
      .attr('transform', 'rotate(-90)')
      .text(this.selectedName);
    // Change line width
    this.graph
      .selectAll('path')
      .style({ stroke: 'black', fill: 'none', 'stroke-width': '1px' });

    // Add the points
    const sparam = this.selectedName;
    const colorScale = this.getColorScaleZ();
    for (const segment of this.data) {
      // Plot only valid points
      const validSegment = segment.filter(
        point =>
          point.parameters &&
          point.parameters[sparam] &&
          point.metricValue !== undefined,
      );
      this.graph
        .append('g')
        .selectAll('dot')
        .data(validSegment)
        .enter()
        .append('circle')
        .attr('cx', function (d) {
          return xAxis(d.generation);
        })
        .attr('cy', function (d) {
          return yAxis(d.parameters[sparam]);
        })
        .attr('r', 2)
        .attr('fill', function (d) {
          return colorScale(d.metricValue);
        });

      if (this.displayTrace && validSegment.length > 0) {
        // Add the lines
        const strokeColor = colorScale(validSegment[0].metricValue);
        this.graph
          .append('path')
          .datum(validSegment)
          .attr('fill', 'none')
          .attr('stroke', strokeColor)
          .attr('stroke-width', 1)
          .attr(
            'd',
            d3.svg
              .line()
              .x(function (d) {
                return xAxis(d.generation);
              })
              .y(function (d) {
                return yAxis(d.parameters[sparam]);
              }),
          );
      }
    }
  }
}
