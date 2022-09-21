import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ExperimentOptimalTrialComponent } from './experiment-optimal-trial.component';

describe('ExperimentOptimalTrialComponent', () => {
  let component: ExperimentOptimalTrialComponent;
  let fixture: ComponentFixture<ExperimentOptimalTrialComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [ExperimentOptimalTrialComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(ExperimentOptimalTrialComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
