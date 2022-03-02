import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ExperimentCreationComponent } from './experiment-creation.component';

describe('ExperimentCreationComponent', () => {
  let component: ExperimentCreationComponent;
  let fixture: ComponentFixture<ExperimentCreationComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [ExperimentCreationComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(ExperimentCreationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
