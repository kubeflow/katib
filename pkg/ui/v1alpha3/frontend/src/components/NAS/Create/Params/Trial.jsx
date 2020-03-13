import React from 'react';
import { connect } from 'react-redux';

import withStyles from '@material-ui/styles/withStyles';
import Grid from '@material-ui/core/Grid';
import Tooltip from '@material-ui/core/Tooltip';
import HelpOutlineIcon from '@material-ui/icons/HelpOutline';
import Typography from '@material-ui/core/Typography';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';

import { filterTemplatesExperiment, changeTemplateName } from '../../../../actions/generalActions';
import { fetchTrialTemplates } from '../../../../actions/templateActions';

const module = 'nasCreate';
const generalModule = 'general';

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
  trialForm: {
    margin: 4,
    width: '100%',
  },
  selectForm: {
    margin: 4,
    width: '20%',
  },
  selectNS: {
    marginRight: 10,
  },
});

class TrialSpecParam extends React.Component {
  componentDidMount() {
    this.props.fetchTrialTemplates();
  }

  onTrialNamespaceChange = event => {
    this.props.filterTemplatesExperiment(event.target.value, '');
  };

  onTrialConfigMapChange = event => {
    this.props.filterTemplatesExperiment(this.props.templateNamespace, event.target.value);
  };

  onTrialTemplateChange = event => {
    this.props.changeTemplateName(event.target.value);
  };

  render() {
    const { classes } = this.props;
    return (
      <div>
        <div className={classes.parameter}>
          <Grid container alignItems={'center'}>
            <Grid item xs={12} sm={3}>
              <Typography variant={'subheading'}>
                <Tooltip title={'Namespace and ConfigMap for Trial Template'}>
                  <HelpOutlineIcon className={classes.help} color={'primary'} />
                </Tooltip>
                {'Namespace and ConfigMapName'}
              </Typography>
            </Grid>
            <Grid item xs={12} sm={8}>
              <FormControl variant="outlined" className={classes.selectForm}>
                <InputLabel>Namespace</InputLabel>
                <Select
                  value={this.props.templateNamespace}
                  onChange={this.onTrialNamespaceChange}
                  className={classes.selectNS}
                  input={<OutlinedInput labelWidth={120} />}
                >
                  {this.props.trialTemplatesList.map((trialTemplate, i) => {
                    return (
                      <MenuItem value={trialTemplate.Namespace} key={i}>
                        {trialTemplate.Namespace}
                      </MenuItem>
                    );
                  })}
                </Select>
              </FormControl>
              <FormControl variant="outlined" className={classes.selectForm}>
                <InputLabel>ConfigMap</InputLabel>
                <Select
                  value={this.props.templateConfigMapName}
                  onChange={this.onTrialConfigMapChange}
                  input={<OutlinedInput labelWidth={120} />}
                >
                  {this.props.currentTemplateConfigMapsList.map((name, i) => {
                    return (
                      <MenuItem value={name} key={i}>
                        {name}
                      </MenuItem>
                    );
                  })}
                </Select>
              </FormControl>
            </Grid>
          </Grid>
        </div>
        <div className={classes.parameter}>
          <Grid container alignItems={'center'}>
            <Grid item xs={12} sm={3}>
              <Typography variant={'subheading'}>
                <Tooltip title={'Trial Template Path in ConfigMap'}>
                  <HelpOutlineIcon className={classes.help} color={'primary'} />
                </Tooltip>
                {'Trial Template Name'}
              </Typography>
            </Grid>
            <Grid item xs={12} sm={8}>
              <FormControl variant="outlined" className={classes.trialForm}>
                <InputLabel>Trial Template</InputLabel>
                <Select
                  value={this.props.templateName}
                  onChange={this.onTrialTemplateChange}
                  input={<OutlinedInput name={'TrialSpec'} labelWidth={100} />}
                >
                  {this.props.currentTemplateNamesList.map((name, i) => {
                    return (
                      <MenuItem value={name} key={i}>
                        {name}
                      </MenuItem>
                    );
                  })}
                </Select>
              </FormControl>
            </Grid>
          </Grid>
        </div>
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    trialTemplatesList: state[generalModule].trialTemplatesList,
    templateNamespace: state[generalModule].templateNamespace,
    templateConfigMapName: state[generalModule].templateConfigMapName,
    currentTemplateConfigMapsList: state[generalModule].currentTemplateConfigMapsList,
    templateName: state[generalModule].templateName,
    currentTemplateNamesList: state[generalModule].currentTemplateNamesList,
  };
};

export default connect(mapStateToProps, {
  filterTemplatesExperiment,
  changeTemplateName,
  fetchTrialTemplates,
})(withStyles(styles)(TrialSpecParam));
