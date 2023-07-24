import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { FormArray, FormControl, FormGroup } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { CommonModule } from '@angular/common';
import { FormModule } from 'kubeflow';
import { MatIconModule } from '@angular/material/icon';
import { ListInputModule } from 'src/app/shared/list-input/list-input.module';
import { MatDividerModule } from '@angular/material/divider';
import { MatCheckboxModule } from '@angular/material/checkbox';

import { FormObjectiveComponent } from './objective.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

describe('FormObjectiveComponent', () => {
  let component: FormObjectiveComponent;
  let fixture: ComponentFixture<FormObjectiveComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          CommonModule,
          BrowserAnimationsModule,
          FormModule,
          MatIconModule,
          ListInputModule,
          MatDividerModule,
          MatCheckboxModule,
        ],
        declarations: [FormObjectiveComponent],
        providers: [{ provide: MatDialog, useValue: {} }],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(FormObjectiveComponent);
    component = fixture.componentInstance;
    let array: FormGroup[] = [];
    array.push(new FormGroup({}));
    let formArray = new FormArray(array);
    component.objectiveForm = new FormGroup({
      strategiesArray: formArray,
      type: new FormControl(),
      metricName: new FormControl(),
      goal: new FormControl(),
      setStrategies: new FormControl(),
      additionalMetricNames: new FormControl([]),
      metricStrategy: new FormControl(),
      metricStrategies: new FormArray([]),
    });
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
