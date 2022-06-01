import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { TrialsGraphEchartsComponent } from './trials-graph-echarts.component';

describe('TrialsGraphEchartsComponent', () => {
  let component: TrialsGraphEchartsComponent;
  let fixture: ComponentFixture<TrialsGraphEchartsComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [TrialsGraphEchartsComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(TrialsGraphEchartsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
