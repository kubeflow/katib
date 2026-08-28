import { TestBed } from '@angular/core/testing';
import { ThemeService } from './theme.service';

describe('ThemeService', () => {
  let service: ThemeService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ThemeService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should toggle theme and update localStorage', () => {
    const initial = localStorage.getItem('katib-theme');
    service.toggle();
    const updated = localStorage.getItem('katib-theme');
    expect(updated).not.toBe(initial);
  });
});
