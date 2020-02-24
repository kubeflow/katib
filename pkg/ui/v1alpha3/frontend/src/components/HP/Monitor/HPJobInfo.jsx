import React from 'react'
import { connect } from 'react-redux'
import { withStyles } from '@material-ui/core';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import { Link } from 'react-router-dom';
import LinearProgress from '@material-ui/core/LinearProgress';

import { fetchHPJobInfo, fetchHPJob } from '../../../actions/hpMonitorActions';

import HPJobPlot from './HPJobPlot';
import HPJobTable from './HPJobTable';
import TrialInfoDialog from './TrialInfoDialog';
import ExperimentInfoDialog from './ExperimentInfoDialog';

const module = "hpMonitor";

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
        textAlign: "center",
        marginBottom: 15
    }
})

class HPJobInfo extends React.Component {

    componentDidMount() {
        this.props.fetchHPJobInfo(
            this.props.match.params.name, this.props.match.params.namespace);
    }

    fetchAndOpenDialogExperiment = (experimentName, experimentNamespace) => (event) => {
      this.props.fetchHPJob(experimentName, experimentNamespace)
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
                    <Typography  className = {classes.header} variant={"h5"}>
                        Experiment Name: {this.props.match.params.name}
                    </Typography>
                    <Typography  className = {classes.header} variant={"h5"}>
                        Experiment Namespace: {this.props.match.params.namespace}
                    </Typography>
                    <div className = {classes.header}>
                        <Button
                          variant={"contained"}
                          color={"primary"}
                          onClick={this.fetchAndOpenDialogExperiment(
                            this.props.match.params.name,
                            this.props.match.params.namespace)}
                        >
                                View Experiment
                        </Button>
                    </div>
                    <HPJobPlot name={this.props.match.params.name} />
                    <HPJobTable namespace={this.props.match.params.namespace} />
                    <ExperimentInfoDialog/>
                    <TrialInfoDialog />
                </div>
                }
            </div>
        )
    }
}

const mapStateToProps = (state) => ({
  loading: state[module].loading
})


export default connect(mapStateToProps, { fetchHPJobInfo, fetchHPJob })(withStyles(styles)(HPJobInfo));
