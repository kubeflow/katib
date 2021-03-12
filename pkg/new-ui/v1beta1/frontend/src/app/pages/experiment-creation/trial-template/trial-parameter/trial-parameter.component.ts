import { Component, OnInit, Input } from '@angular/core';
import { FormGroup } from '@angular/forms';

@Component({
  selector: 'app-trial-parameter',
  templateUrl: './trial-parameter.component.html',
  styleUrls: ['./trial-parameter.component.scss'],
})
export class TrialParameterComponent implements OnInit {
  @Input() formGroup: FormGroup;

  constructor() {}

  ngOnInit() {}
}
