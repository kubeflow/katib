import React from 'react';
import { connect } from 'react-redux';

import { withStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';

import FilterPanel from './FilterPanel';
import ExperimentList from './ExperimentList';

import { fetchExperiments } from '../../../actions/generalActions';
import {
  LINK_HP_MONITOR,
  EXPERIMENT_TYPE_HP,
  EXPERIMENT_TYPE_NAS,
} from '../../../constants/constants';

const styles = () => ({
  root: {
    width: '90%',
    margin: '0 auto',
    marginTop: 10,
  },
  text: {
    marginBottom: 20,
  },
});

class ExperimentMonitor extends React.Component {
  componentDidMount() {
    this.props.fetchExperiments();
  }

  render() {
    const { classes } = this.props;

    return (
      <div className={classes.root}>
        <Typography variant={'h5'} className={classes.text}>
          {'Experiment Monitor'}
        </Typography>
        <FilterPanel />
        <ExperimentList
          experimentType={
            this.props.match.path === LINK_HP_MONITOR ? EXPERIMENT_TYPE_HP : EXPERIMENT_TYPE_NAS
          }
        />
      </div>
    );
  }
}

export default connect(null, { fetchExperiments })(withStyles(styles)(ExperimentMonitor));
