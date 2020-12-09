import React from 'react';
import { connect } from 'react-redux';

import jsyaml from 'js-yaml';

import { withStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import Typography from '@material-ui/core/Typography';

import CommonParametersMeta from './Params/CommonMeta';
import CommonParametersSpec from './Params/CommonSpec';
import Objective from './Params/Objective';
import Algorithm from './Params/Algorithm';
import TrialTemplate from '../../Common/Create/Params/Trial/TrialTemplate';
import NASConfig from './Params/NASConfig';
import MetricsCollectorSpec from '../../Common/Create/Params/MetricsCollector';

import { submitNASJob } from '../../../actions/nasCreateActions';

import { validationError } from '../../../actions/generalActions';
import * as constants from '../../../constants/constants';

const styles = theme => ({
  root: {
    width: '90%',
    margin: '0 auto',
  },
  submit: {
    textAlign: 'center',
    margin: 20,
  },
});

const SectionInTypography = name => {
  return (
    <div>
      <Grid container>
        <Grid item xs={12} sm={12}>
          <Typography variant="h6">{name}</Typography>
          <hr />
        </Grid>
      </Grid>
    </div>
  );
};

// probably get render into a function
const deCapitalizeFirstLetterAndAppend = (source, destination) => {
  source.map((parameter, i) => {
    let value = NaN;
    // Try to get Number from parameter value if it is not empty
    if (parameter.value !== '') {
      value = Number(parameter.value);
    }
    let name = parameter.name.charAt(0).toLowerCase() + parameter.name.slice(1);
    return (destination[name] = isNaN(value) ? parameter.value : value);
  });
};

const addAlgorithmSettings = (spec, destination) => {
  spec.map((parameter, i) => {
    return destination.push(parameter);
  });
};

const addOperations = (source, destination) => {
  source.map((operation, index) => {
    let parameters = [];
    operation.parameters.map((param, i) => {
      let tempParam = {};
      tempParam.name = param.name;
      tempParam.parametertype = param.parameterType;
      tempParam.feasibleSpace = {};
      if (param.feasibleSpace === 'list') {
        tempParam.feasibleSpace.list = param.list.map((param, i) => param.value);
      } else {
        tempParam.feasibleSpace.min = param.min;
        tempParam.feasibleSpace.max = param.max;
        tempParam.feasibleSpace.step = param.step;
      }
      return parameters.push(tempParam);
    });
    return destination.push({
      operationType: operation.operationType,
      parameters: parameters,
    });
  });
};

const NASParameters = props => {
  const submitNASJob = () => {
    let data = {};

    // Add metadata.
    data.metadata = {};
    deCapitalizeFirstLetterAndAppend(props.commonParametersMetadata, data.metadata);

    // Add common parameters.
    data.spec = {};
    deCapitalizeFirstLetterAndAppend(props.commonParametersSpec, data.spec);

    // Add objective.
    data.spec.objective = {};
    deCapitalizeFirstLetterAndAppend(props.objective, data.spec.objective);

    // Add additional metrics.
    data.spec.objective.additionalMetricNames = props.additionalMetricNames;

    // Add metric strategies.
    data.spec.objective.metricStrategies = props.metricStrategies.map(metric => ({
      name: metric.name,
      value: metric.strategy,
    }));

    // Add algorithm.
    data.spec.algorithm = {};
    data.spec.algorithm.algorithmName = props.algorithmName;
    data.spec.algorithm.algorithmSettings = [];
    addAlgorithmSettings(props.algorithmSettings, data.spec.algorithm.algorithmSettings);

    data.spec.nasConfig = {};
    data.spec.nasConfig.graphConfig = {
      numLayers: Number(props.numLayers),
      inputSizes: props.inputSize.map(size => Number(size)),
      outputSizes: props.outputSize.map(size => Number(size)),
    };
    data.spec.nasConfig.operations = [];
    addOperations(props.operations, data.spec.nasConfig.operations);

    // Metrics Collector
    let newMCSpec = JSON.parse(JSON.stringify(props.mcSpec));

    // Delete empty metrics format
    if (
      newMCSpec.source.filter.metricsFormat.length === 0 ||
      newMCSpec.collector.kind === constants.MC_KIND_NONE
    ) {
      delete newMCSpec.source.filter;
    }

    if (
      newMCSpec.collector.kind === constants.MC_KIND_STDOUT ||
      newMCSpec.collector.kind === constants.MC_KIND_NONE
    ) {
      // Delete fileSystemPath and httpGet
      delete newMCSpec.source.fileSystemPath;
      delete newMCSpec.source.httpGet;
    }

    if (
      newMCSpec.collector.kind === constants.MC_KIND_FILE ||
      newMCSpec.collector.kind === constants.MC_KIND_TENSORFLOW_EVENT ||
      newMCSpec.collector.kind === constants.MC_KIND_CUSTOM
    ) {
      // Delete httpGet
      delete newMCSpec.source.httpGet;
      // Delete empty fileSystemPath
      if (newMCSpec.source.fileSystemPath.kind === constants.MC_FILE_SYSTEM_NO_KIND) {
        delete newMCSpec.source.fileSystemPath;
      }
    }

    if (newMCSpec.collector.kind === constants.MC_KIND_PROMETHEUS) {
      // Delete file System Path
      delete newMCSpec.source.fileSystemPath;
      // Delete empty host
      if (newMCSpec.source.httpGet.host === '') {
        delete newMCSpec.source.httpGet.host;
      }
      // Delete empty headers
      if (newMCSpec.source.httpGet.httpHeaders.length === 0) {
        delete newMCSpec.source.httpGet.httpHeaders;
      }
    }

    // Delete empty source
    if (newMCSpec.source !== undefined && Object.keys(newMCSpec.source).length === 0) {
      delete newMCSpec.source;
    }

    // Add Custom Container YAML to the Metrics Collector
    if (
      newMCSpec.collector.kind === constants.MC_KIND_CUSTOM &&
      props.mcCustomContainerYaml !== ''
    ) {
      try {
        let mcCustomContainerJson = jsyaml.load(props.mcCustomContainerYaml);
        newMCSpec.collector.customCollector = mcCustomContainerJson;
      } catch {
        props.validationError('Metrics Collector Custom Container is not valid YAML!');
        return;
      }
    }

    data.spec.metricsCollectorSpec = newMCSpec;

    // Add Trial template.
    // Add Trial specification.
    data.spec.trialTemplate = {};
    deCapitalizeFirstLetterAndAppend(props.trialTemplateSpec, data.spec.trialTemplate);
    if (data.spec.trialTemplate.retain === 'true') {
      data.spec.trialTemplate.retain = true;
    } else if (data.spec.trialTemplate.retain === 'false') {
      data.spec.trialTemplate.retain = false;
    } else {
      props.validationError('Trial template retain parameter must be true or false!');
      return;
    }

    // Remove empty items from PrimaryPodLabels array.
    let filteredPrimaryLabels = props.primaryPodLabels.filter(function (label) {
      return label.key.trim() !== '' && label.value.trim() !== '';
    });

    // If array is not empty add PrimaryPodLabels.
    if (filteredPrimaryLabels.length > 0) {
      data.spec.trialTemplate.primaryPodLabels = {};
      filteredPrimaryLabels.forEach(
        label => (data.spec.trialTemplate.primaryPodLabels[label.key] = label.value),
      );
    }

    // Add Trial Source.
    if (
      props.trialTemplateSource === constants.TEMPLATE_SOURCE_YAML &&
      props.trialTemplateYAML !== ''
    ) {
      // Try to parse template YAML to JSON.
      try {
        let trialTemplateJSON = jsyaml.load(props.trialTemplateYAML);
        data.spec.trialTemplate.trialSpec = trialTemplateJSON;
      } catch {
        props.validationError('Trial Template is not valid YAML!');
        return;
      }
      // Otherwise assign ConfigMap.
    } else {
      data.spec.trialTemplate.configMap = {
        configMapNamespace: props.templateConfigMapNamespace,
        configMapName: props.templateConfigMapName,
        templatePath: props.templateConfigMapPath,
      };
    }

    // Add Trial parameters if it is not empty.
    if (props.trialParameters.length > 0) {
      data.spec.trialTemplate.trialParameters = props.trialParameters;
    }

    props.submitNASJob(data);
  };

  const { classes } = props;

  return (
    <div className={classes.root}>
      {/* Common Metadata */}
      {SectionInTypography('Metadata')}
      <CommonParametersMeta />
      {SectionInTypography('Common Parameters')}
      <CommonParametersSpec />
      {SectionInTypography('Objective')}
      <Objective />
      {SectionInTypography('Algorithm')}
      <Algorithm />
      {SectionInTypography('NAS Config')}
      <NASConfig />
      {SectionInTypography('Metrics Collector Spec')}
      <MetricsCollectorSpec jobType={constants.EXPERIMENT_TYPE_NAS} />
      {SectionInTypography('Trial Template Spec')}
      <TrialTemplate />
      <div className={classes.submit}>
        <Button variant="contained" color={'primary'} onClick={submitNASJob}>
          Deploy
        </Button>
      </div>
    </div>
  );
};

// TODO: think of a better way of passing those
const mapStateToProps = state => {
  let templatesData = state[constants.GENERAL_MODULE].trialTemplatesData;
  let templateCMNamespace = '';
  let templateCMName = '';
  let templateCMPath = '';

  if (state[constants.GENERAL_MODULE].configMapNamespaceIndex !== -1) {
    let nsData = templatesData[state[constants.GENERAL_MODULE].configMapNamespaceIndex];
    let nameData = nsData.ConfigMaps[state[constants.GENERAL_MODULE].configMapNameIndex];
    let pathData = nameData.Templates[state[constants.GENERAL_MODULE].configMapPathIndex];

    templateCMNamespace = nsData.ConfigMapNamespace;
    templateCMName = nameData.ConfigMapName;
    templateCMPath = pathData.Path;
  }
  return {
    commonParametersMetadata: state[constants.NAS_CREATE_MODULE].commonParametersMetadata,
    commonParametersSpec: state[constants.NAS_CREATE_MODULE].commonParametersSpec,
    objective: state[constants.NAS_CREATE_MODULE].objective,
    additionalMetricNames: state[constants.NAS_CREATE_MODULE].additionalMetricNames,
    algorithmName: state[constants.NAS_CREATE_MODULE].algorithmName,
    algorithmSettings: state[constants.NAS_CREATE_MODULE].algorithmSettings,
    metricStrategies: state[constants.NAS_CREATE_MODULE].metricStrategies,
    numLayers: state[constants.NAS_CREATE_MODULE].numLayers,
    inputSize: state[constants.NAS_CREATE_MODULE].inputSize,
    outputSize: state[constants.NAS_CREATE_MODULE].outputSize,
    operations: state[constants.NAS_CREATE_MODULE].operations,
    primaryPodLabels: state[constants.GENERAL_MODULE].primaryPodLabels,
    trialTemplateSpec: state[constants.GENERAL_MODULE].trialTemplateSpec,
    trialTemplateSource: state[constants.GENERAL_MODULE].trialTemplateSource,
    templateConfigMapNamespace: templateCMNamespace,
    templateConfigMapName: templateCMName,
    templateConfigMapPath: templateCMPath,
    trialTemplateYAML: state[constants.GENERAL_MODULE].trialTemplateYAML,
    trialParameters: state[constants.GENERAL_MODULE].trialParameters,
    mcSpec: state[constants.NAS_CREATE_MODULE].mcSpec,
    mcCustomContainerYaml: state[constants.NAS_CREATE_MODULE].mcCustomContainerYaml,
  };
};

export default connect(mapStateToProps, { submitNASJob, validationError })(
  withStyles(styles)(NASParameters),
);
