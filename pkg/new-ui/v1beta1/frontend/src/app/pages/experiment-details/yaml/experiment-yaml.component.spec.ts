import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ExperimentYamlComponent } from './experiment-yaml.component';

describe('ExperimentYamlComponent', () => {
  let component: ExperimentYamlComponent;
  let fixture: ComponentFixture<ExperimentYamlComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ExperimentYamlComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ExperimentYamlComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
