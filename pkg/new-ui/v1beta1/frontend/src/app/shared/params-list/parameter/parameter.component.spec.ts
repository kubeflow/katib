import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { FormControl, FormGroup } from '@angular/forms';
import {
  MatDialogModule,
  MatDialogRef,
  MAT_DIALOG_DATA,
} from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { DetailsListModule, PopoverModule } from 'kubeflow';
import { MatTooltipModule } from '@angular/material/tooltip';

import { ParameterComponent } from './parameter.component';

describe('ParameterComponent', () => {
  let component: ParameterComponent;
  let fixture: ComponentFixture<ParameterComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          MatDialogModule,
          MatIconModule,
          PopoverModule,
          MatTooltipModule,
          DetailsListModule,
        ],
        declarations: [ParameterComponent],
        providers: [
          { provide: MAT_DIALOG_DATA, useValue: {} },
          { provide: MatDialogRef, useValue: {} },
        ],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(ParameterComponent);
    component = fixture.componentInstance;
    component.paramFormGroup = new FormGroup({
      name: new FormControl(),
      parameterType: new FormControl(),
      feasibleSpace: new FormControl({ min: 0, max: 10, step: '' }),
    });
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
