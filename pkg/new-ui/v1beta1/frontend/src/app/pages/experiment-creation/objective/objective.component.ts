import { Component, Input, OnInit, OnDestroy } from '@angular/core';
import { FormArray, FormGroup, FormControl } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { ObjectiveTypes } from 'src/app/constants/objective-types.const';
import { Subscription } from 'rxjs';
import { MetricStrategy } from 'src/app/models/experiment.k8s.model';

@Component({
  selector: 'app-form-objective',
  templateUrl: './objective.component.html',
  styleUrls: ['./objective.component.scss'],
})
export class FormObjectiveComponent implements OnInit, OnDestroy {
  @Input()
  objectiveForm: FormGroup;
  strategiesArray: FormArray;
  objectiveTypes = ObjectiveTypes;
  subs = new Subscription();

  get objectiveStrategy() {
    return this.objectiveForm.get('type').value.substring(0, 3);
  }

  constructor(public dialog: MatDialog) {}

  ngOnInit(): void {
    this.strategiesArray = this.objectiveForm.get(
      'metricStrategies',
    ) as FormArray;

    this.subs.add(
      this.objectiveForm.get('type').valueChanges.subscribe(type => {
        this.objectiveForm
          .get('metricStrategy')
          .setValue(this.objectiveStrategy);

        for (const strategy of this.strategiesArray.controls) {
          strategy.get('strategy').setValue(this.objectiveStrategy);
        }
      }),
    );

    this.subs.add(
      this.objectiveForm
        .get('additionalMetricNames')
        .valueChanges.subscribe(metrics => this.setStrategies(metrics)),
    );

    this.setStrategies(this.objectiveForm.get('additionalMetricNames').value);
  }

  ngOnDestroy(): void {
    this.subs.unsubscribe();
  }

  setStrategies(metrics: MetricStrategy[]) {
    this.strategiesArray.clear();

    for (const additionalMetric of metrics) {
      if (!additionalMetric) {
        continue;
      }

      this.strategiesArray.push(
        new FormGroup(
          {
            metric: new FormControl(additionalMetric),
            strategy: new FormControl(this.objectiveStrategy),
          },
          [],
        ),
      );
    }
  }
}
