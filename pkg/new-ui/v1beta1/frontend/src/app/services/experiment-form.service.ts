import {
  FormBuilder,
  Validators,
  FormGroup,
  FormArray,
  FormControl,
} from '@angular/forms';
import { Injectable } from '@angular/core';
import { ObjectiveTypeEnum } from '../enumerations/objective-type.enum';
import { AlgorithmsEnum } from '../enumerations/algorithms.enum';
import { BehaviorSubject } from 'rxjs';
import { createParameterGroup, createNasOperationGroup } from '../shared/utils';
import { CollectorKind } from '../pages/experiment-creation/metrics-collector/types';
import { K8sObject, SnackBarService, SnackType } from 'kubeflow';
import { dump, load } from 'js-yaml';
import { ObjectiveSpec, AlgorithmSpec } from '../models/experiment.k8s.model';

@Injectable()
export class ExperimentFormService {
  constructor(private builder: FormBuilder, private snack: SnackBarService) {}

  /*
   * functions for creating the Form Controls
   */
  createMetadataForm(namespace) {
    return this.builder.group({
      name: 'random-experiment',
      namespace,
    });
  }

  createTrialThresholdForm() {
    return this.builder.group({
      parallelTrialCount: 3,
      maxTrialCount: 12,
      maxFailedTrialCount: 3,
    });
  }

  createObjectiveForm() {
    return this.builder.group({
      type: 'maximize',
      goal: 0.99,
      metricName: 'accuracy',
      metricStrategy: 'max',
      additionalMetricNames: this.builder.array(['train-accuracy']),
      metricStrategies: this.builder.array([]),
      setStrategies: this.builder.control(false),
    });
  }

  createAlgorithmObjectiveForm() {
    return this.builder.group({
      type: 'hp',
      algorithm: AlgorithmsEnum.RANDOM,
      algorithmSettings: this.builder.array([]),
    });
  }

  createHyperParametersForm() {
    return this.builder.array([
      createParameterGroup({
        name: 'lr',
        type: 'double',
        value: {
          min: '0.01',
          max: '0.03',
          step: '0.01',
        },
      }),
      createParameterGroup({
        name: 'num-layers',
        type: 'int',
        value: {
          min: '1',
          max: '64',
          step: '1',
        },
      }),
      createParameterGroup({
        name: 'optimizer',
        type: 'categorical',
        value: ['sgd', 'adams', 'ftrl'],
      }),
    ]);
  }

  createNasGraphForm() {
    return this.builder.group({
      layers: 8,
      inputSizes: this.builder.array([32, 32, 3]),
      outputSizes: this.builder.array([10]),
    });
  }

  createNasOperationsForm() {
    return this.builder.array([
      createNasOperationGroup({
        type: 'convolution',
        params: [
          {
            name: 'filter_size',
            type: 'categorical',
            value: [3, 5, 7],
          },
          {
            name: 'num_filter',
            type: 'categorical',
            value: [32, 48, 64],
          },
          {
            name: 'stride',
            type: 'categorical',
            value: [1, 2],
          },
        ],
      }),

      createNasOperationGroup({
        type: 'separable_convolution',
        params: [
          {
            name: 'filter_size',
            type: 'categorical',
            value: [3, 5, 7],
          },
          {
            name: 'num_filter',
            type: 'categorical',
            value: [32, 48, 64],
          },
          {
            name: 'stride',
            type: 'categorical',
            value: [1, 2],
          },
          {
            name: 'depth_multiplier',
            type: 'categorical',
            value: [1, 2],
          },
        ],
      }),

      createNasOperationGroup({
        type: 'depthwise_convolution',
        params: [
          {
            name: 'filter_size',
            type: 'categorical',
            value: [3, 5, 7],
          },
          {
            name: 'num_filter',
            type: 'categorical',
            value: [32, 48, 64],
          },
          {
            name: 'depth_multiplier',
            type: 'categorical',
            value: [1, 2],
          },
        ],
      }),

      createNasOperationGroup({
        type: 'reduction',
        params: [
          {
            name: 'reduction_type',
            type: 'categorical',
            value: ['max_pooling', 'avg_pooling'],
          },
          {
            name: 'pool_size',
            type: 'int',
            value: {
              min: '2',
              max: '3',
              step: '1',
            },
          },
        ],
      }),
    ]);
  }

  createMetricsForm() {
    return this.builder.group({
      kind: CollectorKind.STDOUT,
      metricsFile: '/var/log/katib/metrics.log',
      tfDir: '/var/log/katib/tfevent/',
      prometheus: this.builder.group({
        port: this.builder.control('8080', Validators.required),
        path: this.builder.control('/metrics', Validators.required),
        scheme: this.builder.control('HTTP', Validators.required),
        host: this.builder.control(''),
        httpHeaders: this.builder.array([]),
      }),
    });
  }

  createTrialTemplateForm() {
    return this.builder.group({
      type: 'configmap',
      podLabels: this.builder.array([]),
      containerName: 'training-container',
      successCond: 'status.conditions.#(type=="Complete")#|#(status=="True")#',
      failureCond: 'status.conditions.#(type=="Failed")#|#(status=="True")#',
      retain: 'false',
      cmNamespace: '',
      cmName: '',
      cmTrialPath: '',
      yaml: '',
    });
  }

  /*
   * YAML helpers for moving between FormControls and actual YAML strings
   */
  createYamlTemplateForm(): FormControl {
    const k8sObj = {
      apiVersion: 'kubeflow.org/v1beta1',
      kind: 'Experiment',
      metadata: {
        name: '',
        namespace: '',
      },
      spec: {},
      status: {},
    };

    return new FormControl(k8sObj, []);
  }

  metadataFromCtrl(group: FormGroup): any {
    return {
      name: group.get('name').value,
      namespace: group.get('namespace').value,
    };
  }

  objectiveFromCtrl(group: FormGroup): ObjectiveSpec {
    const objMetric = group.get('metricName').value;
    const objStrategy = group.get('type').value.substring(0, 3);

    const objective: any = {
      type: group.get('type').value,
      goal: group.get('goal').value,
      objectiveMetricName: objMetric,
      // metricStrategies: [{ name: objMetric, value: objStrategy }],
      additionalMetricNames: [],
    };

    const additionalMetrics = group.get('additionalMetricNames').value || [];
    additionalMetrics.forEach(metric => {
      if (!metric) {
        return;
      }

      objective.additionalMetricNames.push(metric);
    });

    if (!group.get('setStrategies').value) {
      return objective;
    }

    objective.metricStrategies = [{ name: objMetric, value: objStrategy }];
    group.get('metricStrategies').value.forEach(metricStrategy => {
      objective.metricStrategies.push({
        name: metricStrategy.metric,
        value: metricStrategy.strategy,
      });
    });

    return objective;
  }

  algorithmFromCtrl(group: FormGroup): AlgorithmSpec {
    return {
      algorithmName: group.get('algorithm').value,
      algorithmSettings: group.get('algorithmSettings').value.map(setting => {
        return { name: setting.name, value: `${setting.value}` };
      }),
    };
  }

  hyperParamsFromCtrl(params: FormArray): any {
    return params.controls.map(paramCtrl => {
      let feasibleSpace: any = {};
      const type = paramCtrl.get('type').value;

      if (type === 'int' || type === 'double') {
        feasibleSpace = {
          min: paramCtrl.get('value').value.min,
          max: paramCtrl.get('value').value.max,
          step: paramCtrl.get('value').value.step,
        };

        if (feasibleSpace.step === '') {
          delete feasibleSpace.step;
        }
      } else {
        feasibleSpace = {
          list: paramCtrl.get('value').value,
        };
      }

      return {
        name: paramCtrl.get('name').value,
        parameterType: paramCtrl.get('type').value,
        feasibleSpace,
      };
    });
  }

  metricsCollectorFromCtrl(group: FormGroup): any {
    const kind = group.get('kind').value;
    const metrics: any = {
      source: {
        fileSystemPath: {
          path: '',
          kind: '',
        },
      },
      collector: { kind },
    };

    if (kind === 'StdOut' || kind === 'None') {
      delete metrics.source;
      return metrics;
    }

    if (kind === 'File') {
      metrics.source.fileSystemPath.path = group.get('metricsFile').value;
      metrics.source.fileSystemPath.kind = 'File';
      return metrics;
    }

    if (kind === 'TensorFlowEvent') {
      metrics.source.fileSystemPath.path = group.get('tfDir').value;
      metrics.source.fileSystemPath.kind = 'Directory';
      return metrics;
    }

    if (kind === 'PrometheusMetric') {
      delete metrics.source.fileSystemPath;
      metrics.source.httpGet = group.get('prometheus').value;

      if (!metrics.source.httpGet.host) {
        delete metrics.source.httpGet.host;
      }

      const headers = metrics.source.httpGet.httpHeaders;
      metrics.source.httpGet.httpHeaders = headers.map(header => {
        return { name: header.key, value: header.value };
      });

      return metrics;
    }

    /* TODO(kimwnasptd): We need to handle the Custom case */
    if (kind === 'Custom') {
      return metrics;
    }

    return {};
  }

  trialTemplateFromCtrl(group: FormGroup): any {
    const trialTemplate: any = {};
    const formValue = group.value;

    trialTemplate.primaryContainerName = formValue.containerName;
    trialTemplate.successCondition = formValue.successCond;
    trialTemplate.failureCondition = formValue.failureCond;
    trialTemplate.retain = formValue.retain === 'true' ? true : false;

    if (formValue.podLabels && formValue.podLabels.length) {
      trialTemplate.primaryPodLabels = {};
      formValue.podLabels.map(
        label => (trialTemplate.primaryPodLabels[label.key] = label.value),
      );
    }

    if (formValue.type === 'yaml') {
      try {
        trialTemplate.trialSpec = load(formValue.yaml);
      } catch (e) {
        this.snack.open(`${e.reason}`, SnackType.Warning, 4000);
      }

      return trialTemplate;
    }

    trialTemplate.configMap = {
      configMapName: formValue.cmName,
      configMapNamespace: formValue.cmNamespace,
      templatePath: formValue.cmTrialPath,
    };

    return trialTemplate;
  }
}
