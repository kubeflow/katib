import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { ActivatedRoute, Router } from '@angular/router';
import { KWABackendService } from 'src/app/services/backend.service';

import { ExperimentDetailsComponent } from './experiment-details.component';
import { of } from 'rxjs';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTabsModule } from '@angular/material/tabs';
import {
  ConfirmDialogService,
  NamespaceService,
  TitleActionsToolbarModule,
  LoadingSpinnerModule,
  PanelModule,
} from 'kubeflow';
import { TrialsTableModule } from './trials-table/trials-table.module';
import { ExperimentOverviewModule } from './overview/experiment-overview.module';
import { ExperimentDetailsTabModule } from './details/experiment-details-tab.module';
import { ExperimentYamlModule } from './yaml/experiment-yaml.module';
import { TrialsGraphEchartsModule } from './trials-graph-echarts/trials-graph-echarts.module';
import { ReactiveFormsModule } from '@angular/forms';
import { MatSnackBarModule } from '@angular/material/snack-bar';

let KWABackendServiceStub: Partial<KWABackendService>;
let NamespaceServiceStub: Partial<NamespaceService>;
let ActivatedRouteStub: Partial<ActivatedRoute>;

KWABackendServiceStub = {
  getExperimentTrialsInfo: () =>
    of(
      'Status,trialName,Validation-accuracy,Train-accuracy,lr,num-layers,optimizer\nSucceeded,tpe-05daf02d,0.977807,0.993104,0.023222418198803642,3,sgd',
    ),
  getExperiment: () => of(),
  deleteExperiment: () => of(),
};

NamespaceServiceStub = {
  updateSelectedNamespace: () => {},
};

ActivatedRouteStub = {
  params: of({ namespace: 'kubeflow-user', experimentName: '' }),
  queryParams: of({}),
};

describe('ExperimentDetailsComponent', () => {
  let component: ExperimentDetailsComponent;
  let fixture: ComponentFixture<ExperimentDetailsComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          CommonModule,
          TrialsTableModule,
          MatButtonModule,
          MatTabsModule,
          MatIconModule,
          LoadingSpinnerModule,
          PanelModule,
          ExperimentOverviewModule,
          ExperimentDetailsTabModule,
          ExperimentYamlModule,
          TitleActionsToolbarModule,
          TrialsGraphEchartsModule,
          ReactiveFormsModule,
          MatSnackBarModule,
          TitleActionsToolbarModule,
        ],
        declarations: [ExperimentDetailsComponent],
        providers: [
          { provide: ActivatedRoute, useValue: ActivatedRouteStub },
          { provide: Router, useValue: {} },
          { provide: KWABackendService, useValue: KWABackendServiceStub },
          { provide: ConfirmDialogService, useValue: {} },
          { provide: NamespaceService, useValue: NamespaceServiceStub },
        ],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(ExperimentDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
