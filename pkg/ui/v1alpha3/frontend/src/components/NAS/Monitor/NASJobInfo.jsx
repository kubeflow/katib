import React from 'react'
import { connect } from 'react-redux'
import { withStyles } from '@material-ui/core';
import Button from '@material-ui/core/Button';
import { Link } from 'react-router-dom';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import Typography from '@material-ui/core/Typography';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import LinearProgress from '@material-ui/core/LinearProgress';

import { fetchNASJobInfo } from '../../../actions/nasMonitorActions';

import NASJobStepInfo from './NASJobStepInfo';

const module = "nasMonitor";

const styles = theme => ({
    root: {
        width: '90%',
        margin: '0 auto',
        paddingTop: 20,
    },
    heading: {
        fontSize: theme.typography.pxToRem(15),
        fontWeight: theme.typography.fontWeightRegular,
    },
    panel: {
        width: '100%',
    },
    header: {
        marginTop: 30,
        textAlign: "center"
    }
})


class NASJobInfo extends React.Component {

    componentDidMount() {
        this.props.fetchNASJobInfo(
            this.props.match.params.name, this.props.match.params.namespace);
    }

    render() {
        const { classes } = this.props;
        return (
            <div className={classes.root}>
                <Link to="/katib/nas_monitor">
                    <Button variant={"contained"} color={"primary"}>
                        Back
                    </Button>
                </Link>
                {this.props.loading ?
                    <LinearProgress color={"primary"} className={classes.loading} />
                    :
                    <div>
                        <Typography className={classes.header} variant={"h5"}>
                            Experiment Name: {this.props.match.params.name}
                        </Typography>
                        <Typography className={classes.header} variant={"h5"}>
                            Experiment Namespace: {this.props.match.params.namespace}
                        </Typography>
                        <br />
                        {this.props.steps.map((step, i) => {
                            return (
                                <ExpansionPanel key={i} className={classes.panel}>
                                    <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
                                        <Typography className={classes.heading}>{step.name}</Typography>
                                    </ExpansionPanelSummary>
                                    <ExpansionPanelDetails>
                                        <NASJobStepInfo step={step} id={i + 1} />
                                    </ExpansionPanelDetails>
                                </ExpansionPanel>
                            )
                        })}
                    </div>
                }

            </div>
        )
    }
}

const mapStateToProps = (state) => {
    return {
        steps: state[module].steps,
        loading: state[module].loading,
    }
}


export default connect(mapStateToProps, { fetchNASJobInfo })(withStyles(styles)(NASJobInfo));
