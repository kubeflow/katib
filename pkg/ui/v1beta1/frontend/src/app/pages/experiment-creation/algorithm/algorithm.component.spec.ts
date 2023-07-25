import { CommonModule } from '@angular/common';
import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { FormArray, FormControl, FormGroup } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatSelectModule } from '@angular/material/select';
import { MatRadioModule } from '@angular/material/radio';
import { FormModule } from 'kubeflow';

import { FormAlgorithmComponent } from './algorithm.component';
import { FormAlgorithmModule } from './algorithm.module';

describe('FormAlgorithmComponent', () => {
  let component: FormAlgorithmComponent;
  let fixture: ComponentFixture<FormAlgorithmComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          CommonModule,
          BrowserAnimationsModule,
          MatFormFieldModule,
          MatInputModule,
          MatSelectModule,
          MatRadioModule,
          FormModule,
          FormAlgorithmModule,
        ],
        declarations: [FormAlgorithmComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(FormAlgorithmComponent);
    component = fixture.componentInstance;
    component.algorithmForm = new FormGroup({
      algorithmSettings: new FormArray([]),
      algorithm: new FormControl('tpe'),
      type: new FormControl(),
    });
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
