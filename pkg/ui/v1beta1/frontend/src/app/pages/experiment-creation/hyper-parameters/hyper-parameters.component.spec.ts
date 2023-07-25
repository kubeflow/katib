import { CommonModule } from '@angular/common';
import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { FormArray, FormControl } from '@angular/forms';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { FormModule } from 'kubeflow';
import { ParamsListModule } from 'src/app/shared/params-list/params-list.module';

import { FormHyperParametersComponent } from './hyper-parameters.component';

describe('FormHyperParametersComponent', () => {
  let component: FormHyperParametersComponent;
  let fixture: ComponentFixture<FormHyperParametersComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          CommonModule,
          ParamsListModule,
          FormModule,
          BrowserAnimationsModule,
        ],
        declarations: [FormHyperParametersComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(FormHyperParametersComponent);
    component = fixture.componentInstance;
    component.hyperParamsArray = new FormArray([]);
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
