import { Component, Input, OnInit, OnDestroy } from '@angular/core';
import { FormGroup } from '@angular/forms';
import { NamespaceService } from 'kubeflow';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-form-metadata',
  templateUrl: './metadata.component.html',
  styleUrls: ['./metadata.component.scss'],
})
export class FormMetadataComponent implements OnInit, OnDestroy {
  @Input()
  metadataForm: FormGroup;
  sub: Subscription;

  constructor(private ns: NamespaceService) {}

  ngOnInit() {
    this.metadataForm.get('namespace').disable();

    this.sub = this.ns.getSelectedNamespace().subscribe(namespace => {
      this.metadataForm.get('namespace').enable();
      this.metadataForm.get('namespace').setValue(namespace);
      this.metadataForm.get('namespace').disable();
    });
  }

  ngOnDestroy() {
    this.sub.unsubscribe();
  }
}
