import { HttpClientModule } from '@angular/common/http';
import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { ReactiveFormsModule } from '@angular/forms';
import { MatSnackBarModule } from '@angular/material/snack-bar';

import { TrialDetailsComponent } from './trial-details.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { KWABackendService } from 'src/app/services/backend.service';
import { of } from 'rxjs';
import { ActivatedRoute, Router } from '@angular/router';
import {
  LoadingSpinnerModule,
  NamespaceService,
  TitleActionsToolbarModule,
} from 'kubeflow';
import { NgxEchartsModule } from 'ngx-echarts';

let KWABackendServiceStub: Partial<KWABackendService>;
let NamespaceServiceStub: Partial<NamespaceService>;

KWABackendServiceStub = {
  getTrial: () => of([]),
  getTrialInfo: () => of(),
  getTrialLogs: () => of(),
};

NamespaceServiceStub = {
  getSelectedNamespace: () => of(),
};

describe('TrialDetailsComponent', () => {
  let component: TrialDetailsComponent;
  let fixture: ComponentFixture<TrialDetailsComponent>;
  let activatedRouteSpy;

  beforeEach(
    waitForAsync(() => {
      activatedRouteSpy = {
        snapshot: {
          params: {
            trialName: '',
            experimentName: '',
          },
        },
      };
      TestBed.configureTestingModule({
        imports: [
          HttpClientModule,
          ReactiveFormsModule,
          MatSnackBarModule,
          BrowserAnimationsModule,
          TitleActionsToolbarModule,
          LoadingSpinnerModule,
          NgxEchartsModule.forRoot({
            echarts: () => import('echarts'),
          }),
        ],
        declarations: [TrialDetailsComponent],
        providers: [
          { provide: ActivatedRoute, useValue: activatedRouteSpy },
          { provide: Router, useValue: {} },
          { provide: KWABackendService, useValue: KWABackendServiceStub },
          { provide: NamespaceService, useValue: NamespaceServiceStub },
        ],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(TrialDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
