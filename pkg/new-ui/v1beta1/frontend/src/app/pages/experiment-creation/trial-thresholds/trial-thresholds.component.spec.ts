import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { CommonModule } from '@angular/common';
import { FormModule } from 'kubeflow';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';

import { FormTrialThresholdsComponent } from './trial-thresholds.component';
import { FormControl, FormGroup } from '@angular/forms';

describe('TrialThresholdsComponent', () => {
  let component: FormTrialThresholdsComponent;
  let fixture: ComponentFixture<FormTrialThresholdsComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          CommonModule,
          BrowserAnimationsModule,
          MatFormFieldModule,
          MatInputModule,
          MatSelectModule,
          FormModule,
        ],
        declarations: [FormTrialThresholdsComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(FormTrialThresholdsComponent);
    component = fixture.componentInstance;
    component.formGroup = new FormGroup({
      parallelTrialCount: new FormControl(),
      maxTrialCount: new FormControl(),
      maxFailedTrialCount: new FormControl(),
      resumePolicy: new FormControl(),
    });
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
