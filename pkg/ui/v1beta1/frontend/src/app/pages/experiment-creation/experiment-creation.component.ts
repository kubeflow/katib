import { Component, OnInit } from '@angular/core';
import { FormGroup, FormArray, FormControl } from '@angular/forms';
import { Router } from '@angular/router';

import { MatDialog } from '@angular/material/dialog';

import { ExperimentFormService } from '../../services/experiment-form.service';
import { YamlModalComponent } from './yaml-modal/yaml-modal.component';
import {
  ExperimentK8s,
  ExperimentSpec,
  FeasibleSpaceMinMax,
  EXPERIMENT_APIVERSION,
  EXPERIMENT_KIND,
} from 'src/app/models/experiment.k8s.model';
import { KWABackendService } from 'src/app/services/backend.service';
import {
  NamespaceService,
  SnackType,
  SnackBarService,
  SnackBarConfig,
} from 'kubeflow';
import { pipe } from 'rxjs';
import { take } from 'rxjs/operators';
import { EarlyStoppingAlgorithmsEnum } from 'src/app/enumerations/algorithms.enum';

@Component({
  selector: 'app-experiment-creation',
  templateUrl: './experiment-creation.component.html',
  styleUrls: ['./experiment-creation.component.scss'],
})
export class ExperimentCreationComponent implements OnInit {
  metadataForm: FormGroup;
  trialThresholdsForm: FormGroup;
  objectiveForm: FormGroup;
  algorithmForm: FormGroup;
  earlyStoppingForm: FormGroup;
  hyperParamsArray: FormArray;
  nasGraphForm: FormGroup;
  nasOperationsForm: FormArray;
  metricsForm: FormGroup;
  trialTemplateForm: FormGroup;

  constructor(
    private formSvc: ExperimentFormService,
    private router: Router,
    private dialog: MatDialog,
    private backend: KWABackendService,
    private ns: NamespaceService,
    private snack: SnackBarService,
  ) {}

  ngOnInit() {
    this.metadataForm = this.formSvc.createMetadataForm('kubeflow-user');
    this.trialThresholdsForm = this.formSvc.createTrialThresholdForm();
    this.objectiveForm = this.formSvc.createObjectiveForm();
    this.algorithmForm = this.formSvc.createAlgorithmObjectiveForm();
    this.earlyStoppingForm = this.formSvc.createEarlyStoppingForm();
    this.hyperParamsArray = this.formSvc.createHyperParametersForm();
    this.nasGraphForm = this.formSvc.createNasGraphForm();
    this.nasOperationsForm = this.formSvc.createNasOperationsForm();
    this.metricsForm = this.formSvc.createMetricsForm();
    this.trialTemplateForm = this.formSvc.createTrialTemplateForm();
  }

  /**
   * Create an Experiment CR json obj from the form's values
   */
  getFormYaml(): any {
    const metadata = { name: '', namespace: '' };
    const spec: ExperimentSpec = {};

    metadata.name = this.metadataForm.value.name;

    const thresholds = this.trialThresholdsForm.value;
    spec.maxTrialCount = thresholds.maxTrialCount;
    spec.parallelTrialCount = thresholds.parallelTrialCount;
    spec.maxFailedTrialCount = thresholds.maxFailedTrialCount;
    spec.resumePolicy = thresholds.resumePolicy;

    spec.objective = this.formSvc.objectiveFromCtrl(this.objectiveForm);
    spec.algorithm = this.formSvc.algorithmFromCtrl(this.algorithmForm);

    const algoType = this.algorithmForm.value.type;
    if (algoType === 'hp') {
      spec.parameters = this.formSvc.hyperParamsFromCtrl(this.hyperParamsArray);
    } else if (algoType === 'nas') {
      spec.nasConfig = {
        graphConfig: this.nasGraphForm.value,
        operations: this.formSvc.nasOpsFromCtrl(this.nasOperationsForm),
      };
    }

    const earlyStoppingAlgo = this.earlyStoppingForm.value.algorithmName;
    if (earlyStoppingAlgo !== EarlyStoppingAlgorithmsEnum.NONE) {
      spec.earlyStopping = this.formSvc.earlyStoppingFromCtrl(
        this.earlyStoppingForm,
      );
    }

    spec.metricsCollectorSpec = this.formSvc.metricsCollectorFromCtrl(
      this.metricsForm,
    );
    spec.trialTemplate = this.formSvc.trialTemplateFromCtrl(
      this.trialTemplateForm,
    );

    return {
      apiVersion: EXPERIMENT_APIVERSION,
      kind: EXPERIMENT_KIND,
      metadata,
      spec,
    };
  }

  private submitExperiment(exp: ExperimentK8s) {
    this.ns
      .getSelectedNamespace()
      .pipe(take(1))
      .subscribe(ns => {
        exp.metadata.namespace = ns;
        console.log(exp);

        this.backend.createExperiment(exp).subscribe({
          next: () => {
            const config: SnackBarConfig = {
              data: {
                msg: 'Experiment submitted successfully.',
                snackType: SnackType.Success,
              },
              duration: 3000,
            };
            this.snack.open(config);
            this.returnToExperiments();
          },
          error: err => {
            console.warn('could not submit experiment');
          },
        });
      });
  }

  showYAML() {
    const formYaml = this.getFormYaml();

    // show the dialog
    const ref = this.dialog.open(YamlModalComponent, {
      data: formYaml,
    });

    ref.afterClosed().subscribe((res: ExperimentK8s) => {
      if (!res) {
        return;
      }

      this.submitExperiment(res);
    });
  }

  create() {
    const exp = this.getFormYaml();
    this.submitExperiment(exp);
  }

  returnToExperiments() {
    this.router.navigate(['']);
  }
}
