import { CommonModule } from '@angular/common';
import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { FormArray } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { FormModule } from 'kubeflow';

import { FormNasOperationsComponent } from './nas-operations.component';
import { FormNasOperationsModule } from './nas-operations.module';

describe('FormNasOperationsComponent', () => {
  let component: FormNasOperationsComponent;
  let fixture: ComponentFixture<FormNasOperationsComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          CommonModule,
          BrowserAnimationsModule,
          MatIconModule,
          FormModule,
          FormNasOperationsModule,
        ],
        declarations: [FormNasOperationsComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(FormNasOperationsComponent);
    component = fixture.componentInstance;
    component.formArray = new FormArray([]);
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
