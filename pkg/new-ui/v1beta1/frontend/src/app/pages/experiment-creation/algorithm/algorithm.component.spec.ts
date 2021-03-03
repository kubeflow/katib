import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AlgorithmComponent } from './algorithm.component';

describe('AlgorithmComponent', () => {
  let component: AlgorithmComponent;
  let fixture: ComponentFixture<AlgorithmComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [AlgorithmComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AlgorithmComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
