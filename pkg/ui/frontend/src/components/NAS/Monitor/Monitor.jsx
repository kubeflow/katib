import React from 'react';
import { withStyles } from '@material-ui/core/styles';

import FilterPanel from './Panel';
import JobList from './JobList';

const styles = theme => ({
    root: {
        width: '90%',
        margin: '0 auto',
    },
});

const NASMonitor = (props) => {

    const { classes } = props;
    return (
        <div className={classes.root}>
            <h1>Monitor</h1>
            <FilterPanel />
            <JobList />
        </div>
    )

}

export default withStyles(styles)(NASMonitor);