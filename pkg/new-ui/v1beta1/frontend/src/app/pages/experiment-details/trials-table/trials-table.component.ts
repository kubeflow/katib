import {
  ChangeDetectionStrategy,
  Component,
  EventEmitter,
  Input,
  OnChanges,
  Output,
  SimpleChanges,
} from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import {
  PropertyValue,
  StatusValue,
  ComponentValue,
  TableConfig,
  ActionEvent,
} from 'kubeflow';
import { parseStatus } from '../../experiments/utils';
import lowerCase from 'lodash-es/lowerCase';
import { KfpRunComponent } from './kfp-run/kfp-run.component';
import { Router } from '@angular/router';

@Component({
  selector: 'app-trials-table',
  templateUrl: './trials-table.component.html',
  styleUrls: ['./trials-table.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TrialsTableComponent implements OnChanges {
  @Input()
  displayedColumns = [];

  @Input()
  data = [];

  @Input()
  experimentName = [];

  @Input()
  namespace: string;

  @Input()
  bestTrialName: string;

  @Output()
  mouseOnTrial = new EventEmitter<number>();

  @Output()
  leaveMouseFromTrial = new EventEmitter<void>();

  bestTrialRow: {};

  config: TableConfig = { columns: [] };

  processedData = [];

  constructor(public dialog: MatDialog, private router: Router) {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.displayedColumns && this.displayedColumns.length !== 0) {
      this.displayedColumns = this.displayedColumns.slice(
        0,
        this.displayedColumns.length,
      );
      this.processedData = this.setData(this.data, this.displayedColumns);
      this.config = this.setConfig(this.displayedColumns, this.processedData);
    }

    if (this.data.length > 0 && this.bestTrialName) {
      this.bestTrialRow = this.processedData.find(obj => {
        return obj['trial name'] === this.bestTrialName;
      });
    }
  }

  setData(data: any, displayedColumns: any) {
    const processedData = [];
    for (var i = 0; i < data.length; i++) {
      var list = data[i];
      processedData[i] = {};

      for (var j = 0; j < displayedColumns.length; j++) {
        var key = lowerCase(displayedColumns[j]);
        var value = list[j];
        processedData[i][key] = value;
      }
    }

    return processedData;
  }

  setConfig(displayedColumns: any, processedData: any) {
    const columns = [];
    for (var i = 0; i < displayedColumns.length; i++) {
      if (displayedColumns[i] !== 'Kfp run') {
        if (displayedColumns[i] === 'Trial name') {
          columns.push({
            matHeaderCellDef: displayedColumns[i],
            matColumnDef: 'name',
            value: new PropertyValue({
              field: lowerCase(displayedColumns[i]),
              isLink: true,
            }),
            sort: true,
          });
        } else if (displayedColumns[i] === 'Status') {
          columns.push({
            matHeaderCellDef: displayedColumns[i],
            matColumnDef: displayedColumns[i],
            value: new StatusValue({
              valueFn: parseStatus,
            }),
            sort: true,
          });
        } else {
          columns.push({
            matHeaderCellDef: displayedColumns[i],
            matColumnDef: displayedColumns[i],
            value: new PropertyValue({
              field: lowerCase(displayedColumns[i]),
            }),
            sort: true,
          });
        }
      }
    }

    let kfpRunExists = false;
    for (var i = 0; i < processedData.length; i++) {
      if (processedData[i]['kfp run']) {
        kfpRunExists = true;
      }
    }

    if (kfpRunExists) {
      columns.push({
        matHeaderCellDef: '',
        matColumnDef: 'actions',
        value: new ComponentValue({
          component: KfpRunComponent,
        }),
      });
    }

    return {
      columns,
    };
  }

  // Event handling functions
  reactToAction(a: ActionEvent) {
    switch (a.action) {
      case 'name:link':
        this.openTrialDetails(a.data['trial name']);
        break;
    }
  }

  openTrialDetails(name: string) {
    this.router.navigate([`/experiment/${this.experimentName}/trial/${name}`]);
  }
}
