import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TrialParameterComponent } from './trial-parameter.component';

describe('TrialParameterComponent', () => {
  let component: TrialParameterComponent;
  let fixture: ComponentFixture<TrialParameterComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [TrialParameterComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TrialParameterComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
