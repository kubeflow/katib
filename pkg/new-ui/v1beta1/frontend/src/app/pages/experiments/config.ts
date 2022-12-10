import {
  PropertyValue,
  StatusValue,
  ActionListValue,
  ActionIconValue,
  TableConfig,
  DateTimeValue,
  TemplateValue,
  ChipsListValue,
  ComponentValue,
} from 'kubeflow';
import {
  parseStatus,
  parseSucceededTrials,
  parseRunningTrials,
  parseFailedTrials,
  parseTotalTrials,
} from './utils';
import { ExperimentOptimalTrialComponent } from './optimal-trial/experiment-optimal-trial.component';

export const experimentsTableConfig: TableConfig = {
  columns: [
    {
      matHeaderCellDef: 'Status',
      matColumnDef: 'status',
      value: new StatusValue({
        valueFn: parseStatus,
      }),
      sort: true,
    },
    {
      matHeaderCellDef: 'Name',
      matColumnDef: 'name',
      value: new PropertyValue({
        field: 'name',
        isLink: true,
      }),
      sort: true,
    },
    {
      matHeaderCellDef: 'Created at',
      matColumnDef: 'age',
      textAlignment: 'right',
      value: new DateTimeValue({
        field: 'startTime',
      }),
      sort: true,
    },
    {
      matHeaderCellDef: 'Successful trials',
      matColumnDef: 'successful',
      textAlignment: 'right',
      value: new PropertyValue({
        valueFn: parseSucceededTrials,
      }),
      sort: true,
    },
    {
      matHeaderCellDef: 'Running trials',
      matColumnDef: 'running',
      textAlignment: 'right',
      value: new PropertyValue({
        valueFn: parseRunningTrials,
      }),
      sort: true,
    },
    {
      matHeaderCellDef: 'Failed trials',
      matColumnDef: 'failed',
      textAlignment: 'right',
      value: new PropertyValue({
        valueFn: parseFailedTrials,
      }),
      sort: true,
    },
    {
      matHeaderCellDef: 'Optimal trial',
      matColumnDef: 'optimal',
      value: new ComponentValue({
        component: ExperimentOptimalTrialComponent,
      }),
    },
    {
      matHeaderCellDef: '',
      matColumnDef: 'actions',
      value: new ActionListValue([
        new ActionIconValue({
          name: 'delete',
          tooltip: 'Delete experiment',
          color: '',
          iconReady: 'material:delete',
        }),
      ]),
    },
  ],
};
