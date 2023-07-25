import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { FormArray } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { ListInputComponent } from './list-input.component';

describe('ListInputComponent', () => {
  let component: ListInputComponent;
  let fixture: ComponentFixture<ListInputComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [BrowserAnimationsModule, MatIconModule],
        declarations: [ListInputComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(ListInputComponent);
    component = fixture.componentInstance;
    component.formArray = new FormArray([]);
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
