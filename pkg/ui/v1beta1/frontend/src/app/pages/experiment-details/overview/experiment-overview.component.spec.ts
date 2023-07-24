import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { CommonModule } from '@angular/common';
import { ConditionsTableModule, DetailsListModule } from 'kubeflow';
import { MatSnackBarModule } from '@angular/material/snack-bar';

import { ExperimentOverviewComponent } from './experiment-overview.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

describe('ExperimentOverviewComponent', () => {
  let component: ExperimentOverviewComponent;
  let fixture: ComponentFixture<ExperimentOverviewComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          CommonModule,
          ConditionsTableModule,
          DetailsListModule,
          MatSnackBarModule,
          BrowserAnimationsModule,
        ],
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
