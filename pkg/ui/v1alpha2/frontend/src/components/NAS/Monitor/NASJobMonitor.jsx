import React from 'react';
import { withStyles } from '@material-ui/core/styles';

import FilterPanel from './FilterPanel';
import NASJobList from './NASJobList';

import { connect } from 'react-redux';

import { fetchNASJobs } from '../../../actions/nasMonitorActions';


const styles = theme => ({
    root: {
        width: '90%',
        margin: '0 auto',
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
                <h1>Monitor</h1>
                <FilterPanel />
                <NASJobList />
            </div>
        )
    }
}


export default connect(null, { fetchNASJobs })(withStyles(styles)(NASJobMonitor));
