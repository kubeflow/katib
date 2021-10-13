import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { EarlyStoppingComponent } from './early-stopping.component';

describe('EarlyStoppingComponent', () => {
  let component: EarlyStoppingComponent;
  let fixture: ComponentFixture<EarlyStoppingComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [EarlyStoppingComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(EarlyStoppingComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
