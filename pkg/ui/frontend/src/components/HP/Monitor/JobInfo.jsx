import React from 'react'
import { connect } from 'react-redux'
import { withStyles } from '@material-ui/core';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import { Link } from 'react-router-dom';

import { fetchJobInfo } from '../../../actions/hpMonitorActions';

import HPPlot from './Plot';
import HPTable from './Table';
import PlotDialog from './Dialog';


const styles = theme => ({
    root: {
        width: '90%',
        margin: '0 auto',
        padding: 20,
    }
})

class HPJobInfo extends React.Component {

    // componentDidMount() {
    //     this.props.fetchJobInfo(this.props.match.params.id);
    // }

    render () {
        const { classes } = this.props;
        return (
            <div className={classes.root}>
                <Link to="/hp_monitor">
                    <Button variant={"contained"} color={"primary"}>
                        Back
                    </Button>
                </Link>
                <Typography variant={"h5"}>
                    JOB INFO for {this.props.match.params.id}
                </Typography>
                <br />
                <HPPlot id={this.props.match.params.id} />
                <HPTable />
                <PlotDialog />
            </div>
        )
    }
}

const mapStateToProps = (state) => ({
  
})


export default connect(mapStateToProps, { fetchJobInfo })(withStyles(styles)(HPJobInfo));
