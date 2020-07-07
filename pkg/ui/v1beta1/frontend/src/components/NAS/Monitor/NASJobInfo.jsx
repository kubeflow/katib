import React from 'react';
import { connect } from 'react-redux';

import { Link } from 'react-router-dom';

import { withStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import Typography from '@material-ui/core/Typography';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import LinearProgress from '@material-ui/core/LinearProgress';
import Grid from '@material-ui/core/Grid';

import NASJobStepInfo from './NASJobStepInfo';
import ExperimentInfoDialog from '../../Common/ExperimentInfoDialog';
import SuggestionInfoDialog from '../../Common/SuggestionInfoDialog';

import { fetchNASJobInfo } from '../../../actions/nasMonitorActions';
import { fetchExperiment, fetchSuggestion } from '../../../actions/generalActions';

import { NAS_MONITOR_MODULE } from '../../../constants/constants';

const styles = theme => ({
  root: {
    width: '90%',
    margin: '0 auto',
    padding: 20,
  },
  loading: {
    marginTop: 30,
  },
  heading: {
    fontSize: theme.typography.pxToRem(15),
    fontWeight: theme.typography.fontWeightRegular,
  },
  panel: {
    width: '100%',
  },
  header: {
    marginTop: 10,
    textAlign: 'center',
    marginBottom: 15,
  },
  link: {
    textDecoration: 'none',
  },
  grid: {
    marginBottom: 10,
  },
});

class NASJobInfo extends React.Component {
  componentDidMount() {
    this.props.fetchNASJobInfo(this.props.match.params.name, this.props.match.params.namespace);
  }

  fetchAndOpenDialogExperiment = (experimentName, experimentNamespace) => event => {
    this.props.fetchExperiment(experimentName, experimentNamespace);
  };

  fetchAndOpenDialogSuggestion = (suggestionName, suggestionNamespace) => event => {
    this.props.fetchSuggestion(suggestionName, suggestionNamespace);
  };

  render() {
    const { classes } = this.props;
    return (
      <div className={classes.root}>
        <Link to="/katib/nas_monitor" className={classes.link}>
          <Button variant={'contained'} color={'primary'}>
            Back
          </Button>
        </Link>
        {this.props.loading ? (
          <LinearProgress color={'primary'} className={classes.loading} />
        ) : (
          <div>
            <Typography className={classes.header} variant={'h5'}>
              Experiment Name: {this.props.match.params.name}
            </Typography>
            <Typography className={classes.header} variant={'h5'}>
              Experiment Namespace: {this.props.match.params.namespace}
            </Typography>

            <Grid container className={classes.grid} justify="center" spacing={3}>
              <Grid item>
                <Button
                  variant={'contained'}
                  color={'primary'}
                  onClick={this.fetchAndOpenDialogExperiment(
                    this.props.match.params.name,
                    this.props.match.params.namespace,
                  )}
                >
                  View Experiment
                </Button>
              </Grid>
              <Grid item>
                <Button
                  variant={'contained'}
                  color={'primary'}
                  onClick={this.fetchAndOpenDialogSuggestion(
                    this.props.match.params.name,
                    this.props.match.params.namespace,
                  )}
                >
                  View Suggestion
                </Button>
              </Grid>
            </Grid>
            {this.props.steps.map((step, i) => {
              return (
                <ExpansionPanel key={i} className={classes.panel}>
                  <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
                    <Typography className={classes.heading}>{step.name}</Typography>
                  </ExpansionPanelSummary>
                  <ExpansionPanelDetails>
                    <NASJobStepInfo step={step} id={i + 1} />
                  </ExpansionPanelDetails>
                </ExpansionPanel>
              );
            })}
            <ExperimentInfoDialog />
            <SuggestionInfoDialog />
          </div>
        )}
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    steps: state[NAS_MONITOR_MODULE].steps,
    loading: state[NAS_MONITOR_MODULE].loading,
  };
};

export default connect(mapStateToProps, { fetchNASJobInfo, fetchExperiment, fetchSuggestion })(
  withStyles(styles)(NASJobInfo),
);
