import {
  Component,
  OnDestroy,
  OnInit,
  ViewChild,
  TemplateRef,
} from '@angular/core';
import { environment } from '@app/environment';
import { FormBuilder, FormControl, FormGroup } from '@angular/forms';
import { combineLatest, Subscription } from 'rxjs';
import isEqual from 'lodash-es/isEqual';
import { startWith } from 'rxjs/operators';
import {
  ConfirmDialogService,
  DIALOG_RESP,
  ExponentialBackoff,
  NamespaceService,
  TemplateValue,
  ActionEvent,
  DashboardState,
} from 'kubeflow';

import { KWABackendService } from 'src/app/services/backend.service';
import { Experiments, Experiment } from '../../models/experiment.model';
import { StatusEnum } from '../../enumerations/status.enum';
import { experimentsTableConfig } from './config';
import { getDeleteDialogConfig } from './delete-modal-config';
import { Router } from '@angular/router';

@Component({
  selector: 'app-experiments',
  templateUrl: './experiments.component.html',
  styleUrls: ['./experiments.component.scss'],
})
export class ExperimentsComponent implements OnInit, OnDestroy {
  experiments: Experiments = [];
  currNamespace: string;
  config = experimentsTableConfig;
  env = environment;
  dashboardDisconnectedState = DashboardState.Disconnected;

  private subs = new Subscription();
  private poller: ExponentialBackoff;

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
      case 'newResourceButton': // TODO: could also use enums here
        this.router.navigate(['/new']);
        break;
      case 'name:link':
        this.router.navigate([`/experiment/${exp.name}`]);
        break;
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
        this.backend.getExperiments().subscribe(experiments => {
          // the backend should have proper namespace isolation
          experiments = experiments.filter(
            experiment => experiment.namespace === this.currNamespace,
          );

          if (isEqual(this.experiments, experiments)) {
            return;
          }

          this.experiments = experiments;
          this.poller.reset();
        });
      }),
    );
  }
}
