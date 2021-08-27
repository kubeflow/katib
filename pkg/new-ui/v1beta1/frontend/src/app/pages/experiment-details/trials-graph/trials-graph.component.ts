import {
  ChangeDetectionStrategy,
  Component,
  Input,
  OnChanges,
  SimpleChanges,
} from '@angular/core';
import lowerCase from 'lodash-es/lowerCase';
import { safeDivision } from 'src/app/shared/utils';

import { parcoords } from './d3.parcoords';
declare let d3: any;

@Component({
  selector: 'app-trials-graph',
  templateUrl: './trials-graph.component.html',
  styleUrls: ['./trials-graph.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TrialsGraphComponent implements OnChanges {
  dataLoaded = false;
  graph: any;

  private data = [];

  @Input()
  experimentTrialsCsv: string[] = [];

  @Input()
  hoveredTrial: number;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.experimentTrialsCsv && this.experimentTrialsCsv) {
      this.data = d3.csv.parse(this.experimentTrialsCsv);

      this.data.forEach(trial => {
        delete trial['KFP Run'];
        Object.keys(trial).forEach(key => !trial[key] && (trial[key] = 0));
      });

      this.createOrUpdateGraph();
    }

    if (!changes.hoveredTrial || !this.graph) {
      return;
    }

    // highlight row in chart
    if (isNaN(this.hoveredTrial) || this.hoveredTrial === null) {
      this.graph.unhighlight();
      return;
    }

    const d = this.graph.brushed() || this.data;
    if (d[this.hoveredTrial]) {
      this.graph.highlight([d[this.hoveredTrial]]);
    }
  }

  private createOrUpdateGraph() {
    if (this.data.length && this.data.length > 0) {
      if (!this.graph) {
        this.createGraph(this.data);
        this.dataLoaded = true;
        return;
      }

      this.graph.dimensions({}).data([]).render();

      const dimensions = this.buildDimensions(this.data);
      this.graph
        .data(this.data)
        .hideAxis(['trialName', 'Status', 'id'])
        .dimensions(dimensions)
        .render()
        .updateAxes();
    }
  }

  private createGraph(data) {
    const firstMetric = Object.keys(data[0])[2];
    data.forEach((trial, index) => (trial.id = index));

    // color scale for zscores
    const zcolorscale = d3.scale
      .linear()
      .domain([-1, -0.6, -0.3, 0, 0.3, 0.6, 1])
      .range([
        '#0e0787',
        '#5e01a6',
        '#9b169e',
        '#c9447a',
        '#ea7456',
        '#fba238',
        '#fccc25',
      ])
      .interpolate(d3.interpolateLab);

    this.graph = parcoords()('#trial-graph').alpha(0.6);

    const dimensions = this.buildDimensions(data);

    this.graph
      .data(data)
      .hideAxis(['trialName', 'Status', 'id'])
      .composite('darken')
      .smoothness(0.1)
      .bundlingStrength(0) // set bundling strength
      .dimensions(dimensions)
      .render()
      .shadows()
      .brushMode('1D-axes') // enable brushing
      .interactive() // command line mode
      .reorderable();

    this.changeColor(this.graph, zcolorscale, firstMetric);

    // click label to activate coloring
    this.graph.svg
      .selectAll('.dimension')
      .on('click', this.changeColor.bind(this, this.graph, zcolorscale))
      .selectAll('.label')
      .style('font-size', '14px');
  }

  // update color
  private changeColor(graph, zcolorscale, dimension) {
    graph.svg
      .selectAll('.dimension')
      .style('font-weight', 'normal')
      .filter(d => d === dimension)
      .style('font-weight', 'bold');

    graph.color(this.zcolor(zcolorscale, graph.data(), dimension)).render();
  }

  // return color function based on plot and dimension
  private mean(values: number[]) {
    let sum = 0;
    values.forEach(v => (sum += v));
    return sum / values.length;
  }

  private stdDeviation(values: number[]) {
    const m = this.mean(values);
    return Math.sqrt(
      values.reduce((sq, n) => {
        return sq + Math.pow(n - m, 2);
      }, 0) /
        (values.length - 1),
    );
  }

  private zcolor(zcolorscale, col: any[], dimension) {
    const scores = col.map(trial => trial[dimension]);
    const z = this.zscore(scores.map(parseFloat));
    return d => zcolorscale(z(d[dimension]));
  }

  // color by zscore
  private zscore(col) {
    const mean = this.mean(col);
    const sigma = this.stdDeviation(col);
    return d => (d - mean) / sigma;
  }

  private buildDimensions(data) {
    const range =
      this.graph.height() -
      this.graph.margin().top -
      this.graph.margin().bottom;

    const dimensions = {};
    Object.keys(data[0])
      .filter(key => key !== 'trialName' && key !== 'Status' && key !== 'id')
      .forEach(dimension => {
        const min = d3.min(data, d => {
          const value = +d[dimension];
          if (isNaN(value)) {
            return;
          }

          return value;
        });
        const max = d3.max(data, d => {
          const value = +d[dimension];
          if (isNaN(value)) {
            return;
          }

          return value;
        });

        if (!min && !max) {
          dimensions[dimension] = {
            title: dimension,
            type: 'string',
          };

          return;
        }

        const log = d3.scale
          .linear()
          .domain([
            min - Math.abs(safeDivision(min, 15)),
            max + Math.abs(safeDivision(max, 15)),
          ])
          .range([range, 1]);

        const minMaxLength = min.toString().length;
        const average = safeDivision(min + max, 2);

        const tickValues = [
          min,
          +safeDivision(min + average, 2)
            .toString()
            .slice(0, minMaxLength),
          +average.toString().slice(0, minMaxLength),
          +safeDivision(max + average, 2)
            .toString()
            .slice(0, minMaxLength),
          max,
        ].filter((value, index, array) => value !== array[index - 1]);

        dimensions[dimension] = {
          title: lowerCase(dimension),
          yscale: log,
          innerTickSize: 5,
          outerTickSize: 8,
          type: 'number',
          tickValues,
        };
      });

    return dimensions;
  }
}
