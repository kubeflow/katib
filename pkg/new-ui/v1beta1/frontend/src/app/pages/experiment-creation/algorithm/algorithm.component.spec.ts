import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { AlgorithmComponent } from './algorithm.component';

describe('AlgorithmComponent', () => {
  let component: AlgorithmComponent;
  let fixture: ComponentFixture<AlgorithmComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [AlgorithmComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(AlgorithmComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
