import { Component, OnDestroy, OnInit } from '@angular/core';
import { environment } from '@app/environment';
import { Subscription } from 'rxjs';
import isEqual from 'lodash-es/isEqual';
import {
  ConfirmDialogService,
  DIALOG_RESP,
  ExponentialBackoff,
  NamespaceService,
  ActionEvent,
  DashboardState,
  ToolbarButton,
} from 'kubeflow';

import { KWABackendService } from 'src/app/services/backend.service';
import {
  Experiment,
  ExperimentsProcessed,
} from '../../models/experiment.model';
import { experimentsTableConfig } from './config';
import { getDeleteDialogConfig } from './delete-modal-config';
import { Router } from '@angular/router';

@Component({
  selector: 'app-experiments',
  templateUrl: './experiments.component.html',
  styleUrls: ['./experiments.component.scss'],
})
export class ExperimentsComponent implements OnInit, OnDestroy {
  experiments: ExperimentsProcessed = [];
  currNamespace: string;
  config = experimentsTableConfig;
  env = environment;
  dashboardDisconnectedState = DashboardState.Disconnected;

  private subs = new Subscription();
  private poller: ExponentialBackoff;

  private rawData: Experiment[] = [];

  buttons: ToolbarButton[] = [
    new ToolbarButton({
      text: `New Experiment`,
      icon: 'add',
      stroked: true,
      fn: () => {
        this.router.navigate(['/new']);
      },
    }),
  ];

  constructor(
    private backend: KWABackendService,
    private confirmDialog: ConfirmDialogService,
    private router: Router,
    public ns: NamespaceService,
  ) {}

  ngOnInit() {
    this.startExperimentsPolling();

    // Reset the poller whenever the selected namespace changes
    this.subs.add(
      this.ns.getSelectedNamespace().subscribe(nameSpace => {
        this.currNamespace = nameSpace;
        this.poller.reset();
      }),
    );
  }

  ngOnDestroy(): void {
    this.subs.unsubscribe();
    this.poller.stop();
  }

  trackByFn(index: number, experiment: Experiment) {
    return `${experiment.name}/${experiment.namespace}`;
  }

  reactToAction(a: ActionEvent) {
    const exp = a.data as Experiment;
    switch (a.action) {
      case 'delete':
        this.onDeleteExperiment(exp.name);
        break;
    }
  }

  onDeleteExperiment(name: string) {
    const deleteDialogConfig = getDeleteDialogConfig(name, this.currNamespace);
    const ref = this.confirmDialog.open(name, deleteDialogConfig);

    const delSub = ref.componentInstance.applying$.subscribe(applying => {
      if (!applying) {
        return;
      }

      // Close the open dialog only if the DELETE request succeeded
      this.backend.deleteExperiment(name, this.currNamespace).subscribe({
        next: _ => {
          this.poller.reset();
          ref.close(DIALOG_RESP.ACCEPT);
        },
        error: err => {
          const errorMsg = err;
          deleteDialogConfig.error = errorMsg;
          ref.componentInstance.applying$.next(false);
        },
      });

      // DELETE request has succeeded
      ref.afterClosed().subscribe(res => {
        delSub.unsubscribe();
        if (res !== DIALOG_RESP.ACCEPT) {
          return;
        }

        this.poller.reset();
      });
    });
  }

  private startExperimentsPolling() {
    this.poller = new ExponentialBackoff({ interval: 1000, retries: 3 });

    // Poll for new data and reset the poller if different data is found
    this.subs.add(
      this.poller.start().subscribe(() => {
        if (!this.currNamespace) {
          return;
        }

        this.backend
          .getExperiments(this.currNamespace)
          .subscribe(experiments => {
            if (isEqual(this.rawData, experiments)) {
              return;
            }

            this.experiments = experiments.map(row => {
              return {
                ...row,
                link: {
                  text: row.name,
                  url: `/experiment/${row.namespace}/${row.name}`,
                },
              };
            });

            this.rawData = experiments;
            this.poller.reset();
          });
      }),
    );
  }
}
