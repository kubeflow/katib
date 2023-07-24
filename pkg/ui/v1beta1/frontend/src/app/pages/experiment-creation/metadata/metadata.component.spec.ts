import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { of } from 'rxjs';
import { NamespaceService, FormModule } from 'kubeflow';

import { FormMetadataComponent } from './metadata.component';
import { FormControl, FormGroup } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';

describe('FormMetadataComponent', () => {
  let component: FormMetadataComponent;
  let fixture: ComponentFixture<FormMetadataComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          CommonModule,
          FormModule,
          BrowserAnimationsModule,
          MatFormFieldModule,
          MatInputModule,
        ],
        declarations: [FormMetadataComponent],
        providers: [
          {
            provide: NamespaceService,
            useValue: { getSelectedNamespace: () => of('') },
          },
        ],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(FormMetadataComponent);
    component = fixture.componentInstance;
    component.metadataForm = new FormGroup({
      name: new FormControl(),
      namespace: new FormControl(),
    });
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
