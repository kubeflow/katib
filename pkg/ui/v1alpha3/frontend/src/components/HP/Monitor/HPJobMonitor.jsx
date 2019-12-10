import React from 'react';
import { withStyles } from '@material-ui/core/styles';
 
import FilterPanel from './FilterPanel';
import HPJobList from './HPJobList';

import { fetchHPJobs } from '../../../actions/hpMonitorActions';
import { connect } from 'react-redux'


const styles = theme => ({
    root: {
        width: '90%',
        margin: '0 auto',
        marginTop: 10
    },
});

class HPJobMonitor extends React.Component {
    
    componentDidMount() {
        this.props.fetchHPJobs();
    }

    render () {
        const { classes } = this.props;

        return (
            <div className={classes.root}>
                <h1>Experiment Monitor</h1>
                <FilterPanel />
                <HPJobList />
            </div>
        )
    }
}


export default connect(null, { fetchHPJobs })(withStyles(styles)(HPJobMonitor));
