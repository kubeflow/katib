import { CommonModule } from '@angular/common';
import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { PopoverModule } from 'kubeflow';

import { TrialParameterComponent } from './trial-parameter.component';

describe('TrialParameterComponent', () => {
  let component: TrialParameterComponent;
  let fixture: ComponentFixture<TrialParameterComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          CommonModule,
          BrowserAnimationsModule,
          MatFormFieldModule,
          MatInputModule,
          MatIconModule,
          PopoverModule,
          ReactiveFormsModule,
        ],
        declarations: [TrialParameterComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(TrialParameterComponent);
    component = fixture.componentInstance;
    component.formGroup = new FormGroup({
      name: new FormControl(),
      reference: new FormControl(),
      description: new FormControl(),
    });
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
