import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { EarlyStoppingComponent } from './early-stopping.component';

describe('EarlyStoppingComponent', () => {
  let component: EarlyStoppingComponent;
  let fixture: ComponentFixture<EarlyStoppingComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [EarlyStoppingComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(EarlyStoppingComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
