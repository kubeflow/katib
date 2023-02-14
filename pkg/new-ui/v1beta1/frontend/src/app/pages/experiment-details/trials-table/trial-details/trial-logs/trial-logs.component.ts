import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-trial-logs',
  templateUrl: './trial-logs.component.html',
  styleUrls: ['./trial-logs.component.scss'],
})
export class TrialLogsComponent {
  public logs: string[];

  @Input() logsRequestError: string;

  @Input()
  set trialLogs(trialLogs: string) {
    if (!trialLogs) {
      return;
    }

    this.logs = trialLogs.split('\n');
  }
}
