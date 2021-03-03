import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TrialThresholdsComponent } from './trial-thresholds.component';

describe('TrialThresholdsComponent', () => {
  let component: TrialThresholdsComponent;
  let fixture: ComponentFixture<TrialThresholdsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [TrialThresholdsComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TrialThresholdsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
