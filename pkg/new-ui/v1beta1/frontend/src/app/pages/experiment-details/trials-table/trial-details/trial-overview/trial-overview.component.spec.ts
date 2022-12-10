import { CommonModule } from '@angular/common';
import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import {
  ConditionsTableModule,
  DetailsListModule,
  HeadingSubheadingRowModule,
} from 'kubeflow';
import { TrialMetricsModule } from './metrics/metrics.component.module';

import { TrialOverviewComponent } from './trial-overview.component';

describe('TrialOverviewComponent', () => {
  let component: TrialOverviewComponent;
  let fixture: ComponentFixture<TrialOverviewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [TrialOverviewComponent],
      imports: [
        CommonModule,
        ConditionsTableModule,
        DetailsListModule,
        HeadingSubheadingRowModule,
        TrialMetricsModule,
        MatSnackBarModule,
      ],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TrialOverviewComponent);
    component = fixture.componentInstance;
    component.trial = {
      status: {
        startTime: '2022-06-01T09:58:23Z',
        completionTime: '2022-06-01T10:07:45Z',
        conditions: [
          {
            type: 'Created',
            status: 'True',
            reason: 'TrialCreated',
            message: 'Trial is created',
            lastUpdateTime: '2022-06-01T09:58:23Z',
            lastTransitionTime: '2022-06-01T09:58:23Z',
          },
        ],
        observation: {
          metrics: [
            {
              name: 'Validation-accuracy',
              latest: '0.113854',
              min: '0.113854',
              max: '0.113854',
            },
          ],
        },
      },
    };

    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
