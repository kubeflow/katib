import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { FormArray, FormControl, FormGroup } from '@angular/forms';
import { of } from 'rxjs';
import { KWABackendService } from 'src/app/services/backend.service';
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatDividerModule } from '@angular/material/divider';
import { EditorModule, FormModule, PopoverModule } from 'kubeflow';
import { ListKeyValueModule } from 'src/app/shared/list-key-value/list-key-value.module';

import { FormTrialTemplateComponent } from './trial-template.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

describe('FormTrialTemplateComponent', () => {
  let component: FormTrialTemplateComponent;
  let fixture: ComponentFixture<FormTrialTemplateComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          CommonModule,
          BrowserAnimationsModule,
          FormModule,
          ListKeyValueModule,
          MatDividerModule,
          PopoverModule,
          EditorModule,
        ],
        declarations: [FormTrialTemplateComponent],
        providers: [
          {
            provide: KWABackendService,
            useValue: { getTrialTemplates: () => of() },
          },
        ],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(FormTrialTemplateComponent);
    component = fixture.componentInstance;
    component.formGroup = new FormGroup({
      trialParameters: new FormArray([]),
      podLabels: new FormControl(),
      containerName: new FormControl(),
      successCond: new FormControl(),
      failureCond: new FormControl(),
      retain: new FormControl(),
      type: new FormControl(),
      cmNamespace: new FormControl(),
      cmName: new FormControl(),
      cmTrialPath: new FormControl(),
    });
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
