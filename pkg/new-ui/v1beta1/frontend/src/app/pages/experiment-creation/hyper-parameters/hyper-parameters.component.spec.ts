import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { HyperParametersComponent } from './hyper-parameters.component';

describe('HyperParametersComponent', () => {
  let component: HyperParametersComponent;
  let fixture: ComponentFixture<HyperParametersComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [HyperParametersComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(HyperParametersComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
