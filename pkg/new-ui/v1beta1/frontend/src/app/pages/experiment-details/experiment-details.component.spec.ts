import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ExperimentDetailsComponent } from './experiment-details.component';

describe('ExperimentDetailsComponent', () => {
  let component: ExperimentDetailsComponent;
  let fixture: ComponentFixture<ExperimentDetailsComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [ExperimentDetailsComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(ExperimentDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
