import React from 'react';
import { connect } from 'react-redux';

import { Link } from 'react-router-dom';

import { withStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import LinearProgress from '@material-ui/core/LinearProgress';
import Grid from '@material-ui/core/Grid';

import HPJobPlot from './HPJobPlot';
import HPJobTable from './HPJobTable';
import TrialInfoDialog from './TrialInfoDialog';
import ExperimentInfoDialog from '../../Common/ExperimentInfoDialog';
import SuggestionInfoDialog from '../../Common/SuggestionInfoDialog';

import { fetchHPJobInfo } from '../../../actions/hpMonitorActions';
import { fetchExperiment, fetchSuggestion } from '../../../actions/generalActions';
import { HP_MONITOR_MODULE } from '../../../constants/constants';

const styles = theme => ({
  root: {
    width: '90%',
    margin: '0 auto',
    padding: 20,
  },
  loading: {
    marginTop: 30,
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

class HPJobInfo extends React.Component {
  componentDidMount() {
    this.props.fetchHPJobInfo(this.props.match.params.name, this.props.match.params.namespace);
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
        <Link to="/katib/hp_monitor" className={classes.link}>
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
            <HPJobPlot name={this.props.match.params.name} />
            <HPJobTable namespace={this.props.match.params.namespace} />
            <ExperimentInfoDialog />
            <SuggestionInfoDialog />
            <TrialInfoDialog />
          </div>
        )}
      </div>
    );
  }
}

const mapStateToProps = state => ({
  loading: state[HP_MONITOR_MODULE].loading,
});

export default connect(mapStateToProps, { fetchHPJobInfo, fetchExperiment, fetchSuggestion })(
  withStyles(styles)(HPJobInfo),
);
