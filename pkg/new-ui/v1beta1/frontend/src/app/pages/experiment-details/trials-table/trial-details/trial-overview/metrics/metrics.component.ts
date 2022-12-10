import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-metrics-overview',
  templateUrl: './metrics.component.html',
})
export class TrialMetricsComponent {
  @Input() name: string;
  @Input() latest: string;
  @Input() max: string;
  @Input() min: string;
}
