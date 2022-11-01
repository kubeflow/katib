import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { ExperimentsComponent } from './experiments.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ReactiveFormsModule } from '@angular/forms';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { of } from 'rxjs';
import {
  ConfirmDialogService,
  NamespaceService,
  KubeflowModule,
  PollerService,
} from 'kubeflow';
import { KWABackendService } from 'src/app/services/backend.service';
import { RouterTestingModule } from '@angular/router/testing';

let KWABackendServiceStub: Partial<KWABackendService>;
let NamespaceServiceStub: Partial<NamespaceService>;

KWABackendServiceStub = {
  getExperiments: () => of(),
  deleteExperiment: () => of(),
};

NamespaceServiceStub = {
  getSelectedNamespace: () => of(),
  getSelectedNamespace2: () => of(),
};

describe('ExperimentsComponent', () => {
  let component: ExperimentsComponent;
  let fixture: ComponentFixture<ExperimentsComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          HttpClientTestingModule,
          ReactiveFormsModule,
          MatSnackBarModule,
          KubeflowModule,
          RouterTestingModule,
        ],
        declarations: [ExperimentsComponent],
        providers: [
          { provide: Router, useValue: {} },
          { provide: KWABackendService, useValue: KWABackendServiceStub },
          { provide: ConfirmDialogService, useValue: {} },
          { provide: NamespaceService, useValue: NamespaceServiceStub },
          { provide: PollerService, useValue: {} },
        ],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(ExperimentsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
