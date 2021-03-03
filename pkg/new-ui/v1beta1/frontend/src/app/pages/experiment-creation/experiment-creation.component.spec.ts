import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ExperimentCreationComponent } from './experiment-creation.component';

describe('ExperimentCreationComponent', () => {
  let component: ExperimentCreationComponent;
  let fixture: ComponentFixture<ExperimentCreationComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ExperimentCreationComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ExperimentCreationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
