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
import TrialTemplate from '../../Common/Create/Params/Trial/TrialTemplate';
import Parameters from './Params/Parameters';
import Algorithm from './Params/Algorithm';
import EarlyStopping from '../../Common/Create/Params/EarlyStopping';
import MetricsCollectorSpec from '../../Common/Create/Params/MetricsCollector';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Checkbox from '@material-ui/core/Checkbox';

import { submitHPJob } from '../../../actions/hpCreateActions';

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

const addParameter = (source, destination) => {
  source.map((param, i) => {
    let tempParam = {};
    tempParam.name = param.name;
    tempParam.parameterType = param.parameterType;
    tempParam.feasibleSpace = {};
    if (param.feasibleSpace === 'list') {
      tempParam.feasibleSpace.list = param.list.map((param, i) => param.value);
    } else {
      tempParam.feasibleSpace.min = param.min;
      tempParam.feasibleSpace.max = param.max;
      if (param.step !== '') {
        tempParam.feasibleSpace.step = param.step;
      }
    }
    return destination.push(tempParam);
  });
};

const HPParameters = props => {
  const submitJob = () => {
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

    // Add early stopping if selected.
    if (checkedSetEarlyStopping) {
      data.spec.earlyStopping = {};
      data.spec.earlyStopping.algorithmName = props.earlyStoppingAlgorithm;
      data.spec.earlyStopping.algorithmSettings = [];
      addAlgorithmSettings(props.earlyStoppingSettings, data.spec.earlyStopping.algorithmSettings);
    }

    data.spec.parameters = [];
    addParameter(props.parameters, data.spec.parameters);

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

    props.submitHPJob(data);
  };

  const { classes } = props;

  const [checkedSetEarlyStopping, setCheckedSetEarlyStopping] = React.useState(false);

  const onCheckBoxChange = event => {
    setCheckedSetEarlyStopping(event.target.checked);
  };

  return (
    <div className={classes.root}>
      {/* Common Metadata */}
      {SectionInTypography('Metadata')}
      <br />
      <CommonParametersMeta />
      {SectionInTypography('Common Parameters')}
      <CommonParametersSpec />
      {SectionInTypography('Objective')}
      <Objective />
      {SectionInTypography('Algorithm')}
      <Algorithm />

      <Grid container spacing={3}>
        <Grid item>
          <Typography variant="h6">Early Stopping (Optional)</Typography>
        </Grid>
        <Grid item>
          <FormControlLabel
            control={
              <Checkbox
                checked={checkedSetEarlyStopping}
                onChange={onCheckBoxChange}
                color="primary"
              />
            }
            label="Set"
          />
        </Grid>
      </Grid>
      {checkedSetEarlyStopping && <EarlyStopping />}

      {SectionInTypography('Parameters')}
      <Parameters />
      {SectionInTypography('Metrics Collector Spec')}
      <MetricsCollectorSpec jobType={constants.EXPERIMENT_TYPE_HP} />
      {SectionInTypography('Trial Template Spec')}
      <TrialTemplate />

      <div className={classes.submit}>
        <Button variant="contained" color={'primary'} onClick={submitJob}>
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
    commonParametersMetadata: state[constants.HP_CREATE_MODULE].commonParametersMetadata,
    commonParametersSpec: state[constants.HP_CREATE_MODULE].commonParametersSpec,
    objective: state[constants.HP_CREATE_MODULE].objective,
    additionalMetricNames: state[constants.HP_CREATE_MODULE].additionalMetricNames,
    metricStrategies: state[constants.HP_CREATE_MODULE].metricStrategies,
    algorithmName: state[constants.HP_CREATE_MODULE].algorithmName,
    earlyStoppingAlgorithm: state[constants.GENERAL_MODULE].earlyStoppingAlgorithm,
    earlyStoppingSettings: state[constants.GENERAL_MODULE].earlyStoppingSettings,
    algorithmSettings: state[constants.HP_CREATE_MODULE].algorithmSettings,
    parameters: state[constants.HP_CREATE_MODULE].parameters,
    primaryPodLabels: state[constants.GENERAL_MODULE].primaryPodLabels,
    trialTemplateSpec: state[constants.GENERAL_MODULE].trialTemplateSpec,
    trialTemplateSource: state[constants.GENERAL_MODULE].trialTemplateSource,
    templateConfigMapNamespace: templateCMNamespace,
    templateConfigMapName: templateCMName,
    templateConfigMapPath: templateCMPath,
    trialTemplateYAML: state[constants.GENERAL_MODULE].trialTemplateYAML,
    trialParameters: state[constants.GENERAL_MODULE].trialParameters,
    mcSpec: state[constants.HP_CREATE_MODULE].mcSpec,
    mcCustomContainerYaml: state[constants.HP_CREATE_MODULE].mcCustomContainerYaml,
  };
};

export default connect(mapStateToProps, { submitHPJob, validationError })(
  withStyles(styles)(HPParameters),
);
