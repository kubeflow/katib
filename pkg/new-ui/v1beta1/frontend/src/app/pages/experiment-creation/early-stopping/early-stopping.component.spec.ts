import { CommonModule } from '@angular/common';
import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { FormControl, FormGroup } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { FormModule } from 'kubeflow';

import { EarlyStoppingComponent } from './early-stopping.component';

describe('EarlyStoppingComponent', () => {
  let component: EarlyStoppingComponent;
  let fixture: ComponentFixture<EarlyStoppingComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          CommonModule,
          BrowserAnimationsModule,
          MatFormFieldModule,
          MatInputModule,
          FormModule,
        ],
        declarations: [EarlyStoppingComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(EarlyStoppingComponent);
    component = fixture.componentInstance;
    component.formGroup = new FormGroup({
      algorithmName: new FormControl(),
    });
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
