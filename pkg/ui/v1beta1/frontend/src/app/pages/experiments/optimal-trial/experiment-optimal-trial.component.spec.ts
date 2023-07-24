import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { MatDividerModule } from '@angular/material/divider';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { DetailsListModule, PopoverModule } from 'kubeflow';

import { ExperimentOptimalTrialComponent } from './experiment-optimal-trial.component';

describe('ExperimentOptimalTrialComponent', () => {
  let component: ExperimentOptimalTrialComponent;
  let fixture: ComponentFixture<ExperimentOptimalTrialComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          BrowserAnimationsModule,
          PopoverModule,
          DetailsListModule,
          MatDividerModule,
        ],
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
