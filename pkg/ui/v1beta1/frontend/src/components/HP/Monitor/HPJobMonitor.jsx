import React from 'react';
import { connect } from 'react-redux';

import { withStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';

import FilterPanel from './FilterPanel';
import HPJobList from './HPJobList';

import { fetchHPJobs } from '../../../actions/hpMonitorActions';

const styles = theme => ({
  root: {
    width: '90%',
    margin: '0 auto',
    marginTop: 10,
  },
  text: {
    marginBottom: 20,
  },
});

class HPJobMonitor extends React.Component {
  componentDidMount() {
    this.props.fetchHPJobs();
  }

  render() {
    const { classes } = this.props;

    return (
      <div className={classes.root}>
        <Typography variant={'h5'} className={classes.text}>
          {'Experiment Monitor'}
        </Typography>
        <FilterPanel />
        <HPJobList />
      </div>
    );
  }
}

export default connect(null, { fetchHPJobs })(withStyles(styles)(HPJobMonitor));
