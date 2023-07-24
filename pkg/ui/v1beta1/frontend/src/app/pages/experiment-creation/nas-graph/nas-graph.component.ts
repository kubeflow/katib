import { Component, Input } from '@angular/core';
import { FormGroup } from '@angular/forms';

@Component({
  selector: 'app-form-nas-graph',
  templateUrl: './nas-graph.component.html',
  styleUrls: ['./nas-graph.component.scss'],
})
export class FormNasGraphComponent {
  @Input() formGroup: FormGroup;
}
