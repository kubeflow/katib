import { TestBed } from '@angular/core/testing';

import { BackendService } from './backend.service';

describe('BackendService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: BackendService = TestBed.inject(BackendService);
    expect(service).toBeTruthy();
  });
});
