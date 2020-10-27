import React from 'react';
import { connect } from 'react-redux';

import AceEditor from 'react-ace';
import 'ace-builds/src-noconflict/theme-sqlserver';
import 'ace-builds/src-noconflict/mode-yaml';

import { withStyles } from '@material-ui/core/styles';
import Divider from '@material-ui/core/Divider';
import FormControl from '@material-ui/core/FormControl';
import InputLabel from '@material-ui/core/InputLabel';
import Select from '@material-ui/core/Select';
import Grid from '@material-ui/core/Grid';
import Tooltip from '@material-ui/core/Tooltip';
import HelpOutlineIcon from '@material-ui/icons/HelpOutline';
import Typography from '@material-ui/core/Typography';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import IconButton from '@material-ui/core/IconButton';
import DeleteIcon from '@material-ui/icons/Delete';
import MenuItem from '@material-ui/core/MenuItem';
import TrialConfigMap from './TrialConfigMap';
import TrialParameters from './TrialParameters';

import { GENERAL_MODULE, TEMPLATE_SOURCE_CONFIG_MAP } from '../../../../../constants/constants';
import {
  changeTrialTemplateSource,
  addPrimaryPodLabel,
  changePrimaryPodLabel,
  deletePrimaryPodLabel,
  changeTrialTemplateSpec,
  changeTrialTemplateYAML,
} from '../../../../../actions/generalActions';
import { fetchTrialTemplates } from '../../../../../actions/templateActions';

const styles = () => ({
  parameter: {
    padding: 2,
    marginBottom: 10,
  },
  help: {
    padding: 4 / 2,
    verticalAlign: 'middle',
    marginRight: 5,
  },
  textField: {
    marginLeft: 4,
    marginRight: 4,
    width: '90%',
  },
  button: {
    marginTop: 10,
    marginBottom: 10,
  },
  formSelect: {
    width: '70%',
  },
});

class TrialTemplate extends React.Component {
  onTrialTemplateSourceChange = event => {
    // Change source only if value is changed.
    if (event.target.value !== this.props.trialTemplateSource) {
      this.props.changeTrialTemplateSource(event.target.value);
      // Fetch templates if source is ConfigMap.
      if (event.target.value === TEMPLATE_SOURCE_CONFIG_MAP) {
        this.props.fetchTrialTemplates();
      }
    }
  };

  onPrimaryPodLabelAdd = () => {
    this.props.addPrimaryPodLabel();
  };

  onPrimaryPodLabelChange = (fieldName, index) => event => {
    this.props.changePrimaryPodLabel(fieldName, index, event.target.value);
  };

  onPrimaryPodLabelDelete = index => () => {
    this.props.deletePrimaryPodLabel(index);
  };

  onTrialTemplateSpecChange = name => event => {
    this.props.changeTrialTemplateSpec(name, event.target.value);
  };

  onTrialTemplateYAMLChange = templateYAML => {
    this.props.changeTrialTemplateYAML(templateYAML);
  };

  render() {
    const { classes } = this.props;
    return (
      <div>
        <Grid container alignItems={'center'}>
          <Grid item xs={3}>
            <Typography variant={'subtitle1'}>
              <Tooltip title={'Source type for Trial template'}>
                <HelpOutlineIcon className={classes.help} color={'primary'} />
              </Tooltip>
              {'Source type'}
            </Typography>
          </Grid>
          <Grid item xs={3}>
            <FormControl variant="outlined" className={classes.formSelect}>
              <InputLabel>Source type</InputLabel>
              <Select
                value={this.props.trialTemplateSource}
                onChange={this.onTrialTemplateSourceChange}
                label="Source type"
              >
                {this.props.trialTemplateSourceList.map((source, i) => {
                  return (
                    <MenuItem value={source} key={i}>
                      {source}
                    </MenuItem>
                  );
                })}
              </Select>
            </FormControl>
          </Grid>
        </Grid>
        <Grid container alignItems={'center'}>
          <Grid item xs={3}>
            <Typography variant={'subtitle1'}>
              <Tooltip
                title={
                  'Labels to indicate if pod needs to be injected by Katib sidecar container.\
              If labels are omitted, all created pods are injected'
                }
              >
                <HelpOutlineIcon className={classes.help} color={'primary'} />
              </Tooltip>
              {'PrimaryPodLabels (optional)'}
            </Typography>
            <Button
              variant={'contained'}
              color={'primary'}
              className={classes.button}
              onClick={this.onPrimaryPodLabelAdd}
            >
              Add Label
            </Button>
          </Grid>
        </Grid>
        {this.props.primaryPodLabels.map((label, index) => {
          return (
            <div key={index} className={classes.parameter}>
              <Grid container alignItems={'center'}>
                <Grid item xs={3} />
                <Grid item xs={3}>
                  <TextField
                    label={'Label Key'}
                    className={classes.textField}
                    value={label.key}
                    onChange={this.onPrimaryPodLabelChange('key', index)}
                  />
                </Grid>
                <Grid item xs={3}>
                  <TextField
                    label={'Label Value'}
                    className={classes.textField}
                    value={label.value}
                    onChange={this.onPrimaryPodLabelChange('value', index)}
                  />
                </Grid>
                <Grid item xs={1}>
                  <IconButton
                    aria-label="Close"
                    color={'primary'}
                    onClick={this.onPrimaryPodLabelDelete(index)}
                  >
                    <DeleteIcon />
                  </IconButton>
                </Grid>
              </Grid>
            </div>
          );
        })}
        {this.props.trialTemplateSpec.map((param, i) => {
          return (
            <div key={i} className={classes.parameter}>
              <Grid container alignItems={'center'}>
                <Grid item xs={12} sm={3}>
                  <Typography variant={'subtitle1'}>
                    <Tooltip title={param.description}>
                      <HelpOutlineIcon className={classes.help} color={'primary'} />
                    </Tooltip>
                    {param.name}
                  </Typography>
                </Grid>
                <Grid item xs={12} sm={8}>
                  <TextField
                    className={classes.textField}
                    value={param.value}
                    onChange={this.onTrialTemplateSpecChange(param.name)}
                  />
                </Grid>
              </Grid>
            </div>
          );
        })}
        {this.props.trialTemplateSource === TEMPLATE_SOURCE_CONFIG_MAP ? (
          <TrialConfigMap></TrialConfigMap>
        ) : (
          <Grid container alignItems={'center'}>
            <Grid item xs={3}>
              <Typography variant={'subtitle1'}>
                <Tooltip title={'YAML structure for Trial template'}>
                  <HelpOutlineIcon className={classes.help} color={'primary'} />
                </Tooltip>
                {'Trial template YAML'}
              </Typography>
            </Grid>
            <Grid item xs={6}>
              <AceEditor
                mode="yaml"
                theme="sqlserver"
                value={this.props.trialTemplateYAML}
                tabSize={2}
                fontSize={13}
                width={'100%'}
                showPrintMargin={false}
                autoScrollEditorIntoView={true}
                maxLines={40}
                minLines={20}
                onChange={this.onTrialTemplateYAMLChange}
              />
            </Grid>
          </Grid>
        )}
        <Divider />
        <TrialParameters></TrialParameters>
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    trialTemplateSourceList: state[GENERAL_MODULE].trialTemplateSourceList,
    trialTemplateSource: state[GENERAL_MODULE].trialTemplateSource,
    primaryPodLabels: state[GENERAL_MODULE].primaryPodLabels,
    trialTemplateSpec: state[GENERAL_MODULE].trialTemplateSpec,
    trialTemplateYAML: state[GENERAL_MODULE].trialTemplateYAML,
  };
};

export default connect(mapStateToProps, {
  changeTrialTemplateSource,
  addPrimaryPodLabel,
  changePrimaryPodLabel,
  deletePrimaryPodLabel,
  changeTrialTemplateSpec,
  changeTrialTemplateYAML,
  fetchTrialTemplates,
})(withStyles(styles)(TrialTemplate));
