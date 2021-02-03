import { Component, Input, Output, EventEmitter } from '@angular/core';
import { FormGroup } from '@angular/forms';
import { EventsV1beta1Api } from '@kubernetes/client-node';

@Component({
  selector: 'app-nas-operation',
  templateUrl: './operation.component.html',
  styleUrls: ['./operation.component.scss'],
})
export class OperationComponent {
  @Input() formGroup: FormGroup;
  @Output() removeCtrl = new EventEmitter<boolean>();
}
