import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { FormArray } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';

import { FormModule } from 'kubeflow';

import { ListKeyValueComponent } from './list-key-value.component';

describe('ListKeyValueComponent', () => {
  let component: ListKeyValueComponent;
  let fixture: ComponentFixture<ListKeyValueComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [CommonModule, FormModule, MatIconModule],
        declarations: [ListKeyValueComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(ListKeyValueComponent);
    component = fixture.componentInstance;
    component.formArray = new FormArray([]);
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
