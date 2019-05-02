import React from 'react'
import { connect } from 'react-redux'
import { withStyles } from '@material-ui/core';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import { Link } from 'react-router-dom';
import LinearProgress from '@material-ui/core/LinearProgress';


import { fetchJobInfo } from '../../../actions/hpMonitorActions';

import HPPlot from './Plot';
import HPTable from './Table';
import PlotDialog from './Dialog';

const module = "hpMonitor";

const styles = theme => ({
    root: {
        width: '90%',
        margin: '0 auto',
        padding: 20,
    },
    loading: {
        marginTop: 30,
    }
})

class HPJobInfo extends React.Component {

    componentDidMount() {
        this.props.fetchJobInfo(this.props.match.params.id);
    }

    render () {
        const { classes } = this.props;
        return (
            <div className={classes.root}>
                <Link to="/katib/hp_monitor">
                    <Button variant={"contained"} color={"primary"}>
                        Back
                    </Button>
                </Link>
                {this.props.loading ? 
                <LinearProgress color={"primary"} className={classes.loading} />
                :
                <div>
                    <Typography variant={"h5"}>
                        Study ID: {this.props.match.params.id}
                    </Typography>
                    <br />
                    <HPPlot id={this.props.match.params.id} />
                    <HPTable id={this.props.match.params.id} />
                    <PlotDialog />
                </div>
                }
            </div>
        )
    }
}

const mapStateToProps = (state) => ({
  loading: state[module].loading,
})


export default connect(mapStateToProps, { fetchJobInfo })(withStyles(styles)(HPJobInfo));
