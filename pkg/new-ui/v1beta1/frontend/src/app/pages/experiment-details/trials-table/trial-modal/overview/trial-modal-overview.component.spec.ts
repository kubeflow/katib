import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TrialModalOverviewComponent } from './trial-modal-overview.component';

describe('TrialModalOverviewComponent', () => {
  let component: TrialModalOverviewComponent;
  let fixture: ComponentFixture<TrialModalOverviewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [TrialModalOverviewComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TrialModalOverviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
