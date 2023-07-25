import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { CommonModule } from '@angular/common';

import { MatIconModule } from '@angular/material/icon';
import { MatDividerModule } from '@angular/material/divider';
import { MatRadioModule } from '@angular/material/radio';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';

import { FormModule, PopoverModule, DetailsListModule } from 'kubeflow';
import { ListInputModule } from '../list-input/list-input.module';

import { ParamsListComponent } from './params-list.component';
import { FormArray, FormControl, ReactiveFormsModule } from '@angular/forms';

describe('ParamsListComponent', () => {
  let component: ParamsListComponent;
  let fixture: ComponentFixture<ParamsListComponent>;

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
          ListInputModule,
          ReactiveFormsModule,
        ],
        declarations: [ParamsListComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(ParamsListComponent);
    component = fixture.componentInstance;
    component.paramsArray = new FormArray([]);
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
