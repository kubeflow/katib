import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { CommonModule } from '@angular/common';
import { MatTableModule } from '@angular/material/table';
import { MatIconModule } from '@angular/material/icon';
import { MatDialogModule } from '@angular/material/dialog';
import { NgxChartsModule } from '@swimlane/ngx-charts';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { RouterTestingModule } from '@angular/router/testing';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatButtonModule } from '@angular/material/button';
import { TrialsTableComponent } from './trials-table.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

describe('TrialsTableComponent', () => {
  let component: TrialsTableComponent;
  let fixture: ComponentFixture<TrialsTableComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          CommonModule,
          MatTableModule,
          MatProgressSpinnerModule,
          MatDialogModule,
          MatIconModule,
          NgxChartsModule,
          MatTooltipModule,
          MatButtonModule,
          RouterTestingModule,
          BrowserAnimationsModule,
        ],
        declarations: [TrialsTableComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(TrialsTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
