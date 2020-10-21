import React from 'react';
import { connect } from 'react-redux';

import AceEditor from 'react-ace';
import 'ace-builds/src-noconflict/theme-sqlserver';
import 'ace-builds/src-noconflict/mode-yaml';

import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Typography from '@material-ui/core/Typography';
import Tooltip from '@material-ui/core/Tooltip';
import HelpOutlineIcon from '@material-ui/icons/HelpOutline';
import FormControl from '@material-ui/core/FormControl';
import InputLabel from '@material-ui/core/InputLabel';
import Select from '@material-ui/core/Select';
import MenuItem from '@material-ui/core/MenuItem';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import IconButton from '@material-ui/core/IconButton';
import DeleteIcon from '@material-ui/icons/Delete';
import Divider from '@material-ui/core/Divider';

import * as constants from '../../../../constants/constants';
import {
  changeMCKindHP,
  changeMCFileSystemHP,
  addMCMetricsFormatHP,
  changeMCMetricsFormatHP,
  deleteMCMetricsFormatHP,
  changeMCHttpGetHP,
  addMCHttpGetHeaderHP,
  changeMCHttpGetHeaderHP,
  deleteMCHttpGetHeaderHP,
  changeMCCustomContainerHP,
} from '../../../../actions/hpCreateActions';

import {
  changeMCKindNAS,
  changeMCFileSystemNAS,
  addMCMetricsFormatNAS,
  changeMCMetricsFormatNAS,
  deleteMCMetricsFormatNAS,
  changeMCHttpGetNAS,
  addMCHttpGetHeaderNAS,
  changeMCHttpGetHeaderNAS,
  deleteMCHttpGetHeaderNAS,
  changeMCCustomContainerNAS,
} from '../../../../actions/nasCreateActions';

import {
  GENERAL_MODULE,
  HP_CREATE_MODULE,
  NAS_CREATE_MODULE,
} from '../../../../constants/constants';

const styles = theme => ({
  help: {
    padding: 4 / 2,
    verticalAlign: 'middle',
    marginRight: 5,
  },
  formSelect: {
    width: '70%',
  },
  textField: {
    width: '95%',
  },
  grid: {
    marginTop: 20,
    marginBottom: 30,
  },
  textsList: {
    marginBottom: 15,
  },
  headerButton: {
    marginTop: 10,
  },
});

class MetricsCollectorSpec extends React.Component {
  onMCKindChange = event => {
    this.props.jobType === constants.EXPERIMENT_TYPE_HP
      ? this.props.changeMCKindHP(event.target.value)
      : this.props.changeMCKindNAS(event.target.value);
  };

  onMCFileSystemKindChange = event => {
    this.props.jobType === constants.EXPERIMENT_TYPE_HP
      ? this.props.changeMCFileSystemHP(
          event.target.value,
          this.props.mcSpecHP.source.fileSystemPath.path,
        )
      : this.props.changeMCFileSystemNAS(
          event.target.value,
          this.props.mcSpecNAS.source.fileSystemPath.path,
        );
  };

  onMCFileSystemPathChange = event => {
    this.props.jobType === constants.EXPERIMENT_TYPE_HP
      ? this.props.changeMCFileSystemHP(
          this.props.mcSpecHP.source.fileSystemPath.kind,
          event.target.value,
        )
      : this.props.changeMCFileSystemNAS(
          this.props.mcSpecNAS.source.fileSystemPath.kind,
          event.target.value,
        );
  };

  onMCMetricsFormatAdd = () => {
    this.props.jobType === constants.EXPERIMENT_TYPE_HP
      ? this.props.addMCMetricsFormatHP()
      : this.props.addMCMetricsFormatNAS();
  };

  onMCMetricsFormatChange = index => event => {
    this.props.jobType === constants.EXPERIMENT_TYPE_HP
      ? this.props.changeMCMetricsFormatHP(event.target.value, index)
      : this.props.changeMCMetricsFormatNAS(event.target.value, index);
  };

  onMCMetricsFormatDelete = index => event => {
    this.props.jobType === constants.EXPERIMENT_TYPE_HP
      ? this.props.deleteMCMetricsFormatHP(index)
      : this.props.deleteMCMetricsFormatNAS(index);
  };

  onMCHttpGetPortChange = event => {
    this.props.jobType === constants.EXPERIMENT_TYPE_HP
      ? this.props.changeMCHttpGetHP(
          event.target.value,
          this.props.mcSpecHP.source.httpGet.path,
          this.props.mcSpecHP.source.httpGet.scheme,
          this.props.mcSpecHP.source.httpGet.host,
        )
      : this.props.changeMCHttpGetNAS(
          event.target.value,
          this.props.mcSpecNAS.source.httpGet.path,
          this.props.mcSpecNAS.source.httpGet.scheme,
          this.props.mcSpecNAS.source.httpGet.host,
        );
  };

  onMCHttpGetPathChange = event => {
    this.props.jobType === constants.EXPERIMENT_TYPE_HP
      ? this.props.changeMCHttpGetHP(
          this.props.mcSpecHP.source.httpGet.port,
          event.target.value,
          this.props.mcSpecHP.source.httpGet.scheme,
          this.props.mcSpecHP.source.httpGet.host,
        )
      : this.props.changeMCHttpGetNAS(
          this.props.mcSpecNAS.source.httpGet.port,
          event.target.value,
          this.props.mcSpecNAS.source.httpGet.scheme,
          this.props.mcSpecNAS.source.httpGet.host,
        );
  };

  onMCHttpGetSchemeChange = event => {
    this.props.jobType === constants.EXPERIMENT_TYPE_HP
      ? this.props.changeMCHttpGetHP(
          this.props.mcSpecHP.source.httpGet.port,
          this.props.mcSpecHP.source.httpGet.path,
          event.target.value,
          this.props.mcSpecHP.source.httpGet.host,
        )
      : this.props.changeMCHttpGetNAS(
          this.props.mcSpecNAS.source.httpGet.port,
          this.props.mcSpecNAS.source.httpGet.path,
          event.target.value,
          this.props.mcSpecNAS.source.httpGet.host,
        );
  };

  onMCHttpGetHostChange = event => {
    this.props.jobType === constants.EXPERIMENT_TYPE_HP
      ? this.props.changeMCHttpGetHP(
          this.props.mcSpecHP.source.httpGet.port,
          this.props.mcSpecHP.source.httpGet.path,
          this.props.mcSpecHP.source.httpGet.scheme,
          event.target.value,
        )
      : this.props.changeMCHttpGetNAS(
          this.props.mcSpecNAS.source.httpGet.port,
          this.props.mcSpecNAS.source.httpGet.path,
          this.props.mcSpecNAS.source.httpGet.scheme,
          event.target.value,
        );
  };

  onMCHttpGetHeaderAdd = () => {
    this.props.jobType === constants.EXPERIMENT_TYPE_HP
      ? this.props.addMCHttpGetHeaderHP()
      : this.props.addMCHttpGetHeaderNAS();
  };

  onMCHttpGetHeaderChange = (fieldName, index) => event => {
    this.props.jobType === constants.EXPERIMENT_TYPE_HP
      ? this.props.changeMCHttpGetHeaderHP(fieldName, event.target.value, index)
      : this.props.changeMCHttpGetHeaderNAS(fieldName, event.target.value, index);
  };

  onMCHttpGetHeaderDelete = index => event => {
    this.props.jobType === constants.EXPERIMENT_TYPE_HP
      ? this.props.deleteMCHttpGetHeaderHP(index)
      : this.props.deleteMCHttpGetHeaderNAS(index);
  };

  onMCCustomContainerChange = yamlContainer => {
    this.props.jobType === constants.EXPERIMENT_TYPE_HP
      ? this.props.changeMCCustomContainerHP(yamlContainer)
      : this.props.changeMCCustomContainerNAS(yamlContainer);
  };
  render() {
    const { classes } = this.props;
    return (
      <div>
        <Grid container alignItems={'center'} className={classes.grid}>
          <Grid item xs={3}>
            <Typography variant={'subtitle1'}>
              <Tooltip title={'Kind for the Metrics Collector Spec'}>
                <HelpOutlineIcon className={classes.help} color={'primary'} />
              </Tooltip>
              {'Kind'}
            </Typography>
          </Grid>
          <Grid item xs={3}>
            <FormControl variant="outlined" className={classes.formSelect}>
              <InputLabel>Kind</InputLabel>
              <Select
                value={
                  this.props.jobType === constants.EXPERIMENT_TYPE_HP
                    ? this.props.mcSpecHP.collector.kind
                    : this.props.mcSpecNAS.collector.kind
                }
                onChange={this.onMCKindChange}
                label="Kind"
              >
                {this.props.mcKindsList.map((kind, i) => {
                  return (
                    <MenuItem value={kind} key={i}>
                      {kind}
                    </MenuItem>
                  );
                })}
              </Select>
            </FormControl>
          </Grid>
        </Grid>
        {((this.props.jobType === constants.EXPERIMENT_TYPE_HP &&
          (this.props.mcSpecHP.collector.kind === constants.MC_KIND_FILE ||
            this.props.mcSpecHP.collector.kind === constants.MC_KIND_TENSORFLOW_EVENT ||
            this.props.mcSpecHP.collector.kind === constants.MC_KIND_CUSTOM)) ||
          (this.props.jobType === constants.EXPERIMENT_TYPE_NAS &&
            (this.props.mcSpecNAS.collector.kind === constants.MC_KIND_FILE ||
              this.props.mcSpecNAS.collector.kind === constants.MC_KIND_TENSORFLOW_EVENT ||
              this.props.mcSpecNAS.collector.kind === constants.MC_KIND_CUSTOM))) && (
          <Grid container alignItems={'center'} className={classes.grid}>
            <Grid item xs={3}>
              <Typography variant={'subtitle1'}>
                <Tooltip
                  title={
                    'Kind of the file path and path to the metrics file, path must be absolute'
                  }
                >
                  <HelpOutlineIcon className={classes.help} color={'primary'} />
                </Tooltip>
                {'File System Kind and Path'}
              </Typography>
            </Grid>
            <Grid item xs={3}>
              <FormControl variant="outlined" className={classes.formSelect}>
                <InputLabel>File System Kind</InputLabel>
                <Select
                  value={
                    this.props.jobType === constants.EXPERIMENT_TYPE_HP
                      ? this.props.mcSpecHP.source.fileSystemPath.kind
                      : this.props.mcSpecNAS.source.fileSystemPath.kind
                  }
                  onChange={this.onMCFileSystemKindChange}
                  label="File System Kind"
                >
                  {((this.props.jobType === constants.EXPERIMENT_TYPE_HP &&
                    this.props.mcSpecHP.collector.kind === constants.MC_KIND_FILE) ||
                    (this.props.jobType === constants.EXPERIMENT_TYPE_NAS &&
                      this.props.mcSpecNAS.collector.kind === constants.MC_KIND_FILE)) && (
                    <MenuItem value={constants.MC_FILE_SYSTEM_KIND_FILE} key={0}>
                      {constants.MC_FILE_SYSTEM_KIND_FILE}
                    </MenuItem>
                  )}
                  {((this.props.jobType === constants.EXPERIMENT_TYPE_HP &&
                    this.props.mcSpecHP.collector.kind === constants.MC_KIND_TENSORFLOW_EVENT) ||
                    (this.props.jobType === constants.EXPERIMENT_TYPE_NAS &&
                      this.props.mcSpecNAS.collector.kind ===
                        constants.MC_KIND_TENSORFLOW_EVENT)) && (
                    <MenuItem value={constants.MC_FILE_SYSTEM_KIND_DIRECTORY} key={0}>
                      {constants.MC_FILE_SYSTEM_KIND_DIRECTORY}
                    </MenuItem>
                  )}
                  {((this.props.jobType === constants.EXPERIMENT_TYPE_HP &&
                    this.props.mcSpecHP.collector.kind === constants.MC_KIND_CUSTOM) ||
                    (this.props.jobType === constants.EXPERIMENT_TYPE_NAS &&
                      this.props.mcSpecNAS.collector.kind === constants.MC_KIND_CUSTOM)) &&
                    this.props.mcFileSystemKindsList.map((kind, i) => {
                      return (
                        <MenuItem value={kind} key={i}>
                          {kind}
                        </MenuItem>
                      );
                    })}
                </Select>
              </FormControl>
            </Grid>
            {((this.props.jobType === constants.EXPERIMENT_TYPE_HP &&
              this.props.mcSpecHP.source.fileSystemPath.kind !==
                constants.MC_FILE_SYSTEM_NO_KIND) ||
              (this.props.jobType === constants.EXPERIMENT_TYPE_NAS &&
                this.props.mcSpecNAS.source.fileSystemPath.kind !==
                  constants.MC_FILE_SYSTEM_NO_KIND)) && (
              <Grid item xs={3}>
                <TextField
                  label={'File System Path'}
                  className={classes.textField}
                  value={
                    this.props.jobType === constants.EXPERIMENT_TYPE_HP
                      ? this.props.mcSpecHP.source.fileSystemPath.path
                      : this.props.mcSpecNAS.source.fileSystemPath.path
                  }
                  onChange={this.onMCFileSystemPathChange}
                />
              </Grid>
            )}
          </Grid>
        )}
        {((this.props.jobType === constants.EXPERIMENT_TYPE_HP &&
          this.props.mcSpecHP.collector.kind === constants.MC_KIND_PROMETHEUS) ||
          (this.props.jobType === constants.EXPERIMENT_TYPE_NAS &&
            this.props.mcSpecNAS.collector.kind === constants.MC_KIND_PROMETHEUS)) && (
          <div>
            <Grid container alignItems={'center'} className={classes.grid}>
              <Grid item xs={3}>
                <Typography variant={'subtitle1'}>
                  <Tooltip
                    title={
                      'Port and Path to access on the HTTP server. Port must be a positive integer value. Path must start with "/"'
                    }
                  >
                    <HelpOutlineIcon className={classes.help} color={'primary'} />
                  </Tooltip>
                  {'HttpGet Port and Path'}
                </Typography>
              </Grid>
              <Grid item xs={3}>
                <TextField
                  label={'HttpGet Port'}
                  className={classes.textField}
                  type="number"
                  value={
                    this.props.jobType === constants.EXPERIMENT_TYPE_HP
                      ? this.props.mcSpecHP.source.httpGet.port
                      : this.props.mcSpecNAS.source.httpGet.port
                  }
                  onChange={this.onMCHttpGetPortChange}
                />
              </Grid>
              <Grid item xs={3}>
                <TextField
                  label={'HttpGet Path'}
                  className={classes.textField}
                  value={
                    this.props.jobType === constants.EXPERIMENT_TYPE_HP
                      ? this.props.mcSpecHP.source.httpGet.path
                      : this.props.mcSpecNAS.source.httpGet.path
                  }
                  onChange={this.onMCHttpGetPathChange}
                />
              </Grid>
            </Grid>
            <Grid container alignItems={'center'} className={classes.grid}>
              <Grid item xs={3}>
                <Typography variant={'subtitle1'}>
                  <Tooltip
                    title={
                      'Scheme to use for connecting to the host. Host name to make connection, defaults to the pod IP'
                    }
                  >
                    <HelpOutlineIcon className={classes.help} color={'primary'} />
                  </Tooltip>
                  {'HttpGet Scheme and Host (optional)'}
                </Typography>
              </Grid>
              <Grid item xs={3}>
                <FormControl variant="outlined" className={classes.formSelect}>
                  <InputLabel>Scheme</InputLabel>
                  <Select
                    value={
                      this.props.jobType === constants.EXPERIMENT_TYPE_HP
                        ? this.props.mcSpecHP.source.httpGet.scheme
                        : this.props.mcSpecNAS.source.httpGet.scheme
                    }
                    onChange={this.onMCHttpGetSchemeChange}
                    label="Scheme"
                  >
                    {this.props.mcURISchemesList.map((scheme, i) => {
                      return (
                        <MenuItem value={scheme} key={i}>
                          {scheme}
                        </MenuItem>
                      );
                    })}
                  </Select>
                </FormControl>
              </Grid>
              <Grid item xs={3}>
                <TextField
                  label={'Host name'}
                  className={classes.textField}
                  value={
                    this.props.jobType === constants.EXPERIMENT_TYPE_HP
                      ? this.props.mcSpecHP.source.httpGet.host
                      : this.props.mcSpecNAS.source.httpGet.host
                  }
                  onChange={this.onMCHttpGetHostChange}
                />
              </Grid>
            </Grid>
            <Grid container alignItems={'center'} className={classes.grid}>
              <Grid item xs={3}>
                <Typography variant={'subtitle1'}>
                  <Tooltip title={'Custom headers to set in the request'}>
                    <HelpOutlineIcon className={classes.help} color={'primary'} />
                  </Tooltip>
                  {'HttpGet Headers (optional)'}
                </Typography>
                <Button
                  variant={'contained'}
                  color={'primary'}
                  className={classes.headerButton}
                  onClick={this.onMCHttpGetHeaderAdd}
                >
                  Add Header
                </Button>
              </Grid>
            </Grid>
            {this.props.jobType === constants.EXPERIMENT_TYPE_HP
              ? this.props.mcSpecHP.source.httpGet.httpHeaders.map((header, index) => {
                  return (
                    <div key={index} className={classes.textsList}>
                      <Grid container alignItems={'center'}>
                        <Grid item xs={3} />
                        <Grid item xs={3}>
                          <TextField
                            label={'Header Name'}
                            className={classes.textField}
                            value={header.name}
                            onChange={this.onMCHttpGetHeaderChange('name', index)}
                          />
                        </Grid>
                        <Grid item xs={3}>
                          <TextField
                            label={'Header Value'}
                            className={classes.textField}
                            value={header.value}
                            onChange={this.onMCHttpGetHeaderChange('value', index)}
                          />
                        </Grid>
                        <Grid item xs={1}>
                          <IconButton
                            aria-label="Close"
                            color={'primary'}
                            onClick={this.onMCHttpGetHeaderDelete(index)}
                          >
                            <DeleteIcon />
                          </IconButton>
                        </Grid>
                      </Grid>
                    </div>
                  );
                })
              : this.props.mcSpecNAS.source.httpGet.httpHeaders.map((header, index) => {
                  return (
                    <div key={index} className={classes.textsList}>
                      <Grid container alignItems={'center'}>
                        <Grid item xs={3} />
                        <Grid item xs={3}>
                          <TextField
                            label={'Header Name'}
                            className={classes.textField}
                            value={header.name}
                            onChange={this.onMCHttpGetHeaderChange('name', index)}
                          />
                        </Grid>
                        <Grid item xs={3}>
                          <TextField
                            label={'Header Value'}
                            className={classes.textField}
                            value={header.value}
                            onChange={this.onMCHttpGetHeaderChange('value', index)}
                          />
                        </Grid>
                        <Grid item xs={1}>
                          <IconButton
                            aria-label="Close"
                            color={'primary'}
                            onClick={this.onMCHttpGetHeaderDelete(index)}
                          >
                            <DeleteIcon />
                          </IconButton>
                        </Grid>
                      </Grid>
                    </div>
                  );
                })}
          </div>
        )}
        {((this.props.jobType === constants.EXPERIMENT_TYPE_HP &&
          this.props.mcSpecHP.collector.kind === constants.MC_KIND_CUSTOM) ||
          (this.props.jobType === constants.EXPERIMENT_TYPE_NAS &&
            this.props.mcSpecNAS.collector.kind === constants.MC_KIND_CUSTOM)) && (
          <Grid container alignItems={'center'} className={classes.grid}>
            <Grid item xs={3}>
              <Typography variant={'subtitle1'}>
                <Tooltip title={'Yaml structure for the custom metrics collector container'}>
                  <HelpOutlineIcon className={classes.help} color={'primary'} />
                </Tooltip>
                {'YAML for the Custom Container'}
              </Typography>
            </Grid>
            <Grid item xs={6}>
              <AceEditor
                mode="yaml"
                theme="sqlserver"
                value={
                  this.props.jobType === constants.EXPERIMENT_TYPE_HP
                    ? this.props.mcCustomContainerYamlHP
                    : this.props.mcCustomContainerYamlNAS
                }
                tabSize={2}
                fontSize={13}
                width={'100%'}
                showPrintMargin={false}
                autoScrollEditorIntoView={true}
                maxLines={40}
                minLines={20}
                onChange={this.onMCCustomContainerChange}
              />
            </Grid>
          </Grid>
        )}

        <Divider />
        <Grid container alignItems={'center'} className={classes.grid}>
          <Grid item xs={3}>
            <Button variant={'contained'} color={'primary'} onClick={this.onMCMetricsFormatAdd}>
              Add Metrics Format
            </Button>
          </Grid>
        </Grid>
        {this.props.jobType === constants.EXPERIMENT_TYPE_HP &&
          this.props.mcSpecHP.source !== undefined &&
          this.props.mcSpecHP.source.filter !== undefined &&
          this.props.mcSpecHP.source.filter.metricsFormat.map((format, index) => {
            return (
              <div key={index} className={classes.textsList}>
                <Grid container alignItems={'center'}>
                  <Grid item xs={3} />
                  <Grid item xs={6}>
                    <TextField
                      label={'Metrics Format regular expression'}
                      className={classes.textField}
                      value={format}
                      onChange={this.onMCMetricsFormatChange(index)}
                    />
                  </Grid>
                  <Grid item xs={1}>
                    <IconButton
                      aria-label="Close"
                      color={'primary'}
                      onClick={this.onMCMetricsFormatDelete(index)}
                    >
                      <DeleteIcon />
                    </IconButton>
                  </Grid>
                </Grid>
              </div>
            );
          })}
        {this.props.jobType === constants.EXPERIMENT_TYPE_NAS &&
          this.props.mcSpecNAS.source !== undefined &&
          this.props.mcSpecNAS.source.filter !== undefined &&
          this.props.mcSpecNAS.source.filter.metricsFormat.map((format, index) => {
            return (
              <div key={index} className={classes.textsList}>
                <Grid container alignItems={'center'}>
                  <Grid item xs={3} />
                  <Grid item xs={6}>
                    <TextField
                      label={'Metrics Format regular expression'}
                      className={classes.textField}
                      value={format}
                      onChange={this.onMCMetricsFormatChange(index)}
                    />
                  </Grid>
                  <Grid item xs={1}>
                    <IconButton
                      aria-label="Close"
                      color={'primary'}
                      onClick={this.onMCMetricsFormatDelete(index)}
                    >
                      <DeleteIcon />
                    </IconButton>
                  </Grid>
                </Grid>
              </div>
            );
          })}
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    mcSpecHP: state[HP_CREATE_MODULE].mcSpec,
    mcSpecNAS: state[NAS_CREATE_MODULE].mcSpec,
    mcCustomContainerYamlHP: state[HP_CREATE_MODULE].mcCustomContainerYaml,
    mcCustomContainerYamlNAS: state[NAS_CREATE_MODULE].mcCustomContainerYaml,
    mcKindsList: state[GENERAL_MODULE].mcKindsList,
    mcFileSystemKindsList: state[GENERAL_MODULE].mcFileSystemKindsList,
    mcURISchemesList: state[GENERAL_MODULE].mcURISchemesList,
  };
};

export default connect(mapStateToProps, {
  changeMCKindHP,
  changeMCFileSystemHP,
  addMCMetricsFormatHP,
  changeMCMetricsFormatHP,
  deleteMCMetricsFormatHP,
  changeMCHttpGetHP,
  addMCHttpGetHeaderHP,
  changeMCHttpGetHeaderHP,
  deleteMCHttpGetHeaderHP,
  changeMCCustomContainerHP,
  changeMCKindNAS,
  changeMCFileSystemNAS,
  addMCMetricsFormatNAS,
  changeMCMetricsFormatNAS,
  deleteMCMetricsFormatNAS,
  changeMCHttpGetNAS,
  addMCHttpGetHeaderNAS,
  changeMCHttpGetHeaderNAS,
  deleteMCHttpGetHeaderNAS,
  changeMCCustomContainerNAS,
})(withStyles(styles)(MetricsCollectorSpec));
