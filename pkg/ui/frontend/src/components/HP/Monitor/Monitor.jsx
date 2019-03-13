import React from 'react';
import { withStyles } from '@material-ui/core/styles';
 
import FilterPanel from './Panel';
import JobList from './JobList';

import { fetchHPJobs } from '../../../actions/hpMonitorActions';
import { connect } from 'react-redux'


const styles = theme => ({
    root: {
        width: '90%',
        margin: '0 auto',
    },
});

class HPMonitor extends React.Component {
    
    componentDidMount() {
        // this.props.fetchHPJobs();
    }

    render () {
        const { classes } = this.props;

        return (
            <div className={classes.root}>
                <h1>Monitor</h1>
                <FilterPanel />
                <JobList />
            </div>
        )
    }
}


export default connect(null, { fetchHPJobs })(withStyles(styles)(HPMonitor));