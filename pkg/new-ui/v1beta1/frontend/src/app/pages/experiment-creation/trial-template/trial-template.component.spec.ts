import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { TrialTemplateComponent } from './trial-template.component';

describe('TrialTemplateComponent', () => {
  let component: TrialTemplateComponent;
  let fixture: ComponentFixture<TrialTemplateComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [TrialTemplateComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(TrialTemplateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
