import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class ThemeService {
  private isDark = new BehaviorSubject<boolean>(this.getSavedTheme());
  isDark$ = this.isDark.asObservable();

  private getSavedTheme(): boolean {
    const saved = localStorage.getItem('katib-theme');
    if (saved) return saved === 'dark';
    return window.matchMedia('(prefers-color-scheme: dark)').matches;
  }

  toggle(): void {
    const newValue = !this.isDark.value;
    this.isDark.next(newValue);
    localStorage.setItem('katib-theme', newValue ? 'dark' : 'light');
    document.documentElement.setAttribute('data-theme', newValue ? 'dark' : 'light');
    document.body.classList.toggle('dark-theme', newValue);
    document.body.classList.toggle('light-theme', !newValue);
  }

  initTheme(): void {
    const dark = this.getSavedTheme();
    document.documentElement.setAttribute('data-theme', dark ? 'dark' : 'light');
    document.body.classList.toggle('dark-theme', dark);
    document.body.classList.toggle('light-theme', !dark);
  }
}
