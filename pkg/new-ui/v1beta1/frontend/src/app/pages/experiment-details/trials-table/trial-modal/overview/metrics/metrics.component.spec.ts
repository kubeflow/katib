import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { ConditionsTableModule, DetailsListModule } from 'kubeflow';

import { TrialModalMetricsComponent } from './metrics.component';

describe('TrialModalMetricsComponent', () => {
  let component: TrialModalMetricsComponent;
  let fixture: ComponentFixture<TrialModalMetricsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [TrialModalMetricsComponent],
      imports: [ConditionsTableModule, DetailsListModule, MatSnackBarModule],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TrialModalMetricsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
