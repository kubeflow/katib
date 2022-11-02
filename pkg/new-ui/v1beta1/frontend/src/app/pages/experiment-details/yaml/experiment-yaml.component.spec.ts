import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { CommonModule } from '@angular/common';
import { AceEditorModule } from 'ng2-ace-editor';

import { ExperimentYamlComponent } from './experiment-yaml.component';

describe('ExperimentYamlComponent', () => {
  let component: ExperimentYamlComponent;
  let fixture: ComponentFixture<ExperimentYamlComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [CommonModule, BrowserAnimationsModule, AceEditorModule],
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
