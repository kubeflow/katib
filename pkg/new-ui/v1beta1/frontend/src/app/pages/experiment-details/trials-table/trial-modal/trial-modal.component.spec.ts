import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TrialModalComponent } from './trial-modal.component';

describe('TrialModalComponent', () => {
  let component: TrialModalComponent;
  let fixture: ComponentFixture<TrialModalComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [TrialModalComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TrialModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
