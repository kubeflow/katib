import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HyperParametersComponent } from './hyper-parameters.component';

describe('HyperParametersComponent', () => {
  let component: HyperParametersComponent;
  let fixture: ComponentFixture<HyperParametersComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [HyperParametersComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HyperParametersComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
