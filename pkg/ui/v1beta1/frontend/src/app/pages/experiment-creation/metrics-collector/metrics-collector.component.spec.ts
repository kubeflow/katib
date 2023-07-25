import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { FormControl, FormGroup } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { EditorModule, FormModule } from 'kubeflow';
import { ListKeyValueModule } from 'src/app/shared/list-key-value/list-key-value.module';

import { FormMetricsCollectorComponent } from './metrics-collector.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

describe('FormMetricsCollectorComponent', () => {
  let component: FormMetricsCollectorComponent;
  let fixture: ComponentFixture<FormMetricsCollectorComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          CommonModule,
          BrowserAnimationsModule,
          FormModule,
          ListKeyValueModule,
          EditorModule,
        ],
        declarations: [FormMetricsCollectorComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(FormMetricsCollectorComponent);
    component = fixture.componentInstance;
    component.formGroup = new FormGroup({
      kind: new FormControl(),
      metricsFile: new FormControl(),
      tfDir: new FormControl(),
      port: new FormControl(),
      path: new FormControl(),
      scheme: new FormControl(),
      host: new FormControl(),
      prometheus: new FormControl(),
      customYaml: new FormControl(),
    });
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
