import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { CommonModule } from '@angular/common';

import { ExperimentYamlComponent } from './experiment-yaml.component';
import { EditorModule } from 'kubeflow';

describe('ExperimentYamlComponent', () => {
  let component: ExperimentYamlComponent;
  let fixture: ComponentFixture<ExperimentYamlComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [CommonModule, BrowserAnimationsModule, EditorModule],
        declarations: [ExperimentYamlComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(ExperimentYamlComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
