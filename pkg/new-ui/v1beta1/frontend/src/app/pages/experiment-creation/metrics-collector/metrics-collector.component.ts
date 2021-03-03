import { Component, OnInit, Input } from '@angular/core';
import { FormGroup } from '@angular/forms';
import { CollectorKind } from 'src/app/enumerations/metrics-collector';

@Component({
  selector: 'app-form-metrics-collector',
  templateUrl: './metrics-collector.component.html',
  styleUrls: ['./metrics-collector.component.scss'],
})
export class FormMetricsCollectorComponent implements OnInit {
  @Input() formGroup: FormGroup;
  kind = CollectorKind;
  customYaml =
    'name: metrics-collector\nimage: <collector-image>\nresources: {}';

  constructor() {}

  ngOnInit() {}
}
