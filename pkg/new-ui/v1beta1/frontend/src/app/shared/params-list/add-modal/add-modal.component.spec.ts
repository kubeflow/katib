import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AddModalComponent } from './add-modal.component';

describe('AddModalComponent', () => {
  let component: AddModalComponent;
  let fixture: ComponentFixture<AddModalComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [AddModalComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AddModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
