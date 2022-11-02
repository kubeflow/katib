import { HttpClientModule } from '@angular/common/http';
import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { ReactiveFormsModule } from '@angular/forms';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

import { TrialModalComponent } from './trial-modal.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { KWABackendService } from 'src/app/services/backend.service';
import { of } from 'rxjs';
import { ActivatedRoute, Router } from '@angular/router';
import { NamespaceService, TitleActionsToolbarModule } from 'kubeflow';

let KWABackendServiceStub: Partial<KWABackendService>;
let NamespaceServiceStub: Partial<NamespaceService>;

KWABackendServiceStub = {
  getTrial: () => of([]),
  getTrialInfo: () => of(),
};

NamespaceServiceStub = {
  getSelectedNamespace: () => of(),
};

describe('TrialModalComponent', () => {
  let component: TrialModalComponent;
  let fixture: ComponentFixture<TrialModalComponent>;
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
          MatProgressSpinnerModule,
          BrowserAnimationsModule,
          TitleActionsToolbarModule,
        ],
        declarations: [TrialModalComponent],
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
    fixture = TestBed.createComponent(TrialModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
