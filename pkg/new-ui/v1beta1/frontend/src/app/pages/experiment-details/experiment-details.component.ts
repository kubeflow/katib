import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { MatTabChangeEvent } from '@angular/material/tabs';
import {
  ConfirmDialogService,
  DIALOG_RESP,
  ExponentialBackoff,
  getCondition,
  NamespaceService,
  ToolbarButton,
} from 'kubeflow';

import { KWABackendService } from '../../services/backend.service';
import { StatusEnum } from '../../enumerations/status.enum';
import { Subscription } from 'rxjs';
import {
  numberToExponential,
  transformStringResponses,
} from '../../shared/utils';
import { generateDeleteConfig } from '../experiments/delete-modal-config';
import { ExperimentK8s } from '../../models/experiment.k8s.model';

@Component({
  selector: 'app-experiment-details',
  templateUrl: './experiment-details.component.html',
  styleUrls: ['./experiment-details.component.scss'],
})
export class ExperimentDetailsComponent implements OnInit, OnDestroy {
  name: string;
  namespace: string;
  columns: string[] = [];
  details: string[][] = [];
  experimentTrialsCsv: string;
  hoveredTrial: number;
  experimentDetails: ExperimentK8s;
  showGraph: boolean;
  bestTrialName: string;
  pageLoading = true;
  selectedTab = 0;
  tabs = new Map<string, number>([
    ['overview', 0],
    ['trials', 1],
    ['details', 2],
    ['yaml', 3],
  ]);

  constructor(
    private activatedRoute: ActivatedRoute,
    private router: Router,
    private backendService: KWABackendService,
    private confirmDialog: ConfirmDialogService,
    private namespaceService: NamespaceService,
  ) {}

  buttonsConfig: ToolbarButton[] = [
    new ToolbarButton({
      text: 'DELETE',
      icon: 'delete',
      fn: () => {
        this.deleteExperiment(this.name, this.namespace);
      },
    }),
  ];

  private poller: ExponentialBackoff;

  private subs = new Subscription();

  ngOnInit() {
    this.activatedRoute.params.subscribe(params => {
      this.namespaceService.updateSelectedNamespace(params.namespace);

      this.name = params.experimentName;
      this.namespace = params.namespace;

      this.updateExperimentInfo();
    });

    this.activatedRoute.queryParams.subscribe(queryParams => {
      this.selectedTab = this.tabs.get(queryParams.tab);
    });
  }

  tabChanged(event: MatTabChangeEvent) {
    this.selectedTab = event.index;
  }

  ngOnDestroy(): void {
    this.subs.unsubscribe();
  }

  returnToExperiments() {
    this.router.navigate(['']);
  }

  mouseLeftTrial() {
    this.hoveredTrial = null;
  }

  mouseOverTrial = (index: number) => (this.hoveredTrial = index);

  private updateExperimentInfo() {
    this.backendService
      .getExperimentTrialsInfo(this.name, this.namespace)
      .subscribe(response => {
        this.experimentTrialsCsv = response;
        const data = transformStringResponses(response);
        this.columns = data.types;
        this.details = this.parseTrialsDetails(data.details);
        this.showGraph = this.showGraphFn(this.details);
      });
    this.backendService
      .getExperiment(this.name, this.namespace)
      .subscribe((response: ExperimentK8s) => {
        this.experimentDetails = response;
        this.bestTrialName = response.status.currentOptimalTrial
          ? response.status.currentOptimalTrial.bestTrialName
          : '';

        const status = this.experimentStatus(response);

        if (
          status &&
          !(status === StatusEnum.FAILED || status === StatusEnum.SUCCEEDED)
        ) {
          // if the status of the experiment is not succeeded either failed
          // then start polling the trials
          this.startTrialsPolling();
          this.startExperimentsPolling();
        }

        this.pageLoading = false;
      });
  }

  private deleteExperiment(name: string, namespace: string) {
    const deleteDialogConfig = generateDeleteConfig(name);
    const ref = this.confirmDialog.open(name, deleteDialogConfig);

    const delSub = ref.componentInstance.applying$.subscribe(applying => {
      if (!applying) {
        return;
      }

      // Close the open dialog only if the DELETE request succeeded
      this.backendService.deleteExperiment(name, namespace).subscribe({
        next: _ => {
          ref.close(DIALOG_RESP.ACCEPT);
        },
        error: err => {
          deleteDialogConfig.error = err;
          ref.componentInstance.applying$.next(false);
        },
      });

      // DELETE request has succeeded
      ref.afterClosed().subscribe(res => {
        delSub.unsubscribe();
        if (res !== DIALOG_RESP.ACCEPT) {
          return;
        }
        this.returnToExperiments();
      });
    });
  }

  private startTrialsPolling() {
    this.poller = new ExponentialBackoff({
      interval: 5000,
      retries: 1,
      maxInterval: 5001,
    });

    // Poll for new data and reset the poller if different data is found
    this.subs.add(
      this.poller.start().subscribe(() => {
        this.backendService
          .getExperimentTrialsInfo(this.name, this.namespace)
          .subscribe(trials => {
            this.experimentTrialsCsv = trials;
            const data = transformStringResponses(trials);
            this.columns = data.types;
            this.details = this.parseTrialsDetails(data.details);
            this.showGraph = this.showGraphFn(this.details);
          });
      }),
    );
  }

  private startExperimentsPolling() {
    this.poller = new ExponentialBackoff({
      interval: 5000,
      retries: 1,
      maxInterval: 5001,
    });

    // Poll for new data and reset the poller if different data is found
    this.subs.add(
      this.poller.start().subscribe(() => {
        this.backendService
          .getExperiment(this.name, this.namespace)
          .subscribe(response => {
            this.experimentDetails = response;
            this.bestTrialName = response.status.currentOptimalTrial
              ? response.status.currentOptimalTrial.bestTrialName
              : '';
          });
      }),
    );
  }

  private parseTrialsDetails(details: string[][]): string[][] {
    return details.map((detail, index) => {
      const updatedDetail = detail.map(value =>
        isNaN(+value) || value === '' ? value : numberToExponential(+value, 6),
      );
      updatedDetail.push(index.toString());
      return updatedDetail;
    });
  }

  private experimentStatus(experiment: ExperimentK8s): StatusEnum {
    const succeededCondition = getCondition(experiment, StatusEnum.SUCCEEDED);

    if (succeededCondition && succeededCondition.status === 'True') {
      return StatusEnum.SUCCEEDED;
    }

    const failedCondition = getCondition(experiment, StatusEnum.FAILED);

    if (failedCondition && failedCondition.status === 'True') {
      return StatusEnum.FAILED;
    }

    const runningCondition = getCondition(experiment, StatusEnum.RUNNING);

    if (runningCondition && runningCondition.status === 'True') {
      return StatusEnum.RUNNING;
    }

    const restartingCondition = getCondition(experiment, StatusEnum.RESTARTING);

    if (restartingCondition && restartingCondition.status === 'True') {
      return StatusEnum.RESTARTING;
    }

    const createdCondition = getCondition(experiment, StatusEnum.CREATED);

    if (createdCondition && createdCondition.status === 'True') {
      return StatusEnum.CREATED;
    }
  }

  private showGraphFn(array: string[][]): boolean {
    const found = array.find(a => a[0] === 'Succeeded');

    if (found) {
      return true;
    }
  }
}
