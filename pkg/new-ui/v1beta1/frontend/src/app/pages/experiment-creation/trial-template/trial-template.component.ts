import { Component, OnInit, Input, OnDestroy } from '@angular/core';
import { FormGroup, FormArray, FormControl, Validators } from '@angular/forms';
import { KWABackendService } from 'src/app/services/backend.service';
import {
  TrialTemplateResponse,
  ConfigMapResponse,
  ConfigMapBody,
} from 'src/app/models/trial-templates.model';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-form-trial-template',
  templateUrl: './trial-template.component.html',
  styleUrls: ['./trial-template.component.scss'],
})
export class FormTrialTemplateComponent implements OnInit, OnDestroy {
  public templates: ConfigMapResponse[] = [];
  public configmaps: ConfigMapBody[] = [];
  public paths: string[] = [];
  public trialParameters = new FormArray([]);
  private selectedConfigMap: ConfigMapBody;
  private subs = new Subscription();
  private yamlPrv = '';

  @Input() formGroup: FormGroup;

  get yaml(): string {
    return this.yamlPrv;
  }
  set yaml(str: string) {
    this.formGroup.get('yaml').setValue(str);
    this.yamlPrv = str;
  }

  constructor(private backend: KWABackendService) {}

  ngOnInit() {
    this.subs.add(
      this.formGroup.get('type').valueChanges.subscribe(tp => {
        if (tp === 'yaml') {
          this.yaml = '';
          return;
        }

        if (this.templates.length && this.templates.length) {
          this.formGroup.get('cmNamespace').setValue('kubeflow');
        }
      }),
    );

    this.subs.add(
      this.backend.getTrialTemplates('').subscribe(templates => {
        this.templates = templates.Data;
        this.formGroup.get('cmNamespace').setValue('kubeflow');

        // Use the ConfigMap option if the TrialTemplates were successfully
        // fetched
        if (this.templates && this.templates.length) {
          this.formGroup.get('type').setValue('configmap');
        }
      }),
    );

    this.subs.add(
      this.formGroup.get('cmNamespace').valueChanges.subscribe(ns => {
        const ts = this.templates.filter(t => t.ConfigMapNamespace === ns);
        this.configmaps = ts.map(t => t.ConfigMaps)[0] || [];

        if (this.configmaps.length > 0) {
          this.formGroup
            .get('cmName')
            .setValue(this.configmaps[0].ConfigMapName);
        }
      }),
    );

    this.subs.add(
      this.formGroup.get('cmName').valueChanges.subscribe(nm => {
        const cm = this.configmaps.filter(c => c.ConfigMapName === nm)[0];
        this.paths = cm.Templates.map(t => t.Path);
        this.selectedConfigMap = cm;

        if (this.paths.length > 0) {
          this.formGroup.get('cmTrialPath').setValue(this.paths[0]);
        }
      }),
    );

    this.subs.add(
      this.formGroup.get('cmTrialPath').valueChanges.subscribe(path => {
        const t = this.selectedConfigMap.Templates.filter(
          tpl => tpl.Path === path,
        )[0];

        this.formGroup.get('yaml').setValue(t.Yaml);
        this.yaml = t.Yaml;
      }),
    );
  }

  private getTrialParameters(yaml: string) {
    const params = yaml.match(/\${trialParameters.*}/g);
    if (params === null) {
      return [];
    }

    const parsedParams = [];
    for (const param of params) {
      let parsedParam = param;
      parsedParam = parsedParam.replace('${trialParameters.', '');
      parsedParam = parsedParam.substring(0, parsedParam.length - 1);
      parsedParams.push(parsedParam);
    }

    return parsedParams;
  }

  recalculateTrialParameters(yaml: string) {
    const params = this.getTrialParameters(yaml);

    const arrayCtrl = this.formGroup.get('trialParameters') as FormArray;
    arrayCtrl.clear();
    for (const param of params) {
      arrayCtrl.push(
        new FormGroup({
          name: new FormControl(param, Validators.required),
          reference: new FormControl('', []),
          description: new FormControl('', []),
        }),
      );
    }

    this.trialParameters = arrayCtrl;
  }

  ngOnDestroy() {
    this.subs.unsubscribe();
  }
}
