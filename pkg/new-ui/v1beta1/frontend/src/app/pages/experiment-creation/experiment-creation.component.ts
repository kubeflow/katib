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
} from 'src/app/models/experiment.k8s.model';
import { KWABackendService } from 'src/app/services/backend.service';
import { NamespaceService, SnackType, SnackBarService } from 'kubeflow';
import { pipe } from 'rxjs';
import { take } from 'rxjs/operators';

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
  hyperParamsArray: FormArray;
  nasGraphForm: FormGroup;
  nasOperationsForm: FormArray;
  metricsForm: FormGroup;
  trialTemplateForm: FormGroup;
  yamlTemplateForm: FormControl;

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
    this.hyperParamsArray = this.formSvc.createHyperParametersForm();
    this.nasGraphForm = this.formSvc.createNasGraphForm();
    this.nasOperationsForm = this.formSvc.createNasOperationsForm();
    this.metricsForm = this.formSvc.createMetricsForm();
    this.trialTemplateForm = this.formSvc.createTrialTemplateForm();
    this.yamlTemplateForm = this.formSvc.createYamlTemplateForm();
  }

  /**
   * Create an Experiment CR json obj from the form's values
   */
  getFormYaml(): any {
    const metadata = { name: '', namespace: '' };
    const spec: ExperimentSpec = {};

    metadata.name = this.metadataForm.value.name;

    spec.maxTrialCount = this.trialThresholdsForm.value.maxTrialCount;
    spec.parallelTrialCount = this.trialThresholdsForm.value.parallelTrialCount;
    spec.maxFailedTrialCount = this.trialThresholdsForm.value.maxFailedTrialCount;

    spec.objective = this.formSvc.objectiveFromCtrl(this.objectiveForm);
    spec.algorithm = this.formSvc.algorithmFromCtrl(this.algorithmForm);
    spec.parameters = this.formSvc.hyperParamsFromCtrl(this.hyperParamsArray);
    spec.metricsCollectorSpec = this.formSvc.metricsCollectorFromCtrl(
      this.metricsForm,
    );
    spec.trialTemplate = this.formSvc.trialTemplateFromCtrl(
      this.trialTemplateForm,
    );

    return {
      apiVersion: 'kubeflow.org/v1beta1',
      kind: 'Experiment',
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
            this.snack.open(
              'Experiment submitted successfully.',
              SnackType.Success,
              3000,
            );
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
