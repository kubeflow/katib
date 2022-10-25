import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';

import { FormAlgorithmSettingComponent } from './setting.component';
import { CommonModule } from '@angular/common';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';

describe('FormAlgorithmSettingComponent', () => {
  let component: FormAlgorithmSettingComponent;
  let fixture: ComponentFixture<FormAlgorithmSettingComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          CommonModule,
          BrowserAnimationsModule,
          MatFormFieldModule,
          MatInputModule,
          MatSelectModule,
          ReactiveFormsModule,
        ],
        declarations: [FormAlgorithmSettingComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(FormAlgorithmSettingComponent);
    component = fixture.componentInstance;
    component.setting = new FormGroup({
      value: new FormControl(),
      name: new FormControl(),
      type: new FormControl(),
      values: new FormControl([]),
    });
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
