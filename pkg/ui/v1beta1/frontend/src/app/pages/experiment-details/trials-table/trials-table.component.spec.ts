import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { CommonModule } from '@angular/common';
import { MatTableModule } from '@angular/material/table';
import { MatIconModule } from '@angular/material/icon';
import { MatDialogModule } from '@angular/material/dialog';
import { RouterTestingModule } from '@angular/router/testing';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatButtonModule } from '@angular/material/button';
import { TrialsTableComponent } from './trials-table.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { SimpleChange } from '@angular/core';
import {
  PropertyValue,
  StatusValue,
  ComponentValue,
  LinkValue,
  LinkType,
  KubeflowModule,
} from 'kubeflow';
import { parseStatus } from '../../experiments/utils';
import lowerCase from 'lodash-es/lowerCase';
import { KfpRunComponent } from './kfp-run/kfp-run.component';
import { MatIconTestingModule } from '@angular/material/icon/testing';
import { TrialDetailsModule } from './trial-details/trial-details.module';

describe('TrialsTableComponent', () => {
  let component: TrialsTableComponent;
  let fixture: ComponentFixture<TrialsTableComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          CommonModule,
          MatTableModule,
          MatDialogModule,
          MatIconModule,
          MatTooltipModule,
          MatButtonModule,
          KubeflowModule,
          TrialDetailsModule,
          RouterTestingModule,
          BrowserAnimationsModule,
          KubeflowModule,
          MatIconTestingModule,
        ],
        declarations: [TrialsTableComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(TrialsTableComponent);
    component = fixture.componentInstance;
    component.experimentName = ['open-vaccine'];
    component.displayedColumns = [
      'Status',
      'Trial name',
      'Validation loss',
      'Lr',
      'Batch size',
      'Embed dim',
      'Dropout',
      'Sp dropout',
      'Kfp run',
    ];
    component.data = [
      [
        'Succeeded',
        'open-vaccine-0f37u-6cd03cbf',
        '0.70573',
        '0.0001',
        '32',
        '20',
        '0.2',
        '0.2',
        '9af7c534-689a-48aa-996b-537d13989729',
        '0',
      ],
      [
        'Succeeded',
        'open-vaccine-0f37u-8ec17b8f',
        '0.76401',
        '0.0001',
        '96',
        '20',
        '0.2',
        '0.2',
        '19aed8e0-143c-49d8-8bd3-7ebb464181d8',
        '1',
      ],
    ];
    component.ngOnChanges({
      displayedColumns: new SimpleChange(
        null,
        component.displayedColumns,
        false,
      ),
    });

    fixture.detectChanges();
  });

  it('should create processedData', () => {
    expect(component.processedData).toEqual([
      {
        'trial name': 'open-vaccine-0f37u-6cd03cbf',
        status: 'Succeeded',
        'validation loss': '0.70573',
        lr: '0.0001',
        'batch size': '32',
        'embed dim': '20',
        dropout: '0.2',
        'sp dropout': '0.2',
        'kfp run': '9af7c534-689a-48aa-996b-537d13989729',
        link: {
          text: 'open-vaccine-0f37u-6cd03cbf',
          url: '/experiment/open-vaccine/trial/open-vaccine-0f37u-6cd03cbf',
        },
      },
      {
        'trial name': 'open-vaccine-0f37u-8ec17b8f',
        status: 'Succeeded',
        'validation loss': '0.76401',
        lr: '0.0001',
        'batch size': '96',
        'embed dim': '20',
        dropout: '0.2',
        'sp dropout': '0.2',
        'kfp run': '19aed8e0-143c-49d8-8bd3-7ebb464181d8',
        link: {
          text: 'open-vaccine-0f37u-8ec17b8f',
          url: '/experiment/open-vaccine/trial/open-vaccine-0f37u-8ec17b8f',
        },
      },
    ]);
  });

  it('should create config', () => {
    expect(component.config).toEqual({
      columns: [
        {
          matColumnDef: 'Status',
          matHeaderCellDef: 'Status',
          value: new StatusValue({
            valueFn: parseStatus,
          }),
          sort: true,
        },
        {
          matColumnDef: 'name',
          matHeaderCellDef: 'Trial name',
          style: { width: '25%' },
          value: new LinkValue({
            field: 'link',
            popoverField: 'trial name',
            truncate: true,
            linkType: LinkType.Internal,
          }),
          sort: true,
        },
        {
          matColumnDef: 'Validation loss',
          matHeaderCellDef: 'Validation loss',
          value: new PropertyValue({
            field: lowerCase(component.displayedColumns[2]),
          }),
          sort: true,
          sortingPreprocessorFn: jasmine.any(Function),
        },
        {
          matColumnDef: 'Lr',
          matHeaderCellDef: 'Lr',
          value: new PropertyValue({
            field: lowerCase(component.displayedColumns[3]),
          }),
          sort: true,
          sortingPreprocessorFn: jasmine.any(Function),
        },
        {
          matColumnDef: 'Batch size',
          matHeaderCellDef: 'Batch size',
          value: new PropertyValue({
            field: lowerCase(component.displayedColumns[4]),
          }),
          sort: true,
          sortingPreprocessorFn: jasmine.any(Function),
        },
        {
          matColumnDef: 'Embed dim',
          matHeaderCellDef: 'Embed dim',
          value: new PropertyValue({
            field: lowerCase(component.displayedColumns[5]),
          }),
          sort: true,
          sortingPreprocessorFn: jasmine.any(Function),
        },
        {
          matColumnDef: 'Dropout',
          matHeaderCellDef: 'Dropout',
          value: new PropertyValue({
            field: lowerCase(component.displayedColumns[6]),
          }),
          sort: true,
          sortingPreprocessorFn: jasmine.any(Function),
        },
        {
          matColumnDef: 'Sp dropout',
          matHeaderCellDef: 'Sp dropout',
          value: new PropertyValue({
            field: lowerCase(component.displayedColumns[7]),
          }),
          sort: true,
          sortingPreprocessorFn: jasmine.any(Function),
        },
        {
          matHeaderCellDef: '',
          matColumnDef: 'actions',
          value: new ComponentValue({
            component: KfpRunComponent,
          }),
        },
      ],
    });
  });

  it('should sort numeric metric columns by value, not lexicographically', () => {
    const metricColumn = component.config.columns.find(
      col => col.matColumnDef === 'Validation loss',
    );

    // Values formatted in scientific notation (e.g. by numberToExponential)
    // must sort as numbers: 0.0049713 < 0.070573 < 0.70573.
    const sorted = ['0.70573', '4.9713e-3', '7.0573e-2']
      .map(metricColumn.sortingPreprocessorFn)
      .sort((a, b) => a - b);

    expect(sorted).toEqual([0.0049713, 0.070573, 0.70573]);
  });

  it('should keep a missing metric value as an empty string when sorting', () => {
    const metricColumn = component.config.columns.find(
      col => col.matColumnDef === 'Validation loss',
    );

    // An empty string must stay '' rather than becoming 0 (Number('') === 0),
    // otherwise trials with a missing metric would sort as if it were zero.
    expect(metricColumn.sortingPreprocessorFn('')).toEqual('');
  });

  it('should fall back to the raw value when sorting non-numeric properties', () => {
    const metricColumn = component.config.columns.find(
      col => col.matColumnDef === 'Validation loss',
    );

    expect(metricColumn.sortingPreprocessorFn('not-a-number')).toEqual(
      'not-a-number',
    );
  });
});
