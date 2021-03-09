import { Component, OnInit, Input } from '@angular/core';
import { FormGroup } from '@angular/forms';

@Component({
  selector: 'app-algorithm-setting',
  templateUrl: './setting.component.html',
  styleUrls: ['./setting.component.scss'],
})
export class FormAlgorithmSettingComponent implements OnInit {
  @Input() setting: FormGroup;

  constructor() {}

  ngOnInit() {}
}
