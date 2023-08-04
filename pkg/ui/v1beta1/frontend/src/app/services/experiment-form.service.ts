import { FormBuilder, Validators, FormGroup, FormArray } from '@angular/forms';
import { Injectable } from '@angular/core';
import {
  AlgorithmsEnum,
  EarlyStoppingAlgorithmsEnum,
} from '../enumerations/algorithms.enum';
import { createParameterGroup, createNasOperationGroup } from '../shared/utils';
import { SnackBarConfig, SnackBarService, SnackType } from 'kubeflow';
import { load } from 'js-yaml';
import {
  ObjectiveSpec,
  AlgorithmSpec,
  ParameterSpec,
  FeasibleSpaceMinMax,
  NasOperation,
  AlgorithmSetting,
} from '../models/experiment.k8s.model';
import { CollectorKind } from '../enumerations/metrics-collector';

@Injectable()
export class ExperimentFormService {
  constructor(private builder: FormBuilder, private snack: SnackBarService) {}

  /*
   * functions for creating the Form Controls
   */
  createMetadataForm(namespace): FormGroup {
    return this.builder.group({
      name: 'random-experiment',
      namespace,
    });
  }

  createTrialThresholdForm(): FormGroup {
    return this.builder.group({
      parallelTrialCount: 3,
      maxTrialCount: 12,
      maxFailedTrialCount: 3,
      resumePolicy: 'Never',
    });
  }

  createObjectiveForm(): FormGroup {
    return this.builder.group({
      type: 'maximize',
      goal: 0.99,
      metricName: 'Validation-accuracy',
      metricStrategy: 'max',
      additionalMetricNames: this.builder.array(['Train-accuracy']),
      metricStrategies: this.builder.array([]),
      setStrategies: this.builder.control(false),
    });
  }

  createAlgorithmObjectiveForm(): FormGroup {
    return this.builder.group({
      type: 'hp',
      algorithm: AlgorithmsEnum.RANDOM,
      algorithmSettings: this.builder.array([]),
    });
  }

  createEarlyStoppingForm(): FormGroup {
    return this.builder.group({
      algorithmName: EarlyStoppingAlgorithmsEnum.NONE,
      algorithmSettings: this.builder.array([]),
    });
  }

  createHyperParametersForm(): FormArray {
    return this.builder.array([
      createParameterGroup({
        name: 'lr',
        parameterType: 'double',
        feasibleSpace: {
          min: '0.01',
          max: '0.03',
          step: '0.01',
        },
      }),
      createParameterGroup({
        name: 'num-layers',
        parameterType: 'int',
        feasibleSpace: {
          min: '1',
          max: '64',
          step: '1',
        },
      }),
      createParameterGroup({
        name: 'optimizer',
        parameterType: 'categorical',
        feasibleSpace: { list: ['sgd', 'adam', 'ftrl'] },
      }),
    ]);
  }

  createNasGraphForm(): FormGroup {
    return this.builder.group({
      layers: 8,
      inputSizes: this.builder.array([32, 32, 3]),
      outputSizes: this.builder.array([10]),
    });
  }

  createNasOperationsForm(): FormArray {
    return this.builder.array([
      createNasOperationGroup({
        operationType: 'convolution',
        parameters: [
          {
            name: 'filter_size',
            parameterType: 'categorical',
            feasibleSpace: { list: ['3', '5', '7'] },
          },
          {
            name: 'num_filter',
            parameterType: 'categorical',
            feasibleSpace: { list: ['32', '48', '64'] },
          },
          {
            name: 'stride',
            parameterType: 'categorical',
            feasibleSpace: { list: ['1', '2'] },
          },
        ],
      }),

      createNasOperationGroup({
        operationType: 'separable_convolution',
        parameters: [
          {
            name: 'filter_size',
            parameterType: 'categorical',
            feasibleSpace: { list: ['3', '5', '7'] },
          },
          {
            name: 'num_filter',
            parameterType: 'categorical',
            feasibleSpace: { list: ['32', '48', '64'] },
          },
          {
            name: 'stride',
            parameterType: 'categorical',
            feasibleSpace: { list: ['1', '2'] },
          },
          {
            name: 'depth_multiplier',
            parameterType: 'categorical',
            feasibleSpace: { list: ['1', '2'] },
          },
        ],
      }),

      createNasOperationGroup({
        operationType: 'depthwise_convolution',
        parameters: [
          {
            name: 'filter_size',
            parameterType: 'categorical',
            feasibleSpace: { list: ['3', '5', '7'] },
          },
          {
            name: 'num_filter',
            parameterType: 'categorical',
            feasibleSpace: { list: ['32', '48', '64'] },
          },
          {
            name: 'depth_multiplier',
            parameterType: 'categorical',
            feasibleSpace: { list: ['1', '2'] },
          },
        ],
      }),

      createNasOperationGroup({
        operationType: 'reduction',
        parameters: [
          {
            name: 'reduction_type',
            parameterType: 'categorical',
            feasibleSpace: { list: ['max_pooling', 'avg_pooling'] },
          },
          {
            name: 'pool_size',
            parameterType: 'int',
            feasibleSpace: {
              min: '2',
              max: '3',
              step: '1',
            },
          },
        ],
      }),
    ]);
  }

  createMetricsForm(): FormGroup {
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
      customYaml:
        'name: metrics-collector\nimage: <collector-image>\nresources: {}',
    });
  }

  createTrialTemplateForm(): FormGroup {
    return this.builder.group({
      type: 'yaml',
      podLabels: this.builder.array([]),
      containerName: 'training-container',
      successCond: 'status.conditions.#(type=="Complete")#|#(status=="True")#',
      failureCond: 'status.conditions.#(type=="Failed")#|#(status=="True")#',
      retain: 'false',
      cmNamespace: '',
      cmName: '',
      cmTrialPath: '',
      yaml: '',
      trialParameters: this.builder.array([]),
    });
  }

  /**
   * helpers for parsing the form controls
   */
  parseParam(parameter: ParameterSpec): ParameterSpec {
    // we should omit the step in case it was not defined
    const param = JSON.parse(JSON.stringify(parameter));

    if (
      param.parameterType === 'discrete' ||
      param.parameterType === 'categorical'
    ) {
      return param;
    }

    const step = (param.feasibleSpace as FeasibleSpaceMinMax).step;
    if (step === '' || step === null) {
      delete (param.feasibleSpace as FeasibleSpaceMinMax).step;
    }

    for (const key in param.feasibleSpace) {
      param.feasibleSpace[key] = param.feasibleSpace[key].toString();
    }

    return param;
  }

  /*
   * YAML helpers for moving between FormControls and actual YAML strings
   */
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
    const settings: AlgorithmSetting[] = [];
    group.get('algorithmSettings').value.forEach(setting => {
      if (setting.value === null) {
        return;
      }

      settings.push({ name: setting.name, value: `${setting.value}` });
    });

    return {
      algorithmName: group.get('algorithm').value,
      algorithmSettings: settings,
    };
  }

  earlyStoppingFromCtrl(group: FormGroup): AlgorithmSpec {
    const settings: AlgorithmSetting[] = [];
    group.get('algorithmSettings').value.forEach(setting => {
      if (setting.value === null) {
        return;
      }

      settings.push({ name: setting.name, value: `${setting.value}` });
    });

    return {
      algorithmName: group.get('algorithmName').value,
      algorithmSettings: settings,
    };
  }

  hyperParamsFromCtrl(paramsArray: FormArray): any {
    const params = paramsArray.value as ParameterSpec[];
    return params.map(param => {
      return this.parseParam(param);
    });
  }

  nasOpsFromCtrl(operations: FormArray): NasOperation[] {
    const ops: NasOperation[] = operations.value;

    return ops.map(op => {
      op.parameters = op.parameters.map(param => this.parseParam(param));
      return op;
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

    if (kind === 'Custom') {
      delete metrics.source.fileSystemPath;
      try {
        metrics.collector.customCollector = load(group.get('customYaml').value);
      } catch (e) {
        const config: SnackBarConfig = {
          data: {
            msg: 'Metrics Colletor(Custom): ' + `${e.reason}`,
            snackType: SnackType.Error,
          },
        };
        this.snack.open(config);
      }
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

    trialTemplate.trialParameters = formValue.trialParameters;

    if (formValue.type === 'yaml') {
      try {
        trialTemplate.trialSpec = load(formValue.yaml);
      } catch (e) {
        const config: SnackBarConfig = {
          data: {
            msg: 'Trial Template: ' + `${e.reason}`,
            snackType: SnackType.Error,
          },
        };
        this.snack.open(config);
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
