import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { MetricsCollectorComponent } from './metrics-collector.component';

describe('MetricsCollectorComponent', () => {
  let component: MetricsCollectorComponent;
  let fixture: ComponentFixture<MetricsCollectorComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [MetricsCollectorComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(MetricsCollectorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
