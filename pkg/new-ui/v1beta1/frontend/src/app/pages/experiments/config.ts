import {
  PropertyValue,
  StatusValue,
  ActionListValue,
  ActionIconValue,
  TRUNCATE_TEXT_SIZE,
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
  title: 'Experiments',
  newButtonText: 'NEW EXPERIMENT',
  columns: [
    {
      matHeaderCellDef: 'Status',
      matColumnDef: 'status',
      value: new StatusValue({
        valueFn: parseStatus,
      }),
    },
    {
      matHeaderCellDef: 'Name',
      matColumnDef: 'name',
      value: new PropertyValue({
        field: 'name',
        isLink: true,
      }),
    },
    {
      matHeaderCellDef: 'Age',
      matColumnDef: 'age',
      value: new DateTimeValue({
        field: 'startTime',
      }),
    },
    {
      matHeaderCellDef: 'Successful trials',
      matColumnDef: 'successful',
      textAlignment: 'right',
      value: new PropertyValue({
        valueFn: parseSucceededTrials,
      }),
    },
    {
      matHeaderCellDef: 'Running trials',
      matColumnDef: 'running',
      textAlignment: 'right',
      value: new PropertyValue({
        valueFn: parseRunningTrials,
      }),
    },
    {
      matHeaderCellDef: 'Failed trials',
      matColumnDef: 'failed',
      textAlignment: 'right',
      value: new PropertyValue({
        valueFn: parseFailedTrials,
      }),
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
