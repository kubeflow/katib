import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { CommonModule } from '@angular/common';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { NgxEchartsModule } from 'ngx-echarts';
import { TrialsGraphEchartsComponent } from './trials-graph-echarts.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { SimpleChange } from '@angular/core';

describe('TrialsGraphEchartsComponent', () => {
  let component: TrialsGraphEchartsComponent;
  let fixture: ComponentFixture<TrialsGraphEchartsComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          CommonModule,
          BrowserAnimationsModule,
          MatProgressSpinnerModule,
          NgxEchartsModule.forRoot({
            echarts: () => import('echarts'),
          }),
        ],
        declarations: [TrialsGraphEchartsComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(TrialsGraphEchartsComponent);
    component = fixture.componentInstance;
    component.experimentTrialsCsv =
      'trialName,Status,Validation-accuracy,Train-accuracy,lr,num-layers,optimizer\nbayesian-optimization-52d66cfb,Succeeded,0.113854,0.112390,0.010342079129038503,5,ftrl\nbayesian-optimization-5837197b,Succeeded,0.979100,0.994220,0.024782287934498225,2,sgd';
    component.experiment = {
      spec: {
        objective: {
          type: 'maximize',
          goal: 0.99,
          objectiveMetricName: 'Validation-accuracy',
          additionalMetricNames: ['Train-accuracy'],
          metricStrategies: [
            { name: 'Validation-accuracy', value: 'max' },
            { name: 'Train-accuracy', value: 'max' },
          ],
        },
        parameters: [],
      },
    };
    component.ngOnChanges({
      experimentTrialsCsv: new SimpleChange(
        null,
        component.experimentTrialsCsv,
        false,
      ),
    });

    fixture.detectChanges();
  });

  it('should create dataToDisplay', () => {
    expect(component.dataToDisplay).toEqual([
      ['0.010342079129038503', '5', 'ftrl', '0.112390', '0.113854', '0.113854'],
      ['0.024782287934498225', '2', 'sgd', '0.994220', '0.979100', '0.979100'],
    ]);
  });

  it('should create parallelAxis', () => {
    expect(component.parallelAxis).toEqual([
      {
        dim: 0,
        name: 'lr',
        max: 0.027260516727948048,
        axisLabel: { showMaxLabel: false },
      },
      {
        dim: 1,
        name: 'num layers',
        max: 5.5,
        axisLabel: { showMaxLabel: false },
      },
      {
        dim: 2,
        name: 'optimizer',
        max: NaN,
        axisLabel: { showMaxLabel: false },
      },
      {
        dim: 3,
        name: 'train accuracy',
        max: 1.093642,
        axisLabel: { showMaxLabel: false },
      },
      {
        dim: 4,
        name: 'validation accuracy',
        max: 1.07701,
        axisLabel: { showMaxLabel: false },
      },
      {
        dim: 5,
        max: 1.07701,
        axisTick: { show: false },
        axisLine: { show: false },
        axisLabel: { align: 'right', margin: 1000000 },
      },
    ]);
  });

  it('should create color', () => {
    expect(component.color).toEqual(['#1a2a6c', '#b21f1f', '#fdbb2d']);
  });

  it('should create dataAllInfo', () => {
    expect(component.dataAllInfo).toEqual([
      [
        'bayesian-optimization-52d66cfb',
        'Succeeded',
        '0.010342079129038503',
        '5',
        'ftrl',
        '0.112390',
        '0.113854',
        '0.113854',
      ],
      [
        'bayesian-optimization-5837197b',
        'Succeeded',
        '0.024782287934498225',
        '2',
        'sgd',
        '0.994220',
        '0.979100',
        '0.979100',
      ],
    ]);
  });

  it('should create tooltipDataToDisplay', () => {
    expect(component.tooltipDataToDisplay).toEqual([
      [
        'bayesian-optimization-52d66cfb',
        'Succeeded',
        '0.113854',
        '0.112390',
        '0.010342079129038503',
        '5',
        'ftrl',
      ],
      [
        'bayesian-optimization-5837197b',
        'Succeeded',
        '0.979100',
        '0.994220',
        '0.024782287934498225',
        '2',
        'sgd',
      ],
    ]);
  });

  it('should create tooltipHeaders', () => {
    expect(component.tooltipHeaders).toEqual([
      'trialName',
      'Status',
      'Validation-accuracy',
      'Train-accuracy',
      'lr',
      'num-layers',
      'optimizer',
    ]);
  });
});
