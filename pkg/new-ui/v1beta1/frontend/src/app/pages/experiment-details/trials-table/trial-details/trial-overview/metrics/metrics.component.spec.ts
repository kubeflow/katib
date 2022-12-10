import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { ConditionsTableModule, DetailsListModule } from 'kubeflow';

import { TrialMetricsComponent } from './metrics.component';

describe('TrialMetricsComponent', () => {
  let component: TrialMetricsComponent;
  let fixture: ComponentFixture<TrialMetricsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [TrialMetricsComponent],
      imports: [ConditionsTableModule, DetailsListModule, MatSnackBarModule],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TrialMetricsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
