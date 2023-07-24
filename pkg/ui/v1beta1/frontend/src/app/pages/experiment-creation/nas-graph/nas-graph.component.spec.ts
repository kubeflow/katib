import { CommonModule } from '@angular/common';
import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { FormControl, FormGroup } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatDividerModule } from '@angular/material/divider';
import { MatIconModule } from '@angular/material/icon';
import { ListInputModule } from 'src/app/shared/list-input/list-input.module';
import { FormModule } from 'kubeflow';

import { FormNasGraphComponent } from './nas-graph.component';

describe('FormNasGraphComponent', () => {
  let component: FormNasGraphComponent;
  let fixture: ComponentFixture<FormNasGraphComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          CommonModule,
          BrowserAnimationsModule,
          MatFormFieldModule,
          MatInputModule,
          MatDividerModule,
          MatIconModule,
          ListInputModule,
          FormModule,
        ],
        declarations: [FormNasGraphComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(FormNasGraphComponent);
    component = fixture.componentInstance;
    component.formGroup = new FormGroup({
      layers: new FormControl(),
      inputSizes: new FormControl([]),
      outputSizes: new FormControl([]),
    });
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
