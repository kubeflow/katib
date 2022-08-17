import {
  ChangeDetectionStrategy,
  Component,
  EventEmitter,
  Input,
  OnChanges,
  Output,
  SimpleChanges,
} from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';

@Component({
  selector: 'app-trials-table',
  templateUrl: './trials-table.component.html',
  styleUrls: ['./trials-table.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TrialsTableComponent implements OnChanges {
  @Input()
  displayedColumns = [];

  @Input()
  data = [];

  @Input()
  experimentName = [];

  @Input()
  namespace: string;

  @Input()
  bestTrialName: string;

  @Output()
  mouseOnTrial = new EventEmitter<number>();

  @Output()
  leaveMouseFromTrial = new EventEmitter<void>();

  bestTrialIndex: number;

  constructor(public dialog: MatDialog, private router: Router) {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.displayedColumns) {
      this.displayedColumns = this.displayedColumns.slice(
        0,
        this.displayedColumns.length,
      );
    }

    if (this.data.length > 0 && this.bestTrialName) {
      this.bestTrialIndex = this.data.findIndex(
        trial => trial[0] === this.bestTrialName,
      );
    }
  }

  openTrialModal(name: string) {
    this.router.navigate([`/experiment/${this.experimentName}/trial/${name}`]);
  }

  handleMouseLeave = () => this.leaveMouseFromTrial.emit();

  handleMouseOver = event => this.mouseOnTrial.emit(+event[event.length - 1]);

  goToKfpRun(kfpRun: string) {
    // If the trial does not have a kfp run then do not redirect
    if (!kfpRun) {
      return;
    }

    this.router.navigate([]).then(() => {
      window.open(`/pipeline/#/runs/details/${kfpRun}`);
    });
  }
}
