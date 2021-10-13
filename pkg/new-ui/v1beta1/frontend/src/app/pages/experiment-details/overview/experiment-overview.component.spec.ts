import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ExperimentOverviewComponent } from './experiment-overview.component';

describe('ExperimentOverviewComponent', () => {
  let component: ExperimentOverviewComponent;
  let fixture: ComponentFixture<ExperimentOverviewComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [ExperimentOverviewComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(ExperimentOverviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
