import React from 'react';
import { connect } from 'react-redux';

import { withStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';

import FilterPanel from './FilterPanel';
import NASJobList from './NASJobList';

import { fetchNASJobs } from '../../../actions/nasMonitorActions';

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

class NASJobMonitor extends React.Component {
  componentDidMount() {
    this.props.fetchNASJobs();
  }

  render() {
    const { classes } = this.props;
    return (
      <div className={classes.root}>
        <Typography variant={'h5'} className={classes.text}>
          {'Experiment Monitor'}
        </Typography>
        <FilterPanel />
        <NASJobList />
      </div>
    );
  }
}

export default connect(null, { fetchNASJobs })(withStyles(styles)(NASJobMonitor));
