import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KWABackendService } from './backend.service';
import { ReactiveFormsModule } from '@angular/forms';
import { MatSnackBarModule } from '@angular/material/snack-bar';

describe('KWABackendService', () => {
  beforeEach(() =>
    TestBed.configureTestingModule({
      imports: [
        HttpClientTestingModule,
        ReactiveFormsModule,
        MatSnackBarModule,
      ],
    }),
  );

  it('should be created', () => {
    const service: KWABackendService = TestBed.inject(KWABackendService);
    expect(service).toBeTruthy();
  });
});
