import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { CommonModule } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';
import { MatDividerModule } from '@angular/material/divider';
import { MatRadioModule } from '@angular/material/radio';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';

import { FormModule, PopoverModule, DetailsListModule } from 'kubeflow';
import { AddParamModalComponent } from './add-modal.component';
import {
  MatDialogModule,
  MatDialogRef,
  MAT_DIALOG_DATA,
} from '@angular/material/dialog';
import { FormControl, FormGroup } from '@angular/forms';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

describe('AddParamModalComponent', () => {
  let component: AddParamModalComponent;
  let fixture: ComponentFixture<AddParamModalComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          CommonModule,
          FormModule,
          MatIconModule,
          MatDividerModule,
          MatRadioModule,
          PopoverModule,
          DetailsListModule,
          MatSlideToggleModule,
          MatDialogModule,
          BrowserAnimationsModule,
        ],
        declarations: [AddParamModalComponent],
        providers: [
          {
            provide: MAT_DIALOG_DATA,
            useValue: new FormGroup({
              name: new FormControl(),
              parameterType: new FormControl('int'),
              feasibleSpace: new FormControl({ min: '1', max: '64', step: '' }),
            }),
          },
          { provide: MatDialogRef, useValue: {} },
        ],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(AddParamModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
