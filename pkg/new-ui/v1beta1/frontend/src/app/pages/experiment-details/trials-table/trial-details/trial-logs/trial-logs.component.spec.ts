import { CommonModule } from '@angular/common';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import {
  HeadingSubheadingRowModule,
  KubeflowModule,
  LogsViewerModule,
} from 'kubeflow';

import { TrialLogsComponent } from './trial-logs.component';

describe('TrialLogsComponent', () => {
  let component: TrialLogsComponent;
  let fixture: ComponentFixture<TrialLogsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [TrialLogsComponent],
      imports: [
        CommonModule,
        KubeflowModule,
        HeadingSubheadingRowModule,
        LogsViewerModule,
      ],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(TrialLogsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
