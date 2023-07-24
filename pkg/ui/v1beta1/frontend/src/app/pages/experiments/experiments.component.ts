import { Component, OnDestroy, OnInit } from '@angular/core';
import { environment } from '@app/environment';
import { Subscription } from 'rxjs';
import {
  ConfirmDialogService,
  DIALOG_RESP,
  NamespaceService,
  ActionEvent,
  DashboardState,
  ToolbarButton,
  PollerService,
} from 'kubeflow';

import { KWABackendService } from 'src/app/services/backend.service';
import {
  Experiment,
  ExperimentsProcessed,
} from '../../models/experiment.model';
import { experimentsTableConfig } from './config';
import { generateDeleteConfig } from './delete-modal-config';
import { Router } from '@angular/router';

@Component({
  selector: 'app-experiments',
  templateUrl: './experiments.component.html',
  styleUrls: ['./experiments.component.scss'],
})
export class ExperimentsComponent implements OnInit, OnDestroy {
  env = environment;

  nsSub = new Subscription();
  pollSub = new Subscription();

  currNamespace: string | string[];
  config = experimentsTableConfig;
  experiments: ExperimentsProcessed = [];

  dashboardDisconnectedState = DashboardState.Disconnected;

  private newExperimentButton = new ToolbarButton({
    text: $localize`New Experiment`,
    icon: 'add',
    stroked: true,
    fn: () => {
      this.router.navigate(['/new']);
    },
  });

  buttons: ToolbarButton[] = [this.newExperimentButton];

  constructor(
    private backend: KWABackendService,
    private confirmDialog: ConfirmDialogService,
    private router: Router,
    public ns: NamespaceService,
    public poller: PollerService,
  ) {}

  ngOnInit(): void {
    // Reset the poller whenever the selected namespace changes
    this.nsSub = this.ns.getSelectedNamespace2().subscribe(ns => {
      this.currNamespace = ns;
      this.poll(ns);
      this.newExperimentButton.namespaceChanged(ns, $localize`Experiment`);
    });
  }

  ngOnDestroy() {
    this.nsSub.unsubscribe();
    this.pollSub.unsubscribe();
  }

  public poll(ns: string | string[]) {
    this.pollSub.unsubscribe();
    this.experiments = [];

    const request = this.backend.getExperiments(ns);

    this.pollSub = this.poller.exponential(request).subscribe(experiments => {
      this.experiments = experiments.map(row => {
        return {
          ...row,
          link: {
            text: row.name,
            url: `/experiment/${this.currNamespace}/${row.name}`,
          },
        };
      });
    });
  }

  trackByFn(index: number, experiment: Experiment) {
    return `${experiment.name}/${experiment.namespace}`;
  }

  reactToAction(a: ActionEvent) {
    const exp = a.data as Experiment;
    switch (a.action) {
      case 'delete':
        this.deleteClicked(exp);
        break;
    }
  }

  deleteClicked(exp: Experiment) {
    const deleteDialogConfig = generateDeleteConfig(exp.name);
    const ref = this.confirmDialog.open(exp.name, deleteDialogConfig);

    const delSub = ref.componentInstance.applying$.subscribe(applying => {
      if (!applying) {
        return;
      }

      // Close the open dialog only if the DELETE request succeeded
      this.backend.deleteExperiment(exp.name, exp.namespace).subscribe(
        res => {
          ref.close(DIALOG_RESP.ACCEPT);
        },
        err => {
          deleteDialogConfig.error = err;
          ref.componentInstance.applying$.next(false);
        },
      );

      // DELETE request has succeeded
      ref.afterClosed().subscribe(res => {
        delSub.unsubscribe();
        if (res !== DIALOG_RESP.ACCEPT) {
          return;
        }
      });
    });
  }
}
