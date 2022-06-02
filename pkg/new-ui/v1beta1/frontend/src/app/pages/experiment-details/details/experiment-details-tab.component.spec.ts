import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { CommonModule } from '@angular/common';
import { DetailsListModule, HeadingSubheadingRowModule } from 'kubeflow';
import { MatSnackBarModule } from '@angular/material/snack-bar';

import { ExperimentDetailsTabComponent } from './experiment-details-tab.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

describe('ExperimentDetailsTabComponent', () => {
  let component: ExperimentDetailsTabComponent;
  let fixture: ComponentFixture<ExperimentDetailsTabComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          CommonModule,
          DetailsListModule,
          HeadingSubheadingRowModule,
          MatSnackBarModule,
          BrowserAnimationsModule,
        ],
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
