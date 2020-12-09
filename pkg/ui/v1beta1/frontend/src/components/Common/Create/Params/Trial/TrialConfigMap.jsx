import React from 'react';
import { connect } from 'react-redux';

import AceEditor from 'react-ace';

import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Tooltip from '@material-ui/core/Tooltip';
import HelpOutlineIcon from '@material-ui/icons/HelpOutline';
import Typography from '@material-ui/core/Typography';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';

import { filterTemplatesExperiment } from '../../../../../actions/generalActions';
import { fetchTrialTemplates } from '../../../../../actions/templateActions';
import { GENERAL_MODULE } from '../../../../../constants/constants';

const styles = theme => ({
  help: {
    padding: 4 / 2,
    verticalAlign: 'middle',
    marginRight: 5,
  },
  section: {
    padding: 4,
  },
  parameter: {
    padding: 2,
    marginBottom: 10,
  },
  selectForm: {
    width: '80%',
  },
  yamlHeader: {
    fontSize: 17,
  },
});

class TrialConfigMap extends React.Component {
  componentDidMount() {
    this.props.fetchTrialTemplates();
  }

  onTemplateConfigMapNamespaceChange = event => {
    let nsIndex = this.props.trialTemplatesData.findIndex(function (trialTemplate, i) {
      return trialTemplate.ConfigMapNamespace === event.target.value;
    });

    this.props.filterTemplatesExperiment(
      nsIndex,
      this.props.configMapNameIndex,
      this.props.configMapPathIndex,
    );
  };

  onTemplateConfigMapNameChange = event => {
    let namespacedData = this.props.trialTemplatesData[this.props.configMapNamespaceIndex];
    let nameIndex = namespacedData.ConfigMaps.findIndex(function (configMap, i) {
      return configMap.ConfigMapName === event.target.value;
    });

    this.props.filterTemplatesExperiment(
      this.props.configMapNamespaceIndex,
      nameIndex,
      this.props.configMapPathIndex,
    );
  };

  onTemplateConfigMapPathChange = event => {
    let namespacedData = this.props.trialTemplatesData[this.props.configMapNamespaceIndex];
    let namedConfigMap = namespacedData.ConfigMaps[this.props.configMapNameIndex];

    let pathIndex = namedConfigMap.Templates.findIndex(function (template, i) {
      return template.Path === event.target.value;
    });

    this.props.filterTemplatesExperiment(
      this.props.configMapNamespaceIndex,
      this.props.configMapNameIndex,
      pathIndex,
    );
  };

  render() {
    const { classes } = this.props;
    return this.props.configMapNamespaceIndex !== -1 ? (
      <div>
        <div className={classes.parameter}>
          <Grid container alignItems={'center'}>
            <Grid item xs={3}>
              <Typography variant={'subtitle1'}>
                <Tooltip title={'Trial Template ConfigMap Namespace and Name'}>
                  <HelpOutlineIcon className={classes.help} color={'primary'} />
                </Tooltip>
                {'ConfigMap Namespace and Name'}
              </Typography>
            </Grid>
            <Grid item xs={9}>
              <Grid container>
                <Grid item xs={4}>
                  <FormControl variant="outlined" className={classes.selectForm}>
                    <InputLabel>Namespace</InputLabel>
                    <Select
                      value={
                        this.props.trialTemplatesData[this.props.configMapNamespaceIndex]
                          .ConfigMapNamespace
                      }
                      onChange={this.onTemplateConfigMapNamespaceChange}
                      label="Namespace"
                    >
                      {this.props.trialTemplatesData.map((trialTemplate, i) => {
                        return (
                          <MenuItem value={trialTemplate.ConfigMapNamespace} key={i}>
                            {trialTemplate.ConfigMapNamespace}
                          </MenuItem>
                        );
                      })}
                    </Select>
                  </FormControl>
                </Grid>
                <Grid item xs={4}>
                  <FormControl variant="outlined" className={classes.selectForm}>
                    <InputLabel>Name</InputLabel>
                    <Select
                      value={
                        this.props.trialTemplatesData[this.props.configMapNamespaceIndex]
                          .ConfigMaps[this.props.configMapNameIndex].ConfigMapName
                      }
                      onChange={this.onTemplateConfigMapNameChange}
                      label="Name"
                    >
                      {this.props.trialTemplatesData[
                        this.props.configMapNamespaceIndex
                      ].ConfigMaps.map((configMap, i) => {
                        return (
                          <MenuItem value={configMap.ConfigMapName} key={i}>
                            {configMap.ConfigMapName}
                          </MenuItem>
                        );
                      })}
                    </Select>
                  </FormControl>
                </Grid>
              </Grid>
            </Grid>
          </Grid>
        </div>
        <div className={classes.parameter}>
          <Grid container alignItems={'center'}>
            <Grid item xs={3}>
              <Typography variant={'subtitle1'}>
                <Tooltip title={'Trial Template Path in ConfigMap'}>
                  <HelpOutlineIcon className={classes.help} color={'primary'} />
                </Tooltip>
                {'Trial Template ConfigMap Path'}
              </Typography>
            </Grid>
            <Grid item xs={9}>
              <Grid container>
                <Grid item xs={4}>
                  <FormControl variant="outlined" className={classes.selectForm}>
                    <InputLabel>Template Path</InputLabel>
                    <Select
                      value={
                        this.props.trialTemplatesData[this.props.configMapNamespaceIndex]
                          .ConfigMaps[this.props.configMapNameIndex].Templates[
                          this.props.configMapPathIndex
                        ].Path
                      }
                      onChange={this.onTemplateConfigMapPathChange}
                      label="Template Path"
                    >
                      {this.props.trialTemplatesData[this.props.configMapNamespaceIndex].ConfigMaps[
                        this.props.configMapNameIndex
                      ].Templates.map((template, i) => {
                        return (
                          <MenuItem value={template.Path} key={i}>
                            {template.Path}
                          </MenuItem>
                        );
                      })}
                    </Select>
                  </FormControl>
                </Grid>
                <Grid item xs={7}>
                  <ExpansionPanel>
                    <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
                      <Typography className={classes.yamlHeader}>{'View template'}</Typography>
                    </ExpansionPanelSummary>
                    <ExpansionPanelDetails>
                      <AceEditor
                        mode="yaml"
                        theme="sqlserver"
                        value={
                          this.props.trialTemplatesData[this.props.configMapNamespaceIndex]
                            .ConfigMaps[this.props.configMapNameIndex].Templates[
                            this.props.configMapPathIndex
                          ].Yaml
                        }
                        tabSize={2}
                        fontSize={12}
                        width={'100%'}
                        showPrintMargin={false}
                        autoScrollEditorIntoView={true}
                        readOnly={true}
                        maxLines={160}
                        minLines={10}
                      />
                    </ExpansionPanelDetails>
                  </ExpansionPanel>
                </Grid>
              </Grid>
            </Grid>
          </Grid>
        </div>
      </div>
    ) : (
      <Typography variant="h6">No ConfigMaps with Katib Trial Templates</Typography>
    );
  }
}

const mapStateToProps = state => {
  return {
    configMapNamespaceIndex: state[GENERAL_MODULE].configMapNamespaceIndex,
    configMapNameIndex: state[GENERAL_MODULE].configMapNameIndex,
    configMapPathIndex: state[GENERAL_MODULE].configMapPathIndex,
    trialTemplatesData: state[GENERAL_MODULE].trialTemplatesData,
  };
};

export default connect(mapStateToProps, {
  filterTemplatesExperiment,
  fetchTrialTemplates,
})(withStyles(styles)(TrialConfigMap));
