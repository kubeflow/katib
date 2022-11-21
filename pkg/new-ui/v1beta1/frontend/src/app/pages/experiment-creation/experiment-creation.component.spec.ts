import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { KWABackendService } from 'src/app/services/backend.service';
import { ExperimentFormService } from 'src/app/services/experiment-form.service';
import {
  NamespaceService,
  SnackBarService,
  TitleActionsToolbarModule,
  FormModule,
} from 'kubeflow';

import { ExperimentCreationComponent } from './experiment-creation.component';
import { FormArray, FormControl, FormGroup } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';
import { MatStepperModule } from '@angular/material/stepper';
import { FormMetadataModule } from './metadata/metadata.module';
import { FormTrialThresholdsModule } from './trial-thresholds/trial-thresholds.module';
import { FormObjectiveModule } from './objective/objective.module';
import { FormAlgorithmModule } from './algorithm/algorithm.module';
import { FormHyperParametersModule } from './hyper-parameters/hyper-parameters.module';
import { FormNasGraphModule } from './nas-graph/nas-graph.module';
import { FormNasOperationsModule } from './nas-operations/nas-operations.module';
import { FormMetricsCollectorModule } from './metrics-collector/metrics-collector.module';
import { FormTrialTemplateModule } from './trial-template/trial-template.module';
import { YamlModalModule } from './yaml-modal/yaml-modal.module';
import { FormEarlyStoppingModule } from './early-stopping/early-stopping.module';
import { of } from 'rxjs';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

let ExperimentFormServiceStub: Partial<ExperimentFormService>;
let KWABackendServiceStub: Partial<KWABackendService>;
let NamespaceServiceStub: Partial<NamespaceService>;

ExperimentFormServiceStub = {
  createMetadataForm: () =>
    new FormGroup({
      name: new FormControl(),
      namespace: new FormControl(),
    }),
  createTrialThresholdForm: () =>
    new FormGroup({
      parallelTrialCount: new FormControl(),
      maxTrialCount: new FormControl(),
      maxFailedTrialCount: new FormControl(),
      resumePolicy: new FormControl(),
    }),
  createObjectiveForm: () =>
    new FormGroup({
      strategiesArray: new FormArray([]),
      type: new FormControl(),
      metricName: new FormControl(),
      goal: new FormControl(),
      setStrategies: new FormControl(),
      additionalMetricNames: new FormControl([]),
      metricStrategy: new FormControl('test'),
      metricStrategies: new FormArray([]),
    }),
  createAlgorithmObjectiveForm: () =>
    new FormGroup({
      type: new FormControl(),
      algorithmSettings: new FormArray([]),
      algorithm: new FormControl('tpe'),
    }),
  createEarlyStoppingForm: () => new FormGroup({}),
  createHyperParametersForm: () => new FormArray([]),
  createNasGraphForm: () =>
    new FormGroup({
      layers: new FormControl(),
      inputSizes: new FormControl([]),
      outputSizes: new FormControl([]),
    }),
  createNasOperationsForm: () => new FormArray([]),
  createMetricsForm: () =>
    new FormGroup({
      kind: new FormControl(),
      metricsFile: new FormControl(),
      tfDir: new FormControl(),
      port: new FormControl(),
      path: new FormControl(),
      scheme: new FormControl(),
      host: new FormControl(),
      customYaml: new FormControl(),
    }),
  createTrialTemplateForm: () =>
    new FormGroup({
      trialParameters: new FormArray([]),
      podLabels: new FormControl(),
      containerName: new FormControl(),
      successCond: new FormControl(),
      failureCond: new FormControl(),
      retain: new FormControl(),
      type: new FormControl(),
      cmNamespace: new FormControl(),
      cmName: new FormControl(),
      cmTrialPath: new FormControl(),
    }),
};

KWABackendServiceStub = {
  getTrialTemplates: () => of({ Data: [] }),
};

NamespaceServiceStub = {
  getSelectedNamespace: () => of(),
};

describe('ExperimentCreationComponent', () => {
  let component: ExperimentCreationComponent;
  let fixture: ComponentFixture<ExperimentCreationComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          CommonModule,
          BrowserAnimationsModule,
          MatIconModule,
          MatStepperModule,
          FormModule,
          TitleActionsToolbarModule,
          FormMetadataModule,
          FormTrialThresholdsModule,
          FormObjectiveModule,
          FormAlgorithmModule,
          FormEarlyStoppingModule,
          FormHyperParametersModule,
          FormNasGraphModule,
          FormNasOperationsModule,
          FormMetricsCollectorModule,
          FormTrialTemplateModule,
          YamlModalModule,
        ],
        declarations: [ExperimentCreationComponent],
        providers: [
          {
            provide: ExperimentFormService,
            useValue: ExperimentFormServiceStub,
          },
          { provide: Router, useValue: {} },
          { provide: MatDialog, useValue: {} },
          { provide: KWABackendService, useValue: KWABackendServiceStub },
          { provide: NamespaceService, useValue: NamespaceServiceStub },
          { provide: SnackBarService, useValue: {} },
        ],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(ExperimentCreationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
