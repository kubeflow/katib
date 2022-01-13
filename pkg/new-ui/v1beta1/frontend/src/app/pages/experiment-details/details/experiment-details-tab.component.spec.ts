import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ExperimentDetailsTabComponent } from './experiment-details-tab.component';

describe('ExperimentDetailsTabComponent', () => {
  let component: ExperimentDetailsTabComponent;
  let fixture: ComponentFixture<ExperimentDetailsTabComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [ExperimentDetailsTabComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(ExperimentDetailsTabComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
